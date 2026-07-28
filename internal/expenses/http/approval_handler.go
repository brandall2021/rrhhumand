package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) GetPendingApprovals(c *gin.Context) {
	approvals, err := h.ApprovalSvc.GetPendingApprovals(c.Request.Context(), userID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, approvals)
}

func (h *Handler) ApproveEntity(c *gin.Context) {
	var req struct {
		EntityType string `json:"entity_type" binding:"required"`
		EntityID   string `json:"entity_id" binding:"required"`
		Comment    string `json:"comment"`
	}
	if !bindJSON(c, &req) {
		return
	}
	entityID, err := uuid.Parse(req.EntityID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity_id"})
		return
	}
	if err := h.ApprovalSvc.ApproveEntity(c.Request.Context(), req.EntityType, entityID, userID(c), req.Comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "approved"})
}

func (h *Handler) RejectEntity(c *gin.Context) {
	var req struct {
		EntityType string `json:"entity_type" binding:"required"`
		EntityID   string `json:"entity_id" binding:"required"`
		Reason     string `json:"reason" binding:"required"`
	}
	if !bindJSON(c, &req) {
		return
	}
	entityID, err := uuid.Parse(req.EntityID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity_id"})
		return
	}
	if err := h.ApprovalSvc.RejectEntity(c.Request.Context(), req.EntityType, entityID, userID(c), req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "rejected"})
}

func (h *Handler) ObserveEntity(c *gin.Context) {
	var req struct {
		EntityType  string `json:"entity_type" binding:"required"`
		EntityID    string `json:"entity_id" binding:"required"`
		Observation string `json:"observation" binding:"required"`
	}
	if !bindJSON(c, &req) {
		return
	}
	entityID, err := uuid.Parse(req.EntityID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity_id"})
		return
	}
	if err := h.ApprovalSvc.ObserveEntity(c.Request.Context(), req.EntityType, entityID, userID(c), req.Observation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "observed"})
}
