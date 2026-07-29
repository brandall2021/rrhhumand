package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/onboarding/domain"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

func (h *Handler) ListOnboardings(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	ps, err := h.onbRepo.List(c.Request.Context(), companyID, c.Query("status"), c.Query("employee_id"), c.Query("search"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, ps)
}

func (h *Handler) CreateOnboarding(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req CreateOnboardingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	active, err := h.onbRepo.HasActiveProcess(c.Request.Context(), companyID, req.EmployeeID)
	if err == nil && active {
		response.BadRequest(c, "Active onboarding already exists for this employee")
		return
	}

	p := &domain.OnboardingProcess{
		CompanyID:      companyID,
		EmployeeID:     req.EmployeeID,
		CandidateID:    req.CandidateID,
		ApplicationID:  req.ApplicationID,
		JobOfferID:     req.JobOfferID,
		TemplateID:     req.TemplateID,
		Status:         domain.OnboardingDraft,
		StartDate:      req.StartDate,
		CompletionPolicy: "STRICT",
		CreatedBy:      userID,
	}
	if req.CompletionPolicy != nil {
		p.CompletionPolicy = *req.CompletionPolicy
	}
	if req.EmployeeType != nil {
		et := domain.EmploymentType(*req.EmployeeType)
		p.EmployeeType = &et
	}
	if req.WorkMode != nil {
		wm := domain.WorkMode(*req.WorkMode)
		p.WorkMode = &wm
	}

	response.Created(c, p)
}

func (h *Handler) GetOnboarding(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	p, err := h.onbRepo.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Onboarding process not found")
		return
	}
	response.Success(c, p)
}

func (h *Handler) UpdateOnboarding(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req UpdateOnboardingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	p, err := h.onbRepo.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Onboarding process not found")
		return
	}
	if req.TemplateID != nil {
		p.TemplateID = req.TemplateID
	}
	if req.CompletionPolicy != nil {
		p.CompletionPolicy = *req.CompletionPolicy
	}
	response.Success(c, p)
}

func (h *Handler) StartOnboarding(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.onboardingEngine.Start(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "onboarding started"})
}

func (h *Handler) CompleteOnboarding(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.onboardingEngine.Complete(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "onboarding completed"})
}

func (h *Handler) CancelOnboarding(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "reason is required")
		return
	}
	if err := h.onboardingEngine.Cancel(c.Request.Context(), companyID, c.Param("id"), req.Reason); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "onboarding cancelled"})
}

func (h *Handler) BlockOnboarding(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.onboardingEngine.Block(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "onboarding blocked"})
}

func (h *Handler) ListOnboardingTasks(c *gin.Context) {
	assignments, err := h.taskRepo.ListAssignments(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, assignments)
}

func (h *Handler) CreateOnboardingTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	response.Created(c, req)
}

func (h *Handler) StartTask(c *gin.Context) {
	response.Success(c, gin.H{"message": "task started"})
}

func (h *Handler) CompleteTask(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.onboardingEngine.ExecuteTask(c.Request.Context(), companyID, c.Param("taskId"), tenant.GetUserID(c)); err != nil {
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
	if err := h.onboardingEngine.BlockTask(c.Request.Context(), companyID, c.Param("taskId"), req.Reason); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "task blocked"})
}

func (h *Handler) ListDocuments(c *gin.Context) {
	ds, err := h.docRepo.ListByOnboarding(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, ds)
}

func (h *Handler) CreateDocument(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req struct {
		DocumentType string `json:"document_type" binding:"required"`
		Name         string `json:"name" binding:"required"`
		Required     *bool  `json:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	required := true
	if req.Required != nil {
		required = *req.Required
	}
	doc := &domain.OnboardingDocument{
		CompanyID:    companyID,
		OnboardingID: c.Param("id"),
		DocumentType: req.DocumentType,
		Name:         req.Name,
		Required:     required,
		Status:       domain.DocPending,
	}
	response.Created(c, doc)
}

func (h *Handler) ApproveDocument(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	if err := h.onboardingEngine.ApproveDocument(c.Request.Context(), companyID, c.Param("docId"), userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "document approved"})
}

func (h *Handler) RejectDocument(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	if err := h.onboardingEngine.RejectDocument(c.Request.Context(), companyID, c.Param("docId"), userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "document rejected"})
}

func (h *Handler) ListOnboardingAssets(c *gin.Context) {
	response.Success(c, []interface{}{})
}

func (h *Handler) CreateOnboardingAsset(c *gin.Context) {
	response.Created(c, gin.H{"message": "asset created"})
}

func (h *Handler) AssignOnboardingAsset(c *gin.Context) {
	response.Success(c, gin.H{"message": "asset assigned"})
}

func (h *Handler) DeliverOnboardingAsset(c *gin.Context) {
	response.Success(c, gin.H{"message": "asset delivered"})
}

func (h *Handler) ListAccessRequests(c *gin.Context) {
	response.Success(c, []interface{}{})
}

func (h *Handler) CreateAccessRequest(c *gin.Context) {
	response.Created(c, gin.H{"message": "access request created"})
}

func (h *Handler) GetChecklist(c *gin.Context) {
	response.Success(c, gin.H{"onboarding_id": c.Param("id")})
}

func (h *Handler) CompleteChecklistItem(c *gin.Context) {
	response.Success(c, gin.H{"message": "checklist item completed"})
}

func (h *Handler) ListNotes(c *gin.Context) {
	response.Success(c, []interface{}{})
}

func (h *Handler) CreateNote(c *gin.Context) {
	response.Created(c, gin.H{"message": "note created"})
}

func (h *Handler) UpdateProbation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req ProbationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	if err := h.onbRepo.UpdateProbation(c.Request.Context(), companyID, c.Param("id"), domain.ProbationStatus(req.Status)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "probation updated"})
}

func (h *Handler) AssignBuddy(c *gin.Context) {
	response.Created(c, gin.H{"message": "buddy assigned"})
}

func (h *Handler) GetBuddy(c *gin.Context) {
	response.Success(c, gin.H{"onboarding_id": c.Param("id")})
}

func (h *Handler) AssignTraining(c *gin.Context) {
	response.Created(c, gin.H{"message": "training assigned"})
}

func (h *Handler) GetOnboardingDashboard(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	active, pending, completed, overdue, avgProgress, err := h.onbRepo.GetDashboardStats(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	dash := &OnboardingDashboardResponse{
		ActiveOnboardings:    active,
		PendingOnboardings:   pending,
		CompletedOnboardings: completed,
		OverdueOnboardings:   overdue,
		AverageProgress:      avgProgress,
	}
	response.Success(c, dash)
}

func (h *Handler) GetEmployeeOnboardingDashboard(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	employeeID := c.Query("employee_id")
	if employeeID == "" {
		employeeID = tenant.GetUserID(c)
	}
	p, err := h.onbRepo.GetByID(c.Request.Context(), companyID, employeeID)
	if err != nil {
		response.NotFound(c, "Onboarding process not found")
		return
	}
	total, completed, _ := h.taskRepo.GetCounts(c.Request.Context(), p.ID)
	dash := &EmployeeDashboardResponse{
		Status:         string(p.Status),
		Progress:       p.Progress,
		TasksTotal:     total,
		TasksCompleted: completed,
	}
	response.Success(c, dash)
}

func (h *Handler) GetCandidatesReadyForOnboarding(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	candidates, err := h.atsSvc.GetReadyForOnboarding(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, candidates)
}
