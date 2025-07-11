package users

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/eser/aya.is-services/pkg/ajan/logfx"
	"github.com/eser/aya.is-services/pkg/lib/cursors"
)

var (
	ErrFailedToGetRecord    = errors.New("failed to get record")
	ErrFailedToListRecords  = errors.New("failed to list records")
	ErrFailedToCreateRecord = errors.New("failed to create record")
	ErrFailedToUpdateRecord = errors.New("failed to update record")
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

type Service struct {
	logger      *logfx.Logger
	repo        Repository
	idGenerator RecordIDGenerator

	authProviders map[string]AuthProvider
}

func NewService(
	logger *logfx.Logger,
	repo Repository,
	authProviders map[string]AuthProvider,
) *Service {
	return &Service{
		logger:        logger,
		repo:          repo,
		idGenerator:   DefaultIDGenerator,
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
