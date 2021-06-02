package middlewares

import (
	"net/http"

	"github.com/rs/cors"
)

// Cors returns a new instance of Cors middleware
// which providing cross-control origins rules
func Cors(next http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		MaxAge:           10,
		AllowCredentials: true,
	})
	return c.Handler(next)
}
