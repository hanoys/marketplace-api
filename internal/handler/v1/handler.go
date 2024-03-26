package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/hanoys/marketplace-api/internal/service"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes(router *gin.Engine) {
	router.POST("/ad/post", h.postAd)
	router.POST("/user", h.createUser)
}
