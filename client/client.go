package client

import (
	"encoding/json"
	"os"
	"time"

	"github.com/gobuffalo/envy"
	"github.com/jordan-rash/wasmcloud-go/broker"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	switch envy.Get("LOG_LVL", "ERROR") {
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "TRACE":
		log.SetLevel(log.TraceLevel)
	default:
		log.SetLevel(log.ErrorLevel)
	}
}

type client struct {
	nc       *nats.Conn
	nsprefix string
	timeout  time.Duration
}

// Deprecated: Use `client.New() ClientBuilder` instead.
func New_Old(nc *nats.Conn, prefix string, timeout time.Duration) client {
	return client{
		nc,
		prefix,
		timeout,
	}
}

// NATs topic: ping.hosts
func (c client) GetHosts(timeout time.Duration) []host {
	var hosts []host

	subject := broker.Queries{}.Hosts(c.nsprefix)
	log.Debug(subject)
	hostsRaw := c.printResults(c.nc, subject, nil, &timeout)
	for _, h := range hostsRaw {
		tHost := host{}
		json.Unmarshal([]byte(h), &tHost)
		hosts = append(hosts, tHost)
	}

	return hosts
}

// NATs topic: get.{host}.inv
func (c client) GetHostInventory(hostId string) hostStatus {
	subject := broker.Queries{}.HostInventory(c.nsprefix, hostId)
	log.Debug(subject)

	hoststatus := c.printResults(c.nc, subject, nil, nil)
	hs := hostStatus{}

	if len(hoststatus) < 1 {
		log.Error("Did not find host status")
		return hostStatus{}
	}
	json.Unmarshal([]byte(hoststatus[0]), &hs)

	return hs
}

// NATs topic: get.claims
func (c client) GetClaims() claims {
	subject := broker.Queries{}.Claims(c.nsprefix)
	log.Debug(subject)

	claims := claims{}
	claimsRaw := c.printResults(c.nc, subject, nil, nil)
	json.Unmarshal([]byte(claimsRaw[0]), &claims)

	return claims
}

// NATs topic: auction.actor
func (c client) PerformActorAuction(actorRef string, constraints map[string]string, timeout time.Duration) []string {
	subject := broker.ActorAuctionSubject(c.nsprefix)
	log.Debug(subject)
	data := struct {
		ActorRef    string            `json:"actor_ref"`
		Constraints map[string]string `json:"constraints,omitempty"`
	}{
		ActorRef:    actorRef,
		Constraints: constraints,
	}
	b_data, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return c.printResults(c.nc, subject, &b_data, &timeout)
}

// NATs topic: auction.provider
func (c client) PerformProviderAuction(providerRef string, linkName string, constraints map[string]string, timeout time.Duration) []string {
	subject := broker.ProviderAuctionSubject(c.nsprefix)
	log.Debug(subject)
	data := struct {
		ProviderRef string            `json:"provider_ref"`
		LinkName    string            `json:"link_name"`
		Constraints map[string]string `json:"constraints,omitempty"`
	}{
		ProviderRef: providerRef,
		Constraints: constraints,
		LinkName:    linkName,
	}
	b_data, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return c.printResults(c.nc, subject, &b_data, &timeout)
}

// NATs topic: cmd.{host}.la
func (c client) StartActor(hostID string, actorRef string, count int, annotations map[string]string) []string {
	subject := broker.Commands{}.StartActor(c.nsprefix, hostID)
	log.Debug(subject)
	data := struct {
		ActorRef    string            `json:"actor_ref"`
		HostID      string            `json:"host_id"`
		Count       int               `json:"count"`
		Annotations map[string]string `json:"annotations,omitempty"`
	}{
		ActorRef:    actorRef,
		HostID:      hostID,
		Count:       count,
		Annotations: annotations,
	}

	b_data, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return c.printResults(c.nc, subject, &b_data, nil)
}

// NATs topic: cmd.{host}.lp
func (c client) StartProvider(hostID string, providerRef string, linkName string, annotations map[string]string, providerConfiguration string) []string {
	subject := broker.Commands{}.StartProvider(c.nsprefix, hostID)
	log.Debug(subject)
	data := struct {
		ProviderRef   string            `json:"provider_ref"`
		HostID        string            `json:"host_id"`
		LinkName      string            `json:"link_name"`
		Annotations   map[string]string `json:"annotations,omitempty"`
		Configuration string            `json:"configuration"`
	}{
		ProviderRef:   providerRef,
		HostID:        hostID,
		LinkName:      linkName,
		Annotations:   annotations,
		Configuration: providerConfiguration,
	}

	b_data, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return c.printResults(c.nc, subject, &b_data, nil)
}

