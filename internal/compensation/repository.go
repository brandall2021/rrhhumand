package compensation

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type Repository struct {
	pool *pgxpool.Pool
	log  *zap.Logger
}

func NewRepository(pool *pgxpool.Pool, log *zap.Logger) *Repository {
	return &Repository{pool: pool, log: log}
}

func repoErr(op string, err error) error {
	return fmt.Errorf("compensation_repo.%s: %w", op, err)
}

// ---------------------------------------------------------------------------
// Structures
// ---------------------------------------------------------------------------

func (r *Repository) CreateStructure(ctx context.Context, s *CompensationStructure) error {
	q := `INSERT INTO compensation_structures (id,company_id,name,description,currency,effective_from,effective_to,status,created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := r.pool.Exec(ctx, q, s.ID, s.CompanyID, s.Name, s.Description, s.Currency, s.EffectiveFrom, s.EffectiveTo, s.Status, s.CreatedBy)
	return repoErr("CreateStructure", err)
}

func (r *Repository) UpdateStructure(ctx context.Context, s *CompensationStructure) error {
	q := `UPDATE compensation_structures SET name=COALESCE($3,name),description=COALESCE($4,description),currency=COALESCE($5,currency),effective_from=COALESCE($6,effective_from),effective_to=COALESCE($7,effective_to),status=COALESCE($8,status),updated_at=NOW() WHERE id=$1 AND company_id=$2`
	_, err := r.pool.Exec(ctx, q, s.ID, s.CompanyID, s.Name, s.Description, s.Currency, s.EffectiveFrom, s.EffectiveTo, s.Status)
	return repoErr("UpdateStructure", err)
}

func (r *Repository) GetStructure(ctx context.Context, companyID, id string) (*CompensationStructure, error) {
	q := `SELECT id,company_id,name,description,currency,effective_from,effective_to,status,created_by,created_at,updated_at FROM compensation_structures WHERE id=$1 AND company_id=$2`
	s := &CompensationStructure{}
	err := r.pool.QueryRow(ctx, q, id, companyID).Scan(&s.ID, &s.CompanyID, &s.Name, &s.Description, &s.Currency, &s.EffectiveFrom, &s.EffectiveTo, &s.Status, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetStructure", err)
	}
	return s, nil
}

func (r *Repository) ListStructures(ctx context.Context, companyID string) ([]CompensationStructure, error) {
	q := `SELECT id,company_id,name,description,currency,effective_from,effective_to,status,created_by,created_at,updated_at FROM compensation_structures WHERE company_id=$1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListStructures", err)
	}
	defer rows.Close()
	var res []CompensationStructure
	for rows.Next() {
		var s CompensationStructure
		if err := rows.Scan(&s.ID, &s.CompanyID, &s.Name, &s.Description, &s.Currency, &s.EffectiveFrom, &s.EffectiveTo, &s.Status, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, repoErr("ListStructures", err)
		}
		res = append(res, s)
	}
	return res, nil
}

// ---------------------------------------------------------------------------
// Grades
// ---------------------------------------------------------------------------

func (r *Repository) CreateGrade(ctx context.Context, g *SalaryGrade) error {
	q := `INSERT INTO salary_grades (id,company_id,structure_id,code,name,sort_order,status,created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := r.pool.Exec(ctx, q, g.ID, g.CompanyID, g.StructureID, g.Code, g.Name, g.SortOrder, g.Status, g.CreatedBy)
	return repoErr("CreateGrade", err)
}

func (r *Repository) UpdateGrade(ctx context.Context, g *SalaryGrade) error {
	q := `UPDATE salary_grades SET code=COALESCE($3,code),name=COALESCE($4,name),sort_order=COALESCE($5,sort_order),status=COALESCE($6,status),updated_at=NOW() WHERE id=$1 AND company_id=$2`
	_, err := r.pool.Exec(ctx, q, g.ID, g.CompanyID, g.Code, g.Name, g.SortOrder, g.Status)
	return repoErr("UpdateGrade", err)
}

func (r *Repository) GetGrade(ctx context.Context, companyID, id string) (*SalaryGrade, error) {
	q := `SELECT id,company_id,structure_id,code,name,sort_order,status,created_by,created_at,updated_at FROM salary_grades WHERE id=$1 AND company_id=$2`
	g := &SalaryGrade{}
	err := r.pool.QueryRow(ctx, q, id, companyID).Scan(&g.ID, &g.CompanyID, &g.StructureID, &g.Code, &g.Name, &g.SortOrder, &g.Status, &g.CreatedBy, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetGrade", err)
	}
	return g, nil
}

func (r *Repository) ListGradesByStructure(ctx context.Context, companyID, structureID string) ([]SalaryGrade, error) {
	q := `SELECT id,company_id,structure_id,code,name,sort_order,status,created_by,created_at,updated_at FROM salary_grades WHERE company_id=$1 AND structure_id=$2 ORDER BY sort_order`
	rows, err := r.pool.Query(ctx, q, companyID, structureID)
	if err != nil {
		return nil, repoErr("ListGradesByStructure", err)
	}
	defer rows.Close()
	var res []SalaryGrade
	for rows.Next() {
		var g SalaryGrade
		if err := rows.Scan(&g.ID, &g.CompanyID, &g.StructureID, &g.Code, &g.Name, &g.SortOrder, &g.Status, &g.CreatedBy, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, repoErr("ListGradesByStructure", err)
		}
		res = append(res, g)
	}
	return res, nil
}

func (r *Repository) ListGrades(ctx context.Context, companyID string) ([]SalaryGrade, error) {
	q := `SELECT id,company_id,structure_id,code,name,sort_order,status,created_by,created_at,updated_at FROM salary_grades WHERE company_id=$1 ORDER BY sort_order`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListGrades", err)
	}
	defer rows.Close()
	var res []SalaryGrade
	for rows.Next() {
		var g SalaryGrade
		if err := rows.Scan(&g.ID, &g.CompanyID, &g.StructureID, &g.Code, &g.Name, &g.SortOrder, &g.Status, &g.CreatedBy, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, repoErr("ListGrades", err)
		}
		res = append(res, g)
	}
	return res, nil
}

// ---------------------------------------------------------------------------
// Bands
// ---------------------------------------------------------------------------

func (r *Repository) CreateBand(ctx context.Context, b *SalaryBand) error {
	q := `INSERT INTO salary_bands (id,company_id,structure_id,grade_id,code,name,minimum_amount,midpoint_amount,maximum_amount,currency,status,created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := r.pool.Exec(ctx, q, b.ID, b.CompanyID, b.StructureID, b.GradeID, b.Code, b.Name, b.MinimumAmount, b.MidpointAmount, b.MaximumAmount, b.Currency, b.Status, b.CreatedBy)
	return repoErr("CreateBand", err)
}

func (r *Repository) UpdateBand(ctx context.Context, b *SalaryBand) error {
	q := `UPDATE salary_bands SET name=COALESCE($3,name),grade_id=COALESCE($4,grade_id),minimum_amount=COALESCE($5,minimum_amount),midpoint_amount=COALESCE($6,midpoint_amount),maximum_amount=COALESCE($7,maximum_amount),currency=COALESCE($8,currency),status=COALESCE($9,status),updated_at=NOW() WHERE id=$1 AND company_id=$2`
	_, err := r.pool.Exec(ctx, q, b.ID, b.CompanyID, b.Name, b.GradeID, b.MinimumAmount, b.MidpointAmount, b.MaximumAmount, b.Currency, b.Status)
	return repoErr("UpdateBand", err)
}

func (r *Repository) GetBand(ctx context.Context, companyID, id string) (*SalaryBand, error) {
	q := `SELECT id,company_id,structure_id,grade_id,code,name,minimum_amount,midpoint_amount,maximum_amount,currency,status,created_by,created_at,updated_at FROM salary_bands WHERE id=$1 AND company_id=$2`
	b := &SalaryBand{}
	err := r.pool.QueryRow(ctx, q, id, companyID).Scan(&b.ID, &b.CompanyID, &b.StructureID, &b.GradeID, &b.Code, &b.Name, &b.MinimumAmount, &b.MidpointAmount, &b.MaximumAmount, &b.Currency, &b.Status, &b.CreatedBy, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetBand", err)
	}
	return b, nil
}

func (r *Repository) ListBands(ctx context.Context, companyID, structureID string) ([]SalaryBand, error) {
	q := `SELECT id,company_id,structure_id,grade_id,code,name,minimum_amount,midpoint_amount,maximum_amount,currency,status,created_by,created_at,updated_at FROM salary_bands WHERE company_id=$1 AND structure_id=$2 ORDER BY code`
	rows, err := r.pool.Query(ctx, q, companyID, structureID)
	if err != nil {
		return nil, repoErr("ListBands", err)
	}
	defer rows.Close()
	var res []SalaryBand
	for rows.Next() {
		var b SalaryBand
		if err := rows.Scan(&b.ID, &b.CompanyID, &b.StructureID, &b.GradeID, &b.Code, &b.Name, &b.MinimumAmount, &b.MidpointAmount, &b.MaximumAmount, &b.Currency, &b.Status, &b.CreatedBy, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, repoErr("ListBands", err)
		}
		res = append(res, b)
	}
	return res, nil
}

func (r *Repository) GetAllBands(ctx context.Context, companyID string) ([]SalaryBand, error) {
	q := `SELECT id,company_id,structure_id,grade_id,code,name,minimum_amount,midpoint_amount,maximum_amount,currency,status,created_by,created_at,updated_at FROM salary_bands WHERE company_id=$1 ORDER BY code`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("GetAllBands", err)
	}
	defer rows.Close()
	var res []SalaryBand
	for rows.Next() {
		var b SalaryBand
		if err := rows.Scan(&b.ID, &b.CompanyID, &b.StructureID, &b.GradeID, &b.Code, &b.Name, &b.MinimumAmount, &b.MidpointAmount, &b.MaximumAmount, &b.Currency, &b.Status, &b.CreatedBy, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, repoErr("GetAllBands", err)
		}
		res = append(res, b)
	}
	return res, nil
}

