package v1

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hanoys/marketplace-api/internal/domain/dto"
	"net/http"
)

func (h *Handler) postAd(c *gin.Context) {
	var postAdDTO dto.PostAdvertisementDTO

	if err := c.ShouldBindJSON(&postAdDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "bad json format"})
		return
	}

	ad, err := h.services.Advertisements.Create(context.TODO(), postAdDTO.UserID, postAdDTO.Title, postAdDTO.Body, postAdDTO.ImageURL, postAdDTO.Price)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("can't create post: %v\n", err).Error()})
		return
	}

	c.JSON(http.StatusOK, ad)
}
