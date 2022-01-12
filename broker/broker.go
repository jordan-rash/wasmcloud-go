package broker

import "fmt"

type Commands struct{}
type Queries struct{}

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
