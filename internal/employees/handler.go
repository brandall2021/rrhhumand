package employees

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/models"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

type EmployeeHandler struct {
	service *EmployeeService
}

func NewEmployeeHandler(service *EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: service}
}

func (h *EmployeeHandler) Create(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)

	var req CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	emp, err := h.service.Create(c.Request.Context(), companyID, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, emp)
}

func (h *EmployeeHandler) GetByID(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	emp, err := h.service.GetByID(c.Request.Context(), id, companyID)
	if err != nil {
		if err.Error() == "employee not found" {
			response.NotFound(c, "Employee not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, emp)
}

func (h *EmployeeHandler) List(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	params := models.NewPaginationParams(c)

	filters := EmployeeFilters{
		Search:       c.Query("search"),
		Status:       c.Query("status"),
		DepartmentID: c.Query("department_id"),
		BranchID:     c.Query("branch_id"),
		PositionID:   c.Query("position_id"),
		ManagerID:    c.Query("manager_id"),
		SortBy:       c.Query("sort_by"),
		SortDir:      c.Query("sort_dir"),
	}

	employees, total, err := h.service.List(c.Request.Context(), companyID, params, filters)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := params.ToMeta(total)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    employees,
		"meta":    meta,
	})
}

func (h *EmployeeHandler) Update(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	var req UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	emp, err := h.service.Update(c.Request.Context(), id, companyID, &req)
	if err != nil {
		if err.Error() == "employee not found" {
			response.NotFound(c, "Employee not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, emp)
}

func (h *EmployeeHandler) Delete(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id, companyID); err != nil {
		if err.Error() == "employee not found" {
			response.NotFound(c, "Employee not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.NoContent(c)
}

func (h *EmployeeHandler) GetContacts(c *gin.Context) {
	id := c.Param("id")

	contacts, err := h.service.GetContacts(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, contacts)
}

func (h *EmployeeHandler) UpsertContacts(c *gin.Context) {
	id := c.Param("id")

	var contacts []models.EmployeeContact
	if err := c.ShouldBindJSON(&contacts); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.service.UpsertContacts(c.Request.Context(), id, contacts); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	updated, _ := h.service.GetContacts(c.Request.Context(), id)
	response.Success(c, updated)
}

func (h *EmployeeHandler) GetAddresses(c *gin.Context) {
	id := c.Param("id")

	addresses, err := h.service.GetAddresses(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, addresses)
}

func (h *EmployeeHandler) UpsertAddresses(c *gin.Context) {
	id := c.Param("id")

	var addresses []models.EmployeeAddress
	if err := c.ShouldBindJSON(&addresses); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.service.UpsertAddresses(c.Request.Context(), id, addresses); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	updated, _ := h.service.GetAddresses(c.Request.Context(), id)
	response.Success(c, updated)
}

func (h *EmployeeHandler) GetEmergencyContacts(c *gin.Context) {
	id := c.Param("id")

	contacts, err := h.service.GetEmergencyContacts(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, contacts)
}

func (h *EmployeeHandler) UpsertEmergencyContacts(c *gin.Context) {
	id := c.Param("id")

	var contacts []models.EmployeeEmergencyContact
	if err := c.ShouldBindJSON(&contacts); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.service.UpsertEmergencyContacts(c.Request.Context(), id, contacts); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	updated, _ := h.service.GetEmergencyContacts(c.Request.Context(), id)
	response.Success(c, updated)
}

func (h *EmployeeHandler) GetHistory(c *gin.Context) {
	id := c.Param("id")

	history, err := h.service.GetHistory(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, history)
}
