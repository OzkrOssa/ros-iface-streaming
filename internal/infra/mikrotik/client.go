package mikrotik

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/OzkrOssa/ros-iface-streamer/internal/domain"
	"github.com/OzkrOssa/ros-iface-streamer/pkg/config"
	"github.com/go-routeros/routeros/v3"
)

type Mikrotik struct {
	rosClient *routeros.Client
	config    *config.RouterOS
}

var _ domain.Mikrotik = (*Mikrotik)(nil)

func New(config *config.RouterOS) *Mikrotik {
	return &Mikrotik{
		config: config,
	}
}

func (m *Mikrotik) Connect(host string) error {
	if m.rosClient != nil {
		_ = m.rosClient.Close()
	}

	url := host + ":" + m.config.Port
	rosClient, err := routeros.Dial(url, m.config.User, m.config.Pass)
	if err != nil {
		return fmt.Errorf("failed to connect to mikrotik: %w", err)
	}

	m.rosClient = rosClient
	return nil
}

// TODO: add cache for ifaces to avoid multiple requests to mikrotik
func (m *Mikrotik) GetStreamingTraffic(ctx context.Context, iface string) (chan *domain.Traffic, error) {
	ifaceName, err := m.VerifyIface(ctx, iface)
	if err != nil {
		return nil, errors.New("failed to verify iface")
	}

	errChan := m.rosClient.Async()
	m.rosClient.Queue = 100

	cmd := []string{
		"/interface/monitor-traffic",
		fmt.Sprintf("=interface=%s", ifaceName),
	}

	listener, err := m.rosClient.ListenArgsContext(ctx, cmd)
	if err != nil {
		return nil, errors.New("failed to listen to mikrotik")
	}

	traffic := make(chan *domain.Traffic, 10)

	go func() {
		defer listener.Cancel()
		defer close(traffic)

		for {
			select {
			case <-ctx.Done():
				m.rosClient.Close()
				return
			case sen, ok := <-listener.Chan():
				if !ok {
					return
				}

				data := sen.Map
				rx, _ := strconv.ParseUint(data["rx-bits-per-second"], 10, 64)
				tx, _ := strconv.ParseUint(data["tx-bits-per-second"], 10, 64)

				traffic <- &domain.Traffic{
					Iface: ifaceName,
					Tx:    tx,
					Rx:    rx,
				}
			case err := <-errChan:
				slog.Error("Error in connection", "error", err)
				return
			}
		}
	}()

	return traffic, nil
}

func (m *Mikrotik) VerifyIface(ctx context.Context, iface string) (exactIface string, err error) {
	cmd := []string{
		"/interface/print",
	}

	reply, err := m.rosClient.RunArgsContext(ctx, cmd)
	if err != nil {
		return "", errors.New("failed to get interfaces")
	}

	for _, re := range reply.Re {
		if strings.Contains(re.Map["name"], iface) {
			return re.Map["name"], nil
		}
	}
	return "", fmt.Errorf("iface %s not found", iface)
}

func (m *Mikrotik) Close() error {
	if m.rosClient != nil {
		return m.rosClient.Close()
	}
	return nil
}
