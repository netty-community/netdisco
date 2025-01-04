package factory

const SysDescr = ".1.3.6.1.2.1.1.1.0"
const SysObjectID = ".1.3.6.1.2.1.1.2.0"
const SysUpTime = ".1.3.6.1.2.1.1.3.0"
const SysName = ".1.3.6.1.2.1.1.5.0"

// IF-MIB and related
const IfIndex = ".1.3.6.1.2.1.2.2.1.1"
const IfDescr = ".1.3.6.1.2.1.2.2.1.2"
const IfType = ".1.3.6.1.2.1.2.2.1.3"
const IfMtu = ".1.3.6.1.2.1.2.2.1.4"
const IfSpeed = ".1.3.6.1.2.1.2.2.1.5"
const IfPhysAddr = ".1.3.6.1.2.1.2.2.1.6"
const IfAdminStatus = ".1.3.6.1.2.1.2.2.1.7"
const IfOperStatus = ".1.3.6.1.2.1.2.2.1.8"
const IfLastChange = ".1.3.6.1.2.1.2.2.1.9"
const IfAlias = ".1.3.6.1.2.1.31.1.1.1.18"
const IfHighSpeed = ".1.3.6.1.2.1.31.1.1.1.15"
const IfAdEntIfIndex = ".1.3.6.1.2.1.4.20.1.2"
const IfAdEntNetMask = ".1.3.6.1.2.1.4.20.1.3"

// LLDP-MIB
const LldpLocChassisId = ".1.0.8802.1.1.2.1.3.2.0"
const LldpLocPortId = ".1.0.8802.1.1.2.1.3.7.1.3"
const LldpLocPortDesc = ".1.0.8802.1.1.2.1.3.7.1.4"
const LldpRemChassisIdSubtype = ".1.0.8802.1.1.2.1.4.1.1.4"
const LldpRemChassisId = ".1.0.8802.1.1.2.1.4.1.1.5"
const LldpRemPortIdSubtype = ".1.0.8802.1.1.2.1.4.1.1.6"
const LldpRemPortId = ".1.0.8802.1.1.2.1.4.1.1.7"
const LldpRemPortDesc = ".1.0.8802.1.1.2.1.4.1.1.8"
const LldpRemSysName = ".1.0.8802.1.1.2.1.4.1.1.9"

// ENTITY-MIB

const EntPhysicalDescr = ".1.3.6.1.2.1.47.1.1.1.1.2"

const EntPhysicalClass = ".1.3.6.1.2.1.47.1.1.1.1.5"

const EntPhysicalName = ".1.3.6.1.2.1.47.1.1.1.1.7"
const EntPhysicalHardwareRev = ".1.3.6.1.2.1.47.1.1.1.1.8"
const EntPhysicalFirmwareRev = ".1.3.6.1.2.1.47.1.1.1.1.9"
const EntPhysicalSoftwareRev = ".1.3.6.1.2.1.47.1.1.1.1.10"
const EntPhysicalSerialNum = ".1.3.6.1.2.1.47.1.1.1.1.11"

// -------- stack Cisco ------- #
const CswSwitchRole = ".1.3.6.1.4.1.9.9.500.1.2.1.1.3"
const CswSwitchHwPriority = ".1.3.6.1.4.1.9.9.500.1.2.1.1.5"
const CswSwitchState = ".1.3.6.1.4.1.9.9.500.1.2.1.1.6"
const CswSwitchMacAddress = ".1.3.6.1.4.1.9.9.500.1.2.1.1.7"

// --------- stack H3C --------
const H3cStackMemberID = ".1.3.6.1.4.1.2011.10.2.91.2.1.1"
const H3cStackConfigMemberID = ".1.3.6.1.4.1.2011.10.2.91.2.1.2"
const H3cStackPriority = ".1.3.6.1.4.1.2011.10.2.91.2.1.3"

// ---------- RuiJie stack ----------
const ScMemberMacAddress = ".1.3.6.1.4.1.4881.1.1.10.2.31.1.2.1.1.1"
const ScMemberNumber = ".1.3.6.1.4.1.4881.1.1.10.2.31.1.2.1.1.2"
const ScMemberOperStatus = ".1.3.6.1.4.1.4881.1.1.10.2.31.1.2.1.1.3"
const ScMemberDeviceID = ".1.3.6.1.4.1.4881.1.1.10.2.31.1.2.1.1.4"
const ScMemberRowStatus = ".1.3.6.1.4.1.4881.1.1.10.2.31.1.2.1.1.5"

// PaloAlto Basic Info
const PanSysSwVersion = ".1.3.6.1.4.1.25461.2.1.2.1.1"
const PanSysSerialNumber = ".1.3.6.1.4.1.25461.2.1.2.1.3"

const Dot1dTpFdbAddress = ".1.3.6.1.2.1.17.4.3.1.1"     // The MAC address of the FDB entry. mac_address
const Dot1dTpFdbPort = ".1.3.6.1.2.1.17.4.3.1.2"        // The port number of the FDB entry. int
const Dot1dBasePortIfIndex = ".1.3.6.1.2.1.17.1.4.1.2"  // The ifIndex of the port. int
const IpNetToMediaPhysAddress = ".1.3.6.1.2.1.4.22.1.2" // The MAC address of the port.
const IpNetToMediaType = ".1.3.6.1.2.1.4.22.1.4"        // 1: other(1), 2: invalid(2), 3: dynamic(3), 4: static(4)
// known issue with vrf scan, need add @param to community with suffix, but not support now
const IpAdEntAddr = ".1.3.6.1.2.1.4.20.1.1"
const IpAdEntIfIndex = ".1.3.6.1.2.1.4.20.1.2"
const IpAdEntNetMask = ".1.3.6.1.2.1.4.20.1.3"
