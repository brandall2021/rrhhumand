package scheduling

import (
	"context"
	"fmt"
	"time"

)

type Service struct {
	repo      *Repository
	resolver  *ScheduleResolver
	generator *CalendarGenerator
	conflicts *ConflictDetector
}

func NewService(repo *Repository) *Service {
	resolver := NewScheduleResolver(repo)
	generator := NewCalendarGenerator(repo, resolver)
	conflicts := NewConflictDetector(repo)
	return &Service{
		repo:      repo,
		resolver:  resolver,
		generator: generator,
		conflicts: conflicts,
	}
}

// Work Schedules
func (s *Service) CreateSchedule(ctx context.Context, companyID string, req *CreateScheduleRequest) (*WorkSchedule, error) {
	return s.repo.CreateSchedule(ctx, companyID, req)
}

func (s *Service) GetSchedule(ctx context.Context, companyID, id string) (*WorkSchedule, error) {
	return s.repo.GetSchedule(ctx, companyID, id)
}

func (s *Service) ListSchedules(ctx context.Context, companyID string) ([]WorkSchedule, error) {
	return s.repo.ListSchedules(ctx, companyID)
}

func (s *Service) UpdateSchedule(ctx context.Context, companyID, id string, req *UpdateScheduleRequest) (*WorkSchedule, error) {
	return s.repo.UpdateSchedule(ctx, companyID, id, req)
}

func (s *Service) DeleteSchedule(ctx context.Context, companyID, id string) error {
	return s.repo.DeleteSchedule(ctx, companyID, id)
}

// Shifts
func (s *Service) CreateShift(ctx context.Context, companyID string, req *CreateShiftRequest) (*Shift, error) {
	return s.repo.CreateShift(ctx, companyID, req)
}

func (s *Service) GetShift(ctx context.Context, companyID, id string) (*Shift, error) {
	return s.repo.GetShift(ctx, companyID, id)
}

func (s *Service) ListShifts(ctx context.Context, companyID string) ([]Shift, error) {
	return s.repo.ListShifts(ctx, companyID)
}

func (s *Service) UpdateShift(ctx context.Context, companyID, id string, req *UpdateShiftRequest) (*Shift, error) {
	return s.repo.UpdateShift(ctx, companyID, id, req)
}

func (s *Service) DeleteShift(ctx context.Context, companyID, id string) error {
	return s.repo.DeleteShift(ctx, companyID, id)
}

// Assignments
func (s *Service) AssignSchedule(ctx context.Context, companyID, employeeID string, req *AssignScheduleRequest) (*EmployeeScheduleAssignment, error) {
	return s.repo.AssignSchedule(ctx, companyID, employeeID, req)
}

func (s *Service) GetEmployeeSchedule(ctx context.Context, employeeID string, date time.Time) (*EmployeeScheduleAssignment, error) {
	return s.repo.GetEmployeeScheduleAssignment(ctx, employeeID, date)
}

func (s *Service) AssignShift(ctx context.Context, companyID, employeeID string, req *AssignShiftRequest) (*EmployeeShiftAssignment, error) {
	workDate, _ := time.Parse("2006-01-02", req.WorkDate)

	conflicts, err := s.conflicts.CheckShiftAssignment(ctx, employeeID, req.ShiftID, workDate)
	if err != nil {
		return nil, err
	}
	if len(conflicts) > 0 {
		return nil, fmt.Errorf("conflict: %s", conflicts[0].Description)
	}

	return s.repo.AssignShift(ctx, companyID, employeeID, req)
}

func (s *Service) GetEmployeeShift(ctx context.Context, employeeID string, date time.Time) (*EmployeeShiftAssignment, error) {
	return s.repo.GetEmployeeShiftAssignment(ctx, employeeID, date)
}

// Rotation Templates
func (s *Service) CreateRotationTemplate(ctx context.Context, companyID string, req *CreateRotationTemplateRequest) (*RotationTemplate, error) {
	return s.repo.CreateRotationTemplate(ctx, companyID, req)
}

func (s *Service) GetRotationTemplate(ctx context.Context, companyID, id string) (*RotationTemplate, error) {
	return s.repo.GetRotationTemplate(ctx, companyID, id)
}

func (s *Service) ListRotationTemplates(ctx context.Context, companyID string) ([]RotationTemplate, error) {
	return s.repo.ListRotationTemplates(ctx, companyID)
}

func (s *Service) AssignRotation(ctx context.Context, companyID, employeeID string, req *AssignRotationRequest) (*EmployeeRotationAssignment, error) {
	return s.repo.AssignRotation(ctx, companyID, employeeID, req)
}

// Calendar
func (s *Service) GenerateCalendar(ctx context.Context, companyID, employeeID, fromStr, toStr string) ([]EmployeeWorkCalendar, error) {
	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		return nil, fmt.Errorf("invalid from date")
	}
	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		return nil, fmt.Errorf("invalid to date")
	}

	if to.Sub(from).Hours()/24 > 90 {
		return nil, fmt.Errorf("date range cannot exceed 90 days")
	}

	return s.generator.Generate(ctx, companyID, employeeID, from, to)
}

func (s *Service) ListCalendar(ctx context.Context, companyID string, filters CalendarFilters, page, perPage int) ([]EmployeeWorkCalendar, int64, error) {
	offset := (page - 1) * perPage
	return s.repo.ListCalendarEntries(ctx, companyID, filters, offset, perPage)
}

func (s *Service) GetResolvedSchedule(ctx context.Context, employeeID string, date time.Time) (*ResolvedSchedule, error) {
	return s.resolver.Resolve(ctx, employeeID, date)
}

// Exceptions
func (s *Service) CreateException(ctx context.Context, companyID string, req *CreateExceptionRequest) (*ScheduleException, error) {
	return s.repo.CreateException(ctx, companyID, req)
}

func (s *Service) ListExceptions(ctx context.Context, companyID string) ([]ScheduleException, error) {
	return s.repo.ListExceptions(ctx, companyID)
}

// Shift Swaps
func (s *Service) SwapShift(ctx context.Context, companyID, requesterID string, req *SwapShiftRequest) (*ShiftSwap, error) {
	requesterDate, _ := time.Parse("2006-01-02", req.RequesterDate)
	targetDate, _ := time.Parse("2006-01-02", req.TargetDate)

	conflicts, err := s.conflicts.CheckShiftSwap(ctx, companyID, requesterID, req.TargetEmployeeID, requesterDate, targetDate)
	if err != nil {
		return nil, err
	}
	if len(conflicts) > 0 {
		return nil, fmt.Errorf("conflict: %s", conflicts[0].Description)
	}

	return s.repo.CreateSwap(ctx, companyID, requesterID, req)
}

func (s *Service) ApproveSwap(ctx context.Context, swapID, approvedBy string) error {
	return s.repo.UpdateSwapStatus(ctx, swapID, "APPROVED", approvedBy)
}

func (s *Service) RejectSwap(ctx context.Context, swapID, approvedBy string) error {
	return s.repo.UpdateSwapStatus(ctx, swapID, "REJECTED", approvedBy)
}
