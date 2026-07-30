package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/expenses/domain"
)

type PolicyRepo struct {
	pool *pgxpool.Pool
}

func NewPolicyRepo(pool *pgxpool.Pool) *PolicyRepo {
	return &PolicyRepo{pool: pool}
}

func (r *PolicyRepo) CreatePolicy(ctx context.Context, p *domain.ExpensePolicy) error {
	q := `INSERT INTO expense_policies (id,company_id,name,description,is_active,created_by)
		VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.pool.Exec(ctx, q, p.ID, p.CompanyID, p.Name, p.Description, p.IsActive, p.CreatedBy)
	return repoErr("CreatePolicy", err)
}

func (r *PolicyRepo) GetPolicy(ctx context.Context, companyID, id uuid.UUID) (*domain.ExpensePolicy, error) {
	q := `SELECT id,company_id,name,description,is_active,created_by,created_at,updated_at
		FROM expense_policies WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var p domain.ExpensePolicy
	err := row.Scan(&p.ID, &p.CompanyID, &p.Name, &p.Description, &p.IsActive, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetPolicy", err)
	}
	return &p, nil
}

func (r *PolicyRepo) ListPolicies(ctx context.Context, companyID uuid.UUID) ([]domain.ExpensePolicy, error) {
	q := `SELECT id,company_id,name,description,is_active,created_by,created_at,updated_at
		FROM expense_policies WHERE company_id=$1 ORDER BY name`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListPolicies", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ExpensePolicy, error) {
		var p domain.ExpensePolicy
		err := row.Scan(&p.ID, &p.CompanyID, &p.Name, &p.Description, &p.IsActive, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
		return p, err
	})
}

func (r *PolicyRepo) UpdatePolicy(ctx context.Context, p *domain.ExpensePolicy) error {
	q := `UPDATE expense_policies SET name=$1,description=$2,is_active=$3,updated_at=NOW() WHERE id=$4 AND company_id=$5`
	_, err := r.pool.Exec(ctx, q, p.Name, p.Description, p.IsActive, p.ID, p.CompanyID)
	return repoErr("UpdatePolicy", err)
}

func (r *PolicyRepo) CreateRule(ctx context.Context, rule *domain.ExpensePolicyRule) error {
	q := `INSERT INTO expense_policy_rules (
		id,policy_id,category_id,employee_category,max_amount,currency,
		requires_receipt,requires_approval,allowed_payment_methods,
		daily_allowance_category,conditions,priority,is_active
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`
	_, err := r.pool.Exec(ctx, q,
		rule.ID, rule.PolicyID, rule.CategoryID, rule.EmployeeCategory, rule.MaxAmount, rule.Currency,
		rule.RequiresReceipt, rule.RequiresApproval, rule.AllowedPaymentMethods,
		rule.DailyAllowanceCategory, rule.Conditions, rule.Priority, rule.IsActive)
	return repoErr("CreateRule", err)
}

func (r *PolicyRepo) ListRules(ctx context.Context, policyID uuid.UUID) ([]domain.ExpensePolicyRule, error) {
	q := `SELECT id,policy_id,category_id,employee_category,max_amount,currency,
		requires_receipt,requires_approval,allowed_payment_methods,
		daily_allowance_category,conditions,priority,is_active,created_at,updated_at
		FROM expense_policy_rules WHERE policy_id=$1 ORDER BY created_at`
	rows, err := r.pool.Query(ctx, q, policyID)
	if err != nil {
		return nil, repoErr("ListRules", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ExpensePolicyRule, error) {
		var rule domain.ExpensePolicyRule
		err := row.Scan(&rule.ID, &rule.PolicyID, &rule.CategoryID, &rule.EmployeeCategory, &rule.MaxAmount, &rule.Currency,
			&rule.RequiresReceipt, &rule.RequiresApproval, &rule.AllowedPaymentMethods,
			&rule.DailyAllowanceCategory, &rule.Conditions, &rule.Priority, &rule.IsActive, &rule.CreatedAt, &rule.UpdatedAt)
		return rule, err
	})
}

func (r *PolicyRepo) UpdateRule(ctx context.Context, rule *domain.ExpensePolicyRule) error {
	q := `UPDATE expense_policy_rules SET category_id=$1,employee_category=$2,max_amount=$3,currency=$4,
		requires_receipt=$5,requires_approval=$6,allowed_payment_methods=$7,
		daily_allowance_category=$8,conditions=$9,priority=$10,is_active=$11,updated_at=NOW()
		WHERE id=$12 AND policy_id=$13`
	_, err := r.pool.Exec(ctx, q,
		rule.CategoryID, rule.EmployeeCategory, rule.MaxAmount, rule.Currency,
		rule.RequiresReceipt, rule.RequiresApproval, rule.AllowedPaymentMethods,
		rule.DailyAllowanceCategory, rule.Conditions, rule.Priority, rule.IsActive, rule.ID, rule.PolicyID)
	return repoErr("UpdateRule", err)
}

func (r *PolicyRepo) DeleteRule(ctx context.Context, policyID, ruleID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM expense_policy_rules WHERE id=$1 AND policy_id=$2`, ruleID, policyID)
	return repoErr("DeleteRule", err)
}

func (r *PolicyRepo) GetActiveRules(ctx context.Context, companyID uuid.UUID) ([]domain.ExpensePolicyRule, error) {
	q := `SELECT pr.id,pr.policy_id,pr.category_id,pr.employee_category,pr.max_amount,pr.currency,
		pr.requires_receipt,pr.requires_approval,pr.allowed_payment_methods,
		pr.daily_allowance_category,pr.conditions,pr.priority,pr.is_active,pr.created_at,pr.updated_at
		FROM expense_policy_rules pr
		INNER JOIN expense_policies p ON p.id=pr.policy_id
		WHERE p.company_id=$1 AND pr.is_active=true AND p.is_active=true
		ORDER BY pr.created_at`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("GetActiveRules", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ExpensePolicyRule, error) {
		var rule domain.ExpensePolicyRule
		err := row.Scan(&rule.ID, &rule.PolicyID, &rule.CategoryID, &rule.EmployeeCategory, &rule.MaxAmount, &rule.Currency,
			&rule.RequiresReceipt, &rule.RequiresApproval, &rule.AllowedPaymentMethods,
			&rule.DailyAllowanceCategory, &rule.Conditions, &rule.Priority, &rule.IsActive, &rule.CreatedAt, &rule.UpdatedAt)
		return rule, err
	})
}
