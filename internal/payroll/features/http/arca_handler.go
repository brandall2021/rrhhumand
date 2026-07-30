package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/payroll/features/domain"
)

type ArcaMappingReq struct {
	ConceptID       string   `json:"concept_id" binding:"required"`
	ArcaConceptCode string   `json:"arca_concept_code" binding:"required"`
	ArcaConceptName *string  `json:"arca_concept_name,omitempty"`
	MappingType     string   `json:"mapping_type" binding:"required"`
	Percentage      *float64 `json:"percentage,omitempty"`
	IsTaxable       bool     `json:"is_taxable"`
	IsContributable bool     `json:"is_contributable"`
	Notes           *string  `json:"notes,omitempty"`
	EffectiveFrom   string   `json:"effective_from" binding:"required"`
	EffectiveTo     *string  `json:"effective_to,omitempty"`
	IsActive        bool     `json:"is_active"`
}

func (h *Handler) CreateMappingArca(c *gin.Context) {
	var req domain.ArcaConceptMapping
	if !bindJSON(c, &req) {
		return
	}
	cid := companyID(c)
	uid := userID(c)
	req.CompanyID = cid
	req.CreatedBy = uid
	m, err := h.ArcaSvc.CreateMapping(c.Request.Context(), cid, uid, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, m)
}

func (h *Handler) ListMappingsArca(c *gin.Context) {
	list, err := h.ArcaSvc.ListMappings(c.Request.Context(), companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, list)
}

func (h *Handler) UpdateMappingArca(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req domain.ArcaConceptMapping
	if !bindJSON(c, &req) {
		return
	}
	req.ID = id
	req.CompanyID = companyID(c)
	m, err := h.ArcaSvc.UpdateMapping(c.Request.Context(), companyID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, m)
}

func (h *Handler) DeleteMappingArca(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.ArcaSvc.DeleteMapping(c.Request.Context(), companyID(c), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "deleted"})
}

func (h *Handler) ListExportsArca(c *gin.Context) {
	list, err := h.ArcaSvc.ListExports(c.Request.Context(), companyID(c), uuidPtr(qs(c, "run_id")), qi(c, "limit", 10), qi(c, "offset", 0))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, list)
}

func (h *Handler) GetExportArca(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	e, err := h.ArcaSvc.GetExport(c.Request.Context(), companyID(c), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "export not found"})
		return
	}
	success(c, e)
}

type GenerateArcaExportReq struct {
	RunID      string `json:"run_id" binding:"required"`
	ExportType string `json:"export_type" binding:"required"`
}

func (h *Handler) GenerateExportArca(c *gin.Context) {
	var req GenerateArcaExportReq
	if !bindJSON(c, &req) {
		return
	}
	runID, err := uuid.Parse(req.RunID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid run_id"})
		return
	}
	e, err := h.ArcaSvc.GenerateExport(c.Request.Context(), companyID(c), runID, req.ExportType, userID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, e)
}

func (h *Handler) ValidateExportArca(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.ArcaSvc.ValidateExport(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "validation passed"})
}

func (h *Handler) DownloadExportArca(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	e, err := h.ArcaSvc.GetExport(c.Request.Context(), companyID(c), id)
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
