package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/benefits/domain"
	"github.com/shopspring/decimal"
)

func (h *Handler) CreateCost(c *gin.Context) {
	var req domain.BenefitCost
	if !bindJSON(c, &req) {
		return
	}
	cost, err := h.CostSvc.CreateCost(c.Request.Context(), companyID(c), userID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, cost)
}

func (h *Handler) ListCosts(c *gin.Context) {
	benefitID, err := uuid.Parse(c.Param("benefitId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid benefit id"})
		return
	}
	costs, err := h.CostSvc.ListCosts(c.Request.Context(), companyID(c), benefitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, costs)
}

func (h *Handler) UpdateCost(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req domain.BenefitCost
	if !bindJSON(c, &req) {
		return
	}
	req.ID = id
	cost, err := h.CostSvc.UpdateCost(c.Request.Context(), companyID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, cost)
}

func (h *Handler) CreateSchedule(c *gin.Context) {
	type reqBody struct {
		BenefitID uuid.UUID       `json:"benefit_id" binding:"required"`
		Date      time.Time       `json:"schedule_date" binding:"required"`
		Amount    decimal.Decimal `json:"amount" binding:"required"`
		Currency  string          `json:"currency"`
		Notes     *string         `json:"notes"`
	}
	var req reqBody
	if !bindJSON(c, &req) {
		return
	}
	if req.Currency == "" {
		req.Currency = "ARS"
	}
	schedule, err := h.CostSvc.CreateSchedule(c.Request.Context(), companyID(c), req.BenefitID, req.Date, req.Amount, req.Currency, req.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, schedule)
}

func (h *Handler) ListSchedules(c *gin.Context) {
	benefitID, err := uuid.Parse(c.Param("benefitId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid benefit id"})
		return
	}
	var from, to *time.Time
	if f := c.Query("from"); f != "" {
		t, err := time.Parse(time.RFC3339, f)
		if err == nil {
			from = &t
		}
	}
	if t := c.Query("to"); t != "" {
		parsed, err := time.Parse(time.RFC3339, t)
		if err == nil {
			to = &parsed
		}
	}
	schedules, err := h.CostSvc.ListSchedules(c.Request.Context(), companyID(c), benefitID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, schedules)
}

func (h *Handler) MarkSchedulePaid(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		PaymentReference string `json:"payment_reference" binding:"required"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if err := h.CostSvc.MarkSchedulePaid(c.Request.Context(), id, req.PaymentReference); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "schedule marked as paid"})
}
