package training

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Repository struct {
	pool *pgxpool.Pool
	log  *zap.Logger
}

func NewRepository(pool *pgxpool.Pool, log *zap.Logger) *Repository {
	return &Repository{pool: pool, log: log}
}

func repoErr(op string, err error) error {
	return fmt.Errorf("training_repo.%s: %w", op, err)
}

// ---------------------------------------------------------------------------
// Categories
// ---------------------------------------------------------------------------

func (r *Repository) CreateCategory(ctx context.Context, c *Category) error {
	q := `INSERT INTO training_categories (id,company_id,name,description,parent_id,active,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.pool.Exec(ctx, q, c.ID, c.CompanyID, c.Name, c.Description, c.ParentID, c.Active, c.CreatedBy)
	return repoErr("CreateCategory", err)
}

func (r *Repository) UpdateCategory(ctx context.Context, c *Category) error {
	q := `UPDATE training_categories SET name=COALESCE($3,name), description=COALESCE($4,description),
		active=COALESCE($5,active), updated_at=NOW() WHERE id=$1 AND company_id=$2`
	_, err := r.pool.Exec(ctx, q, c.ID, c.CompanyID, c.Name, c.Description, c.Active)
	return repoErr("UpdateCategory", err)
}

func (r *Repository) GetCategory(ctx context.Context, companyID, id string) (*Category, error) {
	q := `SELECT id,company_id,name,description,parent_id,active,created_by,created_at,updated_at
		FROM training_categories WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	c := &Category{}
	err := row.Scan(&c.ID, &c.CompanyID, &c.Name, &c.Description, &c.ParentID, &c.Active,
		&c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetCategory", err)
	}
	return c, nil
}

func (r *Repository) ListCategories(ctx context.Context, companyID string, parentID *string, active *bool) ([]Category, error) {
	q := `SELECT id,company_id,name,description,parent_id,active,created_by,created_at,updated_at
		FROM training_categories WHERE company_id=$1`
	args := []any{companyID}
	idx := 2
	if parentID != nil {
		q += fmt.Sprintf(" AND parent_id=$%d", idx)
		args = append(args, *parentID)
		idx++
	} else {
		q += " AND parent_id IS NULL"
	}
	if active != nil {
		q += fmt.Sprintf(" AND active=$%d", idx)
		args = append(args, *active)
		idx++
	}
	q += " ORDER BY name"
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListCategories", err)
	}
	defer rows.Close()
	var res []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.CompanyID, &c.Name, &c.Description, &c.ParentID, &c.Active,
			&c.CreatedBy, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, repoErr("ListCategories", err)
		}
		res = append(res, c)
	}
	return res, nil
}

// ---------------------------------------------------------------------------
// Courses
// ---------------------------------------------------------------------------

func (r *Repository) CreateCourse(ctx context.Context, c *Course) error {
	q := `INSERT INTO training_courses (id,company_id,code,name,category_id,short_description,description,
		objectives,difficulty,duration_minutes,modality,status,mandatory,passing_score,certificate_enabled,
		min_attendance_percentage,created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err := r.pool.Exec(ctx, q, c.ID, c.CompanyID, c.Code, c.Name, c.CategoryID, c.ShortDescription,
		c.Description, c.Objectives, c.Difficulty, c.DurationMinutes, c.Modality, c.Status,
		c.Mandatory, c.PassingScore, c.CertificateEnabled, c.MinAttendancePercentage, c.CreatedBy)
	return repoErr("CreateCourse", err)
}

func (r *Repository) UpdateCourse(ctx context.Context, c *Course) error {
	q := `UPDATE training_courses SET name=COALESCE($3,name), category_id=COALESCE($4,category_id),
		short_description=COALESCE($5,short_description), description=COALESCE($6,description),
		objectives=COALESCE($7,objectives), difficulty=COALESCE($8,difficulty),
		duration_minutes=COALESCE($9,duration_minutes), modality=COALESCE($10,modality),
		status=COALESCE($11,status), mandatory=COALESCE($12,mandatory),
		passing_score=COALESCE($13,passing_score), certificate_enabled=COALESCE($14,certificate_enabled),
		min_attendance_percentage=COALESCE($15,min_attendance_percentage), updated_at=NOW()
		WHERE id=$1 AND company_id=$2`
	_, err := r.pool.Exec(ctx, q, c.ID, c.CompanyID, c.Name, c.CategoryID, c.ShortDescription,
		c.Description, c.Objectives, c.Difficulty, c.DurationMinutes, c.Modality, c.Status,
		c.Mandatory, c.PassingScore, c.CertificateEnabled, c.MinAttendancePercentage)
	return repoErr("UpdateCourse", err)
}

func (r *Repository) GetCourse(ctx context.Context, companyID, id string) (*Course, error) {
	q := `SELECT id,company_id,code,name,category_id,short_description,description,objectives,
		difficulty,duration_minutes,modality,status,mandatory,passing_score,certificate_enabled,
		min_attendance_percentage,created_by,published_by,published_at,created_at,updated_at
		FROM training_courses WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	c := &Course{}
	err := row.Scan(&c.ID, &c.CompanyID, &c.Code, &c.Name, &c.CategoryID, &c.ShortDescription,
		&c.Description, &c.Objectives, &c.Difficulty, &c.DurationMinutes, &c.Modality, &c.Status,
		&c.Mandatory, &c.PassingScore, &c.CertificateEnabled, &c.MinAttendancePercentage,
		&c.CreatedBy, &c.PublishedBy, &c.PublishedAt, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetCourse", err)
	}
	return c, nil
}

func (r *Repository) ListCourses(ctx context.Context, companyID string, filter CourseFilter) ([]Course, int, error) {
	where := "WHERE company_id=$1"
	args := []any{companyID}
	idx := 2
	if filter.CategoryID != nil {
		where += fmt.Sprintf(" AND category_id=$%d", idx)
		args = append(args, *filter.CategoryID)
		idx++
	}
	if filter.Status != nil {
		where += fmt.Sprintf(" AND status=$%d", idx)
		args = append(args, *filter.Status)
		idx++
	}
	if filter.Modality != nil {
		where += fmt.Sprintf(" AND modality=$%d", idx)
		args = append(args, *filter.Modality)
		idx++
	}
	if filter.Mandatory != nil {
		where += fmt.Sprintf(" AND mandatory=$%d", idx)
		args = append(args, *filter.Mandatory)
		idx++
	}
	if filter.Difficulty != nil {
		where += fmt.Sprintf(" AND difficulty=$%d", idx)
		args = append(args, *filter.Difficulty)
		idx++
	}
	if filter.Search != nil {
		where += fmt.Sprintf(" AND (name ILIKE $%d OR code ILIKE $%d)", idx, idx)
		args = append(args, "%"+*filter.Search+"%")
		idx++
	}
	var total int
	cq := "SELECT COUNT(*) FROM training_courses " + where
	if err := r.pool.QueryRow(ctx, cq, args...).Scan(&total); err != nil {
		return nil, 0, repoErr("ListCourses.count", err)
	}
	q := `SELECT id,company_id,code,name,category_id,short_description,description,objectives,
		difficulty,duration_minutes,modality,status,mandatory,passing_score,certificate_enabled,
		min_attendance_percentage,created_by,published_by,published_at,created_at,updated_at
		FROM training_courses ` + where
	if filter.SortBy != "" {
		d := "ASC"
		if filter.SortDesc {
			d = "DESC"
		}
		q += fmt.Sprintf(" ORDER BY %s %s", filter.SortBy, d)
	} else {
		q += " ORDER BY created_at DESC"
	}
	if filter.Limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", idx)
		args = append(args, filter.Limit)
		idx++
	}
	if filter.Offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", idx)
		args = append(args, filter.Offset)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, repoErr("ListCourses", err)
	}
	defer rows.Close()
	var res []Course
	for rows.Next() {
		var c Course
		if err := rows.Scan(&c.ID, &c.CompanyID, &c.Code, &c.Name, &c.CategoryID, &c.ShortDescription,
			&c.Description, &c.Objectives, &c.Difficulty, &c.DurationMinutes, &c.Modality, &c.Status,
			&c.Mandatory, &c.PassingScore, &c.CertificateEnabled, &c.MinAttendancePercentage,
			&c.CreatedBy, &c.PublishedBy, &c.PublishedAt, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, 0, repoErr("ListCourses", err)
		}
		res = append(res, c)
	}
	return res, total, nil
}

func (r *Repository) PublishCourse(ctx context.Context, companyID, id, publishedBy string) error {
	q := `UPDATE training_courses SET status='published', published_by=$3, published_at=NOW(),
		updated_at=NOW() WHERE id=$1 AND company_id=$2`
	_, err := r.pool.Exec(ctx, q, id, companyID, publishedBy)
	return repoErr("PublishCourse", err)
}

