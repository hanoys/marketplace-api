package v1

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hanoys/marketplace-api/internal/domain"
	"github.com/hanoys/marketplace-api/internal/handler/dto"
	"github.com/hanoys/marketplace-api/internal/service"
	"net/http"
	"strconv"
)

func (h *Handler) InitAdvertisementRoutes(router *gin.Engine) {
	adGroup := router.Group("/ad")
	{
		adGroup.POST("/", h.verifyToken, h.postAd)
		adGroup.GET("/", h.tryVerifyToken, h.getAd)
	}
}

func (h *Handler) postAd(c *gin.Context) {
	var postAdDTO dto.PostAdvertisementDTO

	if err := c.ShouldBindJSON(&postAdDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "bad json format"})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "no user id"})
		return
	}

	params := service.AdvertisementCreateParams{
		UserID:   userID.(int),
		Title:    postAdDTO.Title,
		Body:     postAdDTO.Body,
		ImageURL: postAdDTO.ImageURL,
		Price:    postAdDTO.Price,
	}

	ad, err := h.services.AdvertisementsService.Create(context.TODO(), params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("can't create post: %v\n", err).Error()})
		return
	}

	c.JSON(http.StatusOK, ad)
}

func convertParams(pageString string, sortString string, dirString string) (int, domain.SortType, domain.DirectionType, error) {
	var page int
	var sort domain.SortType
	var dir domain.DirectionType

	if pageString == "" {
		page = 1
	} else {
		pageConverted, err := strconv.Atoi(pageString)
		if err != nil || pageConverted < 1 {
			return 0, 0, 0, fmt.Errorf("bad request params")
		}

		page = pageConverted
	}

	if sortString == "" {
		sort = domain.DefaultSort
	} else if sortString == "date" {
		sort = domain.DateSort
	} else if sortString == "price" {
		sort = domain.PriceSort
	} else {
		return 0, 0, 0, fmt.Errorf("bad request params")
	}

	if dirString == "" || dirString == "asc" {
		dir = domain.AscDir
	} else if dirString == "desc" {
		dir = domain.DescDir
	} else {
		return 0, 0, 0, fmt.Errorf("bad request params")
	}

	return page, sort, dir, nil
}

func (h *Handler) getAd(c *gin.Context) {
	page, sort, dir, _ := convertParams(c.Query("page"), c.Query("sort"), c.Query("dir"))

	userID, ok := c.Get("userID")
	if _, assertOK := userID.(int); !ok || !assertOK {
		userID = -1
	}

	params := service.AdvertisementSortParams{
		UserID:     userID.(int),
		PageNumber: page,
		Sort:       sort,
		Dir:        dir,
	}

	advertisements, err := h.services.AdvertisementsService.GetAdvertisements(context.TODO(), params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("can't create post: %v\n", err).Error()})
		return
	}

	c.JSON(http.StatusOK, advertisements)
}
