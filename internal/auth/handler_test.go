package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	registerFunc    func(email, password string) error
	loginFunc       func(email, password string) (string, error)
	verifyTokenFunc func(token string) (string, error)
}

func (m *mockService) Register(email, password string) error {
	return m.registerFunc(email, password)
}

func (m *mockService) Login(email, password string) (string, error) {
	return m.loginFunc(email, password)
}

func (m *mockService) VerifyToken(token string) (string, error) {
	return m.verifyTokenFunc(token)
}

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockService{
		registerFunc: func(email, password string) error {
			if email == "existing@example.com" {
				return ErrUserAlreadyExists
			}
			return nil
		},
	}

	handler := NewHandler(mockSvc)

	t.Run("Successful registration", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := bytes.NewBufferString(`{"email":"new@example.com","password":"password123"}`)
		c.Request, _ = http.NewRequest("POST", "/register", body)
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Register(c)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("User already exists", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := bytes.NewBufferString(`{"email":"existing@example.com","password":"password123"}`)
		c.Request, _ = http.NewRequest("POST", "/register", body)
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Register(c)

		assert.Equal(t, http.StatusConflict, w.Code)
	})
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockService{
		loginFunc: func(email, password string) (string, error) {
			if email == "valid@example.com" && password == "password123" {
				return "valid-token", nil
			}
			return "", ErrInvalidCredentials
		},
	}

	handler := NewHandler(mockSvc)

	t.Run("Successful login", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := bytes.NewBufferString(`{"email":"valid@example.com","password":"password123"}`)
		c.Request, _ = http.NewRequest("POST", "/login", body)
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Login(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "valid-token", response["token"])
	})

	t.Run("Invalid credentials", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := bytes.NewBufferString(`{"email":"invalid@example.com","password":"wrongpassword"}`)
		c.Request, _ = http.NewRequest("POST", "/login", body)
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Login(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestVerifyToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockService{
		verifyTokenFunc: func(token string) (string, error) {
			if token == "valid-token" {
				return "user-123", nil
			}
			return "", ErrInvalidCredentials
		},
	}

	handler := NewHandler(mockSvc)

	t.Run("Valid token", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := bytes.NewBufferString(`{"token":"valid-token"}`)
		c.Request, _ = http.NewRequest("POST", "/verify", body)
		c.Request.Header.Set("Content-Type", "application/json")

		handler.VerifyToken(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "user-123", response["user_id"])
	})

	t.Run("Invalid token", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := bytes.NewBufferString(`{"token":"invalid-token"}`)
		c.Request, _ = http.NewRequest("POST", "/verify", body)
		c.Request.Header.Set("Content-Type", "application/json")

		handler.VerifyToken(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
