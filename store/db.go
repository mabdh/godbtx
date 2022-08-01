package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type transactionContextKey struct{}

type Config struct {
	URL string
}

type Client struct {
	DB *sqlx.DB
}

func New(cfg Config) (*Client, error) {
	d, err := sqlx.Open("postgres", cfg.URL)
	if err != nil {
		return nil, err
	}

	if err = d.Ping(); err != nil {
		return nil, err
	}

	return &Client{
		DB: d,
	}, nil
}

func (c *Client) Close() error {
	return c.DB.Close()
}

func (c *Client) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	if err := fn(tx); err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			return fmt.Errorf("rollback transaction error: %v (original error: %w)", txErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

func (c *Client) SetContextWithTransaction(ctx context.Context) context.Context {
	tx, _ := c.DB.BeginTxx(ctx, nil)
	return context.WithValue(ctx, transactionContextKey{}, tx)
}

func getContextFromTransaction(ctx context.Context) *sqlx.Tx {
	if tx, ok := ctx.Value(transactionContextKey{}).(*sqlx.Tx); !ok {
		return nil
	} else {
		return tx
	}
}

func (c *Client) Rollback(ctx context.Context) error {
	if tx := getContextFromTransaction(ctx); tx != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return nil
	}
	return errors.New("no transaction")
}

func (c *Client) Commit(ctx context.Context) error {
	if tx := getContextFromTransaction(ctx); tx != nil {
		if err := tx.Commit(); err != nil {
			return err
		}
		return nil
	}
	return errors.New("no transaction")
}

func (c *Client) GetDB(ctx context.Context) (db sqlx.QueryerContext) {
	db = c.DB
	if tx := getContextFromTransaction(ctx); tx != nil {
		db = tx
	}
	return db
}
