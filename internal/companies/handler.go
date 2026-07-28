package companies

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/models"
	"github.com/rrhhumand/api/pkg/response"
)

type CompanyHandler struct {
	service *CompanyService
}

func NewCompanyHandler(service *CompanyService) *CompanyHandler {
	return &CompanyHandler{service: service}
}

func (h *CompanyHandler) Create(c *gin.Context) {
	var req CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	company, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, company)
}

func (h *CompanyHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	company, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "company not found" {
			response.NotFound(c, "Company not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, company)
}

func (h *CompanyHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	company, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "company not found" {
			response.NotFound(c, "Company not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, company)
}

func (h *CompanyHandler) List(c *gin.Context) {
	params := models.NewPaginationParams(c)

	companies, total, err := h.service.List(c.Request.Context(), params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := params.ToMeta(total)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    companies,
		"meta":    meta,
	})
}
