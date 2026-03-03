package auth

import (
	"errors"
	"keden-service/back/internal/services/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *auth.AuthService
}

func NewAuthHandler(as *auth.AuthService) *AuthHandler {
	return &AuthHandler{authService: as}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req auth.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, auth.ErrEmailExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
			return
		}
		if errors.Is(err, auth.ErrBINExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "BIN already registered"})
			return
		}
		if errors.Is(err, auth.ErrCompanyFieldsRequired) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "company_name, legal_name and bin are required for company account type"})
			return
		}
		if errors.Is(err, auth.ErrCompanyCreateFailed) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create company record, please try again"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		if errors.Is(err, auth.ErrAccountInactive) {
			c.JSON(http.StatusForbidden, gin.H{"error": "account is inactive"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req auth.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidToken) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token refresh failed"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req auth.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_ = h.authService.Logout(c.Request.Context(), req.RefreshToken)
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