// NATs topic: linkdefs.put
func (c client) AdvertiseLink(actorID string, providerID string, contractID string, linkName string, values map[string]string) []string {
	subject := broker.AdvertiseLink(c.nsprefix)
	log.Debug(subject)
	data := struct {
		ActorID    string            `json:"actor_id"`
		ProviderID string            `json:"provider_id"`
		ContractID string            `json:"contract_id"`
		LinkName   string            `json:"link_name"`
		Value      map[string]string `json:"values,omitempty"`
	}{
		ActorID:    actorID,
		ProviderID: providerID,
		ContractID: contractID,
		LinkName:   linkName,
		Value:      values,
	}

	b_data, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return c.printResults(c.nc, subject, &b_data, nil)

}

// NATs topic: linkdefs.del
func (c client) RemoveLink(actorID string, contractID string, linkName string) []string {
	subject := broker.RemoveLink(c.nsprefix)
	log.Debug(subject)
	data := struct {
		ActorID    string `json:"actor_id"`
		ContractID string `json:"contract_id"`
		LinkName   string `json:"link_name"`
	}{
		ActorID:    actorID,
		ContractID: contractID,
		LinkName:   linkName,
	}

	b_data, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return c.printResults(c.nc, subject, &b_data, nil)
}

// NATs topic: get.links
func (c client) QueryLinks() []string {
	subject := broker.Queries{}.LinkDefinitions(c.nsprefix)
	log.Debug(subject)
	return c.printResults(c.nc, subject, nil, nil)
}

// NATs topic: cmd.{host}.upd
func (c client) UpdateActor(hostID string, existingActorID string, newActorRef string, annotations map[string]string) []string {
	subject := broker.Commands{}.UpdateActor(c.nsprefix, hostID)
	log.Debug(subject)
	data := struct {
		HostID          string            `json:"host_id"`
		ExistingActorID string            `json:"actor_id"`
		NewActorRef     string            `json:"new_actor_ref"`
		Annotations     map[string]string `json:"annotations,omitempty"`
	}{
		HostID:          hostID,
		ExistingActorID: existingActorID,
		NewActorRef:     newActorRef,
		Annotations:     annotations,
	}

	b_data, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return c.printResults(c.nc, subject, &b_data, nil)
}

// NATs topic: cmd.{host}.sa
func (c client) StopActor(hostID string, actorRef string, count int, annotations map[string]string) []string {
	subject := broker.Commands{}.StopActor(c.nsprefix, hostID)
	log.Debug(subject)
	data := struct {
		HostID      string            `json:"host_id"`
		ActorRef    string            `json:"actor_ref"`
		Count       int               `json:"count"`
		Annotations map[string]string `json:"annotations,omitempty"`
	}{
		HostID:      hostID,
		ActorRef:    actorRef,
		Count:       count,
		Annotations: annotations,
	}

	b_data, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return c.printResults(c.nc, subject, &b_data, nil)
}

// NATs topic: cmd.{host}.sp
func (c client) StopProvider(hostID string, providerRef string, linkName string, contractID string, annotations map[string]string) []string {
	subject := broker.Commands{}.StopProvider(c.nsprefix, hostID)
	log.Debug(subject)
	data := struct {
		HostID      string            `json:"host_id"`
		ProviderRef string            `json:"provider_ref"`
		LinkName    string            `json:"link_name"`
		ContractID  string            `json:"contract_id"`
		Annotations map[string]string `json:"annotations,omitempty"`
	}{
		HostID:      hostID,
		ProviderRef: providerRef,
		LinkName:    linkName,
		ContractID:  contractID,
		Annotations: annotations,
	}

	b_data, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return c.printResults(c.nc, subject, &b_data, nil)
}

// NATs topic: cmd.{host}.stop
func (c client) StopHost(hostID string, timeout time.Duration) []string {
	subject := broker.Commands{}.StopHost(c.nsprefix, hostID)
	log.Debug(subject)

	data := struct {
		HostID  string `json:"host_id"`
		Timeout int64  `json:"timeout"`
	}{
		HostID:  hostID,
		Timeout: timeout.Milliseconds(),
	}
	b_data, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return c.printResults(c.nc, subject, &b_data, &timeout)
}

func (c client) printResults(nc *nats.Conn, subject string, data *[]byte, timeoutOverride *time.Duration) []string {
	timeout := c.timeout
	if timeoutOverride != nil {
		timeout = *timeoutOverride
	}
	sub := nats.NewInbox()
	ch := make(chan *nats.Msg)
	s, err := nc.ChanSubscribe(sub, ch)
	if err != nil {
		panic(err)
	}

	if data == nil {
		err := nc.PublishRequest(subject, sub, nil)
		if err != nil {
			panic(err)
		}
	} else {
		err := nc.PublishRequest(subject, sub, *data)
		if err != nil {
			panic(err)
		}
	}

	var ret []string
	for {
		select {
		case msg := <-ch:
			ret = append(ret, (string(msg.Data)))
		case <-time.After(timeout):
			s.Unsubscribe()
			s.Drain()
			if envy.Get("PRETTY_PRINT", "false") == "true" {
				PrettyPrint(ret)
			}
			return ret
		}
	}
}

// this is temporary
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		log.Println(string(b))
	}
	return
}
