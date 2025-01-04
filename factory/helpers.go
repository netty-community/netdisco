package factory

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/gosnmp/gosnmp"
	"github.com/netty-community/netdisco/model/device"
	"github.com/samber/lo"
)

var vlanPattern = regexp.MustCompile(`(vl|Vl)[A-Za-z]*(\d+)`)

func ExtractInteger(oid string, values []gosnmp.SnmpPDU) map[string]uint64 {
	result := make(map[string]uint64, len(values))
	for _, v := range values {
		snmpIndex := getSnmpIndex(oid, &v)
		result[snmpIndex] = gosnmp.ToBigInt(v.Value).Uint64()
	}
	return result
}

func ExtractIntegerWithShift(oid string, location int, values []gosnmp.SnmpPDU) map[string]uint64 {
	result := make(map[string]uint64, len(values))
	for _, v := range values {
		snmpIndex := shiftSnmpIndex(getSnmpIndex(oid, &v), location)
		result[snmpIndex] = gosnmp.ToBigInt(v.Value).Uint64()
	}
	return result
}

func extractIfIndex(oid string, values []gosnmp.SnmpPDU) map[string]string {
	result := make(map[string]string, len(values))
	for _, v := range values {
		snmpIndex := getSnmpIndex(oid, &v)
		value := "." + strconv.Itoa(int(gosnmp.ToBigInt(v.Value).Uint64()))
		result[value] = snmpIndex
	}
	return result
}

func ExtractString(oid string, values []gosnmp.SnmpPDU) map[string]string {
	result := make(map[string]string, len(values))
	for _, v := range values {
		snmpIndex := getSnmpIndex(oid, &v)
		result[snmpIndex] = fmt.Sprintf("%s", v.Value)
	}
	return result
}

func ExtractStringWithShift(oid string, location int, values []gosnmp.SnmpPDU) map[string]string {
	result := make(map[string]string, len(values))
	for _, v := range values {
		snmpIndex := shiftSnmpIndex(getSnmpIndex(oid, &v), location)
		result[snmpIndex] = fmt.Sprintf("%s", v.Value)
	}
	return result
}

func ExtractMacAddress(oid string, values []gosnmp.SnmpPDU) map[string]string {
	result := make(map[string]string, len(values))
	for _, v := range values {
		snmpIndex := getSnmpIndex(oid, &v)
		result[snmpIndex] = hex2mac(v.Value.([]byte))
	}
	return result
}

func ExtractMacAddressWithShift(oid string, location int, values []gosnmp.SnmpPDU) map[string]string {
	result := make(map[string]string, len(values))
	for _, v := range values {
		snmpIndex := shiftSnmpIndex(getSnmpIndex(oid, &v), location)
		result[snmpIndex] = hex2mac(v.Value.([]byte))
	}
	return result
}

func getSnmpIndex(oid string, values *gosnmp.SnmpPDU) string {
	return strings.Split(values.Name, oid)[1]
}

// shift snmp index is used for special case for lldp neighbor
// lldpRemChassisID: snmp_index value the last of 2 is local portID's index(last of 1)
func shiftSnmpIndex(index string, location int) string {
	splitValues := strings.Split(index, ".")
	if location < 0 {
		return splitValues[len(splitValues)+location]
	}
	return splitValues[location]
}

func hex2mac(hex []byte) string {
	if len(hex) != 6 {
		return ""
	}
	hexParts := make([]string, 6)
	for i, b := range hex {
		hexParts[i] = fmt.Sprintf("%02x", b)
	}
	return strings.Join(hexParts, ":")
}

func StringToHexMac(snmpIndex string) string {
	parts := strings.Split(strings.TrimPrefix(snmpIndex, "."), ".")
	if len(parts) != 6 {
		return ""
	}
	var macAddress strings.Builder
	for i, part := range parts {
		decimal, err := strconv.Atoi(part)
		if err != nil {
			return ""
		}
		hex := decimalToHex(decimal)
		macAddress.WriteString(hex)
		if i < 5 {
			macAddress.WriteString(":")
		}
	}
	return macAddress.String()
}