// ---------------------------------------------------------------------------
// Course Versions
// ---------------------------------------------------------------------------

func (r *Repository) CreateVersion(ctx context.Context, v *CourseVersion) error {
	q := `INSERT INTO training_course_versions (id,course_id,version,description,created_by)
		VALUES ($1,$2,$3,$4,$5)`
	_, err := r.pool.Exec(ctx, q, v.ID, v.CourseID, v.Version, v.Description, v.CreatedBy)
	return repoErr("CreateVersion", err)
}

func (r *Repository) ListVersions(ctx context.Context, courseID string) ([]CourseVersion, error) {
	q := `SELECT id,course_id,version,description,is_published,created_by,created_at
		FROM training_course_versions WHERE course_id=$1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, q, courseID)
	if err != nil {
		return nil, repoErr("ListVersions", err)
	}
	defer rows.Close()
	var res []CourseVersion
	for rows.Next() {
		var v CourseVersion
		if err := rows.Scan(&v.ID, &v.CourseID, &v.Version, &v.Description, &v.IsPublished,
			&v.CreatedBy, &v.CreatedAt); err != nil {
			return nil, repoErr("ListVersions", err)
		}
		res = append(res, v)
	}
	return res, nil
}

func (r *Repository) PublishVersion(ctx context.Context, id string) error {
	q := `UPDATE training_course_versions SET is_published=true WHERE id=$1`
	_, err := r.pool.Exec(ctx, q, id)
	return repoErr("PublishVersion", err)
}

func (r *Repository) GetActiveVersion(ctx context.Context, courseID string) (*CourseVersion, error) {
	q := `SELECT id,course_id,version,description,is_published,created_by,created_at
		FROM training_course_versions WHERE course_id=$1 AND is_published=true ORDER BY created_at DESC LIMIT 1`
	row := r.pool.QueryRow(ctx, q, courseID)
	v := &CourseVersion{}
	err := row.Scan(&v.ID, &v.CourseID, &v.Version, &v.Description, &v.IsPublished,
		&v.CreatedBy, &v.CreatedAt)
	if err != nil {
		return nil, repoErr("GetActiveVersion", err)
	}
	return v, nil
}

// ---------------------------------------------------------------------------
// Course Contents
// ---------------------------------------------------------------------------

func (r *Repository) CreateContent(ctx context.Context, c *CourseContent) error {
	q := `INSERT INTO training_course_contents (id,course_version_id,title,description,content_type,
		external_url,duration_seconds,sort_order,required,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.pool.Exec(ctx, q, c.ID, c.CourseVersionID, c.Title, c.Description, c.ContentType,
		c.ExternalURL, c.DurationSeconds, c.SortOrder, c.Required, c.CreatedBy)
	return repoErr("CreateContent", err)
}

func (r *Repository) UpdateContent(ctx context.Context, c *CourseContent) error {
	q := `UPDATE training_course_contents SET title=COALESCE($3,title), description=COALESCE($4,description),
		content_type=COALESCE($5,content_type), external_url=COALESCE($6,external_url),
		duration_seconds=COALESCE($7,duration_seconds), sort_order=COALESCE($8,sort_order),
		required=COALESCE($9,required), published=COALESCE($10,published), updated_at=NOW()
		WHERE id=$1 AND course_version_id=$2`
	_, err := r.pool.Exec(ctx, q, c.ID, c.CourseVersionID, c.Title, c.Description, c.ContentType,
		c.ExternalURL, c.DurationSeconds, c.SortOrder, c.Required, c.Published)
	return repoErr("UpdateContent", err)
}

func (r *Repository) GetContent(ctx context.Context, id string) (*CourseContent, error) {
	q := `SELECT id,course_version_id,title,description,content_type,external_url,duration_seconds,
		sort_order,required,published,created_by,created_at,updated_at
		FROM training_course_contents WHERE id=$1`
	row := r.pool.QueryRow(ctx, q, id)
	c := &CourseContent{}
	err := row.Scan(&c.ID, &c.CourseVersionID, &c.Title, &c.Description, &c.ContentType,
		&c.ExternalURL, &c.DurationSeconds, &c.SortOrder, &c.Required, &c.Published,
		&c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetContent", err)
	}
	return c, nil
}

func (r *Repository) ListContents(ctx context.Context, versionID string, publishedOnly bool) ([]CourseContent, error) {
	q := `SELECT id,course_version_id,title,description,content_type,external_url,duration_seconds,
		sort_order,required,published,created_by,created_at,updated_at
		FROM training_course_contents WHERE course_version_id=$1`
	args := []any{versionID}
	if publishedOnly {
		q += " AND published=true"
	}
	q += " ORDER BY sort_order ASC"
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, repoErr("ListContents", err)
	}
	defer rows.Close()
	var res []CourseContent
	for rows.Next() {
		var c CourseContent
		if err := rows.Scan(&c.ID, &c.CourseVersionID, &c.Title, &c.Description, &c.ContentType,
			&c.ExternalURL, &c.DurationSeconds, &c.SortOrder, &c.Required, &c.Published,
			&c.CreatedBy, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, repoErr("ListContents", err)
		}
		res = append(res, c)
	}
	return res, nil
}

func (r *Repository) ReorderContents(ctx context.Context, versionID string, contentIDs []string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return repoErr("ReorderContents.begin", err)
	}
	defer tx.Rollback(ctx)
	for i, id := range contentIDs {
		if _, err := tx.Exec(ctx, `UPDATE training_course_contents SET sort_order=$1 WHERE id=$2 AND course_version_id=$3`, i+1, id, versionID); err != nil {
			return repoErr("ReorderContents", err)
		}
	}
	return repoErr("ReorderContents.commit", tx.Commit(ctx))
}

// ---------------------------------------------------------------------------
// Course Offerings
// ---------------------------------------------------------------------------

func (r *Repository) CreateOffering(ctx context.Context, o *CourseOffering) error {
	q := `INSERT INTO training_course_offerings (id,company_id,course_id,course_version_id,name,
		start_date,end_date,capacity,modality,location,meeting_url,instructor_id,provider_id,
		cost_amount,cost_currency,status,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err := r.pool.Exec(ctx, q, o.ID, o.CompanyID, o.CourseID, o.CourseVersionID, o.Name,
		o.StartDate, o.EndDate, o.Capacity, o.Modality, o.Location, o.MeetingURL,
		o.InstructorID, o.ProviderID, o.CostAmount, o.CostCurrency, o.Status, o.CreatedBy)
	return repoErr("CreateOffering", err)
}

func (r *Repository) UpdateOffering(ctx context.Context, o *CourseOffering) error {
	q := `UPDATE training_course_offerings SET name=COALESCE($3,name), start_date=COALESCE($4,start_date),
		end_date=COALESCE($5,end_date), capacity=COALESCE($6,capacity),
		modality=COALESCE($7,modality), location=COALESCE($8,location),
		meeting_url=COALESCE($9,meeting_url), instructor_id=COALESCE($10,instructor_id),
		provider_id=COALESCE($11,provider_id), cost_amount=COALESCE($12,cost_amount),
		cost_currency=COALESCE($13,cost_currency), status=COALESCE($14,status), updated_at=NOW()
		WHERE id=$1 AND company_id=$2`
	_, err := r.pool.Exec(ctx, q, o.ID, o.CompanyID, o.Name, o.StartDate, o.EndDate, o.Capacity,
		o.Modality, o.Location, o.MeetingURL, o.InstructorID, o.ProviderID, o.CostAmount,
		o.CostCurrency, o.Status)
	return repoErr("UpdateOffering", err)
}

func (r *Repository) GetOffering(ctx context.Context, companyID, id string) (*CourseOffering, error) {
	q := `SELECT id,company_id,course_id,course_version_id,name,start_date,end_date,capacity,
		modality,location,meeting_url,instructor_id,provider_id,enrolled_count,cost_amount,
		cost_currency,status,created_by,created_at,updated_at
		FROM training_course_offerings WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	o := &CourseOffering{}
	err := row.Scan(&o.ID, &o.CompanyID, &o.CourseID, &o.CourseVersionID, &o.Name,
		&o.StartDate, &o.EndDate, &o.Capacity, &o.Modality, &o.Location, &o.MeetingURL,
		&o.InstructorID, &o.ProviderID, &o.EnrolledCount, &o.CostAmount, &o.CostCurrency,
		&o.Status, &o.CreatedBy, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetOffering", err)
	}
	return o, nil
}

