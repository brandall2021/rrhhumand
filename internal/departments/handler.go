package departments

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/models"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

type DepartmentHandler struct {
	service *DepartmentService
}

func NewDepartmentHandler(service *DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{service: service}
}

func (h *DepartmentHandler) Create(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)

	var req CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	dept, err := h.service.Create(c.Request.Context(), companyID, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, dept)
}

func (h *DepartmentHandler) GetByID(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	dept, err := h.service.GetByID(c.Request.Context(), id, companyID)
	if err != nil {
		if err.Error() == "department not found" {
			response.NotFound(c, "Department not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, dept)
}

func (h *DepartmentHandler) List(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	params := models.NewPaginationParams(c)
	search := c.Query("search")

	depts, total, err := h.service.List(c.Request.Context(), companyID, params, search)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := params.ToMeta(total)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    depts,
		"meta":    meta,
	})
}

func (h *DepartmentHandler) Update(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	var req UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	dept, err := h.service.Update(c.Request.Context(), id, companyID, &req)
	if err != nil {
		if err.Error() == "department not found" {
			response.NotFound(c, "Department not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, dept)
}

func (h *DepartmentHandler) Delete(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id, companyID); err != nil {
		if err.Error() == "department not found" {
			response.NotFound(c, "Department not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.NoContent(c)
}