func decimalToHex(decimal int) string {
	hex := strconv.FormatInt(int64(decimal), 16)
	if len(hex) == 1 {
		hex = "0" + hex
	}
	return hex
}

func buildOidWithIndex(oid string, index []string) []string {
	results := make([]string, len(index))
	for i, v := range index {
		results[i] = fmt.Sprintf("%s.%s", oid, v)
	}
	return results
}

func netmaskToLength(netmask string) int {
	if strings.Contains(netmask, ".") {
		stringMask := net.IPMask(net.ParseIP(netmask).To4())
		length, _ := stringMask.Size()
		return length
	}
	stringMask := net.IPMask(net.ParseIP(netmask).To16())
	length, _ := stringMask.Size()
	return length
}

// IP-MIB: ipNetToMediaPhysAddress.102.10.90.200.120 = INTEGER: 102
// last 5 is ifIndex
// last 4 is address
func getIfIndexAndAddress(snmpIndex string) (ifIndex uint64, address string) {
	splitData := strings.Split(snmpIndex, ".")
	xLast4 := splitData[len(splitData)-4:]
	xLast5 := splitData[len(splitData)-5]
	address = strings.Join(xLast4, ".")
	_ifIndex, _ := strconv.Atoi(xLast5)
	ifIndex = uint64(_ifIndex)
	return ifIndex, address
}

func extractVlanId(ifName string) (uint32, bool) {
	if match := vlanPattern.FindStringSubmatch(ifName); match != nil {
		if len(match) == 3 {
			if num, err := strconv.Atoi(match[2]); err == nil {
				return uint32(num), true
			}
			return 0, false
		}
	}
	return 0, false
}

func ipAddrToNet(addr *string) (string, error) {
	if addr == nil || *addr == "" {
		return "", nil
	}
	ip, ipNet, err := net.ParseCIDR(*addr)
	if err != nil {
		return "", err
	}
	networkIp := ip.Mask(ipNet.Mask)
	size, _ := ipNet.Mask.Size()
	networkStr := fmt.Sprintf("%s/%d", networkIp.String(), size)
	return networkStr, nil
}

func getIfVlanIpRange(interfaces []*device.DeviceInterface) map[uint64]*VlanIpRange {
	var results = make(map[uint64]*VlanIpRange)
	for _, iface := range interfaces {
		if iface.IfType == "propVirtual" {
			vlanIpRange := &VlanIpRange{}
			vlan, ok := extractVlanId(iface.IfName)
			if ok {
				vlanIpRange.VlanId = vlan
			}
			inet, err := ipAddrToNet(iface.IfIpAddress)
			if err == nil {
				vlanIpRange.Range = inet
				vlanIpRange.Gateway = *iface.IfIpAddress
			}
			if vlan != 0 && inet != "" {
				results[iface.IfIndex] = vlanIpRange
			}
		} else if iface.IfType == "ethernetCsmacd" {
			vlanIpRange := &VlanIpRange{}
			inet, err := ipAddrToNet(iface.IfIpAddress)
			if err == nil {
				vlanIpRange.Range = inet
				vlanIpRange.Gateway = *iface.IfIpAddress
			}
			if inet != "" {
				results[iface.IfIndex] = vlanIpRange
			}
		}
	}
	return results
}

func EnrichArpInfo(arp []*device.ArpItem, interfaces []*device.DeviceInterface) []*device.ArpItem {
	for len(arp) <= 0 || len(interfaces) <= 0 {
		return arp
	}
	ifVlanIpRange := getIfVlanIpRange(interfaces)
	if len(ifVlanIpRange) <= 0 {
		return arp
	}
	for _, item := range arp {
		inner := ifVlanIpRange[item.IfIndex]
		item.VlanId = inner.VlanId
		item.Range = inner.Range
	}
	return arp
}

