package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/payroll/features/domain"
)

func repoErr(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("payroll_features_repo.%s: %w", op, err)
}

type ReceiptRepo struct {
	pool *pgxpool.Pool
}

func NewReceiptRepo(pool *pgxpool.Pool) *ReceiptRepo {
	return &ReceiptRepo{pool: pool}
}

func (r *ReceiptRepo) CreateTemplate(ctx context.Context, t *domain.ReceiptTemplate) error {
	q := `INSERT INTO payroll_receipt_templates (id,company_id,name,description,template_html,template_css,
		orientation,paper_size,show_logo,show_signature,show_qr,show_barcode,font_family,font_size,
		primary_color,secondary_color,margin_top,margin_bottom,margin_left,margin_right,is_default,is_active,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23)`
	_, err := r.pool.Exec(ctx, q, t.ID, t.CompanyID, t.Name, t.Description, t.TemplateHTML, t.TemplateCSS,
		t.Orientation, t.PaperSize, t.ShowLogo, t.ShowSignature, t.ShowQR, t.ShowBarcode, t.FontFamily, t.FontSize,
		t.PrimaryColor, t.SecondaryColor, t.MarginTop, t.MarginBottom, t.MarginLeft, t.MarginRight, t.IsDefault, t.IsActive, t.CreatedBy)
	return repoErr("CreateTemplate", err)
}

func (r *ReceiptRepo) GetTemplate(ctx context.Context, companyID, id uuid.UUID) (*domain.ReceiptTemplate, error) {
	q := `SELECT id,company_id,name,description,template_html,template_css,orientation,paper_size,
		show_logo,show_signature,show_qr,show_barcode,font_family,font_size,primary_color,secondary_color,
		margin_top,margin_bottom,margin_left,margin_right,is_default,is_active,created_by,created_at,updated_at
		FROM payroll_receipt_templates WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	var t domain.ReceiptTemplate
	err := row.Scan(&t.ID, &t.CompanyID, &t.Name, &t.Description, &t.TemplateHTML, &t.TemplateCSS,
		&t.Orientation, &t.PaperSize, &t.ShowLogo, &t.ShowSignature, &t.ShowQR, &t.ShowBarcode, &t.FontFamily, &t.FontSize,
		&t.PrimaryColor, &t.SecondaryColor, &t.MarginTop, &t.MarginBottom, &t.MarginLeft, &t.MarginRight, &t.IsDefault, &t.IsActive, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetTemplate", err)
	}
	return &t, nil
}

func (r *ReceiptRepo) ListTemplates(ctx context.Context, companyID uuid.UUID) ([]domain.ReceiptTemplate, error) {
	q := `SELECT id,company_id,name,description,template_html,template_css,orientation,paper_size,
		show_logo,show_signature,show_qr,show_barcode,font_family,font_size,primary_color,secondary_color,
		margin_top,margin_bottom,margin_left,margin_right,is_default,is_active,created_by,created_at,updated_at
		FROM payroll_receipt_templates WHERE company_id=$1 AND is_active=true ORDER BY name`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListTemplates", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ReceiptTemplate, error) {
		var t domain.ReceiptTemplate
		err := row.Scan(&t.ID, &t.CompanyID, &t.Name, &t.Description, &t.TemplateHTML, &t.TemplateCSS,
			&t.Orientation, &t.PaperSize, &t.ShowLogo, &t.ShowSignature, &t.ShowQR, &t.ShowBarcode, &t.FontFamily, &t.FontSize,
			&t.PrimaryColor, &t.SecondaryColor, &t.MarginTop, &t.MarginBottom, &t.MarginLeft, &t.MarginRight, &t.IsDefault, &t.IsActive, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt)
		return t, err
	})
}

func (r *ReceiptRepo) UpdateTemplate(ctx context.Context, t *domain.ReceiptTemplate) error {
	q := `UPDATE payroll_receipt_templates SET name=$1,description=$2,template_html=$3,template_css=$4,
		orientation=$5,show_logo=$6,show_signature=$7,show_qr=$8,show_barcode=$9,
		font_family=$10,font_size=$11,primary_color=$12,secondary_color=$13,
		margin_top=$14,margin_bottom=$15,margin_left=$16,margin_right=$17,
		is_default=$18,is_active=$19,updated_at=NOW() WHERE id=$20 AND company_id=$21`
	_, err := r.pool.Exec(ctx, q, t.Name, t.Description, t.TemplateHTML, t.TemplateCSS,
		t.Orientation, t.ShowLogo, t.ShowSignature, t.ShowQR, t.ShowBarcode,
		t.FontFamily, t.FontSize, t.PrimaryColor, t.SecondaryColor,
		t.MarginTop, t.MarginBottom, t.MarginLeft, t.MarginRight,
		t.IsDefault, t.IsActive, t.ID, t.CompanyID)
	return repoErr("UpdateTemplate", err)
}

func (r *ReceiptRepo) DeleteTemplate(ctx context.Context, companyID, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM payroll_receipt_templates WHERE id=$1 AND company_id=$2`, id, companyID)
	return repoErr("DeleteTemplate", err)
}

