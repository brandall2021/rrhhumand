package application

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/recruitment/domain"
	"github.com/rrhhumand/api/internal/recruitment/repository"
)

type EmailService struct {
	templateRepo *repository.EmailRepo
	emailLogRepo *repository.EmailRepo
}

func NewEmailService(templateRepo *repository.EmailRepo, emailLogRepo *repository.EmailRepo) *EmailService {
	return &EmailService{
		templateRepo: templateRepo,
		emailLogRepo: emailLogRepo,
	}
}

func (s *EmailService) CreateTemplate(ctx context.Context, companyID string, tpl *domain.EmailTemplate) (*domain.EmailTemplate, error) {
	const op = "CreateEmailTemplate"
	tpl.ID = uuid.New().String()
	tpl.CompanyID = companyID
	now := time.Now()
	tpl.CreatedAt = now
	tpl.UpdatedAt = now
	if !tpl.Active {
		tpl.Active = true
	}
	result, err := s.templateRepo.CreateTemplate(ctx, companyID, tpl)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *EmailService) GetTemplate(ctx context.Context, companyID, id string) (*domain.EmailTemplate, error) {
	const op = "GetEmailTemplate"
	tpl, err := s.templateRepo.GetTemplate(ctx, companyID, id)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return tpl, nil
}

func (s *EmailService) GetByCode(ctx context.Context, companyID, code string) (*domain.EmailTemplate, error) {
	const op = "GetEmailTemplateByCode"
	tpl, err := s.templateRepo.GetTemplateByCode(ctx, companyID, code)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return tpl, nil
}

func (s *EmailService) ListTemplates(ctx context.Context, companyID string) ([]domain.EmailTemplate, error) {
	const op = "ListEmailTemplates"
	return s.templateRepo.ListTemplates(ctx, companyID, "")
}

func (s *EmailService) UpdateTemplate(ctx context.Context, companyID, id string, tpl *domain.EmailTemplate) (*domain.EmailTemplate, error) {
	const op = "UpdateEmailTemplate"
	result, err := s.templateRepo.UpdateTemplate(ctx, companyID, id, tpl)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *EmailService) DeleteTemplate(ctx context.Context, companyID, id string) error {
	const op = "DeleteEmailTemplate"
	return s.templateRepo.DeleteTemplate(ctx, companyID, id)
}

func (s *EmailService) ListEmails(ctx context.Context, companyID, applicationID string) ([]domain.EmailLog, error) {
	const op = "ListEmails"
	result, err := s.emailLogRepo.ListEmails(ctx, companyID, applicationID)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *EmailService) SendEmail(ctx context.Context, templateCode, recipient string, vars map[string]string) error {
	const op = "SendEmail"
	tpl, err := s.templateRepo.GetTemplateByCode(ctx, "", templateCode)
	if err != nil {
		return svcErr(op, err)
	}
	subject := tpl.Subject
	body := tpl.BodyHTML
	for k, v := range vars {
		placeholder := "{{" + k + "}}"
		subject = strings.ReplaceAll(subject, placeholder, v)
		body = strings.ReplaceAll(body, placeholder, v)
	}
	log := &domain.EmailLog{
		ID:             uuid.New().String(),
		TemplateID:     &tpl.ID,
		RecipientEmail: recipient,
		Subject:        subject,
		Body:           &body,
		Status:         "SENT",
		SentAt:         time.Now(),
	}
	_, err = s.emailLogRepo.LogEmail(ctx, log)
	return err
}
