package broker

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const WASMCLOUD_DEFAULT_NSPREFIX string = "default"

type Commands struct{}
type Queries struct{}
type Event struct{}

func prefix(nsprefix string) string {
	if nsprefix == "" {
		nsprefix = "default"
	}
	return fmt.Sprintf("wasmbus.ctl.%s", nsprefix)
}

func ControlEvent(nsprefix string) string {
	if nsprefix == "" {
		nsprefix = "default"
	}
	return fmt.Sprintf("wasmbus.evt.%s", nsprefix)
}

func ProviderAuctionSubject(nsprefix string) string {
	return fmt.Sprintf("%s.auction.provider", prefix(nsprefix))
}
func ActorAuctionSubject(nsprefix string) string {
	return fmt.Sprintf("%s.auction.actor", prefix(nsprefix))
}
func AdvertiseLink(nsprefix string) string {
	return fmt.Sprintf("%s.linkdefs.put", prefix(nsprefix))
}
func RemoveLink(nsprefix string) string {
	return fmt.Sprintf("%s.linkdefs.del", prefix(nsprefix))
}

func (Commands) StartActor(nsprefix, host string) string {
	return fmt.Sprintf("%s.cmd.%s.la", prefix(nsprefix), host)
}
func (Commands) StopActor(nsprefix, host string) string {
	return fmt.Sprintf("%s.cmd.%s.sa", prefix(nsprefix), host)
}
func (Commands) StartProvider(nsprefix, host string) string {
	return fmt.Sprintf("%s.cmd.%s.lp", prefix(nsprefix), host)
}
func (Commands) StopProvider(nsprefix, host string) string {
	return fmt.Sprintf("%s.cmd.%s.sp", prefix(nsprefix), host)
}
func (Commands) UpdateActor(nsprefix, host string) string {
	return fmt.Sprintf("%s.cmd.%s.upd", prefix(nsprefix), host)
}
func (Commands) StopHost(nsprefix, host string) string {
	return fmt.Sprintf("%s.cmd.%s.stop", prefix(nsprefix), host)
}

func (Queries) LinkDefinitions(nsprefix string) string {
	return fmt.Sprintf("%s.get.links", prefix(nsprefix))
}
func (Queries) Claims(nsprefix string) string {
	return fmt.Sprintf("%s.get.claims", prefix(nsprefix))
}
func (Queries) HostInventory(nsprefix, host string) string {
	return fmt.Sprintf("%s.get.%s.inv", prefix(nsprefix), host)
}
func (Queries) Hosts(nsprefix string) string {
	return fmt.Sprintf("%s.ping.hosts", prefix(nsprefix))
}

func (Event) HostStart(nsprefix, hostid string) (string, []byte) {
	ce := cloudevent{
		Spec:        "1.0",
		Id:          uuid.New().String(),
		Source:      hostid,
		Type:        "com.wasmcloud.lattice.host_started",
		ContentType: "application/json",
		Data: cedata{
			Actors:    []string{},
			Providers: []string{},
			Friendly:  "yolo-bro-1234",
			Labels: map[string]string{
				"hostcore.arch":     "armv61",
				"hostcore.os":       "raspbian",
				"hostcore.osfamily": "unix",
				"hostcore.version":  "v0.0.0-wasmcloud_go",
				"hostcore.runtime":  "wazero",
			},
			Version: "v0.0.0-wasmcloud_go",
		},
	}
	ceb, _ := json.Marshal(ce)
	return ControlEvent(nsprefix), ceb
}

func (Event) HostHeartbeat(nsprefix, hostid string, startTime time.Time) (string, []byte) {
	ce := cloudevent{
		Spec:        "1.0",
		Id:          uuid.New().String(),
		Source:      hostid,
		Type:        "com.wasmcloud.lattice.host_heartbeat",
		ContentType: "application/json",
		Data: cedata{
			Actors:    []string{},
			Providers: []string{},
			Friendly:  "yolo-bro-1234",
			Labels: map[string]string{
				"hostcore.arch":     "armv61",
				"hostcore.os":       "raspbian",
				"hostcore.osfamily": "unix",
				"hostcore.version":  "v0.0.0-wasmcloud_go",
				"hostcore.runtime":  "wazero",
			},
			Version: "v0.0.0-wasmcloud_go",
		},
	}

	ceb, _ := json.Marshal(ce)
	return ControlEvent(nsprefix), ceb
}

type cloudevent struct {
	Spec        string `json:"specversion"`
	Id          string `json:"id"`
	Source      string `json:"source"`
	Type        string `json:"type"`
	ContentType string `json:"datacontenttype"`
	Data        cedata `json:"data"`
}

type cedata struct {
	Actors    []string          `json:"actors"`
	Providers []string          `json:"providers"`
	Friendly  string            `json:"friendly_name"`
	Labels    map[string]string `json:"labels"`
	Version   string            `json:"version"`
}
