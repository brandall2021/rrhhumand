package overtime

import (
	"context"
	"fmt"
	"time"
)

type Service struct {
	repo         *Repository
	calculator   *Calculator
	detector     *Detector
	policy       *PolicyEngine
	validator    *Validator
	balance      *BalanceManager
	ledger       *Ledger
	approval     *ApprovalManager
	compensation *CompensationManager
	request      *RequestManager
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo:         repo,
		calculator:   NewCalculator(),
		detector:     NewDetector(repo),
		policy:       NewPolicyEngine(),
		validator:    NewValidator(),
		balance:      NewBalanceManager(repo),
		ledger:       NewLedger(repo),
		approval:     NewApprovalManager(repo),
		compensation: NewCompensationManager(repo),
		request:      NewRequestManager(repo),
	}
}

// Policies
func (s *Service) CreatePolicy(ctx context.Context, companyID string, req *CreateOvertimePolicyRequest) (*OvertimePolicy, error) {
	if err := s.validator.ValidateCreatePolicy(req); err != nil {
		return nil, err
	}
	return s.repo.CreatePolicy(ctx, companyID, req)
}

func (s *Service) GetPolicy(ctx context.Context, companyID, id string) (*OvertimePolicy, error) {
	return s.repo.GetPolicy(ctx, companyID, id)
}

func (s *Service) ListPolicies(ctx context.Context, companyID string) ([]OvertimePolicy, error) {
	return s.repo.ListPolicies(ctx, companyID)
}

func (s *Service) UpdatePolicy(ctx context.Context, companyID, id string, req *UpdateOvertimePolicyRequest) (*OvertimePolicy, error) {
	return s.repo.UpdatePolicy(ctx, companyID, id, req)
}

func (s *Service) DeletePolicy(ctx context.Context, companyID, id string) error {
	return s.repo.DeletePolicy(ctx, companyID, id)
}

// Overtime Records
func (s *Service) ListRecords(ctx context.Context, companyID string, filters OvertimeFilters) ([]OvertimeRecord, error) {
	return s.repo.ListOvertimeRecords(ctx, companyID, filters)
}

func (s *Service) GetRecord(ctx context.Context, companyID, id string) (*OvertimeRecord, error) {
	return s.repo.GetOvertimeRecord(ctx, companyID, id)
}

// Detection
func (s *Service) DetectOvertime(ctx context.Context, companyID, dateFrom, dateTo string) ([]OvertimeRecord, int, error) {
	from, err := time.Parse("2006-01-02", dateFrom)
	if err != nil { return nil, 0, fmt.Errorf("invalid date_from") }
	to, err := time.Parse("2006-01-02", dateTo)
	if err != nil { return nil, 0, fmt.Errorf("invalid date_to") }

	policy, _ := s.repo.GetActivePolicy(ctx, companyID)
	if policy == nil {
		policy = &OvertimePolicy{
			MaxDailyMinutes:        120,
			MaxWeeklyMinutes:       480,
			MaxMonthlyMinutes:      1920,
			MinimumOvertimeMinutes: 0,
			RoundingMinutes:        1,
			NightStart:             "22:00",
			NightEnd:               "06:00",
			WeekendMultiplier:      1.5,
			HolidayMultiplier:      2.0,
			NightMultiplier:        1.5,
		}
	}

	return s.detector.DetectForDateRange(ctx, companyID, from, to, policy)
}

// Requests
func (s *Service) CreateRequest(ctx context.Context, companyID, employeeID string, req *RequestOvertimeRequest) (*OvertimeRequest, error) {
	if err := s.validator.ValidateRequest(req); err != nil {
		return nil, err
	}
	return s.request.CreateRequest(ctx, companyID, employeeID, req)
}

func (s *Service) ListRequests(ctx context.Context, companyID string, filters OvertimeFilters) ([]OvertimeRequest, error) {
	return s.request.ListRequests(ctx, companyID, filters)
}

func (s *Service) GetRequest(ctx context.Context, companyID, id string) (*OvertimeRequest, error) {
	return s.request.GetRequest(ctx, companyID, id)
}

