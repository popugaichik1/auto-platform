package transport_http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	core_domain "user-service/internal/core/domain"

	"github.com/google/uuid"
)

// TestInitRoutes_StaticMeDoesNotConflictWithWildcardID проверяет, что
// статический /api/user/me и параметризованный /api/user/:id могут жить
// на одном роутере без паники при регистрации и резолвятся в разные
// хендлеры (gin приоритизирует статический сегмент над wildcard).
func TestInitRoutes_StaticMeDoesNotConflictWithWildcardID(t *testing.T) {
	meID := uuid.New()
	otherID := uuid.New()

	h := NewHTTPHandler(&fakeService{
		getUserFunc: func(ctx context.Context, id uuid.UUID) (core_domain.User, error) {
			return core_domain.User{ID: id}, nil
		},
	})

	router := h.InitRoutes(testLogger())

	w1 := httptest.NewRecorder()
	req1 := httptest.NewRequest(http.MethodGet, "/api/user/me", nil)
	req1.Header.Set("X-User-Id", meID.String())
	router.ServeHTTP(w1, req1)
	if w1.Code != http.StatusOK {
		t.Fatalf("GET /api/user/me status = %d, want %d, body=%s", w1.Code, http.StatusOK, w1.Body.String())
	}

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/api/user/"+otherID.String(), nil)
	router.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("GET /api/user/:id status = %d, want %d, body=%s", w2.Code, http.StatusOK, w2.Body.String())
	}

	w3 := httptest.NewRecorder()
	req3 := httptest.NewRequest(http.MethodGet, "/api/user/me", nil)
	router.ServeHTTP(w3, req3)
	if w3.Code != http.StatusUnauthorized {
		t.Fatalf("GET /api/user/me without X-User-Id status = %d, want %d", w3.Code, http.StatusUnauthorized)
	}
}
