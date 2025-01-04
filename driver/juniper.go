package driver

import (
	"fmt"

	"github.com/netty-community/netdisco/factory"
	"github.com/netty-community/netdisco/model/device"
	"github.com/netty-community/netdisco/model/snmp"
)

// AP Name
const jnxWAPStatusAccessPoint = ".1.3.6.1.4.1.2636.3.88.1.1.1.1.2"

// AP Type
const jnxWAPStatusType = ".1.3.6.1.4.1.2636.3.88.1.1.1.1.3"
const jnxWAPStatusSerialNumber = ".1.3.6.1.4.1.2636.3.88.1.1.1.1.5"
const jnxWAPStatusFirmwareVersion = ".1.3.6.1.4.1.2636.3.88.1.1.1.1.6"
const jnxWAPStatusEthernetPortMAC = ".1.3.6.1.4.1.2636.3.88.1.1.1.1.12"
const jnxWAPStatusEthernetIPv4 = ".1.3.6.1.4.1.2636.3.88.1.1.1.1.13"

// VLAN MIB: https://mibs.observium.org/mib/JUNIPER-VLAN-MIB/
const jnxExVlanID = ".1.3.6.1.4.1.2636.3.40.1.5.1.5.1.1"
const jnxExVlanName = ".1.3.6.1.4.1.2636.3.40.1.5.1.5.1.2"
const jnxExVlanSnmpIfIndex = ".1.3.6.1.4.1.2636.3.40.1.5.1.6.1.8"

type JuniperDriver struct {
	factory.SnmpDiscovery
}

func NewJuniperDriver(sc *snmp.SnmpConfig) (*JuniperDriver, error) {
	session, err := factory.NewSnmpSession(sc)
	if err != nil {
		return nil, err
	}
	return &JuniperDriver{
		factory.SnmpDiscovery{
			Session:   session,
			IpAddress: session.Target},
	}, nil
}

func (d *JuniperDriver) APs() (ap []*device.Ap, errors []string) {
	apIP, errApIP := d.Session.BulkWalkAll(jnxWAPStatusEthernetIPv4)
	if len(apIP) == 0 || errApIP != nil {
		return nil, []string{fmt.Sprintf("failed to get ap ipAddress from %s", d.IpAddress)}
	}
	apMac, errApMac := d.Session.BulkWalkAll(jnxWAPStatusEthernetPortMAC)
	apName, errApName := d.Session.BulkWalkAll(jnxWAPStatusAccessPoint)
	apType, errApType := d.Session.BulkWalkAll(jnxWAPStatusType)
	apSerialNumber, errApSerialNumber := d.Session.BulkWalkAll(jnxWAPStatusSerialNumber)
	apVersion, errApVersion := d.Session.BulkWalkAll(jnxWAPStatusFirmwareVersion)
	if errApMac != nil || errApName != nil || errApType != nil || errApSerialNumber != nil || errApVersion != nil {
		errors = append(errors, errApMac.Error())
		errors = append(errors, errApName.Error())
		errors = append(errors, errApType.Error())
		errors = append(errors, errApSerialNumber.Error())
		errors = append(errors, errApVersion.Error())
	}
	indexApIP := factory.ExtractString(jnxWAPStatusEthernetIPv4, apIP)
	indexApMac := factory.ExtractMacAddress(jnxWAPStatusEthernetPortMAC, apMac)
	indexApName := factory.ExtractString(jnxWAPStatusAccessPoint, apName)
	indexApType := factory.ExtractString(jnxWAPStatusType, apType)
	indexApSerialNumber := factory.ExtractString(jnxWAPStatusSerialNumber, apSerialNumber)
	indexApVersion := factory.ExtractString(jnxWAPStatusFirmwareVersion, apVersion)
	for i, v := range indexApIP {
		apVer := indexApVersion[i]
		ap = append(ap, &device.Ap{
			Name:         indexApName[i],
			ManagementIp: v,
			MacAddress:   indexApMac[i],
			DeviceModel:  indexApType[i],
			SerialNumber: indexApSerialNumber[i],
			OsVersion:    &apVer,
		})
	}
	return ap, errors
}

// func (d *JuniperDriver) Vlans() (vlan []*device.VlanItem, errors []string) {
// 	l2Vlan, err := d.Session.BulkWalkAll(jnxExVlanID)
// }
