package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

func (h *Handler) MyApplications(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	data, err := h.ApplicationSvc.List(c.Request.Context(), companyID, userID, "", "")
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) MyReferrals(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.ApplicationSvc.List(c.Request.Context(), companyID, "", "", "")
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) CreateReferral(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req struct {
		CandidateID string `json:"candidate_id"`
		PostingID   string `json:"posting_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.ApplicationSvc.Create(c.Request.Context(), companyID, req.CandidateID, req.PostingID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) MyInterviews(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.InterviewSvc.List(c.Request.Context(), companyID, "", "")
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}
