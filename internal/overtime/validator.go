package overtime

import (
	"fmt"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateCreatePolicy(req *CreateOvertimePolicyRequest) error {
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if req.RoundingMinutes != nil && *req.RoundingMinutes <= 0 {
		return fmt.Errorf("rounding_minutes must be positive")
	}
	if req.WeekendMultiplier != nil && *req.WeekendMultiplier < 1.0 {
		return fmt.Errorf("weekend_multiplier must be >= 1.0")
	}
	if req.HolidayMultiplier != nil && *req.HolidayMultiplier < 1.0 {
		return fmt.Errorf("holiday_multiplier must be >= 1.0")
	}
	return nil
}

func (v *Validator) ValidateRequest(req *RequestOvertimeRequest) error {
	if req.RequestedMinutes <= 0 {
		return fmt.Errorf("requested_minutes must be positive")
	}
	if req.Reason == "" {
		return fmt.Errorf("reason is required")
	}
	return nil
}

func (v *Validator) ValidateCompensation(req *RequestCompensationRequest, balance *EmployeeTimeBalance) error {
	if req.Minutes <= 0 {
		return fmt.Errorf("minutes must be positive")
	}
	if balance == nil || balance.BalanceMinutes < req.Minutes {
		return fmt.Errorf("insufficient balance")
	}
	return nil
}

func (v *Validator) ValidateApproval(record *OvertimeRecord, approvedMinutes int) error {
	if approvedMinutes < 0 {
		return fmt.Errorf("approved_minutes cannot be negative")
	}
	if approvedMinutes > record.OvertimeMinutes {
		return fmt.Errorf("approved_minutes cannot exceed overtime_minutes")
	}
	return nil
}
