package middlewares

import (
	"net/http"
	"time"

	"github.com/codeshelldev/secured-signal-api/internals/config"
	. "github.com/codeshelldev/secured-signal-api/internals/proxy/common"
	"golang.org/x/time/rate"
)

var RateLimit Middleware = Middleware{
	Name: "Rate Limiting",
	Use: ratelimitHandler,
}

type TokenLimiter struct {
	limiter *rate.Limiter
}

func NewTokenLimiter(limit int, period time.Duration) *TokenLimiter {
	r := rate.Every(period / time.Duration(limit))

	return &TokenLimiter{
		limiter: rate.NewLimiter(r, limit),
	}
}

func (t *TokenLimiter) Allow() bool {
	return t.limiter.Allow()
}

var tokenLimiters = map[string]*TokenLimiter{}

func ratelimitHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger := GetLogger(req)

		trusted := GetContext[bool](req, TrustedClientKey)

		if trusted {
			next.ServeHTTP(w, req)
			return
		}

		conf := GetConfigByReq(req)

		rateLimiting := conf.SETTINGS.ACCESS.RATE_LIMITING.OptOrEmpty(config.DEFAULT.SETTINGS.ACCESS.RATE_LIMITING)

		if rateLimiting.Period.Duration != 0 && rateLimiting.Limit != 0 {
			token := GetToken(req)

			tokenLimiter, exists := tokenLimiters[token]

			if !exists {
				tokenLimiter = NewTokenLimiter(rateLimiting.Limit, time.Duration(rateLimiting.Period.Duration))
				tokenLimiters[token] = tokenLimiter
			}

			if !tokenLimiter.Allow() {
				logger.Warn("Token exceeded Rate Limit")

				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				w.Header().Set("Retry-After", "60")

				return
			}
		}

		next.ServeHTTP(w, req)
	})
}