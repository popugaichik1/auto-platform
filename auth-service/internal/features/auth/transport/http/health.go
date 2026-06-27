package auth_transport_http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)



//	@Summary	Health check
//	@Tags		auth
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Router		/health [get]
func (h *AuthHTTPHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Health": "ok",
	})
}


//	@Summary	Проверка авторизации
//	@Description	Возвращает 200, если переданный access-токен валиден
//	@Tags		auth
//	@Security	BearerAuth
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Failure	401	{object}	map[string]string
//	@Router		/authorized [get]
func (h *AuthHTTPHandler) Authorized(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Authorized": "OK",
	})
}