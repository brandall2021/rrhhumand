package feed

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/models"
)

type FeedRepository struct {
	pool *pgxpool.Pool
}

func NewFeedRepository(pool *pgxpool.Pool) *FeedRepository {
	return &FeedRepository{pool: pool}
}

var mentionRegex = regexp.MustCompile(`@(\w+)`)

func (r *FeedRepository) CreatePost(ctx context.Context, post *models.Post) error {
	query := `
		INSERT INTO posts (id, company_id, author_id, content, visibility, pinned)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at`
	return r.pool.QueryRow(ctx, query,
		post.ID, post.CompanyID, post.AuthorID, post.Content,
		post.Visibility, post.Pinned,
	).Scan(&post.CreatedAt, &post.UpdatedAt)
}

func (r *FeedRepository) GetPostByID(ctx context.Context, id, companyID string) (*models.Post, error) {
	query := `
		SELECT
			p.id, p.company_id, p.author_id,
			e.first_name || ' ' || e.last_name,
			e.photo_url,
			p.content, p.visibility, p.pinned, p.created_at, p.updated_at
		FROM posts p
		JOIN employees e ON e.id = p.author_id
		WHERE p.id = $1 AND p.company_id = $2`

	post := &models.Post{}
	err := r.pool.QueryRow(ctx, query, id, companyID).Scan(
		&post.ID, &post.CompanyID, &post.AuthorID,
		&post.AuthorName, &post.AuthorPhoto,
		&post.Content, &post.Visibility, &post.Pinned,
		&post.CreatedAt, &post.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (r *FeedRepository) ListPosts(ctx context.Context, companyID string, offset, limit int) ([]models.Post, int64, error) {
	var total int64
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM posts WHERE company_id = $1`, companyID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT
			p.id, p.company_id, p.author_id,
			e.first_name || ' ' || e.last_name,
			e.photo_url,
			p.content, p.visibility, p.pinned, p.created_at, p.updated_at
		FROM posts p
		JOIN employees e ON e.id = p.author_id
		WHERE p.company_id = $1
		ORDER BY p.pinned DESC, p.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, companyID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(
			&p.ID, &p.CompanyID, &p.AuthorID,
			&p.AuthorName, &p.AuthorPhoto,
			&p.Content, &p.Visibility, &p.Pinned,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		posts = append(posts, p)
	}
	return posts, total, nil
}

func (r *FeedRepository) UpdatePost(ctx context.Context, post *models.Post) error {
	query := `
		UPDATE posts SET content=$1, visibility=$2, pinned=$3, updated_at=NOW()
		WHERE id=$4 AND company_id=$5`
	_, err := r.pool.Exec(ctx, query,
		post.Content, post.Visibility, post.Pinned,
		post.ID, post.CompanyID,
	)
	return err
}

func (r *FeedRepository) DeletePost(ctx context.Context, id, companyID string) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM posts WHERE id=$1 AND company_id=$2`, id, companyID,
	)
	return err
}

func (r *FeedRepository) AddMedia(ctx context.Context, media *models.PostMedia) error {
	query := `
		INSERT INTO post_media (id, post_id, media_type, url, filename, file_size)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)
		RETURNING id, created_at`
	return r.pool.QueryRow(ctx, query,
		media.PostID, media.MediaType, media.URL,
		media.Filename, media.FileSize,
	).Scan(&media.ID, &media.CreatedAt)
}

func (r *FeedRepository) GetMediaByPostID(ctx context.Context, postID string) ([]models.PostMedia, error) {
	query := `SELECT id, post_id, media_type, url, filename, file_size FROM post_media WHERE post_id=$1`
	rows, err := r.pool.Query(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var media []models.PostMedia
	for rows.Next() {
		var m models.PostMedia
		if err := rows.Scan(&m.ID, &m.PostID, &m.MediaType, &m.URL, &m.Filename, &m.FileSize); err != nil {
			return nil, err
		}
		media = append(media, m)
	}
	return media, nil
}

func (r *FeedRepository) AddComment(ctx context.Context, comment *models.Comment) error {
	query := `
		INSERT INTO comments (id, post_id, author_id, parent_id, content)
		VALUES (gen_random_uuid(), $1, $2, $3, $4)
		RETURNING id, created_at, updated_at`
	return r.pool.QueryRow(ctx, query,
		comment.PostID, comment.AuthorID, comment.ParentID, comment.Content,
	).Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)
}

func (r *FeedRepository) GetCommentsByPostID(ctx context.Context, postID string) ([]models.Comment, error) {
	query := `
		SELECT
			c.id, c.post_id, c.author_id,
			e.first_name || ' ' || e.last_name,
			e.photo_url,
			c.parent_id, c.content, c.created_at, c.updated_at
		FROM comments c
		JOIN employees e ON e.id = c.author_id
		WHERE c.post_id = $1
		ORDER BY c.created_at ASC`

	rows, err := r.pool.Query(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(
			&c.ID, &c.PostID, &c.AuthorID,
			&c.AuthorName, &c.AuthorPhoto,
			&c.ParentID, &c.Content, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (r *FeedRepository) GetCommentCount(ctx context.Context, postID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM comments WHERE post_id=$1`, postID,
	).Scan(&count)
	return count, err
}

func (r *FeedRepository) AddReaction(ctx context.Context, reaction *models.Reaction) error {
	query := `
		INSERT INTO reactions (id, post_id, employee_id, reaction_type)
		VALUES (gen_random_uuid(), $1, $2, $3)
		ON CONFLICT (post_id, employee_id, reaction_type) DO NOTHING
		RETURNING id, created_at`
	return r.pool.QueryRow(ctx, query,
		reaction.PostID, reaction.EmployeeID, reaction.ReactionType,
	).Scan(&reaction.ID, &reaction.CreatedAt)
}

func (r *FeedRepository) RemoveReaction(ctx context.Context, postID, employeeID, reactionType string) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM reactions WHERE post_id=$1 AND employee_id=$2 AND reaction_type=$3`,
		postID, employeeID, reactionType,
	)
	return err
}

func (r *FeedRepository) GetReactionsByPostID(ctx context.Context, postID string) ([]models.Reaction, error) {
	query := `SELECT id, post_id, employee_id, reaction_type, created_at FROM reactions WHERE post_id=$1`
	rows, err := r.pool.Query(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reactions []models.Reaction
	for rows.Next() {
		var rx models.Reaction
		if err := rows.Scan(&rx.ID, &rx.PostID, &rx.EmployeeID, &rx.ReactionType, &rx.CreatedAt); err != nil {
			return nil, err
		}
		reactions = append(reactions, rx)
	}
	return reactions, nil
}

func (r *FeedRepository) GetReactionCounts(ctx context.Context, postID string) (map[string]int, error) {
	query := `
		SELECT reaction_type, COUNT(*)
		FROM reactions WHERE post_id=$1
		GROUP BY reaction_type`

	rows, err := r.pool.Query(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var rType string
		var count int
		if err := rows.Scan(&rType, &count); err != nil {
			return nil, err
		}
		counts[rType] = count
	}
	return counts, nil
}

func (r *FeedRepository) GetUserReaction(ctx context.Context, postID, employeeID string) (*string, error) {
	var rType string
	err := r.pool.QueryRow(ctx,
		`SELECT reaction_type FROM reactions WHERE post_id=$1 AND employee_id=$2`,
		postID, employeeID,
	).Scan(&rType)
	if err != nil {
		return nil, err
	}
	return &rType, nil
}

func (r *FeedRepository) AddMentions(ctx context.Context, postID, commentID *string, mentionedIDs []string) error {
	for _, empID := range mentionedIDs {
		_, err := r.pool.Exec(ctx,
			`INSERT INTO mentions (id, post_id, comment_id, mentioned_employee_id) VALUES (gen_random_uuid(), $1, $2, $3)`,
			postID, commentID, empID,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func ExtractMentionIDs(content string) []string {
	matches := mentionRegex.FindAllStringSubmatch(content, -1)
	var ids []string
	seen := make(map[string]bool)
	for _, m := range matches {
		if len(m) > 1 && !seen[m[1]] {
			ids = append(ids, m[1])
			seen[m[1]] = true
		}
	}
	return ids
}

func (r *FeedRepository) FindEmployeeByNumber(ctx context.Context, companyID, employeeNumber string) (string, error) {
	var id string
	err := r.pool.QueryRow(ctx,
		`SELECT id FROM employees WHERE company_id=$1 AND employee_number=$2`,
		companyID, employeeNumber,
	).Scan(&id)
	return id, err
}

func BuildPostQuery(baseQuery string, args []interface{}, argIdx int, search string) (string, []interface{}, int) {
	if search != "" {
		baseQuery += fmt.Sprintf(` AND p.content ILIKE $%d`, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	return baseQuery, args, argIdx
}

func (r *FeedRepository) ListPostsByEmployee(ctx context.Context, companyID, employeeID string, offset, limit int) ([]models.Post, int64, error) {
	var total int64
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM posts WHERE company_id=$1 AND author_id=$2`, companyID, employeeID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT
			p.id, p.company_id, p.author_id,
			e.first_name || ' ' || e.last_name,
			e.photo_url,
			p.content, p.visibility, p.pinned, p.created_at, p.updated_at
		FROM posts p
		JOIN employees e ON e.id = p.author_id
		WHERE p.company_id=$1 AND p.author_id=$2
		ORDER BY p.created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.pool.Query(ctx, query, companyID, employeeID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(
			&p.ID, &p.CompanyID, &p.AuthorID,
			&p.AuthorName, &p.AuthorPhoto,
			&p.Content, &p.Visibility, &p.Pinned,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		posts = append(posts, p)
	}
	return posts, total, nil
}

func BuildSearchClause(baseQuery string, args []interface{}, argIdx int, search string) (string, []interface{}, int) {
	if search != "" {
		baseQuery += fmt.Sprintf(` AND p.content ILIKE $%d`, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	return baseQuery, args, argIdx
}

var _ = strings.TrimSpace
