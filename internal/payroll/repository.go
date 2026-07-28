package payroll

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Periods
func (r *Repository) CreatePeriod(ctx context.Context, companyID string, req *CreatePeriodRequest) (*PayrollPeriod, error) {
	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)
	p := &PayrollPeriod{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO payroll_periods (company_id, name, start_date, end_date)
		 VALUES ($1,$2,$3,$4)
		 RETURNING id, company_id, name, start_date, end_date, status, created_at, updated_at`,
		companyID, req.Name, startDate, endDate,
	).Scan(&p.ID, &p.CompanyID, &p.Name, &p.StartDate, &p.EndDate, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repository) GetPeriod(ctx context.Context, companyID, id string) (*PayrollPeriod, error) {
	p := &PayrollPeriod{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, name, start_date, end_date, status, calculated_at, approved_by, approved_at, closed_by, closed_at, created_at, updated_at
		 FROM payroll_periods WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&p.ID, &p.CompanyID, &p.Name, &p.StartDate, &p.EndDate, &p.Status, &p.CalculatedAt,
		&p.ApprovedBy, &p.ApprovedAt, &p.ClosedBy, &p.ClosedAt, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repository) ListPeriods(ctx context.Context, companyID string) ([]PayrollPeriod, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, name, start_date, end_date, status, calculated_at, approved_by, approved_at, closed_by, closed_at, created_at, updated_at
		 FROM payroll_periods WHERE company_id=$1 ORDER BY start_date DESC`, companyID)
	if err != nil { return nil, err }
	defer rows.Close()

	var periods []PayrollPeriod
	for rows.Next() {
		var p PayrollPeriod
		rows.Scan(&p.ID, &p.CompanyID, &p.Name, &p.StartDate, &p.EndDate, &p.Status, &p.CalculatedAt,
			&p.ApprovedBy, &p.ApprovedAt, &p.ClosedBy, &p.ClosedAt, &p.CreatedAt, &p.UpdatedAt)
		periods = append(periods, p)
	}
	return periods, nil
}

func (r *Repository) UpdatePeriod(ctx context.Context, companyID, id string, req *UpdatePeriodRequest) (*PayrollPeriod, error) {
	p := &PayrollPeriod{}
	err := r.pool.QueryRow(ctx,
		`UPDATE payroll_periods SET name=COALESCE($3,name), start_date=COALESCE($4,start_date), end_date=COALESCE($5,end_date), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, name, start_date, end_date, status, created_at, updated_at`,
		companyID, id, req.Name, req.StartDate, req.EndDate,
	).Scan(&p.ID, &p.CompanyID, &p.Name, &p.StartDate, &p.EndDate, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repository) UpdatePeriodStatus(ctx context.Context, companyID, id, status string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE payroll_periods SET status=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, status)
	return err
}

func (r *Repository) SetPeriodCalculated(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE payroll_periods SET status='REVIEW', calculated_at=NOW(), updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id)
	return err
}

func (r *Repository) ApprovePeriod(ctx context.Context, companyID, id, approvedBy string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE payroll_periods SET status='APPROVED', approved_by=$3, approved_at=NOW(), updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, approvedBy)
	return err
}

func (r *Repository) ClosePeriod(ctx context.Context, companyID, id, closedBy string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE payroll_periods SET status='CLOSED', closed_by=$3, closed_at=NOW(), updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, closedBy)
	return err
}

