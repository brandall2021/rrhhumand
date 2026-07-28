package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/rrhhumand/api/internal/payroll/features/domain"
	"github.com/rrhhumand/api/internal/payroll/features/repository"
)

type ArcaService struct {
	arcaRepo *repository.ArcaRepo
}

func NewArcaService(arcaRepo *repository.ArcaRepo) *ArcaService {
	return &ArcaService{arcaRepo: arcaRepo}
}

func arcaSvcErr(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("arca_svc.%s: %w", op, err)
}

func (s *ArcaService) CreateMapping(ctx context.Context, companyID uuid.UUID, userID uuid.UUID, m *domain.ArcaConceptMapping) (*domain.ArcaConceptMapping, error) {
	m.ID = uuid.New()
	m.CompanyID = companyID
	m.CreatedBy = userID
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	if err := s.arcaRepo.CreateMapping(ctx, m); err != nil {
		return nil, arcaSvcErr("CreateMapping", err)
	}
	return m, nil
}

func (s *ArcaService) GetMapping(ctx context.Context, companyID, id uuid.UUID) (*domain.ArcaConceptMapping, error) {
	return s.arcaRepo.GetMapping(ctx, companyID, id)
}

func (s *ArcaService) ListMappings(ctx context.Context, companyID uuid.UUID) ([]domain.ArcaConceptMapping, error) {
	return s.arcaRepo.ListMappings(ctx, companyID)
}

func (s *ArcaService) UpdateMapping(ctx context.Context, companyID uuid.UUID, m *domain.ArcaConceptMapping) (*domain.ArcaConceptMapping, error) {
	m.CompanyID = companyID
	m.UpdatedAt = time.Now()
	if err := s.arcaRepo.UpdateMapping(ctx, m); err != nil {
		return nil, arcaSvcErr("UpdateMapping", err)
	}
	return m, nil
}

func (s *ArcaService) DeleteMapping(ctx context.Context, companyID, id uuid.UUID) error {
	return s.arcaRepo.DeleteMapping(ctx, companyID, id)
}

func (s *ArcaService) GetActiveMappings(ctx context.Context, companyID, conceptID uuid.UUID, date time.Time) ([]domain.ArcaConceptMapping, error) {
	return s.arcaRepo.GetActiveMappingsForConcept(ctx, companyID, conceptID, date)
}

func (s *ArcaService) CreateExport(ctx context.Context, companyID uuid.UUID, userID uuid.UUID, e *domain.ArcaExport) (*domain.ArcaExport, error) {
	e.ID = uuid.New()
	e.CompanyID = companyID
	e.GeneratedBy = userID
	e.GeneratedAt = time.Now()
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()
	if err := s.arcaRepo.CreateExport(ctx, e); err != nil {
		return nil, arcaSvcErr("CreateExport", err)
	}
	return e, nil
}

func (s *ArcaService) GetExport(ctx context.Context, companyID, id uuid.UUID) (*domain.ArcaExport, error) {
	return s.arcaRepo.GetExport(ctx, companyID, id)
}

func (s *ArcaService) ListExports(ctx context.Context, companyID uuid.UUID, runID *uuid.UUID, limit, offset int) ([]domain.ArcaExport, error) {
	return s.arcaRepo.ListExports(ctx, companyID, runID, limit, offset)
}

func (s *ArcaService) GenerateExport(ctx context.Context, companyID uuid.UUID, runID uuid.UUID, exportType string, userID uuid.UUID) (*domain.ArcaExport, error) {
	e := &domain.ArcaExport{
		ID:            uuid.New(),
		CompanyID:     companyID,
		RunID:         runID,
		ExportType:    exportType,
		FileName:      fmt.Sprintf("ARCA_%s_%s.txt", exportType, time.Now().Format("20060102_150405")),
		Status:        "GENERATED",
		EmployeeCount: 0,
		TotalAmount:   decimal.Zero,
		GeneratedBy:   userID,
		GeneratedAt:   time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := s.arcaRepo.CreateExport(ctx, e); err != nil {
		return nil, arcaSvcErr("GenerateExport", err)
	}
	return e, nil
}

func (s *ArcaService) ValidateExport(ctx context.Context, exportID uuid.UUID) error {
	fields := map[string]any{
		"status": "VALIDATED",
	}
	return s.arcaRepo.UpdateExportStatus(ctx, exportID, "VALIDATED", fields)
}

func (s *ArcaService) GenerateSICOSSFile(ctx context.Context, companyID uuid.UUID, runID uuid.UUID) (*domain.ArcaExport, error) {
	e := &domain.ArcaExport{
		ID:            uuid.New(),
		CompanyID:     companyID,
		RunID:         runID,
		ExportType:    "SICOSS",
		FileName:      fmt.Sprintf("SICOSS_%s.txt", time.Now().Format("20060102_150405")),
		Status:        "GENERATED",
		EmployeeCount: 0,
		TotalAmount:   decimal.Zero,
		GeneratedBy:   uuid.Nil,
		GeneratedAt:   time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := s.arcaRepo.CreateExport(ctx, e); err != nil {
		return nil, arcaSvcErr("GenerateSICOSSFile", err)
	}
	return e, nil
}

func (s *ArcaService) GenerateSIAPFile(ctx context.Context, companyID uuid.UUID, runID uuid.UUID) (*domain.ArcaExport, error) {
	e := &domain.ArcaExport{
		ID:            uuid.New(),
		CompanyID:     companyID,
		RunID:         runID,
		ExportType:    "SIAP",
		FileName:      fmt.Sprintf("SIAP_%s.txt", time.Now().Format("20060102_150405")),
		Status:        "GENERATED",
		EmployeeCount: 0,
		TotalAmount:   decimal.Zero,
		GeneratedBy:   uuid.Nil,
		GeneratedAt:   time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := s.arcaRepo.CreateExport(ctx, e); err != nil {
		return nil, arcaSvcErr("GenerateSIAPFile", err)
	}
	return e, nil
}
