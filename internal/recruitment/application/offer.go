package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/recruitment/domain"
	"github.com/rrhhumand/api/internal/recruitment/repository"
)

type CreateOfferReq struct {
	ApplicationID    string     `json:"application_id"`
	PositionTitle    string     `json:"position_title"`
	DepartmentID     *string    `json:"department_id"`
	OfferType        string     `json:"offer_type"`
	StartDate        *time.Time `json:"start_date"`
	EmploymentType   *string    `json:"employment_type"`
	WorkMode         *string    `json:"work_mode"`
	SalaryAmount     *float64   `json:"salary_amount"`
	SalaryCurrency   *string    `json:"salary_currency"`
	SalaryPeriod     *string    `json:"salary_period"`
	VariableComp     *string    `json:"variable_compensation"`
	BenefitsSummary  *string    `json:"benefits_summary"`
	EquityTerms      *string    `json:"equity_terms"`
	Conditions       *string    `json:"conditions"`
	Notes            *string    `json:"notes"`
	ResponseDeadline *time.Time `json:"response_deadline"`
}

type OfferService struct {
	offerRepo       *repository.OfferRepo
	applicationRepo *repository.ApplicationRepo
}

func NewOfferService(offerRepo *repository.OfferRepo, applicationRepo *repository.ApplicationRepo) *OfferService {
	return &OfferService{
		offerRepo:       offerRepo,
		applicationRepo: applicationRepo,
	}
}

func (s *OfferService) Create(ctx context.Context, companyID, createdBy string, req *CreateOfferReq) (*domain.Offer, error) {
	const op = "CreateOffer"
	now := time.Now()
	o := &domain.Offer{
		ID:               uuid.New().String(),
		CompanyID:        companyID,
		ApplicationID:    req.ApplicationID,
		PositionTitle:    req.PositionTitle,
		DepartmentID:     req.DepartmentID,
		OfferType:        req.OfferType,
		StartDate:        req.StartDate,
		EmploymentType:   req.EmploymentType,
		WorkMode:         req.WorkMode,
		SalaryAmount:     req.SalaryAmount,
		SalaryCurrency:   req.SalaryCurrency,
		SalaryPeriod:     req.SalaryPeriod,
		VariableComp:     req.VariableComp,
		BenefitsSummary:  req.BenefitsSummary,
		EquityTerms:      req.EquityTerms,
		Conditions:       req.Conditions,
		Notes:            req.Notes,
		ResponseDeadline: req.ResponseDeadline,
		Status:           domain.OfferStatusDraft,
		CreatedBy:        &createdBy,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	result, err := s.offerRepo.Create(ctx, companyID, o)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *OfferService) GetByID(ctx context.Context, companyID, id string) (*domain.Offer, error) {
	const op = "GetOffer"
	return s.offerRepo.GetByID(ctx, companyID, id)
}

func (s *OfferService) List(ctx context.Context, companyID, applicationID, status string) ([]domain.Offer, error) {
	const op = "ListOffers"
	return s.offerRepo.List(ctx, companyID, applicationID, status)
}

func (s *OfferService) Update(ctx context.Context, companyID, id string, req *domain.Offer) (*domain.Offer, error) {
	const op = "UpdateOffer"
	req.UpdatedAt = time.Now()
	result, err := s.offerRepo.Update(ctx, companyID, id, req)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *OfferService) SubmitForApproval(ctx context.Context, companyID, id string) error {
	const op = "SubmitOfferForApproval"
	return s.offerRepo.UpdateStatus(ctx, companyID, id, string(domain.OfferStatusPendingApproval))
}

func (s *OfferService) Approve(ctx context.Context, companyID, id string) error {
	const op = "ApproveOffer"
	return s.offerRepo.UpdateStatus(ctx, companyID, id, string(domain.OfferStatusApproved))
}

func (s *OfferService) Send(ctx context.Context, companyID, id string) error {
	const op = "SendOffer"
	return s.offerRepo.UpdateStatus(ctx, companyID, id, string(domain.OfferStatusSent))
}

func (s *OfferService) Accept(ctx context.Context, companyID, id string) error {
	const op = "AcceptOffer"
	return s.offerRepo.UpdateStatus(ctx, companyID, id, string(domain.OfferStatusAccepted))
}

func (s *OfferService) Reject(ctx context.Context, companyID, id string) error {
	const op = "RejectOffer"
	return s.offerRepo.UpdateStatus(ctx, companyID, id, string(domain.OfferStatusRejected))
}

func (s *OfferService) Withdraw(ctx context.Context, companyID, id string) error {
	const op = "WithdrawOffer"
	return s.offerRepo.UpdateStatus(ctx, companyID, id, string(domain.OfferStatusWithdrawn))
}

func (s *OfferService) AddNegotiation(ctx context.Context, companyID, offerID string, neg domain.OfferNegotiation) (*domain.OfferNegotiation, error) {
	const op = "AddNegotiation"
	neg.ID = uuid.New().String()
	neg.OfferID = offerID
	neg.CreatedAt = time.Now()
	result, err := s.offerRepo.AddNegotiation(ctx, &neg)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *OfferService) UpdateNegotiation(ctx context.Context, companyID, offerID string, neg domain.OfferNegotiation) error {
	const op = "UpdateNegotiation"
	_, err := s.offerRepo.UpdateNegotiation(ctx, neg.ID, &neg)
	return err
}

func (s *OfferService) ListNegotiations(ctx context.Context, companyID, offerID string) ([]domain.OfferNegotiation, error) {
	const op = "ListNegotiations"
	return s.offerRepo.ListNegotiations(ctx, offerID)
}

func (s *OfferService) AddDocument(ctx context.Context, companyID, offerID string, doc domain.OfferDocument) (*domain.OfferDocument, error) {
	const op = "AddOfferDocument"
	doc.ID = uuid.New().String()
	doc.OfferID = offerID
	doc.CreatedAt = time.Now()
	result, err := s.offerRepo.AddDocument(ctx, &doc)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *OfferService) ListDocuments(ctx context.Context, companyID, offerID string) ([]domain.OfferDocument, error) {
	const op = "ListOfferDocuments"
	return s.offerRepo.ListDocuments(ctx, offerID)
}
