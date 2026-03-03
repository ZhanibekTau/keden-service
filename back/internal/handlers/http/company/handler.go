package company

import (
	"errors"
	"keden-service/back/internal/middleware"
	"keden-service/back/internal/services/company"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CompanyHandler struct {
	companyService *company.CompanyService
}

func NewCompanyHandler(cs *company.CompanyService) *CompanyHandler {
	return &CompanyHandler{companyService: cs}
}

func (h *CompanyHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	profile, err := h.companyService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, company.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *CompanyHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req company.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.companyService.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *CompanyHandler) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req company.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.companyService.ChangePassword(c.Request.Context(), userID, req)
	if err != nil {
		if errors.Is(err, company.ErrWrongPassword) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "current password is incorrect"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to change password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}

func (h *CompanyHandler) GetAllClients(c *gin.Context) {
	clients, err := h.companyService.GetAllClients(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get clients"})
		return
	}

	c.JSON(http.StatusOK, clients)
}

func (h *CompanyHandler) UpdateUserStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.companyService.UpdateUserStatus(c.Request.Context(), uint(id), req.IsActive); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user status updated"})
}