// Concepts
func (r *Repository) CreateConcept(ctx context.Context, companyID string, req *CreateConceptRequest) (*PayrollConcept, error) {
	calcType := "FIXED"
	taxable := false
	if req.CalculationType != nil { calcType = *req.CalculationType }
	if req.Taxable != nil { taxable = *req.Taxable }

	c := &PayrollConcept{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO payroll_concepts (company_id, code, name, type, calculation_type, taxable)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id, company_id, code, name, type, calculation_type, taxable, active, created_at`,
		companyID, req.Code, req.Name, req.Type, calcType, taxable,
	).Scan(&c.ID, &c.CompanyID, &c.Code, &c.Name, &c.Type, &c.CalculationType, &c.Taxable, &c.Active, &c.CreatedAt)
	return c, err
}

func (r *Repository) GetConcept(ctx context.Context, companyID, id string) (*PayrollConcept, error) {
	c := &PayrollConcept{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, code, name, type, calculation_type, taxable, active, created_at
		 FROM payroll_concepts WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&c.ID, &c.CompanyID, &c.Code, &c.Name, &c.Type, &c.CalculationType, &c.Taxable, &c.Active, &c.CreatedAt)
	return c, err
}

func (r *Repository) GetConceptByCode(ctx context.Context, companyID, code string) (*PayrollConcept, error) {
	c := &PayrollConcept{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, code, name, type, calculation_type, taxable, active, created_at
		 FROM payroll_concepts WHERE company_id=$1 AND code=$2`, companyID, code,
	).Scan(&c.ID, &c.CompanyID, &c.Code, &c.Name, &c.Type, &c.CalculationType, &c.Taxable, &c.Active, &c.CreatedAt)
	return c, err
}

func (r *Repository) ListConcepts(ctx context.Context, companyID string) ([]PayrollConcept, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, code, name, type, calculation_type, taxable, active, created_at
		 FROM payroll_concepts WHERE company_id=$1 ORDER BY code`, companyID)
	if err != nil { return nil, err }
	defer rows.Close()

	var concepts []PayrollConcept
	for rows.Next() {
		var c PayrollConcept
		rows.Scan(&c.ID, &c.CompanyID, &c.Code, &c.Name, &c.Type, &c.CalculationType, &c.Taxable, &c.Active, &c.CreatedAt)
		concepts = append(concepts, c)
	}
	return concepts, nil
}

func (r *Repository) UpdateConcept(ctx context.Context, companyID, id string, req *UpdateConceptRequest) (*PayrollConcept, error) {
	c := &PayrollConcept{}
	err := r.pool.QueryRow(ctx,
		`UPDATE payroll_concepts SET name=COALESCE($3,name), type=COALESCE($4,type), calculation_type=COALESCE($5,calculation_type),
		 taxable=COALESCE($6,taxable), active=COALESCE($7,active)
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, code, name, type, calculation_type, taxable, active, created_at`,
		companyID, id, req.Name, req.Type, req.CalculationType, req.Taxable, req.Active,
	).Scan(&c.ID, &c.CompanyID, &c.Code, &c.Name, &c.Type, &c.CalculationType, &c.Taxable, &c.Active, &c.CreatedAt)
	return c, err
}

// Compensation
func (r *Repository) SetCompensation(ctx context.Context, companyID string, req *SetCompensationRequest) (*EmployeeCompensation, error) {
	effectiveFrom, _ := time.Parse("2006-01-02", req.EffectiveFrom)
	currency := "USD"
	if req.Currency != nil { currency = *req.Currency }

	ec := &EmployeeCompensation{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO employee_compensations (company_id, employee_id, base_amount, currency, effective_from, reason)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id, company_id, employee_id, base_amount, currency, effective_from, effective_to, reason, created_at`,
		companyID, req.EmployeeID, req.BaseAmount, currency, effectiveFrom, req.Reason,
	).Scan(&ec.ID, &ec.CompanyID, &ec.EmployeeID, &ec.BaseAmount, &ec.Currency, &ec.EffectiveFrom, &ec.EffectiveTo, &ec.Reason, &ec.CreatedAt)
	return ec, err
}

func (r *Repository) GetCompensation(ctx context.Context, companyID, employeeID string) (*EmployeeCompensation, error) {
	ec := &EmployeeCompensation{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, employee_id, base_amount, currency, effective_from, effective_to, reason, created_at
		 FROM employee_compensations WHERE company_id=$1 AND employee_id=$2
		 AND effective_from <= CURRENT_DATE AND (effective_to IS NULL OR effective_to >= CURRENT_DATE)
		 ORDER BY effective_from DESC LIMIT 1`, companyID, employeeID,
	).Scan(&ec.ID, &ec.CompanyID, &ec.EmployeeID, &ec.BaseAmount, &ec.Currency, &ec.EffectiveFrom, &ec.EffectiveTo, &ec.Reason, &ec.CreatedAt)
	return ec, err
}

func (r *Repository) GetCompensationHistory(ctx context.Context, companyID, employeeID string) ([]EmployeeCompensation, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, employee_id, base_amount, currency, effective_from, effective_to, reason, created_at
		 FROM employee_compensations WHERE company_id=$1 AND employee_id=$2 ORDER BY effective_from DESC`, companyID, employeeID)
	if err != nil { return nil, err }
	defer rows.Close()

	var comps []EmployeeCompensation
	for rows.Next() {
		var ec EmployeeCompensation
		rows.Scan(&ec.ID, &ec.CompanyID, &ec.EmployeeID, &ec.BaseAmount, &ec.Currency, &ec.EffectiveFrom, &ec.EffectiveTo, &ec.Reason, &ec.CreatedAt)
		comps = append(comps, ec)
	}
	return comps, nil
}

