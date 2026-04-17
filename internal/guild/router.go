package guild

import (
	"crypto/ecdsa"
	"net/http"

	"github.com/Barms1218/nagl/internal/auth"
	"github.com/go-chi/chi/v5"
)

func Routes(service *GuildService, secret *ecdsa.PrivateKey) http.Handler {
	r := chi.NewRouter()

	// Public
	r.Post("/register", RegisterGuild(service, secret))
	r.Post("/login", Login(service, secret))

	r.Group(func(r chi.Router) {
		r.Use(auth.JWTMiddleware(&secret.PublicKey))

	})

	return r

}
