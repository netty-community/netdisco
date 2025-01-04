package driver

import (
	"fmt"

	"github.com/netty-community/netdisco/factory"
	"github.com/netty-community/netdisco/model/device"
	"github.com/netty-community/netdisco/model/snmp"
)

// WLAN AP-MIB: https://www.h3c.com/cn/d_202401/2021242_30005_0.htm#_Toc146378887
// need more test for devices
// HH3C-DOT11-APMT-MIB:

const hh3cDot11CurrAPIPAddress string = ".1.3.6.1.4.1.25506.2.75.2.1.2.1.2"
const hh3cDot11CurrAPMacAddress string = ".1.3.6.1.4.1.25506.2.75.2.1.2.1.3"
const hh3cDot11CurrAPName string = ".1.3.6.1.4.1.25506.2.75.2.1.2.1.8"
const hh3cDot11CurrAPModelName string = ".1.3.6.1.4.1.25506.2.75.2.1.2.1.9"
const hh3cDot11CurrAPTemplateName string = ".1.3.6.1.4.1.25506.2.75.2.1.2.1.6"
const hh3cDot11CurrAPSoftwareVersion string = ".1.3.6.1.4.1.25506.2.75.2.1.2.1.11"

const hh3cDot11StationMAC = ".1.3.6.1.4.1.25506.2.75.3.1.1.1.1"
const hh3cDot11StationIPAddress = ".1.3.6.1.4.1.25506.2.75.3.1.1.1.2"
const hh3cDot11StationUserName = ".1.3.6.1.4.1.25506.2.75.3.1.1.1.3"
const hh3cDot11StationSSIDName = ".1.3.6.1.4.1.25506.2.75.3.1.1.1.12"
const hh3cDot11StationVlanId = ".1.3.6.1.4.1.25506.2.75.3.1.1.1.11"
const hh3cDot11StationRxSNR = ".1.3.6.1.4.1.25506.2.75.3.1.1.1.17"
const hh3cDot11StationAssociateAPMACAddressCM = ".1.3.6.1.4.1.25506.2.75.3.1.1.1.40"
const hh3cDot11StationRSSI = ".1.3.6.1.4.1.25506.2.75.3.1.1.1.7"
const hh3cDot11StationChannel = ".1.3.6.1.4.1.25506.2.75.3.1.1.1.8"
const hh3cDot11StationTxSpeed = ".1.3.6.1.4.1.25506.2.75.3.1.1.1.24"
const hh3cDot11StationRxSpeed = ".1.3.6.1.4.1.25506.2.75.3.1.1.1.25"
const hh3cDot11StationMaxRate = ". 1.3.6.1.4.1.25506.2.75.3.1.1.1.33" // (Mbps)
const hh3cDot11StationAssTime = ".1.3.6.1.4.1.25506.2.75.3.1.1.1.30"
const hh3cDot11StationVendorName = ".1.3.6.1.4.1.25506.2.75.3.1.1.1.20"

type H3CDriver struct {
	factory.SnmpDiscovery
}

func NewH3CDriver(sc *snmp.SnmpConfig) (*H3CDriver, error) {
	session, err := factory.NewSnmpSession(sc)
	if err != nil {
		return nil, err
	}
	return &H3CDriver{
		factory.SnmpDiscovery{
			Session:   session,
			IpAddress: session.Target},
	}, nil
}

func (hd *H3CDriver) APs() (ap []*device.Ap, errors []string) {
	apIp, errApIp := hd.Session.BulkWalkAll(hh3cDot11CurrAPIPAddress)
	if len(apIp) == 0 || errApIp != nil {
		return nil, []string{fmt.Sprintf("failed to get ap ipAddress from %s", hd.IpAddress)}
	}
	apMac, errApMac := hd.Session.BulkWalkAll(hh3cDot11CurrAPMacAddress)
	apName, errApName := hd.Session.BulkWalkAll(hh3cDot11CurrAPName)
	apModel, errApModel := hd.Session.BulkWalkAll(hh3cDot11CurrAPModelName)
	apGroup, errApGroup := hd.Session.BulkWalkAll(hh3cDot11CurrAPTemplateName)
	apVersion, errApVersion := hd.Session.BulkWalkAll(hh3cDot11CurrAPSoftwareVersion)
	if errApMac != nil || errApName != nil || errApModel != nil || errApGroup != nil || errApVersion != nil {
		errors = append(errors, errApMac.Error())
		errors = append(errors, errApName.Error())
		errors = append(errors, errApModel.Error())
		errors = append(errors, errApGroup.Error())
	}
	indexApIP := factory.ExtractString(hh3cDot11CurrAPIPAddress, apIp)
	indexApMac := factory.ExtractMacAddress(hh3cDot11CurrAPMacAddress, apMac)
	indexApName := factory.ExtractString(hh3cDot11CurrAPName, apName)
	indexApModel := factory.ExtractString(hh3cDot11CurrAPModelName, apModel)
	indexApGroup := factory.ExtractString(hh3cDot11CurrAPTemplateName, apGroup)
	indexApVersion := factory.ExtractString(hh3cDot11CurrAPSoftwareVersion, apVersion)
	for i, v := range indexApIP {
		group_name := indexApGroup[i]
		acIp := hd.IpAddress
		apVer := indexApVersion[i]

		ap = append(ap, &device.Ap{
			Name:            indexApName[i],
			ManagementIp:    v,
			MacAddress:      indexApMac[i],
			DeviceModel:     indexApModel[i],
			GroupName:       &group_name,
			WlanACIpAddress: &acIp,
			OsVersion:       &apVer,
		})
	}
	return ap, errors
}
