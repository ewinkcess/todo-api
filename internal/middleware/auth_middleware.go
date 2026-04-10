package middleware

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"todo-api/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimiter(limit int, period time.Duration) gin.HandlerFunc {
	rate := limiter.Rate{
		Period: period,
		Limit:  int64(limit),
	}
	store := memory.NewStore()
	instance := limiter.New(store, rate)

	return func(c *gin.Context) {
		context, err := instance.Get(c.Request.Context(), c.ClientIP())
		if err != nil {
			c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(
				"Internal server error",
				"Gagal memproses rate limit",
			))
			c.Abort()
			return
		}
		if context.Reached {
			c.JSON(http.StatusTooManyRequests, domain.NewErrorResponse(
				"Terlalu banyak request",
				"Anda telah melebihi batas request, coba lagi beberapa saat",
			))
			c.Abort()
			return
		}
		c.Header("X-RateLimit-Limit", strconv.FormatInt(context.Limit, 10))
		c.Header("x-RateLimit-Remaining", strconv.FormatInt(context.Remaining, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(context.Reset, 10))

		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, domain.NewErrorResponse(
				"Akses ditolak",
				"Token tidak di temukan di header Authorization",
			))
			c.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, domain.NewErrorResponse(
				"Akses ditolak",
				"Format token tidak valid, gunakan format: Bearer <token>",
			))
			c.Abort()
			return
		}
		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, domain.NewErrorResponse(
				"Akses ditolak",
				"Token tidak valid atau sudah expired",
			))
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, domain.NewErrorResponse(
				"Akses ditolak",
				"Token tidak valid",
			))
			c.Abort()
			return
		}

		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Next()
	}

}
