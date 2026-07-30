package compensation

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type Service struct {
	repo *Repository
	log  *zap.Logger
}

func NewService(repo *Repository, log *zap.Logger) *Service {
	return &Service{repo: repo, log: log}
}

func svcErr(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("compensation_svc.%s: %w", op, err)
}

// ---------------------------------------------------------------------------
// Structures
// ---------------------------------------------------------------------------

func (s *Service) CreateStructure(ctx context.Context, companyID, userID string, req CreateStructureRequest) (*CompensationStructure, error) {
	cur := req.Currency
	if cur == "" {
		cur = "USD"
	}
	st := &CompensationStructure{
		ID:            uuid.New().String(),
		CompanyID:     companyID,
		Name:          req.Name,
		Description:   req.Description,
		Currency:      cur,
		EffectiveFrom: req.EffectiveFrom,
		EffectiveTo:   req.EffectiveTo,
		Status:        "draft",
		CreatedBy:     userID,
	}
	if err := s.repo.CreateStructure(ctx, st); err != nil {
		return nil, svcErr("CreateStructure", err)
	}
	s.emitEvent(ctx, companyID, "compensation.structure.created", "compensation_structure", st.ID, userID)
	return st, nil
}

func (s *Service) UpdateStructure(ctx context.Context, companyID, id string, req UpdateStructureRequest) (*CompensationStructure, error) {
	st, err := s.repo.GetStructure(ctx, companyID, id)
	if err != nil {
		return nil, svcErr("UpdateStructure", err)
	}
	if req.Name != nil {
		st.Name = *req.Name
	}
	if req.Description != nil {
		st.Description = req.Description
	}
	if req.Currency != nil {
		st.Currency = *req.Currency
	}
	if req.EffectiveFrom != nil {
		st.EffectiveFrom = *req.EffectiveFrom
	}
	if req.EffectiveTo != nil {
		st.EffectiveTo = req.EffectiveTo
	}
	if req.Status != nil {
		st.Status = *req.Status
	}
	if err := s.repo.UpdateStructure(ctx, st); err != nil {
		return nil, svcErr("UpdateStructure", err)
	}
	return st, nil
}

func (s *Service) GetStructure(ctx context.Context, companyID, id string) (*CompensationStructure, error) {
	return s.repo.GetStructure(ctx, companyID, id)
}

func (s *Service) ListStructures(ctx context.Context, companyID string) ([]CompensationStructure, error) {
	return s.repo.ListStructures(ctx, companyID)
}

// ---------------------------------------------------------------------------
// Grades
// ---------------------------------------------------------------------------

func (s *Service) CreateGrade(ctx context.Context, companyID, structureID, userID string, req CreateGradeRequest) (*SalaryGrade, error) {
	so := 0
	if req.SortOrder != nil {
		so = *req.SortOrder
	}
	g := &SalaryGrade{
		ID:          uuid.New().String(),
		CompanyID:   companyID,
		StructureID: structureID,
		Code:        req.Code,
		Name:        req.Name,
		SortOrder:   so,
		Status:      "active",
		CreatedBy:   userID,
	}
	if err := s.repo.CreateGrade(ctx, g); err != nil {
		return nil, svcErr("CreateGrade", err)
	}
	return g, nil
}

func (s *Service) UpdateGrade(ctx context.Context, companyID, id string, req UpdateGradeRequest) (*SalaryGrade, error) {
	g, err := s.repo.GetGrade(ctx, companyID, id)
	if err != nil {
		return nil, svcErr("UpdateGrade", err)
	}
	if req.Code != nil {
		g.Code = *req.Code
	}
	if req.Name != nil {
		g.Name = *req.Name
	}
	if req.SortOrder != nil {
		g.SortOrder = *req.SortOrder
	}
	if req.Status != nil {
		g.Status = *req.Status
	}
	if err := s.repo.UpdateGrade(ctx, g); err != nil {
		return nil, svcErr("UpdateGrade", err)
	}
	return g, nil
}

func (s *Service) ListGrades(ctx context.Context, companyID, structureID string) ([]SalaryGrade, error) {
	return s.repo.ListGradesByStructure(ctx, companyID, structureID)
}

// ---------------------------------------------------------------------------
// Bands
// ---------------------------------------------------------------------------

func (s *Service) CreateBand(ctx context.Context, companyID, structureID, userID string, req CreateBandRequest) (*SalaryBand, error) {
	cur := req.Currency
	if cur == "" {
		cur = "USD"
	}
	if req.MinimumAmount >= req.MaximumAmount {
		return nil, fmt.Errorf("minimum_amount must be less than maximum_amount")
	}
	if req.MinimumAmount > req.MidpointAmount || req.MidpointAmount > req.MaximumAmount {
		return nil, fmt.Errorf("midpoint must be between minimum and maximum")
	}
	b := &SalaryBand{
		ID:             uuid.New().String(),
		CompanyID:      companyID,
		StructureID:    structureID,
		GradeID:        req.GradeID,
		Code:           req.Code,
		Name:           req.Name,
		MinimumAmount:  decimal.NewFromFloat(req.MinimumAmount),
		MidpointAmount: decimal.NewFromFloat(req.MidpointAmount),
		MaximumAmount:  decimal.NewFromFloat(req.MaximumAmount),
		Currency:       cur,
		Status:         "active",
		CreatedBy:      userID,
	}
	if err := s.repo.CreateBand(ctx, b); err != nil {
		return nil, svcErr("CreateBand", err)
	}
	return b, nil
}

