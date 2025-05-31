package router

import (
	"LostAndFound/internal/delivery/http/handler"
	m "LostAndFound/internal/delivery/http/middleware"

	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"

	_ "LostAndFound/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(h *handler.Handler, redisClient *redis.Client) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.URLFormat)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/auth", func(r chi.Router) {
		r.With(m.RateLimitByUserID(redisClient, 3, 5*time.Minute)).Post("/register", h.Register)
		r.With(m.RateLimitByUserID(redisClient, 5, 1*time.Minute)).Post("/login", h.Login)
		r.Group(func(r chi.Router) {
			r.Use(m.AuthMiddleware(h.TokenManager))
			r.With(m.RateLimitByUserID(redisClient, 10, 1*time.Minute)).Post("/logout", h.Logout)
		})
	})

	r.Route("/users", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(m.AuthMiddleware(h.TokenManager))
			r.With(m.RateLimitByUserID(redisClient, 30, 1*time.Minute)).Get("/", h.GetProfile)
			r.With(m.RateLimitByUserID(redisClient, 3, 5*time.Minute)).Put("/update", h.UpdateProfile)
		})
		r.With(m.RateLimitByUserID(redisClient, 20, 1*time.Minute)).Get("/profile", h.GetProfileByID)
	})

	r.Route("/cards", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(m.AuthMiddleware(h.TokenManager))
			r.With(m.RateLimitByUserID(redisClient, 5, 1*time.Minute)).Post("/", h.CreateCard)
			r.With(m.RateLimitByUserID(redisClient, 5, 1*time.Minute)).Put("/{id}", h.UpdateCard)
			r.With(m.RateLimitByUserID(redisClient, 5, 1*time.Minute)).Delete("/{id}", h.DeleteCard)
		})

		r.With(m.RateLimitByUserID(redisClient, 60, 1*time.Minute)).Get("/all", h.GetAllCards)
		r.With(m.RateLimitByUserID(redisClient, 60, 1*time.Minute)).Get("/{id}", h.GetCardByID)
		r.With(m.RateLimitByUserID(redisClient, 60, 1*time.Minute)).Get("/near", h.GetCardsNear)
	})

	r.Route("/files", func(r chi.Router) {
		r.Use(m.AuthMiddleware(h.TokenManager))
		r.With(m.RateLimitByUserID(redisClient, 5, 10*time.Minute)).Post("/", h.UploadFile)
		r.With(m.RateLimitByUserID(redisClient, 5, 10*time.Minute)).Delete("/{key}", h.DeleteFile)
	})

	return r
}
