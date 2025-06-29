package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/dtg-lucifer/everato/internal/utils"
)

func TimeoutMiddleware(duration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), duration)
			defer cancel()

			// Replace original request with one that has a timeout context
			r = r.WithContext(ctx)

			wr := utils.NewHttpWriter(w, r) // HttpWriter

			// Channel to signal completion
			done := make(chan struct{})
			go func() {
				next.ServeHTTP(w, r)
				close(done)
			}()

			select {
			case <-ctx.Done():
				// Timeout exceeded
				if ctx.Err() == context.DeadlineExceeded {
					wr.Error(
						ctx.Err(),
						http.StatusRequestTimeout,
					)
					return
				}
			case <-done:
				// Handler finished before timeout
			}
		})
	}
}
