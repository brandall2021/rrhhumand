package payroll

import (
	"fmt"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateCreatePeriod(req *CreatePeriodRequest) error {
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if req.StartDate == "" || req.EndDate == "" {
		return fmt.Errorf("start_date and end_date are required")
	}
	return nil
}

func (v *Validator) ValidateCanCalculate(period *PayrollPeriod) error {
	if period.Status != "OPEN" {
		return fmt.Errorf("period must be OPEN to calculate, current status: %s", period.Status)
	}
	return nil
}

func (v *Validator) ValidateCanApprove(period *PayrollPeriod) error {
	if period.Status != "REVIEW" {
		return fmt.Errorf("period must be in REVIEW to approve, current status: %s", period.Status)
	}
	return nil
}

func (v *Validator) ValidateCanClose(period *PayrollPeriod) error {
	if period.Status != "APPROVED" {
		return fmt.Errorf("period must be APPROVED to close, current status: %s", period.Status)
	}
	return nil
}

func (v *Validator) ValidateNoErrors(review *PayrollReview) error {
	if review.Errors > 0 {
		return fmt.Errorf("cannot approve with %d errors", review.Errors)
	}
	return nil
}

func (v *Validator) ValidateCompensation(req *SetCompensationRequest) error {
	if req.EmployeeID == "" {
		return fmt.Errorf("employee_id is required")
	}
	if req.BaseAmount < 0 {
		return fmt.Errorf("base_amount cannot be negative")
	}
	if req.EffectiveFrom == "" {
		return fmt.Errorf("effective_from is required")
	}
	return nil
}
