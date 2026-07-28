package leave

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rrhhumand/api/internal/models"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateLeaveType(ctx context.Context, companyID string, req *CreateLeaveTypeRequest) (*models.LeaveType, error) {
	lt := &models.LeaveType{}
	requiresApproval := true
	requiresDocument := false
	affectsBalance := true
	isPaid := true
	if req.RequiresApproval != nil {
		requiresApproval = *req.RequiresApproval
	}
	if req.RequiresDocument != nil {
		requiresDocument = *req.RequiresDocument
	}
	if req.AffectsBalance != nil {
		affectsBalance = *req.AffectsBalance
	}
	if req.IsPaid != nil {
		isPaid = *req.IsPaid
	}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO leave_types (company_id, name, code, description, category, requires_approval, requires_document, affects_balance, is_paid)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id, company_id, name, code, description, category, requires_approval, requires_document, affects_balance, is_paid, is_active, created_at, updated_at`,
		companyID, req.Name, req.Code, req.Description, req.Category, requiresApproval, requiresDocument, affectsBalance, isPaid,
	).Scan(&lt.ID, &lt.CompanyID, &lt.Name, &lt.Code, &lt.Description, &lt.Category, &lt.RequiresApproval, &lt.RequiresDocument, &lt.AffectsBalance, &lt.IsPaid, &lt.IsActive, &lt.CreatedAt, &lt.UpdatedAt)
	return lt, err
}

func (r *Repository) GetLeaveType(ctx context.Context, companyID, id string) (*models.LeaveType, error) {
	lt := &models.LeaveType{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, name, code, description, category, requires_approval, requires_document, affects_balance, is_paid, is_active, created_at, updated_at
		 FROM leave_types WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&lt.ID, &lt.CompanyID, &lt.Name, &lt.Code, &lt.Description, &lt.Category, &lt.RequiresApproval, &lt.RequiresDocument, &lt.AffectsBalance, &lt.IsPaid, &lt.IsActive, &lt.CreatedAt, &lt.UpdatedAt)
	return lt, err
}

func (r *Repository) ListLeaveTypes(ctx context.Context, companyID string) ([]models.LeaveType, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, name, code, description, category, requires_approval, requires_document, affects_balance, is_paid, is_active, created_at, updated_at
		 FROM leave_types WHERE company_id=$1 ORDER BY name`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []models.LeaveType
	for rows.Next() {
		var lt models.LeaveType
		if err := rows.Scan(&lt.ID, &lt.CompanyID, &lt.Name, &lt.Code, &lt.Description, &lt.Category, &lt.RequiresApproval, &lt.RequiresDocument, &lt.AffectsBalance, &lt.IsPaid, &lt.IsActive, &lt.CreatedAt, &lt.UpdatedAt); err != nil {
			return nil, err
		}
		types = append(types, lt)
	}
	return types, nil
}

func (r *Repository) UpdateLeaveType(ctx context.Context, companyID, id string, req *UpdateLeaveTypeRequest) (*models.LeaveType, error) {
	lt := &models.LeaveType{}
	err := r.pool.QueryRow(ctx,
		`UPDATE leave_types SET
		 name=COALESCE($3,name), description=COALESCE($4,description), category=COALESCE($5,category),
		 requires_approval=COALESCE($6,requires_approval), requires_document=COALESCE($7,requires_document),
		 affects_balance=COALESCE($8,affects_balance), is_paid=COALESCE($9,is_paid), is_active=COALESCE($10,is_active), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, name, code, description, category, requires_approval, requires_document, affects_balance, is_paid, is_active, created_at, updated_at`,
		companyID, id, req.Name, req.Description, req.Category, req.RequiresApproval, req.RequiresDocument, req.AffectsBalance, req.IsPaid, req.IsActive,
	).Scan(&lt.ID, &lt.CompanyID, &lt.Name, &lt.Code, &lt.Description, &lt.Category, &lt.RequiresApproval, &lt.RequiresDocument, &lt.AffectsBalance, &lt.IsPaid, &lt.IsActive, &lt.CreatedAt, &lt.UpdatedAt)
	return lt, err
}