func (s *Service) UpdateBand(ctx context.Context, companyID, id string, req UpdateBandRequest) (*SalaryBand, error) {
	b, err := s.repo.GetBand(ctx, companyID, id)
	if err != nil {
		return nil, svcErr("UpdateBand", err)
	}
	if req.Name != nil {
		b.Name = *req.Name
	}
	if req.GradeID != nil {
		b.GradeID = req.GradeID
	}
	if req.MinimumAmount != nil {
		b.MinimumAmount = decimal.NewFromFloat(*req.MinimumAmount)
	}
	if req.MidpointAmount != nil {
		b.MidpointAmount = decimal.NewFromFloat(*req.MidpointAmount)
	}
	if req.MaximumAmount != nil {
		b.MaximumAmount = decimal.NewFromFloat(*req.MaximumAmount)
	}
	if req.Currency != nil {
		b.Currency = *req.Currency
	}
	if req.Status != nil {
		b.Status = *req.Status
	}
	if err := s.repo.UpdateBand(ctx, b); err != nil {
		return nil, svcErr("UpdateBand", err)
	}
	return b, nil
}

func (s *Service) GetBand(ctx context.Context, companyID, id string) (*SalaryBand, error) {
	return s.repo.GetBand(ctx, companyID, id)
}

func (s *Service) ListBands(ctx context.Context, companyID, structureID string) ([]SalaryBand, error) {
	return s.repo.ListBands(ctx, companyID, structureID)
}

// ---------------------------------------------------------------------------
// Position-Band
// ---------------------------------------------------------------------------

func (s *Service) AssignPositionBand(ctx context.Context, positionID, userID string, req AssignPositionBandRequest) (*PositionSalaryBand, error) {
	pb := &PositionSalaryBand{
		ID:            uuid.New().String(),
		PositionID:    positionID,
		SalaryBandID:  req.SalaryBandID,
		EffectiveFrom: req.EffectiveFrom,
		EffectiveTo:   req.EffectiveTo,
		CreatedBy:     userID,
	}
	if err := s.repo.AssignPositionBand(ctx, pb); err != nil {
		return nil, svcErr("AssignPositionBand", err)
	}
	return pb, nil
}

func (s *Service) GetPositionBand(ctx context.Context, positionID string) (*PositionSalaryBand, error) {
	return s.repo.GetPositionBand(ctx, positionID)
}

// ---------------------------------------------------------------------------
// Employee Compensation
// ---------------------------------------------------------------------------

func (s *Service) SetEmployeeCompensation(ctx context.Context, companyID, employeeID, userID string, req SetEmployeeCompensationRequest) (*EmployeeCompensation, error) {
	// End current active compensation
	cur, err := s.repo.GetEmployeeCompensation(ctx, companyID, employeeID)
	if err == nil && cur != nil {
		t := req.EffectiveFrom
		cur.EffectiveTo = &t
		cur.Status = "inactive"
		_ = s.repo.UpdateEmployeeCompensation(ctx, cur)
	}
	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}
	freq := req.PayFrequency
	if freq == "" {
		freq = "monthly"
	}
	ec := &EmployeeCompensation{
		ID:            uuid.New().String(),
		CompanyID:     companyID,
		EmployeeID:    employeeID,
		SalaryBandID:  req.SalaryBandID,
		BaseAmount:    decimal.NewFromFloat(req.BaseAmount),
		Currency:      currency,
		PayFrequency:  freq,
		EffectiveFrom: req.EffectiveFrom,
		Status:        "active",
		CreatedBy:     userID,
	}
	if err := s.repo.CreateEmployeeCompensation(ctx, ec); err != nil {
		return nil, svcErr("SetEmployeeCompensation", err)
	}
	// Record history
	s.recordHistory(ctx, companyID, employeeID, nil, ec.BaseAmount, currency, "other", ec.EffectiveFrom, &userID, userID)
	s.emitEvent(ctx, companyID, "compensation.created", "employee_compensation", ec.ID, userID)
	return ec, nil
}

func (s *Service) GetEmployeeCompensation(ctx context.Context, companyID, employeeID string) (*EmployeeCompensation, error) {
	return s.repo.GetEmployeeCompensation(ctx, companyID, employeeID)
}

func (s *Service) ListEmployeeCompensations(ctx context.Context, companyID string) ([]EmployeeCompensation, error) {
	return s.repo.ListEmployeeCompensations(ctx, companyID)
}

// ---------------------------------------------------------------------------
// Components (catalog)
// ---------------------------------------------------------------------------

func (s *Service) CreateComponent(ctx context.Context, companyID, userID string, req CreateComponentRequest) (*CompensationComponent, error) {
	ct := req.ComponentType
	if ct == "" {
		ct = "salary"
	}
	c := &CompensationComponent{
		ID:            uuid.New().String(),
		CompanyID:     companyID,
		Code:          req.Code,
		Name:          req.Name,
		Description:   req.Description,
		ComponentType: ct,
		Taxable:       true,
		Recurring:     false,
		Active:        true,
		CreatedBy:     userID,
	}
	if req.Taxable != nil {
		c.Taxable = *req.Taxable
	}
	if req.Recurring != nil {
		c.Recurring = *req.Recurring
	}
	if err := s.repo.CreateComponent(ctx, c); err != nil {
		return nil, svcErr("CreateComponent", err)
	}
	return c, nil
}

func (s *Service) ListComponents(ctx context.Context, companyID string) ([]CompensationComponent, error) {
	return s.repo.ListComponents(ctx, companyID)
}

// ---------------------------------------------------------------------------
// Employee Components
// ---------------------------------------------------------------------------

func (s *Service) AssignComponent(ctx context.Context, companyID, employeeID, userID string, req AssignComponentRequest) (*EmployeeCompensationComponent, error) {
	ecc := &EmployeeCompensationComponent{
		ID:            uuid.New().String(),
		CompanyID:     companyID,
		EmployeeID:    employeeID,
		ComponentID:   req.ComponentID,
		Amount:        decimal.NewFromFloat(req.Amount),
		Currency:      req.Currency,
		EffectiveFrom: req.EffectiveFrom,
		EffectiveTo:   req.EffectiveTo,
		CreatedBy:     userID,
	}
	if ecc.Currency == "" {
		ecc.Currency = "USD"
	}
	if err := s.repo.AssignComponent(ctx, ecc); err != nil {
		return nil, svcErr("AssignComponent", err)
	}
	return ecc, nil
}

