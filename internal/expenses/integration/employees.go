package integration

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EmployeeData struct {
	ID              uuid.UUID  `json:"id"`
	CompanyID       uuid.UUID  `json:"company_id"`
	UserID          uuid.UUID  `json:"user_id"`
	DepartmentID    *uuid.UUID `json:"department_id,omitempty"`
	PositionID      *uuid.UUID `json:"position_id,omitempty"`
	EmploymentType  string     `json:"employment_type"`
	CostCenter      *string    `json:"cost_center,omitempty"`
	ManagerID       *uuid.UUID `json:"manager_id,omitempty"`
	AdmissionDate   time.Time  `json:"admission_date"`
	Status          string     `json:"status"`
}

type EmployeeAdapter struct {
	pool *pgxpool.Pool
}

func NewEmployeeAdapter(pool *pgxpool.Pool) *EmployeeAdapter {
	return &EmployeeAdapter{pool: pool}
}

func (a *EmployeeAdapter) GetEmployee(ctx context.Context, id uuid.UUID) (*EmployeeData, error) {
	var emp EmployeeData
	err := a.pool.QueryRow(ctx, `
		SELECT e.id,e.company_id,e.user_id,e.department_id,e.position_id,e.employment_type,
			   e.cost_center,e.manager_id,e.admission_date,e.status
		FROM employees e WHERE e.id=$1`, id).
		Scan(&emp.ID, &emp.CompanyID, &emp.UserID, &emp.DepartmentID, &emp.PositionID,
			&emp.EmploymentType, &emp.CostCenter, &emp.ManagerID, &emp.AdmissionDate, &emp.Status)
	if err != nil {
		return nil, integErr("GetEmployee", err)
	}
	return &emp, nil
}

func (a *EmployeeAdapter) GetManager(ctx context.Context, employeeID uuid.UUID) (*uuid.UUID, error) {
	var managerID uuid.UUID
	err := a.pool.QueryRow(ctx, `SELECT manager_id FROM employees WHERE id=$1`, employeeID).Scan(&managerID)
	if err != nil {
		return nil, integErr("GetManager", err)
	}
	return &managerID, nil
}
