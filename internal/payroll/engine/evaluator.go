package engine

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/rrhhumand/api/internal/payroll/domain"
	"github.com/shopspring/decimal"
)

func EvaluateConcept(ctx context.Context, concept domain.PayrollConcept, rules []domain.PayrollRule, contextData map[string]any) (*domain.PayrollItem, error) {
	if len(rules) == 0 {
		return nil, nil
	}

	rule := rules[0]

	baseSalary, _ := getDecimalContext(contextData, "base_salary")
	baseAmount := baseSalary

	if concept.BaseConceptID != nil {
		baseVal, ok := getDecimalContext(contextData, *concept.BaseConceptID)
		if ok {
			baseAmount = baseVal
		}
	}

	var rate *decimal.Decimal

	switch concept.CalculationType {
	case "AMOUNT":
		amt, err := evalAmount(rule, baseAmount)
		if err != nil {
			return nil, err
		}
		item := &domain.PayrollItem{
			ConceptID:  concept.ID,
			BaseAmount: baseAmount,
			Amount:     amt,
			Quantity:   decimal.NewFromInt(1),
			UnitValue:  baseAmount,
		}
		return item, nil

	case "PERCENTAGE":
		pct, err := evalPercentage(rule, baseAmount)
		if err != nil {
			return nil, err
		}
		r := pct.rate
		rate = &r
		item := &domain.PayrollItem{
			ConceptID:  concept.ID,
			BaseAmount: pct.base,
			Rate:       rate,
			Amount:     pct.result,
			Quantity:   decimal.NewFromInt(1),
			UnitValue:  pct.base,
		}
		return item, nil

	case "HOURLY":
		qty, val, err := evalHourly(rule, baseSalary, contextData)
		if err != nil {
			return nil, err
		}
		amount := qty.Mul(val)
		item := &domain.PayrollItem{
			ConceptID:  concept.ID,
			BaseAmount: baseSalary,
			Amount:     amount,
			Quantity:   qty,
			UnitValue:  val,
		}
		return item, nil

	case "DAILY":
		qty, val, err := evalDaily(rule, baseSalary, contextData)
		if err != nil {
			return nil, err
		}
		amount := qty.Mul(val)
		item := &domain.PayrollItem{
			ConceptID:  concept.ID,
			BaseAmount: baseSalary,
			Amount:     amount,
			Quantity:   qty,
			UnitValue:  val,
		}
		return item, nil

	case "UNIT":
		qty, val, err := evalUnit(rule, contextData)
		if err != nil {
			return nil, err
		}
		amount := qty.Mul(val)
		item := &domain.PayrollItem{
			ConceptID:  concept.ID,
			BaseAmount: baseSalary,
			Amount:     amount,
			Quantity:   qty,
			UnitValue:  val,
		}
		return item, nil

	default:
		return nil, fmt.Errorf("unsupported calculation type: %s", concept.CalculationType)
	}
}

type percentageResult struct {
	base   decimal.Decimal
	rate   decimal.Decimal
	result decimal.Decimal
}

func evalAmount(rule domain.PayrollRule, baseAmount decimal.Decimal) (decimal.Decimal, error) {
	if val, ok := GetParameter(rule.Parameters, "amount"); ok {
		return ApplyAmount(val, nil), nil
	}
	if val, ok := GetParameter(rule.Parameters, "value"); ok {
		return ApplyAmount(val, nil), nil
	}
	return baseAmount, nil
}

func evalPercentage(rule domain.PayrollRule, baseAmount decimal.Decimal) (*percentageResult, error) {
	rate := decimal.NewFromInt(100)
	if val, ok := GetParameter(rule.Parameters, "rate"); ok {
		rate = val
	}
	if val, ok := GetParameter(rule.Parameters, "percentage"); ok {
		rate = val
	}

	base := baseAmount
	if val, ok := GetParameter(rule.Parameters, "base_amount"); ok {
		base = val
	}

	limits := extractLimitsFromParams(rule.Parameters)
	result := ApplyPercentage(base, rate, limits)

	return &percentageResult{base: base, rate: rate, result: result}, nil
}

