package auth_transport_http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TestValidateAuth_NoToken(t *testing.T) {
	h := NewAuthHTTPHandler(&fakeService{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/validate", nil)

	h.ValidateAuth(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("ValidateAuth() status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestValidateAuth_ValidToken(t *testing.T) {
	h := NewAuthHTTPHandler(&fakeService{
		validateTokenFunc: func(tokenString string) (jwt.MapClaims, error) {
			return jwt.MapClaims{"sub": "11111111-1111-1111-1111-111111111111"}, nil
		},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/validate", nil)
	c.Request.Header.Set("Authorization", "Bearer sometoken")

	h.ValidateAuth(c)

	if w.Code != http.StatusOK {
		t.Fatalf("ValidateAuth() status = %d, want %d", w.Code, http.StatusOK)
	}
	if got := w.Header().Get("X-User-Id"); got != "11111111-1111-1111-1111-111111111111" {
		t.Fatalf("X-User-Id header = %q", got)
	}
}

func TestValidateAuth_InvalidToken(t *testing.T) {
	h := NewAuthHTTPHandler(&fakeService{
		validateTokenFunc: func(tokenString string) (jwt.MapClaims, error) {
			return nil, errors.New("invalid token")
		},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/validate", nil)
	c.Request.Header.Set("Authorization", "Bearer bad")

	h.ValidateAuth(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("ValidateAuth() status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}
