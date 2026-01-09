package middlewares

import (
	"net/http"
	"strings"
	"time"

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
		logger := getLogger(req)

		trusted := getContext[bool](req, trustedClientKey)

		if trusted {
			next.ServeHTTP(w, req)
			return
		}

		conf := getConfigByReq(req)

		rateLimiting := conf.SETTINGS.ACCESS.RATE_LIMITING

		limit := rateLimiting.Limit

		if limit == 0 {
			limit = getConfig("").SETTINGS.ACCESS.RATE_LIMITING.Limit
		}

		periodStr := rateLimiting.Period

		if strings.TrimSpace(periodStr) == "" {
			periodStr = conf.SETTINGS.ACCESS.RATE_LIMITING.Period
		}

		if strings.TrimSpace(periodStr) != "" && limit != 0 {
			period, err := time.ParseDuration(periodStr)

			if err != nil {
				logger.Error("Could not parse Duration: ", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			token := getToken(req)

			tokenLimiter, exists := tokenLimiters[token]

			if !exists {
				tokenLimiter = NewTokenLimiter(limit, period)
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