func (r *Repository) ListOfferings(ctx context.Context, companyID string, filter OfferingFilter) ([]CourseOffering, int, error) {
	where := "WHERE company_id=$1"
	args := []any{companyID}
	idx := 2
	if filter.CourseID != nil {
		where += fmt.Sprintf(" AND course_id=$%d", idx)
		args = append(args, *filter.CourseID)
		idx++
	}
	if filter.Status != nil {
		where += fmt.Sprintf(" AND status=$%d", idx)
		args = append(args, *filter.Status)
		idx++
	}
	if filter.InstructorID != nil {
		where += fmt.Sprintf(" AND instructor_id=$%d", idx)
		args = append(args, *filter.InstructorID)
		idx++
	}
	if filter.FromDate != nil {
		where += fmt.Sprintf(" AND start_date>=$%d", idx)
		args = append(args, *filter.FromDate)
		idx++
	}
	if filter.ToDate != nil {
		where += fmt.Sprintf(" AND start_date<=$%d", idx)
		args = append(args, *filter.ToDate)
		idx++
	}
	var total int
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM training_course_offerings "+where, args...).Scan(&total); err != nil {
		return nil, 0, repoErr("ListOfferings.count", err)
	}
	q := `SELECT id,company_id,course_id,course_version_id,name,start_date,end_date,capacity,
		modality,location,meeting_url,instructor_id,provider_id,enrolled_count,cost_amount,
		cost_currency,status,created_by,created_at,updated_at
		FROM training_course_offerings ` + where + " ORDER BY start_date ASC"
	if filter.Limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", idx)
		args = append(args, filter.Limit)
		idx++
	}
	if filter.Offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", idx)
		args = append(args, filter.Offset)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, repoErr("ListOfferings", err)
	}
	defer rows.Close()
	var res []CourseOffering
	for rows.Next() {
		var o CourseOffering
		if err := rows.Scan(&o.ID, &o.CompanyID, &o.CourseID, &o.CourseVersionID, &o.Name,
			&o.StartDate, &o.EndDate, &o.Capacity, &o.Modality, &o.Location, &o.MeetingURL,
			&o.InstructorID, &o.ProviderID, &o.EnrolledCount, &o.CostAmount, &o.CostCurrency,
			&o.Status, &o.CreatedBy, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, 0, repoErr("ListOfferings", err)
		}
		res = append(res, o)
	}
	return res, total, nil
}

// ---------------------------------------------------------------------------
// Sessions
// ---------------------------------------------------------------------------

