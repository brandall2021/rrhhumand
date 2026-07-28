package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

func (h *Handler) CreatePosition(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.PositionSvc.Create(c.Request.Context(), companyID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) ListPositions(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.PositionSvc.List(c.Request.Context(), companyID, c.Request.URL.Query())
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) GetPosition(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.PositionSvc.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Position not found")
		return
	}
	response.Success(c, data)
}

func (h *Handler) UpdatePosition(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.PositionSvc.Update(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) ClosePosition(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.PositionSvc.Close(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "position closed"})
}

func (h *Handler) ListPositionSkills(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.PositionSvc.ListSkills(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) AddPositionSkill(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.PositionSvc.AddSkill(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) RemovePositionSkill(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.PositionSvc.RemoveSkill(c.Request.Context(), companyID, c.Param("id"), c.Param("skillId")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "skill removed"})
}
