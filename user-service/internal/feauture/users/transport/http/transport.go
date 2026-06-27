package transport_http

import (
	"context"
	core_domain "user-service/internal/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "user-service/docs" // swagger docs
)

type HTTPHandler struct {
	service Service
}

func NewHTTPHandler(
	userRepo Service,
) *HTTPHandler {
	return &HTTPHandler{
		service: userRepo,
	}
}

type Service interface {
	SaveUser(
		ctx context.Context,
		ID uuid.UUID,
		username string,
		phoneNumber string,
	) (error)

	GetUser(
		ctx context.Context,
		id uuid.UUID,
	) (core_domain.User, error)
}


func (h *HTTPHandler) InitRoutes() *gin.Engine {
	router := gin.Default()


	routes := router.Group("/api/user")
	{
		routes.GET("/health", h.HealthCheck)
		routes.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		routes.GET("/:id", h.GetProfile)
		//routes.POST("/save", h.SaveUser)
	}
	return router
}