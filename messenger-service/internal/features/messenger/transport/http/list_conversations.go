package transport_http

import (
	core_logger "messenger-service/internal/core/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//	@Summary	Список моих тредов переписки
//	@Tags		messenger
//	@Security	BearerAuth
//	@Produce	json
//	@Success	200	{array}		ConversationResponse
//	@Failure	401	{object}	map[string]string
//	@Router		/conversations [get]
func (h *Handler) ListConversations(c *gin.Context) {
	log := core_logger.FromContext(c.Request.Context())

	userID, ok := c.MustGet("user_id").(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	conversations, err := h.service.ListConversations(c.Request.Context(), userID)
	if err != nil {
		log.Error("list conversations error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	items := make([]ConversationResponse, 0, len(conversations))
	for _, conv := range conversations {
		items = append(items, toConversationResponse(conv))
	}

	c.JSON(http.StatusOK, items)
}
