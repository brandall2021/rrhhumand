package integration

import (
	"context"
	"log"
)

type AccessProvisioningAdapter struct{}

func NewAccessProvisioningAdapter() *AccessProvisioningAdapter {
	return &AccessProvisioningAdapter{}
}

func (a *AccessProvisioningAdapter) CreateAccount(ctx context.Context, employeeID, systemName, accessType string) error {
	log.Printf("[AccessProvisioningAdapter] CreateAccount employee=%s system=%s type=%s", employeeID, systemName, accessType)
	return nil
}

func (a *AccessProvisioningAdapter) DisableAccount(ctx context.Context, employeeID, systemName string) error {
	log.Printf("[AccessProvisioningAdapter] DisableAccount employee=%s system=%s", employeeID, systemName)
	return nil
}

func (a *AccessProvisioningAdapter) RevokeAccess(ctx context.Context, employeeID, systemName string) error {
	log.Printf("[AccessProvisioningAdapter] RevokeAccess employee=%s system=%s", employeeID, systemName)
	return nil
}
