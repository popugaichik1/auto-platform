package transport_http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)



//	@Summary	Health check
//	@Tags		user
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Router		/health [get]
func (h *HTTPHandler) HealthCheck (c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Health": "ok",
	})
}