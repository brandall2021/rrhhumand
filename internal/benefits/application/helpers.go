package application

import "fmt"

func svcErr(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("benefits_svc.%s: %w", op, err)
}
