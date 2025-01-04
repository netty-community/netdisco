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

package netdisco_test

import (
	"testing"

	"github.com/gosnmp/gosnmp"
	"github.com/netty-community/netdisco"
	"github.com/netty-community/netdisco/helpers/network"
	"github.com/netty-community/netdisco/logger"
	"github.com/netty-community/netdisco/model/snmp"
)

func TestDriver(t *testing.T) {
	loggerConfig := logger.LogConfig{
		Formatter: "text",
	}
	logger.InitLogger(&loggerConfig)
	community := "public"
	target := snmp.SnmpConfig{
		IpAddress:      "127.0.0.1",
		Version:        gosnmp.Version2c,
		Port:           161,
		Community:      &community,
		Timeout:        3,
		MaxRepetitions: 10,
	}
	disco, _, err := netdisco.NewNetDisco(&target).Driver()
	if err != nil {
		t.Fatalf("failed to create netdisco: %s", err)
	}
	response := disco.BasicInfo()

	t.Logf("response: %v", response)
}

func TestDiscovery(t *testing.T) {
	loggerConfig := logger.LogConfig{
		Formatter: "text",
	}
	logger.InitLogger(&loggerConfig)
	community := "public"
	cidr := "192.168.0.0/24"
	targets, err := network.CIDR2IpStrings(cidr)
	if err != nil {
		t.Fatalf("received wrong ip range %s", cidr)
	}
	snmpConfigs := make([]*snmp.SnmpConfig, 0)
	port := 161
	timeout := 3
	maxRepetitions := 10
	for _, target := range targets {
		snmpConfigs = append(snmpConfigs, &snmp.SnmpConfig{
			IpAddress:      target,
			Port:           uint16(port),
			Timeout:        uint8(timeout),
			Community:      &community,
			Version:        gosnmp.Version2c,
			MaxRepetitions: maxRepetitions,
		})
	}

	result := netdisco.NetworkDeviceDiscovery(snmpConfigs)

	t.Logf("result: %v", result)
}
