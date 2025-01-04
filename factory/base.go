package factory

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/netty-community/netdisco/logger"
	"github.com/netty-community/netdisco/model/device"
	"github.com/netty-community/netdisco/model/snmp"
	probing "github.com/prometheus-community/pro-bing"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type SnmpDiscovery struct {
	Session   *gosnmp.GoSNMP
	IpAddress string
}

func NewSnmpSession(config *snmp.SnmpConfig) (*gosnmp.GoSNMP, error) {
	var snmpSession *gosnmp.GoSNMP
	if !config.Validate() {
		return nil, fmt.Errorf("invalid snmp config parameters for %s", config.IpAddress)
	}
	snmpSession = &gosnmp.GoSNMP{
		Target:   config.IpAddress,
		Port:     config.Port,
		Timeout:  time.Duration(config.Timeout) * time.Second,
		Retries:  2,
		MaxOids:  int(config.MaxRepetitions),
		Version:  config.Version,
		MsgFlags: gosnmp.AuthPriv,
	}

	switch config.Version {
	case gosnmp.Version2c:
		snmpSession.Community = *config.Community
	case gosnmp.Version3:
		snmpSession.SecurityParameters = &gosnmp.UsmSecurityParameters{
			UserName:                 *config.V3Params.SecurityName,
			AuthenticationProtocol:   config.V3Params.AuthProtocol,
			AuthenticationPassphrase: *config.V3Params.AuthPassword,
			PrivacyProtocol:          config.V3Params.PrivProtocol,
			PrivacyPassphrase:        *config.V3Params.PrivPassword,
		}
	default:
		return nil, fmt.Errorf("unsupported snmp version: %d", config.Version)
	}

	err := snmpSession.Connect()
	if err != nil {
		logger.Logger.Info("snmp connect error", zap.String("ip", config.IpAddress), zap.Error(err))
		return nil, err
	}
	return snmpSession, nil
}

func NewSnmpDiscovery(sc *snmp.SnmpConfig) (*SnmpDiscovery, error) {
	session, err := NewSnmpSession(sc)
	if err != nil {
		return nil, err
	}
	return &SnmpDiscovery{
		Session:   session,
		IpAddress: session.Target,
	}, nil
}

// IcmpReachable returns true if the target is reachable via ICMP.
// linux need privilege for udp
func IcmpReachable(address string) bool {
	pinger, err := probing.NewPinger(address)
	if err != nil {
		logger.Logger.Info("IcmpReachable failed", zap.String("target", address), zap.Error(err))
		return false
	}
	pinger.Count = 2
	pinger.Interval = time.Duration(100) * time.Millisecond
	pinger.Timeout = time.Second
	err = pinger.Run()
	if err != nil {
		return false
	}
	return pinger.Statistics().PacketsRecv > 0
}

// SshReachable returns true if the target is reachable via SSH. default port is 22
func SshReachable(address string) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", address+":22", timeout)
	if err != nil {
		logger.Logger.Info("SshReachable failed", zap.String("target", address), zap.Error(err))
		return false
	}
	if conn == nil {
		logger.Logger.Info("SshReachable failed", zap.String("target", address), zap.String("reason", "connection is nil"))
		return false
	}
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()
	return true
}

// SysDescr returns the system description from RFC-1213MIB
func (sd *SnmpDiscovery) SysDescr() (string, error) {
	oid := []string{SysDescr}
	result, err := sd.Session.Get(oid)
	if err != nil {
		return "", fmt.Errorf("failed to get system description from %s: %w", sd.IpAddress, err)
	}

	for _, variable := range result.Variables {
		if variable.Type == gosnmp.OctetString {
			return fmt.Sprintf("%s", variable.Value), nil
		}
	}

	return "", fmt.Errorf("failed to get system description from %s: no variable found", sd.IpAddress)
}

// SysObjectID returns the sysObjectID from RFC-1213MIB
func (sd *SnmpDiscovery) SysObjectID() (string, error) {
	oid := []string{SysObjectID}
	result, err := sd.Session.Get(oid)
	if err != nil {
		return "", err
	}

	for _, variable := range result.Variables {
		if variable.Type == gosnmp.OctetString {
			return fmt.Sprintf("%s", variable.Value), nil
		}
	}
	return "", fmt.Errorf("failed to get sysObjectID from %s for OIDs: %v", sd.IpAddress, oid)
}

