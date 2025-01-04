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

var MaxApDiscoveryChannelPoolSize = 50

func WlanApDiscovery(targets []*snmp.SnmpConfig) []*factory.DispatchApScanResponse {
	var responses []*factory.DispatchApScanResponse
	var wg sync.WaitGroup
	ch := make(chan struct{}, MaxApDiscoveryChannelPoolSize)
	for _, target := range targets {
		ch <- struct{}{}
		wg.Add(1)
		go func(target *snmp.SnmpConfig) {
			defer wg.Done()
			targetResponse := discoverAp(target)

			responses = append(responses, targetResponse)
			<-ch
		}(target)
	}
	wg.Wait()
	return responses
}

func discoverAp(config *snmp.SnmpConfig) *factory.DispatchApScanResponse {
	var response = &factory.DispatchApScanResponse{}
	response.IpAddress = config.IpAddress
	driver, deviceModel, err := NewNetDisco(config).Driver()
	if err != nil || driver == nil {
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
	if driver != nil {
		aps, errors := driver.APs()
		response.DeviceModel = deviceModel
		response.Errors = errors
		response.Data = aps
	}
	return response

}