func (r *Repository) DeleteLeaveType(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM leave_types WHERE company_id=$1 AND id=$2`, companyID, id)
	return err
}

func (r *Repository) CreateLeavePolicy(ctx context.Context, companyID string, req *CreateLeavePolicyRequest) (*models.LeavePolicy, error) {
	p := &models.LeavePolicy{}
	allowNegative := false
	useBusiness := true
	requireManager := true
	requireHR := false
	minDays := 0
	if req.AllowNegativeBalance != nil {
		allowNegative = *req.AllowNegativeBalance
	}
	if req.UseBusinessDays != nil {
		useBusiness = *req.UseBusinessDays
	}
	if req.RequiresManagerApproval != nil {
		requireManager = *req.RequiresManagerApproval
	}
	if req.RequiresHRApproval != nil {
		requireHR = *req.RequiresHRApproval
	}
	if req.MinimumDaysBeforeRequest != nil {
		minDays = *req.MinimumDaysBeforeRequest
	}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO leave_policies (company_id, leave_type_id, name, days_per_year, minimum_days_before_request, maximum_days_per_request, maximum_accumulated_days, allow_negative_balance, use_business_days, requires_manager_approval, requires_hr_approval)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		 RETURNING id, company_id, leave_type_id, name, days_per_year, minimum_days_before_request, maximum_days_per_request, maximum_accumulated_days, allow_negative_balance, use_business_days, requires_manager_approval, requires_hr_approval, is_active, created_at, updated_at`,
		companyID, req.LeaveTypeID, req.Name, req.DaysPerYear, minDays, req.MaximumDaysPerRequest, req.MaximumAccumulatedDays, allowNegative, useBusiness, requireManager, requireHR,
	).Scan(&p.ID, &p.CompanyID, &p.LeaveTypeID, &p.Name, &p.DaysPerYear, &p.MinimumDaysBeforeRequest, &p.MaximumDaysPerRequest, &p.MaximumAccumulatedDays, &p.AllowNegativeBalance, &p.UseBusinessDays, &p.RequiresManagerApproval, &p.RequiresHRApproval, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repository) GetLeavePolicy(ctx context.Context, companyID, leaveTypeID string) (*models.LeavePolicy, error) {
	p := &models.LeavePolicy{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, leave_type_id, name, days_per_year, minimum_days_before_request, maximum_days_per_request, maximum_accumulated_days, allow_negative_balance, use_business_days, requires_manager_approval, requires_hr_approval, is_active, created_at, updated_at
		 FROM leave_policies WHERE company_id=$1 AND leave_type_id=$2 AND is_active=true LIMIT 1`, companyID, leaveTypeID,
	).Scan(&p.ID, &p.CompanyID, &p.LeaveTypeID, &p.Name, &p.DaysPerYear, &p.MinimumDaysBeforeRequest, &p.MaximumDaysPerRequest, &p.MaximumAccumulatedDays, &p.AllowNegativeBalance, &p.UseBusinessDays, &p.RequiresManagerApproval, &p.RequiresHRApproval, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repository) ListLeavePolicies(ctx context.Context, companyID string) ([]models.LeavePolicy, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT p.id, p.company_id, p.leave_type_id, lt.name, p.name, p.days_per_year, p.minimum_days_before_request, p.maximum_days_per_request, p.maximum_accumulated_days, p.allow_negative_balance, p.use_business_days, p.requires_manager_approval, p.requires_hr_approval, p.is_active, p.created_at, p.updated_at
		 FROM leave_policies p JOIN leave_types lt ON p.leave_type_id=lt.id
		 WHERE p.company_id=$1 ORDER BY p.name`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var policies []models.LeavePolicy
	for rows.Next() {
		var p models.LeavePolicy
		if err := rows.Scan(&p.ID, &p.CompanyID, &p.LeaveTypeID, &p.LeaveTypeName, &p.Name, &p.DaysPerYear, &p.MinimumDaysBeforeRequest, &p.MaximumDaysPerRequest, &p.MaximumAccumulatedDays, &p.AllowNegativeBalance, &p.UseBusinessDays, &p.RequiresManagerApproval, &p.RequiresHRApproval, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		policies = append(policies, p)
	}
	return policies, nil
}

func (r *Repository) CreateHoliday(ctx context.Context, companyID string, req *CreateHolidayRequest) (*models.Holiday, error) {
	date, _ := time.Parse("2006-01-02", req.Date)
	isRecurring := false
	if req.IsRecurring != nil {
		isRecurring = *req.IsRecurring
	}
	h := &models.Holiday{}
	cid := &companyID
	err := r.pool.QueryRow(ctx,
		`INSERT INTO holidays (company_id, date, name, is_recurring) VALUES ($1,$2,$3,$4) RETURNING id, company_id, date, name, is_recurring, created_at`,
		cid, date, req.Name, isRecurring,
	).Scan(&h.ID, &h.CompanyID, &h.Date, &h.Name, &h.IsRecurring, &h.CreatedAt)
	return h, err
}

func (r *Repository) GetHolidays(ctx context.Context, companyID string, from, to time.Time) ([]models.Holiday, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, date, name, is_recurring, created_at
		 FROM holidays
		 WHERE (company_id=$1 OR company_id IS NULL) AND date>=$2 AND date<=$3
		 ORDER BY date`, companyID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var holidays []models.Holiday
	for rows.Next() {
		var h models.Holiday
		if err := rows.Scan(&h.ID, &h.CompanyID, &h.Date, &h.Name, &h.IsRecurring, &h.CreatedAt); err != nil {
			return nil, err
		}
		holidays = append(holidays, h)
	}
	return holidays, nil
}

func (r *Repository) DeleteHoliday(ctx context.Context, companyID, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM holidays WHERE id=$1 AND (company_id=$2 OR company_id IS NULL)`, id, companyID)
	return err
}

func (r *Repository) GetBalanceForUpdate(ctx context.Context, tx pgx.Tx, companyID, employeeID, leaveTypeID string, year int) (*models.LeaveBalance, error) {
	b := &models.LeaveBalance{}
	err := tx.QueryRow(ctx,
		`SELECT id, company_id, employee_id, leave_type_id, year, allocated_days, carried_over_days, adjustment_days, used_days, reserved_days, created_at, updated_at
		 FROM leave_balances WHERE employee_id=$1 AND leave_type_id=$2 AND year=$3 FOR UPDATE`,
		employeeID, leaveTypeID, year,
	).Scan(&b.ID, &b.CompanyID, &b.EmployeeID, &b.LeaveTypeID, &b.Year, &b.AllocatedDays, &b.CarriedOverDays, &b.AdjustmentDays, &b.UsedDays, &b.ReservedDays, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}
	b.AvailableDays = b.AllocatedDays + b.CarriedOverDays + b.AdjustmentDays - b.UsedDays - b.ReservedDays
	return b, nil
}

func (r *Repository) GetOrCreateBalance(ctx context.Context, tx pgx.Tx, companyID, employeeID, leaveTypeID string, year int) (*models.LeaveBalance, error) {
	b, err := r.GetBalanceForUpdate(ctx, tx, companyID, employeeID, leaveTypeID, year)
	if err == nil {
		return b, nil
	}
	// Create new balance
	b = &models.LeaveBalance{
		ID:         uuid.New().String(),
		CompanyID:  companyID,
		EmployeeID: employeeID,
		LeaveTypeID: leaveTypeID,
		Year:       year,
	}
	err = tx.QueryRow(ctx,
		`INSERT INTO leave_balances (id, company_id, employee_id, leave_type_id, year) VALUES ($1,$2,$3,$4,$5)
		 ON CONFLICT (employee_id, leave_type_id, year) DO UPDATE SET updated_at=NOW()
		 RETURNING id, company_id, employee_id, leave_type_id, year, allocated_days, carried_over_days, adjustment_days, used_days, reserved_days, created_at, updated_at`,
		b.ID, b.CompanyID, b.EmployeeID, b.LeaveTypeID, b.Year,
	).Scan(&b.ID, &b.CompanyID, &b.EmployeeID, &b.LeaveTypeID, &b.Year, &b.AllocatedDays, &b.CarriedOverDays, &b.AdjustmentDays, &b.UsedDays, &b.ReservedDays, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}
	b.AvailableDays = b.AllocatedDays + b.CarriedOverDays + b.AdjustmentDays - b.UsedDays - b.ReservedDays
	return b, nil
}

func (r *Repository) ListBalances(ctx context.Context, companyID, employeeID string, year int) ([]models.LeaveBalance, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT b.id, b.company_id, b.employee_id, b.leave_type_id, lt.name, b.year, b.allocated_days, b.carried_over_days, b.adjustment_days, b.used_days, b.reserved_days, b.created_at, b.updated_at
		 FROM leave_balances b JOIN leave_types lt ON b.leave_type_id=lt.id
		 WHERE b.company_id=$1 AND b.employee_id=$2 AND b.year=$3 ORDER BY lt.name`, companyID, employeeID, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []models.LeaveBalance
	for rows.Next() {
		var b models.LeaveBalance
		if err := rows.Scan(&b.ID, &b.CompanyID, &b.EmployeeID, &b.LeaveTypeID, &b.LeaveTypeName, &b.Year, &b.AllocatedDays, &b.CarriedOverDays, &b.AdjustmentDays, &b.UsedDays, &b.ReservedDays, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		b.AvailableDays = b.AllocatedDays + b.CarriedOverDays + b.AdjustmentDays - b.UsedDays - b.ReservedDays
		balances = append(balances, b)
	}
	return balances, nil
}

func (r *Repository) AdjustBalance(ctx context.Context, tx pgx.Tx, companyID, employeeID, leaveTypeID string, year int, adjustmentDays float64, reason string, performedBy string) error {
	b, err := r.GetOrCreateBalance(ctx, tx, companyID, employeeID, leaveTypeID, year)
	if err != nil {
		return err
	}
	newAdjustment := b.AdjustmentDays + adjustmentDays
	_, err = tx.Exec(ctx,
		`UPDATE leave_balances SET adjustment_days=$1, updated_at=NOW() WHERE id=$2`, newAdjustment, b.ID)
	return err
}

func (r *Repository) CreateLeaveRequest(ctx context.Context, req *models.LeaveRequest) (*models.LeaveRequest, error) {
	lr := &models.LeaveRequest{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO leave_requests (id, company_id, employee_id, leave_type_id, start_date, end_date, requested_days, reason, status, document_id)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		 RETURNING id, company_id, employee_id, leave_type_id, start_date, end_date, requested_days, reason, status, document_id, created_at, updated_at`,
		req.ID, req.CompanyID, req.EmployeeID, req.LeaveTypeID, req.StartDate, req.EndDate, req.RequestedDays, req.Reason, req.Status, req.DocumentID,
	).Scan(&lr.ID, &lr.CompanyID, &lr.EmployeeID, &lr.LeaveTypeID, &lr.StartDate, &lr.EndDate, &lr.RequestedDays, &lr.Reason, &lr.Status, &lr.DocumentID, &lr.CreatedAt, &lr.UpdatedAt)
	return lr, err
}

