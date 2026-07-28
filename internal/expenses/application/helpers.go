package application

import (
	"context"
	"fmt"

	"github.com/rrhhumand/api/internal/expenses/domain"
)

type AuditRepository interface {
	Log(ctx context.Context, entry *domain.ExpenseAuditLog) error
}

func svcErr(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("expenses_svc.%s: %w", op, err)
}
