package transport_http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	core_errors "messenger-service/internal/core/errors"
	core_logger "messenger-service/internal/core/logger"
)

//	@Summary		Создать или получить тред переписки
//	@Description	Создаёт тред по объявлению с указанным собеседником, либо возвращает уже существующий (требует авторизации)
//	@Tags			messenger
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateConversationRequest	true	"Объявление и собеседник"
//	@Success		200		{object}	ConversationResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Router			/conversations [post]
func (h *Handler) CreateConversation(c *gin.Context) {
	log := core_logger.FromContext(c.Request.Context())
	
	userID, ok := c.MustGet("user_id").(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conv, err := h.service.CreateOrGetConversation(c.Request.Context(), req.ListingID, userID, req.RecipientID)
	if err != nil {
		switch {
		case errors.Is(err, core_errors.ErrInvalidArgument):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, core_errors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "listing not found"})
		default:
			log.Error("create conversation error")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, toConversationResponse(conv))
}
