package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

func (h *Handler) MyApplications(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	data, err := h.ApplicationSvc.ListByUser(c.Request.Context(), companyID, userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) MyReferrals(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	data, err := h.ApplicationSvc.ListReferralsByUser(c.Request.Context(), companyID, userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) CreateReferral(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.ApplicationSvc.CreateReferral(c.Request.Context(), companyID, userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) MyInterviews(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	data, err := h.InterviewSvc.ListByUser(c.Request.Context(), companyID, userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}
