package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/payroll/features/domain"
)

func (h *Handler) CreateTemplate(c *gin.Context) {
	var req domain.ReceiptTemplate
	if !bindJSON(c, &req) {
		return
	}
	t, err := h.ReceiptSvc.CreateTemplate(c.Request.Context(), companyID(c), userID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, t)
}

func (h *Handler) ListTemplates(c *gin.Context) {
	list, err := h.ReceiptSvc.ListTemplates(c.Request.Context(), companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, list)
}

func (h *Handler) GetTemplate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	t, err := h.ReceiptSvc.GetTemplate(c.Request.Context(), companyID(c), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		return
	}
	success(c, t)
}

func (h *Handler) UpdateTemplate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req domain.ReceiptTemplate
	if !bindJSON(c, &req) {
		return
	}
	req.ID = id
	t, err := h.ReceiptSvc.UpdateTemplate(c.Request.Context(), companyID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, t)
}

func (h *Handler) DeleteTemplate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.ReceiptSvc.DeleteTemplate(c.Request.Context(), companyID(c), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "deleted"})
}

type GenerateReceiptsReq struct {
	RunID       string   `json:"run_id" binding:"required"`
	EmployeeIDs []string `json:"employee_ids" binding:"required"`
}

func (h *Handler) GenerateReceipts(c *gin.Context) {
	var req GenerateReceiptsReq
	if !bindJSON(c, &req) {
		return
	}
	runID, err := uuid.Parse(req.RunID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid run_id"})
		return
	}
	employeeIDs := make([]uuid.UUID, len(req.EmployeeIDs))
	for i, eid := range req.EmployeeIDs {
		employeeIDs[i], err = uuid.Parse(eid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee_id"})
			return
		}
	}
	receipts, err := h.ReceiptSvc.GenerateReceipts(c.Request.Context(), companyID(c), runID, employeeIDs, userID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, receipts)
}

func (h *Handler) ListReceipts(c *gin.Context) {
	list, err := h.ReceiptSvc.ListReceipts(c.Request.Context(), companyID(c), uuidPtr(qs(c, "run_id")), uuidPtr(qs(c, "employee_id")), qi(c, "limit", 20), qi(c, "offset", 0))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, list)
}

func (h *Handler) GetReceipt(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	r, err := h.ReceiptSvc.GetReceipt(c.Request.Context(), companyID(c), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "receipt not found"})
		return
	}
	success(c, r)
}

func (h *Handler) GetReceiptItems(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	items, err := h.ReceiptSvc.GetReceiptItems(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "items not found"})
		return
	}
	success(c, items)
}

func (h *Handler) AcknowledgeReceipt(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.ReceiptSvc.AcknowledgeReceipt(c.Request.Context(), id, c.ClientIP()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "acknowledged"})
}

func (h *Handler) MyReceipts(c *gin.Context) {
	eid := employeeID(c)
	list, err := h.ReceiptSvc.ListReceipts(c.Request.Context(), companyID(c), uuidPtr(qs(c, "run_id")), &eid, qi(c, "limit", 20), qi(c, "offset", 0))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, list)
}
