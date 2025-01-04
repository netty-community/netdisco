package driver

import (
	"fmt"

	"github.com/netty-community/netdisco/factory"
	"github.com/netty-community/netdisco/model/device"
	"github.com/netty-community/netdisco/model/snmp"
)

const fnSysSerial string = ".1.3.6.1.4.1.12356.100.1.1.1.0"
const fgsSysVersion string = ".1.3.6.1.4.1.12356.101.4.1.1.0"
const hardwareModelName string = ".1.3.6.1.2.1.47.1.2.1.1.2.1"

type FortiNetDriver struct {
	factory.SnmpDiscovery
}

func NewFortiNetDriver(sc *snmp.SnmpConfig) (*FortiNetDriver, error) {
	session, err := factory.NewSnmpSession(sc)
	if err != nil {
		return nil, err
	}
	return &FortiNetDriver{
		factory.SnmpDiscovery{
			Session:   session,
			IpAddress: session.Target},
	}, nil
}

func (f *FortiNetDriver) ChassisId() (string, error) {
	ifName, err := f.Session.BulkWalkAll(factory.IfDescr)
	if err != nil {
		return "", fmt.Errorf("get chassisId failed from %s", f.IpAddress)
	}
	// since fortinet not support lldp local chassis id, use mgmt1 mac address replaced
	indexIfName := factory.ExtractString(factory.IfDescr, ifName)
	index := ""
	for i, v := range indexIfName {
		if v == "mgmt1" {
			index = i
			break
		}
	}
	localChassis, err := f.Session.Get([]string{factory.IfPhysAddr + index})
	if err != nil {
		return "", fmt.Errorf("get chassisId failed from %s", f.IpAddress)
	}
	localChassisId := fmt.Sprintf("%s", localChassis.Variables[0].Value)

	return localChassisId, nil
}

func (f *FortiNetDriver) Entities() (entities []*device.Entity, errors []string) {
	version, err := f.Session.Get([]string{fgsSysVersion})
	if err != nil {
		return nil, []string{fmt.Sprintf("get OsVersion failed from %s", f.IpAddress)}
	}

	serial, err := f.Session.Get([]string{fnSysSerial})
	if err != nil {
		return nil, []string{fmt.Sprintf("get SerialNumber failed from %s", f.IpAddress)}
	}

	model, err := f.Session.Get([]string{hardwareModelName})
	if err != nil {
		return nil, []string{fmt.Sprintf("get hardwareModelName failed from %s", f.IpAddress)}
	}

	entities = append(entities, &device.Entity{
		EntityPhysicalSoftwareRev: fmt.Sprintf("%s", version.Variables[0].Value),
		EntityPhysicalSerialNum:   fmt.Sprintf("%s", serial.Variables[0].Value),
		EntityPhysicalName:        fmt.Sprintf("%s", model.Variables[0].Value),
	})
	return entities, []string{}
}

func (f *FortiNetDriver) LldpNeighbors() (lldp []*device.LldpNeighbor, errors []string) {
	localChassisId, err := f.ChassisId()
	if err != nil {
		errors = append(errors, err.Error())
		return nil, errors
	}
	localHostname, _ := f.SysName()
	localIfName, errIfName := f.Session.BulkWalkAll(factory.LldpLocPortId)
	localIfDescr, errIfDescr := f.Session.BulkWalkAll(factory.LldpLocPortDesc)
	remoteChassisId, errRemChassisId := f.Session.BulkWalkAll(factory.LldpRemChassisId)
	remoteHostname, errRemHostname := f.Session.BulkWalkAll(factory.LldpRemSysName)
	remoteIfName, errRemIfName := f.Session.BulkWalkAll(factory.LldpRemPortId)
	remoteIfDescr, errRemIfDescr := f.Session.BulkWalkAll(factory.LldpRemPortDesc)
	if errIfName != nil || errIfDescr != nil || errRemChassisId != nil || errRemIfName != nil || errRemIfDescr != nil {
		errors = append(errors, errIfName.Error())
		errors = append(errors, errIfDescr.Error())
		errors = append(errors, errRemHostname.Error())
		errors = append(errors, errRemChassisId.Error())
		errors = append(errors, errRemIfName.Error())
		errors = append(errors, errRemIfDescr.Error())
	}
	indexIfName := factory.ExtractString(factory.LldpLocPortId, localIfName)
	indexIfDescr := factory.ExtractString(factory.LldpLocPortDesc, localIfDescr)
	indexRemChassisId := factory.ExtractMacAddressWithShift(factory.LldpRemChassisId, -2, remoteChassisId)
	indexRemoteHostname := factory.ExtractStringWithShift(factory.LldpRemSysName, -2, remoteHostname)
	indexRemoteIfName := factory.ExtractStringWithShift(factory.LldpRemPortId, -2, remoteIfName)
	indexRemoteIfDescr := factory.ExtractStringWithShift(factory.LldpRemPortDesc, -2, remoteIfDescr)

	for i, v := range indexRemChassisId {
		if v == localChassisId {
			lldp = append(lldp, &device.LldpNeighbor{
				LocalChassisId:  localChassisId,
				LocalHostname:   localHostname,
				LocalIfName:     indexIfName[i],
				LocalIfDescr:    indexIfDescr[i],
				RemoteChassisId: indexRemChassisId[i],
				RemoteHostname:  indexRemoteHostname[i],
				RemoteIfName:    indexRemoteIfName[i],
				RemoteIfDescr:   indexRemoteIfDescr[i],
			})
		}
	}
	return lldp, errors
}

func (f *FortiNetDriver) BasicInfo() *factory.DiscoveryBasicResponse {
	sysDescr, sysError := f.SysDescr()
	sysName, sysNameError := f.SysName()
	chassisId, chassisIdError := f.ChassisId()

	response := &factory.DiscoveryBasicResponse{
		Hostname:  sysName,
		SysDescr:  sysDescr,
		ChassisId: chassisId,
	}

	if sysError != nil {
		response.Errors = append(response.Errors, sysError.Error())
	}
	if sysNameError != nil {
		response.Errors = append(response.Errors, sysNameError.Error())
	}
	if chassisIdError != nil {
		response.Errors = append(response.Errors, chassisIdError.Error())
	}
	return response
}
