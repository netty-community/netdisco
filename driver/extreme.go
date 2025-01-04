package driver

import (
	"github.com/netty-community/netdisco/factory"
	"github.com/netty-community/netdisco/model/snmp"
)

type ExtremeDriver struct {
	factory.SnmpDiscovery
}

func NewExtremeDriver(sc *snmp.SnmpConfig) (*ExtremeDriver, error) {
	session, err := factory.NewSnmpSession(sc)
	if err != nil {
		return nil, err
	}
	return &ExtremeDriver{
		factory.SnmpDiscovery{
			Session:   session,
			IpAddress: session.Target},
	}, nil
}
