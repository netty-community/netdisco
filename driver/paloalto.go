package driver

import (
	"github.com/netty-community/netdisco/factory"
	"github.com/netty-community/netdisco/model/snmp"
)

type PaloAltoDriver struct {
	factory.SnmpDiscovery
}

func NewPaloAltoDriver(sc *snmp.SnmpConfig) (*PaloAltoDriver, error) {
	session, err := factory.NewSnmpSession(sc)
	if err != nil {
		return nil, err
	}
	return &PaloAltoDriver{
		factory.SnmpDiscovery{
			Session:   session,
			IpAddress: session.Target},
	}, nil
}
