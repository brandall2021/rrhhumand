package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/payroll/features/domain"
)

func (h *Handler) CreateTemplateReport(c *gin.Context) {
	var req domain.ReportTemplate
	if !bindJSON(c, &req) {
		return
	}
	t, err := h.ReportSvc.CreateTemplate(c.Request.Context(), companyID(c), userID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, t)
}

func (h *Handler) ListTemplatesReport(c *gin.Context) {
	list, err := h.ReportSvc.ListTemplates(c.Request.Context(), companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, list)
}

func (h *Handler) GetTemplateReport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	t, err := h.ReportSvc.GetTemplate(c.Request.Context(), companyID(c), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		return
	}
	success(c, t)
}

func (h *Handler) UpdateTemplateReport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req domain.ReportTemplate
	if !bindJSON(c, &req) {
		return
	}
	req.ID = id
	req.CompanyID = companyID(c)
	t, err := h.ReportSvc.UpdateTemplate(c.Request.Context(), companyID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, t)
}

func (h *Handler) DeleteTemplateReport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.ReportSvc.DeleteTemplate(c.Request.Context(), companyID(c), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "deleted"})
}

type GenerateReportReq struct {
	RunID      string `json:"run_id" binding:"required"`
	TemplateID string `json:"template_id" binding:"required"`
	Format     string `json:"format" binding:"required"`
}

func (h *Handler) GenerateReport(c *gin.Context) {
	var req GenerateReportReq
	if !bindJSON(c, &req) {
		return
	}
	runID, err := uuid.Parse(req.RunID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid run_id"})
		return
	}
	templateID, err := uuid.Parse(req.TemplateID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template_id"})
		return
	}
	e, err := h.ReportSvc.GenerateReport(c.Request.Context(), companyID(c), runID, templateID, req.Format, userID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, e)
}

func (h *Handler) ListExportsReport(c *gin.Context) {
	list, err := h.ReportSvc.ListReportExports(c.Request.Context(), companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, list)
}

func (h *Handler) GetExportReport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	e, err := h.ReportSvc.GetReportExport(c.Request.Context(), companyID(c), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "export not found"})
		return
	}
	success(c, e)
}
