package client

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-logr/logr"
	"github.com/jordan-rash/wasmcloud-go/broker"
	"github.com/jordan-rash/wasmcloud-go/kv"
	"github.com/jordan-rash/wasmcloud-go/models"
	"github.com/nats-io/nats.go"

	core "github.com/wasmcloud/interfaces/core/tinygo"
)

// Lattice control interface client
type Client struct {
	nc             *nats.Conn
	topicPrefix    string
	nsPrefix       string
	timeout        time.Duration
	auctionTimeout time.Duration
	jsDomain       string
	kvstore        nats.KeyValue
	logger         logr.Logger
}

// Deprecated: Use `client.New() ClientBuilder` instead.
func New_Old(nc *nats.Conn, prefix string, timeout time.Duration) Client {
	return Client{
		nc:          nc,
		topicPrefix: prefix,
		timeout:     timeout,
	}
}

// NATs topic: ping.hosts
func (c Client) GetHosts(timeout time.Duration) (*models.Hosts, error) {
	hosts := models.Hosts{}

	subject := broker.Queries{}.Hosts(c.nsPrefix)
	msgs, err := c.CollectTimeout(subject, nil, &timeout)
	if err != nil {
		return nil, err
	}

	for _, m := range msgs {
		if m.Header.Get("Status") == "503" {
			return nil, nats.ErrNoResponders
		}
		tHost := models.Host{}
		err := json.Unmarshal(m.Data, &tHost)
		if err != nil {
			return nil, err
		}
		hosts = append(hosts, tHost)
	}

	return &hosts, nil
}

// NATs topic: get.{host}.inv
func (c Client) GetHostInventory(hostId string) (*models.HostInventory, error) {
	subject := broker.Queries{}.HostInventory(c.nsPrefix, hostId)
	hostInventoryRaw, err := c.CollectTimeout(subject, nil, nil)
	if err != nil {
		return nil, err
	}

	if len(hostInventoryRaw) < 1 {
		return nil, errors.New("did not find host status")
	}

	for _, h := range hostInventoryRaw {
		tHostInventory := models.HostInventory{}
		err := json.Unmarshal(h.Data, &tHostInventory)
		if err != nil {
			return nil, err
		}
		if tHostInventory.HostId == hostId {
			return &tHostInventory, nil
		}
	}

	return nil, errors.New("did not find host status")
}

// NATs topic: get.claims
func (c Client) GetClaims() (*models.Claims, error) {
	if c.kvstore != nil {
		ret := models.Claims{}
		resp, err := kv.GetClaims(c.kvstore)
		if err != nil {
			return nil, err
		}
		claims := resp.Claims
		for _, kvm := range claims {
			for _, v := range kvm {
				tClaim := models.Claim{}
				err := json.Unmarshal([]byte(v), &tClaim)
				if err != nil {
					return nil, err
				}
				ret.Claims = append(ret.Claims, tClaim)
			}
		}
		return &ret, nil
	} else {
		subject := broker.Queries{}.Claims(c.nsPrefix)

		claims := models.Claims{}
		claimsRaw, err := c.CollectTimeout(subject, nil, nil)
		if err != nil {
			return nil, err
		}

		for _, c := range claimsRaw {
			tClaim := models.Claim{}
			err := json.Unmarshal(c.Data, &tClaim)
			if err != nil {
				return nil, err
			}
			claims.Claims = append(claims.Claims, tClaim)
		}

		return &claims, nil
	}
}

// NATs topic: auction.actor
func (c Client) PerformActorAuction(actorRef string, constraints map[string]string, timeout time.Duration) ([]*models.ActorAucutionAck, error) {
	subject := broker.ActorAuctionSubject(c.nsPrefix)
	data := models.ActorAuctionRequest{
		ActorRef:    actorRef,
		Constraints: constraints,
	}
	data_bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	acks := []*models.ActorAucutionAck{}
	rawAcks, err := c.CollectTimeout(subject, &data_bytes, &timeout)
	if err != nil {
		return nil, err
	}
	for _, a := range rawAcks {
		tAck := models.ActorAucutionAck{}
		err := json.Unmarshal(a.Data, &tAck)
		if err != nil {
			return nil, err
		}
		acks = append(acks, &tAck)
	}

	return acks, nil

}

// NATs topic: auction.provider
func (c Client) PerformProviderAuction(providerRef string, linkName string, constraints map[string]string, timeout time.Duration) ([]*models.ProviderAuctionAck, error) {
	subject := broker.ProviderAuctionSubject(c.nsPrefix)
	data := models.ProviderAuctionRequest{
		ProviderRef: providerRef,
		Constraints: constraints,
		LinkName:    linkName,
	}

	data_bytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	acks := []*models.ProviderAuctionAck{}
	rawAcks, err := c.CollectTimeout(subject, &data_bytes, &timeout)
	if err != nil {
		return nil, err
	}
	for _, a := range rawAcks {
		pAck := models.ProviderAuctionAck{}
		err := json.Unmarshal(a.Data, &pAck)
		if err != nil {
			return nil, err
		}
		acks = append(acks, &pAck)
	}
	return acks, nil
}