// Items
func (r *Repository) CreateItem(ctx context.Context, item *PayrollItem) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO payroll_items (id, payroll_period_id, employee_id, concept_id, quantity, unit_amount, amount, metadata)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		item.ID, item.PayrollPeriodID, item.EmployeeID, item.ConceptID, item.Quantity, item.UnitAmount, item.Amount, nil)
	return err
}

func (r *Repository) DeleteItemsByPeriod(ctx context.Context, periodID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM payroll_items WHERE payroll_period_id=$1`, periodID)
	return err
}

func (r *Repository) ListItemsByPeriod(ctx context.Context, periodID string) ([]PayrollItem, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT pi.id, pi.payroll_period_id, pi.employee_id, pi.concept_id, pc.code, pc.name, pi.quantity, pi.unit_amount, pi.amount, pi.created_at
		 FROM payroll_items pi
		 JOIN payroll_concepts pc ON pi.concept_id=pc.id
		 WHERE pi.payroll_period_id=$1 ORDER BY pi.employee_id, pc.code`, periodID)
	if err != nil { return nil, err }
	defer rows.Close()

	var items []PayrollItem
	for rows.Next() {
		var item PayrollItem
		rows.Scan(&item.ID, &item.PayrollPeriodID, &item.EmployeeID, &item.ConceptID, &item.ConceptCode, &item.ConceptName, &item.Quantity, &item.UnitAmount, &item.Amount, &item.CreatedAt)
		items = append(items, item)
	}
	return items, nil
}

func (r *Repository) ListItemsByPeriodAndEmployee(ctx context.Context, periodID, employeeID string) ([]PayrollItem, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT pi.id, pi.payroll_period_id, pi.employee_id, pi.concept_id, pc.code, pc.name, pi.quantity, pi.unit_amount, pi.amount, pi.created_at
		 FROM payroll_items pi
		 JOIN payroll_concepts pc ON pi.concept_id=pc.id
		 WHERE pi.payroll_period_id=$1 AND pi.employee_id=$2 ORDER BY pc.code`, periodID, employeeID)
	if err != nil { return nil, err }
	defer rows.Close()

	var items []PayrollItem
	for rows.Next() {
		var item PayrollItem
		rows.Scan(&item.ID, &item.PayrollPeriodID, &item.EmployeeID, &item.ConceptID, &item.ConceptCode, &item.ConceptName, &item.Quantity, &item.UnitAmount, &item.Amount, &item.CreatedAt)
		items = append(items, item)
	}
	return items, nil
}

// Benefits
func (r *Repository) CreateBenefit(ctx context.Context, companyID string, req *CreateBenefitRequest) (*Benefit, error) {
	b := &Benefit{}
	defAmount := 0.0
	if req.DefaultAmount != nil { defAmount = *req.DefaultAmount }
	err := r.pool.QueryRow(ctx,
		`INSERT INTO benefits (company_id, code, name, description, benefit_type, default_amount)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id, company_id, code, name, description, benefit_type, default_amount, active, created_at`,
		companyID, req.Code, req.Name, req.Description, req.BenefitType, defAmount,
	).Scan(&b.ID, &b.CompanyID, &b.Code, &b.Name, &b.Description, &b.BenefitType, &b.DefaultAmount, &b.Active, &b.CreatedAt)
	return b, err
}

func (r *Repository) ListBenefits(ctx context.Context, companyID string) ([]Benefit, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, code, name, description, benefit_type, default_amount, active, created_at
		 FROM benefits WHERE company_id=$1 ORDER BY code`, companyID)
	if err != nil { return nil, err }
	defer rows.Close()

	var benefits []Benefit
	for rows.Next() {
		var b Benefit
		rows.Scan(&b.ID, &b.CompanyID, &b.Code, &b.Name, &b.Description, &b.BenefitType, &b.DefaultAmount, &b.Active, &b.CreatedAt)
		benefits = append(benefits, b)
	}
	return benefits, nil
}

