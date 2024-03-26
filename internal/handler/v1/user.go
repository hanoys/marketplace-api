package v1

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hanoys/marketplace-api/internal/domain/dto"
	"net/http"
)

func (h *Handler) createUser(c *gin.Context) {
	var userDTO dto.CreateUserDTO

	if err := c.ShouldBindJSON(&userDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "bad json format"})
		return
	}

	user, err := h.services.Users.Create(context.TODO(), userDTO.Login, userDTO.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("can't create user: %v\n", err).Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
