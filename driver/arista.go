package driver

import (
	"github.com/netty-community/netdisco/factory"
	"github.com/netty-community/netdisco/model/snmp"
)

type AristaDriver struct {
	factory.SnmpDiscovery
}

func NewAristaDriver(sc *snmp.SnmpConfig) (*AristaDriver, error) {
	session, err := factory.NewSnmpSession(sc)
	if err != nil {
		return nil, err
	}
	return &AristaDriver{
		factory.SnmpDiscovery{
			Session:   session,
			IpAddress: session.Target},
	}, nil
}