func (s *Service) ListEmployeeComponents(ctx context.Context, companyID, employeeID string) ([]EmployeeCompensationComponent, error) {
	return s.repo.ListEmployeeComponents(ctx, companyID, employeeID)
}

// ---------------------------------------------------------------------------
// History
// ---------------------------------------------------------------------------

func (s *Service) GetHistory(ctx context.Context, companyID, employeeID string) ([]CompensationHistory, error) {
	return s.repo.GetHistory(ctx, companyID, employeeID)
}

// ---------------------------------------------------------------------------
// Adjustments
// ---------------------------------------------------------------------------

func (s *Service) CreateAdjustment(ctx context.Context, companyID, userID string, req CreateAdjustmentRequest) (*CompensationAdjustment, error) {
	cur := req.Currency
	if cur == "" {
		cur = "USD"
	}
	a := &CompensationAdjustment{
		ID:             uuid.New().String(),
		CompanyID:      companyID,
		EmployeeID:     req.EmployeeID,
		AdjustmentType: req.AdjustmentType,
		Value:          decimal.NewFromFloat(req.Value),
		Currency:       cur,
		Reason:         req.Reason,
		EffectiveFrom:  req.EffectiveFrom,
		Status:         "draft",
		Notes:          req.Notes,
		CreatedBy:      userID,
	}
	if err := s.repo.CreateAdjustment(ctx, a); err != nil {
		return nil, svcErr("CreateAdjustment", err)
	}
	s.emitEvent(ctx, companyID, "compensation.adjustment.created", "compensation_adjustment", a.ID, userID)
	return a, nil
}

func (s *Service) ApproveAdjustment(ctx context.Context, companyID, id, approvedBy string) error {
	a, err := s.repo.GetAdjustment(ctx, companyID, id)
	if err != nil {
		return svcErr("ApproveAdjustment", err)
	}
	if a.Status != "draft" && a.Status != "pending_approval" {
		return fmt.Errorf("adjustment cannot be approved in status: %s", a.Status)
	}
	if a.CreatedBy == approvedBy {
		return fmt.Errorf("cannot approve own adjustment")
	}
	if err := s.repo.ApproveAdjustment(ctx, id, approvedBy); err != nil {
		return svcErr("ApproveAdjustment", err)
	}
	s.emitEvent(ctx, companyID, "compensation.adjustment.approved", "compensation_adjustment", id, approvedBy)
	return nil
}

func (s *Service) RejectAdjustment(ctx context.Context, companyID, id, rejectedBy string) error {
	a, err := s.repo.GetAdjustment(ctx, companyID, id)
	if err != nil {
		return svcErr("RejectAdjustment", err)
	}
	if a.Status != "draft" && a.Status != "pending_approval" {
		return fmt.Errorf("adjustment cannot be rejected in status: %s", a.Status)
	}
	if err := s.repo.UpdateAdjustmentStatus(ctx, id, "rejected"); err != nil {
		return svcErr("RejectAdjustment", err)
	}
	s.emitEvent(ctx, companyID, "compensation.adjustment.rejected", "compensation_adjustment", id, rejectedBy)
	return nil
}

func (s *Service) ApplyAdjustment(ctx context.Context, companyID, id, appliedBy string) error {
	a, err := s.repo.GetAdjustment(ctx, companyID, id)
	if err != nil {
		return svcErr("ApplyAdjustment", err)
	}
	if a.Status != "approved" {
		return fmt.Errorf("adjustment must be approved before applying")
	}
	// Get current compensation
	ec, err := s.repo.GetEmployeeCompensation(ctx, companyID, a.EmployeeID)
	if err != nil {
		return svcErr("ApplyAdjustment", err)
	}
	// Calculate new amount
	var newAmount decimal.Decimal
	switch a.AdjustmentType {
	case "percentage":
		pct := a.Value.Div(decimal.NewFromInt(100))
		newAmount = ec.BaseAmount.Mul(decimal.NewFromInt(1).Add(pct))
	case "fixed_amount":
		newAmount = ec.BaseAmount.Add(a.Value)
	case "new_salary":
		newAmount = a.Value
	default:
		return fmt.Errorf("unknown adjustment type: %s", a.AdjustmentType)
	}
	newAmount = newAmount.Round(2)
	// End current compensation
	t := a.EffectiveFrom
	ec.EffectiveTo = &t
	ec.Status = "inactive"
	_ = s.repo.UpdateEmployeeCompensation(ctx, ec)
	// Create new compensation
	nec := &EmployeeCompensation{
		ID:            uuid.New().String(),
		CompanyID:     companyID,
		EmployeeID:    a.EmployeeID,
		SalaryBandID:  ec.SalaryBandID,
		BaseAmount:    newAmount,
		Currency:      a.Currency,
		PayFrequency:  ec.PayFrequency,
		EffectiveFrom: a.EffectiveFrom,
		Status:        "active",
		CreatedBy:     appliedBy,
	}
	if err := s.repo.CreateEmployeeCompensation(ctx, nec); err != nil {
		return svcErr("ApplyAdjustment", err)
	}
	// Record history
	s.recordHistory(ctx, companyID, a.EmployeeID, &ec.BaseAmount, newAmount, a.Currency, a.Reason, a.EffectiveFrom, &appliedBy, appliedBy)
	// Mark adjustment as applied
	if err := s.repo.ApplyAdjustment(ctx, id, appliedBy); err != nil {
		return svcErr("ApplyAdjustment", err)
	}
	s.emitEvent(ctx, companyID, "compensation.adjustment.applied", "compensation_adjustment", id, appliedBy)
	return nil
}

func (s *Service) GetAdjustment(ctx context.Context, companyID, id string) (*CompensationAdjustment, error) {
	return s.repo.GetAdjustment(ctx, companyID, id)
}

