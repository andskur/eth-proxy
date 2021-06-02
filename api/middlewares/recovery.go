package middlewares

import (
	"net/http"

	"github.com/dre1080/recovr"
	"github.com/justinas/alice"
)

// Recovery returns a new instance of Recovery middleware which traps panics
func Recovery() alice.Constructor {
	return func(next http.Handler) http.Handler {
		recovery := recovr.New()
		return recovery(next)
	}
}
