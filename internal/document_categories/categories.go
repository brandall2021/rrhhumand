package document_categories

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/models"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

type CategoryRepository struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{pool: pool}
}

func (r *CategoryRepository) Create(ctx context.Context, cat *models.DocumentCategory) error {
	query := `
		INSERT INTO document_categories (id, company_id, name, description, parent_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at`
	return r.pool.QueryRow(ctx, query,
		cat.ID, cat.CompanyID, cat.Name, cat.Description, cat.ParentID,
	).Scan(&cat.CreatedAt, &cat.UpdatedAt)
}

func (r *CategoryRepository) GetByID(ctx context.Context, id, companyID string) (*models.DocumentCategory, error) {
	query := `
		SELECT id, company_id, name, description, parent_id, is_active, created_at, updated_at
		FROM document_categories WHERE id=$1 AND company_id=$2`
	cat := &models.DocumentCategory{}
	err := r.pool.QueryRow(ctx, query, id, companyID).Scan(
		&cat.ID, &cat.CompanyID, &cat.Name, &cat.Description, &cat.ParentID,
		&cat.IsActive, &cat.CreatedAt, &cat.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return cat, nil
}

func (r *CategoryRepository) List(ctx context.Context, companyID string) ([]models.DocumentCategory, error) {
	query := `
		SELECT id, company_id, name, description, parent_id, is_active, created_at, updated_at
		FROM document_categories WHERE company_id=$1 AND is_active=true
		ORDER BY name`
	rows, err := r.pool.Query(ctx, query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []models.DocumentCategory
	for rows.Next() {
		var c models.DocumentCategory
		if err := rows.Scan(&c.ID, &c.CompanyID, &c.Name, &c.Description, &c.ParentID,
			&c.IsActive, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, nil
}

func (r *CategoryRepository) Update(ctx context.Context, cat *models.DocumentCategory) error {
	query := `
		UPDATE document_categories SET name=$1, description=$2, parent_id=$3, is_active=$4, updated_at=NOW()
		WHERE id=$5 AND company_id=$6`
	_, err := r.pool.Exec(ctx, query,
		cat.Name, cat.Description, cat.ParentID, cat.IsActive, cat.ID, cat.CompanyID,
	)
	return err
}

func (r *CategoryRepository) Delete(ctx context.Context, id, companyID string) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM document_categories WHERE id=$1 AND company_id=$2`, id, companyID,
	)
	return err
}

type CategoryService struct {
	repo *CategoryRepository
}

func NewCategoryService(repo *CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(ctx context.Context, companyID string, req *CreateCategoryRequest) (*models.DocumentCategory, error) {
	cat := &models.DocumentCategory{
		ID:          uuid.New().String(),
		CompanyID:   companyID,
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
		IsActive:    true,
	}

	if err := s.repo.Create(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *CategoryService) GetByID(ctx context.Context, id, companyID string) (*models.DocumentCategory, error) {
	cat, err := s.repo.GetByID(ctx, id, companyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return cat, nil
}

func (s *CategoryService) List(ctx context.Context, companyID string) ([]models.DocumentCategory, error) {
	return s.repo.List(ctx, companyID)
}

func (s *CategoryService) Update(ctx context.Context, id, companyID string, req *UpdateCategoryRequest) (*models.DocumentCategory, error) {
	cat, err := s.repo.GetByID(ctx, id, companyID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	if req.Name != nil {
		cat.Name = *req.Name
	}
	if req.Description != nil {
		cat.Description = req.Description
	}
	if req.ParentID != nil {
		cat.ParentID = req.ParentID
	}
	if req.IsActive != nil {
		cat.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *CategoryService) Delete(ctx context.Context, id, companyID string) error {
	_, err := s.repo.GetByID(ctx, id, companyID)
	if err != nil {
		return errors.New("category not found")
	}
	return s.repo.Delete(ctx, id, companyID)
}

type CreateCategoryRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description,omitempty"`
	ParentID    *string `json:"parent_id,omitempty"`
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	ParentID    *string `json:"parent_id,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type CategoryHandler struct {
	service *CategoryService
}

func NewCategoryHandler(service *CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) Create(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)

	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	cat, err := h.service.Create(c.Request.Context(), companyID, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, cat)
}

func (h *CategoryHandler) List(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)

	cats, err := h.service.List(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, cats)
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	cat, err := h.service.GetByID(c.Request.Context(), id, companyID)
	if err != nil {
		if err.Error() == "category not found" {
			response.NotFound(c, "Category not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, cat)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	cat, err := h.service.Update(c.Request.Context(), id, companyID, &req)
	if err != nil {
		if err.Error() == "category not found" {
			response.NotFound(c, "Category not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, cat)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id, companyID); err != nil {
		if err.Error() == "category not found" {
			response.NotFound(c, "Category not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.NoContent(c)
}
