package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/recruitment/application"
	"github.com/rrhhumand/api/internal/recruitment/domain"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

func (h *Handler) CreateRequisition(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req application.CreateRequisitionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.RequisitionSvc.Create(c.Request.Context(), companyID, userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) ListRequisitions(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.RequisitionSvc.List(c.Request.Context(), companyID, c.Query("status"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) GetRequisition(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.RequisitionSvc.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Requisition not found")
		return
	}
	response.Success(c, data)
}

func (h *Handler) UpdateRequisition(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req domain.Requisition
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.RequisitionSvc.Update(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) SubmitRequisition(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.RequisitionSvc.Submit(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "requisition submitted for approval"})
}

func (h *Handler) ApproveRequisition(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.RequisitionSvc.Approve(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "requisition approved"})
}

func (h *Handler) OpenRequisition(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.RequisitionSvc.Open(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "requisition opened"})
}

func (h *Handler) CloseRequisition(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.RequisitionSvc.Close(c.Request.Context(), companyID, c.Param("id"), c.Query("reason")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "requisition closed"})
}

func (h *Handler) CancelRequisition(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.RequisitionSvc.Cancel(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "requisition cancelled"})
}

func (h *Handler) ListRequisitionSkills(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.RequisitionSvc.ListSkills(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) AddRequisitionSkill(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req domain.RequisitionSkill
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.RequisitionSvc.AddSkill(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) RemoveRequisitionSkill(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.RequisitionSvc.RemoveSkill(c.Request.Context(), companyID, c.Param("id"), c.Param("skillId")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "skill removed"})
}
