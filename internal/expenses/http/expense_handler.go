package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/rrhhumand/api/internal/expenses/domain"
)

func (h *Handler) CreateExpense(c *gin.Context) {
	var req domain.Expense
	if !bindJSON(c, &req) {
		return
	}
	expense, err := h.ExpenseSvc.CreateExpense(c.Request.Context(), companyID(c), employeeID(c), userID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, expense)
}

func (h *Handler) ListExpenses(c *gin.Context) {
	empID := qs(c, "employee_id")
	var eid *uuid.UUID
	if empID != nil {
		id, err := uuid.Parse(*empID)
		if err == nil {
			eid = &id
		}
	}
	travelID := qs(c, "travel_id")
	var tid *uuid.UUID
	if travelID != nil {
		id, err := uuid.Parse(*travelID)
		if err == nil {
			tid = &id
		}
	}
	status := qs(c, "status")
	expenses, err := h.ExpenseSvc.ListExpenses(c.Request.Context(), companyID(c), eid, tid, status, nil, nil, 0, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, expenses)
}

func (h *Handler) GetExpense(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	expense, err := h.ExpenseSvc.GetExpense(c.Request.Context(), companyID(c), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	success(c, expense)
}

func (h *Handler) UpdateExpense(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req domain.Expense
	if !bindJSON(c, &req) {
		return
	}
	req.ID = id
	expense, err := h.ExpenseSvc.UpdateExpense(c.Request.Context(), companyID(c), userID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, expense)
}

func (h *Handler) SubmitExpense(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.ExpenseSvc.SubmitExpense(c.Request.Context(), id, userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	expense, _ := h.ExpenseSvc.GetExpense(c.Request.Context(), companyID(c), id)
	success(c, expense)
}

func (h *Handler) ApproveExpense(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Comment *string `json:"comment"`
	}
	if !bindJSON(c, &req) {
		return
	}
	comment := ""
	if req.Comment != nil {
		comment = *req.Comment
	}
	if err := h.ExpenseSvc.ApproveExpense(c.Request.Context(), id, userID(c), comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	expense, _ := h.ExpenseSvc.GetExpense(c.Request.Context(), companyID(c), id)
	success(c, expense)
}

func (h *Handler) RejectExpense(c *gin.Context) {
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
	if err := h.ExpenseSvc.RejectExpense(c.Request.Context(), id, userID(c), req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	expense, _ := h.ExpenseSvc.GetExpense(c.Request.Context(), companyID(c), id)
	success(c, expense)
}

func (h *Handler) ObserveExpense(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Observation string `json:"observation" binding:"required"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if err := h.ExpenseSvc.ObserveExpense(c.Request.Context(), id, userID(c), req.Observation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	expense, _ := h.ExpenseSvc.GetExpense(c.Request.Context(), companyID(c), id)
	success(c, expense)
}

func (h *Handler) CancelExpense(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.ExpenseSvc.CancelExpense(c.Request.Context(), id, userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	expense, _ := h.ExpenseSvc.GetExpense(c.Request.Context(), companyID(c), id)
	success(c, expense)
}

func (h *Handler) DeleteExpense(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.ExpenseSvc.DeleteExpense(c.Request.Context(), companyID(c), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
