package v1

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hanoys/marketplace-api/internal/handler/dto"
	"net/http"
)

// TODO: move handler to user.go

func (h *Handler) logInUser(c *gin.Context) {
	var logInDTO dto.LogInDTO

	if err := c.ShouldBindJSON(&logInDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "bad json format"})
		return
	}

	tokenPair, err := h.services.Authorization.LogIn(context.TODO(), logInDTO.Login, logInDTO.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("login error: %v", err).Error())
		return
	}

	c.JSON(http.StatusOK, tokenPair)
}

func (h *Handler) logOutUser(c *gin.Context) {
	token, ok := c.Request.Header["Authorization"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "no authorization header"})
		return
	}

	err := h.services.Authorization.LogOut(context.TODO(), token[0])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("logout error: %v", err).Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "successful logout"})
}

func (h *Handler) verifyToken(c *gin.Context) {
	token, ok := c.Request.Header["Authorization"]
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "no authorization header"})
		return
	}

	userID, err := h.services.Authorization.VerifyToken(context.TODO(), token[0])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "incorrect token"})
		return
	}

	c.Set("userID", userID)
}

func (h *Handler) tryVerifyToken(c *gin.Context) {
	token, ok := c.Request.Header["Authorization"]
	if !ok {
		return
	}

	userID, err := h.services.Authorization.VerifyToken(context.TODO(), token[0])
	if err != nil {
		return
	}

	c.Set("userID", userID)
}

// RefreshToken TODO
func (h *Handler) refreshToken(c *gin.Context) {

}
