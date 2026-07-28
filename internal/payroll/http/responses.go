package http

import "github.com/rrhhumand/api/internal/payroll/domain"

type PeriodResponse struct {
	Data *domain.PayrollPeriod `json:"data"`
}

type PeriodListResponse struct {
	Data []domain.PayrollPeriod `json:"data"`
}

type RunResponse struct {
	Data *domain.PayrollRun `json:"data"`
}

type RunListResponse struct {
	Data []domain.PayrollRun `json:"data"`
}

type ConceptResponse struct {
	Data *domain.PayrollConcept `json:"data"`
}

type ConceptListResponse struct {
	Data []domain.PayrollConcept `json:"data"`
}

type RuleResponse struct {
	Data *domain.PayrollRule `json:"data"`
}

type RuleListResponse struct {
	Data []domain.PayrollRule `json:"data"`
}

type NoveltyResponse struct {
	Data *domain.PayrollNovelty `json:"data"`
}

type NoveltyListResponse struct {
	Data []domain.PayrollNovelty `json:"data"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