func (r *Repository) AssignBenefit(ctx context.Context, companyID string, req *AssignBenefitRequest) (*EmployeeBenefit, error) {
	effectiveFrom, _ := time.Parse("2006-01-02", req.EffectiveFrom)
	currency := "USD"
	if req.Currency != nil { currency = *req.Currency }

	eb := &EmployeeBenefit{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO employee_benefits (company_id, employee_id, benefit_id, amount, currency, effective_from)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id, company_id, employee_id, benefit_id, amount, currency, effective_from, effective_to, status, created_at`,
		companyID, req.EmployeeID, req.BenefitID, req.Amount, currency, effectiveFrom,
	).Scan(&eb.ID, &eb.CompanyID, &eb.EmployeeID, &eb.BenefitID, &eb.Amount, &eb.Currency, &eb.EffectiveFrom, &eb.EffectiveTo, &eb.Status, &eb.CreatedAt)
	return eb, err
}

func (r *Repository) GetEmployeeBenefits(ctx context.Context, companyID, employeeID string, date time.Time) ([]EmployeeBenefit, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT eb.id, eb.company_id, eb.employee_id, b.name, eb.benefit_id, eb.amount, eb.currency, eb.effective_from, eb.effective_to, eb.status, eb.created_at
		 FROM employee_benefits eb
		 JOIN benefits b ON eb.benefit_id=b.id
		 WHERE eb.company_id=$1 AND eb.employee_id=$2 AND eb.status='ACTIVE'
		 AND eb.effective_from<=$3 AND (eb.effective_to IS NULL OR eb.effective_to>=$3)`, companyID, employeeID, date)
	if err != nil { return nil, err }
	defer rows.Close()

	var benefits []EmployeeBenefit
	for rows.Next() {
		var eb EmployeeBenefit
		rows.Scan(&eb.ID, &eb.CompanyID, &eb.EmployeeID, &eb.BenefitName, &eb.BenefitID, &eb.Amount, &eb.Currency, &eb.EffectiveFrom, &eb.EffectiveTo, &eb.Status, &eb.CreatedAt)
		benefits = append(benefits, eb)
	}
	return benefits, nil
}

