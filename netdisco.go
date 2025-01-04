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
	"fmt"

	"github.com/gosnmp/gosnmp"
	"github.com/netty-community/netdisco/devicemodel"
	"github.com/netty-community/netdisco/factory"
	"github.com/netty-community/netdisco/logger"
	"github.com/netty-community/netdisco/model/snmp"
	"go.uber.org/zap"
)

type Dispatcher struct {
	Target *snmp.SnmpConfig
}

// SnmpReachable returns true if the target is reachable via SNMP.
func (d *Dispatcher) snmpReachable(session *gosnmp.GoSNMP) bool {
	// Try to get the sysName from the target, if it fails, it's not reachable.
	result, err := session.Get([]string{factory.SysName})
	if err != nil || result == nil {
		logger.Logger.Info("[netdisco]: SnmpReachable failed", zap.Error(err), zap.String("target", session.Target))
		return false
	}
	// If the result is not empty, the target is reachable.
	return len(result.Variables) > 0
}

// Session returns a new GoSNMP session.
func (d *Dispatcher) session(config *snmp.SnmpConfig) (*gosnmp.GoSNMP, error) {
	session, err := factory.NewSnmpSession(config)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// SysObjectID returns the sysObjectID from RFC-1213MIB
func (d *Dispatcher) sysObjectID(session *gosnmp.GoSNMP) string {
	oid := []string{factory.SysObjectID}
	result, err := session.Get(oid)
	if err != nil {
		return ""
	}

	for _, variable := range result.Variables {
		if variable.Type == gosnmp.ObjectIdentifier {
			return fmt.Sprintf("%s", variable.Value)
		}
	}
	return ""
}

// Driver returns a new SnmpDriver.
func (d *Dispatcher) Driver() (SnmpDriver, *devicemodel.DeviceModel, error) {
	snmpReachable := true
	session, err := d.session(d.Target)
	if err != nil || session == nil {
		snmpReachable = false
	} else {
		snmpReachable = d.snmpReachable(session)
	}
	result := map[string]string{
		"snmpReachable": fmt.Sprintf("%t", snmpReachable),
	}

	if !snmpReachable {
		return nil, nil, fmt.Errorf("[netdisco]: target %s snmp unreachable: %v", d.Target.IpAddress, result)
	}
	if d.Target.Platform != nil && *d.Target.Platform != "" {
		driver, err := getFactory(*d.Target.Platform, d.Target)
		sysObjectId := d.sysObjectID(session)
		deviceModel := getDeviceModel(sysObjectId)
		return driver, deviceModel, err
	}
	sysObjectId := d.sysObjectID(session)
	if sysObjectId == "" {
		return nil, nil, fmt.Errorf("[netdisco]: Unable to get sysObjectId, target manufacturer is not supported: %s", d.Target.IpAddress)
	}
	deviceModel := getDeviceModel(sysObjectId)
	d.Target.SysObjectId = sysObjectId
	driver, err := getFactory(deviceModel.Platform, d.Target)
	if driver == nil {
		return nil, nil, fmt.Errorf("[netdisco]: Unable to get driver, target manufacturer is not supported: %s, sysObjectId: %s", d.Target.IpAddress, sysObjectId)
	}
	return driver, deviceModel, err
}

func NewNetDisco(config *snmp.SnmpConfig) *Dispatcher {
	return &Dispatcher{Target: config}
}
