package transport_http

import (
	"errors"
	"net/http"
	core_errors "user-service/internal/core/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//	@Summary	Получить свой профиль
//	@Description	Возвращает профиль авторизованного пользователя (id берётся из токена через gateway)
//	@Tags		user
//	@Security	BearerAuth
//	@Produce	json
//	@Success	200	{object}	UserDTO
//	@Failure	401	{object}	map[string]string
//	@Failure	404	{object}	map[string]string
//	@Router		/me [get]
func (h *HTTPHandler) GetMyProfile(c *gin.Context) {
	id := c.MustGet("user_id").(uuid.UUID)

	user, err := h.service.GetUser(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get user",
		})
		return
	}

	c.JSON(http.StatusOK, UserDTO{
		ID:          user.ID,
		Version:     user.Version,
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
	})
}
