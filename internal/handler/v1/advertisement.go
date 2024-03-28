package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hanoys/marketplace-api/internal/domain"
	"github.com/hanoys/marketplace-api/internal/handler/dto"
	"github.com/hanoys/marketplace-api/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type getAdvertisementQueryParams struct {
	pageString string
	sortString string
	dirString  string
	minString  string
	maxString  string
}

type getAdvertisementConvertedParams struct {
	page     int
	sort     domain.SortType
	dir      domain.DirectionType
	minPrice float64
	maxPrice float64
}

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

	ad, err := h.services.AdvertisementsService.Create(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("can't create post: %v\n", err).Error()})
		return
	}

	c.JSON(http.StatusOK, ad)
}

func (h *Handler) getAd(c *gin.Context) {
	convertedParams, err := convertParams(getAdvertisementQueryParams{
		pageString: c.Query("page"),
		sortString: c.Query("sort"),
		dirString:  c.Query("dir"),
		minString:  c.Query("min"),
		maxString:  c.Query("max"),
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("params error: %v\n", err).Error()})
		return
	}

	userID, ok := c.Get("userID")
	if _, assertOK := userID.(int); !ok || !assertOK {
		userID = -1
	}

	serviceParams := service.AdvertisementSortParams{
		UserID:     userID.(int),
		PageNumber: convertedParams.page,
		Sort:       convertedParams.sort,
		Dir:        convertedParams.dir,
		MinPrice:   convertedParams.minPrice,
		MaxPrice:   convertedParams.maxPrice,
	}

	advertisements, err := h.services.AdvertisementsService.GetAdvertisements(c.Request.Context(), serviceParams)
	if advertisements == nil {
		advertisements = make([]domain.AdvertisementEntry, 0)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("can't create post: %v\n", err).Error()})
		return
	}

	c.JSON(http.StatusOK, advertisements)
}

func convertParams(queryParams getAdvertisementQueryParams) (getAdvertisementConvertedParams, error) {
	var convertedParams getAdvertisementConvertedParams

	if queryParams.minString == "" {
		convertedParams.minPrice = 0
	} else {
		minPriceConverted, err := strconv.ParseFloat(queryParams.minString, 64)
		if err != nil || minPriceConverted < 0 {
			return getAdvertisementConvertedParams{}, fmt.Errorf("bad request params: min price")
		}

		convertedParams.minPrice = minPriceConverted
	}

	if queryParams.maxString == "" {
		convertedParams.maxPrice = 0
	} else {
		maxPriceConverted, err := strconv.ParseFloat(queryParams.maxString, 64)
		if err != nil || maxPriceConverted < 0 {
			return getAdvertisementConvertedParams{}, fmt.Errorf("bad request params: max price")
		}

		convertedParams.maxPrice = maxPriceConverted
	}

	if queryParams.pageString == "" {
		convertedParams.page = 1
	} else {
		pageConverted, err := strconv.Atoi(queryParams.pageString)
		if err != nil || pageConverted < 1 {
			return getAdvertisementConvertedParams{}, fmt.Errorf("bad request params: page")
		}

		convertedParams.page = pageConverted
	}

	if queryParams.sortString == "" {
		convertedParams.sort = domain.DefaultSort
	} else if strings.ToLower(queryParams.sortString) == "date" {
		convertedParams.sort = domain.DateSort
	} else if strings.ToLower(queryParams.sortString) == "price" {
		convertedParams.sort = domain.PriceSort
	} else {
		return getAdvertisementConvertedParams{}, fmt.Errorf("bad request params: sort")
	}

	if queryParams.dirString == "" || strings.ToLower(queryParams.dirString) == "asc" {
		convertedParams.dir = domain.AscDir
	} else if strings.ToLower(queryParams.dirString) == "desc" {
		convertedParams.dir = domain.DescDir
	} else {
		return getAdvertisementConvertedParams{}, fmt.Errorf("bad request params: dir")
	}

	return convertedParams, nil
}
