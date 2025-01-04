package sysobjectid

import (
	"github.com/netty-community/netdisco/devicemodel"
	manufacturer "github.com/netty-community/netdisco/model/manufacturer"
	platform "github.com/netty-community/netdisco/model/platform"
)

func NetgearDeviceModel(sysObjId string) *devicemodel.DeviceModel {
	stringPlatform := string(platform.Netgear)
	oidMap := map[string]map[string]string{
		".1.3.6.1.4.1.4526.1.1":       {"platform": stringPlatform, "model": "FSM700S"},
		".1.3.6.1.4.1.4526.100.1.27":  {"platform": stringPlatform, "model": "M4300-8X8F"},
		".1.3.6.1.4.1.4526.100.1.28":  {"platform": stringPlatform, "model": "M4300-12X12F"},
		".1.3.6.1.4.1.4526.100.1.29":  {"platform": stringPlatform, "model": "M4300-24X24F"},
		".1.3.6.1.4.1.4526.100.10.16": {"platform": stringPlatform, "model": "S3300-28X"},
		".1.3.6.1.4.1.4526.100.10.18": {"platform": stringPlatform, "model": "S3300-52X"},
		".1.3.6.1.4.1.4526.100.11.5":  {"platform": stringPlatform, "model": "GSM7224V2"},
		".1.3.6.1.4.1.4526.100.2.1":   {"platform": stringPlatform, "model": "GSM7212"},
		".1.3.6.1.4.1.4526.100.2.3":   {"platform": stringPlatform, "model": "GSM7212"},
		".1.3.6.1.4.1.4526.100.4.1":   {"platform": stringPlatform, "model": "GS748T"},
		".1.3.6.1.4.1.4526.100.4.10":  {"platform": stringPlatform, "model": "GS724TP"},
		".1.3.6.1.4.1.4526.100.4.11":  {"platform": stringPlatform, "model": "GS748TP"},
		".1.3.6.1.4.1.4526.100.4.12":  {"platform": stringPlatform, "model": "GS724TR"},
		".1.3.6.1.4.1.4526.100.4.13":  {"platform": stringPlatform, "model": "GS748TR"},
		".1.3.6.1.4.1.4526.100.4.16":  {"platform": stringPlatform, "model": "GS716TV2"},
		".1.3.6.1.4.1.4526.100.4.17":  {"platform": stringPlatform, "model": "GS724TV3"},
		".1.3.6.1.4.1.4526.100.4.18":  {"platform": stringPlatform, "model": "GS108TV2"},
		".1.3.6.1.4.1.4526.100.4.19":  {"platform": stringPlatform, "model": "GS110TP"},
		".1.3.6.1.4.1.4526.100.4.2":   {"platform": stringPlatform, "model": "FS726T"},
		".1.3.6.1.4.1.4526.100.4.20":  {"platform": stringPlatform, "model": "FS728TPV2"},
		".1.3.6.1.4.1.4526.100.4.3":   {"platform": stringPlatform, "model": "GS716T"},
		".1.3.6.1.4.1.4526.100.4.30":  {"platform": stringPlatform, "model": "XS712T"},
		".1.3.6.1.4.1.4526.100.4.31":  {"platform": stringPlatform, "model": "GS716TV3"},
		".1.3.6.1.4.1.4526.100.4.32":  {"platform": stringPlatform, "model": "GS724TV4"},
		".1.3.6.1.4.1.4526.100.4.33":  {"platform": stringPlatform, "model": "GS748TV5"},
		".1.3.6.1.4.1.4526.100.4.38":  {"platform": stringPlatform, "model": "XS716T"},
		".1.3.6.1.4.1.4526.100.4.39":  {"platform": stringPlatform, "model": "XS708T"},
		".1.3.6.1.4.1.4526.100.4.4":   {"platform": stringPlatform, "model": "FS750T"},
		".1.3.6.1.4.1.4526.100.4.41":  {"platform": stringPlatform, "model": "GS418TPP"},
		".1.3.6.1.4.1.4526.100.4.42":  {"platform": stringPlatform, "model": "GS510TLP"},
		".1.3.6.1.4.1.4526.100.4.43":  {"platform": stringPlatform, "model": "GS510TPP"},
		".1.3.6.1.4.1.4526.100.4.47":  {"platform": stringPlatform, "model": "XS712TV2"},
		".1.3.6.1.4.1.4526.100.4.48":  {"platform": stringPlatform, "model": "GS310TP"},
		".1.3.6.1.4.1.4526.100.4.5":   {"platform": stringPlatform, "model": "GS724T"},
		".1.3.6.1.4.1.4526.100.4.6":   {"platform": stringPlatform, "model": "FS726TP"},
		".1.3.6.1.4.1.4526.100.4.7":   {"platform": stringPlatform, "model": "FS728TP"},
		".1.3.6.1.4.1.4526.100.4.8":   {"platform": stringPlatform, "model": "GS108T"},
		".1.3.6.1.4.1.4526.100.4.9":   {"platform": stringPlatform, "model": "GS108TP"},
		".1.3.6.1.4.1.4526.100.6.11":  {"platform": stringPlatform, "model": "SRX5308"},
	}

	data, ok := oidMap[sysObjId]
	if !ok {
		return &devicemodel.DeviceModel{
			Platform:     platform.Netgear,
			Manufacturer: manufacturer.Netgear,
			DeviceModel:  devicemodel.UnknownDeviceModel,
		}
	}

	return &devicemodel.DeviceModel{
		Platform:     platform.Platform(data["platform"]),
		Manufacturer: manufacturer.Netgear,
		DeviceModel:  data["model"],
	}

}
