package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

func (h *Handler) GetDashboard(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.DashboardSvc.GetDashboard(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) GetFunnel(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.DashboardSvc.GetFunnel(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) GetTimeToHire(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.DashboardSvc.GetTimeToHire(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}