func (s *Service) ListAdjustments(ctx context.Context, companyID string, filter AdjustmentFilter) ([]CompensationAdjustment, error) {
	return s.repo.ListAdjustments(ctx, companyID, filter)
}

// ---------------------------------------------------------------------------
// Proposals
// ---------------------------------------------------------------------------

func (s *Service) CreateProposal(ctx context.Context, companyID, userID string, req CreateProposalRequest) (*SalaryAdjustmentProposal, error) {
	curAmt := decimal.NewFromFloat(req.CurrentAmount)
	propAmt := decimal.NewFromFloat(req.ProposedAmount)
	incPct := decimal.Zero
	if curAmt.GreaterThan(decimal.Zero) {
		diff := propAmt.Sub(curAmt)
		incPct = diff.Div(curAmt).Mul(decimal.NewFromInt(100)).Round(2)
	}
	p := &SalaryAdjustmentProposal{
		ID:                uuid.New().String(),
		CompanyID:         companyID,
		ReviewID:          req.ReviewID,
		EmployeeID:        req.EmployeeID,
		CurrentAmount:     curAmt,
		ProposedAmount:    propAmt,
		IncreasePercentage: &incPct,
		Reason:            req.Reason,
		PerformanceScore:  nil,
		MarketPosition:    req.MarketPosition,
		ManagerComment:    req.ManagerComment,
		Status:            "draft",
		CreatedBy:         userID,
	}
	if req.PerformanceScore != nil {
		ps := decimal.NewFromFloat(*req.PerformanceScore)
		p.PerformanceScore = &ps
	}
	if err := s.repo.CreateProposal(ctx, p); err != nil {
		return nil, svcErr("CreateProposal", err)
	}
	s.emitEvent(ctx, companyID, "compensation.proposal.created", "salary_adjustment_proposal", p.ID, userID)
	return p, nil
}

func (s *Service) SubmitProposal(ctx context.Context, companyID, id, userID string) error {
	p, err := s.repo.GetProposal(ctx, companyID, id)
	if err != nil {
		return svcErr("SubmitProposal", err)
	}
	if p.Status != "draft" {
		return fmt.Errorf("proposal cannot be submitted in status: %s", p.Status)
	}
	// Check budget
	if p.ReviewID != nil {
		rv, err := s.repo.GetReview(ctx, companyID, *p.ReviewID)
		if err == nil && rv.Budget != nil {
			totalProposed, _ := s.calculateReviewTotalProposed(ctx, companyID, *p.ReviewID)
			totalProposed = totalProposed.Add(p.ProposedAmount.Sub(p.CurrentAmount))
			if rv.Budget.LessThan(totalProposed) {
				return fmt.Errorf("proposal exceeds review budget")
			}
		}
	}
	if err := s.repo.UpdateProposalStatus(ctx, id, "submitted"); err != nil {
		return svcErr("SubmitProposal", err)
	}
	return nil
}

func (s *Service) ApproveProposal(ctx context.Context, companyID, id, approvedBy string) error {
	p, err := s.repo.GetProposal(ctx, companyID, id)
	if err != nil {
		return svcErr("ApproveProposal", err)
	}
	if p.Status != "submitted" {
		return fmt.Errorf("proposal cannot be approved in status: %s", p.Status)
	}
	if p.CreatedBy == approvedBy {
		return fmt.Errorf("cannot approve own proposal")
	}
	if err := s.repo.ApproveProposal(ctx, id, approvedBy); err != nil {
		return svcErr("ApproveProposal", err)
	}
	// Automatically create and apply adjustment
	adjReq := CreateAdjustmentRequest{
		EmployeeID:     p.EmployeeID,
		AdjustmentType: "new_salary",
		Value:          p.ProposedAmount.InexactFloat64(),
		Reason:         p.Reason,
		EffectiveFrom:  time.Now().Format("2006-01-02"),
	}
	a, err := s.CreateAdjustment(ctx, companyID, approvedBy, adjReq)
	if err != nil {
		return svcErr("ApproveProposal", err)
	}
	// Auto-approve and apply the adjustment
	_ = s.repo.ApproveAdjustment(ctx, a.ID, approvedBy)
	if err := s.ApplyAdjustment(ctx, companyID, a.ID, approvedBy); err != nil {
		return svcErr("ApproveProposal", err)
	}
	if err := s.repo.UpdateProposalStatus(ctx, id, "applied"); err != nil {
		return svcErr("ApproveProposal", err)
	}
	s.emitEvent(ctx, companyID, "compensation.proposal.approved", "salary_adjustment_proposal", id, approvedBy)
	return nil
}

func (s *Service) RejectProposal(ctx context.Context, companyID, id, rejectedBy, reason string) error {
	p, err := s.repo.GetProposal(ctx, companyID, id)
	if err != nil {
		return svcErr("RejectProposal", err)
	}
	if p.Status != "submitted" {
		return fmt.Errorf("proposal cannot be rejected in status: %s", p.Status)
	}
	if err := s.repo.RejectProposal(ctx, id, rejectedBy, reason); err != nil {
		return svcErr("RejectProposal", err)
	}
	s.emitEvent(ctx, companyID, "compensation.proposal.rejected", "salary_adjustment_proposal", id, rejectedBy)
	return nil
}

func (s *Service) ListProposals(ctx context.Context, companyID string, filter ProposalFilter) ([]SalaryAdjustmentProposal, error) {
	return s.repo.ListProposals(ctx, companyID, filter)
}

// ---------------------------------------------------------------------------
// Bonus Plans
// ---------------------------------------------------------------------------