// NATs topic: cmd.{host}.la
func (c Client) StartActor(hostID string, actorRef string, count uint16, annotations map[string]string) (*models.CtlOperationAck, error) {
	subject := broker.Commands{}.StartActor(c.nsPrefix, hostID)
	data := models.StartActorCommand{
		ActorRef:    actorRef,
		HostId:      hostID,
		Count:       count,
		Annotations: annotations,
	}

	data_bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	rawAcks, err := c.CollectTimeout(subject, &data_bytes, nil)
	if err != nil {
		return nil, err
	}
	for _, a := range rawAcks {
		cAck := models.CtlOperationAck{}
		err := json.Unmarshal(a.Data, &cAck)
		if err != nil {
			return nil, err
		}
		if cAck.Accepted {
			return &cAck, nil
		}
	}
	return nil, errors.New("did not receive ack")
}

// NATs topic: cmd.{host}.lp
func (c Client) StartProvider(hostID string, providerRef string, linkName string, annotations map[string]string, providerConfiguration string) (*models.CtlOperationAck, error) {
	subject := broker.Commands{}.StartProvider(c.nsPrefix, hostID)
	startCmd := models.StartProviderCommand{
		ProviderRef:   providerRef,
		HostId:        hostID,
		LinkName:      linkName,
		Annotations:   annotations,
		Configuration: providerConfiguration,
	}

	if hostID == "" {
		subject := broker.ProviderAuctionSubject(c.nsPrefix)
		proReq := models.ProviderAuctionRequest{
			ProviderRef: providerRef,
			LinkName:    linkName,
			// TODO: where are these contrainsts??
			Constraints: models.ConstraintMap{},
		}
		data_bytes, err := json.Marshal(proReq)
		if err != nil {
			return nil, err
		}
		acks, err := c.CollectTimeout(subject, &data_bytes, nil)
		if err != nil {
			return nil, err
		}
		if len(acks) < 1 {
			return nil, errors.New("no host detected to start provider")
		}

		for _, a := range acks {
			tAck := models.ProviderAuctionAck{}
			err := json.Unmarshal(a.Data, &tAck)
			if err != nil {
				return nil, err
			}
			startCmd.HostId = tAck.HostId
		}
	}

	data_bytes, err := json.Marshal(startCmd)
	if err != nil {
		return nil, err
	}

	rawAcks, err := c.CollectTimeout(subject, &data_bytes, nil)
	if err != nil {
		return nil, err
	}
	for _, a := range rawAcks {
		cAck := models.CtlOperationAck{}
		err := json.Unmarshal(a.Data, &cAck)
		if err != nil {
			return nil, err
		}
		if cAck.Accepted {
			return &cAck, nil
		}
	}

	return nil, errors.New("did not receive ack")
}

// NATs topic: linkdefs.put
func (c Client) AdvertiseLink(actorID string, providerID string, contractID string, linkName string, values map[string]string) (*models.CtlOperationAck, error) {
	ld := core.LinkDefinition{
		ActorId:    actorID,
		ProviderId: providerID,
		ContractId: contractID,
		LinkName:   linkName,
		Values:     values,
	}

	if c.kvstore != nil {
		err := kv.PutLink(c.kvstore, ld)
		if err != nil {
			return nil, err
		}
		return &models.CtlOperationAck{Accepted: true, Error: ""}, nil
	} else {

		subject := broker.AdvertiseLink(c.nsPrefix)

		data_bytes, err := json.Marshal(ld)
		if err != nil {
			panic(err)
		}

		rawAcks, err := c.CollectTimeout(subject, &data_bytes, nil)
		if err != nil {
			return nil, err
		}
		for _, a := range rawAcks {
			cAck := models.CtlOperationAck{}
			err := json.Unmarshal(a.Data, &cAck)
			if err != nil {
				return nil, err
			}
			if cAck.Accepted {
				return &cAck, nil

			}
		}

		return nil, errors.New("did not receive ack")
	}
}

