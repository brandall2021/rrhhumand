package positions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/models"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

type PositionHandler struct {
	service *PositionService
}

func NewPositionHandler(service *PositionService) *PositionHandler {
	return &PositionHandler{service: service}
}

func (h *PositionHandler) Create(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)

	var req CreatePositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	pos, err := h.service.Create(c.Request.Context(), companyID, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, pos)
}

func (h *PositionHandler) GetByID(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	pos, err := h.service.GetByID(c.Request.Context(), id, companyID)
	if err != nil {
		if err.Error() == "position not found" {
			response.NotFound(c, "Position not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, pos)
}

func (h *PositionHandler) List(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	params := models.NewPaginationParams(c)
	search := c.Query("search")
	departmentID := c.Query("department_id")

	positions, total, err := h.service.List(c.Request.Context(), companyID, params, search, departmentID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := params.ToMeta(total)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    positions,
		"meta":    meta,
	})
}

func (h *PositionHandler) Update(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	var req UpdatePositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	pos, err := h.service.Update(c.Request.Context(), id, companyID, &req)
	if err != nil {
		if err.Error() == "position not found" {
			response.NotFound(c, "Position not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, pos)
}

func (h *PositionHandler) Delete(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id, companyID); err != nil {
		if err.Error() == "position not found" {
			response.NotFound(c, "Position not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.NoContent(c)
}
