package profile

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type ProfileService struct {
	repo *ProfileRepository
}

func NewProfileService(repo *ProfileRepository) *ProfileService {
	return &ProfileService{repo: repo}
}

type UpdateProfileRequest struct {
	Phone    *string `json:"phone,omitempty"`
	PhotoURL *string `json:"photo_url,omitempty"`
}

func (s *ProfileService) Get(ctx context.Context, userID, companyID string) (*EmployeeProfile, error) {
	profile, err := s.repo.GetByUserID(ctx, userID, companyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("employee profile not found for this user in this company")
		}
		return nil, err
	}
	return profile, nil
}

func (s *ProfileService) Update(ctx context.Context, userID, companyID string, req *UpdateProfileRequest) (*EmployeeProfile, error) {
	empID, err := s.repo.GetEmployeeIDByUser(ctx, userID, companyID)
	if err != nil {
		return nil, errors.New("employee not found for this user")
	}

	if err := s.repo.UpdateAllowedFields(ctx, empID, companyID, req.Phone, req.PhotoURL); err != nil {
		return nil, err
	}

	return s.repo.GetByUserID(ctx, userID, companyID)
}
