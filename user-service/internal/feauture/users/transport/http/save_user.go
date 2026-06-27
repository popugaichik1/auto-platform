package transport_http

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// type UserDTOResponse UserDTO


// func (h *HTTPHandler) SaveUser(c *gin.Context) {
// 	op := "User.Transport.SaveUser"
// 	var req SaveUserRequest

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			op: err.Error(),
// 		})
// 		return
// 	}

// 	saved, err := h.service.SaveUser(
// 		c.Request.Context(),
// 		req.ID,
// 		req.Username,
// 		req.PhoneNumber,
// 	)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			op: err.Error(),
// 		})
// 	}
// 	resp := UserDTOResponse{
// 		ID: 			saved.ID,
// 		Version: 		saved.Version,
// 		Username: 		saved.Username,
// 		PhoneNumber: 	saved.PhoneNumber,
// 	}

// 	c.JSON(http.StatusOK, resp)
// }