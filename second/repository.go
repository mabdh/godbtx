package second

import (
	"context"
	"database/sql"

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

func (r *UserRepository) CreateWithFn(ctx context.Context, u user.User, postProcessFn func(id string) error) (string, error) {
	var id string

	if err := r.dbc.WithTransaction(ctx, func(tx *sql.Tx) error {
		if err := tx.QueryRowContext(ctx, "INSERT INTO users (name,email) VALUES (?,?)", u.Name, u.Email).Scan(&id); err != nil {
			return err
		}

		if err := postProcessFn(id); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return "", err
	}

	return id, nil
}