// Bonuses
func (r *Repository) CreateBonus(ctx context.Context, companyID string, req *CreateBonusRequest) (*PayrollBonus, error) {
	currency := "USD"
	if req.Currency != nil { currency = *req.Currency }
	b := &PayrollBonus{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO payroll_bonuses (company_id, employee_id, bonus_type, amount, currency, reason, period_start, period_end)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 RETURNING id, company_id, employee_id, bonus_type, amount, currency, reason, period_start, period_end, status, created_at`,
		companyID, req.EmployeeID, req.BonusType, req.Amount, currency, req.Reason, req.PeriodStart, req.PeriodEnd,
	).Scan(&b.ID, &b.CompanyID, &b.EmployeeID, &b.BonusType, &b.Amount, &b.Currency, &b.Reason, &b.PeriodStart, &b.PeriodEnd, &b.Status, &b.CreatedAt)
	return b, err
}

func (r *Repository) ListBonuses(ctx context.Context, companyID string, filters PayrollFilters) ([]PayrollBonus, error) {
	query := `SELECT pb.id, pb.company_id, pb.employee_id, e.first_name||' '||e.last_name, pb.bonus_type, pb.amount, pb.currency, pb.reason, pb.period_start, pb.period_end, pb.status, pb.approved_by, pb.approved_at, pb.created_at
		 FROM payroll_bonuses pb LEFT JOIN employees e ON pb.employee_id=e.id WHERE pb.company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if filters.EmployeeID != "" {
		query += fmt.Sprintf(" AND pb.employee_id=$%d", argIdx)
		args = append(args, filters.EmployeeID)
		argIdx++
	}
	if filters.Status != "" {
		query += fmt.Sprintf(" AND pb.status=$%d", argIdx)
		args = append(args, filters.Status)
		argIdx++
	}
	query += " ORDER BY pb.created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil { return nil, err }
	defer rows.Close()

	var bonuses []PayrollBonus
	for rows.Next() {
		var b PayrollBonus
		rows.Scan(&b.ID, &b.CompanyID, &b.EmployeeID, &b.EmployeeName, &b.BonusType, &b.Amount, &b.Currency, &b.Reason, &b.PeriodStart, &b.PeriodEnd, &b.Status, &b.ApprovedBy, &b.ApprovedAt, &b.CreatedAt)
		bonuses = append(bonuses, b)
	}
	return bonuses, nil
}

func (r *Repository) ApproveBonus(ctx context.Context, id, approvedBy string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE payroll_bonuses SET status='APPROVED', approved_by=$2, approved_at=NOW() WHERE id=$1 AND status='PENDING'`,
		id, approvedBy)
	return err
}

func (r *Repository) GetBonusesForPeriod(ctx context.Context, companyID, employeeID string, start, end time.Time) ([]PayrollBonus, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, employee_id, bonus_type, amount, currency, reason, period_start, period_end, status, created_at
		 FROM payroll_bonuses WHERE company_id=$1 AND employee_id=$2 AND status='APPROVED'
		 AND ((period_start IS NULL AND period_end IS NULL) OR (period_start<=$3 AND (period_end IS NULL OR period_end>=$4)))`,
		companyID, employeeID, end, start)
	if err != nil { return nil, err }
	defer rows.Close()

	var bonuses []PayrollBonus
	for rows.Next() {
		var b PayrollBonus
		rows.Scan(&b.ID, &b.CompanyID, &b.EmployeeID, &b.BonusType, &b.Amount, &b.Currency, &b.Reason, &b.PeriodStart, &b.PeriodEnd, &b.Status, &b.CreatedAt)
		bonuses = append(bonuses, b)
	}
	return bonuses, nil
}

// Advances
func (r *Repository) CreateAdvance(ctx context.Context, companyID string, req *CreateAdvanceRequest) (*PayrollAdvance, error) {
	reqDate, _ := time.Parse("2006-01-02", req.RequestDate)
	currency := "USD"
	if req.Currency != nil { currency = *req.Currency }

	a := &PayrollAdvance{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO payroll_advances (company_id, employee_id, amount, currency, request_date, reason)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id, company_id, employee_id, amount, currency, request_date, reason, status, created_at`,
		companyID, req.EmployeeID, req.Amount, currency, reqDate, req.Reason,
	).Scan(&a.ID, &a.CompanyID, &a.EmployeeID, &a.Amount, &a.Currency, &a.RequestDate, &a.Reason, &a.Status, &a.CreatedAt)
	return a, err
}

func (r *Repository) ListAdvances(ctx context.Context, companyID string, filters PayrollFilters) ([]PayrollAdvance, error) {
	query := `SELECT pa.id, pa.company_id, pa.employee_id, e.first_name||' '||e.last_name, pa.amount, pa.currency, pa.request_date, pa.reason, pa.status, pa.approved_by, pa.approved_at, pa.created_at
		 FROM payroll_advances pa LEFT JOIN employees e ON pa.employee_id=e.id WHERE pa.company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if filters.EmployeeID != "" {
		query += fmt.Sprintf(" AND pa.employee_id=$%d", argIdx)
		args = append(args, filters.EmployeeID)
		argIdx++
	}
	if filters.Status != "" {
		query += fmt.Sprintf(" AND pa.status=$%d", argIdx)
		args = append(args, filters.Status)
		argIdx++
	}
	query += " ORDER BY pa.created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil { return nil, err }
	defer rows.Close()

	var advances []PayrollAdvance
	for rows.Next() {
		var a PayrollAdvance
		rows.Scan(&a.ID, &a.CompanyID, &a.EmployeeID, &a.EmployeeName, &a.Amount, &a.Currency, &a.RequestDate, &a.Reason, &a.Status, &a.ApprovedBy, &a.ApprovedAt, &a.CreatedAt)
		advances = append(advances, a)
	}
	return advances, nil
}

func (r *Repository) ApproveAdvance(ctx context.Context, id, approvedBy string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE payroll_advances SET status='APPROVED', approved_by=$2, approved_at=NOW() WHERE id=$1 AND status='PENDING'`,
		id, approvedBy)
	return err
}

func (r *Repository) GetAdvancesForPeriod(ctx context.Context, companyID, employeeID string, start, end time.Time) ([]PayrollAdvance, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, employee_id, amount, currency, request_date, reason, status, created_at
		 FROM payroll_advances WHERE company_id=$1 AND employee_id=$2 AND status='APPROVED'
		 AND request_date>=$3 AND request_date<=$4`,
		companyID, employeeID, start, end)
	if err != nil { return nil, err }
	defer rows.Close()

	var advances []PayrollAdvance
	for rows.Next() {
		var a PayrollAdvance
		rows.Scan(&a.ID, &a.CompanyID, &a.EmployeeID, &a.Amount, &a.Currency, &a.RequestDate, &a.Reason, &a.Status, &a.CreatedAt)
		advances = append(advances, a)
	}
	return advances, nil
}

