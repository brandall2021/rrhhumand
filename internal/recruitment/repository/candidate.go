package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/recruitment/domain"
)

type CandidateRepo struct {
	pool *pgxpool.Pool
}

func NewCandidateRepo(pool *pgxpool.Pool) *CandidateRepo {
	return &CandidateRepo{pool: pool}
}

func (r *CandidateRepo) Create(ctx context.Context, companyID string, req *domain.Candidate) (*domain.Candidate, error) {
	c := &domain.Candidate{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO candidates (company_id, first_name, last_name, email, phone, phone_country_code, document_type, document_number, birth_date, location, nationality, gender, linkedin_url, portfolio_url, github_url, personal_website, current_company, current_position, notice_period, salary_expectation_min, salary_expectation_max, salary_currency, availability, source, source_detail, is_employee_referral, referrer_employee_id, blacklisted, blacklist_reason, tags, notes)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31)
		 RETURNING id, company_id, first_name, last_name, email, phone, phone_country_code, document_type, document_number, birth_date, location, nationality, gender, linkedin_url, portfolio_url, github_url, personal_website, current_company, current_position, notice_period, salary_expectation_min, salary_expectation_max, salary_currency, availability, source, source_detail, is_employee_referral, referrer_employee_id, status, blacklisted, blacklist_reason, tags, notes, created_at, updated_at`,
		companyID, req.FirstName, req.LastName, req.Email, req.Phone, req.PhoneCountryCode,
		req.DocumentType, req.DocumentNumber, req.BirthDate, req.Location, req.Nationality,
		req.Gender, req.LinkedInURL, req.PortfolioURL, req.GithubURL, req.PersonalWebsite,
		req.CurrentCompany, req.CurrentPosition, req.NoticePeriod, req.SalaryExpectMin,
		req.SalaryExpectMax, req.SalaryCurrency, req.Availability, req.Source, req.SourceDetail,
		req.IsReferral, req.ReferrerEmployeeID, req.Blacklisted, req.BlacklistReason, req.Tags, req.Notes,
	).Scan(&c.ID, &c.CompanyID, &c.FirstName, &c.LastName, &c.Email, &c.Phone, &c.PhoneCountryCode,
		&c.DocumentType, &c.DocumentNumber, &c.BirthDate, &c.Location, &c.Nationality, &c.Gender,
		&c.LinkedInURL, &c.PortfolioURL, &c.GithubURL, &c.PersonalWebsite, &c.CurrentCompany,
		&c.CurrentPosition, &c.NoticePeriod, &c.SalaryExpectMin, &c.SalaryExpectMax, &c.SalaryCurrency,
		&c.Availability, &c.Source, &c.SourceDetail, &c.IsReferral, &c.ReferrerEmployeeID,
		&c.Status, &c.Blacklisted, &c.BlacklistReason, &c.Tags, &c.Notes, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *CandidateRepo) GetByID(ctx context.Context, companyID, id string) (*domain.Candidate, error) {
	c := &domain.Candidate{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, first_name, last_name, email, phone, phone_country_code, document_type, document_number, birth_date, location, nationality, gender, linkedin_url, portfolio_url, github_url, personal_website, current_company, current_position, notice_period, salary_expectation_min, salary_expectation_max, salary_currency, availability, source, source_detail, is_employee_referral, referrer_employee_id, status, blacklisted, blacklist_reason, tags, notes, created_at, updated_at
		 FROM candidates WHERE company_id=$1 AND id=$2`, companyID, id,
	).Scan(&c.ID, &c.CompanyID, &c.FirstName, &c.LastName, &c.Email, &c.Phone, &c.PhoneCountryCode,
		&c.DocumentType, &c.DocumentNumber, &c.BirthDate, &c.Location, &c.Nationality, &c.Gender,
		&c.LinkedInURL, &c.PortfolioURL, &c.GithubURL, &c.PersonalWebsite, &c.CurrentCompany,
		&c.CurrentPosition, &c.NoticePeriod, &c.SalaryExpectMin, &c.SalaryExpectMax, &c.SalaryCurrency,
		&c.Availability, &c.Source, &c.SourceDetail, &c.IsReferral, &c.ReferrerEmployeeID,
		&c.Status, &c.Blacklisted, &c.BlacklistReason, &c.Tags, &c.Notes, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *CandidateRepo) List(ctx context.Context, companyID string, status, source string) ([]domain.Candidate, error) {
	query := `SELECT id, company_id, first_name, last_name, email, phone, phone_country_code, document_type, document_number, birth_date, location, nationality, gender, linkedin_url, portfolio_url, github_url, personal_website, current_company, current_position, notice_period, salary_expectation_min, salary_expectation_max, salary_currency, availability, source, source_detail, is_employee_referral, referrer_employee_id, status, blacklisted, blacklist_reason, tags, notes, created_at, updated_at
		 FROM candidates WHERE company_id=$1`
	args := []interface{}{companyID}
	argIdx := 2

	if status != "" {
		query += fmt.Sprintf(" AND status=$%d", argIdx)
		args = append(args, status)
		argIdx++
	}
	if source != "" {
		query += fmt.Sprintf(" AND source=$%d", argIdx)
		args = append(args, source)
		argIdx++
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var candidates []domain.Candidate
	for rows.Next() {
		var c domain.Candidate
		rows.Scan(&c.ID, &c.CompanyID, &c.FirstName, &c.LastName, &c.Email, &c.Phone, &c.PhoneCountryCode,
			&c.DocumentType, &c.DocumentNumber, &c.BirthDate, &c.Location, &c.Nationality, &c.Gender,
			&c.LinkedInURL, &c.PortfolioURL, &c.GithubURL, &c.PersonalWebsite, &c.CurrentCompany,
			&c.CurrentPosition, &c.NoticePeriod, &c.SalaryExpectMin, &c.SalaryExpectMax, &c.SalaryCurrency,
			&c.Availability, &c.Source, &c.SourceDetail, &c.IsReferral, &c.ReferrerEmployeeID,
			&c.Status, &c.Blacklisted, &c.BlacklistReason, &c.Tags, &c.Notes, &c.CreatedAt, &c.UpdatedAt)
		candidates = append(candidates, c)
	}
	return candidates, nil
}

func (r *CandidateRepo) Update(ctx context.Context, companyID, id string, req *domain.Candidate) (*domain.Candidate, error) {
	c := &domain.Candidate{}
	err := r.pool.QueryRow(ctx,
		`UPDATE candidates SET
		 first_name=COALESCE($3,first_name), last_name=COALESCE($4,last_name), phone=COALESCE($5,phone),
		 phone_country_code=COALESCE($6,phone_country_code), document_type=COALESCE($7,document_type),
		 document_number=COALESCE($8,document_number), birth_date=COALESCE($9,birth_date),
		 location=COALESCE($10,location), nationality=COALESCE($11,nationality),
		 linkedin_url=COALESCE($12,linkedin_url), portfolio_url=COALESCE($13,portfolio_url),
		 github_url=COALESCE($14,github_url), personal_website=COALESCE($15,personal_website),
		 current_company=COALESCE($16,current_company), current_position=COALESCE($17,current_position),
		 notice_period=COALESCE($18,notice_period), salary_expectation_min=COALESCE($19,salary_expectation_min),
		 salary_expectation_max=COALESCE($20,salary_expectation_max), salary_currency=COALESCE($21,salary_currency),
		 availability=COALESCE($22,availability), source=COALESCE($23,source), source_detail=COALESCE($24,source_detail),
		 notes=COALESCE($25,notes), updated_at=NOW()
		 WHERE company_id=$1 AND id=$2
		 RETURNING id, company_id, first_name, last_name, email, phone, phone_country_code, document_type, document_number, birth_date, location, nationality, gender, linkedin_url, portfolio_url, github_url, personal_website, current_company, current_position, notice_period, salary_expectation_min, salary_expectation_max, salary_currency, availability, source, source_detail, is_employee_referral, referrer_employee_id, status, blacklisted, blacklist_reason, tags, notes, created_at, updated_at`,
		companyID, id, req.FirstName, req.LastName, req.Phone, req.PhoneCountryCode,
		req.DocumentType, req.DocumentNumber, req.BirthDate, req.Location, req.Nationality,
		req.LinkedInURL, req.PortfolioURL, req.GithubURL, req.PersonalWebsite,
		req.CurrentCompany, req.CurrentPosition, req.NoticePeriod, req.SalaryExpectMin,
		req.SalaryExpectMax, req.SalaryCurrency, req.Availability, req.Source, req.SourceDetail, req.Notes,
	).Scan(&c.ID, &c.CompanyID, &c.FirstName, &c.LastName, &c.Email, &c.Phone, &c.PhoneCountryCode,
		&c.DocumentType, &c.DocumentNumber, &c.BirthDate, &c.Location, &c.Nationality, &c.Gender,
		&c.LinkedInURL, &c.PortfolioURL, &c.GithubURL, &c.PersonalWebsite, &c.CurrentCompany,
		&c.CurrentPosition, &c.NoticePeriod, &c.SalaryExpectMin, &c.SalaryExpectMax, &c.SalaryCurrency,
		&c.Availability, &c.Source, &c.SourceDetail, &c.IsReferral, &c.ReferrerEmployeeID,
		&c.Status, &c.Blacklisted, &c.BlacklistReason, &c.Tags, &c.Notes, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *CandidateRepo) UpdateStatus(ctx context.Context, companyID, id, status string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE candidates SET status=$3, updated_at=NOW() WHERE company_id=$1 AND id=$2`,
		companyID, id, status)
	return err
}

func (r *CandidateRepo) AddEducation(ctx context.Context, req *domain.CandidateEducation) (*domain.CandidateEducation, error) {
	e := &domain.CandidateEducation{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO candidate_education (candidate_id, institution, degree, field_of_study, start_date, end_date, is_current, grade, description)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		 RETURNING id, candidate_id, institution, degree, field_of_study, start_date, end_date, is_current, grade, description, created_at`,
		req.CandidateID, req.Institution, req.Degree, req.FieldOfStudy, req.StartDate, req.EndDate,
		req.IsCurrent, req.Grade, req.Description,
	).Scan(&e.ID, &e.CandidateID, &e.Institution, &e.Degree, &e.FieldOfStudy, &e.StartDate, &e.EndDate,
		&e.IsCurrent, &e.Grade, &e.Description, &e.CreatedAt)
	return e, err
}

func (r *CandidateRepo) UpdateEducation(ctx context.Context, id string, req *domain.CandidateEducation) (*domain.CandidateEducation, error) {
	e := &domain.CandidateEducation{}
	err := r.pool.QueryRow(ctx,
		`UPDATE candidate_education SET institution=COALESCE($2,institution), degree=COALESCE($3,degree),
		 field_of_study=COALESCE($4,field_of_study), start_date=COALESCE($5,start_date),
		 end_date=COALESCE($6,end_date), is_current=COALESCE($7,is_current),
		 grade=COALESCE($8,grade), description=COALESCE($9,description) WHERE id=$1
		 RETURNING id, candidate_id, institution, degree, field_of_study, start_date, end_date, is_current, grade, description, created_at`,
		id, req.Institution, req.Degree, req.FieldOfStudy, req.StartDate, req.EndDate,
		req.IsCurrent, req.Grade, req.Description,
	).Scan(&e.ID, &e.CandidateID, &e.Institution, &e.Degree, &e.FieldOfStudy, &e.StartDate, &e.EndDate,
		&e.IsCurrent, &e.Grade, &e.Description, &e.CreatedAt)
	return e, err
}

func (r *CandidateRepo) DeleteEducation(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM candidate_education WHERE id=$1`, id)
	return err
}

func (r *CandidateRepo) ListEducation(ctx context.Context, candidateID string) ([]domain.CandidateEducation, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, candidate_id, institution, degree, field_of_study, start_date, end_date, is_current, grade, description, created_at
		 FROM candidate_education WHERE candidate_id=$1 ORDER BY end_date DESC NULLS LAST`, candidateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var edu []domain.CandidateEducation
	for rows.Next() {
		var e domain.CandidateEducation
		rows.Scan(&e.ID, &e.CandidateID, &e.Institution, &e.Degree, &e.FieldOfStudy, &e.StartDate, &e.EndDate,
			&e.IsCurrent, &e.Grade, &e.Description, &e.CreatedAt)
		edu = append(edu, e)
	}
	return edu, nil
}

func (r *CandidateRepo) AddExperience(ctx context.Context, req *domain.CandidateExperience) (*domain.CandidateExperience, error) {
	e := &domain.CandidateExperience{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO candidate_experience (candidate_id, company, position, location, start_date, end_date, is_current, description, achievements, industry, salary, salary_currency)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		 RETURNING id, candidate_id, company, position, location, start_date, end_date, is_current, description, achievements, industry, salary, salary_currency, created_at`,
		req.CandidateID, req.Company, req.Position, req.Location, req.StartDate, req.EndDate,
		req.IsCurrent, req.Description, req.Achievements, req.Industry, req.Salary, req.SalaryCurrency,
	).Scan(&e.ID, &e.CandidateID, &e.Company, &e.Position, &e.Location, &e.StartDate, &e.EndDate,
		&e.IsCurrent, &e.Description, &e.Achievements, &e.Industry, &e.Salary, &e.SalaryCurrency, &e.CreatedAt)
	return e, err
}

func (r *CandidateRepo) UpdateExperience(ctx context.Context, id string, req *domain.CandidateExperience) (*domain.CandidateExperience, error) {
	e := &domain.CandidateExperience{}
	err := r.pool.QueryRow(ctx,
		`UPDATE candidate_experience SET company=COALESCE($2,company), position=COALESCE($3,position),
		 location=COALESCE($4,location), start_date=COALESCE($5,start_date), end_date=COALESCE($6,end_date),
		 is_current=COALESCE($7,is_current), description=COALESCE($8,description),
		 achievements=COALESCE($9,achievements), industry=COALESCE($10,industry),
		 salary=COALESCE($11,salary), salary_currency=COALESCE($12,salary_currency) WHERE id=$1
		 RETURNING id, candidate_id, company, position, location, start_date, end_date, is_current, description, achievements, industry, salary, salary_currency, created_at`,
		id, req.Company, req.Position, req.Location, req.StartDate, req.EndDate,
		req.IsCurrent, req.Description, req.Achievements, req.Industry, req.Salary, req.SalaryCurrency,
	).Scan(&e.ID, &e.CandidateID, &e.Company, &e.Position, &e.Location, &e.StartDate, &e.EndDate,
		&e.IsCurrent, &e.Description, &e.Achievements, &e.Industry, &e.Salary, &e.SalaryCurrency, &e.CreatedAt)
	return e, err
}

func (r *CandidateRepo) DeleteExperience(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM candidate_experience WHERE id=$1`, id)
	return err
}

func (r *CandidateRepo) ListExperience(ctx context.Context, candidateID string) ([]domain.CandidateExperience, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, candidate_id, company, position, location, start_date, end_date, is_current, description, achievements, industry, salary, salary_currency, created_at
		 FROM candidate_experience WHERE candidate_id=$1 ORDER BY end_date DESC NULLS LAST`, candidateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exp []domain.CandidateExperience
	for rows.Next() {
		var e domain.CandidateExperience
		rows.Scan(&e.ID, &e.CandidateID, &e.Company, &e.Position, &e.Location, &e.StartDate, &e.EndDate,
			&e.IsCurrent, &e.Description, &e.Achievements, &e.Industry, &e.Salary, &e.SalaryCurrency, &e.CreatedAt)
		exp = append(exp, e)
	}
	return exp, nil
}

func (r *CandidateRepo) AddSkill(ctx context.Context, req *domain.CandidateSkill) (*domain.CandidateSkill, error) {
	s := &domain.CandidateSkill{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO candidate_skills (candidate_id, skill, category, proficiency, years_experience, is_primary)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 ON CONFLICT (candidate_id, skill) DO UPDATE SET proficiency=$4, years_experience=$5, is_primary=$6
		 RETURNING id, candidate_id, skill, category, proficiency, years_experience, is_primary, created_at`,
		req.CandidateID, req.Skill, req.Category, req.Proficiency, req.YearsExp, req.IsPrimary,
	).Scan(&s.ID, &s.CandidateID, &s.Skill, &s.Category, &s.Proficiency, &s.YearsExp, &s.IsPrimary, &s.CreatedAt)
	return s, err
}