func evalHourly(rule domain.PayrollRule, baseSalary decimal.Decimal, ctx map[string]any) (decimal.Decimal, decimal.Decimal, error) {
	hours := decimal.NewFromInt(0)
	if val, ok := GetParameter(rule.Parameters, "hours"); ok {
		hours = val
	}
	if val, ok := getDecimalContext(ctx, "worked_hours"); ok {
		hours = val
	}

	workDays := 30
	if val, ok := getDecimalContext(ctx, "work_days"); ok {
		workDays = int(val.IntPart())
	}
	if workDays <= 0 {
		workDays = 30
	}

	hourlyValue := CalcHourlyValue(baseSalary, workDays)
	if val, ok := GetParameter(rule.Parameters, "hourly_value"); ok {
		hourlyValue = val
	}

	return hours, hourlyValue, nil
}

func evalDaily(rule domain.PayrollRule, baseSalary decimal.Decimal, ctx map[string]any) (decimal.Decimal, decimal.Decimal, error) {
	days := decimal.NewFromInt(0)
	if val, ok := GetParameter(rule.Parameters, "days"); ok {
		days = val
	}
	if val, ok := getDecimalContext(ctx, "worked_days"); ok {
		days = val
	}

	workDays := 30
	if val, ok := getDecimalContext(ctx, "work_days"); ok {
		workDays = int(val.IntPart())
	}
	if workDays <= 0 {
		workDays = 30
	}

	dailyValue := CalcDailyValue(baseSalary, workDays)
	if val, ok := GetParameter(rule.Parameters, "daily_value"); ok {
		dailyValue = val
	}

	return days, dailyValue, nil
}

func evalUnit(rule domain.PayrollRule, ctx map[string]any) (decimal.Decimal, decimal.Decimal, error) {
	quantity := decimal.NewFromInt(1)
	if val, ok := GetParameter(rule.Parameters, "quantity"); ok {
		quantity = val
	}
	if val, ok := getDecimalContext(ctx, "unit_quantity"); ok {
		quantity = val
	}

	unitValue := decimal.Zero
	if val, ok := GetParameter(rule.Parameters, "unit_value"); ok {
		unitValue = val
	}

	return quantity, unitValue, nil
}

func ApplyPercentage(base decimal.Decimal, rate decimal.Decimal, limits []domain.PayrollLimit) decimal.Decimal {
	result := base.Mul(rate.Div(decimal.NewFromInt(100)))
	return ApplyCap(result, limits)
}

func ApplyAmount(amount decimal.Decimal, limits []domain.PayrollLimit) decimal.Decimal {
	return ApplyCap(amount, limits)
}

var formulaRe = regexp.MustCompile(`\{\{([^}]+)\}\}`)

func ApplyFormula(formula string, context map[string]decimal.Decimal) decimal.Decimal {
	result := formulaRe.ReplaceAllStringFunc(formula, func(match string) string {
		key := strings.TrimSpace(match[2 : len(match)-2])
		if val, ok := context[key]; ok {
			return val.String()
		}
		return "0"
	})

	tokens := tokenizeExpression(result)
	return parseExpression(tokens)
}

func GetParameter(params map[string]any, key string) (decimal.Decimal, bool) {
	if params == nil {
		return decimal.Zero, false
	}
	val, ok := params[key]
	if !ok {
		return decimal.Zero, false
	}
	switch v := val.(type) {
	case float64:
		return decimal.NewFromFloat(v), true
	case int:
		return decimal.NewFromInt(int64(v)), true
	case string:
		d, err := decimal.NewFromString(v)
		if err != nil {
			return decimal.Zero, false
		}
		return d, true
	case decimal.Decimal:
		return v, true
	default:
		return decimal.Zero, false
	}
}

func ApplyCap(amount decimal.Decimal, limits []domain.PayrollLimit) decimal.Decimal {
	result := amount
	for _, limit := range limits {
		if limit.MinimumAmount != nil && result.LessThan(*limit.MinimumAmount) {
			result = *limit.MinimumAmount
		}
		if limit.MaximumAmount != nil && result.GreaterThan(*limit.MaximumAmount) {
			result = *limit.MaximumAmount
		}
	}
	return result
}

