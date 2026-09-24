package main

import (
	"context"
	"log"
	"net/http"

	"kayam-be/internal/config"
	"kayam-be/internal/db"
	"kayam-be/internal/handler"
	"kayam-be/internal/migrations"
	"kayam-be/internal/repository"
	"kayam-be/internal/router"
	"kayam-be/internal/service"
	"kayam-be/internal/storage"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	// Self-migrate on boot: every deploy (including this one, against a
	// brand-new Supabase database) applies any pending schema changes
	// before serving traffic — see internal/migrations/embed.go for why.
	if err := migrations.Run(cfg.DatabaseURL); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	queries := repository.New(pool)

	presigner, err := storage.NewPresigner(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to init storage presigner: %v", err)
	}

	authService := &service.AuthService{Queries: queries}
	checkinService := &service.CheckinService{Pool: pool}
	redemptionService := &service.RedemptionService{Pool: pool}

	h := router.Handlers{
		Config:     &handler.ConfigHandler{Queries: queries},
		Auth:       &handler.AuthHandler{Service: authService},
		Challenge:  &handler.ChallengeHandler{Queries: queries},
		Checkin:    &handler.CheckinHandler{Service: checkinService},
		Wallet:     &handler.WalletHandler{Queries: queries},
		Redemption: &handler.RedemptionHandler{Service: redemptionService, Queries: queries},
		Insight:    &handler.InsightHandler{Queries: queries},
		Community:  &handler.CommunityHandler{Queries: queries},
		Upload:     &handler.UploadHandler{Presigner: presigner},
		Survey:     &handler.SurveyHandler{Queries: queries},
	}

	r := router.New(h, queries)

	log.Printf("kayam-be listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