func (r *Repository) CreateSession(ctx context.Context, s *OfferingSession) error {
	q := `INSERT INTO training_offering_sessions (id,offering_id,title,session_date,start_time,end_time,
		location,meeting_url,instructor_id,created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.pool.Exec(ctx, q, s.ID, s.OfferingID, s.Title, s.SessionDate, s.StartTime, s.EndTime,
		s.Location, s.MeetingURL, s.InstructorID, s.CreatedBy)
	return repoErr("CreateSession", err)
}

func (r *Repository) UpdateSession(ctx context.Context, s *OfferingSession) error {
	q := `UPDATE training_offering_sessions SET title=COALESCE($3,title), session_date=COALESCE($4,session_date),
		start_time=COALESCE($5,start_time), end_time=COALESCE($6,end_time),
		location=COALESCE($7,location), meeting_url=COALESCE($8,meeting_url),
		instructor_id=COALESCE($9,instructor_id), updated_at=NOW() WHERE id=$1 AND offering_id=$2`
	_, err := r.pool.Exec(ctx, q, s.ID, s.OfferingID, s.Title, s.SessionDate, s.StartTime, s.EndTime,
		s.Location, s.MeetingURL, s.InstructorID)
	return repoErr("UpdateSession", err)
}

func (r *Repository) ListSessionsByOffering(ctx context.Context, offeringID string) ([]OfferingSession, error) {
	q := `SELECT id,offering_id,title,session_date,start_time,end_time,location,meeting_url,
		instructor_id,created_by,created_at,updated_at
		FROM training_offering_sessions WHERE offering_id=$1 ORDER BY session_date ASC, start_time ASC`
	rows, err := r.pool.Query(ctx, q, offeringID)
	if err != nil {
		return nil, repoErr("ListSessionsByOffering", err)
	}
	defer rows.Close()
	var res []OfferingSession
	for rows.Next() {
		var s OfferingSession
		if err := rows.Scan(&s.ID, &s.OfferingID, &s.Title, &s.SessionDate, &s.StartTime, &s.EndTime,
			&s.Location, &s.MeetingURL, &s.InstructorID, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, repoErr("ListSessionsByOffering", err)
		}
		res = append(res, s)
	}
	return res, nil
}

// ---------------------------------------------------------------------------
// Enrollments
// ---------------------------------------------------------------------------

func (r *Repository) Enroll(ctx context.Context, e *Enrollment) error {
	q := `INSERT INTO training_enrollments (id,company_id,offering_id,employee_id,assignment_type,status,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.pool.Exec(ctx, q, e.ID, e.CompanyID, e.OfferingID, e.EmployeeID, e.AssignmentType, e.Status, e.CreatedBy)
	return repoErr("Enroll", err)
}

func (r *Repository) UpdateEnrollmentStatus(ctx context.Context, id, status string) error {
	q := `UPDATE training_enrollments SET status=$2, updated_at=NOW() WHERE id=$1`
	_, err := r.pool.Exec(ctx, q, id, status)
	return repoErr("UpdateEnrollmentStatus", err)
}

func (r *Repository) GetEnrollment(ctx context.Context, companyID, id string) (*Enrollment, error) {
	q := `SELECT id,company_id,offering_id,employee_id,assignment_type,status,progress_percentage,
		completed_at,certificate_url,certificate_issued_at,created_by,created_at,updated_at
		FROM training_enrollments WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	e := &Enrollment{}
	err := row.Scan(&e.ID, &e.CompanyID, &e.OfferingID, &e.EmployeeID, &e.AssignmentType, &e.Status,
		&e.ProgressPercentage, &e.CompletedAt, &e.CertificateURL, &e.CertificateIssuedAt,
		&e.CreatedBy, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetEnrollment", err)
	}
	return e, nil
}

func (r *Repository) ListEnrollments(ctx context.Context, companyID string, filter EnrollmentFilter) ([]Enrollment, int, error) {
	where := "WHERE e.company_id=$1"
	args := []any{companyID}
	idx := 2
	if filter.OfferingID != nil {
		where += fmt.Sprintf(" AND e.offering_id=$%d", idx)
		args = append(args, *filter.OfferingID)
		idx++
	}
	if filter.EmployeeID != nil {
		where += fmt.Sprintf(" AND e.employee_id=$%d", idx)
		args = append(args, *filter.EmployeeID)
		idx++
	}
	if filter.Status != nil {
		where += fmt.Sprintf(" AND e.status=$%d", idx)
		args = append(args, *filter.Status)
		idx++
	}
	var total int
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM training_enrollments e "+where, args...).Scan(&total); err != nil {
		return nil, 0, repoErr("ListEnrollments.count", err)
	}
	q := `SELECT e.id,e.company_id,e.offering_id,e.employee_id,e.assignment_type,e.status,
		e.progress_percentage,e.completed_at,e.certificate_url,e.certificate_issued_at,
		e.created_by,e.created_at,e.updated_at
		FROM training_enrollments e ` + where + " ORDER BY e.created_at DESC"
	if filter.Limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", idx)
		args = append(args, filter.Limit)
		idx++
	}
	if filter.Offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", idx)
		args = append(args, filter.Offset)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, repoErr("ListEnrollments", err)
	}
	defer rows.Close()
	var res []Enrollment
	for rows.Next() {
		var e Enrollment
		if err := rows.Scan(&e.ID, &e.CompanyID, &e.OfferingID, &e.EmployeeID, &e.AssignmentType,
			&e.Status, &e.ProgressPercentage, &e.CompletedAt, &e.CertificateURL, &e.CertificateIssuedAt,
			&e.CreatedBy, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, 0, repoErr("ListEnrollments", err)
		}
		res = append(res, e)
	}
	return res, total, nil
}

// ---------------------------------------------------------------------------
// Content Progress
// ---------------------------------------------------------------------------

func (r *Repository) UpsertContentProgress(ctx context.Context, p *ContentProgress) error {
	q := `INSERT INTO training_content_progress (id,enrollment_id,content_id,progress_percentage,time_spent_seconds,last_position) VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (enrollment_id,content_id) DO UPDATE SET progress_percentage=$4,time_spent_seconds=$5,last_position=$6,updated_at=NOW()`
	_, err := r.pool.Exec(ctx, q, p.ID, p.EnrollmentID, p.ContentID, p.ProgressPercentage, p.TimeSpentSeconds, p.LastPosition)
	return repoErr("UpsertContentProgress", err)
}

func (r *Repository) GetContentProgress(ctx context.Context, enrollmentID, contentID string) (*ContentProgress, error) {
	q := `SELECT id,enrollment_id,content_id,progress_percentage,time_spent_seconds,last_position,created_at,updated_at
		FROM training_content_progress WHERE enrollment_id=$1 AND content_id=$2`
	row := r.pool.QueryRow(ctx, q, enrollmentID, contentID)
	p := &ContentProgress{}
	err := row.Scan(&p.ID, &p.EnrollmentID, &p.ContentID, &p.ProgressPercentage, &p.TimeSpentSeconds, &p.LastPosition, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetContentProgress", err)
	}
	return p, nil
}

func (r *Repository) ListContentProgressByEnrollment(ctx context.Context, enrollmentID string) ([]ContentProgress, error) {
	q := `SELECT id,enrollment_id,content_id,progress_percentage,time_spent_seconds,last_position,created_at,updated_at
		FROM training_content_progress WHERE enrollment_id=$1 ORDER BY created_at`
	rows, err := r.pool.Query(ctx, q, enrollmentID)
	if err != nil {
		return nil, repoErr("ListContentProgressByEnrollment", err)
	}
	defer rows.Close()
	var res []ContentProgress
	for rows.Next() {
		var p ContentProgress
		if err := rows.Scan(&p.ID, &p.EnrollmentID, &p.ContentID, &p.ProgressPercentage, &p.TimeSpentSeconds, &p.LastPosition, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, repoErr("ListContentProgressByEnrollment", err)
		}
		res = append(res, p)
	}
	return res, nil
}

// ---------------------------------------------------------------------------
// Assignments (massive)
// ---------------------------------------------------------------------------

func (r *Repository) CreateAssignment(ctx context.Context, a *Assignment) error {
	q := `INSERT INTO training_assignments (id,company_id,course_id,assignee_type,assignee_id,assignment_type,due_date,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := r.pool.Exec(ctx, q, a.ID, a.CompanyID, a.CourseID, a.AssigneeType, a.AssigneeID, a.AssignmentType, a.DueDate, a.CreatedBy)
	return repoErr("CreateAssignment", err)
}

func (r *Repository) ProcessAssignment(ctx context.Context, a *Assignment) error {
	// Assign courses to all employees matching the criteria
	q := `INSERT INTO training_enrollments (id,company_id,offering_id,employee_id,assignment_type,status,created_by)
		SELECT $1,$2,$3,$4,$5,'assigned',$6 FROM employees e WHERE e.id=$4 AND e.company_id=$2
		ON CONFLICT (company_id,offering_id,employee_id) DO NOTHING`
	_, err := r.pool.Exec(ctx, q, uuid.New().String(), a.CompanyID, a.OfferingID, a.AssigneeID, a.CreatedBy)
	return repoErr("ProcessAssignment", err)
}

// ---------------------------------------------------------------------------
// Assignment Rules
// ---------------------------------------------------------------------------

func (r *Repository) CreateAssignmentRule(ctx context.Context, ar *AssignmentRule) error {
	q := `INSERT INTO training_assignment_rules (id,company_id,name,criteria_field,criteria_value,course_id,assignment_type,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := r.pool.Exec(ctx, q, ar.ID, ar.CompanyID, ar.Name, ar.CriteriaField, ar.CriteriaValue, ar.CourseID, ar.AssignmentType, ar.CreatedBy)
	return repoErr("CreateAssignmentRule", err)
}

func (r *Repository) ListAssignmentRules(ctx context.Context, companyID string) ([]AssignmentRule, error) {
	q := `SELECT id,company_id,name,criteria_field,criteria_value,course_id,assignment_type,active,last_run_at,created_by,created_at,updated_at
		FROM training_assignment_rules WHERE company_id=$1 ORDER BY name`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListAssignmentRules", err)
	}
	defer rows.Close()
	var res []AssignmentRule
	for rows.Next() {
		var ar AssignmentRule
		if err := rows.Scan(&ar.ID, &ar.CompanyID, &ar.Name, &ar.CriteriaField, &ar.CriteriaValue,
			&ar.CourseID, &ar.AssignmentType, &ar.Active, &ar.LastRunAt, &ar.CreatedBy, &ar.CreatedAt, &ar.UpdatedAt); err != nil {
			return nil, repoErr("ListAssignmentRules", err)
		}
		res = append(res, ar)
	}
	return res, nil
}

func (r *Repository) ExecuteAssignmentRule(ctx context.Context, companyID, ruleID, courseID string, employeeIDs []string, assignmentType, createdBy string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return repoErr("ExecuteAssignmentRule.begin", err)
	}
	defer tx.Rollback(ctx)
	offeringID := uuid.New().String()
	// Create a default offering
	q := `INSERT INTO training_course_offerings (id,company_id,course_id,name,status,created_by)
		VALUES ($1,$2,$3,'Auto-assigned: '||$3,'active',$4)`
	if _, err := tx.Exec(ctx, q, offeringID, companyID, courseID, createdBy); err != nil {
		return repoErr("ExecuteAssignmentRule.offering", err)
	}
	for _, empID := range employeeIDs {
		eID := uuid.New().String()
		if _, err := tx.Exec(ctx,
			`INSERT INTO training_enrollments (id,company_id,offering_id,employee_id,assignment_type,status,created_by) VALUES ($1,$2,$3,$4,$5,'assigned',$6) ON CONFLICT DO NOTHING`,
			eID, companyID, offeringID, empID, assignmentType, createdBy); err != nil {
			return repoErr("ExecuteAssignmentRule.enroll", err)
		}
	}
	if _, err := tx.Exec(ctx, `UPDATE training_assignment_rules SET last_run_at=NOW() WHERE id=$1`, ruleID); err != nil {
		return repoErr("ExecuteAssignmentRule.update", err)
	}
	return repoErr("ExecuteAssignmentRule.commit", tx.Commit(ctx))
}

// ---------------------------------------------------------------------------
// Assessments
// ---------------------------------------------------------------------------

func (r *Repository) CreateAssessment(ctx context.Context, a *Assessment) error {
	q := `INSERT INTO training_assessments (id,company_id,course_id,title,description,assessment_type,
		attempts_allowed,passing_score,time_limit_minutes,randomize_questions,show_results,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := r.pool.Exec(ctx, q, a.ID, a.CompanyID, a.CourseID, a.Title, a.Description, a.AssessmentType,
		a.AttemptsAllowed, a.PassingScore, a.TimeLimitMinutes, a.RandomizeQuestions, a.ShowResults, a.CreatedBy)
	return repoErr("CreateAssessment", err)
}

func (r *Repository) UpdateAssessment(ctx context.Context, a *Assessment) error {
	q := `UPDATE training_assessments SET title=COALESCE($3,title), description=COALESCE($4,description),
		attempts_allowed=COALESCE($5,attempts_allowed), passing_score=COALESCE($6,passing_score),
		time_limit_minutes=COALESCE($7,time_limit_minutes), randomize_questions=COALESCE($8,randomize_questions),
		show_results=COALESCE($9,show_results), status=COALESCE($10,status), updated_at=NOW()
		WHERE id=$1 AND company_id=$2`
	_, err := r.pool.Exec(ctx, q, a.ID, a.CompanyID, a.Title, a.Description, a.AttemptsAllowed,
		a.PassingScore, a.TimeLimitMinutes, a.RandomizeQuestions, a.ShowResults, a.Status)
	return repoErr("UpdateAssessment", err)
}

func (r *Repository) GetAssessment(ctx context.Context, companyID, id string) (*Assessment, error) {
	q := `SELECT id,company_id,course_id,title,description,assessment_type,attempts_allowed,passing_score,
		time_limit_minutes,randomize_questions,show_results,status,created_by,created_at,updated_at
		FROM training_assessments WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	a := &Assessment{}
	err := row.Scan(&a.ID, &a.CompanyID, &a.CourseID, &a.Title, &a.Description, &a.AssessmentType,
		&a.AttemptsAllowed, &a.PassingScore, &a.TimeLimitMinutes, &a.RandomizeQuestions, &a.ShowResults,
		&a.Status, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetAssessment", err)
	}
	return a, nil
}

func (r *Repository) ListAssessments(ctx context.Context, companyID, courseID string) ([]Assessment, error) {
	q := `SELECT id,company_id,course_id,title,description,assessment_type,attempts_allowed,passing_score,
		time_limit_minutes,randomize_questions,show_results,status,created_by,created_at,updated_at
		FROM training_assessments WHERE company_id=$1 AND course_id=$2 ORDER BY created_at`
	rows, err := r.pool.Query(ctx, q, companyID, courseID)
	if err != nil {
		return nil, repoErr("ListAssessments", err)
	}
	defer rows.Close()
	var res []Assessment
	for rows.Next() {
		var a Assessment
		if err := rows.Scan(&a.ID, &a.CompanyID, &a.CourseID, &a.Title, &a.Description, &a.AssessmentType,
			&a.AttemptsAllowed, &a.PassingScore, &a.TimeLimitMinutes, &a.RandomizeQuestions, &a.ShowResults,
			&a.Status, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, repoErr("ListAssessments", err)
		}
		res = append(res, a)
	}
	return res, nil
}

// ---------------------------------------------------------------------------
// Questions
// ---------------------------------------------------------------------------

func (r *Repository) CreateQuestion(ctx context.Context, q *Question) error {
	qry := `INSERT INTO training_assessment_questions (id,assessment_id,question,question_type,points,sort_order)
		VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.pool.Exec(ctx, qry, q.ID, q.AssessmentID, q.Question, q.QuestionType, q.Points, q.SortOrder)
	return repoErr("CreateQuestion", err)
}

func (r *Repository) CreateOption(ctx context.Context, o *QuestionOption) error {
	qry := `INSERT INTO training_question_options (id,question_id,option_text,is_correct,sort_order)
		VALUES ($1,$2,$3,$4,$5)`
	_, err := r.pool.Exec(ctx, qry, o.ID, o.QuestionID, o.OptionText, o.IsCorrect, o.SortOrder)
	return repoErr("CreateOption", err)
}

func (r *Repository) ListQuestions(ctx context.Context, assessmentID string) ([]Question, error) {
	qry := `SELECT id,assessment_id,question,question_type,points,sort_order,created_at,updated_at
		FROM training_assessment_questions WHERE assessment_id=$1 ORDER BY sort_order`
	rows, err := r.pool.Query(ctx, qry, assessmentID)
	if err != nil {
		return nil, repoErr("ListQuestions", err)
	}
	defer rows.Close()
	var res []Question
	for rows.Next() {
		var q Question
		if err := rows.Scan(&q.ID, &q.AssessmentID, &q.Question, &q.QuestionType, &q.Points, &q.SortOrder, &q.CreatedAt, &q.UpdatedAt); err != nil {
			return nil, repoErr("ListQuestions", err)
		}
		res = append(res, q)
	}
	return res, nil
}

func (r *Repository) ListOptionsByQuestion(ctx context.Context, questionID string) ([]QuestionOption, error) {
	qry := `SELECT id,question_id,option_text,is_correct,sort_order,created_at
		FROM training_question_options WHERE question_id=$1 ORDER BY sort_order`
	rows, err := r.pool.Query(ctx, qry, questionID)
	if err != nil {
		return nil, repoErr("ListOptionsByQuestion", err)
	}
	defer rows.Close()
	var res []QuestionOption
	for rows.Next() {
		var o QuestionOption
		if err := rows.Scan(&o.ID, &o.QuestionID, &o.OptionText, &o.IsCorrect, &o.SortOrder, &o.CreatedAt); err != nil {
			return nil, repoErr("ListOptionsByQuestion", err)
		}
		res = append(res, o)
	}
	return res, nil
}

// ---------------------------------------------------------------------------
// Attempts
// ---------------------------------------------------------------------------

func (r *Repository) CreateAttempt(ctx context.Context, a *Attempt) error {
	qry := `INSERT INTO training_assessment_attempts (id,enrollment_id,assessment_id,attempt_number,started_at)
		VALUES ($1,$2,$3,$4,$5)`
	_, err := r.pool.Exec(ctx, qry, a.ID, a.EnrollmentID, a.AssessmentID, a.AttemptNumber, a.StartedAt)
	return repoErr("CreateAttempt", err)
}

func (r *Repository) SubmitAttempt(ctx context.Context, id string, score, totalPoints float64, status string) error {
	qry := `UPDATE training_assessment_attempts SET score=$2,total_points=$3,status=$4,completed_at=NOW()
		WHERE id=$1`
	_, err := r.pool.Exec(ctx, qry, id, score, totalPoints, status)
	return repoErr("SubmitAttempt", err)
}

func (r *Repository) GetAttempt(ctx context.Context, id string) (*Attempt, error) {
	qry := `SELECT id,enrollment_id,assessment_id,attempt_number,score,total_points,status,started_at,completed_at
		FROM training_assessment_attempts WHERE id=$1`
	row := r.pool.QueryRow(ctx, qry, id)
	a := &Attempt{}
	err := row.Scan(&a.ID, &a.EnrollmentID, &a.AssessmentID, &a.AttemptNumber, &a.Score, &a.TotalPoints,
		&a.Status, &a.StartedAt, &a.CompletedAt)
	if err != nil {
		return nil, repoErr("GetAttempt", err)
	}
	return a, nil
}

func (r *Repository) ListAttemptsByEnrollment(ctx context.Context, enrollmentID string) ([]Attempt, error) {
	qry := `SELECT id,enrollment_id,assessment_id,attempt_number,score,total_points,status,started_at,completed_at
		FROM training_assessment_attempts WHERE enrollment_id=$1 ORDER BY attempt_number`
	rows, err := r.pool.Query(ctx, qry, enrollmentID)
	if err != nil {
		return nil, repoErr("ListAttemptsByEnrollment", err)
	}
	defer rows.Close()
	var res []Attempt
	for rows.Next() {
		var a Attempt
		if err := rows.Scan(&a.ID, &a.EnrollmentID, &a.AssessmentID, &a.AttemptNumber, &a.Score, &a.TotalPoints,
			&a.Status, &a.StartedAt, &a.CompletedAt); err != nil {
			return nil, repoErr("ListAttemptsByEnrollment", err)
		}
		res = append(res, a)
	}
	return res, nil
}

func (r *Repository) CreateAnswer(ctx context.Context, a *Answer) error {
	qry := `INSERT INTO training_attempt_answers (id,attempt_id,question_id,selected_option_id,text_answer,numeric_answer,is_correct,points_earned)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := r.pool.Exec(ctx, qry, a.ID, a.AttemptID, a.QuestionID, a.SelectedOptionID, a.TextAnswer, a.NumericAnswer, a.IsCorrect, a.PointsEarned)
	return repoErr("CreateAnswer", err)
}

func (r *Repository) GetLatestAttempt(ctx context.Context, enrollmentID, assessmentID string) (*Attempt, error) {
	qry := `SELECT id,enrollment_id,assessment_id,attempt_number,score,total_points,status,started_at,completed_at
		FROM training_assessment_attempts WHERE enrollment_id=$1 AND assessment_id=$2 ORDER BY attempt_number DESC LIMIT 1`
	row := r.pool.QueryRow(ctx, qry, enrollmentID, assessmentID)
	a := &Attempt{}
	err := row.Scan(&a.ID, &a.EnrollmentID, &a.AssessmentID, &a.AttemptNumber, &a.Score, &a.TotalPoints,
		&a.Status, &a.StartedAt, &a.CompletedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, repoErr("GetLatestAttempt", err)
	}
	return a, nil
}

// ---------------------------------------------------------------------------
// Instructors
// ---------------------------------------------------------------------------

func (r *Repository) CreateInstructor(ctx context.Context, i *Instructor) error {
	q := `INSERT INTO training_instructors (id,company_id,employee_id,instructor_type,name,email,phone,specialization,bio,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.pool.Exec(ctx, q, i.ID, i.CompanyID, i.EmployeeID, i.InstructorType, i.Name, i.Email, i.Phone, i.Specialization, i.Bio, i.CreatedBy)
	return repoErr("CreateInstructor", err)
}

func (r *Repository) GetInstructor(ctx context.Context, companyID, id string) (*Instructor, error) {
	q := `SELECT id,company_id,employee_id,instructor_type,name,email,phone,specialization,bio,active,rating,created_by,created_at,updated_at
		FROM training_instructors WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	i := &Instructor{}
	err := row.Scan(&i.ID, &i.CompanyID, &i.EmployeeID, &i.InstructorType, &i.Name, &i.Email, &i.Phone,
		&i.Specialization, &i.Bio, &i.Active, &i.Rating, &i.CreatedBy, &i.CreatedAt, &i.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetInstructor", err)
	}
	return i, nil
}

func (r *Repository) ListInstructors(ctx context.Context, companyID string) ([]Instructor, error) {
	q := `SELECT id,company_id,employee_id,instructor_type,name,email,phone,specialization,bio,active,rating,created_by,created_at,updated_at
		FROM training_instructors WHERE company_id=$1 ORDER BY name`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListInstructors", err)
	}
	defer rows.Close()
	var res []Instructor
	for rows.Next() {
		var i Instructor
		if err := rows.Scan(&i.ID, &i.CompanyID, &i.EmployeeID, &i.InstructorType, &i.Name, &i.Email, &i.Phone,
			&i.Specialization, &i.Bio, &i.Active, &i.Rating, &i.CreatedBy, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, repoErr("ListInstructors", err)
		}
		res = append(res, i)
	}
	return res, nil
}

// ---------------------------------------------------------------------------
// Providers
// ---------------------------------------------------------------------------

func (r *Repository) CreateProvider(ctx context.Context, p *TrainingProvider) error {
	q := `INSERT INTO training_providers (id,company_id,name,tax_id,email,phone,website,contact_name,notes,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.pool.Exec(ctx, q, p.ID, p.CompanyID, p.Name, p.TaxID, p.Email, p.Phone, p.Website, p.ContactName, p.Notes, p.CreatedBy)
	return repoErr("CreateProvider", err)
}

func (r *Repository) ListProviders(ctx context.Context, companyID string) ([]TrainingProvider, error) {
	q := `SELECT id,company_id,name,tax_id,email,phone,website,contact_name,notes,active,created_by,created_at,updated_at
		FROM training_providers WHERE company_id=$1 ORDER BY name`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListProviders", err)
	}
	defer rows.Close()
	var res []TrainingProvider
	for rows.Next() {
		var p TrainingProvider
		if err := rows.Scan(&p.ID, &p.CompanyID, &p.Name, &p.TaxID, &p.Email, &p.Phone, &p.Website,
			&p.ContactName, &p.Notes, &p.Active, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, repoErr("ListProviders", err)
		}
		res = append(res, p)
	}
	return res, nil
}

// ---------------------------------------------------------------------------
// Competencies
// ---------------------------------------------------------------------------

func (r *Repository) CreateCompetency(ctx context.Context, c *Competency) error {
	q := `INSERT INTO training_competencies (id,company_id,name,description,competency_type,created_by)
		VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.pool.Exec(ctx, q, c.ID, c.CompanyID, c.Name, c.Description, c.CompetencyType, c.CreatedBy)
	return repoErr("CreateCompetency", err)
}

func (r *Repository) CreateCompetencyLevel(ctx context.Context, l *CompetencyLevel) error {
	q := `INSERT INTO training_competency_levels (id,competency_id,level,label,description,created_by)
		VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.pool.Exec(ctx, q, l.ID, l.CompetencyID, l.Level, l.Label, l.Description, l.CreatedBy)
	return repoErr("CreateCompetencyLevel", err)
}

func (r *Repository) GetCompetency(ctx context.Context, companyID, id string) (*Competency, error) {
	q := `SELECT id,company_id,name,description,competency_type,created_by,created_at,updated_at
		FROM training_competencies WHERE id=$1 AND company_id=$2`
	row := r.pool.QueryRow(ctx, q, id, companyID)
	c := &Competency{}
	err := row.Scan(&c.ID, &c.CompanyID, &c.Name, &c.Description, &c.CompetencyType, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetCompetency", err)
	}
	return c, nil
}

func (r *Repository) ListCompetencyLevels(ctx context.Context, competencyID string) ([]CompetencyLevel, error) {
	q := `SELECT id,competency_id,level,label,description,created_by,created_at
		FROM training_competency_levels WHERE competency_id=$1 ORDER BY level`
	rows, err := r.pool.Query(ctx, q, competencyID)
	if err != nil {
		return nil, repoErr("ListCompetencyLevels", err)
	}
	defer rows.Close()
	var res []CompetencyLevel
	for rows.Next() {
		var l CompetencyLevel
		if err := rows.Scan(&l.ID, &l.CompetencyID, &l.Level, &l.Label, &l.Description, &l.CreatedBy, &l.CreatedAt); err != nil {
			return nil, repoErr("ListCompetencyLevels", err)
		}
		res = append(res, l)
	}
	return res, nil
}

func (r *Repository) ListCompetencies(ctx context.Context, companyID string) ([]Competency, error) {
	q := `SELECT id,company_id,name,description,competency_type,created_by,created_at,updated_at
		FROM training_competencies WHERE company_id=$1 ORDER BY name`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListCompetencies", err)
	}
	defer rows.Close()
	var res []Competency
	for rows.Next() {
		var c Competency
		if err := rows.Scan(&c.ID, &c.CompanyID, &c.Name, &c.Description, &c.CompetencyType, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, repoErr("ListCompetencies", err)
		}
		res = append(res, c)
	}
	return res, nil
}

// ---------------------------------------------------------------------------
// Employee Competencies
// ---------------------------------------------------------------------------

func (r *Repository) UpsertEmployeeCompetency(ctx context.Context, ec *EmployeeCompetency) error {
	q := `INSERT INTO training_employee_competencies (id,company_id,employee_id,competency_id,level,source,verified,verified_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		ON CONFLICT (company_id,employee_id,competency_id) DO UPDATE SET level=$5,source=$6,verified=$7,verified_by=$8,updated_at=NOW()`
	_, err := r.pool.Exec(ctx, q, ec.ID, ec.CompanyID, ec.EmployeeID, ec.CompetencyID, ec.Level, ec.Source, ec.Verified, ec.VerifiedBy)
	return repoErr("UpsertEmployeeCompetency", err)
}

func (r *Repository) ListEmployeeCompetencies(ctx context.Context, companyID, employeeID string) ([]EmployeeCompetency, error) {
	q := `SELECT id,company_id,employee_id,competency_id,level,source,verified,verified_by,created_at,updated_at
		FROM training_employee_competencies WHERE company_id=$1 AND employee_id=$2`
	rows, err := r.pool.Query(ctx, q, companyID, employeeID)
	if err != nil {
		return nil, repoErr("ListEmployeeCompetencies", err)
	}
	defer rows.Close()
	var res []EmployeeCompetency
	for rows.Next() {
		var ec EmployeeCompetency
		if err := rows.Scan(&ec.ID, &ec.CompanyID, &ec.EmployeeID, &ec.CompetencyID, &ec.Level,
			&ec.Source, &ec.Verified, &ec.VerifiedBy, &ec.CreatedAt, &ec.UpdatedAt); err != nil {
			return nil, repoErr("ListEmployeeCompetencies", err)
		}
		res = append(res, ec)
	}
	return res, nil
}

// ---------------------------------------------------------------------------
// Course-Competency mapping
// ---------------------------------------------------------------------------

func (r *Repository) AddCourseCompetency(ctx context.Context, cc *CourseCompetency) error {
	q := `INSERT INTO training_course_competencies (id,course_id,competency_id,expected_level,weight)
		VALUES ($1,$2,$3,$4,$5)`
	_, err := r.pool.Exec(ctx, q, cc.ID, cc.CourseID, cc.CompetencyID, cc.ExpectedLevel, cc.Weight)
	return repoErr("AddCourseCompetency", err)
}

func (r *Repository) ListCourseCompetencies(ctx context.Context, courseID string) ([]CourseCompetency, error) {
	q := `SELECT id,course_id,competency_id,expected_level,weight FROM training_course_competencies WHERE course_id=$1`
	rows, err := r.pool.Query(ctx, q, courseID)
	if err != nil {
		return nil, repoErr("ListCourseCompetencies", err)
	}
	defer rows.Close()
	var res []CourseCompetency
	for rows.Next() {
		var cc CourseCompetency
		if err := rows.Scan(&cc.ID, &cc.CourseID, &cc.CompetencyID, &cc.ExpectedLevel, &cc.Weight); err != nil {
			return nil, repoErr("ListCourseCompetencies", err)
		}
		res = append(res, cc)
	}
	return res, nil
}

// ---------------------------------------------------------------------------
// Training Needs
// ---------------------------------------------------------------------------

func (r *Repository) CreateTrainingNeed(ctx context.Context, n *TrainingNeed) error {
	q := `INSERT INTO training_needs (id,company_id,employee_id,competency_id,title,description,priority,source,source_id,status,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.pool.Exec(ctx, q, n.ID, n.CompanyID, n.EmployeeID, n.CompetencyID, n.Title, n.Description,
		n.Priority, n.Source, n.SourceID, n.Status, n.CreatedBy)
	return repoErr("CreateTrainingNeed", err)
}

func (r *Repository) ListTrainingNeeds(ctx context.Context, companyID string) ([]TrainingNeed, error) {
	q := `SELECT id,company_id,employee_id,competency_id,title,description,priority,source,source_id,status,created_by,created_at,updated_at
		FROM training_needs WHERE company_id=$1 ORDER BY priority DESC, created_at DESC`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListTrainingNeeds", err)
	}
	defer rows.Close()
	var res []TrainingNeed
	for rows.Next() {
		var n TrainingNeed
		if err := rows.Scan(&n.ID, &n.CompanyID, &n.EmployeeID, &n.CompetencyID, &n.Title, &n.Description,
			&n.Priority, &n.Source, &n.SourceID, &n.Status, &n.CreatedBy, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, repoErr("ListTrainingNeeds", err)
		}
		res = append(res, n)
	}
	return res, nil
}

// ---------------------------------------------------------------------------
// Training Plans
// ---------------------------------------------------------------------------

func (r *Repository) CreatePlan(ctx context.Context, p *TrainingPlan) error {
	q := `INSERT INTO training_plans (id,company_id,employee_id,name,description,objectives,period_start,period_end,
		budget_amount,budget_currency,status,created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := r.pool.Exec(ctx, q, p.ID, p.CompanyID, p.EmployeeID, p.Name, p.Description, p.Objectives,
		p.PeriodStart, p.PeriodEnd, p.BudgetAmount, p.BudgetCurrency, p.Status, p.CreatedBy)
	return repoErr("CreatePlan", err)
}

func (r *Repository) ListPlans(ctx context.Context, companyID string) ([]TrainingPlan, error) {
	q := `SELECT id,company_id,employee_id,name,description,objectives,period_start,period_end,
		budget_amount,budget_currency,status,created_by,created_at,updated_at
		FROM training_plans WHERE company_id=$1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListPlans", err)
	}
	defer rows.Close()
	var res []TrainingPlan
	for rows.Next() {
		var p TrainingPlan
		if err := rows.Scan(&p.ID, &p.CompanyID, &p.EmployeeID, &p.Name, &p.Description, &p.Objectives,
			&p.PeriodStart, &p.PeriodEnd, &p.BudgetAmount, &p.BudgetCurrency, &p.Status,
			&p.CreatedBy, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, repoErr("ListPlans", err)
		}
		res = append(res, p)
	}
	return res, nil
}

// ---------------------------------------------------------------------------
// Learning Paths
// ---------------------------------------------------------------------------

func (r *Repository) CreateLearningPath(ctx context.Context, lp *LearningPath) error {
	q := `INSERT INTO training_learning_paths (id,company_id,name,description,objectives,duration_days,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.pool.Exec(ctx, q, lp.ID, lp.CompanyID, lp.Name, lp.Description, lp.Objectives, lp.DurationDays, lp.CreatedBy)
	return repoErr("CreateLearningPath", err)
}

func (r *Repository) ListLearningPaths(ctx context.Context, companyID string) ([]LearningPath, error) {
	q := `SELECT id,company_id,name,description,objectives,duration_days,status,created_by,created_at,updated_at
		FROM training_learning_paths WHERE company_id=$1 ORDER BY name`
	rows, err := r.pool.Query(ctx, q, companyID)
	if err != nil {
		return nil, repoErr("ListLearningPaths", err)
	}
	defer rows.Close()
	var res []LearningPath
	for rows.Next() {
		var lp LearningPath
		if err := rows.Scan(&lp.ID, &lp.CompanyID, &lp.Name, &lp.Description, &lp.Objectives,
			&lp.DurationDays, &lp.Status, &lp.CreatedBy, &lp.CreatedAt, &lp.UpdatedAt); err != nil {
			return nil, repoErr("ListLearningPaths", err)
		}
		res = append(res, lp)
	}
	return res, nil
}

func (r *Repository) AddPathCourse(ctx context.Context, pathID, courseID string, sortOrder int) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO training_learning_path_courses (id,path_id,course_id,sort_order) VALUES ($1,$2,$3,$4)`,
		uuid.New().String(), pathID, courseID, sortOrder)
	return repoErr("AddPathCourse", err)
}

func (r *Repository) ListPathCourses(ctx context.Context, pathID string) ([]string, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT course_id FROM training_learning_path_courses WHERE path_id=$1 ORDER BY sort_order`, pathID)
	if err != nil {
		return nil, repoErr("ListPathCourses", err)
	}
	defer rows.Close()
	var res []string
	for rows.Next() {
		var cid string
		if err := rows.Scan(&cid); err != nil {
			return nil, repoErr("ListPathCourses", err)
		}
		res = append(res, cid)
	}
	return res, nil
}

// ---------------------------------------------------------------------------
// Feedback
// ---------------------------------------------------------------------------

func (r *Repository) CreateFeedback(ctx context.Context, f *Feedback) error {
	q := `INSERT INTO training_feedback (id,enrollment_id,instructor_rating,content_rating,
		organization_rating,platform_rating,overall_rating,comments,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := r.pool.Exec(ctx, q, f.ID, f.EnrollmentID, f.InstructorRating, f.ContentRating,
		f.OrganizationRating, f.PlatformRating, f.OverallRating, f.Comments, f.CreatedBy)
	return repoErr("CreateFeedback", err)
}

func (r *Repository) GetFeedbackByEnrollment(ctx context.Context, enrollmentID string) (*Feedback, error) {
	q := `SELECT id,enrollment_id,instructor_rating,content_rating,organization_rating,platform_rating,
		overall_rating,comments,created_by,created_at
		FROM training_feedback WHERE enrollment_id=$1`
	row := r.pool.QueryRow(ctx, q, enrollmentID)
	f := &Feedback{}
	err := row.Scan(&f.ID, &f.EnrollmentID, &f.InstructorRating, &f.ContentRating, &f.OrganizationRating,
		&f.PlatformRating, &f.OverallRating, &f.Comments, &f.CreatedBy, &f.CreatedAt)
	if err != nil {
		return nil, repoErr("GetFeedbackByEnrollment", err)
	}
	return f, nil
}

// ---------------------------------------------------------------------------
// Attendance
// ---------------------------------------------------------------------------

func (r *Repository) CreateAttendance(ctx context.Context, a *Attendance) error {
	q := `INSERT INTO training_attendance (id,enrollment_id,session_id,status,check_in,check_out,notes,created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := r.pool.Exec(ctx, q, a.ID, a.EnrollmentID, a.SessionID, a.Status, a.CheckIn, a.CheckOut, a.Notes, a.CreatedBy)
	return repoErr("CreateAttendance", err)
}

func (r *Repository) GetAttendanceBySession(ctx context.Context, enrollmentID, sessionID string) (*Attendance, error) {
	q := `SELECT id,enrollment_id,session_id,status,check_in,check_out,notes,created_by,created_at,updated_at
		FROM training_attendance WHERE enrollment_id=$1 AND session_id=$2`
	row := r.pool.QueryRow(ctx, q, enrollmentID, sessionID)
	a := &Attendance{}
	err := row.Scan(&a.ID, &a.EnrollmentID, &a.SessionID, &a.Status, &a.CheckIn, &a.CheckOut,
		&a.Notes, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, repoErr("GetAttendanceBySession", err)
	}
	return a, nil
}

func (r *Repository) ListAttendanceByEnrollment(ctx context.Context, enrollmentID string) ([]Attendance, error) {
	q := `SELECT id,enrollment_id,session_id,status,check_in,check_out,notes,created_by,created_at,updated_at
		FROM training_attendance WHERE enrollment_id=$1 ORDER BY created_at`
	rows, err := r.pool.Query(ctx, q, enrollmentID)
	if err != nil {
		return nil, repoErr("ListAttendanceByEnrollment", err)
	}
	defer rows.Close()
	var res []Attendance
	for rows.Next() {
		var a Attendance
		if err := rows.Scan(&a.ID, &a.EnrollmentID, &a.SessionID, &a.Status, &a.CheckIn, &a.CheckOut,
			&a.Notes, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, repoErr("ListAttendanceByEnrollment", err)
		}
		res = append(res, a)
	}
	return res, nil
}

// ---------------------------------------------------------------------------
// Certificates
// ---------------------------------------------------------------------------

func (r *Repository) UpdateCertificate(ctx context.Context, enrollmentID, url string) error {
	q := `UPDATE training_enrollments SET certificate_url=$2, certificate_issued_at=NOW(), updated_at=NOW() WHERE id=$1`
	_, err := r.pool.Exec(ctx, q, enrollmentID, url)
	return repoErr("UpdateCertificate", err)
}

func (r *Repository) ListCertificatesByEmployee(ctx context.Context, companyID, employeeID string) ([]Enrollment, error) {
	q := `SELECT id,company_id,offering_id,employee_id,assignment_type,status,progress_percentage,
		completed_at,certificate_url,certificate_issued_at,created_by,created_at,updated_at
		FROM training_enrollments WHERE company_id=$1 AND employee_id=$2 AND certificate_url IS NOT NULL
		ORDER BY certificate_issued_at DESC`
	rows, err := r.pool.Query(ctx, q, companyID, employeeID)
	if err != nil {
		return nil, repoErr("ListCertificatesByEmployee", err)
	}
	defer rows.Close()
	var res []Enrollment
	for rows.Next() {
		var e Enrollment
		if err := rows.Scan(&e.ID, &e.CompanyID, &e.OfferingID, &e.EmployeeID, &e.AssignmentType,
			&e.Status, &e.ProgressPercentage, &e.CompletedAt, &e.CertificateURL, &e.CertificateIssuedAt,
			&e.CreatedBy, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, repoErr("ListCertificatesByEmployee", err)
		}
		res = append(res, e)
	}
	return res, nil
}

// ---------------------------------------------------------------------------
// Dashboard stats
// ---------------------------------------------------------------------------

func (r *Repository) GetDashboardStats(ctx context.Context, companyID string) (*DashboardStats, error) {
	ds := &DashboardStats{}
	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM training_courses WHERE company_id=$1 AND status='published'`, companyID).Scan(&ds.TotalCourses); err != nil {
		return nil, repoErr("GetDashboardStats.courses", err)
	}
	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM training_enrollments WHERE company_id=$1`, companyID).Scan(&ds.TotalEnrollments); err != nil {
		return nil, repoErr("GetDashboardStats.enrollments", err)
	}
	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM training_enrollments WHERE company_id=$1 AND status='completed'`, companyID).Scan(&ds.CompletedEnrollments); err != nil {
		return nil, repoErr("GetDashboardStats.completed", err)
	}
	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM training_enrollments WHERE company_id=$1 AND status='in_progress'`, companyID).Scan(&ds.InProgressEnrollments); err != nil {
		return nil, repoErr("GetDashboardStats.in_progress", err)
	}
	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM training_offerings WHERE company_id=$1 AND status='active'`, companyID).Scan(&ds.ActiveOfferings); err != nil {
		return nil, repoErr("GetDashboardStats.offerings", err)
	}
	if err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(ROUND(AVG(overall_rating),2),0) FROM training_feedback f
		JOIN training_enrollments e ON e.id=f.enrollment_id WHERE e.company_id=$1`, companyID).Scan(&ds.AverageRating); err != nil {
		return nil, repoErr("GetDashboardStats.rating", err)
	}
	return ds, nil
}

// ---------------------------------------------------------------------------
// Employee stats
// ---------------------------------------------------------------------------

func (r *Repository) GetEmployeeStats(ctx context.Context, companyID, employeeID string) (*EmployeeStats, error) {
	es := &EmployeeStats{}
	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM training_enrollments WHERE company_id=$1 AND employee_id=$2`, companyID, employeeID).Scan(&es.TotalEnrollments); err != nil {
		return nil, repoErr("GetEmployeeStats.enrollments", err)
	}
	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM training_enrollments WHERE company_id=$1 AND employee_id=$2 AND status='completed'`, companyID, employeeID).Scan(&es.CompletedCourses); err != nil {
		return nil, repoErr("GetEmployeeStats.completed", err)
	}
	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM training_enrollments WHERE company_id=$1 AND employee_id=$2 AND status='in_progress'`, companyID, employeeID).Scan(&es.InProgress); err != nil {
		return nil, repoErr("GetEmployeeStats.in_progress", err)
	}
	if err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(duration_minutes),0) FROM training_enrollments e
		JOIN training_course_offerings o ON o.id=e.offering_id
		JOIN training_courses c ON c.id=o.course_id
		WHERE e.company_id=$1 AND e.employee_id=$2 AND e.status='completed'`, companyID, employeeID).Scan(&es.TotalTrainingHours); err != nil {
		return nil, repoErr("GetEmployeeStats.hours", err)
	}
	es.TotalTrainingHours = es.TotalTrainingHours / 60
	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM training_enrollments WHERE company_id=$1 AND employee_id=$2 AND certificate_url IS NOT NULL`, companyID, employeeID).Scan(&es.CertificatesCount); err != nil {
		return nil, repoErr("GetEmployeeStats.certificates", err)
	}
	return es, nil
}

// ---------------------------------------------------------------------------
// AI Recommendations (materialized)
// ---------------------------------------------------------------------------

func (r *Repository) SaveRecommendations(ctx context.Context, employeeID string, recs []AIRecommendation) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return repoErr("SaveRecommendations.begin", err)
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, `DELETE FROM training_ai_recommendations WHERE employee_id=$1`, employeeID); err != nil {
		return repoErr("SaveRecommendations.delete", err)
	}
	for _, rec := range recs {
		id := uuid.New().String()
		if _, err := tx.Exec(ctx,
			`INSERT INTO training_ai_recommendations (id,employee_id,course_id,reason,competency_id,expected_level,priority)
			VALUES ($1,$2,$3,$4,$5,$6,$7)`,
			id, employeeID, rec.CourseID, rec.Reason, rec.CompetencyID, rec.ExpectedLevel, rec.Priority); err != nil {
			return repoErr("SaveRecommendations.insert", err)
		}
	}
	return repoErr("SaveRecommendations.commit", tx.Commit(ctx))
}

func (r *Repository) GetRecommendations(ctx context.Context, employeeID string) ([]AIRecommendation, error) {
	q := `SELECT r.course_id, c.name, r.reason, r.competency_id, COALESCE(tc.name,''), r.expected_level, r.priority
		FROM training_ai_recommendations r
		JOIN training_courses c ON c.id=r.course_id
		LEFT JOIN training_competencies tc ON tc.id=r.competency_id
		WHERE r.employee_id=$1 ORDER BY r.priority DESC`
	rows, err := r.pool.Query(ctx, q, employeeID)
	if err != nil {
		return nil, repoErr("GetRecommendations", err)
	}
	defer rows.Close()
	var res []AIRecommendation
	for rows.Next() {
		var rec AIRecommendation
		if err := rows.Scan(&rec.CourseID, &rec.CourseName, &rec.Reason, &rec.CompetencyID, &rec.CompetencyName, &rec.ExpectedLevel, &rec.Priority); err != nil {
			return nil, repoErr("GetRecommendations", err)
		}
		res = append(res, rec)
	}
	return res, nil
}

// ---------------------------------------------------------------------------
// Notification events
// ---------------------------------------------------------------------------

func (r *Repository) CreateEvent(ctx context.Context, event *TrainingEvent) error {
	q := `INSERT INTO training_events (id,company_id,event_type,title,description,employee_id,enrollment_id,
		offering_id,severity,scheduled_for,metadata,created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := r.pool.Exec(ctx, q, event.ID, event.CompanyID, event.EventType, event.Title, event.Description,
		event.EmployeeID, event.EnrollmentID, event.OfferingID, event.Severity, event.ScheduledFor,
		event.Metadata, event.CreatedBy)
	return repoErr("CreateEvent", err)
}