// ---------------------------------------------------------------------------
// Position-Band
// ---------------------------------------------------------------------------

func (r *Repository) AssignPositionBand(ctx context.Context, pb *PositionSalaryBand) error {
	q := `INSERT INTO position_salary_bands (id,position_id,salary_band_id,effective_from,effective_to,created_by) VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.pool.Exec(ctx, q, pb.ID, pb.PositionID, pb.SalaryBandID, pb.EffectiveFrom, pb.EffectiveTo, pb.CreatedBy)
	return repoErr("AssignPositionBand", err)
}

func (r *Repository) GetPositionBand(ctx context.Context, positionID string) (*PositionSalaryBand, error) {
	q := `SELECT id,position_id,salary_band_id,effective_from,effective_to,created_by,created_at FROM position_salary_bands WHERE position_id=$1 AND (effective_to IS NULL OR effective_to>=CURRENT_DATE) ORDER BY effective_from DESC LIMIT 1`
	pb := &PositionSalaryBand{}
	err := r.pool.QueryRow(ctx, q, positionID).Scan(&pb.ID, &pb.PositionID, &pb.SalaryBandID, &pb.EffectiveFrom, &pb.EffectiveTo, &pb.CreatedBy, &pb.CreatedAt)
	if err != nil {
		return nil, repoErr("GetPositionBand", err)
	}
	return pb, nil
}

func (r *Repository) ListPositionBands(ctx context.Context, positionID string) ([]PositionSalaryBand, error) {
	q := `SELECT id,position_id,salary_band_id,effective_from,effective_to,created_by,created_at FROM position_salary_bands WHERE position_id=$1 ORDER BY effective_from DESC`
	rows, err := r.pool.Query(ctx, q, positionID)
	if err != nil {
		return nil, repoErr("ListPositionBands", err)
	}
	defer rows.Close()
	var res []PositionSalaryBand
	for rows.Next() {
		var pb PositionSalaryBand
		if err := rows.Scan(&pb.ID, &pb.PositionID, &pb.SalaryBandID, &pb.EffectiveFrom, &pb.EffectiveTo, &pb.CreatedBy, &pb.CreatedAt); err != nil {
			return nil, repoErr("ListPositionBands", err)
		}
		res = append(res, pb)
	}
	return res, nil
}

// ---------------------------------------------------------------------------
// Employee Compensation
// ---------------------------------------------------------------------------

func (r *Repository) CreateEmployeeCompensation(ctx context.Context, ec *EmployeeCompensation) error {
	q := `INSERT INTO employee_compensations (id,company_id,employee_id,salary_band_id,base_amount,currency,pay_frequency,effective_from,effective_to,status,created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.pool.Exec(ctx, q, ec.ID, ec.CompanyID, ec.EmployeeID, ec.SalaryBandID, ec.BaseAmount, ec.Currency, ec.PayFrequency, ec.EffectiveFrom, ec.EffectiveTo, ec.Status, ec.CreatedBy)
	return repoErr("CreateEmployeeCompensation", err)
}

func (r *Repository) UpdateEmployeeCompensation(ctx context.Context, ec *EmployeeCompensation) error {
	q := `UPDATE employee_compensations SET base_amount=COALESCE($3,base_amount),salary_band_id=COALESCE($4,salary_band_id),currency=COALESCE($5,currency),pay_frequency=COALESCE($6,pay_frequency),status=COALESCE($7,status),updated_at=NOW() WHERE id=$1 AND company_id=$2`
	_, err := r.pool.Exec(ctx, q, ec.ID, ec.CompanyID, ec.BaseAmount, ec.SalaryBandID, ec.Currency, ec.PayFrequency, ec.Status)
	return repoErr("UpdateEmployeeCompensation", err)
}

func (r *Repository) GetEmployeeCompensation(ctx context.Context, companyID, employeeID string) (*EmployeeCompensation, error) {
	q := `SELECT id,company_id,employee_id,salary_band_id,base_amount,currency,pay_frequency,effective_from,effective_to,status,created_by,created_at,updated_at FROM employee_compensations WHERE company_id=$1 AND employee_id=$2 AND status='active' ORDER BY effective_from DESC LIMIT 1`
	ec := &EmployeeCompensation{}
	err := r.pool.QueryRow(ctx, q, companyID, employeeID).Scan(&ec.ID, &ec.CompanyID, &ec.EmployeeID, &ec.SalaryBandID, &ec.BaseAmount, &ec.Currency, &ec.PayFrequency, &ec.EffectiveFrom, &ec.EffectiveTo, &ec.Status, &ec.CreatedBy, &ec.CreatedAt, &ec.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetEmployeeCompensation", err)
	}
	return ec, nil
}

func (r *Repository) ListEmployeeCompensations(ctx context.Context, companyID string) ([]EmployeeCompensation, error) {
	q := `SELECT id,company_id,employee_id,salary_band_id,base_amount,currency,pay_frequency,effective_from,effective_to,status,created_by,created_at,updated_at FROM employee_compensations WHERE company_id=$1 AND status='active' ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListEmployeeCompensations", err)
	}
	defer rows.Close()
	var res []EmployeeCompensation
	for rows.Next() {
		var ec EmployeeCompensation
		if err := rows.Scan(&ec.ID, &ec.CompanyID, &ec.EmployeeID, &ec.SalaryBandID, &ec.BaseAmount, &ec.Currency, &ec.PayFrequency, &ec.EffectiveFrom, &ec.EffectiveTo, &ec.Status, &ec.CreatedBy, &ec.CreatedAt, &ec.UpdatedAt); err != nil {
			return nil, repoErr("ListEmployeeCompensations", err)
		}
		res = append(res, ec)
	}
	return res, nil
}

// ---------------------------------------------------------------------------
// Components (catalog)
// ---------------------------------------------------------------------------

func (r *Repository) CreateComponent(ctx context.Context, c *CompensationComponent) error {
	q := `INSERT INTO compensation_components (id,company_id,code,name,description,component_type,taxable,recurring,active,created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.pool.Exec(ctx, q, c.ID, c.CompanyID, c.Code, c.Name, c.Description, c.ComponentType, c.Taxable, c.Recurring, c.Active, c.CreatedBy)
	return repoErr("CreateComponent", err)
}

func (r *Repository) UpdateComponent(ctx context.Context, c *CompensationComponent) error {
	q := `UPDATE compensation_components SET name=COALESCE($3,name),description=COALESCE($4,description),component_type=COALESCE($5,component_type),taxable=COALESCE($6,taxable),recurring=COALESCE($7,recurring),active=COALESCE($8,active),updated_at=NOW() WHERE id=$1 AND company_id=$2`
	_, err := r.pool.Exec(ctx, q, c.ID, c.CompanyID, c.Name, c.Description, c.ComponentType, c.Taxable, c.Recurring, c.Active)
	return repoErr("UpdateComponent", err)
}

func (r *Repository) GetComponent(ctx context.Context, companyID, id string) (*CompensationComponent, error) {
	q := `SELECT id,company_id,code,name,description,component_type,taxable,recurring,active,created_by,created_at,updated_at FROM compensation_components WHERE id=$1 AND company_id=$2`
	c := &CompensationComponent{}
	err := r.pool.QueryRow(ctx, q, id, companyID).Scan(&c.ID, &c.CompanyID, &c.Code, &c.Name, &c.Description, &c.ComponentType, &c.Taxable, &c.Recurring, &c.Active, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetComponent", err)
	}
	return c, nil
}

func (r *Repository) ListComponents(ctx context.Context, companyID string) ([]CompensationComponent, error) {
	q := `SELECT id,company_id,code,name,description,component_type,taxable,recurring,active,created_by,created_at,updated_at FROM compensation_components WHERE company_id=$1 ORDER BY code`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListComponents", err)
	}
	defer rows.Close()
	var res []CompensationComponent
	for rows.Next() {
		var c CompensationComponent
		if err := rows.Scan(&c.ID, &c.CompanyID, &c.Code, &c.Name, &c.Description, &c.ComponentType, &c.Taxable, &c.Recurring, &c.Active, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, repoErr("ListComponents", err)
		}
		res = append(res, c)
	}
	return res, nil
}

// ---------------------------------------------------------------------------
// Employee Components
// ---------------------------------------------------------------------------

func (r *Repository) AssignComponent(ctx context.Context, ecc *EmployeeCompensationComponent) error {
	q := `INSERT INTO employee_compensation_components (id,company_id,employee_id,component_id,amount,currency,effective_from,effective_to,created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := r.pool.Exec(ctx, q, ecc.ID, ecc.CompanyID, ecc.EmployeeID, ecc.ComponentID, ecc.Amount, ecc.Currency, ecc.EffectiveFrom, ecc.EffectiveTo, ecc.CreatedBy)
	return repoErr("AssignComponent", err)
}

func (r *Repository) ListEmployeeComponents(ctx context.Context, companyID, employeeID string) ([]EmployeeCompensationComponent, error) {
	q := `SELECT id,company_id,employee_id,component_id,amount,currency,effective_from,effective_to,created_by,created_at,updated_at FROM employee_compensation_components WHERE company_id=$1 AND employee_id=$2 ORDER BY effective_from DESC`
	rows, err := r.pool.Query(ctx, q, companyID, employeeID)
	if err != nil {
		return nil, repoErr("ListEmployeeComponents", err)
	}
	defer rows.Close()
	var res []EmployeeCompensationComponent
	for rows.Next() {
		var ecc EmployeeCompensationComponent
		if err := rows.Scan(&ecc.ID, &ecc.CompanyID, &ecc.EmployeeID, &ecc.ComponentID, &ecc.Amount, &ecc.Currency, &ecc.EffectiveFrom, &ecc.EffectiveTo, &ecc.CreatedBy, &ecc.CreatedAt, &ecc.UpdatedAt); err != nil {
			return nil, repoErr("ListEmployeeComponents", err)
		}
		res = append(res, ecc)
	}
	return res, nil
}

// ---------------------------------------------------------------------------
// History
// ---------------------------------------------------------------------------

func (r *Repository) AddHistory(ctx context.Context, h *CompensationHistory) error {
	q := `INSERT INTO compensation_history (id,company_id,employee_id,previous_amount,new_amount,currency,reason,effective_from,approved_by,notes,created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.pool.Exec(ctx, q, h.ID, h.CompanyID, h.EmployeeID, h.PreviousAmount, h.NewAmount, h.Currency, h.Reason, h.EffectiveFrom, h.ApprovedBy, h.Notes, h.CreatedBy)
	return repoErr("AddHistory", err)
}

func (r *Repository) GetHistory(ctx context.Context, companyID, employeeID string) ([]CompensationHistory, error) {
	q := `SELECT id,company_id,employee_id,previous_amount,new_amount,currency,reason,effective_from,approved_by,notes,created_by,created_at FROM compensation_history WHERE company_id=$1 AND employee_id=$2 ORDER BY effective_from DESC`
	rows, err := r.pool.Query(ctx, q, companyID, employeeID)
	if err != nil {
		return nil, repoErr("GetHistory", err)
	}
	defer rows.Close()
	var res []CompensationHistory
	for rows.Next() {
		var h CompensationHistory
		if err := rows.Scan(&h.ID, &h.CompanyID, &h.EmployeeID, &h.PreviousAmount, &h.NewAmount, &h.Currency, &h.Reason, &h.EffectiveFrom, &h.ApprovedBy, &h.Notes, &h.CreatedBy, &h.CreatedAt); err != nil {
			return nil, repoErr("GetHistory", err)
		}
		res = append(res, h)
	}
	return res, nil
}

// ---------------------------------------------------------------------------
// Adjustments
// ---------------------------------------------------------------------------

func (r *Repository) CreateAdjustment(ctx context.Context, a *CompensationAdjustment) error {
	q := `INSERT INTO compensation_adjustments (id,company_id,employee_id,adjustment_type,value,currency,reason,effective_from,status,notes,created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.pool.Exec(ctx, q, a.ID, a.CompanyID, a.EmployeeID, a.AdjustmentType, a.Value, a.Currency, a.Reason, a.EffectiveFrom, a.Status, a.Notes, a.CreatedBy)
	return repoErr("CreateAdjustment", err)
}

func (r *Repository) GetAdjustment(ctx context.Context, companyID, id string) (*CompensationAdjustment, error) {
	q := `SELECT id,company_id,employee_id,adjustment_type,value,currency,reason,effective_from,status,approved_by,approved_at,applied_by,applied_at,notes,created_by,created_at,updated_at FROM compensation_adjustments WHERE id=$1 AND company_id=$2`
	a := &CompensationAdjustment{}
	err := r.pool.QueryRow(ctx, q, id, companyID).Scan(&a.ID, &a.CompanyID, &a.EmployeeID, &a.AdjustmentType, &a.Value, &a.Currency, &a.Reason, &a.EffectiveFrom, &a.Status, &a.ApprovedBy, &a.ApprovedAt, &a.AppliedBy, &a.AppliedAt, &a.Notes, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetAdjustment", err)
	}
	return a, nil
}

func (r *Repository) ListAdjustments(ctx context.Context, companyID string, filter AdjustmentFilter) ([]CompensationAdjustment, error) {
	q := `SELECT id,company_id,employee_id,adjustment_type,value,currency,reason,effective_from,status,approved_by,approved_at,applied_by,applied_at,notes,created_by,created_at,updated_at FROM compensation_adjustments WHERE company_id=$1`
	args := []any{companyID}
	idx := 2
	if filter.EmployeeID != nil {
		q += fmt.Sprintf(" AND employee_id=$%d", idx)
		args = append(args, *filter.EmployeeID)
		idx++
	}
	if filter.Status != nil {
		q += fmt.Sprintf(" AND status=$%d", idx)
		args = append(args, *filter.Status)
		idx++
	}
	q += " ORDER BY created_at DESC"
	if filter.Limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", idx)
		args = append(args, filter.Limit)
		idx++
	}
	if filter.Offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", idx)
		args = append(args, filter.Offset)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListAdjustments", err)
	}
	defer rows.Close()
	var res []CompensationAdjustment
	for rows.Next() {
		var a CompensationAdjustment
		if err := rows.Scan(&a.ID, &a.CompanyID, &a.EmployeeID, &a.AdjustmentType, &a.Value, &a.Currency, &a.Reason, &a.EffectiveFrom, &a.Status, &a.ApprovedBy, &a.ApprovedAt, &a.AppliedBy, &a.AppliedAt, &a.Notes, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, repoErr("ListAdjustments", err)
		}
		res = append(res, a)
	}
	return res, nil
}

func (r *Repository) UpdateAdjustmentStatus(ctx context.Context, id, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE compensation_adjustments SET status=$2,updated_at=NOW() WHERE id=$1`, id, status)
	return repoErr("UpdateAdjustmentStatus", err)
}

func (r *Repository) ApproveAdjustment(ctx context.Context, id, approvedBy string) error {
	_, err := r.pool.Exec(ctx, `UPDATE compensation_adjustments SET status='approved',approved_by=$2,approved_at=NOW(),updated_at=NOW() WHERE id=$1`, id, approvedBy)
	return repoErr("ApproveAdjustment", err)
}

func (r *Repository) ApplyAdjustment(ctx context.Context, id, appliedBy string) error {
	_, err := r.pool.Exec(ctx, `UPDATE compensation_adjustments SET status='applied',applied_by=$2,applied_at=NOW(),updated_at=NOW() WHERE id=$1`, id, appliedBy)
	return repoErr("ApplyAdjustment", err)
}

// ---------------------------------------------------------------------------
// Proposals
// ---------------------------------------------------------------------------

func (r *Repository) CreateProposal(ctx context.Context, p *SalaryAdjustmentProposal) error {
	q := `INSERT INTO salary_adjustment_proposals (id,company_id,review_id,employee_id,current_amount,proposed_amount,increase_percentage,reason,performance_score,market_position,manager_comment,status,created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`
	_, err := r.pool.Exec(ctx, q, p.ID, p.CompanyID, p.ReviewID, p.EmployeeID, p.CurrentAmount, p.ProposedAmount, p.IncreasePercentage, p.Reason, p.PerformanceScore, p.MarketPosition, p.ManagerComment, p.Status, p.CreatedBy)
	return repoErr("CreateProposal", err)
}

func (r *Repository) GetProposal(ctx context.Context, companyID, id string) (*SalaryAdjustmentProposal, error) {
	q := `SELECT id,company_id,review_id,employee_id,current_amount,proposed_amount,increase_percentage,reason,performance_score,market_position,manager_comment,hr_comment,status,submitted_by,approved_by,approved_at,rejected_by,rejected_at,rejection_reason,created_by,created_at,updated_at FROM salary_adjustment_proposals WHERE id=$1 AND company_id=$2`
	p := &SalaryAdjustmentProposal{}
	err := r.pool.QueryRow(ctx, q, id, companyID).Scan(&p.ID, &p.CompanyID, &p.ReviewID, &p.EmployeeID, &p.CurrentAmount, &p.ProposedAmount, &p.IncreasePercentage, &p.Reason, &p.PerformanceScore, &p.MarketPosition, &p.ManagerComment, &p.HRComment, &p.Status, &p.SubmittedBy, &p.ApprovedBy, &p.ApprovedAt, &p.RejectedBy, &p.RejectedAt, &p.RejectionReason, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetProposal", err)
	}
	return p, nil
}

func (r *Repository) ListProposals(ctx context.Context, companyID string, filter ProposalFilter) ([]SalaryAdjustmentProposal, error) {
	q := `SELECT id,company_id,review_id,employee_id,current_amount,proposed_amount,increase_percentage,reason,performance_score,market_position,manager_comment,hr_comment,status,submitted_by,approved_by,approved_at,rejected_by,rejected_at,rejection_reason,created_by,created_at,updated_at FROM salary_adjustment_proposals WHERE company_id=$1`
	args := []any{companyID}
	idx := 2
	if filter.ReviewID != nil {
		q += fmt.Sprintf(" AND review_id=$%d", idx)
		args = append(args, *filter.ReviewID)
		idx++
	}
	if filter.EmployeeID != nil {
		q += fmt.Sprintf(" AND employee_id=$%d", idx)
		args = append(args, *filter.EmployeeID)
		idx++
	}
	if filter.Status != nil {
		q += fmt.Sprintf(" AND status=$%d", idx)
		args = append(args, *filter.Status)
		idx++
	}
	q += " ORDER BY created_at DESC"
	if filter.Limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", idx)
		args = append(args, filter.Limit)
		idx++
	}
	if filter.Offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", idx)
		args = append(args, filter.Offset)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListProposals", err)
	}
	defer rows.Close()
	var res []SalaryAdjustmentProposal
	for rows.Next() {
		var p SalaryAdjustmentProposal
		if err := rows.Scan(&p.ID, &p.CompanyID, &p.ReviewID, &p.EmployeeID, &p.CurrentAmount, &p.ProposedAmount, &p.IncreasePercentage, &p.Reason, &p.PerformanceScore, &p.MarketPosition, &p.ManagerComment, &p.HRComment, &p.Status, &p.SubmittedBy, &p.ApprovedBy, &p.ApprovedAt, &p.RejectedBy, &p.RejectedAt, &p.RejectionReason, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, repoErr("ListProposals", err)
		}
		res = append(res, p)
	}
	return res, nil
}