func (r *ReceiptRepo) CreateReceipt(ctx context.Context, rec *domain.Receipt) error {
	q := `INSERT INTO payroll_receipts (id,company_id,run_id,run_employee_id,employee_id,template_id,
		receipt_number,cuit,employee_cuil,period_name,period_start,period_end,payment_date,
		gross_remunerative,gross_non_remunerative,deductions_total,contributions_total,net_amount,employer_cost,
		currency,amount_in_words,digital_token,qr_code,barcode,status,generated_by,storage_path)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27)`
	_, err := r.pool.Exec(ctx, q, rec.ID, rec.CompanyID, rec.RunID, rec.RunEmployeeID, rec.EmployeeID, rec.TemplateID,
		rec.ReceiptNumber, rec.CUIT, rec.EmployeeCUIL, rec.PeriodName, rec.PeriodStart, rec.PeriodEnd, rec.PaymentDate,
		rec.GrossRemunerative, rec.GrossNonRemunerative, rec.DeductionsTotal, rec.ContributionsTotal, rec.NetAmount, rec.EmployerCost,
		rec.Currency, rec.AmountInWords, rec.DigitalToken, rec.QRCode, rec.Barcode, rec.Status, rec.GeneratedBy, rec.StoragePath)
	return repoErr("CreateReceipt", err)
}

