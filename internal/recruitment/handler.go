package recruitment

import (
	"net/http"

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

// Requisitions
func (h *Handler) CreateRequisition(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req CreateRequisitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	rec, err := h.service.CreateRequisition(c.Request.Context(), companyID, userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, rec)
}

func (h *Handler) GetRequisition(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	rec, err := h.service.GetRequisition(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Requisition not found")
		return
	}
	response.Success(c, rec)
}

func (h *Handler) ListRequisitions(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	filters := RecruitmentFilters{Status: c.Query("status")}
	recs, err := h.service.ListRequisitions(c.Request.Context(), companyID, filters)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, recs)
}

func (h *Handler) UpdateRequisition(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req UpdateRequisitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	rec, err := h.service.UpdateRequisition(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, rec)
}

func (h *Handler) SubmitRequisition(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.SubmitRequisition(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "requisition submitted for approval"})
}

func (h *Handler) ApproveRequisition(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.ApproveRequisition(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "requisition approved"})
}

func (h *Handler) OpenRequisition(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.OpenRequisition(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "requisition opened"})
}

func (h *Handler) CloseRequisition(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.CloseRequisition(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "requisition closed"})
}

// Postings
func (h *Handler) CreatePosting(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreatePostingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	p, err := h.service.CreatePosting(c.Request.Context(), companyID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, p)
}

func (h *Handler) GetPosting(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	p, err := h.service.GetPosting(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Posting not found")
		return
	}
	response.Success(c, p)
}

func (h *Handler) ListPostings(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	filters := RecruitmentFilters{Status: c.Query("status")}
	postings, err := h.service.ListPostings(c.Request.Context(), companyID, filters)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, postings)
}

func (h *Handler) PublishPosting(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.PublishPosting(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "posting published"})
}

func (h *Handler) ClosePosting(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.ClosePosting(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "posting closed"})
}

// Candidates
func (h *Handler) CreateCandidate(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateCandidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	cand, err := h.service.CreateCandidate(c.Request.Context(), companyID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, cand)
}

func (h *Handler) GetCandidate(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	cand, err := h.service.GetCandidate(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Candidate not found")
		return
	}
	response.Success(c, cand)
}

func (h *Handler) ListCandidates(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	filters := RecruitmentFilters{Status: c.Query("status"), Source: c.Query("source")}
	cands, err := h.service.ListCandidates(c.Request.Context(), companyID, filters)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, cands)
}

func (h *Handler) UpdateCandidate(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req UpdateCandidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	cand, err := h.service.UpdateCandidate(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, cand)
}

// Applications
func (h *Handler) CreateApplication(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	app, err := h.service.CreateApplication(c.Request.Context(), companyID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, app)
}

func (h *Handler) GetApplication(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	app, err := h.service.GetApplication(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Application not found")
		return
	}
	response.Success(c, app)
}

func (h *Handler) ListApplications(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	filters := RecruitmentFilters{
		CandidateID: c.Query("candidate_id"),
		PostingID:   c.Query("posting_id"),
		Status:      c.Query("status"),
	}
	apps, err := h.service.ListApplications(c.Request.Context(), companyID, filters)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, apps)
}

func (h *Handler) MoveStage(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req MoveStageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	history, err := h.service.MoveStage(c.Request.Context(), companyID, c.Param("id"), userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, history)
}

func (h *Handler) RejectApplication(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req RejectApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	if err := h.service.RejectApplication(c.Request.Context(), companyID, c.Param("id"), &req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "application rejected"})
}

func (h *Handler) WithdrawApplication(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.WithdrawApplication(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "application withdrawn"})
}

func (h *Handler) GetStageHistory(c *gin.Context) {
	history, err := h.service.GetStageHistory(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, history)
}

// Interviews
func (h *Handler) CreateInterview(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	iv, err := h.service.CreateInterview(c.Request.Context(), companyID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, iv)
}

func (h *Handler) GetInterview(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	iv, err := h.service.GetInterview(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Interview not found")
		return
	}
	response.Success(c, iv)
}

func (h *Handler) ListInterviews(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	filters := RecruitmentFilters{
		ApplicationID: c.Query("application_id"),
		InterviewerID: c.Query("interviewer_id"),
		Status:        c.Query("status"),
	}
	ivs, err := h.service.ListInterviews(c.Request.Context(), companyID, filters)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, ivs)
}

func (h *Handler) UpdateInterview(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req UpdateInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	iv, err := h.service.UpdateInterview(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, iv)
}

func (h *Handler) CreateInterviewFeedback(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req CreateInterviewFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	fb, err := h.service.CreateInterviewFeedback(c.Request.Context(), companyID, c.Param("id"), userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, fb)
}

func (h *Handler) ListInterviewFeedback(c *gin.Context) {
	feedbacks, err := h.service.ListInterviewFeedback(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, feedbacks)
}

// Assessments
func (h *Handler) CreateAssessment(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateAssessmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	a, err := h.service.CreateAssessment(c.Request.Context(), companyID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, a)
}

func (h *Handler) ListAssessments(c *gin.Context) {
	assessments, err := h.service.ListAssessments(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, assessments)
}

// Offers
func (h *Handler) CreateOffer(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req CreateOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	o, err := h.service.CreateOffer(c.Request.Context(), companyID, userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, o)
}

func (h *Handler) GetOffer(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	o, err := h.service.GetOffer(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Offer not found")
		return
	}
	response.Success(c, o)
}

func (h *Handler) SendOffer(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.SendOffer(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "offer sent"})
}

func (h *Handler) AcceptOffer(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.AcceptOffer(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "offer accepted, candidate hired"})
}

func (h *Handler) RejectOffer(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.RejectOffer(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "offer rejected"})
}

// Referrals
func (h *Handler) CreateReferral(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req CreateReferralRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	ref, err := h.service.CreateReferral(c.Request.Context(), companyID, userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, ref)
}

func (h *Handler) ListReferrals(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	refs, err := h.service.ListReferrals(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, refs)
}

// Screening
func (h *Handler) CreateScreeningQuestion(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateScreeningQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	q, err := h.service.CreateScreeningQuestion(c.Request.Context(), companyID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, q)
}

func (h *Handler) ListScreeningQuestions(c *gin.Context) {
	questions, err := h.service.ListScreeningQuestions(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, questions)
}

// Hire (candidate → employee conversion)
func (h *Handler) HireCandidate(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	result, err := h.service.HireCandidate(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// Dashboard
func (h *Handler) GetDashboard(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	dash, err := h.service.GetDashboard(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    dash,
	})
}
