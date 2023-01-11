package broker

import (
	"fmt"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
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

func (Event) HostStart(nsprefix, hostid string) (string, cloudevents.Event) {
	event := cloudevents.NewEvent()
	event.SetData(cloudevents.ApplicationJSON, map[string]interface{}{"friendly_name": "yolo-bro-1234", "labels": map[string]string{"hostcore.arch": "armv61", "hostcore.os": "raspbian", "hostcore.osfamily": "unix", "hostcore.version": "v0.0.0-wasmcloud_go", "hostcore.runtime": "wazero"}, "version": "v0.0.0-wasmcloud_go"})
	event.SetID(uuid.New().String())
	event.SetSource(hostid) // "N" key
	event.SetType("com.wasmcloud.lattice.host_started")

	return ControlEvent(nsprefix), event
}

func (Event) HostHeartbeat(nsprefix, hostid string, startTime time.Time) (string, cloudevents.Event) {
	event := cloudevents.NewEvent()
	event.SetData(cloudevents.ApplicationJSON, map[string]interface{}{"actors": []string{}, "providers": []string{}, "uptime_human": "forever seconds", "uptime_seconds": int(time.Since(startTime).Seconds()), "version": "v0.0.0-wasmcloud_go", "friendly_name": "yolo-bro-1234", "labels": map[string]string{"hostcore.arch": "armv61", "hostcore.os": "raspbian", "hostcore.osfamily": "unix", "hostcore.version": "v0.0.0-wasmcloud_go", "hostcore.runtime": "wazero"}})
	event.SetID(uuid.New().String())
	event.SetSource(hostid) // "N" key
	event.SetType("com.wasmcloud.lattice.host_heartbeat")

	return ControlEvent(nsprefix), event
}
