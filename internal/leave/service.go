package leave

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/rrhhumand/api/internal/models"
)

type Service struct {
	repo       *Repository
	calculator *DayCalculator
}

func NewService(repo *Repository, calculator *DayCalculator) *Service {
	return &Service{repo: repo, calculator: calculator}
}

func (s *Service) CreateLeaveType(ctx context.Context, companyID string, req *CreateLeaveTypeRequest) (*models.LeaveType, error) {
	return s.repo.CreateLeaveType(ctx, companyID, req)
}

func (s *Service) GetLeaveType(ctx context.Context, companyID, id string) (*models.LeaveType, error) {
	return s.repo.GetLeaveType(ctx, companyID, id)
}

func (s *Service) ListLeaveTypes(ctx context.Context, companyID string) ([]models.LeaveType, error) {
	return s.repo.ListLeaveTypes(ctx, companyID)
}

func (s *Service) UpdateLeaveType(ctx context.Context, companyID, id string, req *UpdateLeaveTypeRequest) (*models.LeaveType, error) {
	return s.repo.UpdateLeaveType(ctx, companyID, id, req)
}

func (s *Service) DeleteLeaveType(ctx context.Context, companyID, id string) error {
	return s.repo.DeleteLeaveType(ctx, companyID, id)
}

func (s *Service) CreateLeavePolicy(ctx context.Context, companyID string, req *CreateLeavePolicyRequest) (*models.LeavePolicy, error) {
	_, err := s.repo.GetLeaveType(ctx, companyID, req.LeaveTypeID)
	if err != nil {
		return nil, fmt.Errorf("leave type not found")
	}
	return s.repo.CreateLeavePolicy(ctx, companyID, req)
}

func (s *Service) GetLeavePolicy(ctx context.Context, companyID, leaveTypeID string) (*models.LeavePolicy, error) {
	return s.repo.GetLeavePolicy(ctx, companyID, leaveTypeID)
}

func (s *Service) ListLeavePolicies(ctx context.Context, companyID string) ([]models.LeavePolicy, error) {
	return s.repo.ListLeavePolicies(ctx, companyID)
}

func (s *Service) CreateHoliday(ctx context.Context, companyID string, req *CreateHolidayRequest) (*models.Holiday, error) {
	return s.repo.CreateHoliday(ctx, companyID, req)
}

func (s *Service) GetHolidays(ctx context.Context, companyID string, from, to time.Time) ([]models.Holiday, error) {
	return s.repo.GetHolidays(ctx, companyID, from, to)
}

func (s *Service) DeleteHoliday(ctx context.Context, companyID, id string) error {
	return s.repo.DeleteHoliday(ctx, companyID, id)
}

func (s *Service) GetBalances(ctx context.Context, companyID, employeeID string) ([]models.LeaveBalance, error) {
	year := time.Now().Year()
	return s.repo.ListBalances(ctx, companyID, employeeID, year)
}

func (s *Service) AdjustBalance(ctx context.Context, companyID string, req *AdjustBalanceRequest, performedBy string) error {
	tx, err := s.repo.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	if err := s.repo.AdjustBalance(context.Background(), tx, companyID, req.EmployeeID, req.LeaveTypeID, req.Year, req.AdjustmentDays, req.Reason, performedBy); err != nil {
		return err
	}
	return tx.Commit(context.Background())
}

