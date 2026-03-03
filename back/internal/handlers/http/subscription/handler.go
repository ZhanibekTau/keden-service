package subscription

import (
	"errors"
	"keden-service/back/internal/middleware"
	"keden-service/back/internal/models"
	"keden-service/back/internal/services/subscription"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	subService *subscription.SubscriptionService
}

func NewSubscriptionHandler(ss *subscription.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{subService: ss}
}

func (h *SubscriptionHandler) RequestSubscription(c *gin.Context) {
	userID := middleware.GetUserID(c)

	sub, err := h.subService.RequestSubscription(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, subscription.ErrAlreadyPending) {
			c.JSON(http.StatusConflict, gin.H{"error": "you already have a pending subscription request"})
			return
		}
		if errors.Is(err, subscription.ErrAlreadyActive) {
			c.JSON(http.StatusConflict, gin.H{"error": "you already have an active subscription"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create subscription request"})
		return
	}

	c.JSON(http.StatusCreated, sub)
}

func (h *SubscriptionHandler) GetCurrentSubscription(c *gin.Context) {
	userID := middleware.GetUserID(c)

	sub, err := h.subService.GetCurrentSubscription(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get subscription"})
		return
	}

	if sub == nil {
		c.JSON(http.StatusOK, gin.H{"status": "none", "message": "no subscription"})
		return
	}

	c.JSON(http.StatusOK, sub)
}

func (h *SubscriptionHandler) GetSubscriptionHistory(c *gin.Context) {
	userID := middleware.GetUserID(c)

	subs, err := h.subService.GetSubscriptionHistory(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get history"})
		return
	}

	c.JSON(http.StatusOK, subs)
}

// GetActiveRequests returns all open subscription requests (pending + in_progress + invoice_sent)
// enriched with company data (BIN/IIN) for the admin panel.
func (h *SubscriptionHandler) GetActiveRequests(c *gin.Context) {
	details, err := h.subService.GetActiveRequests(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get subscription requests"})
		return
	}

	c.JSON(http.StatusOK, details)
}

// UpdateStatus moves a subscription through the workflow:
// pending → in_progress → invoice_sent → active (or any → rejected).
func (h *SubscriptionHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription id"})
		return
	}

	adminID := middleware.GetUserID(c)

	var req struct {
		Status  string `json:"status" binding:"required"`
		Comment string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status is required"})
		return
	}

	// Validate that the requested status is a known value
	known := map[string]bool{
		models.SubscriptionStatusInProgress:  true,
		models.SubscriptionStatusInvoiceSent: true,
		models.SubscriptionStatusActive:      true,
		models.SubscriptionStatusRejected:    true,
	}
	if !known[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown status"})
		return
	}

	if err := h.subService.UpdateStatus(c.Request.Context(), uint(id), adminID, req.Status, req.Comment); err != nil {
		if errors.Is(err, subscription.ErrSubscriptionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
			return
		}
		if errors.Is(err, subscription.ErrCannotTransition) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status transition"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}
