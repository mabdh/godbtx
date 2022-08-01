package third

import (
	"context"

	"github.com/mabdh/godbtx/audit"
	"github.com/mabdh/godbtx/user"
)

type UserService struct {
	repo         UserRepository
	auditService audit.Service
}

func NewUserService(repo UserRepository, as audit.Service) *UserService {
	return &UserService{
		repo:         repo,
		auditService: as,
	}
}

func (s *UserService) Create(ctx context.Context, u user.User) (string, error) {
	ctxTx := s.repo.WithTransaction(ctx)

	id, err := s.Create(ctxTx, u)
	if err != nil {
		if err := s.repo.Rollback(ctxTx, err); err != nil {
			return "", err
		}
		return "", err
	}

	if _, err := s.auditService.Create(ctxTx, audit.Audit{
		Action:   "create",
		Domain:   "user",
		DomainID: id,
	}); err != nil {
		if err := s.repo.Rollback(ctxTx, err); err != nil {
			return "", err
		}
		return "", err
	}

	if err := s.repo.Commit(ctxTx); err != nil {
		return "", err
	}
	return id, nil
}
