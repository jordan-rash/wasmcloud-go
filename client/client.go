package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jordan-rash/wasmcloud-go/broker"
	"github.com/nats-io/nats.go"
)

type client struct {
	nc       *nats.Conn
	nsprefix string
	timeout  time.Duration
}

func New(nc *nats.Conn, prefix string, timeout time.Duration) client {
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
	fmt.Println(subject)
	hostsRaw := printResults(c.nc, subject, nil)
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
	fmt.Println(subject)

	hoststatus := printResults(c.nc, subject, nil)
	hs := hostStatus{}

	json.Unmarshal([]byte(hoststatus[0]), &hs)

	return hs
}

// NATs topic: get.claims
func (c client) GetClaims() {
	subject := broker.Queries{}.Claims(c.nsprefix)
	fmt.Println(subject)
	printResults(c.nc, subject, nil)
}

// NATs topic: auction.actor
func (c client) PerformActorAuction(actorRef string, constraints map[string]interface{}, timeout time.Duration) {
	subject := broker.ActorAuctionSubject(c.nsprefix)
	fmt.Println(subject)
	data := struct {
		ActorRef    string                 `json:"actor_ref"`
		Constraints map[string]interface{} `json:"constraints"`
	}{
		ActorRef:    actorRef,
		Constraints: constraints,
	}
	b_data, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	printResults(c.nc, subject, &b_data)
}

// NATs topic: auction.provider
func (c client) PerformProviderAuction(providerRef string, linkName string, constraints map[string]interface{}, timeout time.Duration) {
	subject := broker.ProviderAuctionSubject(c.nsprefix)
	fmt.Println(subject)
	data := struct {
		ProviderRef string                 `json:"provider_ref"`
		LinkName    string                 `json:"link_name"`
		Constraints map[string]interface{} `json:"constraints"`
	}{
		ProviderRef: providerRef,
		Constraints: constraints,
		LinkName:    linkName,
	}
	b_data, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	printResults(c.nc, subject, &b_data)
}

// NATs topic: cmd.{host}.la
func (c client) StartActor(hostID string, actorRef string, count int, annotations map[string]interface{}) {
	subject := broker.Commands{}.StartActor(c.nsprefix, hostID)
	fmt.Println(subject)
	data := struct {
		ActorRef    string                 `json:"actor_ref"`
		HostID      string                 `json:"host_id"`
		Count       int                    `json:"count"`
		Annotations map[string]interface{} `json:"annotations"`
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

	printResults(c.nc, subject, &b_data)
}

// NATs topic: cmd.{host}.lp
func (c client) StartProvider(hostID string, providerRef string, linkName string, annotations map[string]interface{}, providerConfiguration string) {
	subject := broker.Commands{}.StartProvider(c.nsprefix, hostID)
	fmt.Println(subject)
	data := struct {
		ProviderRef   string                 `json:"provider_ref"`
		HostID        string                 `json:"host_id"`
		LinkName      string                 `json:"link_name"`
		Annotations   map[string]interface{} `json:"annotations"`
		Configuration string                 `json:"configuration"`
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

	printResults(c.nc, subject, &b_data)
}

// NATs topic: linkdefs.put
func (c client) AdvertiseLink(actorID string, providerID string, contractID string, linkName string, values map[string]interface{}) {
	subject := broker.AdvertiseLink(c.nsprefix)
	fmt.Println(subject)
	data := struct {
		ActorID    string                 `json:"actor_id"`
		ProviderID string                 `json:"provider_id"`
		ContractID string                 `json:"contract_id"`
		LinkName   string                 `json:"link_name"`
		Value      map[string]interface{} `json:"values"`
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

	printResults(c.nc, subject, &b_data)

}

// NATs topic: linkdefs.del
func (c client) RemoveLink(actorID string, contractID string, linkName string) {}

// NATs topic: get.links
func (c client) QueryLinks() {}

// NATs topic: cmd.{host}.upd
func (c client) UpdateActor(hostID string, existingActorID string, newActorRef string, annotations map[string]string) {
}

// NATs topic: cmd.{host}.sa
func (c client) StopActor(hostID string, actorRef string, count int, annotations map[string]string) {}

// NATs topic: cmd.{host}.sp
func (c client) StopProvider(hostID string, providerRef string, linkName string, contractID string, annotations map[string]string) {
}

// NATs topic: cmd.{host}.stop
func (c client) StopHost(hostID string, timeout time.Duration) {}

func printResults(nc *nats.Conn, subject string, data *[]byte) []string {
	timeout := time.Second * 1
	sub := nats.NewInbox()
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

	ch := make(chan *nats.Msg)
	s, err := nc.ChanSubscribe(sub, ch)
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
}