func (r *CandidateRepo) UpdateSkill(ctx context.Context, id string, req *domain.CandidateSkill) (*domain.CandidateSkill, error) {
	s := &domain.CandidateSkill{}
	err := r.pool.QueryRow(ctx,
		`UPDATE candidate_skills SET skill=COALESCE($2,skill), category=COALESCE($3,category),
		 proficiency=COALESCE($4,proficiency), years_experience=COALESCE($5,years_experience),
		 is_primary=COALESCE($6,is_primary) WHERE id=$1
		 RETURNING id, candidate_id, skill, category, proficiency, years_experience, is_primary, created_at`,
		id, req.Skill, req.Category, req.Proficiency, req.YearsExp, req.IsPrimary,
	).Scan(&s.ID, &s.CandidateID, &s.Skill, &s.Category, &s.Proficiency, &s.YearsExp, &s.IsPrimary, &s.CreatedAt)
	return s, err
}

func (r *CandidateRepo) DeleteSkill(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM candidate_skills WHERE id=$1`, id)
	return err
}

func (r *CandidateRepo) ListSkills(ctx context.Context, candidateID string) ([]domain.CandidateSkill, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, candidate_id, skill, category, proficiency, years_experience, is_primary, created_at
		 FROM candidate_skills WHERE candidate_id=$1 ORDER BY is_primary DESC, skill`, candidateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []domain.CandidateSkill
	for rows.Next() {
		var s domain.CandidateSkill
		rows.Scan(&s.ID, &s.CandidateID, &s.Skill, &s.Category, &s.Proficiency, &s.YearsExp, &s.IsPrimary, &s.CreatedAt)
		skills = append(skills, s)
	}
	return skills, nil
}