// Deductions
func (r *Repository) CreateDeduction(ctx context.Context, companyID string, req *CreateDeductionRequest) (*PayrollDeduction, error) {
	currency := "USD"
	if req.Currency != nil { currency = *req.Currency }

	d := &PayrollDeduction{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO payroll_deductions (company_id, employee_id, concept, amount, currency, reason, period_start, period_end)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 RETURNING id, company_id, employee_id, concept, amount, currency, reason, period_start, period_end, status, created_at`,
		companyID, req.EmployeeID, req.Concept, req.Amount, currency, req.Reason, req.PeriodStart, req.PeriodEnd,
	).Scan(&d.ID, &d.CompanyID, &d.EmployeeID, &d.Concept, &d.Amount, &d.Currency, &d.Reason, &d.PeriodStart, &d.PeriodEnd, &d.Status, &d.CreatedAt)
	return d, err
}

func (r *Repository) GetDeductionsForPeriod(ctx context.Context, companyID, employeeID string, start, end time.Time) ([]PayrollDeduction, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, employee_id, concept, amount, currency, reason, period_start, period_end, status, created_at
		 FROM payroll_deductions WHERE company_id=$1 AND employee_id=$2 AND status='ACTIVE'
		 AND (period_start IS NULL OR period_start<=$3) AND (period_end IS NULL OR period_end>=$4)`,
		companyID, employeeID, end, start)
	if err != nil { return nil, err }
	defer rows.Close()

	var deductions []PayrollDeduction
	for rows.Next() {
		var d PayrollDeduction
		rows.Scan(&d.ID, &d.CompanyID, &d.EmployeeID, &d.Concept, &d.Amount, &d.Currency, &d.Reason, &d.PeriodStart, &d.PeriodEnd, &d.Status, &d.CreatedAt)
		deductions = append(deductions, d)
	}
	return deductions, nil
}

// Overtime integration
func (r *Repository) GetApprovedOvertimeForPeriod(ctx context.Context, companyID, employeeID string, start, end time.Time) (float64, error) {
	var total float64
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(approved_minutes),0) FROM overtime_records
		 WHERE company_id=$1 AND employee_id=$2 AND status='APPROVED' AND work_date>=$3 AND work_date<=$4`,
		companyID, employeeID, start, end).Scan(&total)
	return total, err
}

func (r *Repository) GetOvertimeRate(ctx context.Context, companyID, employeeID string) (float64, error) {
	var rate float64
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(op.weekend_multiplier, 1.5) FROM overtime_policies op WHERE op.company_id=$1 AND op.status='ACTIVE' LIMIT 1`,
		companyID).Scan(&rate)
	if err != nil { rate = 1.5 }
	return rate, nil
}

// Ledger
func (r *Repository) CreateLedgerEntry(ctx context.Context, entry *PayrollLedgerEntry) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO payroll_ledger (id, company_id, payroll_period_id, employee_id, transaction_type, concept_code, amount, balance_after, description, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		entry.ID, entry.CompanyID, entry.PayrollPeriodID, entry.EmployeeID, entry.TransactionType,
		entry.ConceptCode, entry.Amount, entry.BalanceAfter, entry.Description, entry.CreatedBy)
	return err
}

