package driver

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/netty-community/netdisco/factory"
	"github.com/netty-community/netdisco/model/device"
	"github.com/netty-community/netdisco/model/snmp"
)

const hwL2VlanDescr = ".1.3.6.1.4.1.2011.5.25.42.3.1.1.1.1.2"
const hwL2VlanIfIndex = ".1.3.6.1.4.1.2011.5.25.42.3.1.1.1.1.6"

const hwWlanApMac = ".1.3.6.1.4.1.2011.6.139.13.3.3.1.1"
const hwWlanApIpAddress = ".1.3.6.1.4.1.2011.6.139.13.3.3.1.13"
const hwWlanApName = ".1.3.6.1.4.1.2011.6.139.13.3.3.1.4"
const hwWlanAPGroup = ".1.3.6.1.4.1.2011.6.139.13.3.3.1.5"
const hwWlanApTypeInfo = ".1.3.6.1.4.1.2011.6.139.13.3.3.1.3"
const hwWlanApSn = ".1.3.6.1.4.1.2011.6.139.13.3.3.1.2"
const hwWlanApSoftwareVersion = ".1.3.6.1.4.1.2011.6.139.13.3.3.1.7"

const hwStackRun = ".1.3.6.1.4.1.2011.5.25.183.1.1" // 0-8
const hwMemberCurrentStackId = ".1.3.6.1.4.1.2011.5.25.183.1.20.1.1"
const hwMemberStackPriority = ".1.3.6.1.4.1.2011.5.25.183.1.20.1.2" // 1-255, default 100
const hwMemberStackRole = ".1.3.6.1.4.1.2011.5.25.183.1.20.1.3"
const hwMemberStackMacAddress = ".1.3.6.1.4.1.2011.5.25.183.1.20.1.4"
const hwStackPortName = ".1.3.6.1.4.1.2011.5.25.183.1.21.1.3"
const hwStackPortStatus = ".1.3.6.1.4.1.2011.5.25.183.1.21.1.5"

// HUAWEI-L2IF-MIB
const hwL2IfPortIfIndex = ".1.3.6.1.4.1.2011.5.25.42.1.1.1.3.1.2"
const hwL2IfPortType = ".1.3.6.1.4.1.2011.5.25.42.1.1.1.3.1.3"
const hwL2IfPVID = ".1.3.6.1.4.1.2011.5.25.42.1.1.1.3.1.4"

const hwWlanStaMac = ".1.3.6.1.4.1.2011.6.139.18.1.2.1.1"
const hwWlanStaUsername = ".1.3.6.1.4.1.2011.6.139.18.1.2.1.2"
const hwWlanStaApMac = ".1.3.6.1.4.1.2011.6.139.18.1.2.1.3"
const hwWlanStaIP = ".1.3.6.1.4.1.2011.6.139.18.1.2.1.25"
const hwWlanStaApName = ".1.3.6.1.4.1.2011.6.139.18.1.2.1.4"
const hwWlanStaAssocBand = ".1.3.6.1.4.1.2011.6.139.18.1.2.1.7"
const hwWlanStaAccessChannel = ".1.3.6.1.4.1.2011.6.139.18.1.2.1.9"
const hwWlanStaEssName = ".1.3.6.1.4.1.2011.6.139.18.1.2.1.16"
const hwWlanStaVlan = ".1.3.6.1.4.1.2011.6.139.18.1.2.1.24"
const hwWlanStaRssi = ".1.3.6.1.4.1.2011.6.139.18.1.2.1.42"
const hwWlanStaSnrUs = ".1.3.6.1.4.1.2011.6.139.18.1.2.1.44"
const hwWlanStaWirelessTxBytes = ".1.3.6.1.4.1.2011.6.139.18.1.2.1.37"
const hwWlanStaWirelessRxBytes = ".1.3.6.1.4.1.2011.6.139.18.1.2.1.34"
const hwWlanStaOnlineTime = ".1.3.6.1.4.1.2011.6.139.18.1.2.1.30"

type HuaweiDriver struct {
	factory.SnmpDiscovery
}

func NewHuaweiDriver(sc *snmp.SnmpConfig) (*HuaweiDriver, error) {
	session, err := factory.NewSnmpSession(sc)
	if err != nil {
		return nil, err
	}
	return &HuaweiDriver{
		factory.SnmpDiscovery{
			Session:   session,
			IpAddress: session.Target},
	}, nil
}