func (r *CandidateRepo) AddCertification(ctx context.Context, req *domain.CandidateCertification) (*domain.CandidateCertification, error) {
	cc := &domain.CandidateCertification{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO candidate_certifications (candidate_id, name, issuer, issue_date, expiry_date, credential_id, credential_url)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)
		 RETURNING id, candidate_id, name, issuer, issue_date, expiry_date, credential_id, credential_url, created_at`,
		req.CandidateID, req.Name, req.Issuer, req.IssueDate, req.ExpiryDate, req.CredentialID, req.CredentialURL,
	).Scan(&cc.ID, &cc.CandidateID, &cc.Name, &cc.Issuer, &cc.IssueDate, &cc.ExpiryDate, &cc.CredentialID, &cc.CredentialURL, &cc.CreatedAt)
	return cc, err
}

func (r *CandidateRepo) DeleteCertification(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM candidate_certifications WHERE id=$1`, id)
	return err
}

func (r *CandidateRepo) ListCertifications(ctx context.Context, candidateID string) ([]domain.CandidateCertification, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, candidate_id, name, issuer, issue_date, expiry_date, credential_id, credential_url, created_at
		 FROM candidate_certifications WHERE candidate_id=$1 ORDER BY issue_date DESC NULLS LAST`, candidateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var certs []domain.CandidateCertification
	for rows.Next() {
		var cc domain.CandidateCertification
		rows.Scan(&cc.ID, &cc.CandidateID, &cc.Name, &cc.Issuer, &cc.IssueDate, &cc.ExpiryDate, &cc.CredentialID, &cc.CredentialURL, &cc.CreatedAt)
		certs = append(certs, cc)
	}
	return certs, nil
}

func (r *CandidateRepo) AddLanguage(ctx context.Context, req *domain.CandidateLanguage) (*domain.CandidateLanguage, error) {
	l := &domain.CandidateLanguage{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO candidate_languages (candidate_id, language, proficiency, is_native)
		 VALUES ($1,$2,$3,$4)
		 ON CONFLICT (candidate_id, language) DO UPDATE SET proficiency=$3, is_native=$4
		 RETURNING id, candidate_id, language, proficiency, is_native, created_at`,
		req.CandidateID, req.Language, req.Proficiency, req.IsNative,
	).Scan(&l.ID, &l.CandidateID, &l.Language, &l.Proficiency, &l.IsNative, &l.CreatedAt)
	return l, err
}

