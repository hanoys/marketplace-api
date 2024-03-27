package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
