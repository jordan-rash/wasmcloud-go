package client

import (
	"time"

	"github.com/nats-io/nats.go"
)

/// Collect results until timeout has elapsed
func (c Client) CollectTimeout(subject string, data *[]byte, timeoutOverride *time.Duration) ([]*nats.Msg, error) {
	timeout := c.timeout
	if timeoutOverride != nil {
		timeout = *timeoutOverride
	}

	sub := c.nc.NewRespInbox()
	var ret []*nats.Msg
	s, err := c.nc.Subscribe(sub, func(msg *nats.Msg) {
		ret = append(ret, msg)
	})
	if err != nil {
		return nil, err
	}

	if data == nil {
		err := c.nc.PublishRequest(subject, sub, nil)
		if err != nil {
			return nil, err
		}
	} else {
		err := c.nc.PublishRequest(subject, sub, *data)
		if err != nil {
			return nil, err
		}
	}

	time.Sleep(timeout)
	err = s.Drain()
	if err != nil {
		return nil, err
	}

	return ret, nil
}