func (r *CandidateRepo) UpdateLanguage(ctx context.Context, id string, req *domain.CandidateLanguage) (*domain.CandidateLanguage, error) {
	l := &domain.CandidateLanguage{}
	err := r.pool.QueryRow(ctx,
		`UPDATE candidate_languages SET language=COALESCE($2,language), proficiency=COALESCE($3,proficiency),
		 is_native=COALESCE($4,is_native) WHERE id=$1
		 RETURNING id, candidate_id, language, proficiency, is_native, created_at`,
		id, req.Language, req.Proficiency, req.IsNative,
	).Scan(&l.ID, &l.CandidateID, &l.Language, &l.Proficiency, &l.IsNative, &l.CreatedAt)
	return l, err
}

func (r *CandidateRepo) DeleteLanguage(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM candidate_languages WHERE id=$1`, id)
	return err
}

func (r *CandidateRepo) ListLanguages(ctx context.Context, candidateID string) ([]domain.CandidateLanguage, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, candidate_id, language, proficiency, is_native, created_at
		 FROM candidate_languages WHERE candidate_id=$1`, candidateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var langs []domain.CandidateLanguage
	for rows.Next() {
		var l domain.CandidateLanguage
		rows.Scan(&l.ID, &l.CandidateID, &l.Language, &l.Proficiency, &l.IsNative, &l.CreatedAt)
		langs = append(langs, l)
	}
	return langs, nil
}

