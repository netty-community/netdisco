package driver

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/netty-community/netdisco/factory"
	"github.com/netty-community/netdisco/model/device"
	"github.com/netty-community/netdisco/model/snmp"
)

const vtpVlanName = ".1.3.6.1.4.1.9.9.46.1.3.1.1.4"
const vtpVlanIfIndex = ".1.3.6.1.4.1.9.9.46.1.3.1.1.18"

// https://pastebin.com/PaP5yfYU
// https://mibs.observium.org/mib/AIRESPACE-WIRELESS-MIB/
// const bsnAPDot3MacAddress = ".1.3.6.1.4.1.14179.2.2.1.1.1" // Radio mac address for clients associated with the AP

const bsnAPEthernetMacAddress = ".1.3.6.1.4.1.14179.2.2.1.1.33" // Ethernet mac address connected to switch
const bsnApIpAddress = ".1.3.6.1.4.1.14179.2.2.1.1.19"
const bsnAPName = ".1.3.6.1.4.1.14179.2.2.1.1.3"
const bsnAPModel = ".1.3.6.1.4.1.14179.2.2.1.1.16"
const bsnAPSerialNumber = ".1.3.6.1.4.1.14179.2.2.1.1.17"

// const bsnAPPrimaryMwarName = ".1.3.6.1.4.1.14179.2.2.1.1.10"
const bsnAPSoftwareVersion = ".1.3.6.1.4.1.14179.2.2.1.1.8"

// CISCO-VLAN-MEMBERSHIP-MIB && CISCO-VTP-MIB
const vmVlanType = ".1.3.6.1.4.1.9.9.68.1.2.2.1.1"
const vmVlan = ".1.3.6.1.4.1.9.9.68.1.2.2.1.2"
const vlanTrunkPortIfIndex = ".1.3.6.1.4.1.9.9.46.1.6.1.1.1"
const vlanTrunkPortDynamicStatus = ".1.3.6.1.4.1.9.9.46.1.6.1.1.14"

// vlanPortIslOperStatus

// some oid has bug if version lower than ios-xe 17.9.5
const bsnMobileStationSsid = ".1.3.6.1.4.1.14179.2.1.4.1.7"
const bsnMobileStationAPMacAddr = ".1.3.6.1.4.1.14179.2.1.4.1.4"
const bsnMobileStationIpAddress = ".1.3.6.1.4.1.14179.2.1.4.1.2"
const bsnMobileStationUserName = ".1.3.6.1.4.1.14179.2.1.4.1.3"
const bsnMobileStationRSSI = ".1.3.6.1.4.1.14179.2.1.6.1.1"
const bsnMobileStationBytesReceived = ".1.3.6.1.4.1.14179.2.1.6.1.2"
const bsnMobileStationSnr = ".1.3.6.1.4.1.14179.2.1.6.1.26"
const bsnMobileStationBytesSent = ".1.3.6.1.4.1.14179.2.1.6.1.3"
const cldcClientUpTime = ".1.3.6.1.4.1.9.9.599.1.3.1.1.15"

const cldcClientDeviceType = ".1.3.6.1.4.1.9.9.599.1.3.1.1.44"
const bsnMobileStationVlanId = ".1.3.6.1.4.1.14179.2.1.4.1.29"
const cldcClientChannel = ".1.3.6.1.4.1.9.9.599.1.3.1.1.35"

// const bsnMobileStationProtocol = ".1.3.6.1.4.1.14179.2.1.4.1.25"

type CiscoBaseDriver struct {
	factory.SnmpDiscovery
}

type CiscoIosDriver struct {
	CiscoBaseDriver
}

type CiscoIosXRDriver struct {
	CiscoBaseDriver
}

type CiscoNexusOSDriver struct {
	CiscoBaseDriver
}

type CiscoIosXEDriver struct {
	CiscoBaseDriver
}

func NewCiscoIosDriver(sc *snmp.SnmpConfig) (*CiscoIosDriver, error) {
	session, err := factory.NewSnmpSession(sc)
	if err != nil {
		return nil, err
	}
	return &CiscoIosDriver{
		CiscoBaseDriver{
			factory.SnmpDiscovery{
				Session:   session,
				IpAddress: session.Target},
		},
	}, nil
}

func NewCiscoIosXRDriver(sc *snmp.SnmpConfig) (*CiscoIosXRDriver, error) {
	session, err := factory.NewSnmpSession(sc)
	if err != nil {
		return nil, err
	}
	return &CiscoIosXRDriver{
		CiscoBaseDriver{
			factory.SnmpDiscovery{
				Session:   session,
				IpAddress: session.Target},
		},
	}, nil
}

func NewCiscoNexusOSDriver(sc *snmp.SnmpConfig) (*CiscoNexusOSDriver, error) {
	session, err := factory.NewSnmpSession(sc)
	if err != nil {
		return nil, err
	}
	return &CiscoNexusOSDriver{
		CiscoBaseDriver{
			factory.SnmpDiscovery{
				Session:   session,
				IpAddress: session.Target},
		},
	}, nil
}

