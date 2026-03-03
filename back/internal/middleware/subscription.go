package middleware

import (
	"keden-service/back/internal/services/subscription"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckActiveSubscription(subService *subscription.SubscriptionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		if userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		active, err := subService.CheckActiveSubscription(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check subscription"})
			c.Abort()
			return
		}

		if !active {
			c.JSON(http.StatusForbidden, gin.H{"error": "active subscription required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
