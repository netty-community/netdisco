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
	"github.com/netty-community/netdisco/factory"
	"github.com/netty-community/netdisco/model/device"
)

type SnmpDriver interface {
	SysDescr() (string, error)
	SysObjectID() (string, error)
	SysUpTime() (uint64, error)
	SysName() (string, error)
	ChassisId() (string, error)
	IfPortMode() map[uint64]string
	Interfaces() (interfaces []*device.DeviceInterface, errors []string)
	LldpNeighbors() (lldp []*device.LldpNeighbor, errors []string)
	Entities() (entities []*device.Entity, errors []string)
	MacAddressTable() (macTable *map[uint64][]string, errors []string)
	ArpTable() (arp []*device.ArpItem, errors []string)
	Vlans() (vlans []*device.VlanItem, errors []string)
	APs() (ap []*device.Ap, errors []string)
	BasicInfo() *factory.DiscoveryBasicResponse
}
