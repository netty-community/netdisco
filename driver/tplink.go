package driver

import (
	"github.com/netty-community/netdisco/factory"
	"github.com/netty-community/netdisco/model/snmp"
)

type TPLinkDriver struct {
	factory.SnmpDiscovery
}

func NewTPLinkDriver(sc *snmp.SnmpConfig) (*TPLinkDriver, error) {
	session, err := factory.NewSnmpSession(sc)
	if err != nil {
		return nil, err
	}
	return &TPLinkDriver{
		factory.SnmpDiscovery{
			Session:   session,
			IpAddress: session.Target},
	}, nil
}
