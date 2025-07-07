package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
)

func TimeoutMiddleware(time_str string) func(http.Handler) http.Handler {
	duration, err := time.ParseDuration(time_str)
	if err != nil {
		// If the duration string is invalid, default to 30 seconds
		logger := pkg.NewLogger()
		logger.Error("Invalid time string passed into timeout middleware, fallig back to", "timeout", "10s")
		duration = 10 * time.Second
	}

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
