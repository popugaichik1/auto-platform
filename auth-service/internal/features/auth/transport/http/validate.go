package auth_transport_http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ValidateAuth — цель nginx-аннотации auth-url на "защищённом" Ingress
// (см. helm/.../templates/ingress.yaml). Какие пути считаются защищёнными
// решает сам Ingress через раздельные path-правила: сюда долетают только
// запросы к уже отобранным protected-путям, поэтому здесь только проверка
// Bearer-токена — без какой-либо публичный/защищённый логики.
//
//	@Summary		Проверка авторизации для API Gateway (nginx auth_request)
//	@Tags			auth
//	@Produce		json
//	@Success		200
//	@Failure		401	{object}	map[string]string
//	@Router			/validate [get]
func (h *AuthHTTPHandler) ValidateAuth(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
		return
	}

	claims, err := h.authService.ValidateToken(parts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}

	c.Header("X-User-Id", userID)
	c.Status(http.StatusOK)
}
