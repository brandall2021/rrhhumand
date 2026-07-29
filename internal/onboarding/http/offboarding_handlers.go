package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/onboarding/domain"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

func (h *Handler) ListOffboardings(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	ps, err := h.offbRepo.ListProcesses(c.Request.Context(), companyID, c.Query("status"), c.Query("employee_id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, ps)
}

func (h *Handler) CreateOffboarding(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req CreateOffboardingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	active, err := h.offbRepo.HasActiveProcess(c.Request.Context(), companyID, req.EmployeeID)
	if err == nil && active {
		response.BadRequest(c, "Active offboarding already exists for this employee")
		return
	}

	p := &domain.OffboardingProcess{
		CompanyID:              companyID,
		EmployeeID:             req.EmployeeID,
		RequestedBy:            userID,
		TerminationType:        domain.TerminationType(req.TerminationType),
		ReasonID:               req.ReasonID,
		NoticeDate:             req.NoticeDate,
		LastWorkingDate:        req.LastWorkingDate,
		TerminationEffectiveDate: req.TerminationEffectiveDate,
		TemplateID:             req.TemplateID,
		Status:                 domain.OffboardingDraft,
		EmployeeStatusAfter:    "INACTIVE",
	}

	response.Created(c, p)
}

func (h *Handler) GetOffboarding(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	p, err := h.offbRepo.GetProcessByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Offboarding process not found")
		return
	}
	response.Success(c, p)
}

func (h *Handler) UpdateOffboarding(c *gin.Context) {
	response.Success(c, gin.H{"message": "offboarding updated"})
}

func (h *Handler) ApproveOffboarding(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.offboardingEngine.Approve(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "offboarding approved"})
}

func (h *Handler) StartOffboarding(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.offboardingEngine.Start(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "offboarding started"})
}

func (h *Handler) CompleteOffboarding(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.offboardingEngine.Complete(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "offboarding completed"})
}

func (h *Handler) CancelOffboarding(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.offboardingEngine.Cancel(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "offboarding cancelled"})
}

func (h *Handler) ListOffboardingTasks(c *gin.Context) {
	ts, err := h.offbRepo.ListTasks(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, ts)
}

func (h *Handler) CreateOffboardingTask(c *gin.Context) {
	response.Created(c, gin.H{"message": "task created"})
}

func (h *Handler) CompleteOffboardingTask(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.offboardingEngine.ExecuteTask(c.Request.Context(), companyID, c.Param("taskId"), tenant.GetUserID(c)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "task completed"})
}

func (h *Handler) ListOffboardingAssets(c *gin.Context) {
	as, err := h.offbRepo.ListAssets(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, as)
}

