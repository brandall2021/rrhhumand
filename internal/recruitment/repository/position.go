package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/recruitment/domain"
)

type PositionRepo struct {
	pool *pgxpool.Pool
}

func NewPositionRepo(pool *pgxpool.Pool) *PositionRepo {
	return &PositionRepo{pool: pool}
}

func (r *PositionRepo) Create(ctx context.Context, companyID string, req *domain.Position) (*domain.Position, error) {
	p := &domain.Position{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO job_positions_ats (company_id, requisition_id, title, department_id, location_id, employment_type, work_mode, description, requirements, responsibilities, benefits, salary_min, salary_max, currency, vacancies)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
		 RETURNING id, company_id, requisition_id, title, department_id, location_id, employment_type, work_mode, description, requirements, responsibilities, benefits, salary_min, salary_max, currency, vacancies, vacancies_filled, status, created_at, updated_at`,
		companyID, req.RequisitionID, req.Title, req.DepartmentID, req.LocationID,
		req.EmploymentType, req.WorkMode, req.Description, req.Requirements, req.Responsibilities,
		req.Benefits, req.SalaryMin, req.SalaryMax, req.Currency, req.Vacancies,
	).Scan(&p.ID, &p.CompanyID, &p.RequisitionID, &p.Title, &p.DepartmentID, &p.LocationID,
		&p.EmploymentType, &p.WorkMode, &p.Description, &p.Requirements, &p.Responsibilities,
		&p.Benefits, &p.SalaryMin, &p.SalaryMax, &p.Currency, &p.Vacancies, &p.VacanciesFilled,
		&p.Status, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *PositionRepo) GetByID(ctx context.Context, companyID, id string) (*domain.Position, error) {
	p := &domain.Position{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, requisition_id, title, department_id, location_id, employment_type, work_mode, description, requirements, responsibilities, benefits, salary_min, salary_max, currency, vacancies, vacancies_filled, status, created_at, updated_at
		 FROM job_positions_ats WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&p.ID, &p.CompanyID, &p.RequisitionID, &p.Title, &p.DepartmentID, &p.LocationID,
		&p.EmploymentType, &p.WorkMode, &p.Description, &p.Requirements, &p.Responsibilities,
		&p.Benefits, &p.SalaryMin, &p.SalaryMax, &p.Currency, &p.Vacancies, &p.VacanciesFilled,
		&p.Status, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *PositionRepo) List(ctx context.Context, companyID string, status string) ([]domain.Position, error) {
	query := `SELECT id, company_id, requisition_id, title, department_id, location_id, employment_type, work_mode, description, requirements, responsibilities, benefits, salary_min, salary_max, currency, vacancies, vacancies_filled, status, created_at, updated_at
		 FROM job_positions_ats WHERE company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if status != "" {
		query += fmt.Sprintf(" AND status=$%d", argIdx)
		args = append(args, status)
		argIdx++
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []domain.Position
	for rows.Next() {
		var p domain.Position
		rows.Scan(&p.ID, &p.CompanyID, &p.RequisitionID, &p.Title, &p.DepartmentID, &p.LocationID,
			&p.EmploymentType, &p.WorkMode, &p.Description, &p.Requirements, &p.Responsibilities,
			&p.Benefits, &p.SalaryMin, &p.SalaryMax, &p.Currency, &p.Vacancies, &p.VacanciesFilled,
			&p.Status, &p.CreatedAt, &p.UpdatedAt)
		positions = append(positions, p)
	}
	return positions, nil
}

func (r *PositionRepo) Update(ctx context.Context, companyID, id string, req *domain.Position) (*domain.Position, error) {
	p := &domain.Position{}
	err := r.pool.QueryRow(ctx,
		`UPDATE job_positions_ats SET
		 title=COALESCE($3,title), description=COALESCE($4,description), requirements=COALESCE($5,requirements),
		 responsibilities=COALESCE($6,responsibilities), benefits=COALESCE($7,benefits),
		 employment_type=COALESCE($8,employment_type), work_mode=COALESCE($9,work_mode),
		 salary_min=COALESCE($10,salary_min), salary_max=COALESCE($11,salary_max), currency=COALESCE($12,currency),
		 vacancies=COALESCE($13,vacancies), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, requisition_id, title, department_id, location_id, employment_type, work_mode, description, requirements, responsibilities, benefits, salary_min, salary_max, currency, vacancies, vacancies_filled, status, created_at, updated_at`,
		companyID, id, req.Title, req.Description, req.Requirements, req.Responsibilities,
		req.Benefits, req.EmploymentType, req.WorkMode, req.SalaryMin, req.SalaryMax, req.Currency, req.Vacancies,
	).Scan(&p.ID, &p.CompanyID, &p.RequisitionID, &p.Title, &p.DepartmentID, &p.LocationID,
		&p.EmploymentType, &p.WorkMode, &p.Description, &p.Requirements, &p.Responsibilities,
		&p.Benefits, &p.SalaryMin, &p.SalaryMax, &p.Currency, &p.Vacancies, &p.VacanciesFilled,
		&p.Status, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *PositionRepo) UpdateStatus(ctx context.Context, companyID, id, status string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE job_positions_ats SET status=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, status)
	return err
}

func (r *PositionRepo) UpdateFilled(ctx context.Context, companyID, id string, filled int) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE job_positions_ats SET vacancies_filled=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, filled)
	return err
}

func (r *PositionRepo) AddSkill(ctx context.Context, req *domain.PositionSkill) (*domain.PositionSkill, error) {
	s := &domain.PositionSkill{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO job_position_skills (position_id, skill, category, required, min_years, weight)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id, position_id, skill, category, required, min_years, weight`,
		req.PositionID, req.Skill, req.Category, req.Required, req.MinYears, req.Weight,
	).Scan(&s.ID, &s.PositionID, &s.Skill, &s.Category, &s.Required, &s.MinYears, &s.Weight)
	return s, err
}

func (r *PositionRepo) RemoveSkill(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM job_position_skills WHERE id=$1`, id)
	return err
}

func (r *PositionRepo) ListSkills(ctx context.Context, positionID string) ([]domain.PositionSkill, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, position_id, skill, category, required, min_years, weight
		 FROM job_position_skills WHERE position_id=$1`, positionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []domain.PositionSkill
	for rows.Next() {
		var s domain.PositionSkill
		rows.Scan(&s.ID, &s.PositionID, &s.Skill, &s.Category, &s.Required, &s.MinYears, &s.Weight)
		skills = append(skills, s)
	}
	return skills, nil
}
