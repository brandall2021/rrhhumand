package feed

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rrhhumand/api/internal/models"
)

type FeedService struct {
	repo *FeedRepository
}

func NewFeedService(repo *FeedRepository) *FeedService {
	return &FeedService{repo: repo}
}

type CreatePostRequest struct {
	Content    string          `json:"content" validate:"required"`
	Visibility *string         `json:"visibility,omitempty"`
	Pinned     *bool           `json:"pinned,omitempty"`
	Media      []MediaRequest  `json:"media,omitempty"`
}

type MediaRequest struct {
	MediaType string  `json:"media_type" validate:"required"`
	URL       string  `json:"url" validate:"required"`
	Filename  *string `json:"filename,omitempty"`
	FileSize  *int    `json:"file_size,omitempty"`
}

type UpdatePostRequest struct {
	Content    *string `json:"content,omitempty"`
	Visibility *string `json:"visibility,omitempty"`
	Pinned     *bool   `json:"pinned,omitempty"`
}

type CreateCommentRequest struct {
	Content  string  `json:"content" validate:"required"`
	ParentID *string `json:"parent_id,omitempty"`
}

type ReactionRequest struct {
	ReactionType string `json:"reaction_type" validate:"required"`
}

func (s *FeedService) CreatePost(ctx context.Context, companyID, authorID string, req *CreatePostRequest) (*models.Post, error) {
	visibility := "company"
	if req.Visibility != nil {
		visibility = *req.Visibility
	}
	pinned := false
	if req.Pinned != nil {
		pinned = *req.Pinned
	}

	post := &models.Post{
		ID:         uuid.New().String(),
		CompanyID:  companyID,
		AuthorID:   authorID,
		Content:    req.Content,
		Visibility: visibility,
		Pinned:     pinned,
	}

	if err := s.repo.CreatePost(ctx, post); err != nil {
		return nil, err
	}

	for _, m := range req.Media {
		media := &models.PostMedia{
			PostID:    post.ID,
			MediaType: m.MediaType,
			URL:       m.URL,
			Filename:  m.Filename,
			FileSize:  m.FileSize,
		}
		_ = s.repo.AddMedia(ctx, media)
	}

	mentionIDs := ExtractMentionIDs(req.Content)
	if len(mentionIDs) > 0 {
		_ = s.repo.AddMentions(ctx, &post.ID, nil, mentionIDs)
	}

	return s.repo.GetPostByID(ctx, post.ID, companyID)
}

func (s *FeedService) GetPostByID(ctx context.Context, id, companyID, currentUserID string) (*models.Post, error) {
	post, err := s.repo.GetPostByID(ctx, id, companyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	media, _ := s.repo.GetMediaByPostID(ctx, post.ID)
	post.Media = media

	comments, _ := s.repo.GetCommentsByPostID(ctx, post.ID)
	post.Comments = comments

	count, _ := s.repo.GetCommentCount(ctx, post.ID)
	post.CommentCount = count

	counts, _ := s.repo.GetReactionCounts(ctx, post.ID)
	post.ReactionCounts = counts

	empID, _ := s.repo.FindEmployeeByNumber(ctx, companyID, currentUserID)
	if empID != "" {
		userReaction, _ := s.repo.GetUserReaction(ctx, post.ID, empID)
		post.UserReaction = userReaction
	}

	return post, nil
}

func (s *FeedService) ListPosts(ctx context.Context, companyID string, params *models.PaginationParams, search string) ([]models.Post, int64, error) {
	posts, total, err := s.repo.ListPosts(ctx, companyID, params.Offset, params.PerPage)
	if err != nil {
		return nil, 0, err
	}

	for i := range posts {
		count, _ := s.repo.GetCommentCount(ctx, posts[i].ID)
		posts[i].CommentCount = count

		counts, _ := s.repo.GetReactionCounts(ctx, posts[i].ID)
		posts[i].ReactionCounts = counts
	}

	return posts, total, nil
}

func (s *FeedService) UpdatePost(ctx context.Context, id, companyID, authorID string, req *UpdatePostRequest) (*models.Post, error) {
	post, err := s.repo.GetPostByID(ctx, id, companyID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	if post.AuthorID != authorID {
		return nil, errors.New("you can only edit your own posts")
	}

	if req.Content != nil {
		post.Content = *req.Content
	}
	if req.Visibility != nil {
		post.Visibility = *req.Visibility
	}
	if req.Pinned != nil {
		post.Pinned = *req.Pinned
	}

	if err := s.repo.UpdatePost(ctx, post); err != nil {
		return nil, err
	}

	return s.repo.GetPostByID(ctx, post.ID, companyID)
}

func (s *FeedService) DeletePost(ctx context.Context, id, companyID, authorID string) error {
	post, err := s.repo.GetPostByID(ctx, id, companyID)
	if err != nil {
		return errors.New("post not found")
	}

	if post.AuthorID != authorID {
		return fmt.Errorf("you can only delete your own posts")
	}

	return s.repo.DeletePost(ctx, id, companyID)
}

func (s *FeedService) AddComment(ctx context.Context, postID, companyID, authorID string, req *CreateCommentRequest) (*models.Comment, error) {
	_, err := s.repo.GetPostByID(ctx, postID, companyID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	comment := &models.Comment{
		PostID:   postID,
		AuthorID: authorID,
		ParentID: req.ParentID,
		Content:  req.Content,
	}

	if err := s.repo.AddComment(ctx, comment); err != nil {
		return nil, err
	}

	mentionIDs := ExtractMentionIDs(req.Content)
	if len(mentionIDs) > 0 {
		_ = s.repo.AddMentions(ctx, nil, &comment.ID, mentionIDs)
	}

	return comment, nil
}

func (s *FeedService) AddReaction(ctx context.Context, postID, companyID, employeeID string, req *ReactionRequest) error {
	_, err := s.repo.GetPostByID(ctx, postID, companyID)
	if err != nil {
		return errors.New("post not found")
	}

	reaction := &models.Reaction{
		PostID:       postID,
		EmployeeID:   employeeID,
		ReactionType: req.ReactionType,
	}

	return s.repo.AddReaction(ctx, reaction)
}

func (s *FeedService) RemoveReaction(ctx context.Context, postID, companyID, employeeID, reactionType string) error {
	_, err := s.repo.GetPostByID(ctx, postID, companyID)
	if err != nil {
		return errors.New("post not found")
	}

	return s.repo.RemoveReaction(ctx, postID, employeeID, reactionType)
}
