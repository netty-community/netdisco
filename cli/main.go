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
package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/gosnmp/gosnmp"
	"github.com/netty-community/netdisco"
	"github.com/netty-community/netdisco/helpers/network"
	"github.com/netty-community/netdisco/logger"
	"github.com/netty-community/netdisco/model/snmp"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func main() {
	loggerConfig := logger.LogConfig{
		Formatter: "text",
	}
	logger.InitLogger(&loggerConfig)
	app := &cli.App{
		Name:        "netdisco",
		Usage:       "SNMP discovery tool",
		Version:     "1.0.0",
		HideVersion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "v",
				Aliases: []string{"version"},
				Value:   "2c",
				Usage:   "SNMP version (2c/v3)",
			},
			&cli.StringFlag{
				Name:    "c",
				Aliases: []string{"community"},
				Value:   "public",
				Usage:   "SNMP v2 community string, default public",
			},
			&cli.StringFlag{
				Name:    "u",
				Aliases: []string{"snmpv3-security-name"},
				Value:   "",
				Usage:   "set security name",
			},
			&cli.StringFlag{
				Name:    "l",
				Aliases: []string{"snmpv3-security-level"},
				Value:   "noAuthNoPriv",
				Usage:   "set security level (noAuthNoPriv|authNoPriv|authPriv)",
			},
			&cli.StringFlag{
				Name:    "a",
				Aliases: []string{"auth-protocol"},
				Value:   "sha",
				Usage:   "set authentication protocol (MD5|SHA)",
			},
			&cli.StringFlag{
				Name:    "x",
				Aliases: []string{"privacy-protocol"},
				Value:   "DES",
				Usage:   "set privacy protocol (DES|AES)",
			},
			&cli.StringFlag{
				Name:    "X",
				Aliases: []string{"privacy-protocol-pass-phrase"},
				Value:   "",
				Usage:   "set privacy protocol pass phrase",
			},
			&cli.IntFlag{
				Name:    "p",
				Aliases: []string{"port"},
				Value:   161,
				Usage:   "SNMP port, default 161",
			},
			&cli.IntFlag{
				Name:    "t",
				Aliases: []string{"timeout"},
				Value:   2,
				Usage:   "Timeout in seconds, default 2",
			},
			&cli.IntFlag{
				Name:    "r",
				Aliases: []string{"max-repetitions"},
				Value:   10,
				Usage:   "snmpwalk/snmpbulk Max repetitions, default 2",
			},
			&cli.StringFlag{
				Name:    "discovery",
				Aliases: []string{"d"},
				Value:   "basic",
				Usage:   "Discovery mode (basic, extend, ap, all)",
			},
			&cli.StringFlag{
				Name:    "network",
				Aliases: []string{"n"},
				Value:   "192.168.0.0/24",
				Usage:   "Network to scan",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   "output.json",
				Usage:   "Output file",
			},
		},
		Action: func(c *cli.Context) error {
			return runDiscovery(c)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Logger.Fatal("App run failed: %v", zap.Error(err))
	}
}

func runDiscovery(c *cli.Context) error {
	version := c.String("version")
	if version != "2c" {
		return fmt.Errorf("unsupported SNMP version: %s", version)
	}
	community := c.String("community")
	port := c.Uint("port")
	timeout := c.Int("timeout")
	maxRepetitions := c.Int("max-repetitions")
	discoveryMode := strings.ToLower(c.String("discovery"))
	cidr := c.String("network")
	outputFile := c.String("output")
	println(fmt.Sprintf("discovery mode :%s\nversion :%s\ncommunity :%s\nport :%d\ntimeout :%d\nmax-repetitions :%d\ndiscovery :%s\nnetwork :%s\noutput :%s\n", discoveryMode, version, community, port, timeout, maxRepetitions, discoveryMode, cidr, outputFile))
	targets, err := network.CIDR2IpStrings(cidr)
	if err != nil {
		logger.Logger.Panic(fmt.Sprintf("[ScanDeviceBasicInfo]: received wrong ip range %s", cidr), zap.Error(err))
	}
	snmpConfigs := make([]*snmp.SnmpConfig, 0)
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

	var output []byte
	switch discoveryMode {
	case "basic":
		result := netdisco.NetworkDeviceDiscovery(snmpConfigs)
		output, err = json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %v", err)
		}
	case "extend":
		result := netdisco.EnrichDevice(snmpConfigs)
		output, err = json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %v", err)
		}
	case "all":
		result := netdisco.EnrichDeviceFull(snmpConfigs)
		output, err = json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %v", err)
		}
	case "ap":
		result := netdisco.WlanApDiscovery(snmpConfigs)
		output, err = json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %v", err)
		}
	default:
		return fmt.Errorf("unsupported discovery mode: %s", discoveryMode)
	}

	err = os.WriteFile(outputFile, output, 0644)
	if err != nil {
		return fmt.Errorf("failed to write output file: %v", err)
	}
	fmt.Printf("Discovery completed. Results saved to %s\n", outputFile)
	return nil
}

func parseCIDR(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	// Exclude network address and broadcast address
	return ips[1 : len(ips)-1], nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
