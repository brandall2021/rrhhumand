package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/recruitment/domain"
	"github.com/rrhhumand/api/internal/recruitment/repository"
)

type CreateCandidateReq struct {
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	Email            string     `json:"email"`
	Phone            *string    `json:"phone"`
	PhoneCountryCode *string    `json:"phone_country_code"`
	DocumentType     *string    `json:"document_type"`
	DocumentNumber   *string    `json:"document_number"`
	BirthDate        *time.Time `json:"birth_date"`
	Location         *string    `json:"location"`
	Nationality      *string    `json:"nationality"`
	Gender           *string    `json:"gender"`
	LinkedInURL      *string    `json:"linkedin_url"`
	PortfolioURL     *string    `json:"portfolio_url"`
	GithubURL        *string    `json:"github_url"`
	PersonalWebsite  *string    `json:"personal_website"`
	CurrentCompany   *string    `json:"current_company"`
	CurrentPosition  *string    `json:"current_position"`
	NoticePeriod     *int       `json:"notice_period"`
	SalaryExpectMin  *float64   `json:"salary_expectation_min"`
	SalaryExpectMax  *float64   `json:"salary_expectation_max"`
	SalaryCurrency   *string    `json:"salary_currency"`
	Availability     *string    `json:"availability"`
	Source           *string    `json:"source"`
	SourceDetail     *string    `json:"source_detail"`
	Tags             []string   `json:"tags"`
	Notes            *string    `json:"notes"`
}

type CandidateService struct {
	candidateRepo *repository.CandidateRepo
}

func NewCandidateService(candidateRepo *repository.CandidateRepo) *CandidateService {
	return &CandidateService{candidateRepo: candidateRepo}
}

