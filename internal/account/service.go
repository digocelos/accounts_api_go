package account

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  time.Now().UTC,
	}
}

func (s *Service) Create(ctx context.Context, in CreateInput) (Account, bool /*created*/, error) {
	in.Name = strings.TrimSpace(in.Name)
	in.Document = strings.TrimSpace(in.Document)
	if in.Email != nil {
		e := strings.TrimSpace(*in.Email)
		in.Email = &e
	}

	if in.Name == "" {
		return Account{}, false, fmt.Errorf("%w: name is required", ErrValidation)
	}

	if in.Document == "" {
		return Account{}, false, fmt.Errorf("%w: document is required", ErrValidation)
	}

	if len(in.Document) > 32 {
		return Account{}, false, fmt.Errorf("%w: document too long", ErrValidation)
	}

	now := s.now()
	acc := Account{
		ID:        uuid.NewString(),
		Name:      in.Name,
		Document:  in.Document,
		Email:     in.Email,
		Version:   1,
		CreatedAt: now,
		UpdatedAt: now,
	}

	inserted, err := s.repo.InsertIdempotentByDocument(ctx, acc)

	if err != nil {
		return Account{}, false, err
	}

	if inserted {
		return acc, true, nil
	}

	// Already exists => return existing idempotency
	existing, err := s.repo.GetByDocument(ctx, in.Document)
	if err != nil {
		return Account{}, false, err
	}

	return existing, false, nil
}

func (s *Service) Get(ctx context.Context, id string) (Account, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Account{}, fmt.Errorf("%w: id is required", ErrValidation)
	}
	return s.repo.GetById(ctx, id)
}

func (s *Service) Update(ctx context.Context, id string, in UpdateInput) (Account, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Account{}, fmt.Errorf("%w: id is required", ErrValidation)
	}

	if in.ExpectedVersion <= 0 {
		return Account{}, fmt.Errorf("%w: expected_version must be > 0", ErrValidation)
	}

	if in.Name != nil {
		n := strings.TrimSpace(*in.Name)
		if n == "" {
			return Account{}, fmt.Errorf("%w: name cannot be empty", ErrValidation)
		}
	}

	if in.Email != nil {
		e := strings.TrimSpace(*in.Email)
		if e == "" {
			return Account{}, fmt.Errorf("%w: email cannot be empty", ErrValidation)
		}
	}

	return s.repo.UpdateWithOptimisticLock(ctx, id, in.ExpectedVersion, in.Name, in.Email)
}
