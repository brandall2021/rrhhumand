package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

type ProfileHandler struct {
	service *ProfileService
}

func NewProfileHandler(service *ProfileService) *ProfileHandler {
	return &ProfileHandler{service: service}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	companyID := tenant.GetCompanyID(c)

	profile, err := h.service.Get(c.Request.Context(), userID.(string), companyID)
	if err != nil {
		if err.Error() == "employee profile not found for this user in this company" {
			response.NotFound(c, "No employee profile found for your account in this company")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, profile)
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	companyID := tenant.GetCompanyID(c)

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	profile, err := h.service.Update(c.Request.Context(), userID.(string), companyID, &req)
	if err != nil {
		if err.Error() == "employee not found for this user" {
			response.NotFound(c, "No employee found for your account")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profile,
	})
}
