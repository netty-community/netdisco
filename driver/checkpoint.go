package driver

import (
	"github.com/netty-community/netdisco/factory"
	"github.com/netty-community/netdisco/model/snmp"
)

type CheckPointDriver struct {
	factory.SnmpDiscovery
}

func NewCheckPointDriver(sc *snmp.SnmpConfig) (*CheckPointDriver, error) {
	session, err := factory.NewSnmpSession(sc)
	if err != nil {
		return nil, err
	}
	return &CheckPointDriver{
		factory.SnmpDiscovery{
			Session:   session,
			IpAddress: session.Target},
	}, nil
}
