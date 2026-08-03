package employees

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/models"
)

type EmployeeRepository struct {
	pool *pgxpool.Pool
}

func NewEmployeeRepository(pool *pgxpool.Pool) *EmployeeRepository {
	return &EmployeeRepository{pool: pool}
}

func (r *EmployeeRepository) Create(ctx context.Context, emp *models.Employee) error {
	query := `
		INSERT INTO employees (
			company_id, employee_number, first_name, last_name, dni, email, phone,
			birth_date, photo_url, branch_id, department_id, position_id, manager_id,
			hire_date, status
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
		RETURNING id, created_at, updated_at`
	return r.pool.QueryRow(ctx, query,
		emp.CompanyID, emp.EmployeeNumber, emp.FirstName, emp.LastName,
		emp.DNI, emp.Email, emp.Phone, emp.BirthDate, emp.PhotoURL,
		emp.BranchID, emp.DepartmentID, emp.PositionID, emp.ManagerID,
		emp.HireDate, emp.Status,
	).Scan(&emp.ID, &emp.CreatedAt, &emp.UpdatedAt)
}

func (r *EmployeeRepository) FindByID(ctx context.Context, id, companyID string) (*models.Employee, error) {
	query := `
		SELECT
			e.id, e.company_id, e.employee_number, e.first_name, e.last_name,
			e.dni, e.email, e.phone, e.birth_date::text, e.photo_url,
			e.branch_id, e.department_id, e.position_id, e.manager_id,
			e.hire_date::text, e.termination_date::text, e.status, e.created_at, e.updated_at,
			b.name, d.name, p.name,
			CASE WHEN m.id IS NOT NULL THEN m.first_name || ' ' || m.last_name ELSE NULL END
		FROM employees e
		LEFT JOIN branches b ON b.id = e.branch_id
		LEFT JOIN departments d ON d.id = e.department_id
		LEFT JOIN positions p ON p.id = e.position_id
		LEFT JOIN employees m ON m.id = e.manager_id
		WHERE e.id = $1 AND e.company_id = $2`

	emp := &models.Employee{}
	err := r.pool.QueryRow(ctx, query, id, companyID).Scan(
		&emp.ID, &emp.CompanyID, &emp.EmployeeNumber, &emp.FirstName, &emp.LastName,
		&emp.DNI, &emp.Email, &emp.Phone, &emp.BirthDate, &emp.PhotoURL,
		&emp.BranchID, &emp.DepartmentID, &emp.PositionID, &emp.ManagerID,
		&emp.HireDate, &emp.TerminationDate, &emp.Status, &emp.CreatedAt, &emp.UpdatedAt,
		&emp.BranchName, &emp.DepartmentName, &emp.PositionName, &emp.ManagerName,
	)
	if err != nil {
		return nil, err
	}
	return emp, nil
}

