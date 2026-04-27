package connection

import (
	"errors"
	"fmt"
	"log/slog"
	"rdpms25-go-rpc-service/pkg/config"
	"time"

	"github.com/nats-io/nats.go"
)

func NewNatsConnection(clientId string, config *config.NatsConf) (*nats.Conn, error) {
	if config == nil {
		return nil, errors.New("nats config is nil")
	}
	natsURL := fmt.Sprintf("://%s:%d", config.Host, config.Port)
	if config.Tls {
		natsURL = "tls" + natsURL
	} else {
		natsURL = "tcp" + natsURL
	}
	return connectNats(clientId, natsURL, config)
}

func connectNats(clientId, natsURL string, config *config.NatsConf) (*nats.Conn, error) {
	opts := []nats.Option{
		nats.Name(clientId),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(5 * time.Second),
		nats.RetryOnFailedConnect(true),
		nats.PingInterval(2 * time.Second),
		nats.ConnectHandler(func(nc *nats.Conn) {
			slog.Info("connected to nats", "url", nc.ConnectedUrl())
		}),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			slog.Error("disconnected from nats, will reconnect", "err", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			slog.Info("reconnected to nats", "url", nc.ConnectedUrl())
		}),
		nats.ReconnectErrHandler(func(nc *nats.Conn, err error) {
			slog.Error("reconnection to nats failed", "err", err)
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			slog.Warn("nats connection closed", "err", nc.LastError())
		}),
	}
	if config.TLSConf != nil {
		opts = append(opts, nats.Secure(config.TLSConf))
	}
	if config.User != "" && config.Password != "" {
		opts = append(opts, nats.UserInfo(config.User, config.Password))
	}
	return nats.Connect(natsURL, opts...)
}
