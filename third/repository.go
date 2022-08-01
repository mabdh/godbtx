package third

import (
	"context"

	"github.com/mabdh/godbtx/store"
	"github.com/mabdh/godbtx/user"
)

type Transactor interface {
	WithTransaction(ctx context.Context) context.Context
	Rollback(ctx context.Context, err error) error
	Commit(ctx context.Context) error
}

type UserRepository struct {
	Transactor
	dbc *store.Client
}

func NewUserRepository(dbc *store.Client) *UserRepository {
	return &UserRepository{
		dbc: dbc,
	}
}

func (r *UserRepository) Create(ctx context.Context, u user.User) (string, error) {

	db := r.dbc.GetDBFromCtx(ctx)
	var id string
	err := db.QueryRowContext(ctx, "INSERT INTO users (name,email) VALUES (?,?)", u.Name, u.Email).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}
