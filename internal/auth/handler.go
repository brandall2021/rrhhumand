package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/pkg/response"
)

type Handler struct {
	service *AuthService
}

func NewAuthHandler(service *AuthService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	companyID := c.Query("company_id")
	if companyID == "" && req.CompanySlug != "" {
		company, err := h.service.companyResolver.FindBySlug(c.Request.Context(), req.CompanySlug)
		if err != nil {
			response.BadRequest(c, "Invalid company")
			return
		}
		companyID = company.ID
	}

	if companyID == "" {
		response.BadRequest(c, "company_id query parameter or company_slug in body required")
		return
	}

	resp, err := h.service.Login(c.Request.Context(), &req, companyID)
	if err != nil {
		switch err {
		case ErrInvalidCredentials:
			response.Unauthorized(c, "Invalid email or password")
		default:
			response.InternalError(c, "Login failed")
		}
		return
	}

	response.Success(c, resp)
}

func (h *Handler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	companyID := c.GetString("company_id")
	if companyID == "" {
		companyID = c.Query("company_id")
	}

	tokenPair, err := h.service.RefreshToken(c.Request.Context(), req.RefreshToken, companyID)
	if err != nil {
		switch err {
		case ErrInvalidToken, ErrTokenRevoked:
			response.Unauthorized(c, "Invalid refresh token")
		case ErrTokenExpired:
			response.Unauthorized(c, "Refresh token expired")
		default:
			response.InternalError(c, "Token refresh failed")
		}
		return
	}

	response.Success(c, tokenPair)
}

func (h *Handler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.service.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		response.InternalError(c, "Logout failed")
		return
	}

	response.NoContent(c)
}

func (h *Handler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "Not authenticated")
		return
	}

	_ = http.MethodGet

	user, err := h.service.userRepo.FindByID(c.Request.Context(), userID.(string))
	if err != nil {
		response.NotFound(c, "User not found")
		return
	}

	companyID, _ := c.Get("company_id")
	roles, _ := h.service.userRepo.GetRolesByCompany(c.Request.Context(), user.ID, companyID.(string))

	response.Success(c, gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"company_id": companyID,
		"roles":      roles,
	})
}
