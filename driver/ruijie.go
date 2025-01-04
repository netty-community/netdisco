package driver

import (
	"github.com/netty-community/netdisco/factory"
	"github.com/netty-community/netdisco/model/device"
	"github.com/netty-community/netdisco/model/snmp"
)

type RuiJieDriver struct {
	factory.SnmpDiscovery
}

const ruijieApMacAddr = ".1.3.6.1.4.1.4881.1.1.10.2.56.2.1.1.1.1"
const ruijieApApName = ".1.3.6.1.4.1.4881.1.1.10.2.56.2.1.1.1.2"
const ruijieApApgName = ".1.3.6.1.4.1.4881.1.1.10.2.56.2.1.1.1.3"
const ruijieApApSn = ".1.3.6.1.4.1.4881.1.1.10.2.56.2.1.1.1.32"
const ruijieApApIp = ".1.3.6.1.4.1.4881.1.1.10.2.56.2.1.1.1.33"
const ruijieApApSwVer = ".1.3.6.1.4.1.4881.1.1.10.2.56.2.1.1.1.37"
const ruijieApApPID = ".1.3.6.1.4.1.4881.1.1.10.2.56.2.1.1.1.39"

const ruijieStaMacAddr = ".1.3.6.1.4.1.4881.1.1.10.2.56.5.1.1.1.1"
const ruijieStaApMacAddr = ".1.3.6.1.4.1.4881.1.1.10.2.56.5.1.1.1.2"
const ruijieStaVlan = ".1.3.6.1.4.1.4881.1.1.10.2.56.5.1.1.1.3"
const ruijieStaIp = ".1.3.6.1.4.1.4881.1.1.10.2.56.5.1.1.1.5"
const ruijieStaSsid = ".1.3.6.1.4.1.4881.1.1.10.2.56.5.1.1.1.7"
const ruijieStaLinkRate = ".1.3.6.1.4.1.4881.1.1.10.2.56.5.1.1.1.18"
const ruijieStaCurChan = ".1.3.6.1.4.1.4881.1.1.10.2.56.5.1.1.1.19"
const ruijieStaRssi = ".1.3.6.1.4.1.4881.1.1.10.2.56.5.1.1.1.21"
const ruijieStaUsername = ".1.3.6.1.4.1.4881.1.1.10.2.56.5.1.1.1.22"
const ruijieStaOnlineTime = ".1.3.6.1.4.1.4881.1.1.10.2.56.5.1.1.1.24"

func NewRuiJieDriver(sc *snmp.SnmpConfig) (*RuiJieDriver, error) {
	session, err := factory.NewSnmpSession(sc)
	if err != nil {
		return nil, err
	}
	return &RuiJieDriver{
		factory.SnmpDiscovery{
			Session:   session,
			IpAddress: session.Target,
		},
	}, nil
}

func (d *RuiJieDriver) APs() (ap []*device.Ap, errors []string) {
	apIP, errApIP := d.Session.BulkWalkAll(ruijieApApIp)
	apMac, errApMac := d.Session.BulkWalkAll(ruijieApMacAddr)
	apName, errApName := d.Session.BulkWalkAll(ruijieApApName)
	apGroup, errApGroup := d.Session.BulkWalkAll(ruijieApApgName)
	apSerialNumber, errApSerialNumber := d.Session.BulkWalkAll(ruijieApApSn)
	apVersion, errApVersion := d.Session.BulkWalkAll(ruijieApApSwVer)
	apPID, errApPID := d.Session.BulkWalkAll(ruijieApApPID)
	if errApMac != nil || errApName != nil || errApGroup != nil || errApSerialNumber != nil || errApVersion != nil || errApPID != nil || errApIP != nil {
		errors = append(errors, errApMac.Error())
		errors = append(errors, errApName.Error())
		errors = append(errors, errApGroup.Error())
		errors = append(errors, errApSerialNumber.Error())
		errors = append(errors, errApVersion.Error())
		errors = append(errors, errApPID.Error())
		errors = append(errors, errApIP.Error())
	}
	indexApIP := factory.ExtractString(ruijieApApIp, apIP)
	indexApMac := factory.ExtractMacAddress(ruijieApMacAddr, apMac)
	indexApName := factory.ExtractString(ruijieApApName, apName)
	indexApGroup := factory.ExtractString(ruijieApApgName, apGroup)
	indexApSerialNumber := factory.ExtractString(ruijieApApSn, apSerialNumber)
	indexApVersion := factory.ExtractString(ruijieApApSwVer, apVersion)
	indexApPID := factory.ExtractString(ruijieApApPID, apPID)
	for i, v := range indexApIP {
		apVer := indexApVersion[i]
		group := indexApGroup[i]
		ap = append(ap, &device.Ap{
			Name:         indexApName[i],
			ManagementIp: v,
			MacAddress:   indexApMac[i],
			DeviceModel:  indexApPID[i],
			SerialNumber: indexApSerialNumber[i],
			OsVersion:    &apVer,
			GroupName:    &group,
		})
	}
	return ap, errors
}
