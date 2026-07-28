package payroll

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Service struct {
	repo      *Repository
	calc      *Calculator
	validator *Validator
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo:      repo,
		calc:      NewCalculator(repo),
		validator: NewValidator(),
	}
}

// Periods
func (s *Service) CreatePeriod(ctx context.Context, companyID string, req *CreatePeriodRequest) (*PayrollPeriod, error) {
	if err := s.validator.ValidateCreatePeriod(req); err != nil {
		return nil, err
	}
	return s.repo.CreatePeriod(ctx, companyID, req)
}

func (s *Service) GetPeriod(ctx context.Context, companyID, id string) (*PayrollPeriod, error) {
	return s.repo.GetPeriod(ctx, companyID, id)
}

func (s *Service) ListPeriods(ctx context.Context, companyID string) ([]PayrollPeriod, error) {
	return s.repo.ListPeriods(ctx, companyID)
}

func (s *Service) UpdatePeriod(ctx context.Context, companyID, id string, req *UpdatePeriodRequest) (*PayrollPeriod, error) {
	period, err := s.repo.GetPeriod(ctx, companyID, id)
	if err != nil { return nil, err }
	if period.Status != "OPEN" {
		return nil, fmt.Errorf("can only update OPEN periods")
	}
	return s.repo.UpdatePeriod(ctx, companyID, id, req)
}

// Concepts
func (s *Service) CreateConcept(ctx context.Context, companyID string, req *CreateConceptRequest) (*PayrollConcept, error) {
	return s.repo.CreateConcept(ctx, companyID, req)
}

func (s *Service) GetConcept(ctx context.Context, companyID, id string) (*PayrollConcept, error) {
	return s.repo.GetConcept(ctx, companyID, id)
}

func (s *Service) ListConcepts(ctx context.Context, companyID string) ([]PayrollConcept, error) {
	return s.repo.ListConcepts(ctx, companyID)
}

func (s *Service) UpdateConcept(ctx context.Context, companyID, id string, req *UpdateConceptRequest) (*PayrollConcept, error) {
	return s.repo.UpdateConcept(ctx, companyID, id, req)
}

// Compensation
func (s *Service) SetCompensation(ctx context.Context, companyID string, req *SetCompensationRequest) (*EmployeeCompensation, error) {
	if err := s.validator.ValidateCompensation(req); err != nil {
		return nil, err
	}
	return s.repo.SetCompensation(ctx, companyID, req)
}

func (s *Service) GetCompensation(ctx context.Context, companyID, employeeID string) (*EmployeeCompensation, error) {
	return s.repo.GetCompensation(ctx, companyID, employeeID)
}

func (s *Service) GetCompensationHistory(ctx context.Context, companyID, employeeID string) ([]EmployeeCompensation, error) {
	return s.repo.GetCompensationHistory(ctx, companyID, employeeID)
}

// Benefits
func (s *Service) CreateBenefit(ctx context.Context, companyID string, req *CreateBenefitRequest) (*Benefit, error) {
	return s.repo.CreateBenefit(ctx, companyID, req)
}

func (s *Service) ListBenefits(ctx context.Context, companyID string) ([]Benefit, error) {
	return s.repo.ListBenefits(ctx, companyID)
}

func (s *Service) AssignBenefit(ctx context.Context, companyID string, req *AssignBenefitRequest) (*EmployeeBenefit, error) {
	return s.repo.AssignBenefit(ctx, companyID, req)
}

func (s *Service) GetEmployeeBenefits(ctx context.Context, companyID, employeeID string) ([]EmployeeBenefit, error) {
	return s.repo.GetEmployeeBenefits(ctx, companyID, employeeID, time.Now())
}

// Bonuses
func (s *Service) CreateBonus(ctx context.Context, companyID string, req *CreateBonusRequest) (*PayrollBonus, error) {
	return s.repo.CreateBonus(ctx, companyID, req)
}

func (s *Service) ListBonuses(ctx context.Context, companyID string, filters PayrollFilters) ([]PayrollBonus, error) {
	return s.repo.ListBonuses(ctx, companyID, filters)
}

func (s *Service) ApproveBonus(ctx context.Context, companyID, id, approvedBy string) error {
	return s.repo.ApproveBonus(ctx, id, approvedBy)
}

// Advances
func (s *Service) CreateAdvance(ctx context.Context, companyID string, req *CreateAdvanceRequest) (*PayrollAdvance, error) {
	return s.repo.CreateAdvance(ctx, companyID, req)
}

func (s *Service) ListAdvances(ctx context.Context, companyID string, filters PayrollFilters) ([]PayrollAdvance, error) {
	return s.repo.ListAdvances(ctx, companyID, filters)
}

func (s *Service) ApproveAdvance(ctx context.Context, companyID, id, approvedBy string) error {
	return s.repo.ApproveAdvance(ctx, id, approvedBy)
}

// Deductions
func (s *Service) CreateDeduction(ctx context.Context, companyID string, req *CreateDeductionRequest) (*PayrollDeduction, error) {
	return s.repo.CreateDeduction(ctx, companyID, req)
}

// Calculate
func (s *Service) CalculatePeriod(ctx context.Context, companyID, periodID string) (*PayrollReview, error) {
	period, err := s.repo.GetPeriod(ctx, companyID, periodID)
	if err != nil { return nil, err }

	if err := s.validator.ValidateCanCalculate(period); err != nil {
		return nil, err
	}

	s.repo.DeleteItemsByPeriod(ctx, periodID)

	review, err := s.calc.CalculateForAllEmployees(ctx, companyID, period)
	if err != nil { return nil, err }

	s.repo.SetPeriodCalculated(ctx, companyID, periodID)

	return review, nil
}

