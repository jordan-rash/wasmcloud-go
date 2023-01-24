package client

import (
	"time"

	"github.com/nats-io/nats.go"
)

/// Collect results until timeout has elapsed
func (c Client) CollectTimeout(subject string, data *[]byte, timeoutOverride *time.Duration) []string {
	timeout := c.timeout
	if timeoutOverride != nil {
		timeout = *timeoutOverride
	}
	sub := nats.NewInbox()
	ch := make(chan *nats.Msg)
	s, err := c.nc.ChanSubscribe(sub, ch)
	if err != nil {
		panic(err)
	}

	err = c.nc.PublishRequest(subject, sub, *data)
	if err != nil {
		panic(err)
	}

	var ret []string
	for {
		select {
		case msg := <-ch:
			ret = append(ret, (string(msg.Data)))
		case <-time.After(timeout):
			s.Unsubscribe()
			s.Drain()
			return ret
		}
	}
	return nil
}