func (r *Repository) ListPendingEvents(ctx context.Context, companyID string, limit int) ([]TrainingEvent, error) {
	q := `SELECT id,company_id,event_type,title,description,employee_id,enrollment_id,offering_id,
		severity,scheduled_for,processed_at,metadata,created_by,created_at
		FROM training_events WHERE company_id=$1 AND processed_at IS NULL AND
		(scheduled_for IS NULL OR scheduled_for<=NOW()) ORDER BY severity DESC, created_at ASC LIMIT $2`
	rows, err := r.pool.Query(ctx, q, companyID, limit)
	if err != nil {
		return nil, repoErr("ListPendingEvents", err)
	}
	defer rows.Close()
	var res []TrainingEvent
	for rows.Next() {
		var e TrainingEvent
		if err := rows.Scan(&e.ID, &e.CompanyID, &e.EventType, &e.Title, &e.Description, &e.EmployeeID,
			&e.EnrollmentID, &e.OfferingID, &e.Severity, &e.ScheduledFor, &e.ProcessedAt,
			&e.Metadata, &e.CreatedBy, &e.CreatedAt); err != nil {
			return nil, repoErr("ListPendingEvents", err)
		}
		res = append(res, e)
	}
	return res, nil
}

func (r *Repository) MarkEventProcessed(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `UPDATE training_events SET processed_at=NOW() WHERE id=$1`, id)
	return repoErr("MarkEventProcessed", err)
}

