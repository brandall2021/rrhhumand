package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

func (h *Handler) CreateEmailTemplate(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.EmailSvc.CreateTemplate(c.Request.Context(), companyID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) ListEmailTemplates(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.EmailSvc.ListTemplates(c.Request.Context(), companyID, c.Request.URL.Query())
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) GetEmailTemplate(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.EmailSvc.GetTemplate(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Email template not found")
		return
	}
	response.Success(c, data)
}

func (h *Handler) UpdateEmailTemplate(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.EmailSvc.UpdateTemplate(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) DeleteEmailTemplate(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.EmailSvc.DeleteTemplate(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "email template deleted"})
}

func (h *Handler) SendEmail(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	if err := h.EmailSvc.Send(c.Request.Context(), companyID, req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "email sent"})
}

func (h *Handler) ListSentEmails(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.EmailSvc.ListEmails(c.Request.Context(), companyID, c.Request.URL.Query())
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}
