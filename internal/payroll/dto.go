package payroll

type CreatePeriodRequest struct {
	Name      string `json:"name" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

type UpdatePeriodRequest struct {
	Name      *string `json:"name"`
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
}

type CreateConceptRequest struct {
	Code            string `json:"code" binding:"required"`
	Name            string `json:"name" binding:"required"`
	Type            string `json:"type" binding:"required"`
	CalculationType *string `json:"calculation_type"`
	Taxable         *bool   `json:"taxable"`
}

type UpdateConceptRequest struct {
	Name            *string `json:"name"`
	Type            *string `json:"type"`
	CalculationType *string `json:"calculation_type"`
	Taxable         *bool   `json:"taxable"`
	Active          *bool   `json:"active"`
}

type SetCompensationRequest struct {
	EmployeeID    string  `json:"employee_id" binding:"required"`
	BaseAmount    float64 `json:"base_amount" binding:"required"`
	Currency      *string `json:"currency"`
	EffectiveFrom string  `json:"effective_from" binding:"required"`
	Reason        *string `json:"reason"`
}

type CreateBenefitRequest struct {
	Code          string   `json:"code" binding:"required"`
	Name          string   `json:"name" binding:"required"`
	Description   *string  `json:"description"`
	BenefitType   *string  `json:"benefit_type"`
	DefaultAmount *float64 `json:"default_amount"`
}

type AssignBenefitRequest struct {
	EmployeeID    string  `json:"employee_id" binding:"required"`
	BenefitID     string  `json:"benefit_id" binding:"required"`
	Amount        *float64 `json:"amount"`
	Currency      *string  `json:"currency"`
	EffectiveFrom string   `json:"effective_from" binding:"required"`
}

type CreateBonusRequest struct {
	EmployeeID  string  `json:"employee_id" binding:"required"`
	BonusType   string  `json:"bonus_type" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	Currency    *string `json:"currency"`
	Reason      *string `json:"reason"`
	PeriodStart *string `json:"period_start"`
	PeriodEnd   *string `json:"period_end"`
}

type ApproveBonusRequest struct {
	Comments *string `json:"comments"`
}

type CreateAdvanceRequest struct {
	EmployeeID  string  `json:"employee_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	Currency    *string `json:"currency"`
	RequestDate string  `json:"request_date" binding:"required"`
	Reason      *string `json:"reason"`
}

type ApproveAdvanceRequest struct {
	Comments *string `json:"comments"`
}

type CreateDeductionRequest struct {
	EmployeeID  string  `json:"employee_id" binding:"required"`
	Concept     string  `json:"concept" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	Currency    *string `json:"currency"`
	Reason      *string `json:"reason"`
	PeriodStart *string `json:"period_start"`
	PeriodEnd   *string `json:"period_end"`
}

type CreateAdjustmentRequest struct {
	EmployeeID string  `json:"employee_id" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
	Reason     string  `json:"reason" binding:"required"`
	Type       *string `json:"type"`
}

type CalculatePeriodRequest struct {
	DateFrom string `json:"date_from" binding:"required"`
	DateTo   string `json:"date_to" binding:"required"`
}
