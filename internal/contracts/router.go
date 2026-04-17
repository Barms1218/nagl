package contracts

import (
	"crypto/ecdsa"
	"github.com/Barms1218/nagl/internal/middleware"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Routes(service *ContractService, secret *ecdsa.PublicKey) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.JWTMiddleware(secret))

	return r
}
