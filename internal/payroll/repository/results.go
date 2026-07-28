package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/rrhhumand/api/internal/payroll/domain"
)

func (r *Repository) BulkCreateItems(ctx context.Context, items []domain.PayrollItem) error {
	if len(items) == 0 {
		return nil
	}
	q := `INSERT INTO payroll_items (id,run_employee_id,concept_id,quantity,unit_value,base_amount,rate,amount,
		is_remunerative,is_deduction,is_employer_contribution,calculation_detail,sort_order) VALUES `
	args := []any{}
	n := 1
	for _, it := range items {
		detail, _ := json.Marshal(it.CalculationDetail)
		q += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),", n, n+1, n+2, n+3, n+4, n+5, n+6, n+7, n+8, n+9, n+10, n+11, n+12)
		args = append(args, it.ID, it.RunEmployeeID, it.ConceptID, it.Quantity, it.UnitValue, it.BaseAmount, it.Rate, it.Amount,
			it.IsRemunerative, it.IsDeduction, it.IsEmployerContribution, detail, it.SortOrder)
		n += 13
	}
	q = q[:len(q)-1]
	_, err := r.pool.Exec(ctx, q, args...)
	return repoErr("BulkCreateItems", err)
}

func (r *Repository) ListItems(ctx context.Context, runEmployeeID string) ([]domain.PayrollItem, error) {
	q := `SELECT id,run_employee_id,concept_id,quantity,unit_value,base_amount,rate,amount,
		is_remunerative,is_deduction,is_employer_contribution,calculation_detail,sort_order,created_at
		FROM payroll_items WHERE run_employee_id=$1 ORDER BY sort_order`
	rows, err := r.pool.Query(ctx, q, runEmployeeID)
	if err != nil {
		return nil, repoErr("ListItems", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.PayrollItem, error) {
		var it domain.PayrollItem
		var detail []byte
		err := row.Scan(&it.ID, &it.RunEmployeeID, &it.ConceptID, &it.Quantity, &it.UnitValue, &it.BaseAmount, &it.Rate,
			&it.Amount, &it.IsRemunerative, &it.IsDeduction, &it.IsEmployerContribution, &detail, &it.SortOrder, &it.CreatedAt)
		if err == nil && detail != nil {
			json.Unmarshal(detail, &it.CalculationDetail)
		}
		return it, err
	})
}

func (r *Repository) DeleteItemsForRunEmployee(ctx context.Context, runEmployeeID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM payroll_items WHERE run_employee_id=$1`, runEmployeeID)
	return repoErr("DeleteItemsForRunEmployee", err)
}

func (r *Repository) BulkCreateBases(ctx context.Context, bases []domain.PayrollBase) error {
	if len(bases) == 0 {
		return nil
	}
	q := `INSERT INTO payroll_bases (id,run_employee_id,base_type,base_amount,concept_ids,calculation_detail) VALUES `
	args := []any{}
	n := 1
	for _, b := range bases {
		detail, _ := json.Marshal(b.CalculationDetail)
		q += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d),", n, n+1, n+2, n+3, n+4, n+5)
		args = append(args, b.ID, b.RunEmployeeID, b.BaseType, b.BaseAmount, b.ConceptIDs, detail)
		n += 6
	}
	q = q[:len(q)-1]
	_, err := r.pool.Exec(ctx, q, args...)
	return repoErr("BulkCreateBases", err)
}

func (r *Repository) ListBases(ctx context.Context, runEmployeeID string) ([]domain.PayrollBase, error) {
	q := `SELECT id,run_employee_id,base_type,base_amount,concept_ids,calculation_detail,created_at
		FROM payroll_bases WHERE run_employee_id=$1 ORDER BY base_type`
	rows, err := r.pool.Query(ctx, q, runEmployeeID)
	if err != nil {
		return nil, repoErr("ListBases", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.PayrollBase, error) {
		var b domain.PayrollBase
		var detail []byte
		err := row.Scan(&b.ID, &b.RunEmployeeID, &b.BaseType, &b.BaseAmount, &b.ConceptIDs, &detail, &b.CreatedAt)
		if err == nil && detail != nil {
			json.Unmarshal(detail, &b.CalculationDetail)
		}
		return b, err
	})
}

func (r *Repository) BulkCreateDeductions(ctx context.Context, deductions []domain.PayrollDeduction) error {
	if len(deductions) == 0 {
		return nil
	}
	q := `INSERT INTO payroll_deductions (id,run_employee_id,concept_id,base_amount,rate,amount,cap_amount,exceeded_amount) VALUES `
	args := []any{}
	n := 1
	for _, d := range deductions {
		q += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),", n, n+1, n+2, n+3, n+4, n+5, n+6, n+7)
		args = append(args, d.ID, d.RunEmployeeID, d.ConceptID, d.BaseAmount, d.Rate, d.Amount, d.CapAmount, d.ExceededAmount)
		n += 8
	}
	q = q[:len(q)-1]
	_, err := r.pool.Exec(ctx, q, args...)
	return repoErr("BulkCreateDeductions", err)
}

func (r *Repository) BulkCreateContributions(ctx context.Context, contributions []domain.PayrollContribution) error {
	if len(contributions) == 0 {
		return nil
	}
	q := `INSERT INTO payroll_contributions (id,run_employee_id,concept_id,base_amount,rate,amount,cap_amount,exceeded_amount) VALUES `
	args := []any{}
	n := 1
	for _, c := range contributions {
		q += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),", n, n+1, n+2, n+3, n+4, n+5, n+6, n+7)
		args = append(args, c.ID, c.RunEmployeeID, c.ConceptID, c.BaseAmount, c.Rate, c.Amount, c.CapAmount, c.ExceededAmount)
		n += 8
	}
	q = q[:len(q)-1]
	_, err := r.pool.Exec(ctx, q, args...)
	return repoErr("BulkCreateContributions", err)
}
