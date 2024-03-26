package handler

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/hanoys/marketplace-api/internal/handler/v1"
	"github.com/hanoys/marketplace-api/internal/service"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{services: services}
}

func (r *Handler) Init() *gin.Engine {
	router := gin.New()
	handler := v1.NewHandler(r.services)
	handler.InitRoutes(router)
	return router
}