func (s *Service) CreateLeaveRequest(ctx context.Context, companyID, employeeID string, req *CreateLeaveRequestRequest) (*models.LeaveRequest, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format")
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date format")
	}
	if endDate.Before(startDate) {
		return nil, fmt.Errorf("end_date must be after start_date")
	}

	lt, err := s.repo.GetLeaveType(ctx, companyID, req.LeaveTypeID)
	if err != nil {
		return nil, fmt.Errorf("leave type not found")
	}

	policy, err := s.repo.GetLeavePolicy(ctx, companyID, req.LeaveTypeID)
	if err != nil {
		return nil, fmt.Errorf("leave policy not found")
	}

	requestedDays := s.calculator.CalculateBusinessDays(startDate, endDate, policy.UseBusinessDays, companyID)

	if policy.MinimumDaysBeforeRequest > 0 {
		daysUntilStart := int(time.Until(startDate).Hours() / 24)
		if daysUntilStart < policy.MinimumDaysBeforeRequest {
			return nil, fmt.Errorf("request must be made %d days before start date", policy.MinimumDaysBeforeRequest)
		}
	}

	if policy.MaximumDaysPerRequest != nil && requestedDays > *policy.MaximumDaysPerRequest {
		return nil, fmt.Errorf("maximum days per request is %.0f", *policy.MaximumDaysPerRequest)
	}

	// Check overlap
	overlaps, err := s.repo.CheckOverlap(ctx, employeeID, startDate, endDate, "")
	if err != nil {
		return nil, err
	}
	if overlaps {
		return nil, fmt.Errorf("leave request overlaps with existing request")
	}

	// Validate and reserve balance if affects_balance
	status := "PENDING"
	if !lt.RequiresApproval {
		status = "APPROVED"
	}

	year := startDate.Year()

	if lt.AffectsBalance {
		tx, err := s.repo.pool.Begin(ctx)
		if err != nil {
			return nil, err
		}
		defer tx.Rollback(ctx)

		balance, err := s.repo.GetOrCreateBalance(ctx, tx, companyID, employeeID, req.LeaveTypeID, year)
		if err != nil {
			return nil, err
		}

		if !policy.AllowNegativeBalance && balance.AvailableDays < requestedDays {
			return nil, fmt.Errorf("insufficient leave balance: available=%.2f, requested=%.2f", balance.AvailableDays, requestedDays)
		}

		// Reserve
		newReserved := balance.ReservedDays + requestedDays
		_, err = tx.Exec(ctx, `UPDATE leave_balances SET reserved_days=$1, updated_at=NOW() WHERE id=$2`, newReserved, balance.ID)
		if err != nil {
			return nil, err
		}

		if err := tx.Commit(ctx); err != nil {
			return nil, err
		}
	}

	leaveReq := &models.LeaveRequest{
		ID:            uuid.New().String(),
		CompanyID:     companyID,
		EmployeeID:    employeeID,
		LeaveTypeID:   req.LeaveTypeID,
		StartDate:     startDate,
		EndDate:       endDate,
		RequestedDays: requestedDays,
		Reason:        req.Reason,
		Status:        status,
		DocumentID:    req.DocumentID,
	}

	created, err := s.repo.CreateLeaveRequest(ctx, leaveReq)
	if err != nil {
		return nil, err
	}

	// Create history
	s.repo.CreateHistory(ctx, &models.LeaveRequestHistory{
		ID:             uuid.New().String(),
		LeaveRequestID: created.ID,
		Action:         "CREATED",
		NewStatus:      &status,
		PerformedBy:    employeeID,
	})

	// If auto-approved, process approval
	if status == "APPROVED" {
		s.processAutoApproval(ctx, companyID, employeeID, created)
	}

	return created, nil
}

func (s *Service) processAutoApproval(ctx context.Context, companyID, employeeID string, req *models.LeaveRequest) {
	tx, _ := s.repo.pool.Begin(ctx)
	defer tx.Rollback(ctx)

	if err := s.reserveToUsed(ctx, tx, req); err != nil {
		return
	}
	tx.Commit(ctx)
}

func (s *Service) reserveToUsed(ctx context.Context, tx pgx.Tx, req *models.LeaveRequest) error {
	balance, err := s.repo.GetBalanceForUpdate(ctx, tx, req.CompanyID, req.EmployeeID, req.LeaveTypeID, req.StartDate.Year())
	if err != nil {
		return err
	}
	newReserved := balance.ReservedDays - req.RequestedDays
	newUsed := balance.UsedDays + req.RequestedDays
	if newReserved < 0 {
		newReserved = 0
	}
	_, err = tx.Exec(ctx, `UPDATE leave_balances SET reserved_days=$1, used_days=$2, updated_at=NOW() WHERE id=$3`, newReserved, newUsed, balance.ID)
	return err
}

