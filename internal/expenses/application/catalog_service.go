package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/expenses/domain"
)

type CatalogRepository interface {
	CreateCategory(ctx context.Context, category *domain.ExpenseCategory) error
	GetCategory(ctx context.Context, companyID, id uuid.UUID) (*domain.ExpenseCategory, error)
	ListCategories(ctx context.Context, companyID uuid.UUID) ([]domain.ExpenseCategory, error)
	UpdateCategory(ctx context.Context, category *domain.ExpenseCategory) error
	DeleteCategory(ctx context.Context, companyID, id uuid.UUID) error

	CreatePaymentMethod(ctx context.Context, pm *domain.ExpensePaymentMethod) error
	ListPaymentMethods(ctx context.Context, companyID uuid.UUID) ([]domain.ExpensePaymentMethod, error)
	UpdatePaymentMethod(ctx context.Context, pm *domain.ExpensePaymentMethod) error
}

type CatalogService struct {
	catalogRepo CatalogRepository
	auditRepo   AuditRepository
}

func NewCatalogService(catalogRepo CatalogRepository, auditRepo AuditRepository) *CatalogService {
	return &CatalogService{catalogRepo: catalogRepo, auditRepo: auditRepo}
}

func (s *CatalogService) CreateCategory(ctx context.Context, companyID uuid.UUID, userID uuid.UUID, category *domain.ExpenseCategory) (*domain.ExpenseCategory, error) {
	const op = "CreateCategory"
	now := time.Now()
	category.ID = uuid.New()
	category.CompanyID = companyID
	category.CreatedAt = now
	category.UpdatedAt = now
	if err := s.catalogRepo.CreateCategory(ctx, category); err != nil {
		return nil, svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: companyID, UserID: userID,
		Action: "category.created", EntityType: "expense_category", EntityID: category.ID, CreatedAt: now,
	})
	return category, nil
}

func (s *CatalogService) GetCategory(ctx context.Context, companyID, id uuid.UUID) (*domain.ExpenseCategory, error) {
	const op = "GetCategory"
	cat, err := s.catalogRepo.GetCategory(ctx, companyID, id)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return cat, nil
}

func (s *CatalogService) ListCategories(ctx context.Context, companyID uuid.UUID) ([]domain.ExpenseCategory, error) {
	const op = "ListCategories"
	cats, err := s.catalogRepo.ListCategories(ctx, companyID)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return cats, nil
}

func (s *CatalogService) UpdateCategory(ctx context.Context, companyID uuid.UUID, userID uuid.UUID, category *domain.ExpenseCategory) (*domain.ExpenseCategory, error) {
	const op = "UpdateCategory"
	existing, err := s.catalogRepo.GetCategory(ctx, companyID, category.ID)
	if err != nil {
		return nil, svcErr(op, err)
	}
	category.CompanyID = companyID
	category.CreatedAt = existing.CreatedAt
	category.UpdatedAt = time.Now()
	if err := s.catalogRepo.UpdateCategory(ctx, category); err != nil {
		return nil, svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: companyID, UserID: userID,
		Action: "category.updated", EntityType: "expense_category", EntityID: category.ID, CreatedAt: time.Now(),
	})
	return category, nil
}

func (s *CatalogService) DeleteCategory(ctx context.Context, companyID, id uuid.UUID) error {
	const op = "DeleteCategory"
	if err := s.catalogRepo.DeleteCategory(ctx, companyID, id); err != nil {
		return svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: companyID,
		Action: "category.deleted", EntityType: "expense_category", EntityID: id, CreatedAt: time.Now(),
	})
	return nil
}

func (s *CatalogService) CreatePaymentMethod(ctx context.Context, companyID uuid.UUID, userID uuid.UUID, pm *domain.ExpensePaymentMethod) (*domain.ExpensePaymentMethod, error) {
	const op = "CreatePaymentMethod"
	pm.ID = uuid.New()
	pm.CompanyID = companyID
	pm.CreatedAt = time.Now()
	if err := s.catalogRepo.CreatePaymentMethod(ctx, pm); err != nil {
		return nil, svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: companyID, UserID: userID,
		Action: "payment_method.created", EntityType: "expense_payment_method", EntityID: pm.ID, CreatedAt: time.Now(),
	})
	return pm, nil
}

func (s *CatalogService) ListPaymentMethods(ctx context.Context, companyID uuid.UUID) ([]domain.ExpensePaymentMethod, error) {
	const op = "ListPaymentMethods"
	pms, err := s.catalogRepo.ListPaymentMethods(ctx, companyID)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return pms, nil
}

func (s *CatalogService) UpdatePaymentMethod(ctx context.Context, companyID uuid.UUID, userID uuid.UUID, pm *domain.ExpensePaymentMethod) (*domain.ExpensePaymentMethod, error) {
	const op = "UpdatePaymentMethod"
	pm.CompanyID = companyID
	if err := s.catalogRepo.UpdatePaymentMethod(ctx, pm); err != nil {
		return nil, svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: companyID, UserID: userID,
		Action: "payment_method.updated", EntityType: "expense_payment_method", EntityID: pm.ID, CreatedAt: time.Now(),
	})
	return pm, nil
}
