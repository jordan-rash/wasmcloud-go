package kv

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jordan-rash/wasmcloud-go/models"
	"github.com/nats-io/nats.go"
	core "github.com/wasmcloud/interfaces/core/tinygo"
)

const (
	LINKDEF_PREFIX = "LINKDEF_"
	CLAIMS_PREFIX  = "CLAIMS_"
	BUCKET_PREFIX  = "LATTICEDATA_"
)

var (
	DEFAULT_LEAFNODE_OPTIONS nats.JSOpt = nats.PublishAsyncMaxPending(256)
)

func GetKVStore(nc *nats.Conn, latticePrefix, jsDomain string) (nats.KeyValue, error) {
	var js nats.JetStreamContext
	var err error

	if jsDomain != "" {
		js, err = nc.JetStream(nats.Domain(jsDomain), DEFAULT_LEAFNODE_OPTIONS)
		if err != nil {
			return nil, err
		}
	} else {
		js, err = nc.JetStream(DEFAULT_LEAFNODE_OPTIONS)
		if err != nil {
			return nil, err
		}
	}

	bucket := fmt.Sprintf("%s%s", BUCKET_PREFIX, latticePrefix)

	return js.KeyValue(bucket)
}

func GetClaims(store nats.KeyValue) (*models.GetClaimsResponse, error) {
	claims := models.CtlKVList{}
	entries, err := store.Keys()
	if err != nil {
		return nil, err
	}

	for _, c := range entries {
		if strings.HasPrefix(c, CLAIMS_PREFIX) {
			entry, err := store.Get(c)
			if err != nil {
				return nil, err
			}

			d := models.KeyValueMap{}
			err = json.Unmarshal(entry.Value(), &d)
			if err != nil {
				return nil, err
			}

			claims = append(claims, d)
		}
	}

	return &models.GetClaimsResponse{Claims: claims}, nil
}

func GetLinks(store nats.KeyValue) (*models.LinkDefinitionList, error) {
	links := core.ActorLinks{}
	entries, err := store.Keys()
	if err != nil {
		return nil, err
	}

	for _, c := range entries {
		if strings.HasPrefix(c, LINKDEF_PREFIX) {
			entry, err := store.Get(c)
			if err != nil {
				return nil, err
			}

			d := core.LinkDefinition{}
			err = json.Unmarshal(entry.Value(), &d)
			if err != nil {
				return nil, err
			}

			links = append(links, d)
		}
	}

	return &models.LinkDefinitionList{Links: links}, nil
}

func PutLink(store nats.KeyValue, ld core.LinkDefinition) error {
	id, err := LDHash(&ld)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s%s", LINKDEF_PREFIX, id)
	ldBytes, err := json.Marshal(ld)
	if err != nil {
		return err
	}

	_, err = store.Put(key, ldBytes)
	return err
}

func DeleteLink(store nats.KeyValue, linkRemoveRequest models.RemoveLinkDefinationRequest) error {
	rawHash, err := LDHashRaw(linkRemoveRequest.ActorId, linkRemoveRequest.ContractId, linkRemoveRequest.LinkName)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s%s", LINKDEF_PREFIX, rawHash)
	return store.Delete(key)
}

func LDHash(ld *core.LinkDefinition) (string, error) {
	return LDHashRaw(ld.ActorId, ld.ContractId, ld.LinkName)
}

// Performs a hash function against the link definition key fields. The corresponding
// Elixir hash function can be found in https://github.com/wasmcloud/wasmcloud-otp/ in the
// host_core/lib/linkdefs/manager.ex file, which uses Erlang's :crypto
func LDHashRaw(actorId, contractId, linkName string) (string, error) {
	var buf bytes.Buffer
	_, err := buf.WriteString(actorId)
	if err != nil {
		return "", err
	}
	_, err = buf.WriteString(contractId)
	if err != nil {
		return "", err
	}
	_, err = buf.WriteString(linkName)
	if err != nil {
		return "", err
	}

	digest := sha256.Sum256(buf.Bytes())
	return strings.ToUpper(fmt.Sprintf("%x", digest)), nil
}
