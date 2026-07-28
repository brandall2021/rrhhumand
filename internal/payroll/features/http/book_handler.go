package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) ListEntries(c *gin.Context) {
	runID, err := uuid.Parse(c.Query("run_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid run_id"})
		return
	}
	list, err := h.BookSvc.ListEntries(c.Request.Context(), runID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, list)
}

func (h *Handler) GetEntry(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	e, err := h.BookSvc.GetEntry(c.Request.Context(), companyID(c), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
		return
	}
	success(c, e)
}

type GenerateBookEntriesReq struct {
	RunID string `json:"run_id" binding:"required"`
}

func (h *Handler) GenerateEntries(c *gin.Context) {
	var req GenerateBookEntriesReq
	if !bindJSON(c, &req) {
		return
	}
	runID, err := uuid.Parse(req.RunID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid run_id"})
		return
	}
	entries, err := h.BookSvc.GenerateEntries(c.Request.Context(), companyID(c), runID, userID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, entries)
}

func (h *Handler) ListExportsBook(c *gin.Context) {
	list, err := h.BookSvc.ListExports(c.Request.Context(), companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, list)
}

func (h *Handler) GetExportBook(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	e, err := h.BookSvc.GetExport(c.Request.Context(), companyID(c), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "export not found"})
		return
	}
	success(c, e)
}

type ExportBookReq struct {
	PeriodID string `json:"period_id" binding:"required"`
	Year     int    `json:"year" binding:"required"`
	Month    int    `json:"month" binding:"required"`
	Format   string `json:"format" binding:"required"`
}

func (h *Handler) ExportBook(c *gin.Context) {
	var req ExportBookReq
	if !bindJSON(c, &req) {
		return
	}
	periodID, err := uuid.Parse(req.PeriodID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period_id"})
		return
	}
	e, err := h.BookSvc.ExportBook(c.Request.Context(), companyID(c), periodID, req.Year, req.Month, req.Format, userID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, e)
}

func (h *Handler) GetBookSummary(c *gin.Context) {
	runID, err := uuid.Parse(c.Query("run_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid run_id"})
		return
	}
	s, err := h.BookSvc.GetBookSummary(c.Request.Context(), runID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, s)
}