// SysUpTime returns the sysUpTime from RFC-1213MIB
func (sd *SnmpDiscovery) SysUpTime() (uint64, error) {
	ticksVariables, err := sd.Session.Get([]string{SysUpTime})
	if err != nil {
		return 0, err
	}

	var ticks uint64
	for _, variable := range ticksVariables.Variables {
		if variable.Type == gosnmp.TimeTicks {
			switch v := variable.Value.(type) {
			case uint32:
				ticks = uint64(v)
			case uint16:
				ticks = uint64(v)
			case uint64:
				ticks = v
			default:
				return 0, fmt.Errorf("unexpected type for sysUpTime: %T", variable.Value)
			}
			return ticks, nil
		}
	}

	return 0, fmt.Errorf("failed to get sysUpTime from %s", sd.IpAddress)
}

// SysName returns the sysName from RFC-1213MIB
func (sd *SnmpDiscovery) SysName() (string, error) {
	sysNameVars, err := sd.Session.Get([]string{SysName})
	if err != nil {
		return "", err
	}

	for _, variable := range sysNameVars.Variables {
		if variable.Type == gosnmp.OctetString {
			return fmt.Sprintf("%s", variable.Value), nil
		}
	}

	return "", fmt.Errorf("failed to get sysName from %s", sd.IpAddress)
}

// ChassisId returns the chassisId from RFC-1213MIB
func (sd *SnmpDiscovery) ChassisId() (string, error) {
	result, err := sd.Session.Get([]string{LldpLocChassisId})
	if err != nil {
		return "", err
	}
	for _, variable := range result.Variables {
		if variable.Type == gosnmp.OctetString {
			return hex2mac(variable.Value.([]byte)), nil
		}
	}
	return "", fmt.Errorf("failed to get ChassisId from %s", sd.IpAddress)
}

func (sd *SnmpDiscovery) IfPortMode() map[uint64]string {
	// TODO: implement me
	return nil
}

