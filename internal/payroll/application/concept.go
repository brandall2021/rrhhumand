package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/rrhhumand/api/internal/payroll/domain"
	"github.com/rrhhumand/api/internal/payroll/repository"
)

type CreateConceptInput struct {
	Code            string
	Name            string
	Description     *string
	ConceptType     string
	Taxability      string
	CalculationType string
	BaseConceptID   *string
	SortOrder       int
}

type UpdateConceptInput struct {
	Name            *string
	Description     *string
	ConceptType     *string
	Taxability      *string
	CalculationType *string
	BaseConceptID   *string
	Active          *bool
	SortOrder       *int
}

type CreateRuleInput struct {
	ConceptID     string
	RuleType      string
	Formula       *string
	Parameters    map[string]any
	Priority      int
	EffectiveFrom string
	EffectiveTo   *string
}

type UpdateRuleInput struct {
	RuleType    *string
	Formula     *string
	Parameters  map[string]any
	Priority    *int
	Active      *bool
	EffectiveTo *string
}

func (s *PayrollService) CreateConcept(ctx context.Context, companyID, userID string, req CreateConceptInput) (*domain.PayrollConcept, error) {
	c := &domain.PayrollConcept{
		ID:              uuid.NewString(),
		CompanyID:       companyID,
		Code:            req.Code,
		Name:            req.Name,
		Description:     req.Description,
		ConceptType:     req.ConceptType,
		Taxability:      req.Taxability,
		CalculationType: req.CalculationType,
		BaseConceptID:   req.BaseConceptID,
		Active:          true,
		EffectiveFrom:   time.Now().AddDate(0, -1, 0),
		SortOrder:       req.SortOrder,
		CreatedBy:       userID,
	}
	if err := s.repo.CreateConcept(ctx, c); err != nil {
		return nil, fmt.Errorf("create concept: %w", err)
	}
	return c, nil
}

func (s *PayrollService) UpdateConcept(ctx context.Context, companyID, id string, req UpdateConceptInput) (*domain.PayrollConcept, error) {
	c, err := s.repo.GetConcept(ctx, companyID, id)
	if err != nil {
		return nil, fmt.Errorf("update concept: get: %w", err)
	}
	if req.Name != nil {
		c.Name = *req.Name
	}
	if req.Description != nil {
		c.Description = req.Description
	}
	if req.ConceptType != nil {
		c.ConceptType = *req.ConceptType
	}
	if req.Taxability != nil {
		c.Taxability = *req.Taxability
	}
	if req.CalculationType != nil {
		c.CalculationType = *req.CalculationType
	}
	if req.BaseConceptID != nil {
		c.BaseConceptID = req.BaseConceptID
	}
	if req.Active != nil {
		c.Active = *req.Active
	}
	if req.SortOrder != nil {
		c.SortOrder = *req.SortOrder
	}
	if err := s.repo.UpdateConcept(ctx, c); err != nil {
		return nil, fmt.Errorf("update concept: save: %w", err)
	}
	return c, nil
}

func (s *PayrollService) GetConcept(ctx context.Context, companyID, id string) (*domain.PayrollConcept, error) {
	c, err := s.repo.GetConcept(ctx, companyID, id)
	if err != nil {
		return nil, fmt.Errorf("get concept: %w", err)
	}
	return c, nil
}

func (s *PayrollService) ListConcepts(ctx context.Context, companyID string, conceptType, taxability *string, active *bool) ([]domain.PayrollConcept, error) {
	concepts, err := s.repo.ListConcepts(ctx, companyID, repository.ConceptFilter{
		ConceptType: conceptType,
		Taxability:  taxability,
		Active:      active,
	})
	if err != nil {
		return nil, fmt.Errorf("list concepts: %w", err)
	}
	return concepts, nil
}

func (s *PayrollService) CreateRule(ctx context.Context, companyID, userID string, req CreateRuleInput) (*domain.PayrollRule, error) {
	effFrom, err := time.Parse("2006-01-02", req.EffectiveFrom)
	if err != nil {
		return nil, fmt.Errorf("create rule: parse effective_from: %w", err)
	}
	var effTo *time.Time
	if req.EffectiveTo != nil {
		t, err := time.Parse("2006-01-02", *req.EffectiveTo)
		if err != nil {
			return nil, fmt.Errorf("create rule: parse effective_to: %w", err)
		}
		effTo = &t
	}
	if req.Parameters == nil {
		req.Parameters = make(map[string]any)
	}
	rule := &domain.PayrollRule{
		ID:            uuid.NewString(),
		CompanyID:     companyID,
		Country:       "AR",
		ConceptID:     req.ConceptID,
		RuleType:      req.RuleType,
		Formula:       req.Formula,
		Parameters:    req.Parameters,
		Priority:      req.Priority,
		EffectiveFrom: effFrom,
		EffectiveTo:   effTo,
		Version:       1,
		Active:        true,
		CreatedBy:     userID,
	}
	if err := s.repo.CreateRule(ctx, rule); err != nil {
		return nil, fmt.Errorf("create rule: %w", err)
	}
	return rule, nil
}

func (s *PayrollService) UpdateRule(ctx context.Context, companyID, id string, req UpdateRuleInput) (*domain.PayrollRule, error) {
	rule, err := s.repo.GetRule(ctx, companyID, id)
	if err != nil {
		return nil, fmt.Errorf("update rule: get: %w", err)
	}
	if req.RuleType != nil {
		rule.RuleType = *req.RuleType
	}
	if req.Formula != nil {
		rule.Formula = req.Formula
	}
	if req.Parameters != nil {
		rule.Parameters = req.Parameters
	}
	if req.Priority != nil {
		rule.Priority = *req.Priority
	}
	if req.Active != nil {
		rule.Active = *req.Active
	}
	if req.EffectiveTo != nil {
		t, err := time.Parse("2006-01-02", *req.EffectiveTo)
		if err != nil {
			return nil, fmt.Errorf("update rule: parse effective_to: %w", err)
		}
		rule.EffectiveTo = &t
	}
	if err := s.repo.UpdateRule(ctx, rule); err != nil {
		return nil, fmt.Errorf("update rule: save: %w", err)
	}
	return rule, nil
}

func (s *PayrollService) GetRule(ctx context.Context, companyID, id string) (*domain.PayrollRule, error) {
	rule, err := s.repo.GetRule(ctx, companyID, id)
	if err != nil {
		return nil, fmt.Errorf("get rule: %w", err)
	}
	return rule, nil
}

func (s *PayrollService) ListRules(ctx context.Context, companyID string) ([]domain.PayrollRule, error) {
	rules, err := s.repo.ListRules(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("list rules: %w", err)
	}
	return rules, nil
}