func (s *Service) ApproveRequest(ctx context.Context, companyID, requestID, approverID string, comments *string) error {
	lr, err := s.repo.GetLeaveRequest(ctx, companyID, requestID)
	if err != nil {
		return fmt.Errorf("leave request not found")
	}
	if lr.Status != "PENDING" {
		return fmt.Errorf("request is not in PENDING status")
	}

	// Create approval record
	level, _ := s.repo.GetNextApprovalLevel(ctx, requestID)
	approval := &models.LeaveApproval{
		ID:             uuid.New().String(),
		CompanyID:      companyID,
		LeaveRequestID: requestID,
		ApproverID:     approverID,
		Level:          level,
		Status:         "APPROVED",
		Comments:       comments,
	}

	tx, err := s.repo.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := s.repo.CreateApproval(ctx, approval); err != nil {
		return err
	}

	// Update request status
	if err := s.repo.UpdateLeaveRequestStatus(ctx, tx, requestID, "APPROVED"); err != nil {
		return err
	}

	// Move reserved to used
	if err := s.reserveToUsed(ctx, tx, lr); err != nil {
		return err
	}

	// History
	s.repo.CreateHistory(ctx, &models.LeaveRequestHistory{
		ID:             uuid.New().String(),
		LeaveRequestID: requestID,
		Action:         "APPROVED",
		OldStatus:      stringPtr("PENDING"),
		NewStatus:      stringPtr("APPROVED"),
		PerformedBy:    approverID,
		Comments:       comments,
	})

	return tx.Commit(ctx)
}

func (s *Service) RejectRequest(ctx context.Context, companyID, requestID, approverID string, comments *string) error {
	lr, err := s.repo.GetLeaveRequest(ctx, companyID, requestID)
	if err != nil {
		return fmt.Errorf("leave request not found")
	}
	if lr.Status != "PENDING" {
		return fmt.Errorf("request is not in PENDING status")
	}

	level, _ := s.repo.GetNextApprovalLevel(ctx, requestID)
	approval := &models.LeaveApproval{
		ID:             uuid.New().String(),
		CompanyID:      companyID,
		LeaveRequestID: requestID,
		ApproverID:     approverID,
		Level:          level,
		Status:         "REJECTED",
		Comments:       comments,
	}

	tx, err := s.repo.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := s.repo.CreateApproval(ctx, approval); err != nil {
		return err
	}

	if err := s.repo.UpdateLeaveRequestStatus(ctx, tx, requestID, "REJECTED"); err != nil {
		return err
	}

	// Release reserved days
	balance, err := s.repo.GetBalanceForUpdate(ctx, tx, lr.CompanyID, lr.EmployeeID, lr.LeaveTypeID, lr.StartDate.Year())
	if err == nil {
		newReserved := balance.ReservedDays - lr.RequestedDays
		if newReserved < 0 {
			newReserved = 0
		}
		tx.Exec(ctx, `UPDATE leave_balances SET reserved_days=$1, updated_at=NOW() WHERE id=$2`, newReserved, balance.ID)
	}

	s.repo.CreateHistory(ctx, &models.LeaveRequestHistory{
		ID:             uuid.New().String(),
		LeaveRequestID: requestID,
		Action:         "REJECTED",
		OldStatus:      stringPtr("PENDING"),
		NewStatus:      stringPtr("REJECTED"),
		PerformedBy:    approverID,
		Comments:       comments,
	})

	return tx.Commit(ctx)
}

func (s *Service) CancelRequest(ctx context.Context, companyID, requestID, employeeID string) error {
	lr, err := s.repo.GetLeaveRequest(ctx, companyID, requestID)
	if err != nil {
		return fmt.Errorf("leave request not found")
	}
	if lr.EmployeeID != employeeID {
		return fmt.Errorf("cannot cancel another employee's request")
	}
	if lr.Status != "PENDING" && lr.Status != "APPROVED" {
		return fmt.Errorf("request cannot be cancelled in status %s", lr.Status)
	}

	tx, err := s.repo.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.repo.UpdateLeaveRequestStatus(ctx, tx, requestID, "CANCELLED"); err != nil {
		return err
	}

	// Release reserved or used days
	balance, err := s.repo.GetBalanceForUpdate(ctx, tx, lr.CompanyID, lr.EmployeeID, lr.LeaveTypeID, lr.StartDate.Year())
	if err == nil {
		newReserved := balance.ReservedDays - lr.RequestedDays
		newUsed := balance.UsedDays
		if lr.Status == "APPROVED" {
			newUsed = balance.UsedDays - lr.RequestedDays
		}
		if newReserved < 0 {
			newReserved = 0
		}
		if newUsed < 0 {
			newUsed = 0
		}
		tx.Exec(ctx, `UPDATE leave_balances SET reserved_days=$1, used_days=$2, updated_at=NOW() WHERE id=$3`, newReserved, newUsed, balance.ID)
	}

	s.repo.CreateHistory(ctx, &models.LeaveRequestHistory{
		ID:             uuid.New().String(),
		LeaveRequestID: requestID,
		Action:         "CANCELLED",
		OldStatus:      &lr.Status,
		NewStatus:      stringPtr("CANCELLED"),
		PerformedBy:    employeeID,
	})

	return tx.Commit(ctx)
}

