package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

func (h *Handler) CreateApplication(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.ApplicationSvc.Create(c.Request.Context(), companyID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) ListApplications(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.ApplicationSvc.List(c.Request.Context(), companyID, c.Request.URL.Query())
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) GetApplication(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.ApplicationSvc.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Application not found")
		return
	}
	response.Success(c, data)
}

func (h *Handler) MoveStage(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.ApplicationSvc.MoveStage(c.Request.Context(), companyID, c.Param("id"), userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) RejectApplication(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	if err := h.ApplicationSvc.Reject(c.Request.Context(), companyID, c.Param("id"), req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "application rejected"})
}

func (h *Handler) WithdrawApplication(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.ApplicationSvc.Withdraw(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "application withdrawn"})
}

func (h *Handler) GetStageHistory(c *gin.Context) {
	data, err := h.ApplicationSvc.GetStageHistory(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) ListApplicationNotes(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.ApplicationSvc.ListNotes(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) AddApplicationNote(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.ApplicationSvc.AddNote(c.Request.Context(), companyID, c.Param("id"), userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) UpdateApplicationNote(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.ApplicationSvc.UpdateNote(c.Request.Context(), companyID, c.Param("noteId"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) ListApplicationRatings(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.ApplicationSvc.ListRatings(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) AddApplicationRating(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.ApplicationSvc.AddRating(c.Request.Context(), companyID, c.Param("id"), userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}
