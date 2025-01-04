package driver

import (
	"fmt"

	"github.com/netty-community/netdisco/factory"
	"github.com/netty-community/netdisco/model/device"
	"github.com/netty-community/netdisco/model/snmp"
)

// const wlanApMacAddress string = ".1.3.6.1.4.1.14823.2.2.1.5.2.1.4.1.1" not implement by aruba, replace with snmpIndex

const wlanAPIpAddress string = ".1.3.6.1.4.1.14823.2.2.1.5.2.1.4.1.2"
const wlanAPName string = ".1.3.6.1.4.1.14823.2.2.1.5.2.1.4.1.3"
const wlanAPGroupName string = ".1.3.6.1.4.1.14823.2.2.1.5.2.1.4.1.4"
const wlanAPModel string = ".1.3.6.1.4.1.14823.2.2.1.5.2.1.4.1.5"
const wlanAPSerialNumber string = ".1.3.6.1.4.1.14823.2.2.1.5.2.1.4.1.6"
const wlanAPSwitchIpAddress string = ".1.3.6.1.4.1.14823.2.2.1.5.2.1.4.1.39"
const wlanAPSwVersion string = ".1.3.6.1.4.1.14823.2.2.1.5.2.1.4.1.34"
const wlanAPBssidAPMacAddress = ".1.3.6.1.4.1.14823.2.2.1.5.2.1.7.1.13"

// WLSX-SYSTEMEXT-MIB
const wlsxSysExtHostname string = ".1.3.6.1.4.1.14823.2.2.1.2.1.2.0"
const wlsxSysExtModelName string = ".1.3.6.1.4.1.14823.2.2.1.2.1.3.0"
const wlsxSysExtSwVersion string = ".1.3.6.1.4.1.14823.2.2.1.2.1.28.0"
const wlsxSysExtSerialNumber string = ".1.3.6.1.4.1.14823.2.2.1.2.1.29.0"

const nUserName = ".1.3.6.1.4.1.14823.2.2.1.4.1.2.1.3"
const nUserAssignedVlan = ".1.3.6.1.4.1.14823.2.2.1.4.1.2.1.17"
const nUserApBSSID = ".1.3.6.1.4.1.14823.2.2.1.4.1.2.1.11"
const wlanStaAccessPointESSID = ".1.3.6.1.4.1.14823.2.2.1.5.2.2.1.1.12"
const wlanStaRSSI = ".1.3.6.1.4.1.14823.2.2.1.5.2.2.1.1.14"
const wlanStaUpTime = ".1.3.6.1.4.1.14823.2.2.1.5.2.2.1.1.15"
const wlanStaChannel = ".1.3.6.1.4.1.14823.2.2.1.5.2.2.1.1.6"
const wlanStaTxBytes = ".1.3.6.1.4.1.14823.2.2.1.5.3.2.1.1.3"
const wlanStaRxBytes = ".1.3.6.1.4.1.14823.2.2.1.5.3.2.1.1.5"
const nUserDeviceType = ".1.3.6.1.4.1.14823.2.2.1.4.1.2.1.39"

type ArubaDriver struct {
	factory.SnmpDiscovery
}

type ArubaOSSwitchDriver struct {
	factory.SnmpDiscovery
}

func NewArubaDriver(sc *snmp.SnmpConfig) (*ArubaDriver, error) {
	session, err := factory.NewSnmpSession(sc)
	if err != nil {
		return nil, err
	}
	return &ArubaDriver{
		factory.SnmpDiscovery{
			Session:   session,
			IpAddress: session.Target},
	}, nil
}

func NewArubaOSSwitchDriver(sc *snmp.SnmpConfig) (*ArubaOSSwitchDriver, error) {
	session, err := factory.NewSnmpSession(sc)
	if err != nil {
		return nil, err
	}
	return &ArubaOSSwitchDriver{
		factory.SnmpDiscovery{
			Session:   session,
			IpAddress: session.Target},
	}, nil
}