func (r *Repository) UpdateProposalStatus(ctx context.Context, id, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE salary_adjustment_proposals SET status=$2,updated_at=NOW() WHERE id=$1`, id, status)
	return repoErr("UpdateProposalStatus", err)
}

func (r *Repository) ApproveProposal(ctx context.Context, id, approvedBy string) error {
	_, err := r.pool.Exec(ctx, `UPDATE salary_adjustment_proposals SET status='approved',approved_by=$2,approved_at=NOW(),updated_at=NOW() WHERE id=$1`, id, approvedBy)
	return repoErr("ApproveProposal", err)
}

func (r *Repository) RejectProposal(ctx context.Context, id, rejectedBy, reason string) error {
	_, err := r.pool.Exec(ctx, `UPDATE salary_adjustment_proposals SET status='rejected',rejected_by=$2,rejected_at=NOW(),rejection_reason=$3,updated_at=NOW() WHERE id=$1`, id, rejectedBy, reason)
	return repoErr("RejectProposal", err)
}

// ---------------------------------------------------------------------------
// Bonus Plans
// ---------------------------------------------------------------------------

func (r *Repository) CreateBonusPlan(ctx context.Context, bp *BonusPlan) error {
	q := `INSERT INTO bonus_plans (id,company_id,name,description,period,target_percentage,maximum_percentage,eligibility_rules,status,created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.pool.Exec(ctx, q, bp.ID, bp.CompanyID, bp.Name, bp.Description, bp.Period, bp.TargetPercentage, bp.MaximumPercentage, bp.EligibilityRules, bp.Status, bp.CreatedBy)
	return repoErr("CreateBonusPlan", err)
}

func (r *Repository) GetBonusPlan(ctx context.Context, companyID, id string) (*BonusPlan, error) {
	q := `SELECT id,company_id,name,description,period,target_percentage,maximum_percentage,eligibility_rules,status,created_by,created_at,updated_at FROM bonus_plans WHERE id=$1 AND company_id=$2`
	bp := &BonusPlan{}
	err := r.pool.QueryRow(ctx, q, id, companyID).Scan(&bp.ID, &bp.CompanyID, &bp.Name, &bp.Description, &bp.Period, &bp.TargetPercentage, &bp.MaximumPercentage, &bp.EligibilityRules, &bp.Status, &bp.CreatedBy, &bp.CreatedAt, &bp.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetBonusPlan", err)
	}
	return bp, nil
}

func (r *Repository) ListBonusPlans(ctx context.Context, companyID string) ([]BonusPlan, error) {
	q := `SELECT id,company_id,name,description,period,target_percentage,maximum_percentage,eligibility_rules,status,created_by,created_at,updated_at FROM bonus_plans WHERE company_id=$1 ORDER BY name`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListBonusPlans", err)
	}
	defer rows.Close()
	var res []BonusPlan
	for rows.Next() {
		var bp BonusPlan
		if err := rows.Scan(&bp.ID, &bp.CompanyID, &bp.Name, &bp.Description, &bp.Period, &bp.TargetPercentage, &bp.MaximumPercentage, &bp.EligibilityRules, &bp.Status, &bp.CreatedBy, &bp.CreatedAt, &bp.UpdatedAt); err != nil {
			return nil, repoErr("ListBonusPlans", err)
		}
		res = append(res, bp)
	}
	return res, nil
}