// collect interfaces via IF-MIB
func (sd *SnmpDiscovery) Interfaces() (interfaces []*device.DeviceInterface, errors []string) {
	ifIndex, errIfIndex := sd.Session.BulkWalkAll(IfIndex)
	if len(ifIndex) == 0 || errIfIndex != nil {
		return nil, []string{fmt.Sprintf("failed to get ifIndex from %s", sd.IpAddress)}
	}
	ifName, errIfName := sd.Session.BulkWalkAll(IfDescr)
	ifDesc, errIfDesc := sd.Session.BulkWalkAll(IfAlias)
	ifMtu, errIfMtu := sd.Session.BulkWalkAll(IfMtu)
	ifSpeed, errIfSpeed := sd.Session.BulkWalkAll(IfHighSpeed)
	ifHighSpeed, errIfHighSpeed := sd.Session.BulkWalkAll(IfHighSpeed)
	ifPhysAddr, errIfPhysAddr := sd.Session.BulkWalkAll(IfPhysAddr)
	ifType, errIfType := sd.Session.BulkWalkAll(IfType)
	ifAdminStatus, errIfAdminStatus := sd.Session.BulkWalkAll(IfAdminStatus)
	ifOperStatus, errIfOperStatus := sd.Session.BulkWalkAll(IfOperStatus)
	ifLastChange, errIfLastChange := sd.Session.BulkWalkAll(IfLastChange)
	ifAddrIndex, errIfAddrIndex := sd.Session.BulkWalkAll(IfAdEntIfIndex)
	ifAddrNetMask, errIfAddrNetMask := sd.Session.BulkWalkAll(IfAdEntNetMask)

	if errIfName != nil || errIfDesc != nil ||
		errIfMtu != nil || errIfSpeed != nil || errIfHighSpeed != nil ||
		errIfPhysAddr != nil || errIfType != nil || errIfAdminStatus != nil || errIfLastChange != nil ||
		errIfOperStatus != nil || errIfAddrIndex != nil || errIfAddrNetMask != nil {
		errors = append(errors, errIfName.Error())
		errors = append(errors, errIfDesc.Error())
		errors = append(errors, errIfMtu.Error())
		errors = append(errors, errIfSpeed.Error())
		errors = append(errors, errIfHighSpeed.Error())
		errors = append(errors, errIfPhysAddr.Error())
		errors = append(errors, errIfType.Error())
		errors = append(errors, errIfAdminStatus.Error())
		errors = append(errors, errIfOperStatus.Error())
		errors = append(errors, errIfLastChange.Error())
		errors = append(errors, errIfAddrIndex.Error())
		errors = append(errors, errIfAddrNetMask.Error())
	}
	indexIfIndex := ExtractInteger(IfIndex, ifIndex)
	indexIfName := ExtractString(IfDescr, ifName)
	indexIfDesc := ExtractString(IfAlias, ifDesc)
	indexIfMtu := ExtractInteger(IfMtu, ifMtu)
	indexIfSpeed := ExtractInteger(IfHighSpeed, ifSpeed)
	indexIfHighSpeed := ExtractInteger(IfHighSpeed, ifHighSpeed)
	indexIfPhysAddr := ExtractMacAddress(IfPhysAddr, ifPhysAddr)
	indexIfType := ExtractInteger(IfType, ifType)
	indexIfAdminStatus := ExtractInteger(IfAdminStatus, ifAdminStatus)
	indexIfOperStatus := ExtractInteger(IfOperStatus, ifOperStatus)
	indexIfLastChange := ExtractInteger(IfLastChange, ifLastChange)
	indexIfAddrIndex := extractIfIndex(IfAdEntIfIndex, ifAddrIndex)
	indexIfAddrNetMask := ExtractString(IfAdEntNetMask, ifAddrNetMask)
	for i, v := range indexIfIndex {
		var _ifAddrIndex string
		ifTypeString := GetIfTypeValue(indexIfType[i])
		if lo.Contains([]string{"other", "softwareLoopback"}, ifTypeString) {
			continue
		}
		itemIfAddrIndex := indexIfAddrIndex[i]
		itemIfAddrNetMask := indexIfAddrNetMask[itemIfAddrIndex]
		if itemIfAddrIndex == "" || itemIfAddrNetMask == "" {
			_ifAddrIndex = ""
		} else {
			_ifAddrIndex = strings.TrimPrefix(indexIfAddrIndex[i], ".") + "/" + strconv.Itoa(netmaskToLength(itemIfAddrNetMask))
		}
		physAddr := indexIfPhysAddr[i]
		iface := device.DeviceInterface{
			IfIndex:       v,
			IfName:        indexIfName[i],
			IfDescr:       indexIfDesc[i],
			IfType:        ifTypeString,
			IfMtu:         indexIfMtu[i],
			IfSpeed:       indexIfSpeed[i],
			IfPhysAddr:    &physAddr,
			IfAdminStatus: GetIfAdminStatusValue(indexIfAdminStatus[i]),
			IfOperStatus:  GetIfOperStatusValue(indexIfOperStatus[i]),
			IfLastChange:  indexIfLastChange[i],
			IfHighSpeed:   indexIfHighSpeed[i],
			IfIpAddress:   &_ifAddrIndex,
		}
		interfaces = append(interfaces, &iface)
	}

	return interfaces, errors
}

func (sd *SnmpDiscovery) LldpNeighbors() (lldp []*device.LldpNeighbor, errors []string) {

	localChassisId, err := sd.ChassisId()
	if err != nil {
		errors = append(errors, err.Error())
		return nil, errors
	}
	hostname, _ := sd.SysName()
	localIfName, errIfName := sd.Session.BulkWalkAll(LldpLocPortId)
	localIfDescr, errIfDescr := sd.Session.BulkWalkAll(LldpLocPortDesc)
	remoteChassisId, errRemChassisId := sd.Session.BulkWalkAll(LldpRemChassisId)
	remoteHostname, errRemHostname := sd.Session.BulkWalkAll(LldpRemSysName)
	remoteIfName, errRemIfName := sd.Session.BulkWalkAll(LldpRemPortId)
	remoteIfDescr, errRemIfDescr := sd.Session.BulkWalkAll(LldpRemPortDesc)
	if errIfName != nil || errIfDescr != nil || errRemChassisId != nil || errRemIfName != nil || errRemIfDescr != nil {
		errors = append(errors, errIfName.Error())
		errors = append(errors, errIfDescr.Error())
		errors = append(errors, errRemHostname.Error())
		errors = append(errors, errRemChassisId.Error())
		errors = append(errors, errRemIfName.Error())
		errors = append(errors, errRemIfDescr.Error())
	}
	IndexIfName := ExtractString(LldpLocPortId, localIfName)
	IndexIfDescr := ExtractString(LldpLocPortDesc, localIfDescr)
	IndexRemChassisId := ExtractMacAddressWithShift(LldpRemChassisId, -2, remoteChassisId)
	IndexRemoteHostname := ExtractStringWithShift(LldpRemSysName, -2, remoteHostname)
	IndexRemoteIfName := ExtractStringWithShift(LldpRemPortId, -2, remoteIfName)
	IndexRemoteIfDescr := ExtractStringWithShift(LldpRemPortDesc, -2, remoteIfDescr)

	for i, v := range IndexRemChassisId {
		remoteHostname := IndexRemoteHostname[i]
		if remoteHostname == "" {
			continue
		}
		neighbor := &device.LldpNeighbor{
			LocalChassisId:  localChassisId,
			LocalHostname:   hostname,
			LocalIfName:     IndexIfName["."+i],
			LocalIfDescr:    IndexIfDescr["."+i],
			RemoteChassisId: v,
			RemoteHostname:  remoteHostname,
			RemoteIfName:    IndexRemoteIfName[i],
			RemoteIfDescr:   IndexRemoteIfDescr[i],
		}
		lldp = append(lldp, neighbor)
	}
	return lldp, errors
}

