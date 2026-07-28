package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

func (h *Handler) CreateCandidate(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.CandidateSvc.Create(c.Request.Context(), companyID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) ListCandidates(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.CandidateSvc.List(c.Request.Context(), companyID, c.Request.URL.Query())
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) GetCandidate(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.CandidateSvc.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Candidate not found")
		return
	}
	response.Success(c, data)
}

func (h *Handler) UpdateCandidate(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.CandidateSvc.Update(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) BlacklistCandidate(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.CandidateSvc.Blacklist(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "candidate blacklisted"})
}

func (h *Handler) UnblacklistCandidate(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.CandidateSvc.Unblacklist(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "candidate unblacklisted"})
}

func (h *Handler) SearchCandidates(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	skills := c.Query("skills")
	data, err := h.CandidateSvc.Search(c.Request.Context(), companyID, skills)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) ListCandidateEducation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.CandidateSvc.ListEducation(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) AddCandidateEducation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.CandidateSvc.AddEducation(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) UpdateCandidateEducation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.CandidateSvc.UpdateEducation(c.Request.Context(), companyID, c.Param("eduId"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) DeleteCandidateEducation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.CandidateSvc.DeleteEducation(c.Request.Context(), companyID, c.Param("eduId")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "education deleted"})
}

func (h *Handler) ListCandidateExperience(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.CandidateSvc.ListExperience(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) AddCandidateExperience(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.CandidateSvc.AddExperience(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) UpdateCandidateExperience(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.CandidateSvc.UpdateExperience(c.Request.Context(), companyID, c.Param("expId"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) DeleteCandidateExperience(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.CandidateSvc.DeleteExperience(c.Request.Context(), companyID, c.Param("expId")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "experience deleted"})
}

func (h *Handler) ListCandidateSkills(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.CandidateSvc.ListSkills(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) AddCandidateSkill(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.CandidateSvc.AddSkill(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) UpdateCandidateSkill(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.CandidateSvc.UpdateSkill(c.Request.Context(), companyID, c.Param("skillId"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) DeleteCandidateSkill(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.CandidateSvc.DeleteSkill(c.Request.Context(), companyID, c.Param("skillId")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "skill deleted"})
}

func (h *Handler) ListCandidateCertifications(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.CandidateSvc.ListCertifications(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) AddCandidateCertification(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.CandidateSvc.AddCertification(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) DeleteCandidateCertification(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.CandidateSvc.DeleteCertification(c.Request.Context(), companyID, c.Param("certId")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "certification deleted"})
}

func (h *Handler) ListCandidateLanguages(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.CandidateSvc.ListLanguages(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) AddCandidateLanguage(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.CandidateSvc.AddLanguage(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) UpdateCandidateLanguage(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.CandidateSvc.UpdateLanguage(c.Request.Context(), companyID, c.Param("langId"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) DeleteCandidateLanguage(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.CandidateSvc.DeleteLanguage(c.Request.Context(), companyID, c.Param("langId")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "language deleted"})
}

func (h *Handler) ListCandidateDocuments(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.CandidateSvc.ListDocuments(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) AddCandidateDocument(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.CandidateSvc.AddDocument(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}
