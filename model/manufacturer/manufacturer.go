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

// package manufacturer define built-in supported manufacturer

package manufacturer

import "github.com/netty-community/netdisco/model/platform"

type Manufacturer string

const (
	Cisco      Manufacturer = "Cisco"
	Huawei     Manufacturer = "Huawei"
	Aruba      Manufacturer = "Aruba"
	Arista     Manufacturer = "Arista"
	H3C        Manufacturer = "H3C"
	RuiJie     Manufacturer = "Ruijie"
	PaloAlto   Manufacturer = "Palo Alto"
	FortiNet   Manufacturer = "Fortinet"
	Netgear    Manufacturer = "Netgear"
	TPLink     Manufacturer = "TP-Link"
	Ruckus     Manufacturer = "Ruckus"
	Juniper    Manufacturer = "Juniper"
	CheckPoint Manufacturer = "Check Point"
	F5         Manufacturer = "F5 Networks"
	Extreme    Manufacturer = "Extreme"
	MikroTik   Manufacturer = "MikroTik"
	Unknown    Manufacturer = "Unknown"
)

func SupportManufacturer() []Manufacturer {
	return []Manufacturer{
		Cisco, Huawei, Aruba, Arista, H3C, RuiJie, PaloAlto,
		FortiNet, Netgear, TPLink, Ruckus, Juniper, CheckPoint,
		F5, Extreme, MikroTik, Unknown,
	}
}

func GetManufacturerByEnterpriseId(etpId string) Manufacturer {
	idMapping := map[string]Manufacturer{
		"2011":  Huawei,
		"56813": Huawei,
		"9":     Cisco,
		"5771":  Cisco,
		"5842":  Cisco,
		"53683": Cisco,
		"14823": Aruba,
		"30065": Arista,
		"4881":  RuiJie,
		"61878": H3C,
		"25506": H3C,
		"25461": PaloAlto,
		"12356": FortiNet,
		"2636":  Juniper,
		"4562":  Netgear,
		"11863": TPLink,
		"25053": Ruckus,
		"2620":  CheckPoint,
		"12276": F5,
		"1916":  Extreme,
		"14988": MikroTik,
	}
	if etpId == "" {
		return Unknown
	}
	mf, ok := idMapping[etpId]
	if !ok {
		return Unknown
	}
	return mf
}

func GetAllManufacturerPlatform() map[Manufacturer][]platform.Platform {
	return map[Manufacturer][]platform.Platform{
		Cisco:      {platform.CiscoIos, platform.CiscoIosXE, platform.CiscoIosXR, platform.CiscoNexusOS},
		Huawei:     {platform.Huawei, platform.HuaweiCE, platform.HuaweiFM},
		Aruba:      {platform.Aruba, platform.ArubaOSSwitch},
		Arista:     {platform.Arista},
		RuiJie:     {platform.RuiJie},
		H3C:        {platform.H3C},
		PaloAlto:   {platform.PaloAlto},
		FortiNet:   {platform.FortiNet},
		Netgear:    {platform.Netgear},
		TPLink:     {platform.TPLink},
		Ruckus:     {platform.Ruckus},
		Juniper:    {platform.Juniper},
		CheckPoint: {platform.CheckPoint},
		F5:         {platform.F5},
		Extreme:    {platform.Extreme},
		MikroTik:   {platform.MikroTik},
		Unknown:    {platform.Unknown},
	}
}
