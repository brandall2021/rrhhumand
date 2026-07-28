package training

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/tenant"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) userID(c *gin.Context) string {
	return c.GetString("user_id")
}

func (h *Handler) companyID(c *gin.Context) string {
	return tenant.GetCompanyID(c)
}

func (h *Handler) bindJSON(c *gin.Context, obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}
	return true
}

func (h *Handler) bindQuery(c *gin.Context, obj any) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}
	return true
}

func queryStr(c *gin.Context, key string) *string {
	v := c.Query(key)
	if v == "" {
		return nil
	}
	return &v
}

func queryBool(c *gin.Context, key string) *bool {
	v := c.Query(key)
	if v == "" {
		return nil
	}
	b := v == "true" || v == "1"
	return &b
}

func queryInt(c *gin.Context, key string, defaultVal int) int {
	v := c.Query(key)
	if v == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}
	return i
}

// ---------------------------------------------------------------------------
// Categories
// ---------------------------------------------------------------------------

func (h *Handler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if !h.bindJSON(c, &req) {
		return
	}
	cat, err := h.svc.CreateCategory(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cat)
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	var req UpdateCategoryRequest
	if !h.bindJSON(c, &req) {
		return
	}
	cat, err := h.svc.UpdateCategory(c.Request.Context(), h.companyID(c), c.Param("id"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cat)
}

func (h *Handler) GetCategory(c *gin.Context) {
	cat, err := h.svc.GetCategory(c.Request.Context(), h.companyID(c), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}
	c.JSON(http.StatusOK, cat)
}

func (h *Handler) ListCategories(c *gin.Context) {
	cats, err := h.svc.ListCategories(c.Request.Context(), h.companyID(c), queryStr(c, "parent_id"), queryBool(c, "active"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cats)
}

// ---------------------------------------------------------------------------
// Courses
// ---------------------------------------------------------------------------

func (h *Handler) CreateCourse(c *gin.Context) {
	var req CreateCourseRequest
	if !h.bindJSON(c, &req) {
		return
	}
	course, err := h.svc.CreateCourse(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, course)
}

func (h *Handler) UpdateCourse(c *gin.Context) {
	var req UpdateCourseRequest
	if !h.bindJSON(c, &req) {
		return
	}
	course, err := h.svc.UpdateCourse(c.Request.Context(), h.companyID(c), c.Param("id"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, course)
}

func (h *Handler) GetCourse(c *gin.Context) {
	course, err := h.svc.GetCourse(c.Request.Context(), h.companyID(c), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
		return
	}
	c.JSON(http.StatusOK, course)
}

func (h *Handler) ListCourses(c *gin.Context) {
	filter := CourseFilter{
		CategoryID: queryStr(c, "category_id"),
		Status:     queryStr(c, "status"),
		Modality:   queryStr(c, "modality"),
		Mandatory:  queryBool(c, "mandatory"),
		Difficulty: queryStr(c, "difficulty"),
		Search:     queryStr(c, "search"),
		SortBy:     c.Query("sort_by"),
		SortDesc:   c.Query("sort_desc") == "true",
		Limit:      queryInt(c, "limit", 20),
		Offset:     queryInt(c, "offset", 0),
	}
	courses, total, err := h.svc.ListCourses(c.Request.Context(), h.companyID(c), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": courses, "total": total})
}

func (h *Handler) PublishCourse(c *gin.Context) {
	if err := h.svc.PublishCourse(c.Request.Context(), h.companyID(c), c.Param("id"), h.userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "course published"})
}

func (h *Handler) GetCourseWithDetails(c *gin.Context) {
	details, err := h.svc.GetCourseWithDetails(c.Request.Context(), h.companyID(c), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
		return
	}
	c.JSON(http.StatusOK, details)
}

// ---------------------------------------------------------------------------
// Versions
// ---------------------------------------------------------------------------

func (h *Handler) CreateVersion(c *gin.Context) {
	var req CreateVersionRequest
	if !h.bindJSON(c, &req) {
		return
	}
	v, err := h.svc.CreateVersion(c.Request.Context(), c.Param("course_id"), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, v)
}

func (h *Handler) ListVersions(c *gin.Context) {
	versions, err := h.svc.ListVersions(c.Request.Context(), c.Param("course_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, versions)
}

// ---------------------------------------------------------------------------
// Contents
// ---------------------------------------------------------------------------

func (h *Handler) CreateContent(c *gin.Context) {
	var req CreateContentRequest
	if !h.bindJSON(c, &req) {
		return
	}
	content, err := h.svc.CreateContent(c.Request.Context(), c.Param("version_id"), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, content)
}

func (h *Handler) UpdateContent(c *gin.Context) {
	var req UpdateContentRequest
	if !h.bindJSON(c, &req) {
		return
	}
	content, err := h.svc.UpdateContent(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, content)
}

func (h *Handler) ListContents(c *gin.Context) {
	contents, err := h.svc.ListContents(c.Request.Context(), c.Param("version_id"), queryBool(c, "published") != nil && *queryBool(c, "published"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, contents)
}

func (h *Handler) ReorderContents(c *gin.Context) {
	var req struct {
		ContentIDs []string `json:"content_ids" binding:"required"`
	}
	if !h.bindJSON(c, &req) {
		return
	}
	if err := h.svc.ReorderContents(c.Request.Context(), c.Param("version_id"), req.ContentIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "contents reordered"})
}

// ---------------------------------------------------------------------------
// Offerings
// ---------------------------------------------------------------------------

func (h *Handler) CreateOffering(c *gin.Context) {
	var req CreateOfferingRequest
	if !h.bindJSON(c, &req) {
		return
	}
	offering, err := h.svc.CreateOffering(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, offering)
}

func (h *Handler) UpdateOffering(c *gin.Context) {
	var req UpdateOfferingRequest
	if !h.bindJSON(c, &req) {
		return
	}
	offering, err := h.svc.UpdateOffering(c.Request.Context(), h.companyID(c), c.Param("id"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, offering)
}

func (h *Handler) GetOffering(c *gin.Context) {
	offering, err := h.svc.GetOffering(c.Request.Context(), h.companyID(c), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "offering not found"})
		return
	}
	c.JSON(http.StatusOK, offering)
}

func (h *Handler) ListOfferings(c *gin.Context) {
	filter := OfferingFilter{
		CourseID:     queryStr(c, "course_id"),
		Status:       queryStr(c, "status"),
		InstructorID: queryStr(c, "instructor_id"),
		FromDate:     queryStr(c, "from_date"),
		ToDate:       queryStr(c, "to_date"),
		Limit:        queryInt(c, "limit", 20),
		Offset:       queryInt(c, "offset", 0),
	}
	offerings, total, err := h.svc.ListOfferings(c.Request.Context(), h.companyID(c), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": offerings, "total": total})
}

// ---------------------------------------------------------------------------
// Sessions
// ---------------------------------------------------------------------------

func (h *Handler) CreateSession(c *gin.Context) {
	var req CreateSessionRequest
	if !h.bindJSON(c, &req) {
		return
	}
	sess, err := h.svc.CreateSession(c.Request.Context(), c.Param("offering_id"), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sess)
}

func (h *Handler) ListSessions(c *gin.Context) {
	sessions, err := h.svc.ListSessions(c.Request.Context(), c.Param("offering_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sessions)
}

// ---------------------------------------------------------------------------
// Enrollments
// ---------------------------------------------------------------------------

func (h *Handler) Enroll(c *gin.Context) {
	var req EnrollRequest
	if !h.bindJSON(c, &req) {
		return
	}
	enrollment, err := h.svc.Enroll(c.Request.Context(), h.companyID(c), c.Param("offering_id"), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, enrollment)
}

func (h *Handler) CompleteEnrollment(c *gin.Context) {
	if err := h.svc.CompleteEnrollment(c.Request.Context(), h.companyID(c), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "enrollment completed"})
}

func (h *Handler) GetEnrollment(c *gin.Context) {
	enrollment, err := h.svc.GetEnrollment(c.Request.Context(), h.companyID(c), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "enrollment not found"})
		return
	}
	c.JSON(http.StatusOK, enrollment)
}

func (h *Handler) ListEnrollments(c *gin.Context) {
	filter := EnrollmentFilter{
		OfferingID: queryStr(c, "offering_id"),
		EmployeeID: queryStr(c, "employee_id"),
		Status:     queryStr(c, "status"),
		Limit:      queryInt(c, "limit", 20),
		Offset:     queryInt(c, "offset", 0),
	}
	enrollments, total, err := h.svc.ListEnrollments(c.Request.Context(), h.companyID(c), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": enrollments, "total": total})
}

// ---------------------------------------------------------------------------
// Progress
// ---------------------------------------------------------------------------

func (h *Handler) UpdateProgress(c *gin.Context) {
	var req UpdateProgressRequest
	if !h.bindJSON(c, &req) {
		return
	}
	if err := h.svc.UpdateProgress(c.Request.Context(), c.Param("enrollment_id"), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "progress updated"})
}

func (h *Handler) GetProgress(c *gin.Context) {
	progress, err := h.svc.GetProgress(c.Request.Context(), c.Param("enrollment_id"), c.Param("content_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "progress not found"})
		return
	}
	c.JSON(http.StatusOK, progress)
}

func (h *Handler) ListProgress(c *gin.Context) {
	progress, err := h.svc.ListProgress(c.Request.Context(), c.Param("enrollment_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, progress)
}

// ---------------------------------------------------------------------------
// Assignments
// ---------------------------------------------------------------------------

func (h *Handler) CreateAssignment(c *gin.Context) {
	var req CreateAssignmentRequest
	if !h.bindJSON(c, &req) {
		return
	}
	assignment, err := h.svc.CreateAssignment(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, assignment)
}

// ---------------------------------------------------------------------------
// Assignment Rules
// ---------------------------------------------------------------------------

func (h *Handler) CreateAssignmentRule(c *gin.Context) {
	var req CreateAssignmentRuleRequest
	if !h.bindJSON(c, &req) {
		return
	}
	rule, err := h.svc.CreateAssignmentRule(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, rule)
}

func (h *Handler) ListAssignmentRules(c *gin.Context) {
	rules, err := h.svc.ListAssignmentRules(c.Request.Context(), h.companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rules)
}

// ---------------------------------------------------------------------------
// Assessments
// ---------------------------------------------------------------------------

func (h *Handler) CreateAssessment(c *gin.Context) {
	var req CreateAssessmentRequest
	if !h.bindJSON(c, &req) {
		return
	}
	assessment, err := h.svc.CreateAssessment(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, assessment)
}

func (h *Handler) UpdateAssessment(c *gin.Context) {
	var req UpdateAssessmentRequest
	if !h.bindJSON(c, &req) {
		return
	}
	assessment, err := h.svc.UpdateAssessment(c.Request.Context(), h.companyID(c), c.Param("id"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assessment)
}

func (h *Handler) GetAssessment(c *gin.Context) {
	assessment, err := h.svc.GetAssessment(c.Request.Context(), h.companyID(c), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "assessment not found"})
		return
	}
	c.JSON(http.StatusOK, assessment)
}

func (h *Handler) ListAssessments(c *gin.Context) {
	assessments, err := h.svc.ListAssessments(c.Request.Context(), h.companyID(c), c.Param("course_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assessments)
}

// ---------------------------------------------------------------------------
// Questions
// ---------------------------------------------------------------------------

func (h *Handler) AddQuestion(c *gin.Context) {
	var req CreateQuestionRequest
	if !h.bindJSON(c, &req) {
		return
	}
	question, err := h.svc.AddQuestion(c.Request.Context(), c.Param("assessment_id"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, question)
}

func (h *Handler) GetQuestions(c *gin.Context) {
	questions, err := h.svc.GetQuestions(c.Request.Context(), c.Param("assessment_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, questions)
}

// ---------------------------------------------------------------------------
// Attempts
// ---------------------------------------------------------------------------

func (h *Handler) StartAttempt(c *gin.Context) {
	attempt, err := h.svc.StartAttempt(c.Request.Context(), c.Param("enrollment_id"), c.Param("assessment_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, attempt)
}

func (h *Handler) SubmitAttempt(c *gin.Context) {
	var req SubmitAttemptRequest
	if !h.bindJSON(c, &req) {
		return
	}
	result, err := h.svc.SubmitAttempt(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetAttempt(c *gin.Context) {
	attempt, err := h.svc.GetAttempt(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "attempt not found"})
		return
	}
	c.JSON(http.StatusOK, attempt)
}

func (h *Handler) ListAttempts(c *gin.Context) {
	attempts, err := h.svc.ListAttempts(c.Request.Context(), c.Param("enrollment_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, attempts)
}

// ---------------------------------------------------------------------------
// Instructors
// ---------------------------------------------------------------------------

func (h *Handler) CreateInstructor(c *gin.Context) {
	var req CreateInstructorRequest
	if !h.bindJSON(c, &req) {
		return
	}
	instructor, err := h.svc.CreateInstructor(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, instructor)
}

func (h *Handler) GetInstructor(c *gin.Context) {
	instructor, err := h.svc.GetInstructor(c.Request.Context(), h.companyID(c), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "instructor not found"})
		return
	}
	c.JSON(http.StatusOK, instructor)
}

func (h *Handler) ListInstructors(c *gin.Context) {
	instructors, err := h.svc.ListInstructors(c.Request.Context(), h.companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, instructors)
}

// ---------------------------------------------------------------------------
// Providers
// ---------------------------------------------------------------------------

func (h *Handler) CreateProvider(c *gin.Context) {
	var req CreateProviderRequest
	if !h.bindJSON(c, &req) {
		return
	}
	provider, err := h.svc.CreateProvider(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, provider)
}

func (h *Handler) ListProviders(c *gin.Context) {
	providers, err := h.svc.ListProviders(c.Request.Context(), h.companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, providers)
}

// ---------------------------------------------------------------------------
// Competencies
// ---------------------------------------------------------------------------

func (h *Handler) CreateCompetency(c *gin.Context) {
	var req CreateCompetencyRequest
	if !h.bindJSON(c, &req) {
		return
	}
	competency, err := h.svc.CreateCompetency(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, competency)
}

func (h *Handler) GetCompetency(c *gin.Context) {
	competency, err := h.svc.GetCompetency(c.Request.Context(), h.companyID(c), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "competency not found"})
		return
	}
	c.JSON(http.StatusOK, competency)
}

func (h *Handler) ListCompetencies(c *gin.Context) {
	competencies, err := h.svc.ListCompetencies(c.Request.Context(), h.companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, competencies)
}

func (h *Handler) AssignCompetency(c *gin.Context) {
	var req AssignCompetencyRequest
	if !h.bindJSON(c, &req) {
		return
	}
	if err := h.svc.AssignCompetency(c.Request.Context(), h.companyID(c), c.Param("employee_id"), c.Param("competency_id"), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "competency assigned"})
}

func (h *Handler) GetEmployeeCompetencies(c *gin.Context) {
	competencies, err := h.svc.GetEmployeeCompetencies(c.Request.Context(), h.companyID(c), c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, competencies)
}

// ---------------------------------------------------------------------------
// Course-Competency mapping
// ---------------------------------------------------------------------------

func (h *Handler) AddCourseCompetency(c *gin.Context) {
	var req struct {
		CompetencyID  string `json:"competency_id" binding:"required"`
		ExpectedLevel int    `json:"expected_level"`
		Weight        int    `json:"weight"`
	}
	if !h.bindJSON(c, &req) {
		return
	}
	if err := h.svc.AddCourseCompetency(c.Request.Context(), c.Param("course_id"), req.CompetencyID, req.ExpectedLevel, req.Weight); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "competency added to course"})
}

func (h *Handler) ListCourseCompetencies(c *gin.Context) {
	competencies, err := h.svc.ListCourseCompetencies(c.Request.Context(), c.Param("course_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, competencies)
}

// ---------------------------------------------------------------------------
// Training Needs
// ---------------------------------------------------------------------------

func (h *Handler) CreateTrainingNeed(c *gin.Context) {
	var req CreateTrainingNeedRequest
	if !h.bindJSON(c, &req) {
		return
	}
	need, err := h.svc.CreateTrainingNeed(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, need)
}

func (h *Handler) ListTrainingNeeds(c *gin.Context) {
	needs, err := h.svc.ListTrainingNeeds(c.Request.Context(), h.companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, needs)
}

// ---------------------------------------------------------------------------
// Training Plans
// ---------------------------------------------------------------------------

func (h *Handler) CreatePlan(c *gin.Context) {
	var req CreateTrainingPlanRequest
	if !h.bindJSON(c, &req) {
		return
	}
	plan, err := h.svc.CreatePlan(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, plan)
}

func (h *Handler) ListPlans(c *gin.Context) {
	plans, err := h.svc.ListPlans(c.Request.Context(), h.companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plans)
}

// ---------------------------------------------------------------------------
// Learning Paths
// ---------------------------------------------------------------------------

func (h *Handler) CreateLearningPath(c *gin.Context) {
	var req CreateLearningPathRequest
	if !h.bindJSON(c, &req) {
		return
	}
	path, err := h.svc.CreateLearningPath(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, path)
}

func (h *Handler) ListLearningPaths(c *gin.Context) {
	paths, err := h.svc.ListLearningPaths(c.Request.Context(), h.companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, paths)
}

// ---------------------------------------------------------------------------
// Feedback
// ---------------------------------------------------------------------------

func (h *Handler) CreateFeedback(c *gin.Context) {
	var req CreateFeedbackRequest
	if !h.bindJSON(c, &req) {
		return
	}
	feedback, err := h.svc.CreateFeedback(c.Request.Context(), c.Param("enrollment_id"), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, feedback)
}

func (h *Handler) GetFeedbackByEnrollment(c *gin.Context) {
	feedback, err := h.svc.GetFeedbackByEnrollment(c.Request.Context(), c.Param("enrollment_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "feedback not found"})
		return
	}
	c.JSON(http.StatusOK, feedback)
}

// ---------------------------------------------------------------------------
// Attendance
// ---------------------------------------------------------------------------

func (h *Handler) CreateAttendance(c *gin.Context) {
	var req CreateAttendanceRequest
	if !h.bindJSON(c, &req) {
		return
	}
	att, err := h.svc.CreateAttendance(c.Request.Context(), c.Param("enrollment_id"), c.Param("session_id"), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, att)
}

func (h *Handler) GetAttendance(c *gin.Context) {
	att, err := h.svc.GetAttendance(c.Request.Context(), c.Param("enrollment_id"), c.Param("session_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "attendance not found"})
		return
	}
	c.JSON(http.StatusOK, att)
}

func (h *Handler) ListAttendance(c *gin.Context) {
	att, err := h.svc.ListAttendance(c.Request.Context(), c.Param("enrollment_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, att)
}

// ---------------------------------------------------------------------------
// Certificates
// ---------------------------------------------------------------------------

func (h *Handler) ListCertificates(c *gin.Context) {
	certs, err := h.svc.ListCertificates(c.Request.Context(), h.companyID(c), c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, certs)
}

// ---------------------------------------------------------------------------
// Dashboard & Stats
// ---------------------------------------------------------------------------

func (h *Handler) DashboardStats(c *gin.Context) {
	stats, err := h.svc.GetDashboardStats(c.Request.Context(), h.companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *Handler) EmployeeStats(c *gin.Context) {
	stats, err := h.svc.GetEmployeeStats(c.Request.Context(), h.companyID(c), c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// ---------------------------------------------------------------------------
// AI Recommendations
// ---------------------------------------------------------------------------

func (h *Handler) GenerateRecommendations(c *gin.Context) {
	var req AIRecommendationRequest
	if !h.bindJSON(c, &req) {
		return
	}
	recs, err := h.svc.GenerateRecommendations(c.Request.Context(), h.companyID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recs)
}

func (h *Handler) GetRecommendations(c *gin.Context) {
	recs, err := h.svc.GetRecommendations(c.Request.Context(), c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recs)
}


