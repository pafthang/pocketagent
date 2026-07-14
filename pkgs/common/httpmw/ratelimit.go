package httpmw

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pafthang/pocketagent/pkgs/common/secrets"
	"golang.org/x/time/rate"
)

// RateLimitConfig controls HTTP rate limiting (gate / auth endpoints).
type RateLimitConfig struct {
	Enabled       bool `mapstructure:"rate_limit_enabled"`
	PerMinute     int  `mapstructure:"rate_limit_per_minute"`
	Burst         int  `mapstructure:"rate_limit_burst"`
	AuthPerMinute int  `mapstructure:"auth_rate_limit_per_minute"`
	AuthBurst     int  `mapstructure:"auth_rate_limit_burst"`
}

// EffectiveEnabled returns true when rate limiting should be active.
func (c RateLimitConfig) EffectiveEnabled() bool {
	if c.Enabled {
		return true
	}
	return secrets.IsProduction()
}

func (c RateLimitConfig) apiBurst() int {
	if c.Burst > 0 {
		return c.Burst
	}
	return 40
}

func (c RateLimitConfig) authBurst() int {
	if c.AuthBurst > 0 {
		return c.AuthBurst
	}
	return 5
}

// RateLimiter returns a per-IP rate limiter middleware.
func RateLimiter(perMinute int, burst int) echo.MiddlewareFunc {
	if perMinute <= 0 {
		perMinute = 120
	}
	if burst <= 0 {
		burst = 40
	}
	return middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Skipper: healthMetricsSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
			Rate:      rate.Limit(float64(perMinute) / 60.0),
			Burst:     burst,
			ExpiresIn: middleware.DefaultRateLimiterMemoryStoreConfig.ExpiresIn,
		}),
		IdentifierExtractor: func(c echo.Context) (string, error) {
			return c.RealIP(), nil
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "rate limiter error"})
		},
		DenyHandler: func(c echo.Context, identifier string, err error) error {
			return c.JSON(http.StatusTooManyRequests, map[string]string{"error": "rate limit exceeded"})
		},
	})
}

// APIRateLimiter applies the default API rate limit from config.
func APIRateLimiter(cfg RateLimitConfig) echo.MiddlewareFunc {
	return RateLimiter(cfg.PerMinute, cfg.apiBurst())
}

// AuthRateLimiter applies a stricter per-IP limit for auth endpoints.
func AuthRateLimiter(cfg RateLimitConfig) echo.MiddlewareFunc {
	perMin := cfg.AuthPerMinute
	if perMin <= 0 {
		perMin = 10
	}
	return RateLimiter(perMin, cfg.authBurst())
}

func healthMetricsSkipper(c echo.Context) bool {
	path := c.Path()
	return path == "/health" || path == "/metrics"
}