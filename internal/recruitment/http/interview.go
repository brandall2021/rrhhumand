package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/recruitment/application"
	"github.com/rrhhumand/api/internal/recruitment/domain"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

func (h *Handler) CreateInterview(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req application.CreateInterviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.InterviewSvc.Create(c.Request.Context(), companyID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) ListInterviews(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.InterviewSvc.List(c.Request.Context(), companyID, c.Query("application_id"), c.Query("status"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) GetInterview(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.InterviewSvc.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Interview not found")
		return
	}
	response.Success(c, data)
}

func (h *Handler) UpdateInterview(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req domain.Interview
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.InterviewSvc.Update(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) CancelInterview(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.InterviewSvc.Cancel(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "interview cancelled"})
}

func (h *Handler) CompleteInterview(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.InterviewSvc.Complete(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "interview completed"})
}

func (h *Handler) AddPanelMember(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req domain.InterviewPanelMember
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.InterviewSvc.AddPanelMember(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) RemovePanelMember(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.InterviewSvc.RemovePanelMember(c.Request.Context(), companyID, c.Query("interview_id"), c.Param("panelId")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "panel member removed"})
}

func (h *Handler) ListPanelMembers(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.InterviewSvc.ListPanelMembers(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) SubmitInterviewFeedback(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req application.SubmitFeedbackReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.InterviewSvc.SubmitFeedback(c.Request.Context(), companyID, c.Param("id"), userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) ListInterviewFeedback(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.InterviewSvc.ListFeedback(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}
