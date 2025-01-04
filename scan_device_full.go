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
	"sync"

	"github.com/netty-community/netdisco/factory"
	"github.com/netty-community/netdisco/model/snmp"
)


func EnrichDeviceFull(targets []*snmp.SnmpConfig) []*factory.DispatchResponse {
	var responses []*factory.DispatchResponse
	var wg sync.WaitGroup
	ch := make(chan struct{}, MaxEnrichDeviceChannelPoolSize)
	for _, target := range targets {
		ch <- struct{}{}
		wg.Add(1)
		go func(target *snmp.SnmpConfig) {
			defer wg.Done()
			targetResponse := enrichDevice(target)

			responses = append(responses, targetResponse)
			<-ch
		}(target)
	}
	wg.Wait()
	return responses
}

func enrichDeviceFull(config *snmp.SnmpConfig) *factory.DispatchResponse {
	var response = &factory.DispatchResponse{}
	response.IpAddress = config.IpAddress
	sd, deviceModel, err := NewNetDisco(config).Driver()
	if err != nil || sd == nil {
		response.SnmpReachable = false
	} else {
		response.SnmpReachable = true
	}
	icmp := factory.IcmpReachable(config.IpAddress)
	ssh := factory.SshReachable(config.IpAddress)
	response.IcmpReachable = icmp
	response.SshReachable = ssh
	if !response.SnmpReachable {
		return response
	}
	response.DeviceModel = deviceModel
	sysDescr, sysError := sd.SysDescr()
	sysUpTime, sysUpTimeError := sd.SysUpTime()
	sysName, sysNameError := sd.SysName()
	chassisId, chassisIdError := sd.ChassisId()
	interfaces, interfacesError := sd.Interfaces()
	entities, entitiesError := sd.Entities()
	lldp, lldpError := sd.LldpNeighbors()
	macAddress, macAddressError := sd.MacAddressTable()
	arp, arpError := sd.ArpTable()
	arp = factory.EnrichArpInfo(arp, interfaces)
	vlan, VlanError := sd.Vlans()
	vlan = factory.EnrichVlanInfo(vlan, interfaces)
	macAddress_ := factory.EnrichMacAddress(macAddress, interfaces, lldp, arp)

	disc := &factory.DiscoveryResponse{
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
		disc.Errors = append(disc.Errors, sysError.Error())
	}
	if sysUpTimeError != nil {
		disc.Errors = append(disc.Errors, sysUpTimeError.Error())
	}
	if sysNameError != nil {
		disc.Errors = append(disc.Errors, sysNameError.Error())
	}
	if chassisIdError != nil {
		disc.Errors = append(disc.Errors, chassisIdError.Error())
	}
	if interfacesError != nil {
		disc.Errors = append(disc.Errors, interfacesError...)
	}
	if entitiesError != nil {
		disc.Errors = append(disc.Errors, entitiesError...)
	}
	if lldpError != nil {
		disc.Errors = append(disc.Errors, lldpError...)
	}
	if macAddressError != nil {
		disc.Errors = append(disc.Errors, macAddressError...)
	}
	if arpError != nil {
		disc.Errors = append(disc.Errors, arpError...)
	}
	if VlanError != nil {
		disc.Errors = append(disc.Errors, VlanError...)
	}
	response.Data = disc
	return response
}
