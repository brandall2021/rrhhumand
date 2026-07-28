package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/expenses/domain"
)

type TravelRepo struct {
	pool *pgxpool.Pool
}

func NewTravelRepo(pool *pgxpool.Pool) *TravelRepo {
	return &TravelRepo{pool: pool}
}

func (r *TravelRepo) CreateTravel(ctx context.Context, t *domain.Travel) error {
	q := `INSERT INTO travels (id,company_id,employee_id,purpose,origin,destination,departure_at,return_at,
		expected_cost,currency,status,notes,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`
	_, err := r.pool.Exec(ctx, q, t.ID, t.CompanyID, t.EmployeeID, t.Purpose, t.Origin, t.Destination,
		t.DepartureAt, t.ReturnAt, t.ExpectedCost, t.Currency, t.Status, t.Notes, t.CreatedBy)
	return repoErr("CreateTravel", err)
}

func (r *TravelRepo) GetTravel(ctx context.Context, companyID, id uuid.UUID) (*domain.Travel, error) {
	q := `SELECT id,company_id,employee_id,purpose,origin,destination,departure_at,return_at,
		expected_cost,currency,status,notes,created_by,created_at,updated_at
		FROM travels WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var t domain.Travel
	err := row.Scan(&t.ID, &t.CompanyID, &t.EmployeeID, &t.Purpose, &t.Origin, &t.Destination,
		&t.DepartureAt, &t.ReturnAt, &t.ExpectedCost, &t.Currency, &t.Status, &t.Notes, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetTravel", err)
	}
	return &t, nil
}

func (r *TravelRepo) ListTravels(ctx context.Context, companyID uuid.UUID, employeeID *uuid.UUID, status *string, from, to *time.Time) ([]domain.Travel, error) {
	q := `SELECT id,company_id,employee_id,purpose,origin,destination,departure_at,return_at,
		expected_cost,currency,status,notes,created_by,created_at,updated_at
		FROM travels WHERE company_id=$1`
	args := []any{companyID}
	n := 2
	if employeeID != nil {
		q += fmt.Sprintf(" AND employee_id=$%d", n)
		args = append(args, *employeeID)
		n++
	}
	if status != nil {
		q += fmt.Sprintf(" AND status=$%d", n)
		args = append(args, *status)
		n++
	}
	if from != nil {
		q += fmt.Sprintf(" AND departure_at>=$%d", n)
		args = append(args, *from)
		n++
	}
	if to != nil {
		q += fmt.Sprintf(" AND departure_at<=$%d", n)
		args = append(args, *to)
		n++
	}
	q += " ORDER BY departure_at DESC"
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListTravels", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Travel, error) {
		var t domain.Travel
		err := row.Scan(&t.ID, &t.CompanyID, &t.EmployeeID, &t.Purpose, &t.Origin, &t.Destination,
			&t.DepartureAt, &t.ReturnAt, &t.ExpectedCost, &t.Currency, &t.Status, &t.Notes, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt)
		return t, err
	})
}

func (r *TravelRepo) UpdateTravel(ctx context.Context, t *domain.Travel) error {
	q := `UPDATE travels SET purpose=$1,origin=$2,destination=$3,departure_at=$4,return_at=$5,
		expected_cost=$6,currency=$7,status=$8,notes=$9,updated_at=NOW() WHERE id=$10 AND company_id=$11`
	_, err := r.pool.Exec(ctx, q, t.Purpose, t.Origin, t.Destination, t.DepartureAt, t.ReturnAt,
		t.ExpectedCost, t.Currency, t.Status, t.Notes, t.ID, t.CompanyID)
	return repoErr("UpdateTravel", err)
}

func (r *TravelRepo) UpdateTravelStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE travels SET status=$1,updated_at=NOW() WHERE id=$2`, status, id)
	return repoErr("UpdateTravelStatus", err)
}

func (r *TravelRepo) AddParticipant(ctx context.Context, p *domain.TravelParticipant) error {
	q := `INSERT INTO travel_participants (id,travel_id,employee_id,role,notes) VALUES ($1,$2,$3,$4,$5)`
	_, err := r.pool.Exec(ctx, q, p.ID, p.TravelID, p.EmployeeID, p.Role, p.Notes)
	return repoErr("AddParticipant", err)
}

func (r *TravelRepo) RemoveParticipant(ctx context.Context, travelID, employeeID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM travel_participants WHERE travel_id=$1 AND employee_id=$2`, travelID, employeeID)
	return repoErr("RemoveParticipant", err)
}

func (r *TravelRepo) ListParticipants(ctx context.Context, travelID uuid.UUID) ([]domain.TravelParticipant, error) {
	q := `SELECT id,travel_id,employee_id,role,notes,created_at FROM travel_participants WHERE travel_id=$1 ORDER BY created_at`
	rows, err := r.pool.Query(ctx, q, travelID)
	if err != nil {
		return nil, repoErr("ListParticipants", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.TravelParticipant, error) {
		var p domain.TravelParticipant
		err := row.Scan(&p.ID, &p.TravelID, &p.EmployeeID, &p.Role, &p.Notes, &p.CreatedAt)
		return p, err
	})
}
