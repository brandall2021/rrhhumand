package integration

import (
	"context"
	"log"
)

type TrainingAdapter struct{}

func NewTrainingAdapter() *TrainingAdapter {
	return &TrainingAdapter{}
}

func (a *TrainingAdapter) AssignCourse(ctx context.Context, companyID, employeeID, courseName string, mandatory bool, dueDate string) error {
	log.Printf("[TrainingAdapter] AssignCourse company=%s employee=%s course=%s mandatory=%v due=%s", companyID, employeeID, courseName, mandatory, dueDate)
	return nil
}

func (a *TrainingAdapter) GetCourseStatus(ctx context.Context, employeeID, courseName string) (string, error) {
	log.Printf("[TrainingAdapter] GetCourseStatus employee=%s course=%s", employeeID, courseName)
	return "PENDING", nil
}