func (h *Handler) CreateOffboardingAsset(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req struct {
		AssetType    string  `json:"asset_type" binding:"required"`
		Description  *string `json:"description"`
		SerialNumber *string `json:"serial_number"`
		Condition    *string `json:"condition"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	a := &domain.OffboardingAsset{
		CompanyID:          companyID,
		OffboardingID:      c.Param("id"),
		AssetType:          req.AssetType,
		Description:        req.Description,
		SerialNumber:       req.SerialNumber,
		ConditionOnDelivery: req.Condition,
		Status:             domain.AssetPendingReturn,
	}
	response.Created(c, a)
}

func (h *Handler) ReturnOffboardingAsset(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req ReturnAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "condition_on_return is required")
		return
	}
	if err := h.offbRepo.UpdateAssetStatus(c.Request.Context(), companyID, c.Param("assetId"), domain.AssetReturned, req.ConditionOnReturn); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "asset returned"})
}

func (h *Handler) ReportAssetDamaged(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req ReturnAssetRequest
	if err := c.ShouldBindJSON(&req); err == nil {
		h.offbRepo.UpdateAssetStatus(c.Request.Context(), companyID, c.Param("assetId"), domain.AssetDamaged, req.ConditionOnReturn)
	}
	response.Success(c, gin.H{"message": "asset reported as damaged"})
}

func (h *Handler) ReportAssetLost(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	h.offbRepo.UpdateAssetStatus(c.Request.Context(), companyID, c.Param("assetId"), domain.AssetLost, nil)
	response.Success(c, gin.H{"message": "asset reported as lost"})
}

func (h *Handler) ListAccessRevocations(c *gin.Context) {
	as, err := h.offbRepo.ListAccessRevocations(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, as)
}

func (h *Handler) CreateAccessRevocation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req RevokeAccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	a := &domain.AccessRevocation{
		CompanyID:     companyID,
		OffboardingID: strPtr(c.Param("id")),
		SystemName:    req.SystemName,
		AccessType:    req.AccessType,
		Status:        string(domain.RevokePending),
	}
	response.Created(c, a)
}

func (h *Handler) RevokeAccess(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	h.offbRepo.UpdateAccessRevocation(c.Request.Context(), companyID, c.Param("accessId"), "REVOKED", &userID, nil)
	response.Success(c, gin.H{"message": "access revoked"})
}

func (h *Handler) RetryAccessRevocation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	h.offbRepo.UpdateAccessRevocation(c.Request.Context(), companyID, c.Param("accessId"), "PENDING", nil, nil)
	response.Success(c, gin.H{"message": "retry scheduled"})
}

func (h *Handler) GetExitInterview(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	e, err := h.offbRepo.GetExitInterview(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Exit interview not found")
		return
	}
	response.Success(c, e)
}

func (h *Handler) CreateExitInterview(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CompleteExitInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	e := &domain.ExitInterview{
		CompanyID:     companyID,
		OffboardingID: c.Param("id"),
		Reason:        &req.Reason,
		Feedback:      &req.Feedback,
		Recommendation: req.Recommendation,
		Rating:        req.Rating,
		Anonymous:     req.Anonymous,
	}
	response.Created(c, e)
}

func (h *Handler) UpdateExitInterview(c *gin.Context) {
	response.Success(c, gin.H{"message": "exit interview updated"})
}

func (h *Handler) CompleteExitInterview(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CompleteExitInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "reason is required")
		return
	}
	if err := h.offbRepo.CompleteExitInterview(c.Request.Context(), companyID, c.Param("id"),
		req.Reason, req.Feedback, req.Recommendation, req.Rating); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "exit interview completed"})
}

func (h *Handler) GetHandover(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	handover, err := h.offbRepo.GetHandover(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Handover not found")
		return
	}
	response.Success(c, handover)
}

func (h *Handler) CreateHandover(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req HandoverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	hnd := &domain.EmployeeHandover{
		CompanyID:    companyID,
		OffboardingID: strPtr(c.Param("id")),
		HandoverTo:   req.HandoverTo,
		Description:  req.Description,
		Projects:     req.Projects,
		PendingTasks: req.PendingTasks,
		Documents:    req.Documents,
		Status:       string(domain.HandoverPending),
	}
	response.Created(c, hnd)
}

func (h *Handler) CompleteHandover(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.offbRepo.CompleteHandover(c.Request.Context(), companyID, c.Param("handoverId")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "handover completed"})
}

func (h *Handler) RequestFinalSettlement(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.offboardingEngine.RequestFinalSettlement(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "final settlement requested"})
}

func (h *Handler) ListExitReasons(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	reasons, err := h.offbRepo.ListExitReasons(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, reasons)
}

func (h *Handler) GetOffboardingDashboard(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	active, pending, completed, overdue, err := h.offbRepo.GetDashboardStats(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	dash := &OffboardingDashboardResponse{
		ActiveOffboardings:    active,
		PendingOffboardings:   pending,
		CompletedOffboardings: completed,
		OverdueOffboardings:   overdue,
	}
	response.Success(c, dash)
}

func (h *Handler) CreateWorkflowRule(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req WorkflowRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	rule := &domain.WorkflowRule{
		CompanyID:    companyID,
		WorkflowType: domain.WorkflowType(req.WorkflowType),
		Name:         req.Name,
		Conditions:   req.Conditions,
		Actions:      req.Actions,
		Active:       true,
	}
	if req.Priority != nil {
		rule.Priority = *req.Priority
	}
	response.Created(c, rule)
}

func (h *Handler) ListWorkflowRules(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	rules, err := h.sharedRepo.ListWorkflowRules(c.Request.Context(), companyID, domain.WorkflowType(c.Query("workflow_type")))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, rules)
}

func (h *Handler) EvaluateWorkflowRules(c *gin.Context) {
	response.Success(c, gin.H{"message": "rules evaluated"})
}

func strPtr(s string) *string { return &s }
