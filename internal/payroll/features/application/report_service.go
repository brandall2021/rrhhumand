package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/payroll/features/domain"
	"github.com/rrhhumand/api/internal/payroll/features/repository"
)

type ReportService struct {
	reportRepo *repository.ReportRepo
}

func NewReportService(reportRepo *repository.ReportRepo) *ReportService {
	return &ReportService{reportRepo: reportRepo}
}

func reportSvcErr(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("report_svc.%s: %w", op, err)
}

func (s *ReportService) CreateTemplate(ctx context.Context, companyID uuid.UUID, userID uuid.UUID, t *domain.ReportTemplate) (*domain.ReportTemplate, error) {
	t.ID = uuid.New()
	t.CompanyID = companyID
	t.CreatedBy = userID
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	if err := s.reportRepo.CreateTemplate(ctx, t); err != nil {
		return nil, reportSvcErr("CreateTemplate", err)
	}
	return t, nil
}

func (s *ReportService) GetTemplate(ctx context.Context, companyID, id uuid.UUID) (*domain.ReportTemplate, error) {
	return s.reportRepo.GetTemplate(ctx, companyID, id)
}

func (s *ReportService) ListTemplates(ctx context.Context, companyID uuid.UUID) ([]domain.ReportTemplate, error) {
	return s.reportRepo.ListTemplates(ctx, companyID)
}

func (s *ReportService) UpdateTemplate(ctx context.Context, companyID uuid.UUID, t *domain.ReportTemplate) (*domain.ReportTemplate, error) {
	t.CompanyID = companyID
	t.UpdatedAt = time.Now()
	if err := s.reportRepo.UpdateTemplate(ctx, t); err != nil {
		return nil, reportSvcErr("UpdateTemplate", err)
	}
	return t, nil
}

func (s *ReportService) DeleteTemplate(ctx context.Context, companyID, id uuid.UUID) error {
	return s.reportRepo.DeleteTemplate(ctx, companyID, id)
}

func (s *ReportService) GenerateReport(ctx context.Context, companyID uuid.UUID, runID uuid.UUID, templateID uuid.UUID, format string, userID uuid.UUID) (*domain.ReportExport, error) {
	e := &domain.ReportExport{
		ID:          uuid.New(),
		CompanyID:   companyID,
		RunID:       &runID,
		TemplateID:  &templateID,
		ReportType:  "CUSTOM",
		FileFormat:  format,
		FileName:    fmt.Sprintf("REPORT_%s.%s", time.Now().Format("20060102_150405"), format),
		Status:      "GENERATED",
		GeneratedBy: userID,
		GeneratedAt: time.Now(),
		CreatedAt:   time.Now(),
	}
	if err := s.reportRepo.CreateExport(ctx, e); err != nil {
		return nil, reportSvcErr("GenerateReport", err)
	}
	return e, nil
}

func (s *ReportService) GetReportExport(ctx context.Context, companyID, exportID uuid.UUID) (*domain.ReportExport, error) {
	return s.reportRepo.GetExport(ctx, companyID, exportID)
}

func (s *ReportService) ListReportExports(ctx context.Context, companyID uuid.UUID) ([]domain.ReportExport, error) {
	return s.reportRepo.ListExports(ctx, companyID, 0, 0)
}