func (sd *SnmpDiscovery) Entities() (entities []*device.Entity, errors []string) {
	entPhysicalClass, err := sd.Session.BulkWalkAll(EntPhysicalClass)
	if len(entPhysicalClass) == 0 || err != nil {
		if err != nil {
			errors = append(errors, err.Error())
		} else {
			errors = append(errors, fmt.Sprintf("get entity physical class failed from %s, No entities found", sd.IpAddress))
		}
		return nil, errors
	}
	IndexEntPhysicalClass := ExtractInteger(EntPhysicalClass, entPhysicalClass)
	FilteredIndexEntPhysicalClass := lo.PickByValues(IndexEntPhysicalClass, []uint64{3})
	chassidIndex := lo.Keys(FilteredIndexEntPhysicalClass)
	if len(chassidIndex) == 0 {
		errors = append(errors, fmt.Sprintf("get entity physical class failed from %s, No entities found", sd.IpAddress))
		return nil, errors
	}
	entityPhysicalDescr, errEntityPhysicalDescr := sd.Session.Get(buildOidWithIndex(EntPhysicalDescr, chassidIndex))
	entityPhysicalName, errEntityPhysicalName := sd.Session.Get(buildOidWithIndex(EntPhysicalName, chassidIndex))
	entityPhysicalSoftwareRev, errEntityPhysicalSoftwareRev := sd.Session.Get(buildOidWithIndex(EntPhysicalSoftwareRev, chassidIndex))
	entityPhysicalSerialNum, errEntityPhysicalSerialNum := sd.Session.Get(buildOidWithIndex(EntPhysicalSerialNum, chassidIndex))

	if errEntityPhysicalDescr != nil || errEntityPhysicalName != nil || errEntityPhysicalSoftwareRev != nil || errEntityPhysicalSerialNum != nil {
		errors = append(errors, errEntityPhysicalDescr.Error())
		errors = append(errors, errEntityPhysicalName.Error())
		errors = append(errors, errEntityPhysicalSoftwareRev.Error())
		errors = append(errors, errEntityPhysicalSerialNum.Error())
	}
	IndexEntityPhysicalDescr := ExtractString(EntPhysicalDescr, entityPhysicalDescr.Variables)
	IndexEntityPhysicalName := ExtractString(EntPhysicalName, entityPhysicalName.Variables)
	IndexEntityPhysicalSoftwareRev := ExtractString(EntPhysicalSoftwareRev, entityPhysicalSoftwareRev.Variables)
	IndexEntityPhysicalSerialNum := ExtractString(EntPhysicalSerialNum, entityPhysicalSerialNum.Variables)
	for i, v := range FilteredIndexEntPhysicalClass {
		entities = append(entities, &device.Entity{
			EntityPhysicalClass:       GetEntPhysicalClassValue(v),
			EntityPhysicalDescr:       IndexEntityPhysicalDescr[i],
			EntityPhysicalName:        IndexEntityPhysicalName[i],
			EntityPhysicalSoftwareRev: IndexEntityPhysicalSoftwareRev[i],
			EntityPhysicalSerialNum:   IndexEntityPhysicalSerialNum[i],
		})
	}
	return entities, errors
}

