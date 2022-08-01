package second

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
	return s.r.CreateWithFn(ctx, u, func(id string) error {
		if _, err := s.auditService.Create(ctx, audit.Audit{
			Action:   "create",
			Domain:   "user",
			DomainID: id,
		}); err != nil {
			return err
		}
		return nil
	})
}