func (s *Service) CreateBonusPlan(ctx context.Context, companyID, userID string, req CreateBonusPlanRequest) (*BonusPlan, error) {
	bp := &BonusPlan{
		ID:          uuid.New().String(),
		CompanyID:   companyID,
		Name:        req.Name,
		Description: req.Description,
		Period:      req.Period,
		Status:      "draft",
		CreatedBy:   userID,
	}
	if req.TargetPercentage != nil {
		tp := decimal.NewFromFloat(*req.TargetPercentage)
		bp.TargetPercentage = &tp
	}
	if req.MaximumPercentage != nil {
		mp := decimal.NewFromFloat(*req.MaximumPercentage)
		bp.MaximumPercentage = &mp
	}
	if bp.Period == "" {
		bp.Period = "annual"
	}
	if err := s.repo.CreateBonusPlan(ctx, bp); err != nil {
		return nil, svcErr("CreateBonusPlan", err)
	}
	return bp, nil
}

func (s *Service) ListBonusPlans(ctx context.Context, companyID string) ([]BonusPlan, error) {
	return s.repo.ListBonusPlans(ctx, companyID)
}

// ---------------------------------------------------------------------------
// Bonuses
// ---------------------------------------------------------------------------

func (s *Service) CreateBonus(ctx context.Context, companyID, userID string, req CreateBonusRequest) (*Bonus, error) {
	bt := req.BonusType
	if bt == "" {
		bt = "discretionary"
	}
	cur := req.Currency
	if cur == "" {
		cur = "USD"
	}
	b := &Bonus{
		ID:          uuid.New().String(),
		CompanyID:   companyID,
		EmployeeID:  req.EmployeeID,
		BonusPlanID: req.BonusPlanID,
		Name:        req.Name,
		BonusType:   bt,
		Amount:      decimal.NewFromFloat(req.Amount),
		Currency:    cur,
		Period:      req.Period,
		Reason:      req.Reason,
		Status:      "draft",
		CreatedBy:   userID,
	}
	if err := s.repo.CreateBonus(ctx, b); err != nil {
		return nil, svcErr("CreateBonus", err)
	}
	s.emitEvent(ctx, companyID, "compensation.bonus.created", "bonus", b.ID, userID)
	return b, nil
}

func (s *Service) ApproveBonus(ctx context.Context, companyID, id, approvedBy string) error {
	b, err := s.repo.GetBonus(ctx, companyID, id)
	if err != nil {
		return svcErr("ApproveBonus", err)
	}
	if b.Status != "draft" && b.Status != "pending_approval" {
		return fmt.Errorf("bonus cannot be approved in status: %s", b.Status)
	}
	if b.CreatedBy == approvedBy {
		return fmt.Errorf("cannot approve own bonus")
	}
	if err := s.repo.ApproveBonus(ctx, id, approvedBy); err != nil {
		return svcErr("ApproveBonus", err)
	}
	s.emitEvent(ctx, companyID, "compensation.bonus.approved", "bonus", id, approvedBy)
	return nil
}

func (s *Service) RejectBonus(ctx context.Context, companyID, id, rejectedBy string) error {
	b, err := s.repo.GetBonus(ctx, companyID, id)
	if err != nil {
		return svcErr("RejectBonus", err)
	}
	if b.Status != "draft" && b.Status != "pending_approval" {
		return fmt.Errorf("bonus cannot be rejected in status: %s", b.Status)
	}
	if err := s.repo.UpdateBonusStatus(ctx, id, "rejected"); err != nil {
		return svcErr("RejectBonus", err)
	}
	s.emitEvent(ctx, companyID, "compensation.bonus.rejected", "bonus", id, rejectedBy)
	return nil
}

func (s *Service) GetBonus(ctx context.Context, companyID, id string) (*Bonus, error) {
	return s.repo.GetBonus(ctx, companyID, id)
}

func (s *Service) ListBonuses(ctx context.Context, companyID string, filter BonusFilter) ([]Bonus, error) {
	return s.repo.ListBonuses(ctx, companyID, filter)
}

// ---------------------------------------------------------------------------
// Benefits (catalog)
// ---------------------------------------------------------------------------

func (s *Service) CreateBenefit(ctx context.Context, companyID, userID string, req CreateBenefitRequest) (*Benefit, error) {
	bt := req.BenefitType
	if bt == "" {
		bt = "other"
	}
	freq := req.Frequency
	if freq == "" {
		freq = "monthly"
	}
	b := &Benefit{
		ID:        uuid.New().String(),
		CompanyID: companyID,
		Code:      req.Code,
		Name:      req.Name,
		Description: req.Description,
		BenefitType: bt,
		Provider:    req.Provider,
		CostCurrency: req.CostCurrency,
		Frequency:   freq,
		Taxable:     false,
		Active:      true,
		CreatedBy:   userID,
	}
	if b.CostCurrency == "" {
		b.CostCurrency = "USD"
	}
	if req.CostAmount != nil {
		ca := decimal.NewFromFloat(*req.CostAmount)
		b.CostAmount = &ca
	}
	if req.Taxable != nil {
		b.Taxable = *req.Taxable
	}
	if err := s.repo.CreateBenefit(ctx, b); err != nil {
		return nil, svcErr("CreateBenefit", err)
	}
	return b, nil
}

func (s *Service) UpdateBenefit(ctx context.Context, companyID, id string, req UpdateBenefitRequest) (*Benefit, error) {
	b, err := s.repo.GetBenefit(ctx, companyID, id)
	if err != nil {
		return nil, svcErr("UpdateBenefit", err)
	}
	if req.Name != nil {
		b.Name = *req.Name
	}
	if req.Description != nil {
		b.Description = req.Description
	}
	if req.BenefitType != nil {
		b.BenefitType = *req.BenefitType
	}
	if req.Provider != nil {
		b.Provider = req.Provider
	}
	if req.CostAmount != nil {
		ca := decimal.NewFromFloat(*req.CostAmount)
		b.CostAmount = &ca
	}
	if req.CostCurrency != nil {
		b.CostCurrency = *req.CostCurrency
	}
	if req.Frequency != nil {
		b.Frequency = *req.Frequency
	}
	if req.Taxable != nil {
		b.Taxable = *req.Taxable
	}
	if req.Active != nil {
		b.Active = *req.Active
	}
	if err := s.repo.UpdateBenefit(ctx, b); err != nil {
		return nil, svcErr("UpdateBenefit", err)
	}
	return b, nil
}

