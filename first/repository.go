package first

import (
	"context"

	"github.com/mabdh/godbtx/store"
	"github.com/mabdh/godbtx/user"
)

type UserRepository struct {
	dbc *store.Client
}

func NewUserRepository(dbc *store.Client) *UserRepository {
	return &UserRepository{
		dbc: dbc,
	}
}

func (r *UserRepository) Create(ctx context.Context, u user.User) (string, error) {
	var id string
	err := r.dbc.DB.QueryRowContext(ctx, "INSERT INTO users (name,email) VALUES (?,?)", u.Name, u.Email).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}