func (r *Repository) GetLeaveRequest(ctx context.Context, companyID, id string) (*models.LeaveRequest, error) {
	lr := &models.LeaveRequest{}
	err := r.pool.QueryRow(ctx,
		`SELECT lr.id, lr.company_id, lr.employee_id, e.first_name || ' ' || e.last_name, lr.leave_type_id, lt.name, lr.start_date, lr.end_date, lr.requested_days, lr.reason, lr.status, lr.document_id, lr.created_at, lr.updated_at
		 FROM leave_requests lr
		 JOIN employees e ON lr.employee_id=e.id
		 JOIN leave_types lt ON lr.leave_type_id=lt.id
		 WHERE lr.company_id=$1 AND lr.id=$2`, companyID, id,
	).Scan(&lr.ID, &lr.CompanyID, &lr.EmployeeID, &lr.EmployeeName, &lr.LeaveTypeID, &lr.LeaveTypeName, &lr.StartDate, &lr.EndDate, &lr.RequestedDays, &lr.Reason, &lr.Status, &lr.DocumentID, &lr.CreatedAt, &lr.UpdatedAt)
	return lr, err
}

func (r *Repository) ListLeaveRequests(ctx context.Context, companyID string, filters LeaveFilters, offset, limit int) ([]models.LeaveRequest, int64, error) {
	query := `SELECT lr.id, lr.company_id, lr.employee_id, e.first_name || ' ' || e.last_name, lr.leave_type_id, lt.name, lr.start_date, lr.end_date, lr.requested_days, lr.reason, lr.status, lr.document_id, lr.created_at, lr.updated_at
		 FROM leave_requests lr
		 JOIN employees e ON lr.employee_id=e.id
		 JOIN leave_types lt ON lr.leave_type_id=lt.id
		 WHERE lr.company_id=$1`
	countQuery := `SELECT COUNT(*) FROM leave_requests lr WHERE lr.company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if filters.EmployeeID != "" {
		query += fmt.Sprintf(" AND lr.employee_id=$%d", argIdx)
		countQuery += fmt.Sprintf(" AND lr.employee_id=$%d", argIdx)
		args = append(args, filters.EmployeeID)
		argIdx++
	}
	if filters.LeaveTypeID != "" {
		query += fmt.Sprintf(" AND lr.leave_type_id=$%d", argIdx)
		countQuery += fmt.Sprintf(" AND lr.leave_type_id=$%d", argIdx)
		args = append(args, filters.LeaveTypeID)
		argIdx++
	}
	if filters.Status != "" {
		query += fmt.Sprintf(" AND lr.status=$%d", argIdx)
		countQuery += fmt.Sprintf(" AND lr.status=$%d", argIdx)
		args = append(args, filters.Status)
		argIdx++
	}
	if filters.DateFrom != "" {
		query += fmt.Sprintf(" AND lr.start_date>=$%d", argIdx)
		countQuery += fmt.Sprintf(" AND lr.start_date>=$%d", argIdx)
		args = append(args, filters.DateFrom)
		argIdx++
	}
	if filters.DateTo != "" {
		query += fmt.Sprintf(" AND lr.end_date<=$%d", argIdx)
		countQuery += fmt.Sprintf(" AND lr.end_date<=$%d", argIdx)
		args = append(args, filters.DateTo)
		argIdx++
	}
	if filters.DepartmentID != "" {
		query += fmt.Sprintf(" AND e.department_id=$%d", argIdx)
		countQuery += fmt.Sprintf(" AND e.department_id=$%d", argIdx)
		args = append(args, filters.DepartmentID)
		argIdx++
	}

	var total int64
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query += fmt.Sprintf(" ORDER BY lr.created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var requests []models.LeaveRequest
	for rows.Next() {
		var lr models.LeaveRequest
		if err := rows.Scan(&lr.ID, &lr.CompanyID, &lr.EmployeeID, &lr.EmployeeName, &lr.LeaveTypeID, &lr.LeaveTypeName, &lr.StartDate, &lr.EndDate, &lr.RequestedDays, &lr.Reason, &lr.Status, &lr.DocumentID, &lr.CreatedAt, &lr.UpdatedAt); err != nil {
			return nil, 0, err
		}
		requests = append(requests, lr)
	}
	return requests, total, nil
}

func (r *Repository) UpdateLeaveRequestStatus(ctx context.Context, tx pgx.Tx, id, status string) error {
	_, err := tx.Exec(ctx, `UPDATE leave_requests SET status=$1, updated_at=NOW() WHERE id=$2`, status, id)
	return err
}

func (r *Repository) CheckOverlap(ctx context.Context, employeeID string, start, end time.Time, excludeID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM leave_requests WHERE employee_id=$1 AND status IN ('PENDING','APPROVED') AND start_date<=$2 AND end_date>=$3)`
	args := []interface{}{employeeID, end, start}
	if excludeID != "" {
		query = `SELECT EXISTS(SELECT 1 FROM leave_requests WHERE employee_id=$1 AND status IN ('PENDING','APPROVED') AND start_date<=$2 AND end_date>=$3 AND id!=$4)`
		args = append(args, excludeID)
	}
	var exists bool
	err := r.pool.QueryRow(ctx, query, args...).Scan(&exists)
	return exists, err
}

