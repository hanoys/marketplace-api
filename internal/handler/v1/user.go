package v1

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hanoys/marketplace-api/internal/domain/dto"
	"net/http"
)

func (h *Handler) InitUserRoutes(router *gin.Engine) {
	userRoutes := router.Group("/user")
	{
		userRoutes.POST("/signup", h.signUpUser)
		userRoutes.POST("/login", h.logInUser)

		authorizedRoutes := userRoutes.Group("/", h.verifyToken)
		{
			authorizedRoutes.POST("/logout", h.logOutUser)
			authorizedRoutes.GET("/refresh", h.refreshToken)
		}
	}
}

func (h *Handler) signUpUser(c *gin.Context) {
	var userDTO dto.SignUpDTO

	if err := c.ShouldBindJSON(&userDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("bad json format: %v", err).Error()})
		return
	}

	user, err := h.services.Users.SignUp(context.TODO(), userDTO.Login, userDTO.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("can't create user: %v\n", err).Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
