package onboarding

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Templates

func (h *Handler) CreateTemplate(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	t, err := h.service.CreateTemplate(c.Request.Context(), companyID, userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, t)
}

func (h *Handler) GetTemplate(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	t, err := h.service.GetTemplate(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Template not found")
		return
	}
	response.Success(c, t)
}

func (h *Handler) GetTemplateWithTasks(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	t, err := h.service.GetTemplateWithTasks(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Template not found")
		return
	}
	response.Success(c, t)
}

func (h *Handler) ListTemplates(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	ts, err := h.service.ListTemplates(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, ts)
}

func (h *Handler) UpdateTemplate(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	t, err := h.service.UpdateTemplate(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, t)
}

func (h *Handler) DeleteTemplate(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.DeleteTemplate(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.NoContent(c)
}

// Template Tasks

func (h *Handler) AddTemplateTask(c *gin.Context) {
	var req CreateTemplateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	t, err := h.service.CreateTemplateTask(c.Request.Context(), c.Param("id"), &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, t)
}

func (h *Handler) ListTemplateTasks(c *gin.Context) {
	ts, err := h.service.ListTemplateTasks(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, ts)
}

func (h *Handler) UpdateTemplateTask(c *gin.Context) {
	var req UpdateTemplateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	t, err := h.service.UpdateTemplateTask(c.Request.Context(), c.Param("taskId"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, t)
}

func (h *Handler) DeleteTemplateTask(c *gin.Context) {
	if err := h.service.DeleteTemplateTask(c.Request.Context(), c.Param("taskId")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.NoContent(c)
}

// Processes

func (h *Handler) CreateOnboarding(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req CreateOnboardingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	p, err := h.service.CreateOnboarding(c.Request.Context(), companyID, userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, p)
}

func (h *Handler) GetProcess(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	p, err := h.service.GetProcessWithDetails(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			response.NotFound(c, "Onboarding process not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, p)
}

func (h *Handler) ListProcesses(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	filters := OnboardingFilters{
		Status:     c.Query("status"),
		EmployeeID: c.Query("employee_id"),
		Search:     c.Query("search"),
	}
	ps, err := h.service.ListProcesses(c.Request.Context(), companyID, filters)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, ps)
}

func (h *Handler) UpdateProcess(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req UpdateOnboardingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	p, err := h.service.UpdateProcess(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, p)
}

func (h *Handler) StartOnboarding(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.StartOnboarding(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "onboarding started"})
}

func (h *Handler) HoldOnboarding(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.HoldOnboarding(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "onboarding on hold"})
}

func (h *Handler) ResumeOnboarding(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.ResumeOnboarding(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "onboarding resumed"})
}

func (h *Handler) CompleteOnboarding(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.CompleteOnboarding(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "onboarding completed"})
}

func (h *Handler) CancelOnboarding(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CancelOnboardingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "reason is required")
		return
	}
	if err := h.service.CancelOnboarding(c.Request.Context(), companyID, c.Param("id"), req.Reason); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "onboarding cancelled"})
}

// Tasks

func (h *Handler) ListTasks(c *gin.Context) {
	ts, err := h.service.ListTasks(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, ts)
}

func (h *Handler) CreateTask(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	t, err := h.service.CreateTask(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, t)
}

func (h *Handler) GetTask(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	t, err := h.service.GetTask(c.Request.Context(), companyID, c.Param("taskId"))
	if err != nil {
		response.NotFound(c, "Task not found")
		return
	}
	response.Success(c, t)
}

func (h *Handler) UpdateTask(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	t, err := h.service.UpdateTask(c.Request.Context(), companyID, c.Param("taskId"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, t)
}

func (h *Handler) StartTask(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.StartTask(c.Request.Context(), companyID, c.Param("taskId")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "task started"})
}

func (h *Handler) CompleteTask(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.CompleteTask(c.Request.Context(), companyID, c.Param("taskId")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "task completed"})
}

func (h *Handler) BlockTask(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "reason is required")
		return
	}
	if err := h.service.BlockTask(c.Request.Context(), companyID, c.Param("taskId"), req.Reason); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "task blocked"})
}

// Documents

func (h *Handler) ListDocuments(c *gin.Context) {
	ds, err := h.service.ListDocuments(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, ds)
}

func (h *Handler) CreateDocumentRequirement(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req UploadDocumentRequest
	if err := c.ShouldBind(&req); err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}
	d, err := h.service.CreateDocumentRequirement(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, d)
}

