package employees

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rrhhumand/api/internal/models"
)

type CycleChecker interface {
	HasCycle(ctx context.Context, employeeID, newManagerID string) (bool, error)
}

type EmployeeService struct {
	repo        *EmployeeRepository
	cycleCheck  CycleChecker
}

func NewEmployeeService(repo *EmployeeRepository, cycleCheck CycleChecker) *EmployeeService {
	return &EmployeeService{repo: repo, cycleCheck: cycleCheck}
}

type CreateEmployeeRequest struct {
	EmployeeNumber string  `json:"employee_number" validate:"required"`
	FirstName      string  `json:"first_name" validate:"required"`
	LastName       string  `json:"last_name" validate:"required"`
	DNI            *string `json:"dni,omitempty"`
	Email          *string `json:"email,omitempty"`
	Phone          *string `json:"phone,omitempty"`
	BirthDate      *string `json:"birth_date,omitempty"`
	PhotoURL       *string `json:"photo_url,omitempty"`
	BranchID       *string `json:"branch_id,omitempty"`
	DepartmentID   *string `json:"department_id,omitempty"`
	PositionID     *string `json:"position_id,omitempty"`
	ManagerID      *string `json:"manager_id,omitempty"`
	HireDate       string  `json:"hire_date" validate:"required"`
}

type UpdateEmployeeRequest struct {
	FirstName      *string `json:"first_name,omitempty"`
	LastName       *string `json:"last_name,omitempty"`
	DNI            *string `json:"dni,omitempty"`
	Email          *string `json:"email,omitempty"`
	Phone          *string `json:"phone,omitempty"`
	BirthDate      *string `json:"birth_date,omitempty"`
	PhotoURL       *string `json:"photo_url,omitempty"`
	BranchID       *string `json:"branch_id,omitempty"`
	DepartmentID   *string `json:"department_id,omitempty"`
	PositionID     *string `json:"position_id,omitempty"`
	ManagerID      *string `json:"manager_id,omitempty"`
	HireDate       *string `json:"hire_date,omitempty"`
	TerminationDate *string `json:"termination_date,omitempty"`
	Status         *string `json:"status,omitempty"`
}

type EmployeeListResponse struct {
	Employees []models.Employee    `json:"employees"`
	Total     int64                `json:"total"`
	Page      int                  `json:"page"`
	PerPage   int                  `json:"per_page"`
}

func (s *EmployeeService) Create(ctx context.Context, companyID string, req *CreateEmployeeRequest) (*models.Employee, error) {
	emp := &models.Employee{
		CompanyID:      companyID,
		EmployeeNumber: req.EmployeeNumber,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		DNI:            req.DNI,
		Email:          req.Email,
		Phone:          req.Phone,
		BirthDate:      req.BirthDate,
		PhotoURL:       req.PhotoURL,
		BranchID:       req.BranchID,
		DepartmentID:   req.DepartmentID,
		PositionID:     req.PositionID,
		ManagerID:      req.ManagerID,
		HireDate:       req.HireDate,
		Status:         "active",
	}

	if err := s.repo.Create(ctx, emp); err != nil {
		return nil, err
	}

	_ = s.repo.AddHistory(ctx, &models.EmployeeHistory{
		EmployeeID:  emp.ID,
		EventType:   "HIRED",
		Description: strPtr(fmt.Sprintf("Employee %s %s hired", emp.FirstName, emp.LastName)),
	})

	return s.repo.FindByID(ctx, emp.ID, companyID)
}

func (s *EmployeeService) GetByID(ctx context.Context, id, companyID string) (*models.Employee, error) {
	emp, err := s.repo.FindByID(ctx, id, companyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("employee not found")
		}
		return nil, err
	}
	return emp, nil
}

func (s *EmployeeService) List(ctx context.Context, companyID string, params *models.PaginationParams, filters EmployeeFilters) ([]models.Employee, int64, error) {
	return s.repo.List(ctx, companyID, params, filters)
}

