package users

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

var (
	ErrFailedToGetRecord     = errors.New("failed to get record")
	ErrFailedToListRecords   = errors.New("failed to list records")
	ErrFailedToCreateRecord  = errors.New("failed to create record")
	ErrFailedToUpdateRecord  = errors.New("failed to update record")
	ErrInvalidToken          = errors.New("invalid token")
	ErrSessionExpired        = errors.New("session expired")
	ErrJWTNotConfigured      = errors.New("JWT not configured")
	ErrInvalidSigningMethod  = errors.New("invalid JWT signing method")
	ErrFailedToGenerateToken = errors.New("failed to generate token")
	ErrFailedToGetUser       = errors.New("failed to get user")
	ErrFailedToUpdateSession = errors.New("failed to update session")
)

type Repository interface {
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	ListUsers(
		ctx context.Context,
		cursor *cursors.Cursor,
	) (cursors.Cursored[[]*User], error)
	CreateUser(ctx context.Context, user *User) error

	CreateSession(ctx context.Context, session *Session) error
	GetSessionByID(ctx context.Context, id string) (*Session, error)
	UpdateSessionLoggedInAt(ctx context.Context, id string, loggedInAt time.Time) error
}

type AuthProvider interface {
	// InitiateOAuth returns the URL to redirect the user to GitHub for login, and the state to track the request.
	InitiateOAuth(
		ctx context.Context,
		redirectURI string,
	) (authURL string, state OAuthState, err error)

	// HandleOAuthCallback exchanges the code for a token, fetches user info, upserts user,
	// creates session, and returns JWT.
	HandleOAuthCallback(ctx context.Context, code string, state string) (AuthResult, error)
}

type TokenService interface {
	// ParseToken validates a JWT token and returns the claims
	ParseToken(tokenStr string) (*JWTClaims, error)

	// GenerateToken creates a new JWT token with the given claims
	GenerateToken(claims *JWTClaims) (string, error)
}

type Service struct {
	logger        *logfx.Logger
	repo          Repository
	idGenerator   RecordIDGenerator
	tokenService  TokenService
	AuthConfig    *AuthConfig
	authProviders map[string]AuthProvider
}

func NewService(
	logger *logfx.Logger,
	repo Repository,
	tokenService TokenService,
	authConfig *AuthConfig,
	authProviders map[string]AuthProvider,
) *Service {
	return &Service{
		logger:        logger,
		repo:          repo,
		idGenerator:   DefaultIDGenerator,
		tokenService:  tokenService,
		AuthConfig:    authConfig,
		authProviders: authProviders,
	}
}

func (s *Service) GetByID(ctx context.Context, id string) (*User, error) {
	record, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w(id: %s): %w", ErrFailedToGetRecord, id, err)
	}

	return record, nil
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*User, error) {
	record, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("%w(email: %s): %w", ErrFailedToGetRecord, email, err)
	}

	return record, nil
}

func (s *Service) List(
	ctx context.Context,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*User], error) {
	records, err := s.repo.ListUsers(ctx, cursor)
	if err != nil {
		return cursors.Cursored[[]*User]{}, fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	}

	return records, nil
}

func (s *Service) Create(ctx context.Context, user *User) error {
	err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToCreateRecord, err)
	}

	return nil
}

func (s *Service) GetSessionByID(ctx context.Context, id string) (*Session, error) {
	session, err := s.repo.GetSessionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w(id: %s): %w", ErrFailedToGetRecord, id, err)
	}

	return session, nil
}

func (s *Service) UpdateSessionLoggedInAt(
	ctx context.Context,
	id string,
	loggedInAt time.Time,
) error {
	err := s.repo.UpdateSessionLoggedInAt(ctx, id, loggedInAt)
	if err != nil {
		return fmt.Errorf("%w(id: %s): %w", ErrFailedToUpdateRecord, id, err)
	}

	return nil
}

func (s *Service) GetAuthProvider(provider string) AuthProvider {
	service, serviceOk := s.authProviders[provider]
	if !serviceOk {
		return nil
	}

	return service
}

// RefreshToken validates the current JWT token and issues a new one with extended expiration.
func (s *Service) RefreshToken( //nolint:funlen
	ctx context.Context,
	tokenStr string,
) (*AuthResult, error) {
	s.logger.DebugContext(ctx, "Attempting to refresh JWT token")

	// Parse and validate current token using the token service
	claims, err := s.tokenService.ParseToken(tokenStr)
	if err != nil {
		s.logger.WarnContext(ctx, "Failed to parse JWT token", slog.String("error", err.Error()))

		return nil, fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}

	s.logger.DebugContext(ctx, "JWT token parsed successfully",
		slog.String("user_id", claims.UserID),
		slog.String("session_id", claims.SessionID))

	// Verify session is still active
	session, err := s.GetSessionByID(ctx, claims.SessionID)
	if err != nil || session.Status != "active" {
		s.logger.WarnContext(ctx, "Session is not active",
			slog.String("session_id", claims.SessionID),
			slog.String("status", session.Status))

		return nil, ErrSessionExpired
	}

	// Get user details
	user, err := s.GetByID(ctx, claims.UserID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get user for token refresh",
			slog.String("user_id", claims.UserID),
			slog.String("error", err.Error()))

		return nil, fmt.Errorf("%w: %w", ErrFailedToGetUser, err)
	}

	// Generate new JWT with extended expiration
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour) //nolint:mnd

	newClaims := &JWTClaims{
		UserID:    claims.UserID,
		SessionID: claims.SessionID,
		ExpiresAt: expiresAt.Unix(),
	}

	tokenString, err := s.tokenService.GenerateToken(newClaims)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to generate new JWT token",
			slog.String("user_id", claims.UserID),
			slog.String("error", err.Error()))

		return nil, fmt.Errorf("%w: %w", ErrFailedToGenerateToken, err)
	}

	// Update session's last activity
	updateErr := s.UpdateSessionLoggedInAt(ctx, claims.SessionID, now)
	if updateErr != nil {
		s.logger.WarnContext(ctx, "Failed to update session logged in time",
			slog.String("session_id", claims.SessionID),
			slog.String("error", updateErr.Error()))
		// Don't fail the whole operation for this
	}

	s.logger.DebugContext(ctx, "JWT token refreshed successfully",
		slog.String("user_id", claims.UserID),
		slog.String("session_id", claims.SessionID))

	return &AuthResult{
		User:      user,
		SessionID: claims.SessionID,
		JWT:       tokenString,
		ExpiresAt: expiresAt,
	}, nil
}
