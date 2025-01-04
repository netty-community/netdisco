// Copyright 2024 wangxin.jeffry@gmail.com
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package netdisco

import (
	"strings"

	manufacturer "github.com/netty-community/netdisco/model/manufacturer"
	"github.com/netty-community/netdisco/model/snmp"

	dt "github.com/netty-community/netdisco/devicemodel"
	s "github.com/netty-community/netdisco/devicemodel/sysobjectid"
	"github.com/netty-community/netdisco/driver"
	"github.com/netty-community/netdisco/factory"
	platform "github.com/netty-community/netdisco/model/platform"
)

func getDeviceModel(sysObjId string) *dt.DeviceModel {
	privateEnterPriseId := strings.Split(strings.Split(sysObjId, ".1.3.6.1.4.1.")[1], ".")[0]
	enterprise := manufacturer.GetManufacturerByEnterpriseId(privateEnterPriseId)
	if enterprise == manufacturer.Unknown {
		return &dt.DeviceModel{
			Platform:     platform.Unknown,
			Manufacturer: manufacturer.Unknown,
			DeviceModel:  dt.UnknownDeviceModel,
		}
	}
	return getDeviceModelFromManufacturer(enterprise, sysObjId)
}

func getDeviceModelFromManufacturer(mf manufacturer.Manufacturer, sysObjId string) *dt.DeviceModel {
	switch mf {
	case manufacturer.Cisco:
		return s.CiscoDeviceModel(sysObjId)
	case manufacturer.Huawei:
		return s.HuaweiDeviceModel(sysObjId)
	case manufacturer.Aruba:
		return s.ArubaDeviceModel(sysObjId)
	case manufacturer.Arista:
		return s.AristaDeviceModel(sysObjId)
	case manufacturer.H3C:
		return s.H3CDeviceModel(sysObjId)
	case manufacturer.RuiJie:
		return s.RuiJieDeviceModel(sysObjId)
	case manufacturer.PaloAlto:
		return s.PaloAltoDeviceModel(sysObjId)
	case manufacturer.FortiNet:
		return s.FortiNetDeviceModel(sysObjId)
	case manufacturer.Netgear:
		return s.NetgearDeviceModel(sysObjId)
	case manufacturer.TPLink:
		return s.TPLinkDeviceModel(sysObjId)
	case manufacturer.Ruckus:
		return s.RuckusDeviceModel(sysObjId)
	case manufacturer.Juniper:
		return s.JuniperDeviceModel(sysObjId)
	case manufacturer.CheckPoint:
		return s.CheckPointDeviceModel(sysObjId)
	case manufacturer.F5:
		return s.F5DeviceModel(sysObjId)
	case manufacturer.Extreme:
		return s.ExtremeDeviceModel(sysObjId)
	case manufacturer.MikroTik:
		return s.MikroTikDeviceModel(sysObjId)
	}

	return &dt.DeviceModel{
		Platform:     platform.Unknown,
		Manufacturer: mf,
		DeviceModel:  dt.UnknownDeviceModel,
	}
}

func getFactory(platformType platform.Platform, sc *snmp.SnmpConfig) (SnmpDriver, error) {
	var snmpDriver SnmpDriver
	var err error

	switch platformType {
	case platform.CiscoIos:
		snmpDriver, err = driver.NewCiscoIosDriver(sc)
	case platform.CiscoIosXE:
		snmpDriver, err = driver.NewCiscoIosXEDriver(sc)
	case platform.CiscoIosXR:
		snmpDriver, err = driver.NewCiscoIosXRDriver(sc)
	case platform.CiscoNexusOS:
		snmpDriver, err = driver.NewCiscoNexusOSDriver(sc)
	case platform.Huawei:
		snmpDriver, err = driver.NewHuaweiDriver(sc)
	case platform.Aruba:
		snmpDriver, err = driver.NewArubaDriver(sc)
	case platform.ArubaOSSwitch:
		snmpDriver, err = driver.NewArubaOSSwitchDriver(sc)
	case platform.Arista:
		snmpDriver, err = driver.NewAristaDriver(sc)
	case platform.RuiJie:
		snmpDriver, err = driver.NewRuiJieDriver(sc)
	case platform.H3C:
		snmpDriver, err = driver.NewH3CDriver(sc)
	case platform.FortiNet:
		snmpDriver, err = driver.NewFortiNetDriver(sc)
	case platform.PaloAlto:
		snmpDriver, err = driver.NewPaloAltoDriver(sc)
	case platform.Juniper:
		snmpDriver, err = driver.NewJuniperDriver(sc)
	case platform.Netgear:
		snmpDriver, err = driver.NewNetgearDriver(sc)
	case platform.TPLink:
		snmpDriver, err = driver.NewTPLinkDriver(sc)
	case platform.Ruckus:
		snmpDriver, err = driver.NewRuckusDriver(sc)
	case platform.F5:
		snmpDriver, err = driver.NewF5Driver(sc)
	case platform.CheckPoint:
		snmpDriver, err = driver.NewCheckPointDriver(sc)
	case platform.Extreme:
		snmpDriver, err = driver.NewExtremeDriver(sc)
	case platform.MikroTik:
		snmpDriver, err = driver.NewMikroTikDriver(sc)
	case platform.Unknown:
		snmpDriver, err = factory.NewSnmpDiscovery(sc)
	default:
		snmpDriver, err = factory.NewSnmpDiscovery(sc)
	}

	if err != nil {
		return nil, err
	}

	return snmpDriver, nil
}