func (h *Handler) ReviewDocument(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req ReviewDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	if req.Status == "APPROVED" {
		if err := h.service.ApproveDocument(c.Request.Context(), companyID, c.Param("id"), userID); err != nil {
			response.BadRequest(c, err.Error())
			return
		}
	} else if req.Status == "REJECTED" {
		reason := ""
		if req.Reason != nil {
			reason = *req.Reason
		}
		if err := h.service.RejectDocument(c.Request.Context(), companyID, c.Param("id"), userID, reason); err != nil {
			response.BadRequest(c, err.Error())
			return
		}
	} else {
		if err := h.service.repo.UpdateDocumentStatus(c.Request.Context(), companyID, c.Param("id"), req.Status); err != nil {
			response.BadRequest(c, err.Error())
			return
		}
	}
	response.Success(c, gin.H{"message": "document updated"})
}

// Assets

func (h *Handler) ListAssets(c *gin.Context) {
	as, err := h.service.ListAssets(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, as)
}

func (h *Handler) CreateAsset(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	a, err := h.service.CreateAsset(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, a)
}

func (h *Handler) AssignAsset(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	if err := h.service.AssignAsset(c.Request.Context(), companyID, c.Param("id"), userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "asset assigned"})
}

func (h *Handler) DeliverAsset(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.DeliverAsset(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "asset delivered"})
}

func (h *Handler) ReturnAsset(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.ReturnAsset(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "asset returned"})
}

// Access Requests

func (h *Handler) ListAccessRequests(c *gin.Context) {
	ars, err := h.service.ListAccessRequests(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, ars)
}

func (h *Handler) CreateAccessRequest(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateAccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	ar, err := h.service.CreateAccessRequest(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, ar)
}

func (h *Handler) ApproveAccess(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	if err := h.service.ApproveAccess(c.Request.Context(), companyID, c.Param("id"), userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "access approved"})
}

func (h *Handler) RejectAccess(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.RejectAccess(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "access rejected"})
}

func (h *Handler) ActivateAccess(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.ActivateAccess(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "access activated"})
}

func (h *Handler) RevokeAccess(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.RevokeAccess(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "access revoked"})
}

// Milestones

func (h *Handler) ListMilestones(c *gin.Context) {
	ms, err := h.service.ListMilestones(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, ms)
}

func (h *Handler) CreateMilestone(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateMilestoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	m, err := h.service.CreateMilestone(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, m)
}

func (h *Handler) UpdateMilestone(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req UpdateMilestoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	m, err := h.service.UpdateMilestone(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, m)
}

func (h *Handler) CompleteMilestone(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.CompleteMilestone(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "milestone completed"})
}

// Feedback

func (h *Handler) ListFeedback(c *gin.Context) {
	fs, err := h.service.ListFeedback(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, fs)
}

func (h *Handler) CreateFeedback(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req CreateFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	f, err := h.service.CreateFeedback(c.Request.Context(), companyID, c.Param("id"), userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, f)
}

// Buddies

func (h *Handler) AssignBuddy(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req AssignBuddyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	b, err := h.service.AssignBuddy(c.Request.Context(), companyID, c.Param("id"), userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, b)
}

func (h *Handler) GetBuddy(c *gin.Context) {
	b, err := h.service.GetBuddy(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.NotFound(c, "Buddy not found")
		return
	}
	response.Success(c, b)
}

// Exceptions

func (h *Handler) ListExceptions(c *gin.Context) {
	es, err := h.service.ListExceptions(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, es)
}

func (h *Handler) CreateException(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req CreateExceptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	e, err := h.service.CreateException(c.Request.Context(), companyID, c.Param("id"), userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, e)
}

// Training

func (h *Handler) ListTraining(c *gin.Context) {
	ts, err := h.service.ListTraining(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, ts)
}

func (h *Handler) CreateTraining(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateTrainingAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	t, err := h.service.CreateTraining(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, t)
}

// Dashboard

func (h *Handler) GetDashboard(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	d, err := h.service.GetDashboard(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, d)
}

func (h *Handler) GetEmployeeDashboard(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	d, err := h.service.GetEmployeeDashboard(c.Request.Context(), companyID, userID)
	if err != nil {
		response.NotFound(c, "Onboarding process not found for this employee")
		return
	}
	response.Success(c, d)
}

// Candidate Hired (FASE 15 integration)
func (h *Handler) HandleCandidateHired(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var event CandidateHiredEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	p, err := h.service.HandleCandidateHired(c.Request.Context(), companyID, &event)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, p)
}

// IA Assistant
func (h *Handler) GenerateTemplateProposal(c *gin.Context) {
	var req IAOnboardingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	proposal, err := h.service.GenerateTemplateProposal(c.Request.Context(), &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, proposal)
}
