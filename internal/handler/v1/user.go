package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hanoys/marketplace-api/internal/handler/dto"
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

	user, err := h.services.UsersService.SignUp(c.Request.Context(), userDTO.Login, userDTO.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("can't create user: %v\n", err).Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) logInUser(c *gin.Context) {
	var logInDTO dto.LogInDTO

	if err := c.ShouldBindJSON(&logInDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "bad json format"})
		return
	}

	tokenPair, err := h.services.Authorization.LogIn(c.Request.Context(), logInDTO.Login, logInDTO.Password)
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

	err := h.services.Authorization.LogOut(c.Request.Context(), token[0])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("logout error: %v", err).Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "successful logout"})
}

func (h *Handler) refreshToken(c *gin.Context) {
	token, ok := c.Request.Header["Authorization"]
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "no authorization header"})
		return
	}

	tokenPair, err := h.services.Authorization.RefreshToken(c.Request.Context(), token[0])
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("login error: %v", err).Error())
		return
	}

	c.JSON(http.StatusOK, tokenPair)
}
