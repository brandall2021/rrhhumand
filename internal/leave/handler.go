package leave

import (
	"net/http"
	"strconv"
	"time"

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

func (h *Handler) CreateLeaveType(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateLeaveTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	lt, err := h.service.CreateLeaveType(c.Request.Context(), companyID, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, lt)
}

func (h *Handler) GetLeaveType(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")
	lt, err := h.service.GetLeaveType(c.Request.Context(), companyID, id)
	if err != nil {
		response.NotFound(c, "Leave type not found")
		return
	}
	response.Success(c, lt)
}

func (h *Handler) ListLeaveTypes(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	types, err := h.service.ListLeaveTypes(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, types)
}

func (h *Handler) UpdateLeaveType(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")
	var req UpdateLeaveTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	lt, err := h.service.UpdateLeaveType(c.Request.Context(), companyID, id, &req)
	if err != nil {
		response.NotFound(c, "Leave type not found")
		return
	}
	response.Success(c, lt)
}

func (h *Handler) DeleteLeaveType(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")
	if err := h.service.DeleteLeaveType(c.Request.Context(), companyID, id); err != nil {
		response.NotFound(c, "Leave type not found")
		return
	}
	response.Success(c, gin.H{"message": "deleted"})
}

func (h *Handler) CreateLeavePolicy(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateLeavePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	p, err := h.service.CreateLeavePolicy(c.Request.Context(), companyID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, p)
}

func (h *Handler) ListLeavePolicies(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	policies, err := h.service.ListLeavePolicies(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, policies)
}

func (h *Handler) CreateHoliday(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateHolidayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	hol, err := h.service.CreateHoliday(c.Request.Context(), companyID, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, hol)
}

func (h *Handler) ListHolidays(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	from := c.DefaultQuery("from", time.Now().Format("2006-01-02"))
	to := c.DefaultQuery("to", time.Now().AddDate(1, 0, 0).Format("2006-01-02"))
	startDate, _ := time.Parse("2006-01-02", from)
	endDate, _ := time.Parse("2006-01-02", to)
	holidays, err := h.service.GetHolidays(c.Request.Context(), companyID, startDate, endDate)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, holidays)
}

func (h *Handler) DeleteHoliday(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")
	if err := h.service.DeleteHoliday(c.Request.Context(), companyID, id); err != nil {
		response.NotFound(c, "Holiday not found")
		return
	}
	response.Success(c, gin.H{"message": "deleted"})
}

func (h *Handler) GetBalance(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)

	targetEmployeeID := c.Query("employee_id")
	if targetEmployeeID == "" {
		targetEmployeeID = userID
	}

	balances, err := h.service.GetBalances(c.Request.Context(), companyID, targetEmployeeID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, balances)
}

func (h *Handler) AdjustBalance(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req AdjustBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	if err := h.service.AdjustBalance(c.Request.Context(), companyID, &req, userID); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "balance adjusted"})
}

func (h *Handler) CreateLeaveRequest(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req CreateLeaveRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	lr, err := h.service.CreateLeaveRequest(c.Request.Context(), companyID, userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, lr)
}

func (h *Handler) GetLeaveRequest(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")
	lr, err := h.service.GetRequest(c.Request.Context(), companyID, id)
	if err != nil {
		response.NotFound(c, "Leave request not found")
		return
	}
	response.Success(c, lr)
}

func (h *Handler) ListLeaveRequests(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	filters := LeaveFilters{
		EmployeeID:   c.Query("employee_id"),
		LeaveTypeID:  c.Query("leave_type_id"),
		Status:       c.Query("status"),
		DateFrom:     c.Query("date_from"),
		DateTo:       c.Query("date_to"),
		DepartmentID: c.Query("department_id"),
	}

	requests, total, err := h.service.ListRequests(c.Request.Context(), companyID, filters, offset, limit)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    requests,
		"meta": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

func (h *Handler) ApproveRequest(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	id := c.Param("id")
	var req ApproveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req = ApproveRequest{}
	}
	if err := h.service.ApproveRequest(c.Request.Context(), companyID, id, userID, req.Comments); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "approved"})
}

func (h *Handler) RejectRequest(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	id := c.Param("id")
	var req RejectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Comments are required for rejection")
		return
	}
	if err := h.service.RejectRequest(c.Request.Context(), companyID, id, userID, req.Comments); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "rejected"})
}

func (h *Handler) CancelRequest(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	id := c.Param("id")
	if err := h.service.CancelRequest(c.Request.Context(), companyID, id, userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "cancelled"})
}

func (h *Handler) GetRequestHistory(c *gin.Context) {
	id := c.Param("id")
	history, err := h.service.GetHistory(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, history)
}

func (h *Handler) GetCalendar(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	filters := CalendarFilters{
		DateFrom:     c.Query("from"),
		DateTo:       c.Query("to"),
		DepartmentID: c.Query("department_id"),
		BranchID:     c.Query("branch_id"),
		EmployeeID:   c.Query("employee_id"),
	}
	days, err := h.service.GetCalendar(c.Request.Context(), companyID, filters)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, days)
}

func (h *Handler) GetTeamCalendar(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	managerID := tenant.GetUserID(c)
	from := c.DefaultQuery("from", time.Now().Format("2006-01-02"))
	to := c.DefaultQuery("to", time.Now().AddDate(0, 1, 0).Format("2006-01-02"))
	startDate, _ := time.Parse("2006-01-02", from)
	endDate, _ := time.Parse("2006-01-02", to)

	requests, err := h.service.GetTeamAbsences(c.Request.Context(), companyID, managerID, startDate, endDate)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, requests)
}

func (h *Handler) GetReport(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	departmentID := c.Query("department_id")
	year := time.Now().Year()
	if y := c.Query("year"); y != "" {
		year, _ = strconv.Atoi(y)
	}

	// Get all requests for the year
	filters := LeaveFilters{DepartmentID: departmentID, DateFrom: strconv.Itoa(year) + "-01-01", DateTo: strconv.Itoa(year) + "-12-31"}
	requests, total, _ := h.service.ListRequests(c.Request.Context(), companyID, filters, 0, 10000)

	response.Success(c, gin.H{
		"year":   year,
		"total":  total,
		"data":   requests,
	})
}
