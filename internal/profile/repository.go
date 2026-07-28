package profile

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EmployeeProfile struct {
	ID             string     `json:"id"`
	EmployeeNumber string     `json:"employee_number"`
	FirstName      string     `json:"first_name"`
	LastName       string     `json:"last_name"`
	DNI            *string    `json:"dni,omitempty"`
	Email          *string    `json:"email,omitempty"`
	Phone          *string    `json:"phone,omitempty"`
	BirthDate      *string    `json:"birth_date,omitempty"`
	PhotoURL       *string    `json:"photo_url,omitempty"`
	HireDate       string     `json:"hire_date"`
	Tenure         string     `json:"tenure"`
	Status         string     `json:"status"`

	BranchName     *string `json:"branch_name,omitempty"`
	DepartmentName *string `json:"department_name,omitempty"`
	PositionName   *string `json:"position_name,omitempty"`
	PositionLevel  *int    `json:"position_level,omitempty"`
	ManagerName    *string `json:"manager_name,omitempty"`

	PhoneContact     *string `json:"phone_contact,omitempty"`
	PersonalEmail    *string `json:"personal_email,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

type ProfileRepository struct {
	pool *pgxpool.Pool
}

func NewProfileRepository(pool *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{pool: pool}
}

func (r *ProfileRepository) GetByUserID(ctx context.Context, userID, companyID string) (*EmployeeProfile, error) {
	query := `
		SELECT
			e.id, e.employee_number, e.first_name, e.last_name,
			e.dni, e.email, e.phone, e.birth_date, e.photo_url,
			e.hire_date, e.status, e.created_at,
			b.name, d.name, p.name, p.level,
			CASE WHEN m.id IS NOT NULL THEN m.first_name || ' ' || m.last_name ELSE NULL END
		FROM employees e
		LEFT JOIN user_companies uc ON uc.company_id = e.company_id
		LEFT JOIN branches b ON b.id = e.branch_id
		LEFT JOIN departments d ON d.id = e.department_id
		LEFT JOIN positions p ON p.id = e.position_id
		LEFT JOIN employees m ON m.id = e.manager_id
		WHERE uc.user_id = $1 AND e.company_id = $2 AND e.status = 'active'
		LIMIT 1`

	profile := &EmployeeProfile{}
	err := r.pool.QueryRow(ctx, query, userID, companyID).Scan(
		&profile.ID, &profile.EmployeeNumber, &profile.FirstName, &profile.LastName,
		&profile.DNI, &profile.Email, &profile.Phone, &profile.BirthDate, &profile.PhotoURL,
		&profile.HireDate, &profile.Status, &profile.CreatedAt,
		&profile.BranchName, &profile.DepartmentName, &profile.PositionName, &profile.PositionLevel,
		&profile.ManagerName,
	)
	if err != nil {
		return nil, err
	}

	hireDate, err := time.Parse("2006-01-02", profile.HireDate)
	if err == nil {
		profile.Tenure = calculateTenure(hireDate)
	}

	return profile, nil
}

func (r *ProfileRepository) GetEmployeeIDByUser(ctx context.Context, userID, companyID string) (string, error) {
	query := `
		SELECT e.id
		FROM employees e
		WHERE e.company_id = $2 AND e.status = 'active'
		AND EXISTS (SELECT 1 FROM user_companies uc WHERE uc.user_id = $1 AND uc.company_id = e.company_id)
		LIMIT 1`

	var empID string
	err := r.pool.QueryRow(ctx, query, userID, companyID).Scan(&empID)
	return empID, err
}

func (r *ProfileRepository) UpdateAllowedFields(ctx context.Context, employeeID, companyID string, phone, photoURL *string) error {
	query := `
		UPDATE employees
		SET phone = COALESCE($1, phone),
			photo_url = COALESCE($2, photo_url),
			updated_at = NOW()
		WHERE id = $3 AND company_id = $4`
	_, err := r.pool.Exec(ctx, query, phone, photoURL, employeeID, companyID)
	return err
}

func calculateTenure(hireDate time.Time) string {
	now := time.Now()
	years := now.Year() - hireDate.Year()
	months := int(now.Month()) - int(hireDate.Month())
	days := now.Day() - hireDate.Day()

	if days < 0 {
		months--
	}
	if months < 0 {
		years--
		months += 12
	}

	if years > 0 {
		return formatTenure(years, months, "year", "month")
	}
	if months > 0 {
		return formatTenure(months, 0, "month", "")
	}

	d := int(now.Sub(hireDate).Hours() / 24)
	return formatTenure(d, 0, "day", "")
}

func formatTenure(a, b int, unitA, unitB string) string {
	if a != 1 {
		unitA = unitA + "s"
	}
	if b > 0 && unitB != "" {
		if b != 1 {
			unitB = unitB + "s"
		}
		return fmt.Sprintf("%d %s, %d %s", a, unitA, b, unitB)
	}
	return fmt.Sprintf("%d %s", a, unitA)
}