// NATs topic: linkdefs.del
func (c Client) RemoveLink(actorID string, contractID string, linkName string) (*models.CtlOperationAck, error) {
	removeLinkReq := models.RemoveLinkDefinationRequest{
		ActorId:    actorID,
		ContractId: contractID,
		LinkName:   linkName,
	}

	if c.kvstore != nil {
		err := kv.DeleteLink(c.kvstore, removeLinkReq)
		if err != nil {
			return nil, err
		}
		return &models.CtlOperationAck{Accepted: true, Error: ""}, nil
	} else {

		subject := broker.RemoveLink(c.nsPrefix)
		data_bytes, err := json.Marshal(removeLinkReq)
		if err != nil {
			panic(err)
		}
		rawAcks, err := c.CollectTimeout(subject, &data_bytes, nil)
		if err != nil {
			return nil, err
		}
		for _, a := range rawAcks {
			cAck := models.CtlOperationAck{}
			err := json.Unmarshal(a.Data, &cAck)
			if err != nil {
				return nil, err
			}
			if cAck.Accepted {
				return &cAck, nil
			}
		}
		return nil, errors.New("did not receive ack")
	}
}

// NATs topic: get.links
func (c Client) QueryLinks() (*models.LinkDefinitionList, error) {
	if c.kvstore != nil {
		return kv.GetLinks(c.kvstore)
	} else {

		subject := broker.Queries{}.LinkDefinitions(c.nsPrefix)
		rawLinks, err := c.CollectTimeout(subject, nil, nil)
		if err != nil {
			return nil, err
		}

		ld := models.LinkDefinitionList{}
		for _, r := range rawLinks {
			tLink := core.LinkDefinition{}
			err := json.Unmarshal(r.Data, &tLink)
			if err != nil {
				return nil, err
			}

			ld.Links = append(ld.Links, tLink)
		}

		return &ld, nil
	}
}

// NATs topic: cmd.{host}.upd
func (c Client) UpdateActor(hostID string, existingActorID string, newActorRef string, annotations map[string]string) (*models.CtlOperationAck, error) {
	subject := broker.Commands{}.UpdateActor(c.nsPrefix, hostID)
	data := models.UpdateActorCommand{
		ActorRef:    existingActorID,
		Annotations: annotations,
		HostId:      hostID,
		NewActorRef: newActorRef,
	}

	data_bytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	rawAcks, err := c.CollectTimeout(subject, &data_bytes, nil)
	if err != nil {
		return nil, err
	}
	for _, a := range rawAcks {
		cAck := models.CtlOperationAck{}
		err := json.Unmarshal(a.Data, &cAck)
		if err != nil {
			return nil, err
		}
		if cAck.Accepted {
			return &cAck, nil
		}
	}

	return nil, errors.New("did not receive ack")
}

// NATs topic: cmd.{host}.sa
func (c Client) StopActor(hostID string, actorRef string, count uint16, annotations map[string]string) (*models.CtlOperationAck, error) {
	subject := broker.Commands{}.StopActor(c.nsPrefix, hostID)
	data := models.StopActorCommand{
		HostId:      hostID,
		ActorRef:    actorRef,
		Count:       count,
		Annotations: annotations,
	}

	data_bytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	rawAcks, err := c.CollectTimeout(subject, &data_bytes, nil)
	if err != nil {
		return nil, err
	}
	for _, a := range rawAcks {
		cAck := models.CtlOperationAck{}
		err := json.Unmarshal(a.Data, &cAck)
		if err != nil {
			return nil, err
		}
		if cAck.Accepted {
			return &cAck, nil
		}
	}

	return nil, errors.New("did not receive ack")
}

// NATs topic: cmd.{host}.sp
func (c Client) StopProvider(hostID string, providerRef string, linkName string, contractID string, annotations map[string]string) (*models.CtlOperationAck, error) {
	subject := broker.Commands{}.StopProvider(c.nsPrefix, hostID)
	data := models.StopProviderCommand{
		HostId:      hostID,
		ProviderRef: providerRef,
		LinkName:    linkName,
		ContractId:  contractID,
		Annotations: annotations,
	}

	data_bytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	rawAcks, err := c.CollectTimeout(subject, &data_bytes, nil)
	if err != nil {
		return nil, err
	}
	for _, a := range rawAcks {
		cAck := models.CtlOperationAck{}
		err := json.Unmarshal(a.Data, &cAck)
		if err != nil {
			return nil, err
		}
		if cAck.Accepted {
			return &cAck, nil
		}
	}

	return nil, errors.New("did not receive ack")
}

// NATs topic: cmd.{host}.stop
func (c Client) StopHost(hostID string, timeout uint64) (*models.CtlOperationAck, error) {
	subject := broker.Commands{}.StopHost(c.nsPrefix, hostID)
	data := models.StopHostCommand{
		HostId:  hostID,
		Timeout: timeout,
	}

	data_bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	rawAcks, err := c.CollectTimeout(subject, &data_bytes, nil)
	if err != nil {
		return nil, err
	}
	for _, a := range rawAcks {
		cAck := models.CtlOperationAck{}
		err := json.Unmarshal(a.Data, &cAck)
		if err != nil {
			return nil, err
		}
		if cAck.Accepted {
			return &cAck, nil
		}
	}
	return nil, errors.New("did not receive ack")
}