func (s *CandidateService) Create(ctx context.Context, companyID string, req *CreateCandidateReq) (*domain.Candidate, error) {
	const op = "CreateCandidate"
	now := time.Now()
	c := &domain.Candidate{
		ID:               uuid.New().String(),
		CompanyID:        companyID,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Email:            req.Email,
		Phone:            req.Phone,
		PhoneCountryCode: req.PhoneCountryCode,
		DocumentType:     req.DocumentType,
		DocumentNumber:   req.DocumentNumber,
		BirthDate:        req.BirthDate,
		Location:         req.Location,
		Nationality:      req.Nationality,
		Gender:           req.Gender,
		LinkedInURL:      req.LinkedInURL,
		PortfolioURL:     req.PortfolioURL,
		GithubURL:        req.GithubURL,
		PersonalWebsite:  req.PersonalWebsite,
		CurrentCompany:   req.CurrentCompany,
		CurrentPosition:  req.CurrentPosition,
		NoticePeriod:     req.NoticePeriod,
		SalaryExpectMin:  req.SalaryExpectMin,
		SalaryExpectMax:  req.SalaryExpectMax,
		SalaryCurrency:   req.SalaryCurrency,
		Availability:     req.Availability,
		Source:           req.Source,
		SourceDetail:     req.SourceDetail,
		Tags:             req.Tags,
		Notes:            req.Notes,
		Status:           domain.CandStatusActive,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	result, err := s.candidateRepo.Create(ctx, companyID, c)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *CandidateService) GetByID(ctx context.Context, companyID, id string) (*domain.Candidate, error) {
	const op = "GetCandidate"
	return s.candidateRepo.GetByID(ctx, companyID, id)
}

func (s *CandidateService) List(ctx context.Context, companyID, status, source string) ([]domain.Candidate, error) {
	const op = "ListCandidates"
	return s.candidateRepo.List(ctx, companyID, status, source)
}

func (s *CandidateService) Update(ctx context.Context, companyID, id string, req *domain.Candidate) (*domain.Candidate, error) {
	const op = "UpdateCandidate"
	req.UpdatedAt = time.Now()
	result, err := s.candidateRepo.Update(ctx, companyID, id, req)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *CandidateService) Blacklist(ctx context.Context, companyID, id, reason string) error {
	const op = "BlacklistCandidate"
	c, err := s.candidateRepo.GetByID(ctx, companyID, id)
	if err != nil {
		return svcErr(op, err)
	}
	c.Status = domain.CandStatusBlacklisted
	c.Blacklisted = true
	c.BlacklistReason = &reason
	c.UpdatedAt = time.Now()
	_, err = s.candidateRepo.Update(ctx, companyID, id, c)
	return err
}

func (s *CandidateService) Unblacklist(ctx context.Context, companyID, id string) error {
	const op = "UnblacklistCandidate"
	c, err := s.candidateRepo.GetByID(ctx, companyID, id)
	if err != nil {
		return svcErr(op, err)
	}
	c.Status = domain.CandStatusActive
	c.Blacklisted = false
	c.BlacklistReason = nil
	c.UpdatedAt = time.Now()
	_, err = s.candidateRepo.Update(ctx, companyID, id, c)
	return err
}

func (s *CandidateService) AddEducation(ctx context.Context, companyID, candidateID string, edu domain.CandidateEducation) (*domain.CandidateEducation, error) {
	const op = "AddEducation"
	edu.ID = uuid.New().String()
	edu.CandidateID = candidateID
	edu.CreatedAt = time.Now()
	result, err := s.candidateRepo.AddEducation(ctx, &edu)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *CandidateService) UpdateEducation(ctx context.Context, companyID, candidateID string, edu domain.CandidateEducation) error {
	const op = "UpdateEducation"
	_, err := s.candidateRepo.UpdateEducation(ctx, edu.ID, &edu)
	return err
}

func (s *CandidateService) DeleteEducation(ctx context.Context, companyID, candidateID, educationID string) error {
	const op = "DeleteEducation"
	return s.candidateRepo.DeleteEducation(ctx, educationID)
}

func (s *CandidateService) ListEducation(ctx context.Context, companyID, candidateID string) ([]domain.CandidateEducation, error) {
	const op = "ListEducation"
	return s.candidateRepo.ListEducation(ctx, candidateID)
}

func (s *CandidateService) AddExperience(ctx context.Context, companyID, candidateID string, exp domain.CandidateExperience) (*domain.CandidateExperience, error) {
	const op = "AddExperience"
	exp.ID = uuid.New().String()
	exp.CandidateID = candidateID
	exp.CreatedAt = time.Now()
	result, err := s.candidateRepo.AddExperience(ctx, &exp)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *CandidateService) UpdateExperience(ctx context.Context, companyID, candidateID string, exp domain.CandidateExperience) error {
	const op = "UpdateExperience"
	_, err := s.candidateRepo.UpdateExperience(ctx, exp.ID, &exp)
	return err
}

func (s *CandidateService) DeleteExperience(ctx context.Context, companyID, candidateID, experienceID string) error {
	const op = "DeleteExperience"
	return s.candidateRepo.DeleteExperience(ctx, experienceID)
}

func (s *CandidateService) ListExperience(ctx context.Context, companyID, candidateID string) ([]domain.CandidateExperience, error) {
	const op = "ListExperience"
	return s.candidateRepo.ListExperience(ctx, candidateID)
}

func (s *CandidateService) AddSkill(ctx context.Context, companyID, candidateID string, skill domain.CandidateSkill) (*domain.CandidateSkill, error) {
	const op = "AddSkill"
	skill.ID = uuid.New().String()
	skill.CandidateID = candidateID
	skill.CreatedAt = time.Now()
	result, err := s.candidateRepo.AddSkill(ctx, &skill)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *CandidateService) UpdateSkill(ctx context.Context, companyID, candidateID string, skill domain.CandidateSkill) error {
	const op = "UpdateSkill"
	_, err := s.candidateRepo.UpdateSkill(ctx, skill.ID, &skill)
	return err
}

func (s *CandidateService) DeleteSkill(ctx context.Context, companyID, candidateID, skillID string) error {
	const op = "DeleteSkill"
	return s.candidateRepo.DeleteSkill(ctx, skillID)
}

func (s *CandidateService) ListSkills(ctx context.Context, companyID, candidateID string) ([]domain.CandidateSkill, error) {
	const op = "ListSkills"
	return s.candidateRepo.ListSkills(ctx, candidateID)
}

func (s *CandidateService) AddCertification(ctx context.Context, companyID, candidateID string, cert domain.CandidateCertification) (*domain.CandidateCertification, error) {
	const op = "AddCertification"
	cert.ID = uuid.New().String()
	cert.CandidateID = candidateID
	cert.CreatedAt = time.Now()
	result, err := s.candidateRepo.AddCertification(ctx, &cert)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *CandidateService) DeleteCertification(ctx context.Context, companyID, candidateID, certID string) error {
	const op = "DeleteCertification"
	return s.candidateRepo.DeleteCertification(ctx, certID)
}

func (s *CandidateService) ListCertifications(ctx context.Context, companyID, candidateID string) ([]domain.CandidateCertification, error) {
	const op = "ListCertifications"
	return s.candidateRepo.ListCertifications(ctx, candidateID)
}

func (s *CandidateService) AddLanguage(ctx context.Context, companyID, candidateID string, lang domain.CandidateLanguage) (*domain.CandidateLanguage, error) {
	const op = "AddLanguage"
	lang.ID = uuid.New().String()
	lang.CandidateID = candidateID
	lang.CreatedAt = time.Now()
	result, err := s.candidateRepo.AddLanguage(ctx, &lang)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *CandidateService) UpdateLanguage(ctx context.Context, companyID, candidateID string, lang domain.CandidateLanguage) error {
	const op = "UpdateLanguage"
	_, err := s.candidateRepo.UpdateLanguage(ctx, lang.ID, &lang)
	return err
}

func (s *CandidateService) DeleteLanguage(ctx context.Context, companyID, candidateID, langID string) error {
	const op = "DeleteLanguage"
	return s.candidateRepo.DeleteLanguage(ctx, langID)
}

func (s *CandidateService) ListLanguages(ctx context.Context, companyID, candidateID string) ([]domain.CandidateLanguage, error) {
	const op = "ListLanguages"
	return s.candidateRepo.ListLanguages(ctx, candidateID)
}

func (s *CandidateService) AddDocument(ctx context.Context, companyID, candidateID string, doc domain.CandidateDocument) (*domain.CandidateDocument, error) {
	const op = "AddDocument"
	doc.ID = uuid.New().String()
	doc.CandidateID = candidateID
	doc.CompanyID = companyID
	doc.CreatedAt = time.Now()
	result, err := s.candidateRepo.AddDocument(ctx, &doc)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *CandidateService) ListDocuments(ctx context.Context, companyID, candidateID string) ([]domain.CandidateDocument, error) {
	const op = "ListDocuments"
	return s.candidateRepo.ListDocuments(ctx, candidateID)
}

func (s *CandidateService) SearchBySkills(ctx context.Context, companyID string, skills []string) ([]domain.Candidate, error) {
	const op = "SearchBySkills"
	return s.candidateRepo.SearchBySkills(ctx, companyID, skills)
}
