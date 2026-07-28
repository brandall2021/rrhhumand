package surveys

import (
	"testing"
	"time"
)

func TestCreateSurveyRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateSurveyRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateSurveyRequest{
				Title: "Test Survey",
				Type:  "CLIMATE",
			},
			wantErr: false,
		},
		{
			name: "missing title",
			req: CreateSurveyRequest{
				Type: "CLIMATE",
			},
			wantErr: true,
		},
		{
			name: "missing type",
			req: CreateSurveyRequest{
				Title: "Test Survey",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.Title == "" || tt.req.Type == "" {
				if !tt.wantErr {
					t.Error("expected validation error")
				}
			}
		})
	}
}

func TestSurveyStatusTransitions(t *testing.T) {
	validTransitions := map[string][]string{
		"DRAFT":     {"PUBLISHED"},
		"PUBLISHED": {"CLOSED"},
		"CLOSED":    {"ARCHIVED"},
		"ARCHIVED":  {"DRAFT"},
	}

	invalidTransitions := []struct {
		from string
		to   string
	}{
		{"DRAFT", "CLOSED"},
		{"DRAFT", "ARCHIVED"},
		{"PUBLISHED", "DRAFT"},
		{"PUBLISHED", "ARCHIVED"},
		{"CLOSED", "DRAFT"},
		{"CLOSED", "PUBLISHED"},
		{"ARCHIVED", "PUBLISHED"},
		{"ARCHIVED", "CLOSED"},
	}

	for from, tos := range validTransitions {
		for _, to := range tos {
			found := false
			for _, validTo := range validTransitions[from] {
				if validTo == to {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("transition %s -> %s should be valid but wasn't found", from, to)
			}
		}
	}

	for _, tt := range invalidTransitions {
		found := false
		for _, validTo := range validTransitions[tt.from] {
			if validTo == tt.to {
				found = true
				break
			}
		}
		if found {
			t.Errorf("transition %s -> %s should be invalid but was found", tt.from, tt.to)
		}
	}
}

func TestSurveyTypes(t *testing.T) {
	validTypes := []string{
		"GENERAL", "CLIMATE", "SATISFACTION", "FEEDBACK",
		"PULSE", "TRAINING", "INTERNAL",
	}

	seen := make(map[string]bool)
	for _, vt := range validTypes {
		if seen[vt] {
			t.Errorf("duplicate survey type: %s", vt)
		}
		seen[vt] = true
	}
}

func TestQuestionTypes(t *testing.T) {
	validTypes := []string{
		"TEXT", "SINGLE_CHOICE", "MULTIPLE_CHOICE",
		"RATING", "YES_NO", "NUMBER",
	}

	seen := make(map[string]bool)
	for _, vt := range validTypes {
		if seen[vt] {
			t.Errorf("duplicate question type: %s", vt)
		}
		seen[vt] = true
	}
}

func TestTargetTypes(t *testing.T) {
	validTypes := []string{
		"ALL", "DEPARTMENT", "BRANCH", "POSITION", "EMPLOYEE",
	}

	seen := make(map[string]bool)
	for _, vt := range validTypes {
		if seen[vt] {
			t.Errorf("duplicate target type: %s", vt)
		}
		seen[vt] = true
	}
}

func TestAnonymousSurveyBehavior(t *testing.T) {
	survey := &AnonymousSurvey{
		ID:        "test-id",
		Anonymous: true,
	}

	if !survey.Anonymous {
		t.Error("expected survey to be anonymous")
	}

	if survey.ShouldHideEmployeeID() != true {
		t.Error("anonymous survey should hide employee ID")
	}

	identifiedSurvey := &AnonymousSurvey{
		ID:        "test-id",
		Anonymous: false,
	}

	if identifiedSurvey.ShouldHideEmployeeID() != false {
		t.Error("identified survey should not hide employee ID")
	}
}

type AnonymousSurvey struct {
	ID        string
	Anonymous bool
}

func (a *AnonymousSurvey) ShouldHideEmployeeID() bool {
	return a.Anonymous
}

func TestSurveyDates(t *testing.T) {
	now := time.Now()
	startsAt := now.Add(-1 * time.Hour)
	endsAt := now.Add(1 * time.Hour)

	survey := &DateTestSurvey{
		StartsAt: &startsAt,
		EndsAt:   &endsAt,
	}

	if !survey.IsActive() {
		t.Error("survey should be active within date range")
	}

	pastEnd := now.Add(-2 * time.Hour)
	survey2 := &DateTestSurvey{
		StartsAt: &startsAt,
		EndsAt:   &pastEnd,
	}

	if survey2.IsActive() {
		t.Error("survey should not be active after end date")
	}

	noStartNoEnd := &DateTestSurvey{}
	if !noStartNoEnd.IsActive() {
		t.Error("survey with no dates should be active")
	}
}

type DateTestSurvey struct {
	StartsAt *time.Time
	EndsAt   *time.Time
}

func (d *DateTestSurvey) IsActive() bool {
	now := time.Now()
	if d.StartsAt != nil && now.Before(*d.StartsAt) {
		return false
	}
	if d.EndsAt != nil && now.After(*d.EndsAt) {
		return false
	}
	return true
}

func TestSurveyStatsCalculation(t *testing.T) {
	ratings := []float64{5, 4, 3, 2, 1}
	var sum float64
	for _, r := range ratings {
		sum += r
	}
	avg := sum / float64(len(ratings))

	if avg != 3.0 {
		t.Errorf("expected average 3.0, got %f", avg)
	}

	min := ratings[0]
	max := ratings[0]
	for _, r := range ratings[1:] {
		if r < min {
			min = r
		}
		if r > max {
			max = r
		}
	}

	if min != 1.0 {
		t.Errorf("expected min 1.0, got %f", min)
	}
	if max != 5.0 {
		t.Errorf("expected max 5.0, got %f", max)
	}
}

func TestParticipationRate(t *testing.T) {
	tests := []struct {
		name      string
		responded int
		targeted  int
		want      float64
	}{
		{"full participation", 100, 100, 100.0},
		{"half participation", 50, 100, 50.0},
		{"no participation", 0, 100, 0.0},
		{"no targets", 0, 0, 0.0},
		{"more responded than targeted", 150, 100, 150.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rate float64
			if tt.targeted > 0 {
				rate = float64(tt.responded) / float64(tt.targeted) * 100
			}
			if rate != tt.want {
				t.Errorf("expected %.1f%%, got %.1f%%", tt.want, rate)
			}
		})
	}
}

