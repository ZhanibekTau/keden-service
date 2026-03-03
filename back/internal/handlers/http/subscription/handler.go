package subscription

import (
	"errors"
	"keden-service/back/internal/middleware"
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

func (h *SubscriptionHandler) GetPendingRequests(c *gin.Context) {
	subs, err := h.subService.GetPendingRequests(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get pending requests"})
		return
	}

	c.JSON(http.StatusOK, subs)
}

func (h *SubscriptionHandler) ApproveSubscription(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription id"})
		return
	}

	adminID := middleware.GetUserID(c)

	var req struct {
		Comment string `json:"comment"`
	}
	_ = c.ShouldBindJSON(&req)

	if err := h.subService.ApproveSubscription(c.Request.Context(), uint(id), adminID, req.Comment); err != nil {
		if errors.Is(err, subscription.ErrSubscriptionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
			return
		}
		if errors.Is(err, subscription.ErrCannotApprove) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "subscription is not in pending status"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to approve subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "subscription approved"})
}

func (h *SubscriptionHandler) RejectSubscription(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription id"})
		return
	}

	adminID := middleware.GetUserID(c)

	var req struct {
		Comment string `json:"comment"`
	}
	_ = c.ShouldBindJSON(&req)

	if err := h.subService.RejectSubscription(c.Request.Context(), uint(id), adminID, req.Comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reject subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "subscription rejected"})
}
