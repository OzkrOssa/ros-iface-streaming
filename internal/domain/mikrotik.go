package domain

import (
	"context"
)

type Mikrotik interface {
	Connect(host string) error
	GetStreamingTraffic(ctx context.Context, iface string) (chan *Traffic, error)
	VerifyIface(ctx context.Context, iface string) (exactIface string, err error)
	Close() error
}