// ---------------------------------------------------------------------------
// Bonuses
// ---------------------------------------------------------------------------

func (r *Repository) CreateBonus(ctx context.Context, b *Bonus) error {
	q := `INSERT INTO bonuses (id,company_id,employee_id,bonus_plan_id,name,bonus_type,amount,currency,period,reason,status,created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := r.pool.Exec(ctx, q, b.ID, b.CompanyID, b.EmployeeID, b.BonusPlanID, b.Name, b.BonusType, b.Amount, b.Currency, b.Period, b.Reason, b.Status, b.CreatedBy)
	return repoErr("CreateBonus", err)
}

func (r *Repository) GetBonus(ctx context.Context, companyID, id string) (*Bonus, error) {
	q := `SELECT id,company_id,employee_id,bonus_plan_id,name,bonus_type,amount,currency,period,reason,status,approved_by,approved_at,paid_at,created_by,created_at,updated_at FROM bonuses WHERE id=$1 AND company_id=$2`
	b := &Bonus{}
	err := r.pool.QueryRow(ctx, q, id, companyID).Scan(&b.ID, &b.CompanyID, &b.EmployeeID, &b.BonusPlanID, &b.Name, &b.BonusType, &b.Amount, &b.Currency, &b.Period, &b.Reason, &b.Status, &b.ApprovedBy, &b.ApprovedAt, &b.PaidAt, &b.CreatedBy, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetBonus", err)
	}
	return b, nil
}

func (r *Repository) ListBonuses(ctx context.Context, companyID string, filter BonusFilter) ([]Bonus, error) {
	q := `SELECT id,company_id,employee_id,bonus_plan_id,name,bonus_type,amount,currency,period,reason,status,approved_by,approved_at,paid_at,created_by,created_at,updated_at FROM bonuses WHERE company_id=$1`
	args := []any{companyID}
	idx := 2
	if filter.EmployeeID != nil {
		q += fmt.Sprintf(" AND employee_id=$%d", idx)
		args = append(args, *filter.EmployeeID)
		idx++
	}
	if filter.Status != nil {
		q += fmt.Sprintf(" AND status=$%d", idx)
		args = append(args, *filter.Status)
		idx++
	}
	q += " ORDER BY created_at DESC"
	if filter.Limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", idx)
		args = append(args, filter.Limit)
		idx++
	}
	if filter.Offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", idx)
		args = append(args, filter.Offset)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListBonuses", err)
	}
	defer rows.Close()
	var res []Bonus
	for rows.Next() {
		var b Bonus
		if err := rows.Scan(&b.ID, &b.CompanyID, &b.EmployeeID, &b.BonusPlanID, &b.Name, &b.BonusType, &b.Amount, &b.Currency, &b.Period, &b.Reason, &b.Status, &b.ApprovedBy, &b.ApprovedAt, &b.PaidAt, &b.CreatedBy, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, repoErr("ListBonuses", err)
		}
		res = append(res, b)
	}
	return res, nil
}

func (r *Repository) UpdateBonusStatus(ctx context.Context, id, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE bonuses SET status=$2,updated_at=NOW() WHERE id=$1`, id, status)
	return repoErr("UpdateBonusStatus", err)
}

