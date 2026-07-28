package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/rrhhumand/api/internal/expenses/domain"
)

func (h *Handler) RequestAdvance(c *gin.Context) {
	var req domain.ExpenseAdvance
	if !bindJSON(c, &req) {
		return
	}
	advance, err := h.AdvanceSvc.RequestAdvance(c.Request.Context(), companyID(c), employeeID(c), userID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, advance)
}

func (h *Handler) ListAdvances(c *gin.Context) {
	empID := qs(c, "employee_id")
	var eid *uuid.UUID
	if empID != nil {
		id, err := uuid.Parse(*empID)
		if err == nil {
			eid = &id
		}
	}
	status := qs(c, "status")
	advances, err := h.AdvanceSvc.ListAdvances(c.Request.Context(), companyID(c), eid, status, 0, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, advances)
}

func (h *Handler) GetAdvance(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	advance, err := h.AdvanceSvc.GetAdvance(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	success(c, advance)
}

func (h *Handler) ApproveAdvance(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		ApprovedAmount *decimal.Decimal `json:"approved_amount"`
	}
	if !bindJSON(c, &req) {
		return
	}
	approvedAmount := decimal.Zero
	if req.ApprovedAmount != nil {
		approvedAmount = *req.ApprovedAmount
	}
	if err := h.AdvanceSvc.ApproveAdvance(c.Request.Context(), id, userID(c), approvedAmount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	advance, _ := h.AdvanceSvc.GetAdvance(c.Request.Context(), id)
	success(c, advance)
}

func (h *Handler) PayAdvance(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.AdvanceSvc.PayAdvance(c.Request.Context(), id, userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	advance, _ := h.AdvanceSvc.GetAdvance(c.Request.Context(), id)
	success(c, advance)
}

func (h *Handler) RejectAdvance(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if err := h.AdvanceSvc.RejectAdvance(c.Request.Context(), id, userID(c), req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	advance, _ := h.AdvanceSvc.GetAdvance(c.Request.Context(), id)
	success(c, advance)
}

func (h *Handler) CancelAdvance(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.AdvanceSvc.CancelAdvance(c.Request.Context(), id, userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	advance, _ := h.AdvanceSvc.GetAdvance(c.Request.Context(), id)
	success(c, advance)
}