func NewCiscoIosXEDriver(sc *snmp.SnmpConfig) (*CiscoIosXEDriver, error) {
	session, err := factory.NewSnmpSession(sc)
	if err != nil {
		return nil, err
	}
	return &CiscoIosXEDriver{
		CiscoBaseDriver{
			factory.SnmpDiscovery{
				Session:   session,
				IpAddress: session.Target},
		},
	}, nil
}

func (cd *CiscoBaseDriver) Vlans() (vlan []*device.VlanItem, errors []string) {
	l2Vlan, err := cd.Session.BulkWalkAll(vtpVlanName)
	if err != nil {
		return nil, []string{fmt.Sprintf("failed to get vlan from %s", cd.IpAddress)}
	}
	l2VlanIfIndex, errIfIndex := cd.Session.BulkWalkAll(vtpVlanIfIndex)
	if errIfIndex != nil {
		errors = append(errors, errIfIndex.Error())
	}
	indexL2Vlan := factory.ExtractString(vtpVlanName, l2Vlan)
	indexVlanIndex := factory.ExtractInteger(vtpVlanIfIndex, l2VlanIfIndex)
	for i, v := range indexL2Vlan {
		vlanIdStrings := strings.Split(i, ".")
		vlanIdString := vlanIdStrings[len(vlanIdStrings)-1]
		vlanId, _ := strconv.Atoi(vlanIdString)
		_vlan := &device.VlanItem{
			VlanId:   uint32(vlanId),
			VlanName: v,
			IfIndex:  indexVlanIndex[i],
		}
		vlan = append(vlan, _vlan)
	}

	return vlan, errors
}

func (cd *CiscoBaseDriver) VlanAssign() (vlan []*factory.VlanAssignItem, errors []string) {
	vlanType, errVlanType := cd.Session.BulkWalkAll(vmVlanType)
	if errVlanType != nil {
		return nil, []string{fmt.Sprintf("failed to get vlan assignment from %s", cd.IpAddress)}
	}

	vlanId, errVlanId := cd.Session.BulkWalkAll(vmVlan)
	if errVlanId != nil {
		return nil, []string{fmt.Sprintf("failed to get vlan assignment from %s", cd.IpAddress)}
	}
	indexVlanType := factory.ExtractInteger(vmVlanType, vlanType)
	indexVlanId := factory.ExtractInteger(vmVlan, vlanId)
	for i, v := range indexVlanType {
		ifIndexString := strings.TrimPrefix(i, ".")
		ifIndex, _ := strconv.Atoi(ifIndexString)
		_vlan := &factory.VlanAssignItem{
			VlanType: factory.GetCiscoVlanMemberShipTypeValue(v),
			IfIndex:  uint64(ifIndex),
			VlanId:   uint32(indexVlanId[i]),
		}
		vlan = append(vlan, _vlan)
	}
	return vlan, errors

}

func (cd *CiscoBaseDriver) APs() (ap []*device.Ap, errors []string) {
	apIP, errApIP := cd.Session.BulkWalkAll(bsnApIpAddress)
	if len(apIP) == 0 || errApIP != nil {
		return nil, []string{fmt.Sprintf("failed to get ap ipAddress from %s", cd.IpAddress)}
	}
	apMac, errApMac := cd.Session.BulkWalkAll(bsnAPEthernetMacAddress)
	apName, errApName := cd.Session.BulkWalkAll(bsnAPName)
	apType, errApType := cd.Session.BulkWalkAll(bsnAPModel)
	apSerialNumber, errApSerialNumber := cd.Session.BulkWalkAll(bsnAPSerialNumber)
	// apPrimaryMwarName, errApPrimaryMwarName := cd.Session.BulkWalkAll(bsnAPPrimaryMwarName)
	apVersion, errApVersion := cd.Session.BulkWalkAll(bsnAPSoftwareVersion)
	if errApMac != nil || errApName != nil || errApType != nil || errApSerialNumber != nil || errApVersion != nil {
		errors = append(errors, errApMac.Error())
		errors = append(errors, errApName.Error())
		errors = append(errors, errApType.Error())
		errors = append(errors, errApSerialNumber.Error())
		errors = append(errors, errApVersion.Error())
	}
	indexApIP := factory.ExtractString(bsnApIpAddress, apIP)
	indexApMac := factory.ExtractMacAddress(bsnAPEthernetMacAddress, apMac)
	indexApName := factory.ExtractString(bsnAPName, apName)
	indexApType := factory.ExtractString(bsnAPModel, apType)
	indexApSerialNumber := factory.ExtractString(bsnAPSerialNumber, apSerialNumber)
	// indexApPrimaryMwarName := factory.ExtractString(bsnAPPrimaryMwarName, apPrimaryMwarName)
	indexApVersion := factory.ExtractString(bsnAPSoftwareVersion, apVersion)
	for i, v := range indexApIP {
		ipAddr := cd.IpAddress
		apVer := indexApVersion[i]
		ap = append(ap, &device.Ap{
			Name:            indexApName[i],
			ManagementIp:    v,
			MacAddress:      indexApMac[i],
			DeviceModel:     indexApType[i],
			SerialNumber:    indexApSerialNumber[i],
			WlanACIpAddress: &ipAddr,
			OsVersion:       &apVer,
		})
	}
	return ap, errors

}
