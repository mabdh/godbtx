package audit

import (
	"context"
	"time"
)

type Client interface {
	Create(context.Context, Audit) (string, error)
	Get(context.Context, string) (Audit, error)
}

type Audit struct {
	ID        string
	Action    string
	Domain    string
	DomainID  string
	CreatedAt time.Time
}
