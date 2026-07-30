package engine

import (
	"context"
	"time"

	"github.com/rrhhumand/api/internal/recruitment/domain"
	"github.com/rrhhumand/api/internal/recruitment/repository"
)

type OfferContent struct {
	Subject         string `json:"subject"`
	Greeting        string `json:"greeting"`
	Body            string `json:"body"`
	SalarySection   string `json:"salary_section"`
	BenefitsSection string `json:"benefits_section"`
	Conditions      string `json:"conditions"`
	DeadlineText    string `json:"deadline_text"`
	Closing         string `json:"closing"`
}

type OfferEngine struct {
	offerRepo       *repository.OfferRepo
	applicationRepo *repository.ApplicationRepo
	candidateRepo   *repository.CandidateRepo
}

func NewOfferEngine(offerRepo *repository.OfferRepo, applicationRepo *repository.ApplicationRepo, candidateRepo *repository.CandidateRepo) *OfferEngine {
	return &OfferEngine{
		offerRepo:       offerRepo,
		applicationRepo: applicationRepo,
		candidateRepo:   candidateRepo,
	}
}

func (e *OfferEngine) GenerateOfferContent(ctx context.Context, companyID, applicationID string, templateVars map[string]string) (*OfferContent, error) {
	app, err := e.applicationRepo.GetByID(ctx, companyID, applicationID)
	if err != nil {
		return nil, err
	}

	candidate, err := e.candidateRepo.GetByID(ctx, app.CompanyID, app.CandidateID)
	if err != nil {
		return nil, err
	}

	content := &OfferContent{
		Subject:         templateVars["subject"],
		Greeting:        "Estimado/a " + candidate.FirstName + " " + candidate.LastName + ",",
		Body:            templateVars["body"],
		SalarySection:   templateVars["salary"],
		BenefitsSection: templateVars["benefits"],
		Conditions:      templateVars["conditions"],
		DeadlineText:    "Esta oferta tiene una validez de " + templateVars["deadline_days"] + " días hábiles.",
		Closing:         "Quedamos atentos a su respuesta.",
	}

	return content, nil
}

func (e *OfferEngine) CalculateDeadline(days int) time.Time {
	return time.Now().AddDate(0, 0, days)
}

func (e *OfferEngine) CheckExpiredOffers() error {
	return nil
}

func (e *OfferEngine) ValidateOffer(offer domain.Offer) error {
	if offer.PositionTitle == "" {
		return domain.ErrInvalidInput
	}
	if offer.ApplicationID == "" {
		return domain.ErrInvalidInput
	}
	return nil
}
