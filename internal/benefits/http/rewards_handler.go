package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/benefits/domain"
)

func (h *Handler) CreateRewardsItem(c *gin.Context) {
	var req domain.TotalRewardsItem
	if !bindJSON(c, &req) {
		return
	}
	item, err := h.RewardsSvc.CreateRewardsItem(c.Request.Context(), companyID(c), userID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, item)
}

func (h *Handler) ListRewardsItems(c *gin.Context) {
	items, err := h.RewardsSvc.ListRewardsItems(c.Request.Context(), companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, items)
}

func (h *Handler) UpdateRewardsItem(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req domain.TotalRewardsItem
	if !bindJSON(c, &req) {
		return
	}
	req.ID = id
	item, err := h.RewardsSvc.UpdateRewardsItem(c.Request.Context(), companyID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, item)
}

func (h *Handler) GenerateSnapshot(c *gin.Context) {
	var req struct {
		EmployeeID uuid.UUID `json:"employee_id" binding:"required"`
		FiscalYear int       `json:"fiscal_year" binding:"required"`
		PeriodName string    `json:"period_name"`
	}
	if !bindJSON(c, &req) {
		return
	}
	snapshot, err := h.RewardsSvc.GenerateSnapshot(c.Request.Context(), companyID(c), req.EmployeeID, userID(c), req.FiscalYear, req.PeriodName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, snapshot)
}

func (h *Handler) GetLatestSnapshot(c *gin.Context) {
	empID, err := uuid.Parse(c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}
	snapshot, err := h.RewardsSvc.GetLatestSnapshot(c.Request.Context(), companyID(c), empID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	success(c, snapshot)
}

func (h *Handler) ListSnapshots(c *gin.Context) {
	var fiscalYear *int
	if fy := c.Query("fiscal_year"); fy != "" {
		if v, err := strconv.Atoi(fy); err == nil {
			fiscalYear = &v
		}
	}
	snapshots, err := h.RewardsSvc.ListSnapshots(c.Request.Context(), companyID(c), fiscalYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, snapshots)
}

func (h *Handler) CreateReportDefinition(c *gin.Context) {
	var req domain.BenefitReportDefinition
	if !bindJSON(c, &req) {
		return
	}
	d, err := h.RewardsSvc.CreateReportDefinition(c.Request.Context(), companyID(c), userID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, d)
}

func (h *Handler) ListReportDefinitions(c *gin.Context) {
	definitions, err := h.RewardsSvc.ListReportDefinitions(c.Request.Context(), companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, definitions)
}

func (h *Handler) LogNotification(c *gin.Context) {
	var req struct {
		EmployeeID *uuid.UUID        `json:"employee_id"`
		NotifType  string            `json:"notification_type" binding:"required"`
		Channel    string            `json:"channel" binding:"required"`
		Title      string            `json:"title" binding:"required"`
		Body       string            `json:"body"`
		Metadata   map[string]any    `json:"metadata"`
	}
	if !bindJSON(c, &req) {
		return
	}
	n, err := h.RewardsSvc.LogNotification(c.Request.Context(), companyID(c), req.EmployeeID, req.NotifType, req.Channel, req.Title, req.Body, req.Metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, n)
}

func (h *Handler) ListNotifications(c *gin.Context) {
	empID, err := uuid.Parse(c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}
	notifType := qs(c, "type")
	limit := qi(c, "limit", 50)
	offset := qi(c, "offset", 0)

	var notifTypePtr *string
	if notifType != nil && *notifType != "" {
		notifTypePtr = notifType
	}
	notifications, err := h.RewardsSvc.ListNotifications(c.Request.Context(), empID, notifTypePtr, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, notifications)
}

func (h *Handler) MarkNotificationRead(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.RewardsSvc.MarkNotificationRead(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "notification marked as read"})
}
