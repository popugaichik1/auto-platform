package transport_http

import (
	"errors"
	"net/http"
	core_errors "user-service/internal/core/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//	@Summary	Получить профиль пользователя
//	@Tags		user
//	@Produce	json
//	@Param		id	path		string	true	"ID пользователя"
//	@Success	200	{object}	UserDTO
//	@Failure	400	{object}	map[string]string
//	@Failure	404	{object}	map[string]string
//	@Router		/{id} [get]
func (h *HTTPHandler) GetProfile(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

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
