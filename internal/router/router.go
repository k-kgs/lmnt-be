package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"kayam-be/internal/handler"
	kmw "kayam-be/internal/middleware"
	"kayam-be/internal/repository"
)

type Handlers struct {
	Config     *handler.ConfigHandler
	Auth       *handler.AuthHandler
	Challenge  *handler.ChallengeHandler
	Checkin    *handler.CheckinHandler
	Wallet     *handler.WalletHandler
	Redemption *handler.RedemptionHandler
	Insight    *handler.InsightHandler
	Community  *handler.CommunityHandler
	Upload     *handler.UploadHandler
}

func New(h Handlers, queries *repository.Queries) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: false,
	}))

	r.Get("/healthz", handler.Health)
	r.Get("/health", handler.Health) // Render's health check path is configured as /health

	r.Route("/api", func(r chi.Router) {
		// public
		r.Get("/config", h.Config.Get)
		r.Post("/auth/login", h.Auth.Login)
		r.Get("/challenges/{id}/leaderboard", h.Community.Leaderboard)

		// authenticated
		r.Group(func(r chi.Router) {
			r.Use(kmw.RequireAuth(queries))

			r.Post("/challenges/{id}/join", h.Challenge.Join)
			r.Get("/user-challenges", h.Challenge.MyChallenges)
			r.Get("/user-challenges/{id}/trend", h.Insight.Trend)
			r.Get("/user-challenges/{id}/adherence", h.Insight.Adherence)

			r.Post("/checkins", h.Checkin.Create)
			r.Post("/uploads/request", h.Upload.RequestURL)

			r.Get("/wallet", h.Wallet.Get)
			r.Post("/redemptions", h.Redemption.Create)
			r.Get("/redemptions", h.Redemption.MyRedemptions)
		})
	})

	return r
}