func (r *EmployeeRepository) List(ctx context.Context, companyID string, params *models.PaginationParams, filters EmployeeFilters) ([]models.Employee, int64, error) {
	where := []string{`e.company_id = $1`}
	args := []interface{}{companyID}
	argIdx := 2

	if filters.Search != "" {
		where = append(where, fmt.Sprintf(`(e.first_name ILIKE $%d OR e.last_name ILIKE $%d OR e.employee_number ILIKE $%d OR e.dni ILIKE $%d)`, argIdx, argIdx, argIdx, argIdx))
		args = append(args, "%"+filters.Search+"%")
		argIdx++
	}
	if filters.Status != "" {
		where = append(where, fmt.Sprintf(`e.status = $%d`, argIdx))
		args = append(args, filters.Status)
		argIdx++
	}
	if filters.DepartmentID != "" {
		where = append(where, fmt.Sprintf(`e.department_id = $%d`, argIdx))
		args = append(args, filters.DepartmentID)
		argIdx++
	}
	if filters.BranchID != "" {
		where = append(where, fmt.Sprintf(`e.branch_id = $%d`, argIdx))
		args = append(args, filters.BranchID)
		argIdx++
	}
	if filters.PositionID != "" {
		where = append(where, fmt.Sprintf(`e.position_id = $%d`, argIdx))
		args = append(args, filters.PositionID)
		argIdx++
	}
	if filters.ManagerID != "" {
		where = append(where, fmt.Sprintf(`e.manager_id = $%d`, argIdx))
		args = append(args, filters.ManagerID)
		argIdx++
	}

	whereClause := strings.Join(where, " AND ")

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM employees e WHERE %s`, whereClause)
	var total int64
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	orderBy := "e.last_name ASC, e.first_name ASC"
	if filters.SortBy != "" {
		validSorts := map[string]string{
			"first_name": "e.first_name",
			"last_name":  "e.last_name",
			"hire_date":  "e.hire_date",
			"status":     "e.status",
			"created_at": "e.created_at",
		}
		if col, ok := validSorts[filters.SortBy]; ok {
			dir := "ASC"
			if filters.SortDir == "desc" {
				dir = "DESC"
			}
			orderBy = fmt.Sprintf(`%s %s`, col, dir)
		}
	}

	listQuery := fmt.Sprintf(`
		SELECT
			e.id, e.company_id, e.employee_number, e.first_name, e.last_name,
			e.dni, e.email, e.phone, e.birth_date::text, e.photo_url,
			e.branch_id, e.department_id, e.position_id, e.manager_id,
			e.hire_date::text, e.termination_date::text, e.status, e.created_at, e.updated_at,
			b.name, d.name, p.name,
			CASE WHEN m.id IS NOT NULL THEN m.first_name || ' ' || m.last_name ELSE NULL END
		FROM employees e
		LEFT JOIN branches b ON b.id = e.branch_id
		LEFT JOIN departments d ON d.id = e.department_id
		LEFT JOIN positions p ON p.id = e.position_id
		LEFT JOIN employees m ON m.id = e.manager_id
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d`,
		whereClause, orderBy, argIdx, argIdx+1)

	args = append(args, params.PerPage, params.Offset)

	rows, err := r.pool.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var employees []models.Employee
	for rows.Next() {
		var e models.Employee
		if err := rows.Scan(
			&e.ID, &e.CompanyID, &e.EmployeeNumber, &e.FirstName, &e.LastName,
			&e.DNI, &e.Email, &e.Phone, &e.BirthDate, &e.PhotoURL,
			&e.BranchID, &e.DepartmentID, &e.PositionID, &e.ManagerID,
			&e.HireDate, &e.TerminationDate, &e.Status, &e.CreatedAt, &e.UpdatedAt,
			&e.BranchName, &e.DepartmentName, &e.PositionName, &e.ManagerName,
		); err != nil {
			return nil, 0, err
		}
		employees = append(employees, e)
	}
	return employees, total, nil
}

func (r *EmployeeRepository) Update(ctx context.Context, emp *models.Employee) error {
	query := `
		UPDATE employees
		SET first_name=$1, last_name=$2, dni=$3, email=$4, phone=$5,
			birth_date=$6, photo_url=$7, branch_id=$8, department_id=$9,
			position_id=$10, manager_id=$11, hire_date=$12, termination_date=$13,
			status=$14, updated_at=NOW()
		WHERE id=$15 AND company_id=$16`
	_, err := r.pool.Exec(ctx, query,
		emp.FirstName, emp.LastName, emp.DNI, emp.Email, emp.Phone,
		emp.BirthDate, emp.PhotoURL, emp.BranchID, emp.DepartmentID,
		emp.PositionID, emp.ManagerID, emp.HireDate, emp.TerminationDate,
		emp.Status, emp.ID, emp.CompanyID,
	)
	return err
}

func (r *EmployeeRepository) Delete(ctx context.Context, id, companyID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM employees WHERE id=$1 AND company_id=$2`, id, companyID)
	return err
}

