package integration

import "github.com/rrhhumand/api/internal/performance/repository"

type Adapter struct {
	CycleRepo       repository.CycleRepository
	EvaluationRepo  repository.EvaluationRepository
	ObjectiveRepo   repository.ObjectiveRepository
}