// Review
func (s *Service) GetReview(ctx context.Context, companyID, periodID string) (*PayrollReview, error) {
	period, err := s.repo.GetPeriod(ctx, companyID, periodID)
	if err != nil { return nil, err }

	employees, _ := s.repo.GetActiveEmployees(ctx, companyID)
	calculated, _ := s.repo.CountItemsByPeriod(ctx, periodID)

	review := &PayrollReview{
		TotalEmployees: len(employees),
		Calculated:     calculated,
		Pending:        len(employees) - calculated,
	}

	if period.Status == "REVIEW" || period.Status == "APPROVED" || period.Status == "CLOSED" {
		gross, _ := s.calc.GetPeriodTotalGross(ctx, periodID)
		net, _ := s.calc.GetPeriodTotalNet(ctx, periodID)
		review.TotalGross = gross
		review.TotalNet = net
	}

	return review, nil
}

// Approve
func (s *Service) ApprovePeriod(ctx context.Context, companyID, periodID, approvedBy string) error {
	period, err := s.repo.GetPeriod(ctx, companyID, periodID)
	if err != nil { return err }

	if err := s.validator.ValidateCanApprove(period); err != nil {
		return err
	}

	return s.repo.ApprovePeriod(ctx, companyID, periodID, approvedBy)
}

// Close
func (s *Service) ClosePeriod(ctx context.Context, companyID, periodID, closedBy string) error {
	period, err := s.repo.GetPeriod(ctx, companyID, periodID)
	if err != nil { return err }

	if err := s.validator.ValidateCanClose(period); err != nil {
		return err
	}

	items, _ := s.repo.ListItemsByPeriod(ctx, periodID)
	snapshot := map[string]interface{}{
		"period_id":   periodID,
		"company_id":  companyID,
		"created_at":  time.Now(),
		"items_count": len(items),
	}
	data, _ := json.Marshal(snapshot)
	s.repo.CreateSnapshot(ctx, companyID, periodID, data)

	return s.repo.ClosePeriod(ctx, companyID, periodID, closedBy)
}

// Adjustments
func (s *Service) CreateAdjustment(ctx context.Context, companyID, periodID string, req *CreateAdjustmentRequest) (*PayrollAdjustment, error) {
	period, err := s.repo.GetPeriod(ctx, companyID, periodID)
	if err != nil { return nil, err }
	if period.Status == "CLOSED" {
		return nil, fmt.Errorf("cannot adjust closed period")
	}
	return s.repo.CreateAdjustment(ctx, companyID, periodID, req)
}

// Snapshot
func (s *Service) GetSnapshot(ctx context.Context, companyID, periodID string) (*PayrollSnapshot, error) {
	return s.repo.GetSnapshot(ctx, companyID, periodID)
}

// Ledger
func (s *Service) GetLedger(ctx context.Context, companyID, periodID string) ([]PayrollLedgerEntry, error) {
	return s.repo.GetLedgerForPeriod(ctx, companyID, periodID)
}

// Employee view
func (s *Service) GetEmployeePayroll(ctx context.Context, companyID, employeeID, periodID string) (*EmployeePayrollSummary, error) {
	compensation, _ := s.repo.GetCompensation(ctx, companyID, employeeID)
	name, _ := s.repo.GetEmployeeName(ctx, employeeID)

	summary := &EmployeePayrollSummary{
		EmployeeID:   employeeID,
		EmployeeName: name,
	}
	if compensation != nil {
		summary.BaseSalary = compensation.BaseAmount
		summary.Currency = compensation.Currency
	}

	items, _ := s.repo.ListItemsByPeriodAndEmployee(ctx, periodID, employeeID)
	summary.Items = items

	for _, item := range items {
		switch {
		case item.Amount > 0 && item.ConceptCode == "OVERTIME":
			summary.OvertimeAmount += item.Amount
		case item.Amount > 0 && item.ConceptCode == "BONUS":
			summary.BonusTotal += item.Amount
		case item.Amount > 0 && item.ConceptCode == "BENEFIT":
			summary.BenefitTotal += item.Amount
		case item.Amount < 0 && item.ConceptCode == "DEDUCTION":
			summary.DeductionTotal += -item.Amount
		case item.Amount < 0 && item.ConceptCode == "ADVANCE":
			summary.AdvanceTotal += -item.Amount
		}

		if item.Amount > 0 {
			summary.GrossTotal += item.Amount
		}
		summary.NetTotal += item.Amount
	}

	return summary, nil
}

// Dashboard
func (s *Service) GetDashboard(ctx context.Context, companyID, periodID string) (*PayrollDashboard, error) {
	items, _ := s.repo.ListItemsByPeriod(ctx, periodID)
	employees, _ := s.repo.GetActiveEmployees(ctx, companyID)

	dash := &PayrollDashboard{
		TotalEmployees: len(employees),
	}

	for _, item := range items {
		switch {
		case item.ConceptCode == "BASE":
			dash.TotalGross += item.Amount
		case item.ConceptCode == "OVERTIME":
			dash.TotalOvertime += item.Amount
			dash.TotalGross += item.Amount
		case item.ConceptCode == "BONUS":
			dash.TotalBonuses += item.Amount
			dash.TotalGross += item.Amount
		case item.ConceptCode == "BENEFIT":
			dash.TotalBenefits += item.Amount
		case item.ConceptCode == "DEDUCTION":
			dash.TotalDeductions += -item.Amount
		case item.ConceptCode == "ADVANCE":
			dash.TotalAdvances += -item.Amount
		}
		dash.TotalNet += item.Amount
	}

	return dash, nil
}