func (sd *SnmpDiscovery) MacAddressTable() (macTable *map[uint64][]string, errors []string) {
	dot1dBasePortIndex, err := sd.Session.BulkWalkAll(Dot1dBasePortIfIndex)
	if err != nil {
		return nil, []string{err.Error()}
	}
	dot1dTpFdbAddress, errDot1dTpFdbAddress := sd.Session.BulkWalkAll(Dot1dTpFdbAddress)
	dot1dTpFdbPort, errDot1dTpFdbPort := sd.Session.BulkWalkAll(Dot1dTpFdbPort)
	if errDot1dTpFdbAddress != nil || errDot1dTpFdbPort != nil {
		errors = append(errors, errDot1dTpFdbAddress.Error())
		errors = append(errors, errDot1dTpFdbPort.Error())
	}

	_IndexDot1dBasePortIfIndex := ExtractInteger(Dot1dBasePortIfIndex, dot1dBasePortIndex)
	IndexDot1dBasePortIfIndex := lo.MapValues(_IndexDot1dBasePortIfIndex, func(x uint64, _ string) string {
		return strconv.FormatUint(x, 10)
	})
	IndexDot1dBasePortIfIndex[".0"] = "0"
	indexDot1dTpFdbAddress := make(map[string]string)
	for _, v := range dot1dTpFdbAddress {
		splitValue := strings.Split(v.Name, ".")
		index := strings.Join(splitValue[len(splitValue)-6:], ".") // 最后六位时MAC地址地址的十进制表示
		if mac := hex2mac(v.Value.([]byte)); mac != "" {
			indexDot1dTpFdbAddress[index] = mac
		}
	}
	indexDot1dTpFdbPort := make(map[string]string)
	for _, v := range dot1dTpFdbPort {
		splitValue := strings.Split(v.Name, ".")
		index := strings.Join(splitValue[len(splitValue)-6:], ".")
		indexDot1dTpFdbPort[index] = strconv.FormatUint(gosnmp.ToBigInt(v.Value).Uint64(), 10)
	}

	result := make(map[string][]string)
	for _, v := range IndexDot1dBasePortIfIndex {
		result[v] = make([]string, 0)
	}

	for portIndex, macAddress := range indexDot1dTpFdbAddress {
		if _, ok := indexDot1dTpFdbPort[portIndex]; ok {
			ifIndex := IndexDot1dBasePortIfIndex["."+indexDot1dTpFdbPort[portIndex]]
			result[ifIndex] = append(result[ifIndex], macAddress)
		}
	}

	_macTable := lo.MapKeys(result, func(_ []string, x string) uint64 {
		index, _ := strconv.ParseUint(x, 10, 64)
		return index
	})

	return &_macTable, errors
}

func (sd *SnmpDiscovery) ArpTable() (arp []*device.ArpItem, errors []string) {
	arpTable, errArpTable := sd.Session.BulkWalkAll(IpNetToMediaPhysAddress)
	arpType, errArpType := sd.Session.BulkWalkAll(IpNetToMediaType)
	if errArpTable != nil || errArpType != nil {
		errors = append(errors, errArpType.Error())
		errors = append(errors, errArpTable.Error())
		return nil, errors
	}
	arpMap := ExtractMacAddress(IpNetToMediaPhysAddress, arpTable)
	arpTypeMap := ExtractInteger(IpNetToMediaType, arpType)
	results := make([]*device.ArpItem, 0)
	for key, value := range arpMap {
		ifIndex, address := getIfIndexAndAddress(key)
		results = append(results, &device.ArpItem{
			IpAddress:  address,
			MacAddress: value,
			Type:       GetArpTypeValue(arpTypeMap[key]),
			IfIndex:    ifIndex,
		})
	}
	return results, nil
}

func (sd *SnmpDiscovery) Vlans() (vlan []*device.VlanItem, errors []string) {
	results := make([]*device.VlanItem, 0)
	return results, nil
}

func (sd *SnmpDiscovery) APs() (ap []*device.Ap, errors []string) {
	// need implement in vendor driver
	results := make([]*device.Ap, 0)
	return results, nil
}

func (sd *SnmpDiscovery) BasicInfo() *DiscoveryBasicResponse {
	sysDescr, sysError := sd.SysDescr()
	sysName, sysNameError := sd.SysName()
	chassisId, chassisIdError := sd.ChassisId()

	response := &DiscoveryBasicResponse{
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
