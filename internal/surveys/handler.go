package surveys

import (
	"encoding/csv"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/models"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

type SurveyHandler struct {
	service      *SurveyService
	statsService *SurveyStatsService
}

func NewSurveyHandler(service *SurveyService, statsService *SurveyStatsService) *SurveyHandler {
	return &SurveyHandler{service: service, statsService: statsService}
}

func (h *SurveyHandler) CreateSurvey(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)

	var req CreateSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	survey, err := h.service.CreateSurvey(c.Request.Context(), companyID, userID, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, survey)
}

func (h *SurveyHandler) ListSurveys(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	params := models.NewPaginationParams(c)

	filters := SurveyFilters{
		Status: c.Query("status"),
		Type:   c.Query("type"),
		Search: c.Query("search"),
	}

	surveys, total, err := h.service.ListSurveys(c.Request.Context(), companyID, filters, params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    surveys,
		"meta":    params.ToMeta(total),
	})
}

func (h *SurveyHandler) GetSurvey(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	survey, err := h.service.GetSurveyByID(c.Request.Context(), id, companyID)
	if err != nil {
		if err.Error() == "survey not found" {
			response.NotFound(c, "Survey not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, survey)
}

func (h *SurveyHandler) UpdateSurvey(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	var req UpdateSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	survey, err := h.service.UpdateSurvey(c.Request.Context(), id, companyID, &req)
	if err != nil {
		switch err.Error() {
		case "survey not found":
			response.NotFound(c, "Survey not found")
		case "survey not editable":
			response.BadRequest(c, "Survey can only be edited in DRAFT status")
		default:
			response.InternalError(c, err.Error())
		}
		return
	}

	response.Success(c, survey)
}

func (h *SurveyHandler) DeleteSurvey(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	if err := h.service.DeleteSurvey(c.Request.Context(), id, companyID); err != nil {
		switch err.Error() {
		case "survey not found":
			response.NotFound(c, "Survey not found")
		case "only draft surveys can be deleted":
			response.BadRequest(c, "Only draft surveys can be deleted")
		default:
			response.InternalError(c, err.Error())
		}
		return
	}

	response.NoContent(c)
}

func (h *SurveyHandler) PublishSurvey(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	if err := h.service.PublishSurvey(c.Request.Context(), id, companyID); err != nil {
		switch err.Error() {
		case "survey not found":
			response.NotFound(c, "Survey not found")
		case "survey invalid transition":
			response.BadRequest(c, "Survey must be in DRAFT status to publish")
		case "survey no questions":
			response.BadRequest(c, "Survey must have at least one question to publish")
		default:
			response.InternalError(c, err.Error())
		}
		return
	}

	response.Success(c, gin.H{"message": "survey published"})
}

func (h *SurveyHandler) CloseSurvey(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	if err := h.service.CloseSurvey(c.Request.Context(), id, companyID); err != nil {
		switch err.Error() {
		case "survey not found":
			response.NotFound(c, "Survey not found")
		case "survey invalid transition":
			response.BadRequest(c, "Survey must be in PUBLISHED status to close")
		default:
			response.InternalError(c, err.Error())
		}
		return
	}

	response.Success(c, gin.H{"message": "survey closed"})
}

func (h *SurveyHandler) ArchiveSurvey(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	if err := h.service.ArchiveSurvey(c.Request.Context(), id, companyID); err != nil {
		switch err.Error() {
		case "survey not found":
			response.NotFound(c, "Survey not found")
		case "survey invalid transition":
			response.BadRequest(c, "Survey must be in CLOSED status to archive")
		default:
			response.InternalError(c, err.Error())
		}
		return
	}

	response.Success(c, gin.H{"message": "survey archived"})
}

func (h *SurveyHandler) AddQuestion(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	surveyID := c.Param("id")

	var req CreateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	q, err := h.service.AddQuestion(c.Request.Context(), surveyID, companyID, &req)
	if err != nil {
		switch err.Error() {
		case "survey not found":
			response.NotFound(c, "Survey not found")
		case "survey not editable":
			response.BadRequest(c, "Survey can only be edited in DRAFT status")
		default:
			response.InternalError(c, err.Error())
		}
		return
	}

	response.Created(c, q)
}

func (h *SurveyHandler) UpdateQuestion(c *gin.Context) {
	questionID := c.Param("id")

	var req UpdateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	q, err := h.service.UpdateQuestion(c.Request.Context(), questionID, &req)
	if err != nil {
		switch err.Error() {
		case "question not found":
			response.NotFound(c, "Question not found")
		case "survey not editable":
			response.BadRequest(c, "Survey can only be edited in DRAFT status")
		case "survey question has responses":
			response.BadRequest(c, "Cannot modify question that already has responses")
		default:
			response.InternalError(c, err.Error())
		}
		return
	}

	response.Success(c, q)
}

func (h *SurveyHandler) DeleteQuestion(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	questionID := c.Param("id")

	if err := h.service.DeleteQuestion(c.Request.Context(), questionID, companyID); err != nil {
		switch err.Error() {
		case "question not found":
			response.NotFound(c, "Question not found")
		case "survey not found":
			response.NotFound(c, "Survey not found")
		case "survey not editable":
			response.BadRequest(c, "Survey can only be edited in DRAFT status")
		default:
			response.InternalError(c, err.Error())
		}
		return
	}

	response.NoContent(c)
}

func (h *SurveyHandler) AddOption(c *gin.Context) {
	questionID := c.Param("id")

	var req CreateOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	o, err := h.service.AddOption(c.Request.Context(), questionID, &req)
	if err != nil {
		switch err.Error() {
		case "question not found":
			response.NotFound(c, "Question not found")
		case "survey not editable":
			response.BadRequest(c, "Survey can only be edited in DRAFT status")
		default:
			response.InternalError(c, err.Error())
		}
		return
	}

	response.Created(c, o)
}

func (h *SurveyHandler) DeleteOption(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	optionID := c.Param("id")

	if err := h.service.DeleteOption(c.Request.Context(), optionID, companyID); err != nil {
		switch err.Error() {
		case "option not found":
			response.NotFound(c, "Option not found")
		case "survey not found":
			response.NotFound(c, "Survey not found")
		case "survey not editable":
			response.BadRequest(c, "Survey can only be edited in DRAFT status")
		default:
			response.InternalError(c, err.Error())
		}
		return
	}

	response.NoContent(c)
}

func (h *SurveyHandler) SetTargets(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	surveyID := c.Param("id")

	var req SetTargetsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.service.SetTargets(c.Request.Context(), surveyID, companyID, &req); err != nil {
		switch err.Error() {
		case "survey not found":
			response.NotFound(c, "Survey not found")
		case "survey not editable":
			response.BadRequest(c, "Survey can only be edited in DRAFT status")
		default:
			response.InternalError(c, err.Error())
		}
		return
	}

	response.Success(c, gin.H{"message": "targets updated"})
}

func (h *SurveyHandler) ListTargets(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	surveyID := c.Param("id")

	targets, err := h.service.ListTargets(c.Request.Context(), surveyID, companyID)
	if err != nil {
		if err.Error() == "survey not found" {
			response.NotFound(c, "Survey not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, targets)
}

func (h *SurveyHandler) ListAvailableSurveys(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	employeeID := c.GetString("employee_id")

	surveys, err := h.service.ListAvailableSurveys(c.Request.Context(), companyID, employeeID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, surveys)
}

func (h *SurveyHandler) RespondSurvey(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	employeeID := c.GetString("employee_id")
	surveyID := c.Param("id")

	var req RespondSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	if err := h.service.RespondSurvey(c.Request.Context(), surveyID, companyID, employeeID, &req, ipAddress, userAgent); err != nil {
		switch err.Error() {
		case "survey not found":
			response.NotFound(c, "Survey not found")
		case "survey closed":
			response.BadRequest(c, "Survey is not accepting responses")
		case "survey not active":
			response.BadRequest(c, "Survey is not within its active period")
		case "employee not targeted":
			response.Forbidden(c, "You are not targeted for this survey")
		case "survey already answered":
			response.Conflict(c, "You have already answered this survey")
		default:
			response.InternalError(c, err.Error())
		}
		return
	}

	response.Success(c, gin.H{"message": "response recorded"})
}

func (h *SurveyHandler) GetResults(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	stats, err := h.service.GetResults(c.Request.Context(), id, companyID, h.statsService)
	if err != nil {
		if err.Error() == "survey not found" {
			response.NotFound(c, "Survey not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, stats)
}

func (h *SurveyHandler) ExportSurvey(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")
	format := c.DefaultQuery("format", "csv")

	if format != "csv" {
		response.BadRequest(c, "Only CSV format is currently supported")
		return
	}

	rows, err := h.service.ExportCSV(c.Request.Context(), id, companyID)
	if err != nil {
		if err.Error() == "survey not found" {
			response.NotFound(c, "Survey not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=survey_export.csv")

	var sb strings.Builder
	writer := csv.NewWriter(&sb)
	for _, row := range rows {
		_ = writer.Write(row)
	}
	writer.Flush()

	c.String(http.StatusOK, sb.String())
}

func (h *SurveyHandler) ListQuestions(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	surveyID := c.Param("id")

	_, err := h.service.GetSurveyByID(c.Request.Context(), surveyID, companyID)
	if err != nil {
		if err.Error() == "survey not found" {
			response.NotFound(c, "Survey not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	questions, err := h.service.repo.ListQuestionsBySurveyID(c.Request.Context(), surveyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	for i := range questions {
		options, _ := h.service.repo.ListOptionsByQuestionID(c.Request.Context(), questions[i].ID)
		questions[i].Options = options
	}

	response.Success(c, questions)
}

func init() {
	_ = strconv.Itoa
}
