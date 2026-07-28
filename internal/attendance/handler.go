package attendance

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

func (h *Handler) ClockIn(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)

	employeeID, err := h.service.repo.GetUserEmployeeID(c.Request.Context(), companyID, userID)
	if err != nil {
		response.NotFound(c, "Employee profile not found")
		return
	}

	var req ClockInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req = ClockInRequest{Source: "WEB"}
	}
	if req.Source == "" {
		req.Source = "WEB"
	}

	record, err := h.service.ClockIn(c.Request.Context(), companyID, employeeID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, record)
}

func (h *Handler) ClockOut(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)

	employeeID, err := h.service.repo.GetUserEmployeeID(c.Request.Context(), companyID, userID)
	if err != nil {
		response.NotFound(c, "Employee profile not found")
		return
	}

	var req ClockOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req = ClockOutRequest{Source: "WEB"}
	}

	record, err := h.service.ClockOut(c.Request.Context(), companyID, employeeID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, record)
}

func (h *Handler) StartBreak(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)

	employeeID, err := h.service.repo.GetUserEmployeeID(c.Request.Context(), companyID, userID)
	if err != nil {
		response.NotFound(c, "Employee profile not found")
		return
	}

	var req BreakStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req = BreakStartRequest{Source: "WEB"}
	}

	if err := h.service.StartBreak(c.Request.Context(), companyID, employeeID, &req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "break started"})
}

func (h *Handler) EndBreak(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)

	employeeID, err := h.service.repo.GetUserEmployeeID(c.Request.Context(), companyID, userID)
	if err != nil {
		response.NotFound(c, "Employee profile not found")
		return
	}

	var req BreakEndRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req = BreakEndRequest{Source: "WEB"}
	}

	if err := h.service.EndBreak(c.Request.Context(), companyID, employeeID, &req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "break ended"})
}

func (h *Handler) GetMyAttendance(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)

	employeeID, err := h.service.repo.GetUserEmployeeID(c.Request.Context(), companyID, userID)
	if err != nil {
		response.NotFound(c, "Employee profile not found")
		return
	}

	fromStr := c.DefaultQuery("from", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	toStr := c.DefaultQuery("to", time.Now().Format("2006-01-02"))
	from, _ := time.Parse("2006-01-02", fromStr)
	to, _ := time.Parse("2006-01-02", toStr)
	to = to.AddDate(0, 0, 1)

	records, err := h.service.GetMyAttendance(c.Request.Context(), companyID, employeeID, from, to)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, records)
}

func (h *Handler) ListRecords(c *gin.Context) {
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

	filters := AttendanceFilters{
		EmployeeID:   c.Query("employee_id"),
		DepartmentID: c.Query("department_id"),
		Status:       c.Query("status"),
		DateFrom:     c.Query("date_from"),
		DateTo:       c.Query("date_to"),
		BranchID:     c.Query("branch_id"),
	}

	records, total, err := h.service.ListRecords(c.Request.Context(), companyID, filters, offset, limit)
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
		"data":    records,
		"meta": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

func (h *Handler) GetRecord(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")
	rec, err := h.service.GetRecordByID(c.Request.Context(), companyID, id)
	if err != nil {
		response.NotFound(c, "Record not found")
		return
	}
	response.Success(c, rec)
}

func (h *Handler) GetDashboard(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	dash, err := h.service.GetDashboard(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, dash)
}

func (h *Handler) GetTeamDashboard(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)

	employeeID, err := h.service.repo.GetUserEmployeeID(c.Request.Context(), companyID, userID)
	if err != nil {
		response.NotFound(c, "Employee profile not found")
		return
	}

	records, err := h.service.GetTeamRecords(c.Request.Context(), companyID, employeeID, time.Now())
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, records)
}

func (h *Handler) CreatePolicy(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreatePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	policy, err := h.service.CreatePolicy(c.Request.Context(), companyID, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, policy)
}

func (h *Handler) CreateLocation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	loc, err := h.service.CreateLocation(c.Request.Context(), companyID, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, loc)
}

func (h *Handler) CreateDevice(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	dev, err := h.service.CreateDevice(c.Request.Context(), companyID, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, dev)
}

func (h *Handler) CreateCorrection(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)

	employeeID, err := h.service.repo.GetUserEmployeeID(c.Request.Context(), companyID, userID)
	if err != nil {
		response.NotFound(c, "Employee profile not found")
		return
	}

	var req CreateCorrectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	corr, err := h.service.CreateCorrection(c.Request.Context(), companyID, employeeID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, corr)
}

func (h *Handler) ListCorrections(c *gin.Context) {
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

	status := c.Query("status")
	corrections, total, err := h.service.ListCorrections(c.Request.Context(), companyID, status, offset, limit)
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
		"data":    corrections,
		"meta": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

func (h *Handler) ApproveCorrection(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	id := c.Param("id")

	if err := h.service.ApproveCorrection(c.Request.Context(), companyID, id, userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "correction approved"})
}

func (h *Handler) RejectCorrection(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	id := c.Param("id")

	if err := h.service.RejectCorrection(c.Request.Context(), companyID, id, userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "correction rejected"})
}

func (h *Handler) GetCalendar(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)

	employeeID := c.Query("employee_id")
	if employeeID == "" {
		var err error
		employeeID, err = h.service.repo.GetUserEmployeeID(c.Request.Context(), companyID, userID)
		if err != nil {
			response.NotFound(c, "Employee profile not found")
			return
		}
	}

	fromStr := c.DefaultQuery("from", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	toStr := c.DefaultQuery("to", time.Now().Format("2006-01-02"))
	from, _ := time.Parse("2006-01-02", fromStr)
	to, _ := time.Parse("2006-01-02", toStr)
	to = to.AddDate(0, 0, 1)

	records, err := h.service.GetCalendar(c.Request.Context(), companyID, employeeID, from, to)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, records)
}

func (h *Handler) ExportCSV(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	filters := AttendanceFilters{
		EmployeeID:   c.Query("employee_id"),
		DepartmentID: c.Query("department_id"),
		Status:       c.Query("status"),
		DateFrom:     c.Query("date_from"),
		DateTo:       c.Query("date_to"),
	}

	records, err := h.service.ExportCSV(c.Request.Context(), companyID, filters)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=attendance.csv")

	csv := "Date,Employee,Status,Worked Minutes,Late Minutes,Early Leave,Overtime\n"
	for _, r := range records {
		csv += r.WorkDate.Format("2006-01-02") + "," + r.EmployeeName + "," + r.Status + "," +
			strconv.Itoa(r.WorkedMinutes) + "," + strconv.Itoa(r.LateMinutes) + "," +
			strconv.Itoa(r.EarlyLeaveMinutes) + "," + strconv.Itoa(r.OvertimeMinutes) + "\n"
	}
	c.String(http.StatusOK, csv)
}
