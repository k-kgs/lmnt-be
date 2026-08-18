package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"kayam-be/internal/repository"
)

var ErrValidation = errors.New("validation error")

// ErrAlreadyCheckedIn: the checkins table's (user_challenge_id, date) unique
// constraint is what actually enforces "one check-in per day" — this just
// gives that constraint violation a clean, client-facing error instead of
// leaking a raw Postgres error through as a 500.
var ErrAlreadyCheckedIn = errors.New("already checked in today")

// ErrNotActiveParticipant: the participant has left or been disqualified —
// their check-ins stopped counting the moment that happened.
var ErrNotActiveParticipant = errors.New("not an active participant in this challenge")

type CheckinService struct {
	Pool *pgxpool.Pool
}

type CreateCheckinResult struct {
	CheckinID     pgtype.UUID
	CurrentStreak int32
	LongestStreak int32
	CoinsEarned   int32
	Disqualified  bool // true if this check-in's gap crossed the challenge's threshold
}

// CreateCheckin validates metricData against the vertical's schema, then
// writes the checkin, recomputes the streak, and credits the coin ledger —
// all in one DB transaction, so a failure at any step leaves zero partial
// writes (root design doc §3, US-3 AC3 in kayam-be requirements.md).
func (s *CheckinService) CreateCheckin(ctx context.Context, userChallengeID pgtype.UUID, metricData map[string]any) (*CreateCheckinResult, error) {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) // no-op if committed

	q := repository.New(tx)

	uc, err := q.GetUserChallengeWithVertical(ctx, userChallengeID)
	if err != nil {
		return nil, fmt.Errorf("user_challenge not found: %w", err)
	}

	if uc.Status != "active" {
		return nil, ErrNotActiveParticipant
	}

	if err := ValidateAgainstSchema(uc.InputSchema, metricData); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrValidation, err.Error())
	}

	metricJSON, err := json.Marshal(metricData)
	if err != nil {
		return nil, fmt.Errorf("marshal metric_data: %w", err)
	}

	checkin, err := q.InsertCheckin(ctx, repository.InsertCheckinParams{
		UserChallengeID: userChallengeID,
		MetricData:      metricJSON,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrAlreadyCheckedIn
		}
		return nil, fmt.Errorf("insert checkin: %w", err)
	}

	current, longest, missedDays, err := recomputeStreak(ctx, q, userChallengeID, checkin.Date)
	if err != nil {
		return nil, fmt.Errorf("recompute streak: %w", err)
	}

	// Auto-disqualification: only for challenges that declared a threshold
	// (nil = never, e.g. Weight/Diet — see seed/seed.sql for why those stay
	// non-punitive). This check-in itself still gets recorded and paid out —
	// the participant did the activity, verified — but their status flips so
	// the ErrNotActiveParticipant guard above blocks any further check-ins.
	disqualified := false
	if uc.DisqualifyAfterMissedDays != nil && missedDays >= *uc.DisqualifyAfterMissedDays {
		if err := q.DisqualifyUserChallenge(ctx, userChallengeID); err != nil {
			return nil, fmt.Errorf("disqualify: %w", err)
		}
		disqualified = true
	}

	coins := computeCoinsForCheckin(current)
	if _, err := q.InsertWalletTransaction(ctx, repository.InsertWalletTransactionParams{
		UserID: uc.UserID,
		Delta:  coins,
		Reason: "checkin",
	}); err != nil {
		return nil, fmt.Errorf("credit wallet: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return &CreateCheckinResult{
		CheckinID:     checkin.ID,
		CurrentStreak: current,
		LongestStreak: longest,
		CoinsEarned:   coins,
		Disqualified:  disqualified,
	}, nil
}

// recomputeStreak also returns missedDays — the number of full calendar days
// between the last check-in and this one, 0 for a first-ever or
// same-streak-continuing check-in. This is what the auto-disqualification
// check above compares against a challenge's threshold.
func recomputeStreak(ctx context.Context, q *repository.Queries, userChallengeID pgtype.UUID, checkinDate pgtype.Date) (current, longest, missedDays int32, err error) {
	existing, err := q.GetStreak(ctx, userChallengeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			current, longest = 1, 1
		} else {
			return 0, 0, 0, err
		}
	} else {
		yesterday := checkinDate.Time.AddDate(0, 0, -1)
		if existing.LastCheckinDate.Valid && existing.LastCheckinDate.Time.Equal(yesterday) {
			current = existing.CurrentStreak + 1
		} else {
			current = 1 // gap in check-ins — streak resets, not punished further than that
			if existing.LastCheckinDate.Valid {
				gapDays := int32(checkinDate.Time.Sub(existing.LastCheckinDate.Time).Hours() / 24)
				missedDays = gapDays - 1 // days strictly between last check-in and this one
			}
		}
		longest = existing.LongestStreak
		if current > longest {
			longest = current
		}
	}

	_, err = q.UpsertStreak(ctx, repository.UpsertStreakParams{
		UserChallengeID: userChallengeID,
		CurrentStreak:   current,
		LongestStreak:   longest,
		LastCheckinDate: checkinDate,
	})
	return current, longest, missedDays, err
}

// computeCoinsForCheckin: flat rate per check-in, plus a milestone bonus
// every 7-day streak. A deliberately simple prototype rule — see kayam-be
// design.md §3 for where this would evolve if reward tuning becomes a
// real product question.
func computeCoinsForCheckin(currentStreak int32) int32 {
	coins := int32(50)
	if currentStreak%7 == 0 {
		coins += 100
	}
	return coins
}