func (r *Repository) CreateApproval(ctx context.Context, approval *models.LeaveApproval) (*models.LeaveApproval, error) {
	a := &models.LeaveApproval{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO leave_approvals (id, company_id, leave_request_id, approver_id, level, status, comments, decided_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 RETURNING id, company_id, leave_request_id, approver_id, level, status, comments, decided_at, created_at`,
		approval.ID, approval.CompanyID, approval.LeaveRequestID, approval.ApproverID, approval.Level, approval.Status, approval.Comments, approval.DecidedAt,
	).Scan(&a.ID, &a.CompanyID, &a.LeaveRequestID, &a.ApproverID, &a.Level, &a.Status, &a.Comments, &a.DecidedAt, &a.CreatedAt)
	return a, err
}

func (r *Repository) UpdateApprovalStatus(ctx context.Context, tx pgx.Tx, id, status string, comments *string) error {
	_, err := tx.Exec(ctx,
		`UPDATE leave_approvals SET status=$1, comments=$2, decided_at=NOW() WHERE id=$3`, status, comments, id)
	return err
}

func (r *Repository) GetPendingApprovalsForRequest(ctx context.Context, leaveRequestID string) ([]models.LeaveApproval, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, leave_request_id, approver_id, level, status, comments, decided_at, created_at
		 FROM leave_approvals WHERE leave_request_id=$1 ORDER BY level`, leaveRequestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var approvals []models.LeaveApproval
	for rows.Next() {
		var a models.LeaveApproval
		if err := rows.Scan(&a.ID, &a.CompanyID, &a.LeaveRequestID, &a.ApproverID, &a.Level, &a.Status, &a.Comments, &a.DecidedAt, &a.CreatedAt); err != nil {
			return nil, err
		}
		approvals = append(approvals, a)
	}
	return approvals, nil
}

func (r *Repository) GetNextApprovalLevel(ctx context.Context, leaveRequestID string) (int, error) {
	var level int
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(MAX(level),0)+1 FROM leave_approvals WHERE leave_request_id=$1`, leaveRequestID).Scan(&level)
	return level, err
}

func (r *Repository) CreateHistory(ctx context.Context, h *models.LeaveRequestHistory) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO leave_request_history (id, leave_request_id, action, old_status, new_status, performed_by, comments)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		h.ID, h.LeaveRequestID, h.Action, h.OldStatus, h.NewStatus, h.PerformedBy, h.Comments)
	return err
}

func (r *Repository) GetHistory(ctx context.Context, leaveRequestID string) ([]models.LeaveRequestHistory, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT h.id, h.leave_request_id, h.action, h.old_status, h.new_status, h.performed_by, u.first_name || ' ' || u.last_name, h.comments, h.created_at
		 FROM leave_request_history h JOIN users u ON h.performed_by=u.id
		 WHERE h.leave_request_id=$1 ORDER BY h.created_at`, leaveRequestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.LeaveRequestHistory
	for rows.Next() {
		var h models.LeaveRequestHistory
		if err := rows.Scan(&h.ID, &h.LeaveRequestID, &h.Action, &h.OldStatus, &h.NewStatus, &h.PerformedBy, &h.PerformedByName, &h.Comments, &h.CreatedAt); err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, nil
}

func (r *Repository) GetExpiringRequests(ctx context.Context, days int) ([]models.LeaveRequest, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT lr.id, lr.company_id, lr.employee_id, lr.leave_type_id, lr.start_date, lr.end_date, lr.requested_days, lr.reason, lr.status, lr.document_id, lr.created_at, lr.updated_at
		 FROM leave_requests lr
		 WHERE lr.status='APPROVED' AND lr.start_date <= CURRENT_DATE + ($1 || ' days')::INTERVAL
		 ORDER BY lr.start_date`, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.LeaveRequest
	for rows.Next() {
		var lr models.LeaveRequest
		if err := rows.Scan(&lr.ID, &lr.CompanyID, &lr.EmployeeID, &lr.LeaveTypeID, &lr.StartDate, &lr.EndDate, &lr.RequestedDays, &lr.Reason, &lr.Status, &lr.DocumentID, &lr.CreatedAt, &lr.UpdatedAt); err != nil {
			return nil, err
		}
		requests = append(requests, lr)
	}
	return requests, nil
}

func (r *Repository) GetTeamAbsences(ctx context.Context, companyID, managerID string, start, end time.Time) ([]models.LeaveRequest, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT lr.id, lr.company_id, lr.employee_id, e.first_name || ' ' || e.last_name, lr.leave_type_id, lt.name, lr.start_date, lr.end_date, lr.requested_days, lr.reason, lr.status, lr.document_id, lr.created_at, lr.updated_at
		 FROM leave_requests lr
		 JOIN employees e ON lr.employee_id=e.id
		 JOIN leave_types lt ON lr.leave_type_id=lt.id
		 WHERE lr.company_id=$1 AND e.manager_id=$2 AND lr.status IN ('PENDING','APPROVED') AND lr.start_date<=$3 AND lr.end_date>=$4
		 ORDER BY lr.start_date`, companyID, managerID, end, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.LeaveRequest
	for rows.Next() {
		var lr models.LeaveRequest
		if err := rows.Scan(&lr.ID, &lr.CompanyID, &lr.EmployeeID, &lr.EmployeeName, &lr.LeaveTypeID, &lr.LeaveTypeName, &lr.StartDate, &lr.EndDate, &lr.RequestedDays, &lr.Reason, &lr.Status, &lr.DocumentID, &lr.CreatedAt, &lr.UpdatedAt); err != nil {
			return nil, err
		}
		requests = append(requests, lr)
	}
	return requests, nil
}

func (r *Repository) GetDepartmentAbsences(ctx context.Context, companyID, departmentID string, start, end time.Time) ([]models.LeaveRequest, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT lr.id, lr.company_id, lr.employee_id, e.first_name || ' ' || e.last_name, lr.leave_type_id, lt.name, lr.start_date, lr.end_date, lr.requested_days, lr.reason, lr.status, lr.document_id, lr.created_at, lr.updated_at
		 FROM leave_requests lr
		 JOIN employees e ON lr.employee_id=e.id
		 JOIN leave_types lt ON lr.leave_type_id=lt.id
		 WHERE lr.company_id=$1 AND e.department_id=$2 AND lr.status IN ('PENDING','APPROVED') AND lr.start_date<=$3 AND lr.end_date>=$4
		 ORDER BY lr.start_date`, companyID, departmentID, end, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.LeaveRequest
	for rows.Next() {
		var lr models.LeaveRequest
		if err := rows.Scan(&lr.ID, &lr.CompanyID, &lr.EmployeeID, &lr.EmployeeName, &lr.LeaveTypeID, &lr.LeaveTypeName, &lr.StartDate, &lr.EndDate, &lr.RequestedDays, &lr.Reason, &lr.Status, &lr.DocumentID, &lr.CreatedAt, &lr.UpdatedAt); err != nil {
			return nil, err
		}
		requests = append(requests, lr)
	}
	return requests, nil
}