func (r *CandidateRepo) AddDocument(ctx context.Context, req *domain.CandidateDocument) (*domain.CandidateDocument, error) {
	d := &domain.CandidateDocument{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO candidate_documents (candidate_id, company_id, document_type, file_name, mime_type, size_bytes, storage_provider, storage_key)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 RETURNING id, candidate_id, company_id, document_type, file_name, mime_type, size_bytes, storage_provider, storage_key, parsed_data, parsed_at, created_at`,
		req.CandidateID, req.CompanyID, req.DocumentType, req.FileName, req.MimeType, req.SizeBytes,
		req.StorageProvider, req.StorageKey,
	).Scan(&d.ID, &d.CandidateID, &d.CompanyID, &d.DocumentType, &d.FileName, &d.MimeType, &d.SizeBytes,
		&d.StorageProvider, &d.StorageKey, &d.ParsedData, &d.ParsedAt, &d.CreatedAt)
	return d, err
}

func (r *CandidateRepo) ListDocuments(ctx context.Context, candidateID string) ([]domain.CandidateDocument, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, candidate_id, company_id, document_type, file_name, mime_type, size_bytes, storage_provider, storage_key, parsed_data, parsed_at, created_at
		 FROM candidate_documents WHERE candidate_id=$1 ORDER BY created_at DESC`, candidateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []domain.CandidateDocument
	for rows.Next() {
		var d domain.CandidateDocument
		rows.Scan(&d.ID, &d.CandidateID, &d.CompanyID, &d.DocumentType, &d.FileName, &d.MimeType, &d.SizeBytes,
			&d.StorageProvider, &d.StorageKey, &d.ParsedData, &d.ParsedAt, &d.CreatedAt)
		docs = append(docs, d)
	}
	return docs, nil
}

func (r *CandidateRepo) CreateParsedData(ctx context.Context, req *domain.CandidateParsedData) (*domain.CandidateParsedData, error) {
	pd := &domain.CandidateParsedData{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO candidate_parsed_data (candidate_id, raw_text, structured_data, skills_found, education_found, experience_found, languages_found, certifications_found, summary, score, parser_version, parsed_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		 ON CONFLICT (candidate_id) DO UPDATE SET raw_text=$2, structured_data=$3, skills_found=$4, education_found=$5,
		 experience_found=$6, languages_found=$7, certifications_found=$8, summary=$9, score=$10, parser_version=$11, parsed_at=$12
		 RETURNING id, candidate_id, raw_text, structured_data, skills_found, education_found, experience_found, languages_found, certifications_found, summary, score, parser_version, parsed_at, created_at`,
		req.CandidateID, req.RawText, req.StructuredData, req.SkillsFound, req.EducationFound,
		req.ExperienceFound, req.LanguagesFound, req.CertsFound, req.Summary, req.Score,
		req.ParserVersion, time.Now(),
	).Scan(&pd.ID, &pd.CandidateID, &pd.RawText, &pd.StructuredData, &pd.SkillsFound, &pd.EducationFound,
		&pd.ExperienceFound, &pd.LanguagesFound, &pd.CertsFound, &pd.Summary, &pd.Score,
		&pd.ParserVersion, &pd.ParsedAt, &pd.CreatedAt)
	return pd, err
}

func (r *CandidateRepo) GetParsedData(ctx context.Context, candidateID string) (*domain.CandidateParsedData, error) {
	pd := &domain.CandidateParsedData{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, candidate_id, raw_text, structured_data, skills_found, education_found, experience_found, languages_found, certifications_found, summary, score, parser_version, parsed_at, created_at
		 FROM candidate_parsed_data WHERE candidate_id=$1`, candidateID,
	).Scan(&pd.ID, &pd.CandidateID, &pd.RawText, &pd.StructuredData, &pd.SkillsFound, &pd.EducationFound,
		&pd.ExperienceFound, &pd.LanguagesFound, &pd.CertsFound, &pd.Summary, &pd.Score,
		&pd.ParserVersion, &pd.ParsedAt, &pd.CreatedAt)
	return pd, err
}

func (r *CandidateRepo) SearchBySkills(ctx context.Context, companyID string, skills []string) ([]domain.Candidate, error) {
	query := `SELECT DISTINCT c.id, c.company_id, c.first_name, c.last_name, c.email, c.phone, c.phone_country_code,
		 c.document_type, c.document_number, c.birth_date, c.location, c.nationality, c.gender,
		 c.linkedin_url, c.portfolio_url, c.github_url, c.personal_website, c.current_company,
		 c.current_position, c.notice_period, c.salary_expectation_min, c.salary_expectation_max,
		 c.salary_currency, c.availability, c.source, c.source_detail, c.is_employee_referral,
		 c.referrer_employee_id, c.status, c.blacklisted, c.blacklist_reason, c.tags, c.notes, c.created_at, c.updated_at
		 FROM candidates c
		 INNER JOIN candidate_skills cs ON c.id = cs.candidate_id
		 WHERE c.company_id=$1 AND cs.skill = ANY($2)
		 ORDER BY c.created_at DESC`
	rows, err := r.pool.Query(ctx, query, companyID, skills)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var candidates []domain.Candidate
	for rows.Next() {
		var c domain.Candidate
		rows.Scan(&c.ID, &c.CompanyID, &c.FirstName, &c.LastName, &c.Email, &c.Phone, &c.PhoneCountryCode,
			&c.DocumentType, &c.DocumentNumber, &c.BirthDate, &c.Location, &c.Nationality, &c.Gender,
			&c.LinkedInURL, &c.PortfolioURL, &c.GithubURL, &c.PersonalWebsite, &c.CurrentCompany,
			&c.CurrentPosition, &c.NoticePeriod, &c.SalaryExpectMin, &c.SalaryExpectMax, &c.SalaryCurrency,
			&c.Availability, &c.Source, &c.SourceDetail, &c.IsReferral, &c.ReferrerEmployeeID,
			&c.Status, &c.Blacklisted, &c.BlacklistReason, &c.Tags, &c.Notes, &c.CreatedAt, &c.UpdatedAt)
		candidates = append(candidates, c)
	}
	return candidates, nil
}
