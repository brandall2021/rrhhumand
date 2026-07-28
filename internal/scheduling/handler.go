package scheduling

import (
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

// Work Schedules
func (h *Handler) CreateSchedule(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	schedule, err := h.service.CreateSchedule(c.Request.Context(), companyID, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, schedule)
}

func (h *Handler) GetSchedule(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	schedule, err := h.service.GetSchedule(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Schedule not found")
		return
	}
	response.Success(c, schedule)
}

func (h *Handler) ListSchedules(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	schedules, err := h.service.ListSchedules(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, schedules)
}

func (h *Handler) UpdateSchedule(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req UpdateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	schedule, err := h.service.UpdateSchedule(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, schedule)
}

func (h *Handler) DeleteSchedule(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.DeleteSchedule(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.NoContent(c)
}

// Shifts
func (h *Handler) CreateShift(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	shift, err := h.service.CreateShift(c.Request.Context(), companyID, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, shift)
}

func (h *Handler) GetShift(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	shift, err := h.service.GetShift(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Shift not found")
		return
	}
	response.Success(c, shift)
}

func (h *Handler) ListShifts(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	shifts, err := h.service.ListShifts(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, shifts)
}

func (h *Handler) UpdateShift(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req UpdateShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	shift, err := h.service.UpdateShift(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, shift)
}

func (h *Handler) DeleteShift(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.DeleteShift(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.NoContent(c)
}

// Assignments
func (h *Handler) AssignSchedule(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	employeeID := c.Param("id")
	var req AssignScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	assignment, err := h.service.AssignSchedule(c.Request.Context(), companyID, employeeID, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, assignment)
}

func (h *Handler) GetEmployeeSchedule(c *gin.Context) {
	employeeID := c.Param("id")
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02"))
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.BadRequest(c, "Invalid date format")
		return
	}
	schedule, err := h.service.GetEmployeeSchedule(c.Request.Context(), employeeID, date)
	if err != nil {
		response.NotFound(c, "No active schedule found")
		return
	}
	response.Success(c, schedule)
}

func (h *Handler) AssignShift(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	employeeID := c.Param("id")
	var req AssignShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	assignment, err := h.service.AssignShift(c.Request.Context(), companyID, employeeID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, assignment)
}

func (h *Handler) GetEmployeeShift(c *gin.Context) {
	employeeID := c.Param("id")
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02"))
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.BadRequest(c, "Invalid date format")
		return
	}
	shift, err := h.service.GetEmployeeShift(c.Request.Context(), employeeID, date)
	if err != nil {
		response.NotFound(c, "No shift assigned")
		return
	}
	response.Success(c, shift)
}

// Rotation Templates
func (h *Handler) CreateRotationTemplate(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateRotationTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	template, err := h.service.CreateRotationTemplate(c.Request.Context(), companyID, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, template)
}

func (h *Handler) GetRotationTemplate(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	template, err := h.service.GetRotationTemplate(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Rotation template not found")
		return
	}
	response.Success(c, template)
}

func (h *Handler) ListRotationTemplates(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	templates, err := h.service.ListRotationTemplates(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, templates)
}

func (h *Handler) AssignRotation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	employeeID := c.Param("id")
	var req AssignRotationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	assignment, err := h.service.AssignRotation(c.Request.Context(), companyID, employeeID, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, assignment)
}

// Calendar
func (h *Handler) GenerateCalendar(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req GenerateCalendarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	entries, err := h.service.GenerateCalendar(c.Request.Context(), companyID, req.EmployeeID, req.From, req.To)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, entries)
}

func (h *Handler) ListCalendar(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	if page < 1 { page = 1 }
	if perPage < 1 || perPage > 100 { perPage = 30 }

	filters := CalendarFilters{
		EmployeeID: c.Query("employee_id"),
		DateFrom:   c.Query("date_from"),
		DateTo:     c.Query("date_to"),
		Status:     c.Query("status"),
	}

	entries, total, err := h.service.ListCalendar(c.Request.Context(), companyID, filters, page, perPage)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 { totalPages++ }
	response.SuccessWithMeta(c, entries, &response.Meta{
		Page: page, PerPage: perPage, Total: total, TotalPages: totalPages,
	})
}

func (h *Handler) GetResolvedSchedule(c *gin.Context) {
	employeeID := c.Param("id")
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02"))
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.BadRequest(c, "Invalid date format")
		return
	}
	resolved, err := h.service.GetResolvedSchedule(c.Request.Context(), employeeID, date)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, resolved)
}

// Exceptions
func (h *Handler) CreateException(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateExceptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	exception, err := h.service.CreateException(c.Request.Context(), companyID, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, exception)
}

func (h *Handler) ListExceptions(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	exceptions, err := h.service.ListExceptions(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, exceptions)
}

// Shift Swaps
func (h *Handler) SwapShift(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)

	var req SwapShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	swap, err := h.service.SwapShift(c.Request.Context(), companyID, userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, swap)
}

func (h *Handler) ApproveSwap(c *gin.Context) {
	userID := tenant.GetUserID(c)
	if err := h.service.ApproveSwap(c.Request.Context(), c.Param("id"), userID); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "swap approved"})
}

func (h *Handler) RejectSwap(c *gin.Context) {
	userID := tenant.GetUserID(c)
	if err := h.service.RejectSwap(c.Request.Context(), c.Param("id"), userID); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "swap rejected"})
}
