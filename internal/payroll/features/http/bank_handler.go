package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/payroll/features/domain"
)

func (h *Handler) CreateBatch(c *gin.Context) {
	var req domain.BankBatch
	if !bindJSON(c, &req) {
		return
	}
	b, err := h.BankSvc.CreateBatch(c.Request.Context(), companyID(c), req.RunID, req.BankCode, req.PaymentType, req.PaymentDate, userID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, b)
}

func (h *Handler) ListBatches(c *gin.Context) {
	list, err := h.BankSvc.ListBatches(c.Request.Context(), companyID(c), uuidPtr(qs(c, "run_id")), c.Query("status"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, list)
}

func (h *Handler) GetBatch(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	b, err := h.BankSvc.GetBatch(c.Request.Context(), companyID(c), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "batch not found"})
		return
	}
	success(c, b)
}

func (h *Handler) GetBatchItems(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	items, err := h.BankSvc.GetBatchItems(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "items not found"})
		return
	}
	success(c, items)
}

type GenerateBankFileReq struct {
	Format string `json:"format" binding:"required"`
}

func (h *Handler) GenerateBankFile(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req GenerateBankFileReq
	if !bindJSON(c, &req) {
		return
	}
	if err := h.BankSvc.GenerateBankFile(c.Request.Context(), id, req.Format); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "file generated"})
}

func (h *Handler) SendBatch(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.BankSvc.SendBatch(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "batch sent"})
}

func (h *Handler) GetBatchSummary(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	s, err := h.BankSvc.GetBatchSummary(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, s)
}