func (r *Repository) GetLedgerForPeriod(ctx context.Context, companyID, periodID string) ([]PayrollLedgerEntry, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, payroll_period_id, employee_id, transaction_type, concept_code, amount, balance_after, description, created_by, created_at
		 FROM payroll_ledger WHERE company_id=$1 AND payroll_period_id=$2 ORDER BY employee_id, created_at`, companyID, periodID)
	if err != nil { return nil, err }
	defer rows.Close()

	var entries []PayrollLedgerEntry
	for rows.Next() {
		var e PayrollLedgerEntry
		rows.Scan(&e.ID, &e.CompanyID, &e.PayrollPeriodID, &e.EmployeeID, &e.TransactionType, &e.ConceptCode, &e.Amount, &e.BalanceAfter, &e.Description, &e.CreatedBy, &e.CreatedAt)
		entries = append(entries, e)
	}
	return entries, nil
}

// Adjustments
func (r *Repository) CreateAdjustment(ctx context.Context, companyID, periodID string, req *CreateAdjustmentRequest) (*PayrollAdjustment, error) {
	adjType := "CREDIT"
	if req.Type != nil { adjType = *req.Type }

	a := &PayrollAdjustment{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO payroll_adjustments (company_id, payroll_period_id, employee_id, amount, reason, type)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id, company_id, payroll_period_id, employee_id, amount, reason, type, created_at`,
		companyID, periodID, req.EmployeeID, req.Amount, req.Reason, adjType,
	).Scan(&a.ID, &a.CompanyID, &a.PayrollPeriodID, &a.EmployeeID, &a.Amount, &a.Reason, &a.Type, &a.CreatedAt)
	return a, err
}

func (r *Repository) GetAdjustmentsForPeriod(ctx context.Context, companyID, periodID string) ([]PayrollAdjustment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, payroll_period_id, employee_id, amount, reason, type, created_at
		 FROM payroll_adjustments WHERE company_id=$1 AND payroll_period_id=$2 ORDER BY created_at`, companyID, periodID)
	if err != nil { return nil, err }
	defer rows.Close()

	var adjustments []PayrollAdjustment
	for rows.Next() {
		var a PayrollAdjustment
		rows.Scan(&a.ID, &a.CompanyID, &a.PayrollPeriodID, &a.EmployeeID, &a.Amount, &a.Reason, &a.Type, &a.CreatedAt)
		adjustments = append(adjustments, a)
	}
	return adjustments, nil
}

// Snapshot
func (r *Repository) CreateSnapshot(ctx context.Context, companyID, periodID string, data []byte) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO payroll_snapshots (company_id, payroll_period_id, snapshot_data) VALUES ($1,$2,$3)`,
		companyID, periodID, data)
	return err
}

func (r *Repository) GetSnapshot(ctx context.Context, companyID, periodID string) (*PayrollSnapshot, error) {
	s := &PayrollSnapshot{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, payroll_period_id, snapshot_data, created_at
		 FROM payroll_snapshots WHERE company_id=$1 AND payroll_period_id=$2`, companyID, periodID,
	).Scan(&s.ID, &s.CompanyID, &s.PayrollPeriodID, &s.SnapshotData, &s.CreatedAt)
	return s, err
}

// Helpers
func (r *Repository) GetActiveEmployees(ctx context.Context, companyID string) ([]string, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id FROM employees WHERE company_id=$1 AND status='active'`, companyID)
	if err != nil { return nil, err }
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		rows.Scan(&id)
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *Repository) GetEmployeeName(ctx context.Context, employeeID string) (string, error) {
	var name string
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(first_name||' '||last_name, 'Unknown') FROM employees WHERE id=$1`, employeeID,
	).Scan(&name)
	return name, err
}

func (r *Repository) CountItemsByPeriod(ctx context.Context, periodID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(DISTINCT employee_id) FROM payroll_items WHERE payroll_period_id=$1`, periodID,
	).Scan(&count)
	return count, err
}
