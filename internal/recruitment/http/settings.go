package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

func (h *Handler) CreateSource(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.SettingsSvc.CreateSource(c.Request.Context(), companyID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) ListSources(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.SettingsSvc.ListSources(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) UpdateSource(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.SettingsSvc.UpdateSource(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) CreateStage(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.SettingsSvc.CreateStage(c.Request.Context(), companyID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) ListStages(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.SettingsSvc.ListStages(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) UpdateStage(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.SettingsSvc.UpdateStage(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) ReorderStages(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	if err := h.SettingsSvc.ReorderStages(c.Request.Context(), companyID, req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "stages reordered"})
}

func (h *Handler) CreateTransition(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.SettingsSvc.CreateTransition(c.Request.Context(), companyID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) ListTransitions(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.SettingsSvc.ListTransitions(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) DeleteTransition(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.SettingsSvc.DeleteTransition(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "transition deleted"})
}

func (h *Handler) CreateRejectionReason(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.SettingsSvc.CreateRejectionReason(c.Request.Context(), companyID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) ListRejectionReasons(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.SettingsSvc.ListRejectionReasons(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) UpdateRejectionReason(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.SettingsSvc.UpdateRejectionReason(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}