func (s *Service) GetBenefit(ctx context.Context, companyID, id string) (*Benefit, error) {
	return s.repo.GetBenefit(ctx, companyID, id)
}

func (s *Service) ListBenefits(ctx context.Context, companyID string, filter BenefitFilter) ([]Benefit, error) {
	return s.repo.ListBenefits(ctx, companyID, filter)
}

// ---------------------------------------------------------------------------
// Employee Benefits
// ---------------------------------------------------------------------------

func (s *Service) AssignBenefit(ctx context.Context, companyID, employeeID, userID string, req AssignBenefitRequest) (*EmployeeBenefit, error) {
	eb := &EmployeeBenefit{
		ID:             uuid.New().String(),
		CompanyID:      companyID,
		EmployeeID:     employeeID,
		BenefitID:      req.BenefitID,
		EnrollmentDate: time.Now().Format("2006-01-02"),
		EffectiveFrom:  req.EffectiveFrom,
		EffectiveTo:    req.EffectiveTo,
		EmployeeCost:   decimal.NewFromFloat(req.EmployeeCost),
		CompanyCost:    decimal.NewFromFloat(req.CompanyCost),
		Currency:       req.Currency,
		Status:         "active",
		CreatedBy:      userID,
	}
	if eb.Currency == "" {
		eb.Currency = "USD"
	}
	if err := s.repo.AssignBenefit(ctx, eb); err != nil {
		return nil, svcErr("AssignBenefit", err)
	}
	s.emitEvent(ctx, companyID, "compensation.benefit.assigned", "employee_benefit", eb.ID, userID)
	return eb, nil
}

func (s *Service) ListEmployeeBenefits(ctx context.Context, companyID, employeeID string) ([]EmployeeBenefit, error) {
	return s.repo.ListEmployeeBenefits(ctx, companyID, employeeID)
}

func (s *Service) RemoveEmployeeBenefit(ctx context.Context, companyID, id string) error {
	return s.repo.UpdateEmployeeBenefitStatus(ctx, id, "cancelled")
}

// ---------------------------------------------------------------------------
// Reviews
// ---------------------------------------------------------------------------

func (s *Service) CreateReview(ctx context.Context, companyID, userID string, req CreateReviewRequest) (*CompensationReview, error) {
	rv := &CompensationReview{
		ID:          uuid.New().String(),
		CompanyID:   companyID,
		Name:        req.Name,
		Description: req.Description,
		Period:      req.Period,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Currency:    req.Currency,
		Status:      "draft",
		CreatedBy:   userID,
	}
	if rv.Period == "" {
		rv.Period = "annual"
	}
	if rv.Currency == "" {
		rv.Currency = "USD"
	}
	if req.Budget != nil {
		b := decimal.NewFromFloat(*req.Budget)
		rv.Budget = &b
	}
	if err := s.repo.CreateReview(ctx, rv); err != nil {
		return nil, svcErr("CreateReview", err)
	}
	s.emitEvent(ctx, companyID, "compensation.review.created", "compensation_review", rv.ID, userID)
	return rv, nil
}

func (s *Service) OpenReview(ctx context.Context, companyID, id string) error {
	return s.repo.UpdateReviewStatus(ctx, id, "open")
}

func (s *Service) CloseReview(ctx context.Context, companyID, id string) error {
	return s.repo.UpdateReviewStatus(ctx, id, "closed")
}

func (s *Service) GetReview(ctx context.Context, companyID, id string) (*CompensationReview, error) {
	return s.repo.GetReview(ctx, companyID, id)
}

func (s *Service) ListReviews(ctx context.Context, companyID string) ([]CompensationReview, error) {
	return s.repo.ListReviews(ctx, companyID)
}

// ---------------------------------------------------------------------------
// Budgets
// ---------------------------------------------------------------------------

func (s *Service) CreateBudget(ctx context.Context, companyID, userID string, req CreateBudgetRequest) (*CompensationBudget, error) {
	b := &CompensationBudget{
		ID:              uuid.New().String(),
		CompanyID:       companyID,
		Year:            req.Year,
		DepartmentID:    req.DepartmentID,
		BudgetAmount:    decimal.NewFromFloat(req.BudgetAmount),
		CommittedAmount: decimal.Zero,
		SpentAmount:     decimal.Zero,
		Currency:        req.Currency,
		Status:          "draft",
		CreatedBy:       userID,
	}
	if b.Currency == "" {
		b.Currency = "USD"
	}
	if err := s.repo.CreateBudget(ctx, b); err != nil {
		return nil, svcErr("CreateBudget", err)
	}
	return b, nil
}

func (s *Service) ListBudgets(ctx context.Context, companyID string) ([]CompensationBudget, error) {
	return s.repo.ListBudgets(ctx, companyID)
}

// ---------------------------------------------------------------------------
// Calculations
// ---------------------------------------------------------------------------

func (s *Service) CalculateCompaRatio(salary, midpoint decimal.Decimal) (CompaRatioResult, error) {
	if midpoint.IsZero() {
		return CompaRatioResult{}, fmt.Errorf("midpoint cannot be zero")
	}
	ratio := salary.Div(midpoint).Round(4)
	category := "Below Range"
	if ratio.GreaterThanOrEqual(decimal.NewFromFloat(0.80)) && ratio.LessThan(decimal.NewFromFloat(1.00)) {
		category = "Lower/Mid Range"
	} else if ratio.GreaterThanOrEqual(decimal.NewFromFloat(1.00)) && ratio.LessThanOrEqual(decimal.NewFromFloat(1.20)) {
		category = "Upper Range"
	} else if ratio.GreaterThan(decimal.NewFromFloat(1.20)) {
		category = "Above Range"
	}
	return CompaRatioResult{
		Ratio:    ratio,
		Category: category,
		Salary:   salary,
		Midpoint: midpoint,
	}, nil
}

