package client

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
)

/// Collect results until timeout has elapsed
func (c Client) CollectTimeout(subject string, data *[]byte, timeoutOverride *time.Duration) ([]*nats.Msg, error) {
	timeout := c.timeout
	if timeoutOverride != nil {
		timeout = *timeoutOverride
	}

	tempInbox := c.nc.NewRespInbox()
	msgs := make(chan (*nats.Msg))
	sub, err := c.nc.ChanSubscribe(tempInbox, msgs)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if data == nil {
		err := c.nc.PublishRequest(subject, tempInbox, nil)
		if err != nil {
			return nil, err
		}
	} else {
		err := c.nc.PublishRequest(subject, tempInbox, *data)
		if err != nil {
			return nil, err
		}
	}

	var ret []*nats.Msg

	for {
		select {
		case <-ctx.Done():
			err = sub.Drain()
			if err != nil {
				return nil, err
			}
			return ret, nil
		case msg := <-msgs:
			ret = append(ret, msg)
		}
	}
}
