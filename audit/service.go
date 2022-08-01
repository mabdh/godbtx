package audit

import "context"

type Service struct {
	client Client
}

func NewService(client Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) Create(ctx context.Context, a Audit) (string, error) {
	return s.client.Create(ctx, a)
}