func (r *ReceiptRepo) GetReceipt(ctx context.Context, companyID, id uuid.UUID) (*domain.Receipt, error) {
	q := `SELECT id,company_id,run_id,run_employee_id,employee_id,template_id,receipt_number,cuit,employee_cuil,
		period_name,period_start,period_end,payment_date,gross_remunerative,gross_non_remunerative,
		deductions_total,contributions_total,net_amount,employer_cost,currency,amount_in_words,
		digital_token,qr_code,barcode,status,acknowledged_at,acknowledged_ip,viewed_at,downloaded_at,
		emailed_at,storage_path,generated_by,generated_at,created_at,updated_at
		FROM payroll_receipts WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	return scanReceipt(row)
}

func (r *ReceiptRepo) ListReceipts(ctx context.Context, companyID uuid.UUID, runID *uuid.UUID, employeeID *uuid.UUID, limit, offset int) ([]domain.Receipt, error) {
	q := `SELECT id,company_id,run_id,run_employee_id,employee_id,template_id,receipt_number,cuit,employee_cuil,
		period_name,period_start,period_end,payment_date,gross_remunerative,gross_non_remunerative,
		deductions_total,contributions_total,net_amount,employer_cost,currency,amount_in_words,
		digital_token,qr_code,barcode,status,acknowledged_at,acknowledged_ip,viewed_at,downloaded_at,
		emailed_at,storage_path,generated_by,generated_at,created_at,updated_at
		FROM payroll_receipts WHERE company_id=$1`
	args := []any{companyID}
	n := 2
	if runID != nil {
		q += fmt.Sprintf(" AND run_id=$%d", n)
		args = append(args, *runID)
		n++
	}
	if employeeID != nil {
		q += fmt.Sprintf(" AND employee_id=$%d", n)
		args = append(args, *employeeID)
		n++
	}
	q += " ORDER BY created_at DESC"
	if limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", n)
		args = append(args, limit)
		n++
	}
	if offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", n)
		args = append(args, offset)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListReceipts", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Receipt, error) {
		return scanReceipt(row)
	})
}

func (r *ReceiptRepo) UpdateReceiptStatus(ctx context.Context, id uuid.UUID, status string, fields map[string]any) error {
	q := `UPDATE payroll_receipts SET status=$1`
	args := []any{status}
	n := 2
	for k, v := range fields {
		q += fmt.Sprintf(",%s=$%d", k, n)
		args = append(args, v)
		n++
	}
	q += fmt.Sprintf(",updated_at=NOW() WHERE id=$%d", n)
	args = append(args, id)
	_, err := r.pool.Exec(ctx, q, args...)
	return repoErr("UpdateReceiptStatus", err)
}

func (r *ReceiptRepo) CreateReceiptItems(ctx context.Context, items []domain.ReceiptItem) error {
	if len(items) == 0 {
		return nil
	}
	q := `INSERT INTO payroll_receipt_items (id,receipt_id,concept_code,concept_name,quantity,unit_value,base_amount,rate,amount,is_remunerative,is_deduction,is_contribution,sort_order) VALUES `
	args := []any{}
	n := 1
	for _, it := range items {
		q += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),", n, n+1, n+2, n+3, n+4, n+5, n+6, n+7, n+8, n+9, n+10, n+11, n+12)
		args = append(args, it.ID, it.ReceiptID, it.ConceptCode, it.ConceptName, it.Quantity, it.UnitValue, it.BaseAmount, it.Rate, it.Amount, it.IsRemunerative, it.IsDeduction, it.IsContribution, it.SortOrder)
		n += 13
	}
	q = q[:len(q)-1]
	_, err := r.pool.Exec(ctx, q, args...)
	return repoErr("CreateReceiptItems", err)
}

func (r *ReceiptRepo) ListReceiptItems(ctx context.Context, receiptID uuid.UUID) ([]domain.ReceiptItem, error) {
	q := `SELECT id,receipt_id,concept_code,concept_name,quantity,unit_value,base_amount,rate,amount,is_remunerative,is_deduction,is_contribution,sort_order,created_at
		FROM payroll_receipt_items WHERE receipt_id=$1 ORDER BY sort_order`
	rows, err := r.pool.Query(ctx, q, receiptID)
	if err != nil {
		return nil, repoErr("ListReceiptItems", err)
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.ReceiptItem, error) {
		var it domain.ReceiptItem
		err := row.Scan(&it.ID, &it.ReceiptID, &it.ConceptCode, &it.ConceptName, &it.Quantity, &it.UnitValue, &it.BaseAmount, &it.Rate, &it.Amount, &it.IsRemunerative, &it.IsDeduction, &it.IsContribution, &it.SortOrder, &it.CreatedAt)
		return it, err
	})
}

func scanReceipt(row pgx.CollectableRow) (domain.Receipt, error) {
	var r domain.Receipt
	err := row.Scan(&r.ID, &r.CompanyID, &r.RunID, &r.RunEmployeeID, &r.EmployeeID, &r.TemplateID,
		&r.ReceiptNumber, &r.CUIT, &r.EmployeeCUIL, &r.PeriodName, &r.PeriodStart, &r.PeriodEnd, &r.PaymentDate,
		&r.GrossRemunerative, &r.GrossNonRemunerative, &r.DeductionsTotal, &r.ContributionsTotal, &r.NetAmount, &r.EmployerCost,
		&r.Currency, &r.AmountInWords, &r.DigitalToken, &r.QRCode, &r.Barcode, &r.Status,
		&r.AcknowledgedAt, &r.AcknowledgedIP, &r.ViewedAt, &r.DownloadedAt, &r.EmailedAt,
		&r.StoragePath, &r.GeneratedBy, &r.GeneratedAt, &r.CreatedAt, &r.UpdatedAt)
	return r, err
}
