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
	q := `INSERT INTO expense_policy_rules (id,policy_id,category_id,rule_type,operator,value,
		currency,description,is_active,effective_from,effective_to,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := r.pool.Exec(ctx, q, rule.ID, rule.PolicyID, rule.CategoryID, rule.RuleType, rule.Operator,
		rule.Value, rule.Currency, rule.Description, rule.IsActive, rule.EffectiveFrom, rule.EffectiveTo, rule.CreatedBy)
	return repoErr("CreateRule", err)
}

func (r *PolicyRepo) ListRules(ctx context.Context, policyID uuid.UUID) ([]domain.ExpensePolicyRule, error) {
	q := `SELECT id,policy_id,category_id,rule_type,operator,value,currency,description,
		is_active,effective_from,effective_to,created_by,created_at,updated_at
		FROM expense_policy_rules WHERE policy_id=$1 ORDER BY created_at`
	rows, err := r.pool.Query(ctx, q, policyID)
	if err != nil {
		return nil, repoErr("ListRules", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ExpensePolicyRule, error) {
		var rule domain.ExpensePolicyRule
		err := row.Scan(&rule.ID, &rule.PolicyID, &rule.CategoryID, &rule.RuleType, &rule.Operator,
			&rule.Value, &rule.Currency, &rule.Description, &rule.IsActive, &rule.EffectiveFrom, &rule.EffectiveTo,
			&rule.CreatedBy, &rule.CreatedAt, &rule.UpdatedAt)
		return rule, err
	})
}

func (r *PolicyRepo) UpdateRule(ctx context.Context, rule *domain.ExpensePolicyRule) error {
	q := `UPDATE expense_policy_rules SET category_id=$1,rule_type=$2,operator=$3,value=$4,
		currency=$5,description=$6,is_active=$7,effective_from=$8,effective_to=$9,updated_at=NOW()
		WHERE id=$10 AND policy_id=$11`
	_, err := r.pool.Exec(ctx, q, rule.CategoryID, rule.RuleType, rule.Operator, rule.Value,
		rule.Currency, rule.Description, rule.IsActive, rule.EffectiveFrom, rule.EffectiveTo, rule.ID, rule.PolicyID)
	return repoErr("UpdateRule", err)
}

func (r *PolicyRepo) DeleteRule(ctx context.Context, policyID, ruleID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM expense_policy_rules WHERE id=$1 AND policy_id=$2`, ruleID, policyID)
	return repoErr("DeleteRule", err)
}

func (r *PolicyRepo) GetActiveRules(ctx context.Context, companyID uuid.UUID, date time.Time) ([]domain.ExpensePolicyRule, error) {
	q := `SELECT pr.id,pr.policy_id,pr.category_id,pr.rule_type,pr.operator,pr.value,pr.currency,pr.description,
		pr.is_active,pr.effective_from,pr.effective_to,pr.created_by,pr.created_at,pr.updated_at
		FROM expense_policy_rules pr
		INNER JOIN expense_policies p ON p.id=pr.policy_id
		WHERE p.company_id=$1 AND pr.is_active=true AND p.is_active=true
		AND (pr.effective_from IS NULL OR pr.effective_from<=$2)
		AND (pr.effective_to IS NULL OR pr.effective_to>=$2)
		ORDER BY pr.created_at`
	rows, err := r.pool.Query(ctx, q, companyID, date)
	if err != nil {
		return nil, repoErr("GetActiveRules", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ExpensePolicyRule, error) {
		var rule domain.ExpensePolicyRule
		err := row.Scan(&rule.ID, &rule.PolicyID, &rule.CategoryID, &rule.RuleType, &rule.Operator,
			&rule.Value, &rule.Currency, &rule.Description, &rule.IsActive, &rule.EffectiveFrom, &rule.EffectiveTo,
			&rule.CreatedBy, &rule.CreatedAt, &rule.UpdatedAt)
		return rule, err
	})
}
