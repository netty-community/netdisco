# Netdisco
[![License](https://img.shields.io/badge/License-BSD_3--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)
> It's a brand new version of [netty_snmp](https://github.com/netty-community/netty-snmp) and re-write with Go for easy running script.
> A lot of new items supported in Go version.

# Network Device information Discovery via SNMP

Netdisco is a flexible,powerful,high performance tool for network device information collection.
It will recognize network device's manufacturer and platform(netmiko driver) automatically through
SNMP `SysObjectId` and collect a lot of network basic information with friendly output format and
exceptions traceback.

Inspirations of SNMP and why not command line tools(CLI):

CLI output is un-structure data, we need use regex/textfsm liked tools to extract the information
from network device output. It's very painful to maintain the string-based regex code and very hard to
extend. SNMP is fast, stable to get and all data is structured which is easy to handle by code, also it's
much easier to maintain the codebase, some of items is generic for all network manufacturers, such
as interface, hostname, lldp neighbors. One time coding and support all vendors, fancy, right?
So, personally use snmp to collect basic information of network device is a better choice than CLI.

## Main Support Items
- **Hostname**: the hostname of network device
- **System Description**: the textual information of the system: include model name, software version, hardware version and etc.
- **Uptime**: network device uptime
- **ChassisID**: the mac address of chassis, the unique identifier of chassis
- **Interfaces**: interface full information
  - Interface Index
  - Interface Name
  - Interface Description
  - Interface Mtu
  - Interface Speed: physical speed
  - Interface High Speed: negotiation speed
  - Interface Type
  - Interface Mac Address
  - Interface Admin status
  - Interface Operational Status
  - Interface IP Addresses: (L3 Interface)
  - Interface Port Mode: access/trunk/hybrid ...
- **LLDP Neighbors**: lldp information will local and remote connection info
  - Local Chassis Id: the unique identifier mac address of chassis
  - Local Hostname: the hostname of local device
  - Local Interface Name
  - Local Interface Description
  - Remote Chassis Id: the unique identifier mac address of lldp neighbor chassis
  - Remote Hostname: the hostname of lldp neighbor
  - Remote Interface Name: the interface name of lldp neighbor
  - Remote interface Description: the interface description of lldp neighbor

- **Entities**
  - Entity Physical Description: the mode information of entity
  - Entity Physical Name: the name of entity
  - Entity Software revision: software version
  - Entity Serial Number: the serial number of entity
- **VLANS**
- **StackWise**
- **Prefixes**
- **ARP Table**
- **Mac Address Table**
- **AccessPoints**
---
## Main support manufactures
### well tested
- Cisco
- Huawei
- Aruba
### basic tested
- H3C
- Ruijie
- Arista
- Fortinet
- PaloAlto
- Juniper
