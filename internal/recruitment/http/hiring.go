package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/recruitment/domain"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

type createHiringRequest struct {
	OfferID   string `json:"offer_id"`
	CreatedBy string `json:"created_by"`
}

func (h *Handler) CreateHiringProcess(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req createHiringRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.HiringSvc.Create(c.Request.Context(), companyID, req.OfferID, req.CreatedBy)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) ListHiringProcesses(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.HiringSvc.ListByCompany(c.Request.Context(), companyID, c.Query("status"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

type updateBackgroundCheckRequest struct {
	Status string `json:"status"`
	Result string `json:"result"`
}

func (h *Handler) UpdateBackgroundCheck(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req updateBackgroundCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	if err := h.HiringSvc.UpdateBackgroundCheck(c.Request.Context(), companyID, c.Param("id"), req.Status, req.Result); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "background check updated"})
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

func (h *Handler) UpdateMedicalCheck(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req updateBackgroundCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	if err := h.HiringSvc.UpdateMedicalCheck(c.Request.Context(), companyID, c.Param("id"), req.Status, req.Result); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "medical check updated"})
}

type updateDocVerificationRequest struct {
	Status string `json:"status"`
}

func (h *Handler) UpdateDocVerification(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req updateDocVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	if err := h.HiringSvc.UpdateDocVerification(c.Request.Context(), companyID, c.Param("id"), req.Status); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "document verification updated"})
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

type addHiringTaskRequest struct {
	TaskType    string  `json:"task_type"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	AssignedTo  *string `json:"assigned_to,omitempty"`
}

func (h *Handler) AddHiringTask(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req addHiringTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	task := domain.HiringProcessTask{
		TaskType:    req.TaskType,
		Title:       req.Title,
		Description: req.Description,
		AssignedTo:  req.AssignedTo,
	}
	data, err := h.HiringSvc.AddTask(c.Request.Context(), companyID, c.Param("id"), task)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) CompleteHiringTask(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	processID := c.Query("process_id")
	if err := h.HiringSvc.CompleteTask(c.Request.Context(), companyID, processID, c.Param("taskId")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "task completed"})
}
