package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/benefits/domain"
)

func (h *Handler) CreateWorkflow(c *gin.Context) {
	var req domain.BenefitWorkflow
	if !bindJSON(c, &req) {
		return
	}
	w, err := h.WorkflowSvc.CreateWorkflow(c.Request.Context(), companyID(c), userID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, w)
}

func (h *Handler) ListWorkflows(c *gin.Context) {
	var benefitID *uuid.UUID
	if bid := c.Query("benefit_id"); bid != "" {
		id, err := uuid.Parse(bid)
		if err == nil {
			benefitID = &id
		}
	}
	workflows, err := h.WorkflowSvc.ListWorkflows(c.Request.Context(), companyID(c), benefitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, workflows)
}

func (h *Handler) GetWorkflow(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	w, err := h.WorkflowSvc.GetWorkflow(c.Request.Context(), companyID(c), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	success(c, w)
}

func (h *Handler) UpdateWorkflow(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req domain.BenefitWorkflow
	if !bindJSON(c, &req) {
		return
	}
	req.ID = id
	w, err := h.WorkflowSvc.UpdateWorkflow(c.Request.Context(), companyID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, w)
}

func (h *Handler) AddStep(c *gin.Context) {
	var req domain.BenefitWorkflowStep
	if !bindJSON(c, &req) {
		return
	}
	step, err := h.WorkflowSvc.AddStep(c.Request.Context(), companyID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, step)
}

func (h *Handler) ListSteps(c *gin.Context) {
	workflowID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workflow id"})
		return
	}
	steps, err := h.WorkflowSvc.ListSteps(c.Request.Context(), workflowID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, steps)
}

func (h *Handler) UpdateStep(c *gin.Context) {
	id, err := uuid.Parse(c.Param("stepId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid step id"})
		return
	}
	var req domain.BenefitWorkflowStep
	if !bindJSON(c, &req) {
		return
	}
	req.ID = id
	step, err := h.WorkflowSvc.UpdateStep(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, step)
}

func (h *Handler) DeleteStep(c *gin.Context) {
	id, err := uuid.Parse(c.Param("stepId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid step id"})
		return
	}
	if err := h.WorkflowSvc.DeleteStep(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
