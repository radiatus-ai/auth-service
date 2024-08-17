package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/radiatus-ai/auth-service/internal/auth"
	"github.com/stretchr/testify/assert"
)

type mockAuthService struct{}

func (m *mockAuthService) Register(email, password string) error {
	return nil
}

func (m *mockAuthService) Login(email, password string) (string, error) {
	return "", nil
}

func (m *mockAuthService) VerifyToken(token string) (string, error) {
	if token == "valid_token" {
		return "user_123", nil
	}
	return "", auth.ErrInvalidToken
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &mockAuthService{}

	t.Run("Valid token", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.Use(AuthMiddleware(mockService))
		r.GET("/test", func(c *gin.Context) {
			userID, exists := c.Get("user_id")
			assert.True(t, exists)
			assert.Equal(t, "user_123", userID)
			c.Status(http.StatusOK)
		})

		c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)
		c.Request.Header.Set("Authorization", "Bearer valid_token")
		r.ServeHTTP(w, c.Request)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid token", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.Use(AuthMiddleware(mockService))
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)
		c.Request.Header.Set("Authorization", "Bearer invalid_token")
		r.ServeHTTP(w, c.Request)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Missing Authorization header", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.Use(AuthMiddleware(mockService))
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)
		r.ServeHTTP(w, c.Request)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
