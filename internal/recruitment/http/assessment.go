package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/recruitment/application"
	"github.com/rrhhumand/api/internal/recruitment/domain"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

func (h *Handler) CreateAssessment(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req application.CreateAssessmentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.AssessmentSvc.Create(c.Request.Context(), companyID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) ListAssessments(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	applicationID := c.Query("application_id")
	status := c.Query("status")
	data, err := h.AssessmentSvc.List(c.Request.Context(), companyID, applicationID, status)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) GetAssessment(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.AssessmentSvc.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Assessment not found")
		return
	}
	response.Success(c, data)
}

func (h *Handler) UpdateAssessment(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req domain.Assessment
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.AssessmentSvc.Update(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) SendAssessment(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.AssessmentSvc.Send(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "assessment sent"})
}

func (h *Handler) ScoreAssessment(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req application.ScoreAssessmentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	if err := h.AssessmentSvc.Score(c.Request.Context(), companyID, c.Param("id"), &req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "assessment scored"})
}

func (h *Handler) CancelAssessment(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.AssessmentSvc.Cancel(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "assessment cancelled"})
}

func (h *Handler) ListAssessmentSections(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.AssessmentSvc.ListSections(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) AddAssessmentSection(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req domain.AssessmentSection
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.AssessmentSvc.AddSection(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) ListAssessmentResults(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.AssessmentSvc.ListResults(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) AddAssessmentResult(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req domain.AssessmentResult
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.AssessmentSvc.AddResult(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}
