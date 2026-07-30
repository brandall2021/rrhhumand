package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/models"
)

type RefreshTokenRepository struct {
	pool *pgxpool.Pool
}

func NewRefreshTokenRepository(pool *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{pool: pool}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, token *models.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, company_id, token_hash, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at`
	return r.pool.QueryRow(ctx, query,
		token.ID, token.UserID, token.CompanyID, token.TokenHash, token.ExpiresAt,
	).Scan(&token.CreatedAt)
}

func (r *RefreshTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	query := `
		SELECT id, user_id, company_id, token_hash, expires_at, revoked, created_at
		FROM refresh_tokens
		WHERE token_hash = $1`

	token := &models.RefreshToken{}
	err := r.pool.QueryRow(ctx, query, tokenHash).Scan(
		&token.ID, &token.UserID, &token.CompanyID, &token.TokenHash,
		&token.ExpiresAt, &token.Revoked, &token.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, tokenHash string) error {
	query := `UPDATE refresh_tokens SET revoked = true WHERE token_hash = $1`
	_, err := r.pool.Exec(ctx, query, tokenHash)
	return err
}

func (r *RefreshTokenRepository) RevokeAllByUser(ctx context.Context, userID string) error {
	query := `UPDATE refresh_tokens SET revoked = true WHERE user_id = $1`
	_, err := r.pool.Exec(ctx, query, userID)
	return err
}

func (r *RefreshTokenRepository) Cleanup(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW() OR revoked = true`
	_, err := r.pool.Exec(ctx, query)
	return err
}

type UserRepositoryInterface interface {
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id string) (*models.User, error)
	UpdateLastLogin(ctx context.Context, userID string) error
	GetRolesByCompany(ctx context.Context, userID, companyID string) ([]string, error)
	GetCompaniesByUser(ctx context.Context, userID string) ([]string, error)
}

type RoleRepositoryInterface interface {
	FindByName(ctx context.Context, name string) (*models.Role, error)
}

type CompanyResolver interface {
	FindBySlug(ctx context.Context, slug string) (*models.Company, error)
}

type AuthService struct {
	userRepo        UserRepositoryInterface
	tokenRepo       *RefreshTokenRepository
	jwtService      *JWTService
	roleRepo        RoleRepositoryInterface
	companyResolver CompanyResolver
}

func NewAuthService(
	userRepo UserRepositoryInterface,
	tokenRepo *RefreshTokenRepository,
	jwtService *JWTService,
	roleRepo RoleRepositoryInterface,
	companyResolver CompanyResolver,
) *AuthService {
	return &AuthService{
		userRepo:        userRepo,
		tokenRepo:       tokenRepo,
		jwtService:      jwtService,
		roleRepo:        roleRepo,
		companyResolver: companyResolver,
	}
}

type LoginRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required"`
	CompanySlug string `json:"company_slug"`
}

type LoginResponse struct {
	TokenPair *TokenPair      `json:"token_pair"`
	User      *UserResponse   `json:"user"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (s *AuthService) Login(ctx context.Context, req *LoginRequest, companyID string) (*LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !user.Active {
		return nil, ErrInvalidCredentials
	}

	if !CheckPassword(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	_ = s.userRepo.UpdateLastLogin(ctx, user.ID)

	roles, err := s.userRepo.GetRolesByCompany(ctx, user.ID, companyID)
	if err != nil {
		roles = []string{}
	}

	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Email, companyID, roles)
	if err != nil {
		return nil, err
	}

	tokenHash := HashToken(tokenPair.RefreshToken)
	expiresAt := time.Now().Add(s.jwtService.RefreshExpiration())

	refreshToken := &models.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		CompanyID: &companyID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}
	_ = s.tokenRepo.Create(ctx, refreshToken)

	return &LoginResponse{
		TokenPair: tokenPair,
		User: &UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenString, companyID string) (*TokenPair, error) {
	tokenHash := HashToken(refreshTokenString)

	refreshToken, err := s.tokenRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	if refreshToken.Revoked {
		return nil, ErrTokenRevoked
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		return nil, ErrTokenExpired
	}

	user, err := s.userRepo.FindByID(ctx, refreshToken.UserID)
	if err != nil {
		return nil, err
	}

	roles, err := s.userRepo.GetRolesByCompany(ctx, user.ID, companyID)
	if err != nil {
		roles = []string{}
	}

	_ = s.tokenRepo.Revoke(ctx, tokenHash)

	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Email, companyID, roles)
	if err != nil {
		return nil, err
	}

	newHash := HashToken(tokenPair.RefreshToken)
	newRefresh := &models.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		CompanyID: &companyID,
		TokenHash: newHash,
		ExpiresAt: time.Now().Add(s.jwtService.RefreshExpiration()),
	}
	_ = s.tokenRepo.Create(ctx, newRefresh)

	return tokenPair, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshTokenString string) error {
	tokenHash := HashToken(refreshTokenString)
	return s.tokenRepo.Revoke(ctx, tokenHash)
}
