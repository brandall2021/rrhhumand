package workflow

import (
	"fmt"

	"github.com/rrhhumand/api/internal/performance/domain"
)

type StateMachine struct {
	transitions map[domain.CycleStatus][]domain.CycleStatus
}

func NewCycleStateMachine() *StateMachine {
	return &StateMachine{
		transitions: map[domain.CycleStatus][]domain.CycleStatus{
			domain.CycleStatusDraft:       {domain.CycleStatusOpen},
			domain.CycleStatusOpen:        {domain.CycleStatusInProgress, domain.CycleStatusCancelled},
			domain.CycleStatusInProgress:  {domain.CycleStatusReview, domain.CycleStatusClosed, domain.CycleStatusCancelled},
			domain.CycleStatusReview:      {domain.CycleStatusCalibration, domain.CycleStatusClosed, domain.CycleStatusCancelled},
			domain.CycleStatusCalibration: {domain.CycleStatusClosed, domain.CycleStatusCancelled},
			domain.CycleStatusClosed:      {},
			domain.CycleStatusCancelled:   {},
		},
	}
}

func (sm *StateMachine) CanTransition(from, to domain.CycleStatus) bool {
	allowed, ok := sm.transitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

type EvaluationStateMachine struct{}

func NewEvaluationStateMachine() *EvaluationStateMachine {
	return &EvaluationStateMachine{}
}

func (sm *EvaluationStateMachine) CanSubmit(status domain.EvaluationStatus) bool {
	return status == domain.EvaluationStatusDraft || status == domain.EvaluationStatusReopened
}

func (sm *EvaluationStateMachine) CanApprove(status domain.EvaluationStatus) bool {
	return status == domain.EvaluationStatusSubmitted
}

func (sm *EvaluationStateMachine) CanReopen(status domain.EvaluationStatus) bool {
	return status == domain.EvaluationStatusSubmitted || status == domain.EvaluationStatusApproved
}

func (sm *EvaluationStateMachine) CanLock(status domain.EvaluationStatus) bool {
	return status == domain.EvaluationStatusApproved
}

type PlanStateMachine struct{}

func NewPlanStateMachine() *PlanStateMachine {
	return &PlanStateMachine{}
}

func (sm *PlanStateMachine) CanActivate(status domain.PlanStatus) bool {
	return status == domain.PlanStatusDraft
}

func (sm *PlanStateMachine) CanComplete(status domain.PlanStatus) bool {
	return status == domain.PlanStatusActive
}

func (sm *PlanStateMachine) CanCancel(status domain.PlanStatus) bool {
	return status == domain.PlanStatusDraft || status == domain.PlanStatusActive
}

func ValidateCycleTransition(from, to domain.CycleStatus) error {
	sm := NewCycleStateMachine()
	if !sm.CanTransition(from, to) {
		return fmt.Errorf("transición inválida de %s a %s", from, to)
	}
	return nil
}

func ValidateEvaluationTransition(current domain.EvaluationStatus, target domain.EvaluationStatus) error {
	valid := map[domain.EvaluationStatus][]domain.EvaluationStatus{
		domain.EvaluationStatusDraft:     {domain.EvaluationStatusSubmitted},
		domain.EvaluationStatusReopened:  {domain.EvaluationStatusSubmitted},
		domain.EvaluationStatusSubmitted: {domain.EvaluationStatusApproved, domain.EvaluationStatusReopened},
		domain.EvaluationStatusApproved:  {domain.EvaluationStatusLocked, domain.EvaluationStatusReopened},
		domain.EvaluationStatusLocked:    {},
	}
	allowed, ok := valid[current]
	if !ok {
		return fmt.Errorf("estado de evaluación inválido: %s", current)
	}
	for _, s := range allowed {
		if s == target {
			return nil
		}
	}
	return fmt.Errorf("transición inválida de evaluación: %s -> %s", current, target)
}
