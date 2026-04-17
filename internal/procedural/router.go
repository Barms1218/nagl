package procedural

import (
	"crypto/ecdsa"
	"net/http"

	"github.com/Barms1218/nagl/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func Routes(p *ProceduralService, secretKey *ecdsa.PublicKey) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.JWTMiddleware(secretKey))
	r.Post("/adventurers", RequestAdventurer(p))

	return r
}
