package overtime

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type RequestManager struct {
	repo *Repository
}

func NewRequestManager(repo *Repository) *RequestManager {
	return &RequestManager{repo: repo}
}

func (rm *RequestManager) CreateRequest(ctx context.Context, companyID, employeeID string, req *RequestOvertimeRequest) (*OvertimeRequest, error) {
	workDate, err := time.Parse("2006-01-02", req.WorkDate)
	if err != nil {
		return nil, fmt.Errorf("invalid date format")
	}

	if req.RequestedMinutes <= 0 {
		return nil, fmt.Errorf("requested_minutes must be positive")
	}

	overtimeReq := &OvertimeRequest{
		ID:               uuid.New().String(),
		CompanyID:        companyID,
		EmployeeID:       employeeID,
		OvertimeRecordID: req.OvertimeRecordID,
		WorkDate:         workDate,
		RequestedMinutes: req.RequestedMinutes,
		Reason:           req.Reason,
		Status:           "PENDING",
		RequestedAt:      time.Now(),
	}

	if err := rm.repo.CreateOvertimeRequest(ctx, overtimeReq); err != nil {
		return nil, err
	}

	return overtimeReq, nil
}

func (rm *RequestManager) GetRequest(ctx context.Context, companyID, requestID string) (*OvertimeRequest, error) {
	return rm.repo.GetOvertimeRequest(ctx, companyID, requestID)
}

func (rm *RequestManager) ListRequests(ctx context.Context, companyID string, filters OvertimeFilters) ([]OvertimeRequest, error) {
	return rm.repo.ListOvertimeRequests(ctx, companyID, filters)
}
