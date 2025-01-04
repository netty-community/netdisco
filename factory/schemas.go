package factory

import (
	"github.com/netty-community/netdisco/devicemodel"
	"github.com/netty-community/netdisco/model/device"
)

type VlanAssignItem struct {
	VlanType string `json:"vlanType"`
	VlanId   uint32 `json:"vlanId"`
	IfIndex  uint64 `json:"ifIndex"`
}

type Route struct{}

type Prefix struct{}

type DiscoveryResponse struct {
	Hostname        string                    `json:"hostname"`
	SysDescr        string                    `json:"sysDescr"`
	Uptime          uint64                    `json:"uptime"`
	ChassisId       string                    `json:"chassisID"`
	Interfaces      []*device.DeviceInterface `json:"interfaces"`
	LldpNeighbors   []*device.LldpNeighbor    `json:"lldpNeighbors"`
	Entities        []*device.Entity          `json:"entities"`
	Stacks          []*device.Stack           `json:"stacks"`
	Vlans           []*device.VlanItem        `json:"vlans"`
	MacAddressTable []*device.MacAddressItem  `json:"macAddressTable"`
	ArpTable        []*device.ArpItem         `json:"arpTable"`
	Errors          []string                  `json:"errors"`
}

type DiscoveryBasicResponse struct {
	Hostname    string   `json:"hostname"`
	SysDescr    string   `json:"sysDescr"`
	ChassisId   string   `json:"chassisID"`
	Errors      []string `json:"errors"`
}

type DispatchResponse struct {
	IpAddress     string                   `json:"ipAddress"`
	Data          *DiscoveryResponse       `json:"data"`
	SnmpReachable bool                     `json:"snmpReachable"`
	IcmpReachable bool                     `json:"icmpReachable"`
	SshReachable  bool                     `json:"sshReachable"`
	SysObjectId   string                   `json:"sysObjectId"`
	DeviceModel   *devicemodel.DeviceModel `json:"deviceModel"`
}

type DispatchBasicResponse struct {
	IpAddress     string                   `json:"ipAddress"`
	Data          *DiscoveryBasicResponse  `json:"data"`
	SnmpReachable bool                     `json:"snmpReachable"`
	IcmpReachable bool                     `json:"icmpReachable"`
	SshReachable  bool                     `json:"sshReachable"`
	SysObjectId   string                   `json:"sysObjectId"`
	DeviceModel   *devicemodel.DeviceModel `json:"deviceModel"`
}

type DispatchApScanResponse struct {
	IpAddress     string                   `json:"ipAddress"`
	Data          []*device.Ap             `json:"data"`
	SnmpReachable bool                     `json:"snmpReachable"`
	IcmpReachable bool                     `json:"icmpReachable"`
	SshReachable  bool                     `json:"sshReachable"`
	SysObjectId   string                   `json:"sysObjectId"`
	DeviceModel   *devicemodel.DeviceModel `json:"deviceModel"`
	Errors        []string                 `json:"errors"`
}

type VlanIpRange struct {
	VlanId  uint32
	Range   string
	Gateway string
}
