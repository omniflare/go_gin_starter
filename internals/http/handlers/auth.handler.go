package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	middleware "github.com/omniflare/go_starter/internals/http/middlewares"
	types "github.com/omniflare/go_starter/internals/http/types/auth"
	service "github.com/omniflare/go_starter/internals/services/auth"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(router *gin.RouterGroup, service service.AuthService) {
	handler := &AuthHandler{
		service: service,
	}

	// Public routes
	router.GET("/get", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
	router.POST("/login", handler.Login)
	router.POST("/register", handler.Register)
	router.POST("/forgot-password", handler.ForgotPassword)

	// Protected routes - using middleware
	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/verify-email", handler.VerifyEmail)
		protected.POST("/reset-password", handler.ResetPassword)
		protected.GET("/me", handler.GetMe)
		protected.POST("/resend-verification", handler.ResendVerification)
		protected.POST("/logout", handler.Logout)
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var loginData types.AuthLogin
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	_, user, err := h.service.Login(c, &loginData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Successfully logged in",
		"data": gin.H{
			"user": user,
		},
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var registerData types.AuthRegister
	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	token, user, err := h.service.Register(c, &registerData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Successfully registered",
		"data": gin.H{
			"token": token,
			"user":  user,
		},
	})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var request struct {
		Code string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	if err := h.service.VerifyEmail(c, request.Code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Email successfully verified",
	})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	if err := h.service.ForgotPassword(c, request.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Password reset email sent",
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var request struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	if err := h.service.ResetPassword(c, request.Token, request.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Password successfully reset",
	})
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "fail",
			"error":  "unauthorized",
		})
		return
	}

	user, err := h.service.GetMe(c, userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user": user,
		},
	})
}

func (h *AuthHandler) ResendVerification(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	if err := h.service.ResendVerification(c, request.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Verification email resent",
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	if err := h.service.Logout(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Successfully logged out",
	})
}