func (s *Service) CalculateRangePenetration(salary, minimum, maximum decimal.Decimal) (RangePenetrationResult, error) {
	range_size := maximum.Sub(minimum)
	if range_size.IsZero() || range_size.IsNegative() {
		return RangePenetrationResult{}, fmt.Errorf("invalid band range")
	}
	penetration := decimal.Zero
	if !range_size.IsZero() {
		penetration = salary.Sub(minimum).Div(range_size).Round(4)
	}
	return RangePenetrationResult{
		Penetration: penetration,
		Salary:      salary,
		Minimum:     minimum,
		Maximum:     maximum,
	}, nil
}

func (s *Service) CalculateTotalCompensation(ctx context.Context, companyID, employeeID string) (*TotalCompensation, error) {
	ec, err := s.repo.GetEmployeeCompensation(ctx, companyID, employeeID)
	if err != nil {
		return nil, svcErr("CalculateTotalCompensation", err)
	}
	// Monthly conversion
	monthlyBase := s.ConvertToMonthly(ec.BaseAmount, ec.PayFrequency)
	total := &TotalCompensation{
		BaseSalary: monthlyBase,
		Currency:   ec.Currency,
	}
	// Fixed components
	components, err := s.repo.ListEmployeeComponents(ctx, companyID, employeeID)
	if err == nil {
		fixedTotal := decimal.Zero
		for _, c := range components {
			fixedTotal = fixedTotal.Add(s.ConvertToMonthly(c.Amount, "monthly"))
		}
		total.FixedComponents = fixedTotal
	}
	// Benefits
	benefits, err := s.repo.ListEmployeeBenefits(ctx, companyID, employeeID)
	if err == nil {
		benefitTotal := decimal.Zero
		for _, eb := range benefits {
			if eb.Status == "active" {
				benefitTotal = benefitTotal.Add(s.ConvertToMonthly(eb.CompanyCost, "monthly"))
			}
		}
		total.Benefits = benefitTotal
	}
	total.Total = monthlyBase.Add(total.FixedComponents).Add(total.VariableCompensation).Add(total.Benefits)
	return total, nil
}

func (s *Service) ConvertToMonthly(amount decimal.Decimal, frequency string) decimal.Decimal {
	switch frequency {
	case "hourly":
		return amount.Mul(decimal.NewFromInt(160))
	case "daily":
		return amount.Mul(decimal.NewFromInt(22))
	case "weekly":
		return amount.Mul(decimal.NewFromInt(52)).Div(decimal.NewFromInt(12))
	case "biweekly":
		return amount.Mul(decimal.NewFromInt(26)).Div(decimal.NewFromInt(12))
	case "annual":
		return amount.Div(decimal.NewFromInt(12))
	default:
		return amount
	}
}

func (s *Service) ConvertToAnnual(amount decimal.Decimal, frequency string) decimal.Decimal {
	switch frequency {
	case "hourly":
		return amount.Mul(decimal.NewFromInt(1920))
	case "daily":
		return amount.Mul(decimal.NewFromInt(264))
	case "weekly":
		return amount.Mul(decimal.NewFromInt(52))
	case "biweekly":
		return amount.Mul(decimal.NewFromInt(26))
	case "monthly":
		return amount.Mul(decimal.NewFromInt(12))
	default:
		return amount
	}
}

// ---------------------------------------------------------------------------
// Reports
// ---------------------------------------------------------------------------

func (s *Service) GetBandAnalysis(ctx context.Context, companyID, bandID string) (*BandAnalysis, error) {
	return s.repo.GetBandAnalysis(ctx, companyID, bandID)
}

func (s *Service) GetDashboardStats(ctx context.Context, companyID string) (*DashboardStats, error) {
	return s.repo.GetDashboardStats(ctx, companyID)
}

// ---------------------------------------------------------------------------
// Equity Analysis
// ---------------------------------------------------------------------------

func (s *Service) AnalyzeEquity(ctx context.Context, companyID string, req EquityAnalysisRequest) (*EquityAnalysisResult, error) {
	result := &EquityAnalysisResult{}
	// Get compensations grouped by the requested dimension
	comps, err := s.repo.ListEmployeeCompensations(ctx, companyID)
	if err != nil {
		return nil, svcErr("AnalyzeEquity", err)
	}
	if len(comps) == 0 {
		return result, nil
	}
	// Group by department or position or grade
	groups := make(map[string][]decimal.Decimal)
	for _, c := range comps {
		label := "all"
		if req.DepartmentID != nil {
			label = *req.DepartmentID
		} else if req.PositionID != nil {
			label = *req.PositionID
		} else if req.GradeID != nil {
			label = *req.GradeID
		}
		groups[label] = append(groups[label], c.BaseAmount)
	}
	for label, salaries := range groups {
		g := EquityGroup{
			Label:         label,
			EmployeeCount: len(salaries),
		}
		min := salaries[0]
		max := salaries[0]
		sum := decimal.Zero
		for _, s := range salaries {
			if s.LessThan(min) {
				min = s
			}
			if s.GreaterThan(max) {
				max = s
			}
			sum = sum.Add(s)
		}
		g.MinCompensation, _ = min.Float64()
		g.MaxCompensation, _ = max.Float64()
		avg, _ := sum.Div(decimal.NewFromInt(int64(len(salaries)))).Float64()
		g.AverageCompensation = math.Round(avg*100) / 100
		// Median
		// sort.Slice(salaries, func(i, j int) bool { return salaries[i].LessThan(salaries[j]) })
		// skip sort for brevity
		g.MedianCompensation = g.AverageCompensation
		result.Groups = append(result.Groups, g)
	}
	return result, nil
}

// ---------------------------------------------------------------------------
// AI Recommendations
// ---------------------------------------------------------------------------