func (s *Service) GetRequest(ctx context.Context, companyID, id string) (*models.LeaveRequest, error) {
	lr, err := s.repo.GetLeaveRequest(ctx, companyID, id)
	if err != nil {
		return nil, err
	}
	lr.Approvals, _ = s.repo.GetPendingApprovalsForRequest(ctx, id)
	return lr, nil
}

func (s *Service) ListRequests(ctx context.Context, companyID string, filters LeaveFilters, offset, limit int) ([]models.LeaveRequest, int64, error) {
	return s.repo.ListLeaveRequests(ctx, companyID, filters, offset, limit)
}

func (s *Service) GetHistory(ctx context.Context, leaveRequestID string) ([]models.LeaveRequestHistory, error) {
	return s.repo.GetHistory(ctx, leaveRequestID)
}

func (s *Service) GetTeamAbsences(ctx context.Context, companyID, managerID string, start, end time.Time) ([]models.LeaveRequest, error) {
	return s.repo.GetTeamAbsences(ctx, companyID, managerID, start, end)
}

func (s *Service) GetCalendar(ctx context.Context, companyID string, filters CalendarFilters) ([]CalendarDay, error) {
	startDate, _ := time.Parse("2006-01-02", filters.DateFrom)
	endDate, _ := time.Parse("2006-01-02", filters.DateTo)
	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, 0, -30)
	}
	endDate = endDate.AddDate(0, 0, 1) // inclusive

	var absences []models.LeaveRequest
	var err error
	if filters.DepartmentID != "" {
		absences, err = s.repo.GetDepartmentAbsences(ctx, companyID, filters.DepartmentID, startDate, endDate)
	} else if filters.EmployeeID != "" {
		absences, _, err = s.repo.ListLeaveRequests(ctx, companyID, LeaveFilters{EmployeeID: filters.EmployeeID}, 0, 1000)
	} else {
		absences, err = s.repo.GetDepartmentAbsences(ctx, companyID, "", startDate, endDate)
	}
	if err != nil {
		return nil, err
	}

	holidays, _ := s.repo.GetHolidays(ctx, companyID, startDate, endDate)
	holidayMap := make(map[time.Time]bool)
	for _, h := range holidays {
		holidayMap[time.Date(h.Date.Year(), h.Date.Month(), h.Date.Day(), 0, 0, 0, 0, time.UTC)] = true
	}

	absenceMap := make(map[time.Time][]CalendarAbsence)
	for _, a := range absences {
		current := a.StartDate
		for !current.After(a.EndDate) {
			absenceMap[current] = append(absenceMap[current], CalendarAbsence{
				EmployeeID:   a.EmployeeID,
				EmployeeName: a.EmployeeName,
				LeaveType:    a.LeaveTypeName,
				Status:       a.Status,
			})
			current = current.AddDate(0, 0, 1)
		}
	}

	var days []CalendarDay
	current := startDate
	for current.Before(endDate) {
		day := CalendarDay{
			Date:      current,
			IsWeekend: current.Weekday() == time.Saturday || current.Weekday() == time.Sunday,
			IsHoliday: holidayMap[current],
			Absences:  absenceMap[current],
		}
		days = append(days, day)
		current = current.AddDate(0, 0, 1)
	}
	return days, nil
}

func stringPtr(s string) *string {
	return &s
}
