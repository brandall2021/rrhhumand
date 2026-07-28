package ai

import (
	"context"
	"time"
)

type ParsedCV struct {
	RawText      string            `json:"raw_text"`
	FirstName    string            `json:"first_name"`
	LastName     string            `json:"last_name"`
	Email        string            `json:"email"`
	Phone        string            `json:"phone"`
	Summary      string            `json:"summary"`
	Skills       []string          `json:"skills"`
	Education    []ParsedEducation `json:"education"`
	Experience   []ParsedExperience `json:"experience"`
	Languages    []string          `json:"languages"`
	Certificates []string          `json:"certificates"`
}

type ParsedEducation struct {
	Institution  string `json:"institution"`
	Degree       string `json:"degree"`
	FieldOfStudy string `json:"field_of_study"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	Grade        string `json:"grade"`
}

type ParsedExperience struct {
	Company      string   `json:"company"`
	Position     string   `json:"position"`
	StartDate    string   `json:"start_date"`
	EndDate      string   `json:"end_date"`
	Description  string   `json:"description"`
	Achievements []string `json:"achievements"`
}

type CVParser struct{}

func NewCVParser() *CVParser {
	return &CVParser{}
}

// TODO: implement AI integration
func (p *CVParser) ParseCV(ctx context.Context, documentContent []byte) (*ParsedCV, error) {
	return &ParsedCV{
		RawText:   string(documentContent),
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Phone:     "+1234567890",
		Summary:   "Experienced professional parsed from CV document.",
		Skills:    []string{"Go", "Python", "SQL", "Project Management"},
		Education: []ParsedEducation{
			{
				Institution:  "University of Example",
				Degree:       "Bachelor",
				FieldOfStudy: "Computer Science",
				StartDate:    "2015",
				EndDate:      "2019",
			},
		},
		Experience: []ParsedExperience{
			{
				Company:     "Example Corp",
				Position:    "Software Engineer",
				StartDate:   "2019",
				EndDate:     "2024",
				Description: "Worked on backend systems",
				Achievements: []string{"Scaled system to 1M users"},
			},
		},
		Languages:    []string{"English", "Spanish"},
		Certificates: []string{"AWS Certified"},
	}, nil
}

// TODO: implement AI integration
func (p *CVParser) ExtractSkills(text string) ([]string, error) {
	return []string{"Go", "Python", "SQL", "JavaScript", "Docker"}, nil
}

// TODO: implement AI integration
func (p *CVParser) ExtractEducation(text string) ([]*ParsedEducation, error) {
	return []*ParsedEducation{
		{
			Institution:  "Parsed Institution",
			Degree:       "Bachelor",
			FieldOfStudy: "Computer Science",
			StartDate:    "2015",
			EndDate:      "2019",
		},
	}, nil
}

// TODO: implement AI integration
func (p *CVParser) ExtractExperience(text string) ([]*ParsedExperience, error) {
	return []*ParsedExperience{
		{
			Company:     "Parsed Company",
			Position:    "Software Developer",
			StartDate:   "2019",
			EndDate:     "2024",
			Description: "Developed software solutions",
		},
	}, nil
}

// TODO: implement AI integration
func (p *CVParser) ExtractLanguages(text string) ([]string, error) {
	return []string{"English", "Spanish"}, nil
}

// GetParserVersion returns the current parser version
func (p *CVParser) GetParserVersion() string {
	return "ai-stub-v1.0"
}

// ParseCVWithProgress is a future hook for streaming parse results
// TODO: implement AI integration
func (p *CVParser) ParseCVWithProgress(ctx context.Context, documentContent []byte, progressFn func(progress float64)) (*ParsedCV, error) {
	progressFn(0.0)
	result, err := p.ParseCV(ctx, documentContent)
	if err != nil {
		return nil, err
	}
	progressFn(1.0)
	return result, nil
}

// validateParsedCV checks that the parsed data meets minimum requirements
func validateParsedCV(p *ParsedCV) bool {
	return p.FirstName != "" || p.LastName != "" || p.Email != ""
}