func (s *EmployeeService) Update(ctx context.Context, id, companyID string, req *UpdateEmployeeRequest) (*models.Employee, error) {
	emp, err := s.repo.FindByID(ctx, id, companyID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	if req.FirstName != nil {
		emp.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		emp.LastName = *req.LastName
	}
	if req.DNI != nil {
		emp.DNI = req.DNI
	}
	if req.Email != nil {
		emp.Email = req.Email
	}
	if req.Phone != nil {
		emp.Phone = req.Phone
	}
	if req.BirthDate != nil {
		emp.BirthDate = req.BirthDate
	}
	if req.PhotoURL != nil {
		emp.PhotoURL = req.PhotoURL
	}
	if req.BranchID != nil {
		emp.BranchID = req.BranchID
	}
	if req.DepartmentID != nil {
		emp.DepartmentID = req.DepartmentID
	}
	if req.PositionID != nil {
		emp.PositionID = req.PositionID
	}
	if req.ManagerID != nil {
		if s.cycleCheck != nil {
			hasCycle, err := s.cycleCheck.HasCycle(ctx, emp.ID, *req.ManagerID)
			if err != nil {
				return nil, fmt.Errorf("failed to check hierarchy cycle: %w", err)
			}
			if hasCycle {
				return nil, errors.New("assigning this manager would create a hierarchy cycle")
			}
		}
		emp.ManagerID = req.ManagerID
	}
	if req.HireDate != nil {
		emp.HireDate = *req.HireDate
	}
	if req.TerminationDate != nil {
		emp.TerminationDate = req.TerminationDate
	}
	if req.Status != nil {
		oldStatus := emp.Status
		emp.Status = *req.Status
		if oldStatus != *req.Status {
			_ = s.repo.AddHistory(ctx, &models.EmployeeHistory{
				EmployeeID:  emp.ID,
				EventType:   "STATUS_CHANGE",
				OldValue:    &oldStatus,
				NewValue:    req.Status,
				Description: strPtr(fmt.Sprintf("Status changed from %s to %s", oldStatus, *req.Status)),
			})
		}
	}

	if err := s.repo.Update(ctx, emp); err != nil {
		return nil, err
	}

	_ = s.repo.AddHistory(ctx, &models.EmployeeHistory{
		EmployeeID:  emp.ID,
		EventType:   "UPDATED",
		Description: strPtr("Employee record updated"),
	})

	return s.repo.FindByID(ctx, emp.ID, companyID)
}

func (s *EmployeeService) Delete(ctx context.Context, id, companyID string) error {
	emp, err := s.repo.FindByID(ctx, id, companyID)
	if err != nil {
		return errors.New("employee not found")
	}

	_ = s.repo.AddHistory(ctx, &models.EmployeeHistory{
		EmployeeID:  emp.ID,
		EventType:   "TERMINATED",
		Description: strPtr(fmt.Sprintf("Employee %s %s terminated", emp.FirstName, emp.LastName)),
	})

	return s.repo.Delete(ctx, id, companyID)
}

func (s *EmployeeService) GetContacts(ctx context.Context, employeeID string) ([]models.EmployeeContact, error) {
	return s.repo.GetContacts(ctx, employeeID)
}

func (s *EmployeeService) UpsertContacts(ctx context.Context, employeeID string, contacts []models.EmployeeContact) error {
	return s.repo.UpsertContacts(ctx, employeeID, contacts)
}

func (s *EmployeeService) GetAddresses(ctx context.Context, employeeID string) ([]models.EmployeeAddress, error) {
	return s.repo.GetAddresses(ctx, employeeID)
}

func (s *EmployeeService) UpsertAddresses(ctx context.Context, employeeID string, addresses []models.EmployeeAddress) error {
	return s.repo.UpsertAddresses(ctx, employeeID, addresses)
}

func (s *EmployeeService) GetEmergencyContacts(ctx context.Context, employeeID string) ([]models.EmployeeEmergencyContact, error) {
	return s.repo.GetEmergencyContacts(ctx, employeeID)
}

func (s *EmployeeService) UpsertEmergencyContacts(ctx context.Context, employeeID string, contacts []models.EmployeeEmergencyContact) error {
	return s.repo.UpsertEmergencyContacts(ctx, employeeID, contacts)
}

func (s *EmployeeService) GetHistory(ctx context.Context, employeeID string) ([]models.EmployeeHistory, error) {
	return s.repo.GetHistory(ctx, employeeID, 50)
}

func strPtr(s string) *string {
	return &s
}
