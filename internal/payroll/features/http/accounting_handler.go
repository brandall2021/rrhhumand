package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/payroll/features/domain"
)

func (h *Handler) CreateMappingAccounting(c *gin.Context) {
	var req domain.AccountingAccountMapping
	if !bindJSON(c, &req) {
		return
	}
	m, err := h.AccountingSvc.CreateMapping(c.Request.Context(), companyID(c), userID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, m)
}

func (h *Handler) ListMappingsAccounting(c *gin.Context) {
	list, err := h.AccountingSvc.ListMappings(c.Request.Context(), companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, list)
}

func (h *Handler) UpdateMappingAccounting(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req domain.AccountingAccountMapping
	if !bindJSON(c, &req) {
		return
	}
	req.ID = id
	req.CompanyID = companyID(c)
	m, err := h.AccountingSvc.UpdateMapping(c.Request.Context(), companyID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, m)
}

func (h *Handler) ListExportsAccounting(c *gin.Context) {
	list, err := h.AccountingSvc.ListExports(c.Request.Context(), companyID(c), uuidPtr(qs(c, "run_id")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, list)
}

func (h *Handler) GetExportAccounting(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	e, err := h.AccountingSvc.GetExport(c.Request.Context(), companyID(c), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "export not found"})
		return
	}
	success(c, e)
}

type GenerateAccountingExportReq struct {
	RunID      string `json:"run_id" binding:"required"`
	ExportType string `json:"export_type" binding:"required"`
	FileFormat string `json:"file_format" binding:"required"`
}

func (h *Handler) GenerateExportAccounting(c *gin.Context) {
	var req GenerateAccountingExportReq
	if !bindJSON(c, &req) {
		return
	}
	runID, err := uuid.Parse(req.RunID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid run_id"})
		return
	}
	e, err := h.AccountingSvc.GenerateExport(c.Request.Context(), companyID(c), runID, req.ExportType, req.FileFormat, userID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, e)
}

func (h *Handler) GetEntriesAccounting(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	entries, err := h.AccountingSvc.GetEntries(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "entries not found"})
		return
	}
	success(c, entries)
}

func (h *Handler) DownloadExportAccounting(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	e, err := h.AccountingSvc.GetExport(c.Request.Context(), companyID(c), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "export not found"})
		return
	}
	content := ""
	if e.FileContent != nil {
		content = *e.FileContent
	}
	c.Data(http.StatusOK, "application/octet-stream", []byte(content))
}
