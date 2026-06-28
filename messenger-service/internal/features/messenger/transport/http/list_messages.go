package transport_http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	core_errors "messenger-service/internal/core/errors"
	core_logger "messenger-service/internal/core/logger"
)

//	@Summary	История сообщений треда
//	@Tags		messenger
//	@Security	BearerAuth
//	@Produce	json
//	@Param		id		path		string	true	"ID треда"
//	@Param		page	query		int		false	"Страница"
//	@Param		limit	query		int		false	"Размер страницы"
//	@Success	200		{array}		MessageResponse
//	@Failure	400		{object}	map[string]string
//	@Failure	401		{object}	map[string]string
//	@Failure	403		{object}	map[string]string
//	@Failure	404		{object}	map[string]string
//	@Router		/conversations/{id}/messages [get]
func (h *Handler) ListMessages(c *gin.Context) {
	log := core_logger.FromContext(c.Request.Context())

	userID, ok := c.MustGet("user_id").(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	conversationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	messages, err := h.service.ListMessages(c.Request.Context(), conversationID, userID, page, limit)
	if err != nil {
		switch {
		case errors.Is(err, core_errors.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		case errors.Is(err, core_errors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
		default:
			log.Error("list message error")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	items := make([]MessageResponse, 0, len(messages))
	for _, msg := range messages {
		items = append(items, toMessageResponse(msg))
	}

	c.JSON(http.StatusOK, items)
}
