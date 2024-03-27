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
	"strings"
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

func convertParams(pageString string, sortString string, dirString string, minString string, maxString string) (int, domain.SortType, domain.DirectionType, float64, float64, error) {
	var page int
	var sort domain.SortType
	var dir domain.DirectionType
	var minPrice, maxPrice float64

	if minString == "" {
		minPrice = 0
	} else {
		minPriceConverted, err := strconv.ParseFloat(minString, 64)
		if err != nil || minPriceConverted < 0 {
			return 0, 0, 0, 0, 0, fmt.Errorf("bad request params: min price")
		}

		minPrice = minPriceConverted
	}

	if maxString == "" {
		maxPrice = 0
	} else {
		maxPriceConverted, err := strconv.ParseFloat(maxString, 64)
		if err != nil || maxPriceConverted < 0 {
			return 0, 0, 0, 0, 0, fmt.Errorf("bad request params: max price")
		}

		maxPrice = maxPriceConverted
	}

	if pageString == "" {
		page = 1
	} else {
		pageConverted, err := strconv.Atoi(pageString)
		if err != nil || pageConverted < 1 {
			return 0, 0, 0, 0, 0, fmt.Errorf("bad request params: page")
		}

		page = pageConverted
	}

	if sortString == "" {
		sort = domain.DefaultSort
	} else if strings.ToLower(sortString) == "date" {
		sort = domain.DateSort
	} else if strings.ToLower(sortString) == "price" {
		sort = domain.PriceSort
	} else {
		return 0, 0, 0, 0, 0, fmt.Errorf("bad request params: sort")
	}

	if dirString == "" || strings.ToLower(dirString) == "asc" {
		dir = domain.AscDir
	} else if strings.ToLower(dirString) == "desc" {
		dir = domain.DescDir
	} else {
		return 0, 0, 0, 0, 0, fmt.Errorf("bad request params: dir")
	}

	return page, sort, dir, minPrice, maxPrice, nil
}

func (h *Handler) getAd(c *gin.Context) {
	page, sort, dir, minPrice, maxPrice, err := convertParams(c.Query("page"), c.Query("sort"),
		c.Query("dir"), c.Query("min"), c.Query("max"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("params error: %v\n", err).Error()})
		return
	}

	userID, ok := c.Get("userID")
	if _, assertOK := userID.(int); !ok || !assertOK {
		userID = -1
	}

	params := service.AdvertisementSortParams{
		UserID:     userID.(int),
		PageNumber: page,
		Sort:       sort,
		Dir:        dir,
		MinPrice:   minPrice,
		MaxPrice:   maxPrice,
	}

	advertisements, err := h.services.AdvertisementsService.GetAdvertisements(context.TODO(), params)
	if advertisements == nil {
		advertisements = make([]domain.AdvertisementEntry, 0)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("can't create post: %v\n", err).Error()})
		return
	}

	c.JSON(http.StatusOK, advertisements)
}
