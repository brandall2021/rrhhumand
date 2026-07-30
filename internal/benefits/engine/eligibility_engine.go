package engine

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/benefits/domain"
	"github.com/shopspring/decimal"
)

type EmployeeData struct {
	ID             uuid.UUID
	CompanyID      uuid.UUID
	PositionID     *uuid.UUID
	DepartmentID   *uuid.UUID
	BranchID       *uuid.UUID
	EmploymentType string
	ContractType   string
	WorkSchedule   string
	AdmissionDate  time.Time
	Salary         decimal.Decimal
	BirthDate      *time.Time
	Gender         *string
	JobLevel       *string
	Location       *string
	Performance    *string
}

type employeeRepo interface {
	GetEmployee(ctx context.Context, id uuid.UUID) (*EmployeeData, error)
}

type EligibilityEngine struct {
	employeeRepo employeeRepo
}

func NewEligibilityEngine(er employeeRepo) *EligibilityEngine {
	return &EligibilityEngine{employeeRepo: er}
}

func engErr(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("benefits_engine.eligibility.%s: %w", op, err)
}

func (e *EligibilityEngine) Evaluate(ctx context.Context, employeeID uuid.UUID, rules []domain.BenefitEligibilityRule) (bool, []string, error) {
	if len(rules) == 0 {
		return true, nil, nil
	}

	emp, err := e.employeeRepo.GetEmployee(ctx, employeeID)
	if err != nil {
		return false, nil, engErr("Evaluate.getEmployee", err)
	}

	groups := make(map[int][]domain.BenefitEligibilityRule)
	for _, rule := range rules {
		groups[rule.LogicGroup] = append(groups[rule.LogicGroup], rule)
	}

	var failures []string

	for groupID, groupRules := range groups {
		groupPass := false
		groupOp := groupRules[0].LogicOperator

		for _, rule := range groupRules {
			pass, err := evaluateRule(emp, rule)
			if err != nil {
				return false, nil, engErr("Evaluate.evaluateRule", err)
			}
			if groupOp == "OR" {
				if pass {
					groupPass = true
					break
				}
			} else {
				if !pass {
					msg := fmt.Sprintf("group_%d: rule %s %s %s failed", groupID, rule.RuleType, rule.Operator, rule.Value)
					if rule.ErrorMessage != nil {
						msg = *rule.ErrorMessage
					}
					failures = append(failures, msg)
					groupPass = false
					break
				}
				groupPass = true
			}
		}

		if groupOp == "OR" && !groupPass {
			msg := fmt.Sprintf("group_%d: no OR rules matched", groupID)
			if len(groupRules) > 0 && groupRules[0].ErrorMessage != nil {
				msg = *groupRules[0].ErrorMessage
			}
			failures = append(failures, msg)
		}

		if !groupPass {
			return false, failures, nil
		}
	}

	if len(failures) > 0 {
		return false, failures, nil
	}
	return true, nil, nil
}

func evaluateRule(emp *EmployeeData, rule domain.BenefitEligibilityRule) (bool, error) {
	var fieldVal string
	switch rule.RuleType {
	case "JOB_LEVEL":
		if emp.JobLevel != nil {
			fieldVal = *emp.JobLevel
		}
	case "DEPARTMENT":
		if emp.DepartmentID != nil {
			fieldVal = emp.DepartmentID.String()
		}
	case "BRANCH":
		if emp.BranchID != nil {
			fieldVal = emp.BranchID.String()
		}
	case "POSITION":
		if emp.PositionID != nil {
			fieldVal = emp.PositionID.String()
		}
	case "EMPLOYMENT_TYPE":
		fieldVal = emp.EmploymentType
	case "SENIORITY":
		fieldVal = fmt.Sprintf("%d", int(time.Since(emp.AdmissionDate).Hours()/24/30))
	case "AGE":
		if emp.BirthDate != nil {
			age := int(time.Since(*emp.BirthDate).Hours() / 24 / 365)
			fieldVal = fmt.Sprintf("%d", age)
		}
	case "CONTRACT_TYPE":
		fieldVal = emp.ContractType
	case "WORK_SCHEDULE":
		fieldVal = emp.WorkSchedule
	case "LOCATION":
		if emp.Location != nil {
			fieldVal = *emp.Location
		}
	case "TENURE":
		fieldVal = fmt.Sprintf("%d", int(time.Since(emp.AdmissionDate).Hours()/24/30))
	case "SALARY_BAND":
		fieldVal = emp.Salary.String()
	case "GENDER":
		if emp.Gender != nil {
			fieldVal = *emp.Gender
		}
	case "PERFORMANCE_RATING":
		if emp.Performance != nil {
			fieldVal = *emp.Performance
		}
	case "CUSTOM":
		fieldVal = rule.Value
	default:
		return false, fmt.Errorf("unknown rule type: %s", rule.RuleType)
	}

	return compareValues(fieldVal, rule.Operator, rule.Value, rule.ValueTo)
}

func compareValues(fieldVal, operator, value string, valueTo *string) (bool, error) {
	switch operator {
	case "EQ":
		return fieldVal == value, nil
	case "NEQ":
		return fieldVal != value, nil
	case "GT":
		return numericCompare(fieldVal, value, func(a, b float64) bool { return a > b })
	case "GTE":
		return numericCompare(fieldVal, value, func(a, b float64) bool { return a >= b })
	case "LT":
		return numericCompare(fieldVal, value, func(a, b float64) bool { return a < b })
	case "LTE":
		return numericCompare(fieldVal, value, func(a, b float64) bool { return a <= b })
	case "IN":
		parts := strings.Split(value, ",")
		for _, p := range parts {
			if strings.TrimSpace(p) == fieldVal {
				return true, nil
			}
		}
		return false, nil
	case "NOT_IN":
		parts := strings.Split(value, ",")
		for _, p := range parts {
			if strings.TrimSpace(p) == fieldVal {
				return false, nil
			}
		}
		return true, nil
	case "BETWEEN":
		if valueTo == nil {
			return false, fmt.Errorf("BETWEEN requires value_to")
		}
		ok, err := numericCompare(fieldVal, value, func(a, b float64) bool { return a >= b })
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
		return numericCompare(fieldVal, *valueTo, func(a, b float64) bool { return a <= b })
	case "CONTAINS":
		return strings.Contains(fieldVal, value), nil
	default:
		return false, fmt.Errorf("unknown operator: %s", operator)
	}
}

func numericCompare(fieldVal, value string, cmp func(float64, float64) bool) (bool, error) {
	a, err := strconv.ParseFloat(fieldVal, 64)
	if err != nil {
		return false, fmt.Errorf("cannot parse field value %q as number: %w", fieldVal, err)
	}
	b, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return false, fmt.Errorf("cannot parse rule value %q as number: %w", value, err)
	}
	return cmp(a, b), nil
}
