package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/benefits/domain"
	"github.com/shopspring/decimal"
)

func (h *Handler) CreateReimbursement(c *gin.Context) {
	var req domain.BenefitReimbursement
	if !bindJSON(c, &req) {
		return
	}
	r, err := h.ReimbursementSvc.CreateReimbursement(c.Request.Context(), companyID(c), employeeID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, r)
}

func (h *Handler) ListReimbursements(c *gin.Context) {
	var employeeID *uuid.UUID
	if eid := c.Query("employee_id"); eid != "" {
		id, err := uuid.Parse(eid)
		if err == nil {
			employeeID = &id
		}
	}
	status := qs(c, "status")
	reimbursements, err := h.ReimbursementSvc.ListReimbursements(c.Request.Context(), companyID(c), employeeID, status)
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
	r, err := h.ReimbursementSvc.GetReimbursement(c.Request.Context(), companyID(c), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	success(c, r)
}

func (h *Handler) ApproveReimbursement(c *gin.Context) {
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
	if err := h.ReimbursementSvc.ApproveReimbursement(c.Request.Context(), id, userID(c), req.ApprovedAmount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "reimbursement approved"})
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
	success(c, gin.H{"message": "reimbursement rejected"})
}

func (h *Handler) PayReimbursement(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Method    string `json:"method" binding:"required"`
		Reference string `json:"reference" binding:"required"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if err := h.ReimbursementSvc.PayReimbursement(c.Request.Context(), id, req.Method, req.Reference); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "reimbursement paid"})
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
	success(c, gin.H{"message": "reimbursement cancelled"})
}

func (h *Handler) UploadDocument(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reimbursement id"})
		return
	}
	var req struct {
		FileName    string `json:"file_name" binding:"required"`
		FileType    string `json:"file_type" binding:"required"`
		StoragePath string `json:"storage_path" binding:"required"`
		FileSize    int    `json:"file_size"`
	}
	if !bindJSON(c, &req) {
		return
	}
	doc, err := h.ReimbursementSvc.UploadDocument(c.Request.Context(), id, userID(c), req.FileName, req.FileType, req.StoragePath, req.FileSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, doc)
}

func (h *Handler) ListDocuments(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reimbursement id"})
		return
	}
	docs, err := h.ReimbursementSvc.ListDocuments(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, docs)
}
