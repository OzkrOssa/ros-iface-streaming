package domain

import "context"

type Manager interface {
	HandleConnection(ctx context.Context) error
}
