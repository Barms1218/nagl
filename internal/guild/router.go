package guild

import (
	"crypto/ecdsa"
	"net/http"

	"github.com/Barms1218/nagl/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func Routes(service *GuildService, secret *ecdsa.PrivateKey) http.Handler {
	r := chi.NewRouter()

	// Public
	r.Post("/register", RegisterGuild(service))
	r.Post("/login", Login(service))

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTMiddleware(secret))

	})

	return r

}
