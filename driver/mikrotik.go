package driver

import (
	"github.com/netty-community/netdisco/factory"
	"github.com/netty-community/netdisco/model/snmp"
)

type MikroTikDriver struct {
	factory.SnmpDiscovery
}

func NewMikroTikDriver(sc *snmp.SnmpConfig) (*MikroTikDriver, error) {
	session, err := factory.NewSnmpSession(sc)
	if err != nil {
		return nil, err
	}
	return &MikroTikDriver{
		factory.SnmpDiscovery{
			Session:   session,
			IpAddress: session.Target},
	}, nil
}
