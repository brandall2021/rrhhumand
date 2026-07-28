package branches

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/models"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

type BranchHandler struct {
	service *BranchService
}

func NewBranchHandler(service *BranchService) *BranchHandler {
	return &BranchHandler{service: service}
}

func (h *BranchHandler) Create(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)

	var req CreateBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	branch, err := h.service.Create(c.Request.Context(), companyID, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, branch)
}

func (h *BranchHandler) GetByID(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	branch, err := h.service.GetByID(c.Request.Context(), id, companyID)
	if err != nil {
		if err.Error() == "branch not found" {
			response.NotFound(c, "Branch not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, branch)
}

func (h *BranchHandler) List(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	params := models.NewPaginationParams(c)
	search := c.Query("search")

	branches, total, err := h.service.List(c.Request.Context(), companyID, params, search)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := params.ToMeta(total)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    branches,
		"meta":    meta,
	})
}

func (h *BranchHandler) Update(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	var req UpdateBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	branch, err := h.service.Update(c.Request.Context(), id, companyID, &req)
	if err != nil {
		if err.Error() == "branch not found" {
			response.NotFound(c, "Branch not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, branch)
}

func (h *BranchHandler) Delete(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id, companyID); err != nil {
		if err.Error() == "branch not found" {
			response.NotFound(c, "Branch not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.NoContent(c)
}