func (r *Repository) ApproveBonus(ctx context.Context, id, approvedBy string) error {
	_, err := r.pool.Exec(ctx, `UPDATE bonuses SET status='approved',approved_by=$2,approved_at=NOW(),updated_at=NOW() WHERE id=$1`, id, approvedBy)
	return repoErr("ApproveBonus", err)
}

// ---------------------------------------------------------------------------
// Benefits (catalog)
// ---------------------------------------------------------------------------

func (r *Repository) CreateBenefit(ctx context.Context, b *Benefit) error {
	q := `INSERT INTO benefits (id,company_id,code,name,description,benefit_type,provider,cost_amount,cost_currency,frequency,taxable,active,created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`
	_, err := r.pool.Exec(ctx, q, b.ID, b.CompanyID, b.Code, b.Name, b.Description, b.BenefitType, b.Provider, b.CostAmount, b.CostCurrency, b.Frequency, b.Taxable, b.Active, b.CreatedBy)
	return repoErr("CreateBenefit", err)
}

func (r *Repository) UpdateBenefit(ctx context.Context, b *Benefit) error {
	q := `UPDATE benefits SET name=COALESCE($3,name),description=COALESCE($4,description),benefit_type=COALESCE($5,benefit_type),provider=COALESCE($6,provider),cost_amount=COALESCE($7,cost_amount),cost_currency=COALESCE($8,cost_currency),frequency=COALESCE($9,frequency),taxable=COALESCE($10,taxable),active=COALESCE($11,active),updated_at=NOW() WHERE id=$1 AND company_id=$2`
	_, err := r.pool.Exec(ctx, q, b.ID, b.CompanyID, b.Name, b.Description, b.BenefitType, b.Provider, b.CostAmount, b.CostCurrency, b.Frequency, b.Taxable, b.Active)
	return repoErr("UpdateBenefit", err)
}

