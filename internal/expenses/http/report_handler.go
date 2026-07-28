package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/rrhhumand/api/internal/expenses/domain"
)

func (h *Handler) CreateReport(c *gin.Context) {
	var req domain.ExpenseReport
	if !bindJSON(c, &req) {
		return
	}
	report, err := h.ReportSvc.CreateReport(c.Request.Context(), companyID(c), employeeID(c), userID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, report)
}

func (h *Handler) ListReports(c *gin.Context) {
	empID := qs(c, "employee_id")
	var eid *uuid.UUID
	if empID != nil {
		id, err := uuid.Parse(*empID)
		if err == nil {
			eid = &id
		}
	}
	status := qs(c, "status")
	reports, err := h.ReportSvc.ListReports(c.Request.Context(), companyID(c), eid, status, 0, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, reports)
}

func (h *Handler) GetReport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	report, err := h.ReportSvc.GetReport(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	success(c, report)
}

func (h *Handler) SubmitReport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.ReportSvc.SubmitReport(c.Request.Context(), id, userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	report, _ := h.ReportSvc.GetReport(c.Request.Context(), id)
	success(c, report)
}

func (h *Handler) ApproveReport(c *gin.Context) {
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
	if err := h.ReportSvc.ApproveReport(c.Request.Context(), id, userID(c), comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	report, _ := h.ReportSvc.GetReport(c.Request.Context(), id)
	success(c, report)
}

func (h *Handler) RejectReport(c *gin.Context) {
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
	if err := h.ReportSvc.RejectReport(c.Request.Context(), id, userID(c), req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	report, _ := h.ReportSvc.GetReport(c.Request.Context(), id)
	success(c, report)
}

func (h *Handler) ObserveReport(c *gin.Context) {
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
	if err := h.ReportSvc.ObserveReport(c.Request.Context(), id, userID(c), req.Observation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	report, _ := h.ReportSvc.GetReport(c.Request.Context(), id)
	success(c, report)
}
