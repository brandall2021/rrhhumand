package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ReceiptTemplate struct {
	ID             uuid.UUID  `json:"id"`
	CompanyID      uuid.UUID  `json:"company_id"`
	Name           string     `json:"name"`
	Description    *string    `json:"description,omitempty"`
	TemplateHTML   string     `json:"template_html"`
	TemplateCSS    *string    `json:"template_css,omitempty"`
	Orientation    string     `json:"orientation"`
	PaperSize      string     `json:"paper_size"`
	ShowLogo       bool       `json:"show_logo"`
	ShowSignature  bool       `json:"show_signature"`
	ShowQR         bool       `json:"show_qr"`
	ShowBarcode    bool       `json:"show_barcode"`
	FontFamily     string     `json:"font_family"`
	FontSize       int        `json:"font_size"`
	PrimaryColor   string     `json:"primary_color"`
	SecondaryColor string     `json:"secondary_color"`
	MarginTop      float64    `json:"margin_top"`
	MarginBottom   float64    `json:"margin_bottom"`
	MarginLeft     float64    `json:"margin_left"`
	MarginRight    float64    `json:"margin_right"`
	IsDefault      bool       `json:"is_default"`
	IsActive       bool       `json:"is_active"`
	CreatedBy      uuid.UUID  `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type Receipt struct {
	ID                   uuid.UUID       `json:"id"`
	CompanyID            uuid.UUID       `json:"company_id"`
	RunID                uuid.UUID       `json:"run_id"`
	RunEmployeeID        uuid.UUID       `json:"run_employee_id"`
	EmployeeID           uuid.UUID       `json:"employee_id"`
	TemplateID           *uuid.UUID      `json:"template_id,omitempty"`
	ReceiptNumber        string          `json:"receipt_number"`
	CUIT                 string          `json:"cuit"`
	EmployeeCUIL         string          `json:"employee_cuil"`
	PeriodName           string          `json:"period_name"`
	PeriodStart          time.Time       `json:"period_start"`
	PeriodEnd            time.Time       `json:"period_end"`
	PaymentDate          *time.Time      `json:"payment_date,omitempty"`
	GrossRemunerative    decimal.Decimal `json:"gross_remunerative"`
	GrossNonRemunerative decimal.Decimal `json:"gross_non_remunerative"`
	DeductionsTotal      decimal.Decimal `json:"deductions_total"`
	ContributionsTotal   decimal.Decimal `json:"contributions_total"`
	NetAmount            decimal.Decimal `json:"net_amount"`
	EmployerCost         decimal.Decimal `json:"employer_cost"`
	Currency             string          `json:"currency"`
	AmountInWords        *string         `json:"amount_in_words,omitempty"`
	DigitalToken         *string         `json:"digital_token,omitempty"`
	QRCode               *string         `json:"qr_code,omitempty"`
	Barcode              *string         `json:"barcode,omitempty"`
	Status               string          `json:"status"`
	AcknowledgedAt       *time.Time      `json:"acknowledged_at,omitempty"`
	AcknowledgedIP       *string         `json:"acknowledged_ip,omitempty"`
	ViewedAt             *time.Time      `json:"viewed_at,omitempty"`
	DownloadedAt         *time.Time      `json:"downloaded_at,omitempty"`
	EmailedAt            *time.Time      `json:"emailed_at,omitempty"`
	StoragePath          *string         `json:"storage_path,omitempty"`
	GeneratedBy          uuid.UUID       `json:"generated_by"`
	GeneratedAt          time.Time       `json:"generated_at"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}

type ReceiptItem struct {
	ID             uuid.UUID       `json:"id"`
	ReceiptID      uuid.UUID       `json:"receipt_id"`
	ConceptCode    string          `json:"concept_code"`
	ConceptName    string          `json:"concept_name"`
	Quantity       decimal.Decimal `json:"quantity"`
	UnitValue      decimal.Decimal `json:"unit_value"`
	BaseAmount     decimal.Decimal `json:"base_amount"`
	Rate           *decimal.Decimal `json:"rate,omitempty"`
	Amount         decimal.Decimal `json:"amount"`
	IsRemunerative bool            `json:"is_remunerative"`
	IsDeduction    bool            `json:"is_deduction"`
	IsContribution bool            `json:"is_contribution"`
	SortOrder      int             `json:"sort_order"`
	CreatedAt      time.Time       `json:"created_at"`
}