func (r *Repository) GetBenefit(ctx context.Context, companyID, id string) (*Benefit, error) {
	q := `SELECT id,company_id,code,name,description,benefit_type,provider,cost_amount,cost_currency,frequency,taxable,active,created_by,created_at,updated_at FROM benefits WHERE id=$1 AND company_id=$2`
	b := &Benefit{}
	err := r.pool.QueryRow(ctx, q, id, companyID).Scan(&b.ID, &b.CompanyID, &b.Code, &b.Name, &b.Description, &b.BenefitType, &b.Provider, &b.CostAmount, &b.CostCurrency, &b.Frequency, &b.Taxable, &b.Active, &b.CreatedBy, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetBenefit", err)
	}
	return b, nil
}

func (r *Repository) ListBenefits(ctx context.Context, companyID string, filter BenefitFilter) ([]Benefit, error) {
	q := `SELECT id,company_id,code,name,description,benefit_type,provider,cost_amount,cost_currency,frequency,taxable,active,created_by,created_at,updated_at FROM benefits WHERE company_id=$1`
	args := []any{companyID}
	idx := 2
	if filter.Active != nil {
		q += fmt.Sprintf(" AND active=$%d", idx)
		args = append(args, *filter.Active)
		idx++
	}
	if filter.BenefitType != nil {
		q += fmt.Sprintf(" AND benefit_type=$%d", idx)
		args = append(args, *filter.BenefitType)
		idx++
	}
	q += " ORDER BY code"
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListBenefits", err)
	}
	defer rows.Close()
	var res []Benefit
	for rows.Next() {
		var b Benefit
		if err := rows.Scan(&b.ID, &b.CompanyID, &b.Code, &b.Name, &b.Description, &b.BenefitType, &b.Provider, &b.CostAmount, &b.CostCurrency, &b.Frequency, &b.Taxable, &b.Active, &b.CreatedBy, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, repoErr("ListBenefits", err)
		}
		res = append(res, b)
	}
	return res, nil
}

// ---------------------------------------------------------------------------
// Employee Benefits
// ---------------------------------------------------------------------------

func (r *Repository) AssignBenefit(ctx context.Context, eb *EmployeeBenefit) error {
	q := `INSERT INTO employee_benefits (id,company_id,employee_id,benefit_id,enrollment_date,effective_from,effective_to,employee_cost,company_cost,currency,status,created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := r.pool.Exec(ctx, q, eb.ID, eb.CompanyID, eb.EmployeeID, eb.BenefitID, eb.EnrollmentDate, eb.EffectiveFrom, eb.EffectiveTo, eb.EmployeeCost, eb.CompanyCost, eb.Currency, eb.Status, eb.CreatedBy)
	return repoErr("AssignBenefit", err)
}

func (r *Repository) GetEmployeeBenefit(ctx context.Context, companyID, id string) (*EmployeeBenefit, error) {
	q := `SELECT id,company_id,employee_id,benefit_id,enrollment_date,effective_from,effective_to,employee_cost,company_cost,currency,status,created_by,created_at,updated_at FROM employee_benefits WHERE id=$1 AND company_id=$2`
	eb := &EmployeeBenefit{}
	err := r.pool.QueryRow(ctx, q, id, companyID).Scan(&eb.ID, &eb.CompanyID, &eb.EmployeeID, &eb.BenefitID, &eb.EnrollmentDate, &eb.EffectiveFrom, &eb.EffectiveTo, &eb.EmployeeCost, &eb.CompanyCost, &eb.Currency, &eb.Status, &eb.CreatedBy, &eb.CreatedAt, &eb.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetEmployeeBenefit", err)
	}
	return eb, nil
}

func (r *Repository) ListEmployeeBenefits(ctx context.Context, companyID, employeeID string) ([]EmployeeBenefit, error) {
	q := `SELECT id,company_id,employee_id,benefit_id,enrollment_date,effective_from,effective_to,employee_cost,company_cost,currency,status,created_by,created_at,updated_at FROM employee_benefits WHERE company_id=$1 AND employee_id=$2 ORDER BY effective_from DESC`
	rows, err := r.pool.Query(ctx, q, companyID, employeeID)
	if err != nil {
		return nil, repoErr("ListEmployeeBenefits", err)
	}
	defer rows.Close()
	var res []EmployeeBenefit
	for rows.Next() {
		var eb EmployeeBenefit
		if err := rows.Scan(&eb.ID, &eb.CompanyID, &eb.EmployeeID, &eb.BenefitID, &eb.EnrollmentDate, &eb.EffectiveFrom, &eb.EffectiveTo, &eb.EmployeeCost, &eb.CompanyCost, &eb.Currency, &eb.Status, &eb.CreatedBy, &eb.CreatedAt, &eb.UpdatedAt); err != nil {
			return nil, repoErr("ListEmployeeBenefits", err)
		}
		res = append(res, eb)
	}
	return res, nil
}

func (r *Repository) UpdateEmployeeBenefitStatus(ctx context.Context, id, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE employee_benefits SET status=$2,updated_at=NOW() WHERE id=$1`, id, status)
	return repoErr("UpdateEmployeeBenefitStatus", err)
}

// ---------------------------------------------------------------------------
// Reviews
// ---------------------------------------------------------------------------

