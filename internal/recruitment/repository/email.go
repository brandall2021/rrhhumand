package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/recruitment/domain"
)

type EmailRepo struct {
	pool *pgxpool.Pool
}

func NewEmailRepo(pool *pgxpool.Pool) *EmailRepo {
	return &EmailRepo{pool: pool}
}

func (r *EmailRepo) CreateTemplate(ctx context.Context, companyID string, req *domain.EmailTemplate) (*domain.EmailTemplate, error) {
	t := &domain.EmailTemplate{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO email_templates (company_id, name, code, subject, body_html, body_text, variables, category, active)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		 RETURNING id, company_id, name, code, subject, body_html, body_text, variables, category, active, created_at, updated_at`,
		companyID, req.Name, req.Code, req.Subject, req.BodyHTML, req.BodyText, req.Variables, req.Category, req.Active,
	).Scan(&t.ID, &t.CompanyID, &t.Name, &t.Code, &t.Subject, &t.BodyHTML, &t.BodyText, &t.Variables, &t.Category, &t.Active, &t.CreatedAt, &t.UpdatedAt)
	return t, err
}

func (r *EmailRepo) GetTemplate(ctx context.Context, companyID, id string) (*domain.EmailTemplate, error) {
	t := &domain.EmailTemplate{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, name, code, subject, body_html, body_text, variables, category, active, created_at, updated_at
		 FROM email_templates WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&t.ID, &t.CompanyID, &t.Name, &t.Code, &t.Subject, &t.BodyHTML, &t.BodyText, &t.Variables, &t.Category, &t.Active, &t.CreatedAt, &t.UpdatedAt)
	return t, err
}

func (r *EmailRepo) GetTemplateByCode(ctx context.Context, companyID, code string) (*domain.EmailTemplate, error) {
	t := &domain.EmailTemplate{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, name, code, subject, body_html, body_text, variables, category, active, created_at, updated_at
		 FROM email_templates WHERE company_id=$1 AND code=$2`, companyID, code,
	).Scan(&t.ID, &t.CompanyID, &t.Name, &t.Code, &t.Subject, &t.BodyHTML, &t.BodyText, &t.Variables, &t.Category, &t.Active, &t.CreatedAt, &t.UpdatedAt)
	return t, err
}

func (r *EmailRepo) ListTemplates(ctx context.Context, companyID string, category string) ([]domain.EmailTemplate, error) {
	query := `SELECT id, company_id, name, code, subject, body_html, body_text, variables, category, active, created_at, updated_at
		 FROM email_templates WHERE company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if category != "" {
		query += fmt.Sprintf(" AND category=$%d", argIdx)
		args = append(args, category)
		argIdx++
	}
	query += " ORDER BY name"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []domain.EmailTemplate
	for rows.Next() {
		var t domain.EmailTemplate
		rows.Scan(&t.ID, &t.CompanyID, &t.Name, &t.Code, &t.Subject, &t.BodyHTML, &t.BodyText, &t.Variables, &t.Category, &t.Active, &t.CreatedAt, &t.UpdatedAt)
		templates = append(templates, t)
	}
	return templates, nil
}

func (r *EmailRepo) UpdateTemplate(ctx context.Context, companyID, id string, req *domain.EmailTemplate) (*domain.EmailTemplate, error) {
	t := &domain.EmailTemplate{}
	err := r.pool.QueryRow(ctx,
		`UPDATE email_templates SET
		 name=COALESCE($3,name), subject=COALESCE($4,subject), body_html=COALESCE($5,body_html),
		 body_text=COALESCE($6,body_text), variables=COALESCE($7,variables),
		 category=COALESCE($8,category), active=COALESCE($9,active), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, name, code, subject, body_html, body_text, variables, category, active, created_at, updated_at`,
		companyID, id, req.Name, req.Subject, req.BodyHTML, req.BodyText, req.Variables, req.Category, req.Active,
	).Scan(&t.ID, &t.CompanyID, &t.Name, &t.Code, &t.Subject, &t.BodyHTML, &t.BodyText, &t.Variables, &t.Category, &t.Active, &t.CreatedAt, &t.UpdatedAt)
	return t, err
}

func (r *EmailRepo) DeleteTemplate(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM email_templates WHERE company_id=$1 AND id=$2`, companyID, id)
	return err
}

func (r *EmailRepo) LogEmail(ctx context.Context, req *domain.EmailLog) (*domain.EmailLog, error) {
	e := &domain.EmailLog{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO email_log (company_id, template_id, application_id, candidate_id, recipient_email, subject, body, status, error_message)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		 RETURNING id, company_id, template_id, application_id, candidate_id, recipient_email, subject, body, status, error_message, sent_at`,
		req.CompanyID, req.TemplateID, req.ApplicationID, req.CandidateID,
		req.RecipientEmail, req.Subject, req.Body, req.Status, req.ErrorMessage,
	).Scan(&e.ID, &e.CompanyID, &e.TemplateID, &e.ApplicationID, &e.CandidateID,
		&e.RecipientEmail, &e.Subject, &e.Body, &e.Status, &e.ErrorMessage, &e.SentAt)
	return e, err
}

func (r *EmailRepo) ListEmails(ctx context.Context, companyID string, applicationID string) ([]domain.EmailLog, error) {
	query := `SELECT id, company_id, template_id, application_id, candidate_id, recipient_email, subject, body, status, error_message, sent_at
		 FROM email_log WHERE company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if applicationID != "" {
		query += fmt.Sprintf(" AND application_id=$%d", argIdx)
		args = append(args, applicationID)
		argIdx++
	}
	query += " ORDER BY sent_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []domain.EmailLog
	for rows.Next() {
		var e domain.EmailLog
		rows.Scan(&e.ID, &e.CompanyID, &e.TemplateID, &e.ApplicationID, &e.CandidateID,
			&e.RecipientEmail, &e.Subject, &e.Body, &e.Status, &e.ErrorMessage, &e.SentAt)
		emails = append(emails, e)
	}
	return emails, nil
}
