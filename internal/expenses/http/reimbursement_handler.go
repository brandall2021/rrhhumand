package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/rrhhumand/api/internal/expenses/domain"
)

func (h *Handler) CreateReimbursement(c *gin.Context) {
	var req domain.ExpenseReimbursement
	if !bindJSON(c, &req) {
		return
	}
	reimb, err := h.ReimbursementSvc.CreateReimbursement(c.Request.Context(), companyID(c), employeeID(c), userID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, reimb)
}

func (h *Handler) ListReimbursements(c *gin.Context) {
	empID := qs(c, "employee_id")
	var eid *uuid.UUID
	if empID != nil {
		id, err := uuid.Parse(*empID)
		if err == nil {
			eid = &id
		}
	}
	status := qs(c, "status")
	reimbursements, err := h.ReimbursementSvc.ListReimbursements(c.Request.Context(), companyID(c), eid, status, 0, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, reimbursements)
}

func (h *Handler) GetReimbursement(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	reimb, err := h.ReimbursementSvc.GetReimbursement(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	success(c, reimb)
}

func (h *Handler) ApproveReimbursement(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.ReimbursementSvc.ApproveReimbursement(c.Request.Context(), id, userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	reimb, _ := h.ReimbursementSvc.GetReimbursement(c.Request.Context(), id)
	success(c, reimb)
}

func (h *Handler) PayReimbursement(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.ReimbursementSvc.PayReimbursement(c.Request.Context(), id, "", nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	reimb, _ := h.ReimbursementSvc.GetReimbursement(c.Request.Context(), id)
	success(c, reimb)
}

func (h *Handler) RejectReimbursement(c *gin.Context) {
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
	if err := h.ReimbursementSvc.RejectReimbursement(c.Request.Context(), id, userID(c), req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	reimb, _ := h.ReimbursementSvc.GetReimbursement(c.Request.Context(), id)
	success(c, reimb)
}

func (h *Handler) CancelReimbursement(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.ReimbursementSvc.CancelReimbursement(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	reimb, _ := h.ReimbursementSvc.GetReimbursement(c.Request.Context(), id)
	success(c, reimb)
}