func (r *Repository) CreateReview(ctx context.Context, rv *CompensationReview) error {
	q := `INSERT INTO compensation_reviews (id,company_id,name,description,period,start_date,end_date,budget,currency,status,created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.pool.Exec(ctx, q, rv.ID, rv.CompanyID, rv.Name, rv.Description, rv.Period, rv.StartDate, rv.EndDate, rv.Budget, rv.Currency, rv.Status, rv.CreatedBy)
	return repoErr("CreateReview", err)
}

func (r *Repository) GetReview(ctx context.Context, companyID, id string) (*CompensationReview, error) {
	q := `SELECT id,company_id,name,description,period,start_date,end_date,budget,currency,status,created_by,created_at,updated_at FROM compensation_reviews WHERE id=$1 AND company_id=$2`
	rv := &CompensationReview{}
	err := r.pool.QueryRow(ctx, q, id, companyID).Scan(&rv.ID, &rv.CompanyID, &rv.Name, &rv.Description, &rv.Period, &rv.StartDate, &rv.EndDate, &rv.Budget, &rv.Currency, &rv.Status, &rv.CreatedBy, &rv.CreatedAt, &rv.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetReview", err)
	}
	return rv, nil
}

func (r *Repository) ListReviews(ctx context.Context, companyID string) ([]CompensationReview, error) {
	q := `SELECT id,company_id,name,description,period,start_date,end_date,budget,currency,status,created_by,created_at,updated_at FROM compensation_reviews WHERE company_id=$1 ORDER BY start_date DESC`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListReviews", err)
	}
	defer rows.Close()
	var res []CompensationReview
	for rows.Next() {
		var rv CompensationReview
		if err := rows.Scan(&rv.ID, &rv.CompanyID, &rv.Name, &rv.Description, &rv.Period, &rv.StartDate, &rv.EndDate, &rv.Budget, &rv.Currency, &rv.Status, &rv.CreatedBy, &rv.CreatedAt, &rv.UpdatedAt); err != nil {
			return nil, repoErr("ListReviews", err)
		}
		res = append(res, rv)
	}
	return res, nil
}

func (r *Repository) UpdateReviewStatus(ctx context.Context, id, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE compensation_reviews SET status=$2,updated_at=NOW() WHERE id=$1`, id, status)
	return repoErr("UpdateReviewStatus", err)
}

// ---------------------------------------------------------------------------
// Budgets
// ---------------------------------------------------------------------------

func (r *Repository) CreateBudget(ctx context.Context, b *CompensationBudget) error {
	q := `INSERT INTO compensation_budgets (id,company_id,year,department_id,budget_amount,committed_amount,spent_amount,currency,status,created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.pool.Exec(ctx, q, b.ID, b.CompanyID, b.Year, b.DepartmentID, b.BudgetAmount, b.CommittedAmount, b.SpentAmount, b.Currency, b.Status, b.CreatedBy)
	return repoErr("CreateBudget", err)
}

func (r *Repository) GetBudget(ctx context.Context, companyID, id string) (*CompensationBudget, error) {
	q := `SELECT id,company_id,year,department_id,budget_amount,committed_amount,spent_amount,currency,status,created_by,created_at,updated_at FROM compensation_budgets WHERE id=$1 AND company_id=$2`
	b := &CompensationBudget{}
	err := r.pool.QueryRow(ctx, q, id, companyID).Scan(&b.ID, &b.CompanyID, &b.Year, &b.DepartmentID, &b.BudgetAmount, &b.CommittedAmount, &b.SpentAmount, &b.Currency, &b.Status, &b.CreatedBy, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetBudget", err)
	}
	return b, nil
}

func (r *Repository) ListBudgets(ctx context.Context, companyID string) ([]CompensationBudget, error) {
	q := `SELECT id,company_id,year,department_id,budget_amount,committed_amount,spent_amount,currency,status,created_by,created_at,updated_at FROM compensation_budgets WHERE company_id=$1 ORDER BY year DESC`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListBudgets", err)
	}
	defer rows.Close()
	var res []CompensationBudget
	for rows.Next() {
		var b CompensationBudget
		if err := rows.Scan(&b.ID, &b.CompanyID, &b.Year, &b.DepartmentID, &b.BudgetAmount, &b.CommittedAmount, &b.SpentAmount, &b.Currency, &b.Status, &b.CreatedBy, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, repoErr("ListBudgets", err)
		}
		res = append(res, b)
	}
	return res, nil
}

func (r *Repository) UpdateBudgetAmounts(ctx context.Context, companyID, budgetID string, committedAmount, spentAmount decimal.Decimal) error {
	_, err := r.pool.Exec(ctx, `UPDATE compensation_budgets SET committed_amount=$3,spent_amount=$4,updated_at=NOW() WHERE id=$1 AND company_id=$2`, budgetID, companyID, committedAmount, spentAmount)
	return repoErr("UpdateBudgetAmounts", err)
}

// ---------------------------------------------------------------------------
// Equity Snapshots
// ---------------------------------------------------------------------------

func (r *Repository) CreateEquitySnapshot(ctx context.Context, es *CompensationEquitySnapshot) error {
	q := `INSERT INTO compensation_equity_snapshots (id,company_id,snapshot_date,department_id,position_id,grade_id,employee_count,median_compensation,average_compensation,min_compensation,max_compensation,currency,metadata,created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`
	_, err := r.pool.Exec(ctx, q, es.ID, es.CompanyID, es.SnapshotDate, es.DepartmentID, es.PositionID, es.GradeID, es.EmployeeCount, es.MedianCompensation, es.AverageCompensation, es.MinCompensation, es.MaxCompensation, es.Currency, es.Metadata, es.CreatedBy)
	return repoErr("CreateEquitySnapshot", err)
}

// ---------------------------------------------------------------------------
// Audit
// ---------------------------------------------------------------------------

func (r *Repository) LogAudit(ctx context.Context, log *CompensationAuditLog) error {
	q := `INSERT INTO compensation_audit_logs (id,company_id,user_id,action,entity_type,entity_id,old_value,new_value,ip_address,user_agent) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.pool.Exec(ctx, q, log.ID, log.CompanyID, log.UserID, log.Action, log.EntityType, log.EntityID, log.OldValue, log.NewValue, log.IPAddress, log.UserAgent)
	return repoErr("LogAudit", err)
}

// ---------------------------------------------------------------------------
// Domain Events
// ---------------------------------------------------------------------------

func (r *Repository) CreateDomainEvent(ctx context.Context, e *CompensationDomainEvent) error {
	q := `INSERT INTO compensation_domain_events (id,company_id,event_type,entity_type,entity_id,payload,created_by) VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.pool.Exec(ctx, q, e.ID, e.CompanyID, e.EventType, e.EntityType, e.EntityID, e.Payload, e.CreatedBy)
	return repoErr("CreateDomainEvent", err)
}

func (r *Repository) ListPendingDomainEvents(ctx context.Context, companyID string) ([]CompensationDomainEvent, error) {
	q := `SELECT id,company_id,event_type,entity_type,entity_id,payload,created_by,processed_at,created_at FROM compensation_domain_events WHERE company_id=$1 AND processed_at IS NULL ORDER BY created_at ASC LIMIT 100`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListPendingDomainEvents", err)
	}
	defer rows.Close()
	var res []CompensationDomainEvent
	for rows.Next() {
		var e CompensationDomainEvent
		if err := rows.Scan(&e.ID, &e.CompanyID, &e.EventType, &e.EntityType, &e.EntityID, &e.Payload, &e.CreatedBy, &e.ProcessedAt, &e.CreatedAt); err != nil {
			return nil, repoErr("ListPendingDomainEvents", err)
		}
		res = append(res, e)
	}
	return res, nil
}

func (r *Repository) MarkDomainEventProcessed(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `UPDATE compensation_domain_events SET processed_at=NOW() WHERE id=$1`, id)
	return repoErr("MarkDomainEventProcessed", err)
}

// ---------------------------------------------------------------------------
// Dashboard stats
// ---------------------------------------------------------------------------

func (r *Repository) GetDashboardStats(ctx context.Context, companyID string) (*DashboardStats, error) {
	ds := &DashboardStats{Currency: "USD"}
	r.pool.QueryRow(ctx, `SELECT COALESCE(SUM(base_amount),0) FROM employee_compensations WHERE company_id=$1 AND status='active'`, companyID).Scan(&ds.TotalSalaryCost)
	r.pool.QueryRow(ctx, `SELECT COALESCE(AVG(base_amount),0) FROM employee_compensations WHERE company_id=$1 AND status='active'`, companyID).Scan(&ds.AverageCompensation)
	r.pool.QueryRow(ctx, `SELECT COALESCE(SUM(company_cost),0) FROM employee_benefits eb JOIN benefits b ON b.id=eb.benefit_id WHERE eb.company_id=$1 AND eb.status='active'`, companyID).Scan(&ds.BenefitCost)
	r.pool.QueryRow(ctx, `SELECT COALESCE(SUM(amount),0) FROM bonuses WHERE company_id=$1 AND status='paid'`, companyID).Scan(&ds.TotalBonuses)
	r.pool.QueryRow(ctx, `SELECT COALESCE(SUM(budget_amount),0) FROM compensation_budgets WHERE company_id=$1 AND status='active'`, companyID).Scan(&ds.BudgetTotal)
	r.pool.QueryRow(ctx, `SELECT COALESCE(SUM(committed_amount+spent_amount),0) FROM compensation_budgets WHERE company_id=$1 AND status='active'`, companyID).Scan(&ds.BudgetUsed)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM salary_adjustment_proposals WHERE company_id=$1 AND status='submitted'`, companyID).Scan(&ds.PendingProposals)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM employee_compensations ec WHERE ec.company_id=$1 AND ec.status='active' AND ec.salary_band_id IS NOT NULL AND (ec.base_amount<(SELECT minimum_amount FROM salary_bands sb WHERE sb.id=ec.salary_band_id) OR ec.base_amount>(SELECT maximum_amount FROM salary_bands sb WHERE sb.id=ec.salary_band_id))`, companyID).Scan(&ds.EmployeesOutOfBand)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM compensation_reviews WHERE company_id=$1 AND status IN ('open','in_review')`, companyID).Scan(&ds.ActiveReviews)
	return ds, nil
}

