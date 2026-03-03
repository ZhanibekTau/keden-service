package admin

import (
	"keden-service/back/internal/services/company"
	"keden-service/back/internal/services/document"
	"keden-service/back/internal/services/subscription"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	companyService *company.CompanyService
	subService     *subscription.SubscriptionService
	docService     *document.DocumentService
}

func NewAdminHandler(
	cs *company.CompanyService,
	ss *subscription.SubscriptionService,
	ds *document.DocumentService,
) *AdminHandler {
	return &AdminHandler{
		companyService: cs,
		subService:     ss,
		docService:     ds,
	}
}

func (h *AdminHandler) GetStats(c *gin.Context) {
	ctx := c.Request.Context()

	totalUsers, _ := h.companyService.GetStats(ctx)
	activeSubs, _ := h.subService.GetActiveCount(ctx)
	pendingSubs, _ := h.subService.GetPendingCount(ctx)
	totalDocs, completedDocs, _ := h.docService.GetStats(ctx)

	c.JSON(http.StatusOK, gin.H{
		"total_companies":       totalUsers,
		"active_subscriptions":  activeSubs,
		"pending_subscriptions": pendingSubs,
		"total_documents":       totalDocs,
		"completed_documents":   completedDocs,
	})
}