func (ad *ArubaDriver) Entities() (entities []*device.Entity, errors []string) {
	hostname, errHostname := ad.Session.Get([]string{wlsxSysExtHostname})
	if errHostname != nil {
		errors = append(errors, errHostname.Error())
	}
	modelName, errModelName := ad.Session.Get([]string{wlsxSysExtModelName})
	swVersion, errSwVersion := ad.Session.Get([]string{wlsxSysExtSwVersion})
	serialNumber, errSerialNumber := ad.Session.Get([]string{wlsxSysExtSerialNumber})
	if errModelName != nil || errSwVersion != nil || errSerialNumber != nil {
		errors = append(errors, errModelName.Error())
		errors = append(errors, errSwVersion.Error())
		errors = append(errors, errSerialNumber.Error())
		return nil, errors
	}
	return []*device.Entity{
		{
			EntityPhysicalClass:       "chassis",
			EntityPhysicalName:        fmt.Sprintf("%s", hostname.Variables[0].Value),
			EntityPhysicalSoftwareRev: fmt.Sprintf("%s", swVersion.Variables[0].Value),
			EntityPhysicalSerialNum:   fmt.Sprintf("%s", serialNumber.Variables[0].Value),
			EntityPhysicalDescr:       fmt.Sprintf("%s", modelName.Variables[0].Value),
		},
	}, nil

}

func (ad *ArubaDriver) APs() (ap []*device.Ap, errors []string) {
	apIp, errApIP := ad.Session.BulkWalkAll(wlanAPIpAddress)
	if len(apIp) == 0 || errApIP != nil {
		return nil, []string{fmt.Sprintf("failed to get ap ipAddress from %s", ad.IpAddress)}
	}
	apName, errApName := ad.Session.BulkWalkAll(wlanAPName)
	apGroupName, errApGroupName := ad.Session.BulkWalkAll(wlanAPGroupName)
	apModel, errApModel := ad.Session.BulkWalkAll(wlanAPModel)
	apSerialNumber, errApSerialNumber := ad.Session.BulkWalkAll(wlanAPSerialNumber)
	switchIp, errSwitchIp := ad.Session.BulkWalkAll(wlanAPSwitchIpAddress)
	apVersion, errApVersion := ad.Session.BulkWalkAll(wlanAPSwVersion)

	if errApName != nil || errApGroupName != nil || errApModel != nil || errApSerialNumber != nil || errSwitchIp != nil || errApVersion != nil {
		errors = append(errors, errApName.Error())
		errors = append(errors, errApGroupName.Error())
		errors = append(errors, errApModel.Error())
		errors = append(errors, errApSerialNumber.Error())
		errors = append(errors, errSwitchIp.Error())
		errors = append(errors, errApVersion.Error())
	}
	indexApIp := factory.ExtractString(wlanAPIpAddress, apIp)
	indexApName := factory.ExtractString(wlanAPName, apName)
	indexApGroupName := factory.ExtractString(wlanAPGroupName, apGroupName)
	indexApModel := factory.ExtractString(wlanAPModel, apModel)
	indexApSerialNumber := factory.ExtractString(wlanAPSerialNumber, apSerialNumber)
	indexSwitchIp := factory.ExtractString(wlanAPSwitchIpAddress, switchIp)
	indexApVersion := factory.ExtractString(wlanAPSwVersion, apVersion)
	for i, v := range indexApIp {
		apMac := factory.StringToHexMac(i)
		apName := indexApName[i]
		groupName := indexApGroupName[i]
		switchIp := indexSwitchIp[i]
		apVersion := indexApVersion[i]
		ap = append(ap, &device.Ap{
			Name:            apName,
			ManagementIp:    v,
			MacAddress:      apMac,
			GroupName:       &groupName,
			DeviceModel:     indexApModel[i],
			SerialNumber:    indexApSerialNumber[i],
			WlanACIpAddress: &switchIp,
			OsVersion:       &apVersion,
		})
	}
	return ap, errors
}