func EnrichVlanInfo(vlan []*device.VlanItem, interfaces []*device.DeviceInterface) []*device.VlanItem {
	if len(vlan) <= 0 || len(interfaces) <= 0 {
		return vlan
	}
	ifVlanIpRange := getIfVlanIpRange(interfaces)
	for _, vl := range vlan {
		ipRange, ok := ifVlanIpRange[vl.IfIndex]
		if !ok {
			continue
		}
		vl.Network = ipRange.Range
		vl.Gateway = ipRange.Gateway
	}
	return vlan
}

func lldpNeighborInterfaces(lldp []*device.LldpNeighbor) []string {
	results := make([]string, len(lldp))
	for _, item := range lldp {
		results = append(results, item.LocalIfName)
	}
	return results
}

func RemoveNonLocalMacAddress(
	mac *map[uint64][]string,
	interfaces []*device.DeviceInterface,
	lldp []*device.LldpNeighbor,
) *map[uint64][]string {
	interfaceNames := make(map[uint64]string)
	interfaceTypes := make(map[uint64]string)
	for _, iface := range interfaces {
		interfaceNames[iface.IfIndex] = iface.IfName
		interfaceTypes[iface.IfIndex] = iface.IfType
	}
	lldpInterfaces := lldpNeighborInterfaces(lldp)
	for index := range *mac {
		if lo.Contains(lldpInterfaces, interfaceNames[index]) || interfaceTypes[index] != "ethernetCsmacd" {
			delete(*mac, index)
		}
	}
	return mac
}

func EnrichMacAddress(
	macAddresses *map[uint64][]string,
	interfaces []*device.DeviceInterface,
	lldpNeighbors []*device.LldpNeighbor,
	arpItems []*device.ArpItem,
) []*device.MacAddressItem {
	var results []*device.MacAddressItem
	interfaceMapping := make(map[uint64]*device.DeviceInterface, len(interfaces))
	arpMapping := make(map[string]*device.ArpItem, len(arpItems))
	for _, iface := range interfaces {
		interfaceMapping[iface.IfIndex] = iface
	}
	for _, arpItem := range arpItems {
		arpMapping[arpItem.MacAddress] = arpItem
	}
	lldpInterfaces := lldpNeighborInterfaces(lldpNeighbors)
	for index, macAddresses := range *macAddresses {
		if len(macAddresses) == 0 {
			continue
		}
		iface := interfaceMapping[index]
		if !shouldIncludeMacAddress(iface, lldpInterfaces) {
			continue
		}
		for _, macAddress := range macAddresses {
			arpItem, ok := arpMapping[macAddress]
			if !ok {
				continue
			}
			results = append(results, &device.MacAddressItem{
				MacAddress: macAddress,
				IfIndex:    index,
				IfName:     iface.IfName,
				IfDescr:    iface.IfDescr,
				IpAddress:  arpItem.IpAddress,
				VlanId:     arpItem.VlanId,
			})
		}
	}
	return results
}

func shouldIncludeMacAddress(iface *device.DeviceInterface, lldpInterfaces []string) bool {
	return iface.IfType == "ethernetCsmacd" && !lo.Contains(lldpInterfaces, iface.IfName)
}

// return mac, ip and macIndex
func SnmpIndexToMacAndIp(snmpIndex string) (string, string, string) {
	parts := strings.Split(snmpIndex, ".")
	maxIndex := "." + strings.Join(parts[1:7], ".")
	macParts := make([]string, 6)
	for i := 0; i < 6; i++ {
		macParts[i] = fmt.Sprintf("%02x", parseInt(parts[i]))
	}
	mac := strings.Join(macParts, ":")
	ip := strings.Join(parts[6:], ".")
	return mac, ip, maxIndex

}
func parseInt(s string) byte {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return byte(n)
}

func ChannelToRadioType(channel uint16) string {
	channel24G := []uint16{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
	channel5G := []uint16{36, 40, 44, 48, 52, 56, 60, 64, 100, 104,
		108, 112, 116, 132, 136, 140, 149, 153, 157, 161, 165}
	if lo.Contains(channel24G, channel) {
		return "2.4GHz"
	} else if lo.Contains(channel5G, channel) {
		return "5GHz"
	} else if channel == 0 {
		return "Unknown"
	} else {
		return "6GHz"
	}
}
