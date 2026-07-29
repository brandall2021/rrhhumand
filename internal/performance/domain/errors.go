package domain

import "errors"

var (
	ErrNotFound              = errors.New("not found")
	ErrInvalidStatus         = errors.New("invalid status transition")
	ErrCycleClosed           = errors.New("cycle is closed or cancelled")
	ErrEvaluationLocked      = errors.New("evaluation is locked")
	ErrEvaluationSubmitted   = errors.New("evaluation already submitted")
	ErrDuplicateEvaluator    = errors.New("evaluator already assigned")
	ErrInvalidWeight         = errors.New("weights must sum to 100")
	ErrInvalidDates          = errors.New("end date must be after start date")
	ErrInsufficientResponses = errors.New("insufficient anonymous responses")
	ErrCalibrationComplete   = errors.New("calibration session already completed")
	ErrPlanLocked            = errors.New("plan is completed or cancelled")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrValidation            = errors.New("validation error")
)
