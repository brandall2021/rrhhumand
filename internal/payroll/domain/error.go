package domain

import "time"

type PayrollError struct {
	ID         string     `json:"id"`
	RunID      string     `json:"run_id"`
	EmployeeID *string    `json:"employee_id,omitempty"`
	Severity   string     `json:"severity"`
	Code       string     `json:"code"`
	Message    string     `json:"message"`
	Field      *string    `json:"field,omitempty"`
	Resolved   bool       `json:"resolved"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

type PayrollAuditLog struct {
	ID         string         `json:"id"`
	CompanyID  string         `json:"company_id"`
	UserID     string         `json:"user_id"`
	Action     string         `json:"action"`
	EntityType string         `json:"entity_type"`
	EntityID   *string        `json:"entity_id,omitempty"`
	OldValues  map[string]any `json:"old_values,omitempty"`
	NewValues  map[string]any `json:"new_values,omitempty"`
	IPAddress  *string        `json:"ip_address,omitempty"`
	UserAgent  *string        `json:"user_agent,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
}
