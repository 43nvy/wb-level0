package service

import "github.com/nats-io/stan.go"

type Config struct {
	URL     string
	Cluster string
	Client  string
}

func NewNATSConn(cfg Config) (stan.Conn, error) {
	sc, err := stan.Connect(cfg.Cluster, cfg.Client, stan.NatsURL(cfg.URL))
	if err != nil {
		return nil, err
	}

	return sc, nil
}