func extractLimitsFromParams(params map[string]any) []domain.PayrollLimit {
	var limits []domain.PayrollLimit
	if params == nil {
		return limits
	}

	if minVal, ok := GetParameter(params, "min_amount"); ok {
		min := minVal
		limits = append(limits, domain.PayrollLimit{
			LimitType:     "MINIMUM",
			MinimumAmount: &min,
		})
	}
	if maxVal, ok := GetParameter(params, "max_amount"); ok {
		max := maxVal
		limits = append(limits, domain.PayrollLimit{
			LimitType:     "MAXIMUM",
			MaximumAmount: &max,
		})
	}
	if capVal, ok := GetParameter(params, "cap_amount"); ok {
		cap := capVal
		limits = append(limits, domain.PayrollLimit{
			LimitType:     "MAXIMUM",
			MaximumAmount: &cap,
		})
	}

	return limits
}

func getDecimalContext(ctx map[string]any, key string) (decimal.Decimal, bool) {
	if ctx == nil {
		return decimal.Zero, false
	}
	val, ok := ctx[key]
	if !ok {
		return decimal.Zero, false
	}
	switch v := val.(type) {
	case decimal.Decimal:
		return v, true
	case float64:
		return decimal.NewFromFloat(v), true
	case string:
		d, err := decimal.NewFromString(v)
		if err != nil {
			return decimal.Zero, false
		}
		return d, true
	case int:
		return decimal.NewFromInt(int64(v)), true
	default:
		return decimal.Zero, false
	}
}

type formulaToken struct {
	kind  tokenKind
	value string
}

type tokenKind int

const (
	tokenNumber tokenKind = iota
	tokenOp
	tokenLParen
	tokenRParen
)

func tokenizeExpression(expr string) []formulaToken {
	var tokens []formulaToken
	i := 0
	runes := []rune(strings.TrimSpace(expr))

	for i < len(runes) {
		ch := runes[i]
		if ch == ' ' {
			i++
			continue
		}

		if ch >= '0' && ch <= '9' || ch == '.' {
			start := i
			for i < len(runes) && (runes[i] >= '0' && runes[i] <= '9' || runes[i] == '.') {
				i++
			}
			tokens = append(tokens, formulaToken{kind: tokenNumber, value: string(runes[start:i])})
			continue
		}

		switch ch {
		case '+', '-', '*', '/':
			tokens = append(tokens, formulaToken{kind: tokenOp, value: string(ch)})
		case '(':
			tokens = append(tokens, formulaToken{kind: tokenLParen, value: "("})
		case ')':
			tokens = append(tokens, formulaToken{kind: tokenRParen, value: ")"})
		}
		i++
	}
	return tokens
}

func parseExpression(tokens []formulaToken) decimal.Decimal {
	if len(tokens) == 0 {
		return decimal.Zero
	}

	output := make([]decimal.Decimal, 0)
	ops := make([]string, 0)

	prec := func(op string) int {
		switch op {
		case "+", "-":
			return 1
		case "*", "/":
			return 2
		}
		return 0
	}

	applyOp := func(op string) {
		if len(output) < 2 {
			return
		}
		b := output[len(output)-1]
		a := output[len(output)-2]
		output = output[:len(output)-2]

		var r decimal.Decimal
		switch op {
		case "+":
			r = a.Add(b)
		case "-":
			r = a.Sub(b)
		case "*":
			r = a.Mul(b)
		case "/":
			if b.IsZero() {
				r = decimal.Zero
			} else {
				r = a.Div(b)
			}
		}
		output = append(output, r)
	}

	for _, tok := range tokens {
		switch tok.kind {
		case tokenNumber:
			d, err := decimal.NewFromString(tok.value)
			if err != nil {
				d = decimal.Zero
			}
			output = append(output, d)
		case tokenOp:
			for len(ops) > 0 && ops[len(ops)-1] != "(" && prec(ops[len(ops)-1]) >= prec(tok.value) {
				applyOp(ops[len(ops)-1])
				ops = ops[:len(ops)-1]
			}
			ops = append(ops, tok.value)
		case tokenLParen:
			ops = append(ops, "(")
		case tokenRParen:
			for len(ops) > 0 && ops[len(ops)-1] != "(" {
				applyOp(ops[len(ops)-1])
				ops = ops[:len(ops)-1]
			}
			if len(ops) > 0 && ops[len(ops)-1] == "(" {
				ops = ops[:len(ops)-1]
			}
		}
	}

	for len(ops) > 0 {
		applyOp(ops[len(ops)-1])
		ops = ops[:len(ops)-1]
	}

	if len(output) == 0 {
		return decimal.Zero
	}
	return output[0]
}


