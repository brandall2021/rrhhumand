package overtime

import (
	"net/http"
	"strconv"

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

// Policies
func (h *Handler) CreatePolicy(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateOvertimePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	policy, err := h.service.CreatePolicy(c.Request.Context(), companyID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, policy)
}

func (h *Handler) GetPolicy(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	policy, err := h.service.GetPolicy(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Policy not found")
		return
	}
	response.Success(c, policy)
}

func (h *Handler) ListPolicies(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	policies, err := h.service.ListPolicies(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, policies)
}

func (h *Handler) UpdatePolicy(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req UpdateOvertimePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	policy, err := h.service.UpdatePolicy(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, policy)
}

func (h *Handler) DeletePolicy(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.DeletePolicy(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.NoContent(c)
}

// Records
func (h *Handler) ListRecords(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	filters := OvertimeFilters{
		EmployeeID:   c.Query("employee_id"),
		Status:       c.Query("status"),
		OvertimeType: c.Query("overtime_type"),
		DateFrom:     c.Query("date_from"),
		DateTo:       c.Query("date_to"),
	}
	records, err := h.service.ListRecords(c.Request.Context(), companyID, filters)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, records)
}

func (h *Handler) GetRecord(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	rec, err := h.service.GetRecord(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Overtime record not found")
		return
	}
	response.Success(c, rec)
}

func (h *Handler) ApproveRecord(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req ApproveOvertimeRequest
	c.ShouldBindJSON(&req)
	approvedMinutes := 0
	if req.ApprovedMinutes != nil {
		approvedMinutes = *req.ApprovedMinutes
	}
	if err := h.service.ApproveRecord(c.Request.Context(), companyID, c.Param("id"), approvedMinutes, userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "overtime approved"})
}

func (h *Handler) RejectRecord(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req RejectOvertimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Reason is required")
		return
	}
	if err := h.service.RejectRecord(c.Request.Context(), companyID, c.Param("id"), req.Reason, userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "overtime rejected"})
}

// Requests
func (h *Handler) CreateRequest(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	employeeID, err := h.service.repo.GetEmployeeIDFromUser(c.Request.Context(), companyID, userID)
	if err != nil {
		response.NotFound(c, "Employee profile not found")
		return
	}
	var req RequestOvertimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	overtimeReq, err := h.service.CreateRequest(c.Request.Context(), companyID, employeeID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, overtimeReq)
}

func (h *Handler) ListRequests(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	filters := OvertimeFilters{
		EmployeeID: c.Query("employee_id"),
		Status:     c.Query("status"),
		DateFrom:   c.Query("date_from"),
		DateTo:     c.Query("date_to"),
	}
	requests, err := h.service.ListRequests(c.Request.Context(), companyID, filters)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, requests)
}

func (h *Handler) GetRequest(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	req, err := h.service.GetRequest(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Request not found")
		return
	}
	response.Success(c, req)
}

func (h *Handler) ApproveRequest(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req ApproveOvertimeRequest
	c.ShouldBindJSON(&req)
	approvedMinutes := 0
	if req.ApprovedMinutes != nil {
		approvedMinutes = *req.ApprovedMinutes
	}
	if err := h.service.ApproveRequest(c.Request.Context(), companyID, c.Param("id"), approvedMinutes, userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "request approved"})
}

func (h *Handler) RejectRequest(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req RejectOvertimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Reason is required")
		return
	}
	if err := h.service.RejectRequest(c.Request.Context(), companyID, c.Param("id"), req.Reason, userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "request rejected"})
}

// Detect
func (h *Handler) DetectOvertime(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req DetectOvertimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	records, count, err := h.service.DetectOvertime(c.Request.Context(), companyID, req.DateFrom, req.DateTo)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"records": records, "count": count})
}

// Compensations
func (h *Handler) RequestCompensation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	employeeID, err := h.service.repo.GetEmployeeIDFromUser(c.Request.Context(), companyID, userID)
	if err != nil {
		response.NotFound(c, "Employee profile not found")
		return
	}
	var req RequestCompensationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	comp, err := h.service.RequestCompensation(c.Request.Context(), companyID, employeeID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, comp)
}

func (h *Handler) ListCompensations(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	filters := OvertimeFilters{
		EmployeeID: c.Query("employee_id"),
		Status:     c.Query("status"),
	}
	comps, err := h.service.ListCompensations(c.Request.Context(), companyID, filters)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, comps)
}

func (h *Handler) ApproveCompensation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	if err := h.service.ApproveCompensation(c.Request.Context(), companyID, c.Param("id"), userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "compensation approved"})
}

func (h *Handler) RejectCompensation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req RejectOvertimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Reason is required")
		return
	}
	if err := h.service.RejectCompensation(c.Request.Context(), companyID, c.Param("id"), req.Reason, userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "compensation rejected"})
}

func (h *Handler) CancelCompensation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.CancelCompensation(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "compensation cancelled"})
}

// Balance
func (h *Handler) GetBalance(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	employeeID := c.Param("id")
	balance, err := h.service.GetBalance(c.Request.Context(), companyID, employeeID)
	if err != nil {
		response.NotFound(c, "Balance not found")
		return
	}
	response.Success(c, balance)
}

func (h *Handler) AdjustBalance(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	employeeID := c.Param("id")
	var req AdjustBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	if err := h.service.AdjustBalance(c.Request.Context(), companyID, employeeID, req.Minutes, req.Reason, userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "balance adjusted"})
}

func (h *Handler) GetBalanceTransactions(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	employeeID := c.Param("id")
	txs, err := h.service.GetBalanceTransactions(c.Request.Context(), companyID, employeeID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, txs)
}

// Dashboard
func (h *Handler) GetDashboard(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	dash, err := h.service.GetDashboard(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, dash)
}

func (h *Handler) GetEmployeeOvertime(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	employeeID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	if page < 1 { page = 1 }
	if perPage < 1 || perPage > 100 { perPage = 30 }

	filters := OvertimeFilters{
		EmployeeID: employeeID,
		DateFrom:   c.Query("date_from"),
		DateTo:     c.Query("date_to"),
		Status:     c.Query("status"),
	}

	records, err := h.service.ListRecords(c.Request.Context(), companyID, filters)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	start := (page - 1) * perPage
	end := start + perPage
	if start > len(records) { start = len(records) }
	if end > len(records) { end = len(records) }

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    records[start:end],
		"meta": gin.H{
			"page":    page,
			"per_page": perPage,
			"total":   len(records),
		},
	})
}
