package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/recruitment/domain"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

func (h *Handler) CreateWorkflow(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req domain.Workflow
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.WorkflowSvc.Create(c.Request.Context(), companyID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) ListWorkflows(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.WorkflowSvc.List(c.Request.Context(), companyID, c.Query("entity_type"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) GetWorkflow(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.WorkflowSvc.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Workflow not found")
		return
	}
	response.Success(c, data)
}

func (h *Handler) UpdateWorkflow(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req domain.Workflow
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.WorkflowSvc.Update(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) ActivateWorkflow(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.WorkflowSvc.Activate(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "workflow activated"})
}

func (h *Handler) DeactivateWorkflow(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.WorkflowSvc.Deactivate(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "workflow deactivated"})
}

func (h *Handler) ListWorkflowStages(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.WorkflowSvc.ListStages(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) AddWorkflowStage(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req domain.WorkflowStage
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.WorkflowSvc.AddStage(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) RemoveWorkflowStage(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.WorkflowSvc.RemoveStage(c.Request.Context(), companyID, c.Param("id"), c.Param("stageId")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "stage removed"})
}

func (h *Handler) ReorderWorkflowStages(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req struct {
		StageIDs []string `json:"stage_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	if err := h.WorkflowSvc.ReorderStages(c.Request.Context(), companyID, c.Param("id"), req.StageIDs); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "stages reordered"})
}

func (h *Handler) ListWorkflowRules(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.WorkflowSvc.ListRules(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) AddWorkflowRule(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req domain.WorkflowRule
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.WorkflowSvc.AddRule(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) UpdateWorkflowRule(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req domain.WorkflowRule
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	if err := h.WorkflowSvc.UpdateRule(c.Request.Context(), companyID, c.Param("ruleId"), req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "rule updated"})
}

func (h *Handler) DeleteWorkflowRule(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.WorkflowSvc.DeleteRule(c.Request.Context(), companyID, c.Param("id"), c.Param("ruleId")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "rule deleted"})
}
