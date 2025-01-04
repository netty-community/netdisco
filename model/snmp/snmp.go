package snmp

import (
	"github.com/gosnmp/gosnmp"
	platform "github.com/netty-community/netdisco/model/platform"
)

type SnmpV3Params struct {
	ContextEngineId *string
	SecurityName    *string
	SecurityLevel   gosnmp.SnmpV3SecurityModel
	AuthProtocol    gosnmp.SnmpV3AuthProtocol
	AuthPassword    *string
	PrivProtocol    gosnmp.SnmpV3PrivProtocol
	PrivPassword    *string
}

type SnmpConfig struct {
	IpAddress      string
	Port           uint16
	Version        gosnmp.SnmpVersion
	Timeout        uint8
	Community      *string
	V3Params       *SnmpV3Params
	MaxRepetitions int
	Platform       *platform.Platform
	SysObjectId    string
}

func (c *SnmpConfig) Validate() bool {
	switch c.Version {
	case gosnmp.Version2c:
		return c.Community != nil
	case gosnmp.Version3:
		return c.V3Params != nil
	default:
		return false
	}
}
