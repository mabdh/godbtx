package first

import (
	"context"

	"github.com/mabdh/godbtx/audit"
	"github.com/mabdh/godbtx/user"
)

type UserService struct {
	r            UserRepository
	auditService audit.Service
}

func NewUserService(r UserRepository, as audit.Service) *UserService {
	return &UserService{
		r:            r,
		auditService: as,
	}
}

func (s *UserService) Create(ctx context.Context, u user.User) (string, error) {
	id, err := s.r.Create(ctx, u)
	if err != nil {
		return "", err
	}

	if _, err := s.auditService.Create(ctx, audit.Audit{
		Action:   "create",
		Domain:   "user",
		DomainID: id,
	}); err != nil {
		return "", err
	}

	return id, nil
}