func (r *EmployeeRepository) GetContacts(ctx context.Context, employeeID string) ([]models.EmployeeContact, error) {
	query := `SELECT id, employee_id, contact_type, contact_value, is_primary FROM employee_contacts WHERE employee_id=$1 ORDER BY is_primary DESC`
	rows, err := r.pool.Query(ctx, query, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []models.EmployeeContact
	for rows.Next() {
		var c models.EmployeeContact
		if err := rows.Scan(&c.ID, &c.EmployeeID, &c.ContactType, &c.ContactValue, &c.IsPrimary); err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}
	return contacts, nil
}

func (r *EmployeeRepository) UpsertContacts(ctx context.Context, employeeID string, contacts []models.EmployeeContact) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM employee_contacts WHERE employee_id=$1`, employeeID)
	if err != nil {
		return err
	}

	for _, c := range contacts {
		_, err := r.pool.Exec(ctx,
			`INSERT INTO employee_contacts (id, employee_id, contact_type, contact_value, is_primary) VALUES (gen_random_uuid(),$1,$2,$3,$4)`,
			employeeID, c.ContactType, c.ContactValue, c.IsPrimary,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *EmployeeRepository) GetAddresses(ctx context.Context, employeeID string) ([]models.EmployeeAddress, error) {
	query := `SELECT id, employee_id, address_type, street, street_number, apartment, city, state, country, postal_code, is_primary FROM employee_addresses WHERE employee_id=$1 ORDER BY is_primary DESC`
	rows, err := r.pool.Query(ctx, query, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []models.EmployeeAddress
	for rows.Next() {
		var a models.EmployeeAddress
		if err := rows.Scan(&a.ID, &a.EmployeeID, &a.AddressType, &a.Street, &a.StreetNumber, &a.Apartment, &a.City, &a.State, &a.Country, &a.PostalCode, &a.IsPrimary); err != nil {
			return nil, err
		}
		addresses = append(addresses, a)
	}
	return addresses, nil
}

func (r *EmployeeRepository) UpsertAddresses(ctx context.Context, employeeID string, addresses []models.EmployeeAddress) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM employee_addresses WHERE employee_id=$1`, employeeID)
	if err != nil {
		return err
	}

	for _, a := range addresses {
		_, err := r.pool.Exec(ctx,
			`INSERT INTO employee_addresses (id, employee_id, address_type, street, street_number, apartment, city, state, country, postal_code, is_primary) VALUES (gen_random_uuid(),$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
			employeeID, a.AddressType, a.Street, a.StreetNumber, a.Apartment, a.City, a.State, a.Country, a.PostalCode, a.IsPrimary,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *EmployeeRepository) GetEmergencyContacts(ctx context.Context, employeeID string) ([]models.EmployeeEmergencyContact, error) {
	query := `SELECT id, employee_id, name, relationship, phone, alt_phone, is_primary FROM employee_emergency_contacts WHERE employee_id=$1 ORDER BY is_primary DESC`
	rows, err := r.pool.Query(ctx, query, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []models.EmployeeEmergencyContact
	for rows.Next() {
		var c models.EmployeeEmergencyContact
		if err := rows.Scan(&c.ID, &c.EmployeeID, &c.Name, &c.Relationship, &c.Phone, &c.AltPhone, &c.IsPrimary); err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}
	return contacts, nil
}

func (r *EmployeeRepository) UpsertEmergencyContacts(ctx context.Context, employeeID string, contacts []models.EmployeeEmergencyContact) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM employee_emergency_contacts WHERE employee_id=$1`, employeeID)
	if err != nil {
		return err
	}

	for _, c := range contacts {
		_, err := r.pool.Exec(ctx,
			`INSERT INTO employee_emergency_contacts (id, employee_id, name, relationship, phone, alt_phone, is_primary) VALUES (gen_random_uuid(),$1,$2,$3,$4,$5,$6)`,
			employeeID, c.Name, c.Relationship, c.Phone, c.AltPhone, c.IsPrimary,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *EmployeeRepository) AddHistory(ctx context.Context, h *models.EmployeeHistory) error {
	query := `
		INSERT INTO employee_history (id, employee_id, event_type, old_value, new_value, description, performed_by)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`
	return r.pool.QueryRow(ctx, query,
		h.EmployeeID, h.EventType, h.OldValue, h.NewValue, h.Description, h.PerformedBy,
	).Scan(&h.ID, &h.CreatedAt)
}

func (r *EmployeeRepository) GetHistory(ctx context.Context, employeeID string, limit int) ([]models.EmployeeHistory, error) {
	query := `SELECT id, employee_id, event_type, old_value, new_value, description, performed_by, created_at FROM employee_history WHERE employee_id=$1 ORDER BY created_at DESC LIMIT $2`
	rows, err := r.pool.Query(ctx, query, employeeID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.EmployeeHistory
	for rows.Next() {
		var h models.EmployeeHistory
		if err := rows.Scan(&h.ID, &h.EmployeeID, &h.EventType, &h.OldValue, &h.NewValue, &h.Description, &h.PerformedBy, &h.CreatedAt); err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, nil
}

type EmployeeFilters struct {
	Search       string
	Status       string
	DepartmentID string
	BranchID     string
	PositionID   string
	ManagerID    string
	SortBy       string
	SortDir      string
}
