package transport_http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	core_domain "user-service/internal/core/domain"
	core_errors "user-service/internal/core/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestGetMyProfile_OK(t *testing.T) {
	id := uuid.New()
	h := NewHTTPHandler(&fakeService{
		getUserFunc: func(ctx context.Context, gotID uuid.UUID) (core_domain.User, error) {
			if gotID != id {
				t.Fatalf("GetUser called with %v, want %v", gotID, id)
			}
			return core_domain.User{ID: id, Username: "ivan"}, nil
		},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/me", nil)
	c.Set("user_id", id)

	h.GetMyProfile(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetMyProfile_NotFound(t *testing.T) {
	h := NewHTTPHandler(&fakeService{
		getUserFunc: func(ctx context.Context, id uuid.UUID) (core_domain.User, error) {
			return core_domain.User{}, core_errors.ErrNotFound
		},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/me", nil)
	c.Set("user_id", uuid.New())

	h.GetMyProfile(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}
