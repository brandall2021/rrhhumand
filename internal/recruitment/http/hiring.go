package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

func (h *Handler) CreateHiringProcess(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.HiringSvc.Create(c.Request.Context(), companyID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) ListHiringProcesses(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.HiringSvc.List(c.Request.Context(), companyID, c.Request.URL.Query())
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) GetHiringProcess(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.HiringSvc.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Hiring process not found")
		return
	}
	response.Success(c, data)
}

func (h *Handler) UpdateBackgroundCheck(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.HiringSvc.UpdateBackgroundCheck(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) UpdateMedicalCheck(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.HiringSvc.UpdateMedicalCheck(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) UpdateDocVerification(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.HiringSvc.UpdateDocVerification(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) CompleteHiringProcess(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.HiringSvc.Complete(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "hiring process completed"})
}

func (h *Handler) CancelHiringProcess(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.HiringSvc.Cancel(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "hiring process cancelled"})
}

func (h *Handler) ListHiringTasks(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.HiringSvc.ListTasks(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) AddHiringTask(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.HiringSvc.AddTask(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) CompleteHiringTask(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.HiringSvc.CompleteTask(c.Request.Context(), companyID, c.Param("taskId")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "task completed"})
}