// ---------------------------------------------------------------------------
// CourseWithDetails helper
// ---------------------------------------------------------------------------

func (r *Repository) GetCourseWithDetails(ctx context.Context, companyID, courseID string) (*CourseWithDetails, error) {
	c, err := r.GetCourse(ctx, companyID, courseID)
	if err != nil {
		return nil, err
	}
	res := &CourseWithDetails{Course: *c}
	av, err := r.GetActiveVersion(ctx, courseID)
	if err != nil && err.Error() != "training_repo.GetActiveVersion: no rows in result set" {
		return nil, err
	}
	if av != nil {
		res.Versions = append(res.Versions, *av)
		contents, err := r.ListContents(ctx, av.ID, false)
		if err != nil {
			return nil, err
		}
		res.Contents = contents
	}
	offerings, _, err := r.ListOfferings(ctx, companyID, OfferingFilter{CourseID: &courseID})
	if err != nil {
		return nil, err
	}
	res.Offerings = offerings
	assessments, err := r.ListAssessments(ctx, companyID, courseID)
	if err != nil {
		return nil, err
	}
	res.Assessments = assessments
	competencies, err := r.ListCourseCompetencies(ctx, courseID)
	if err != nil {
		return nil, err
	}
	res.Competencies = competencies
	return res, nil
}