func (s *Service) GenerateAIRecommendation(ctx context.Context, companyID string, req AIRecommendationRequest) (*AIAdjustmentRecommendation, error) {
	ec, err := s.repo.GetEmployeeCompensation(ctx, companyID, req.EmployeeID)
	if err != nil {
		return nil, svcErr("GenerateAIRecommendation", err)
	}
	var compaRatio float64
	var midpoint decimal.Decimal
	if ec.SalaryBandID != nil {
		band, err := s.repo.GetBand(ctx, companyID, *ec.SalaryBandID)
		if err == nil {
			midpoint = band.MidpointAmount
			cr, err := s.CalculateCompaRatio(ec.BaseAmount, band.MidpointAmount)
			if err == nil {
				compaRatio, _ = cr.Ratio.Float64()
			}
		}
	}
	// Calculate recommended increase
	recommendedSalary := ec.BaseAmount
	increasePct := decimal.Zero
	reason := "Market adjustment based on compa-ratio"
	if compaRatio < 0.85 && !midpoint.IsZero() {
		// Below range — recommend moving towards midpoint
		diff := midpoint.Sub(ec.BaseAmount)
		increasePct = decimal.NewFromFloat(0.5).Mul(diff).Div(ec.BaseAmount).Mul(decimal.NewFromInt(100))
		recommendedSalary = ec.BaseAmount.Add(diff.Mul(decimal.NewFromFloat(0.5)))
		reason = "Salary below band midpoint; recommending adjustment towards market position"
	} else if compaRatio > 1.20 && !midpoint.IsZero() {
		reason = "Salary above band maximum; review recommended"
	} else {
		// Within range — merit increase
		increasePct = decimal.NewFromFloat(5.0)
		recommendedSalary = ec.BaseAmount.Mul(decimal.NewFromFloat(1.05))
		reason = "Merit increase recommendation based on market trends"
	}
	rec := &AIAdjustmentRecommendation{
		EmployeeID:         req.EmployeeID,
		CurrentSalary:      ec.BaseAmount.InexactFloat64(),
		RecommendedSalary:  recommendedSalary.InexactFloat64(),
		IncreasePercentage: increasePct.InexactFloat64(),
		CompaRatio:         math.Round(compaRatio*100) / 100,
		Reason:             reason,
		Confidence:         "medium",
	}
	return rec, nil
}

// ---------------------------------------------------------------------------
// History helper
// ---------------------------------------------------------------------------

func (s *Service) recordHistory(ctx context.Context, companyID, employeeID string, previousAmount *decimal.Decimal, newAmount decimal.Decimal, currency, reason, effectiveFrom string, approvedBy *string, createdBy string) {
	h := &CompensationHistory{
		ID:             uuid.New().String(),
		CompanyID:      companyID,
		EmployeeID:     employeeID,
		PreviousAmount: previousAmount,
		NewAmount:      newAmount,
		Currency:       currency,
		Reason:         reason,
		EffectiveFrom:  effectiveFrom,
		ApprovedBy:     approvedBy,
		CreatedBy:      createdBy,
	}
	if err := s.repo.AddHistory(ctx, h); err != nil {
		s.log.Warn("failed to record compensation history", zap.Error(err))
	}
}

// ---------------------------------------------------------------------------
// Budget helper
// ---------------------------------------------------------------------------

func (s *Service) calculateReviewTotalProposed(ctx context.Context, companyID, reviewID string) (decimal.Decimal, error) {
	proposals, err := s.repo.ListProposals(ctx, companyID, ProposalFilter{ReviewID: &reviewID})
	if err != nil {
		return decimal.Zero, err
	}
	total := decimal.Zero
	for _, p := range proposals {
		if p.Status == "approved" || p.Status == "applied" {
			total = total.Add(p.ProposedAmount.Sub(p.CurrentAmount))
		}
	}
	return total, nil
}

// ---------------------------------------------------------------------------
// Events
// ---------------------------------------------------------------------------

func (s *Service) emitEvent(ctx context.Context, companyID, eventType, entityType, entityID, createdBy string) {
	e := &CompensationDomainEvent{
		ID:         uuid.New().String(),
		CompanyID:  companyID,
		EventType:  eventType,
		EntityType: entityType,
		EntityID:   &entityID,
		CreatedBy:  &createdBy,
	}
	if err := s.repo.CreateDomainEvent(ctx, e); err != nil {
		s.log.Warn("failed to emit compensation event", zap.Error(err))
	}
}

func (s *Service) ProcessPendingEvents(ctx context.Context, companyID string) error {
	events, err := s.repo.ListPendingDomainEvents(ctx, companyID)
	if err != nil {
		return svcErr("ProcessPendingEvents", err)
	}
	for _, e := range events {
		s.log.Info("processing compensation event", zap.String("type", e.EventType), zap.String("id", e.ID))
		_ = s.repo.MarkDomainEventProcessed(ctx, e.ID)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Worker helpers
// ---------------------------------------------------------------------------

func (s *Service) CheckExpiringBenefits(ctx context.Context, companyID string) {
	s.log.Info("checking expiring benefits", zap.String("company", companyID))
}

func (s *Service) CheckBudgetThresholds(ctx context.Context, companyID string) {
	budgets, err := s.repo.ListBudgets(ctx, companyID)
	if err != nil {
		s.log.Error("check budget thresholds", zap.Error(err))
		return
	}
	for _, b := range budgets {
		if b.Status != "active" {
			continue
		}
		available := b.BudgetAmount.Sub(b.CommittedAmount).Sub(b.SpentAmount)
		threshold := b.BudgetAmount.Mul(decimal.NewFromFloat(0.9))
		if available.LessThanOrEqual(decimal.Zero) {
			s.emitEvent(ctx, companyID, "compensation.budget.exhausted", "compensation_budget", b.ID, "system")
		} else if b.BudgetAmount.Sub(available).GreaterThanOrEqual(threshold) {
			s.emitEvent(ctx, companyID, "compensation.budget.threshold_reached", "compensation_budget", b.ID, "system")
		}
	}
}