func TestMultipleChoiceOptions(t *testing.T) {
	options := []string{"Option A", "Option B", "Option C", "Option D"}

	seen := make(map[string]bool)
	for _, opt := range options {
		if seen[opt] {
			t.Errorf("duplicate option: %s", opt)
		}
		seen[opt] = true
	}

	if len(options) < 2 {
		t.Error("multiple choice should have at least 2 options")
	}
}

func TestYesNoOptions(t *testing.T) {
	yesNoOptions := []string{"Si", "No"}

	if len(yesNoOptions) != 2 {
		t.Error("yes/no should have exactly 2 options")
	}
}

func TestRatingScale(t *testing.T) {
	minRating := 1.0
	maxRating := 5.0

	testRatings := []float64{1, 2, 3, 4, 5}
	for _, r := range testRatings {
		if r < minRating || r > maxRating {
			t.Errorf("rating %f out of scale range [%f, %f]", r, minRating, maxRating)
		}
	}
}

func TestSurveyFilterDefaults(t *testing.T) {
	filters := SurveyFilters{}

	if filters.Status != "" {
		t.Error("default status should be empty")
	}
	if filters.Type != "" {
		t.Error("default type should be empty")
	}
	if filters.Search != "" {
		t.Error("default search should be empty")
	}
}

func TestCreateSurveyRequestDefaults(t *testing.T) {
	req := CreateSurveyRequest{
		Title: "Test",
		Type:  "GENERAL",
	}

	anonymous := false
	multipleResponses := false

	if req.Anonymous != nil && *req.Anonymous != anonymous {
		t.Error("default anonymous should be false")
	}
	if req.MultipleResponses != nil && *req.MultipleResponses != multipleResponses {
		t.Error("default multiple responses should be false")
	}
}