// ---------------------------------------------------------------------------
// Band analysis
// ---------------------------------------------------------------------------

func (r *Repository) GetBandAnalysis(ctx context.Context, companyID, bandID string) (*BandAnalysis, error) {
	band, err := r.GetBand(ctx, companyID, bandID)
	if err != nil {
		return nil, err
	}
	ba := &BandAnalysis{Band: *band}
	q := `SELECT COUNT(*),COALESCE(AVG(base_amount),0),COALESCE(MIN(base_amount),0),COALESCE(MAX(base_amount),0) FROM employee_compensations WHERE company_id=$1 AND salary_band_id=$2 AND status='active'`
	r.pool.QueryRow(ctx, q, companyID, bandID).Scan(&ba.EmployeeCount, &ba.AverageSalary, &ba.MinSalary, &ba.MaxSalary)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM employee_compensations WHERE company_id=$1 AND salary_band_id=$2 AND status='active' AND base_amount<$3`, companyID, bandID, band.MinimumAmount).Scan(&ba.BelowRange)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM employee_compensations WHERE company_id=$1 AND salary_band_id=$2 AND status='active' AND base_amount>=$3 AND base_amount<=$4`, companyID, bandID, band.MinimumAmount, band.MaximumAmount).Scan(&ba.InRange)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM employee_compensations WHERE company_id=$1 AND salary_band_id=$2 AND status='active' AND base_amount>$3`, companyID, bandID, band.MaximumAmount).Scan(&ba.AboveRange)
	return ba, nil
}

// ---------------------------------------------------------------------------
// Worker support methods
// ---------------------------------------------------------------------------

func (r *Repository) GetExpiringBenefits(ctx context.Context) ([]EmployeeBenefit, error) {
	q := `SELECT eb.id, eb.employee_id, eb.company_id, eb.benefit_id, eb.effective_date, eb.expiration_date, eb.status, eb.enrollment_data, eb.created_at, eb.updated_at
		FROM employee_benefits eb
		WHERE eb.status='active' AND eb.expiration_date IS NOT NULL AND eb.expiration_date <= NOW() + INTERVAL '14 days' AND eb.expiration_date > NOW()`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []EmployeeBenefit
	for rows.Next() {
		var eb EmployeeBenefit
		if err := rows.Scan(&eb.ID, &eb.EmployeeID, &eb.CompanyID, &eb.BenefitID, &eb.EffectiveDate, &eb.ExpirationDate,
			&eb.Status, &eb.EnrollmentData, &eb.CreatedAt, &eb.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, eb)
	}
	return result, nil
}

func (r *Repository) GetBudgetsNearThreshold(ctx context.Context) ([]CompensationBudget, error) {
	q := `SELECT id, company_id, fiscal_year, period, base_amount, committed_amount, spent_amount, status, created_by, created_at, updated_at
		FROM compensation_budgets
		WHERE COALESCE(committed_amount,0)+COALESCE(spent_amount,0) > base_amount * 0.85 AND status='active'`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []CompensationBudget
	for rows.Next() {
		var b CompensationBudget
		if err := rows.Scan(&b.ID, &b.CompanyID, &b.FiscalYear, &b.Period, &b.BaseAmount, &b.CommittedAmount, &b.SpentAmount,
			&b.Status, &b.CreatedBy, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, b)
	}
	return result, nil
}

func (r *Repository) GetUnprocessedDomainEvents(ctx context.Context, source string) ([]CompensationDomainEvent, error) {
	q := `SELECT id, company_id, event_type, entity_type, entity_id, payload, created_by, processed_at, created_at
		FROM compensation_domain_events WHERE processed_at IS NULL ORDER BY created_at ASC LIMIT 100`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []CompensationDomainEvent
	for rows.Next() {
		var e CompensationDomainEvent
		if err := rows.Scan(&e.ID, &e.CompanyID, &e.EventType, &e.EntityType, &e.EntityID, &e.Payload, &e.CreatedBy,
			&e.ProcessedAt, &e.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, nil
}

func (r *Repository) NotifyBenefitExpiration(ctx context.Context, eb EmployeeBenefit) error {
	r.log.Warn("benefit expiring soon", zap.String("employee_benefit_id", eb.ID), zap.String("employee_id", eb.EmployeeID))
	return nil
}

func (r *Repository) NotifyBudgetAlert(ctx context.Context, b CompensationBudget) error {
	r.log.Warn("budget near threshold", zap.String("budget_id", b.ID), zap.String("company_id", b.CompanyID), zap.Int("fiscal_year", b.FiscalYear))
	return nil
}
