package application

import (
	"go.uber.org/zap"

	"github.com/rrhhumand/api/internal/payroll/engine"
	"github.com/rrhhumand/api/internal/payroll/repository"
)

type PayrollService struct {
	repo   *repository.Repository
	engine engine.RuleEngine
	log    *zap.Logger
}

func NewPayrollService(repo *repository.Repository, engine engine.RuleEngine, log *zap.Logger) *PayrollService {
	return &PayrollService{
		repo:   repo,
		engine: engine,
		log:    log,
	}
}
