package domain

import "errors"

var (
    ErrNotFound           = errors.New("not found")
    ErrInvalidInput       = errors.New("invalid input")
    ErrAlreadyExists      = errors.New("already exists")
    ErrInvalidStatus      = errors.New("invalid status transition")
    ErrStageNotAllowed    = errors.New("stage transition not allowed")
    ErrRequisitionClosed  = errors.New("requisition is closed")
    ErrPostingClosed      = errors.New("posting is closed")
    ErrCandidateBlacklisted = errors.New("candidate is blacklisted")
    ErrDuplicateApplication = errors.New("candidate already applied to this posting")
    ErrInterviewConflict  = errors.New("interviewer has a scheduling conflict")
    ErrOfferExpired       = errors.New("offer has expired")
    ErrOfferAlreadyResponded = errors.New("offer already accepted or rejected")
    ErrInsufficientVacancies = errors.New("no vacancies remaining")
    ErrWorkflowInactive   = errors.New("workflow is inactive")
    ErrPermissionDenied   = errors.New("permission denied")
    ErrScoringModelInactive = errors.New("scoring model is inactive")
    ErrParsingFailed      = errors.New("CV parsing failed")
    ErrMatchingFailed     = errors.New("matching failed")
    ErrIntegrationError   = errors.New("integration error")
    ErrEmailSendFailed    = errors.New("failed to send email")
    ErrProcessInProgress  = errors.New("hiring process already in progress")
)