func (hd *HuaweiDriver) Vlans() (vlan []*device.VlanItem, errors []string) {
	l2Vlan, err := hd.Session.BulkWalkAll(hwL2VlanDescr)
	l2VlanIfIndex, errIfIndex := hd.Session.BulkWalkAll(hwL2VlanIfIndex)
	if err != nil || errIfIndex != nil {
		errors = append(errors, err.Error())
		errors = append(errors, errIfIndex.Error())
	}
	indexL2Vlan := factory.ExtractString(hwL2VlanDescr, l2Vlan)
	indexVlanIndex := factory.ExtractInteger(hwL2VlanIfIndex, l2VlanIfIndex)

	for i, v := range indexL2Vlan {
		vlanIdString := strings.TrimPrefix(i, ".")
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

func (hd *HuaweiDriver) APs() (ap []*device.Ap, errors []string) {
	apIp, errApIp := hd.Session.BulkWalkAll(hwWlanApIpAddress)
	if len(apIp) == 0 || errApIp != nil {
		return nil, []string{fmt.Sprintf("failed to get ap ipAddress from %s", hd.IpAddress)}
	}
	apName, errApName := hd.Session.BulkWalkAll(hwWlanApName)
	apGroupName, errApGroupName := hd.Session.BulkWalkAll(hwWlanAPGroup)
	apModel, errApModel := hd.Session.BulkWalkAll(hwWlanApTypeInfo)
	apSerialNumber, errApSerialNumber := hd.Session.BulkWalkAll(hwWlanApSn)
	apMac, errApMac := hd.Session.BulkWalkAll(hwWlanApMac)
	apVersion, errApVersion := hd.Session.BulkWalkAll(hwWlanApSoftwareVersion)
	if errApName != nil || errApGroupName != nil || errApModel != nil || errApSerialNumber != nil || errApMac != nil || errApVersion != nil {
		errors = append(errors, errApVersion.Error())
		errors = append(errors, errApName.Error())
		errors = append(errors, errApGroupName.Error())
		errors = append(errors, errApModel.Error())
		errors = append(errors, errApSerialNumber.Error())
		errors = append(errors, errApMac.Error())
	}
	indexApIP := factory.ExtractString(hwWlanApIpAddress, apIp)
	indexApMac := factory.ExtractString(hwWlanApMac, apMac)
	indexApName := factory.ExtractString(hwWlanApName, apName)
	indexApGroupName := factory.ExtractString(hwWlanAPGroup, apGroupName)
	indexApModel := factory.ExtractString(hwWlanApTypeInfo, apModel)
	indexApSerialNumber := factory.ExtractString(hwWlanApSn, apSerialNumber)
	indexApVersion := factory.ExtractString(hwWlanApSoftwareVersion, apVersion)
	for i, v := range indexApIP {
		group := indexApGroupName[i]
		acIp := hd.IpAddress
		apVer := indexApVersion[i]
		ap = append(ap, &device.Ap{
			Name:            indexApName[i],
			ManagementIp:    v,
			MacAddress:      indexApMac[i],
			DeviceModel:     indexApModel[i],
			SerialNumber:    indexApSerialNumber[i],
			GroupName:       &group,
			WlanACIpAddress: &acIp,
			OsVersion:       &apVer,
		})
	}
	return ap, errors
}

func (hd *HuaweiDriver) Discovery() *factory.DiscoveryResponse {
	sysDescr, sysError := hd.SysDescr()
	sysUpTime, sysUpTimeError := hd.SysUpTime()
	sysName, sysNameError := hd.SysName()
	chassisId, chassisIdError := hd.ChassisId()
	interfaces, interfacesError := hd.Interfaces()
	entities, entitiesError := hd.Entities()
	lldp, lldpError := hd.LldpNeighbors()
	macAddress, macAddressError := hd.MacAddressTable()
	arp, arpError := hd.ArpTable()
	arp = factory.EnrichArpInfo(arp, interfaces)
	vlan, VlanError := hd.Vlans()
	vlan = factory.EnrichVlanInfo(vlan, interfaces)
	macAddress_ := factory.EnrichMacAddress(macAddress, interfaces, lldp, arp)

	response := &factory.DiscoveryResponse{
		SysDescr:        sysDescr,
		Uptime:          sysUpTime,
		Hostname:        sysName,
		ChassisId:       chassisId,
		Interfaces:      interfaces,
		LldpNeighbors:   lldp,
		Entities:        entities,
		MacAddressTable: macAddress_,
		ArpTable:        arp,
		Vlans:           vlan,
	}
	if sysError != nil {
		response.Errors = append(response.Errors, sysError.Error())
	}
	if sysUpTimeError != nil {
		response.Errors = append(response.Errors, sysUpTimeError.Error())
	}
	if sysNameError != nil {
		response.Errors = append(response.Errors, sysNameError.Error())
	}
	if chassisIdError != nil {
		response.Errors = append(response.Errors, chassisIdError.Error())
	}
	if interfacesError != nil {
		response.Errors = append(response.Errors, interfacesError...)
	}
	if entitiesError != nil {
		response.Errors = append(response.Errors, entitiesError...)
	}
	if lldpError != nil {
		response.Errors = append(response.Errors, lldpError...)
	}
	if macAddressError != nil {
		response.Errors = append(response.Errors, macAddressError...)
	}
	if arpError != nil {
		response.Errors = append(response.Errors, arpError...)
	}
	if VlanError != nil {
		response.Errors = append(response.Errors, VlanError...)
	}
	return response
}