// Approvals
func (s *Service) ApproveRecord(ctx context.Context, companyID, recordID string, approvedMinutes int, approvedBy string) error {
	return s.approval.Approve(ctx, companyID, recordID, approvedMinutes, approvedBy)
}

func (s *Service) RejectRecord(ctx context.Context, companyID, recordID, reason, rejectedBy string) error {
	return s.approval.Reject(ctx, companyID, recordID, reason, rejectedBy)
}

func (s *Service) ApproveRequest(ctx context.Context, companyID, requestID string, approvedMinutes int, approvedBy string) error {
	return s.approval.ApproveRequest(ctx, companyID, requestID, approvedMinutes, approvedBy)
}

func (s *Service) RejectRequest(ctx context.Context, companyID, requestID, reason, rejectedBy string) error {
	return s.approval.RejectRequest(ctx, companyID, requestID, reason, rejectedBy)
}

// Compensations
func (s *Service) RequestCompensation(ctx context.Context, companyID, employeeID string, req *RequestCompensationRequest) (*CompensationRequest, error) {
	return s.compensation.RequestCompensation(ctx, companyID, employeeID, req)
}

func (s *Service) ApproveCompensation(ctx context.Context, companyID, requestID, approvedBy string) error {
	return s.compensation.ApproveCompensation(ctx, companyID, requestID, approvedBy)
}

func (s *Service) RejectCompensation(ctx context.Context, companyID, requestID, reason, rejectedBy string) error {
	return s.compensation.RejectCompensation(ctx, companyID, requestID, reason, rejectedBy)
}

func (s *Service) CancelCompensation(ctx context.Context, companyID, requestID string) error {
	return s.compensation.CancelCompensation(ctx, companyID, requestID)
}

func (s *Service) ListCompensations(ctx context.Context, companyID string, filters OvertimeFilters) ([]CompensationRequest, error) {
	return s.repo.ListCompensationRequests(ctx, companyID, filters)
}

// Balance
func (s *Service) GetBalance(ctx context.Context, companyID, employeeID string) (*EmployeeTimeBalance, error) {
	return s.balance.GetBalance(ctx, companyID, employeeID)
}

func (s *Service) AdjustBalance(ctx context.Context, companyID, employeeID string, minutes int, reason, createdBy string) error {
	return s.balance.Adjust(ctx, companyID, employeeID, minutes, reason, createdBy)
}

func (s *Service) GetBalanceTransactions(ctx context.Context, companyID, employeeID string) ([]TimeBalanceTransaction, error) {
	return s.balance.GetTransactions(ctx, companyID, employeeID)
}

// Dashboard
func (s *Service) GetDashboard(ctx context.Context, companyID string) (*OvertimeDashboard, error) {
	return s.repo.GetOvertimeDashboard(ctx, companyID)
}

func (s *Service) GetEmployeeDashboard(ctx context.Context, companyID, employeeID string) (*OvertimeDashboard, error) {
	records, err := s.repo.ListOvertimeRecords(ctx, companyID, OvertimeFilters{EmployeeID: employeeID})
	if err != nil { return nil, err }

	dash := &OvertimeDashboard{}
	for _, rec := range records {
		dash.TotalMinutes += rec.OvertimeMinutes
		switch rec.Status {
		case "DETECTED":
			dash.TotalDetected += rec.OvertimeMinutes
		case "PENDING", "REQUESTED", "SUBMITTED":
			dash.TotalPending += rec.OvertimeMinutes
		case "APPROVED":
			dash.TotalApproved += rec.ApprovedMinutes
		case "REJECTED":
			dash.TotalRejected += rec.OvertimeMinutes
		}
		dash.TotalCompensated += rec.CompensatedMinutes
		dash.TotalPaid += rec.PaidMinutes
	}

	balance, err := s.balance.GetBalance(ctx, companyID, employeeID)
	if err == nil {
		dash.BalanceMinutes = balance.BalanceMinutes
	}

	dash.Records = records
	return dash, nil
}
