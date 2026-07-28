package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/benefits/domain"
)

type EligibilityRepo struct {
	pool *pgxpool.Pool
}

func NewEligibilityRepo(pool *pgxpool.Pool) *EligibilityRepo {
	return &EligibilityRepo{pool: pool}
}

func (r *EligibilityRepo) CreateRule(ctx context.Context, rule *domain.BenefitEligibilityRule) error {
	q := `INSERT INTO benefit_eligibility_rules (id,company_id,benefit_id,rule_type,operator,value,value_to,
		logic_group,logic_operator,priority,error_message,is_active,effective_from,effective_to,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`
	_, err := r.pool.Exec(ctx, q, rule.ID, rule.CompanyID, rule.BenefitID, rule.RuleType, rule.Operator, rule.Value, rule.ValueTo,
		rule.LogicGroup, rule.LogicOperator, rule.Priority, rule.ErrorMessage, rule.IsActive, rule.EffectiveFrom, rule.EffectiveTo, rule.CreatedBy)
	return repoErr("CreateRule", err)
}

func (r *EligibilityRepo) GetRule(ctx context.Context, companyID, id uuid.UUID) (*domain.BenefitEligibilityRule, error) {
	q := `SELECT id,company_id,benefit_id,rule_type,operator,value,value_to,logic_group,logic_operator,priority,
		error_message,is_active,effective_from,effective_to,created_by,created_at,updated_at
		FROM benefit_eligibility_rules WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var rule domain.BenefitEligibilityRule
	err := row.Scan(&rule.ID, &rule.CompanyID, &rule.BenefitID, &rule.RuleType, &rule.Operator, &rule.Value, &rule.ValueTo,
		&rule.LogicGroup, &rule.LogicOperator, &rule.Priority, &rule.ErrorMessage, &rule.IsActive, &rule.EffectiveFrom, &rule.EffectiveTo, &rule.CreatedBy, &rule.CreatedAt, &rule.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetRule", err)
	}
	return &rule, nil
}

func (r *EligibilityRepo) ListRules(ctx context.Context, benefitID uuid.UUID) ([]domain.BenefitEligibilityRule, error) {
	q := `SELECT id,company_id,benefit_id,rule_type,operator,value,value_to,logic_group,logic_operator,priority,
		error_message,is_active,effective_from,effective_to,created_by,created_at,updated_at
		FROM benefit_eligibility_rules WHERE benefit_id=$1 ORDER BY priority,logic_group`
	rows, err := r.pool.Query(ctx, q, benefitID)
	if err != nil {
		return nil, repoErr("ListRules", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.BenefitEligibilityRule, error) {
		var rule domain.BenefitEligibilityRule
		err := row.Scan(&rule.ID, &rule.CompanyID, &rule.BenefitID, &rule.RuleType, &rule.Operator, &rule.Value, &rule.ValueTo,
			&rule.LogicGroup, &rule.LogicOperator, &rule.Priority, &rule.ErrorMessage, &rule.IsActive, &rule.EffectiveFrom, &rule.EffectiveTo, &rule.CreatedBy, &rule.CreatedAt, &rule.UpdatedAt)
		return rule, err
	})
}

func (r *EligibilityRepo) UpdateRule(ctx context.Context, rule *domain.BenefitEligibilityRule) error {
	q := `UPDATE benefit_eligibility_rules SET rule_type=$1,operator=$2,value=$3,value_to=$4,logic_group=$5,
		logic_operator=$6,priority=$7,error_message=$8,is_active=$9,effective_from=$10,effective_to=$11,updated_at=NOW()
		WHERE id=$12 AND company_id=$13`
	_, err := r.pool.Exec(ctx, q, rule.RuleType, rule.Operator, rule.Value, rule.ValueTo,
		rule.LogicGroup, rule.LogicOperator, rule.Priority, rule.ErrorMessage, rule.IsActive, rule.EffectiveFrom, rule.EffectiveTo, rule.ID, rule.CompanyID)
	return repoErr("UpdateRule", err)
}

func (r *EligibilityRepo) DeleteRule(ctx context.Context, companyID, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM benefit_eligibility_rules WHERE id=$1 AND company_id=$2`, id, companyID)
	return repoErr("DeleteRule", err)
}

func (r *EligibilityRepo) ListByType(ctx context.Context, companyID uuid.UUID, ruleType string) ([]domain.BenefitEligibilityRule, error) {
	q := `SELECT id,company_id,benefit_id,rule_type,operator,value,value_to,logic_group,logic_operator,priority,
		error_message,is_active,effective_from,effective_to,created_by,created_at,updated_at
		FROM benefit_eligibility_rules WHERE company_id=$1 AND rule_type=$2 ORDER BY priority,logic_group`
	rows, err := r.pool.Query(ctx, q, companyID, ruleType)
	if err != nil {
		return nil, repoErr("ListByType", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.BenefitEligibilityRule, error) {
		var rule domain.BenefitEligibilityRule
		err := row.Scan(&rule.ID, &rule.CompanyID, &rule.BenefitID, &rule.RuleType, &rule.Operator, &rule.Value, &rule.ValueTo,
			&rule.LogicGroup, &rule.LogicOperator, &rule.Priority, &rule.ErrorMessage, &rule.IsActive, &rule.EffectiveFrom, &rule.EffectiveTo, &rule.CreatedBy, &rule.CreatedAt, &rule.UpdatedAt)
		return rule, err
	})
}
