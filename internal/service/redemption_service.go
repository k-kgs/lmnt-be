package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"kayam-be/internal/repository"
)

var ErrInsufficientBalance = errors.New("insufficient balance")

type RedemptionService struct {
	Pool *pgxpool.Pool
}

type RedeemResult struct {
	RedemptionID pgtype.UUID
	CodeOrSlot   string
	NewBalance   int64
}

// Redeem locks the user's row (LockUserForRedemption), computes their current
// balance, and — only if it covers the item's cost — inserts the redemption
// and the debiting ledger entry, all in one transaction. The row lock is what
// makes two concurrent redemption requests for the same user serialize
// instead of racing past an insufficient-balance check (root design doc §4,
// kayam-be requirements.md US-5 AC3).
func (s *RedemptionService) Redeem(ctx context.Context, userID pgtype.UUID, itemID pgtype.UUID) (*RedeemResult, error) {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	q := repository.New(tx)

	if err := q.LockUserForRedemption(ctx, userID); err != nil {
		return nil, fmt.Errorf("lock user: %w", err)
	}

	item, err := q.GetRedemptionItem(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("redemption item not found: %w", err)
	}
	if !item.Active {
		return nil, fmt.Errorf("redemption item is not active")
	}

	balance, err := q.GetBalance(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get balance: %w", err)
	}
	if balance < int64(item.CoinCost) {
		return nil, ErrInsufficientBalance
	}

	codeOrSlot := generateCodeOrSlot(item)

	redemption, err := q.InsertRedemption(ctx, repository.InsertRedemptionParams{
		UserID:           userID,
		RedemptionItemID: itemID,
		CodeOrSlot:       &codeOrSlot,
	})
	if err != nil {
		return nil, fmt.Errorf("insert redemption: %w", err)
	}

	if _, err := q.InsertWalletTransaction(ctx, repository.InsertWalletTransactionParams{
		UserID: userID,
		Delta:  -item.CoinCost,
		Reason: "redemption",
	}); err != nil {
		return nil, fmt.Errorf("debit wallet: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return &RedeemResult{
		RedemptionID: redemption.ID,
		CodeOrSlot:   codeOrSlot,
		NewBalance:   balance - int64(item.CoinCost),
	}, nil
}

// generateCodeOrSlot: prototype-simple deterministic-ish code for vouchers;
// consultation "slots" are just a fixed placeholder for now (real slot
// picking is out of scope per root design doc — see requirements.md
// "Out of Scope: real calendar/video integration").
func generateCodeOrSlot(item repository.GetRedemptionItemRow) string {
	if item.Type == "voucher" {
		return "KYM-" + item.ID.String()[0:8]
	}
	return "slot-tbd"
}
