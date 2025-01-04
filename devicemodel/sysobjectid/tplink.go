package sysobjectid

import (
	"github.com/netty-community/netdisco/devicemodel"
	manufacturer "github.com/netty-community/netdisco/model/manufacturer"
	platform "github.com/netty-community/netdisco/model/platform"
)

func TPLinkDeviceModel(sysObjId string) *devicemodel.DeviceModel {
	stringPlatform := string(platform.TPLink)
	oidMap := map[string]map[string]string{
		".1.3.6.1.4.1.11863.5.1":    {"platform": stringPlatform, "model": "TL-SL5428"},
		".1.3.6.1.4.1.11863.5.2":    {"platform": stringPlatform, "model": "TL-SL5428"},
		".1.3.6.1.4.1.11863.5.3":    {"platform": stringPlatform, "model": "TL-SL3424"},
		".1.3.6.1.4.1.11863.5.4":    {"platform": stringPlatform, "model": "TL-SG3216"},
		".1.3.6.1.4.1.11863.5.5":    {"platform": stringPlatform, "model": "TL-SG3210"},
		".1.3.6.1.4.1.11863.5.6":    {"platform": stringPlatform, "model": "TL-SL3428"},
		".1.3.6.1.4.1.11863.5.7":    {"platform": stringPlatform, "model": "TL-SG5428"},
		".1.3.6.1.4.1.11863.5.8":    {"platform": stringPlatform, "model": "TL-SG3424P"},
		".1.3.6.1.4.1.11863.5.9":    {"platform": stringPlatform, "model": "TL-SG5412F"},
		".1.3.6.1.4.1.11863.5.10":   {"platform": stringPlatform, "model": "T2700-28TCT"},
		".1.3.6.1.4.1.11863.5.11":   {"platform": stringPlatform, "model": "TL-SL2428"},
		".1.3.6.1.4.1.11863.5.12":   {"platform": stringPlatform, "model": "TL-SG2216"},
		".1.3.6.1.4.1.11863.5.13":   {"platform": stringPlatform, "model": "TL-SG2424"},
		".1.3.6.1.4.1.11863.5.14":   {"platform": stringPlatform, "model": "TL-SG5428-CN"},
		".1.3.6.1.4.1.11863.5.15":   {"platform": stringPlatform, "model": "TL-SG2452"},
		".1.3.6.1.4.1.11863.5.16":   {"platform": stringPlatform, "model": "TL-SL2218"},
		".1.3.6.1.4.1.11863.5.17":   {"platform": stringPlatform, "model": "TL-SG2424P"},
		".1.3.6.1.4.1.11863.5.18":   {"platform": stringPlatform, "model": "TL-SG2210"},
		".1.3.6.1.4.1.11863.5.19":   {"platform": stringPlatform, "model": "TL-SL2210"},
		".1.3.6.1.4.1.11863.5.20":   {"platform": stringPlatform, "model": "T3700G-28TQ"},
		".1.3.6.1.4.1.11863.5.21":   {"platform": stringPlatform, "model": "TL-SL2226P"},
		".1.3.6.1.4.1.11863.5.22":   {"platform": stringPlatform, "model": "TL-SL2452"},
		".1.3.6.1.4.1.11863.5.23":   {"platform": stringPlatform, "model": "TL-SL2218P"},
		".1.3.6.1.4.1.11863.5.24":   {"platform": stringPlatform, "model": "TL-SG3424-IPV6"},
		".1.3.6.1.4.1.11863.5.25":   {"platform": stringPlatform, "model": "TL-SG2008"},
		".1.3.6.1.4.1.11863.5.26":   {"platform": stringPlatform, "model": "TL-SG2210P"},
		".1.3.6.1.4.1.11863.5.27":   {"platform": stringPlatform, "model": "T2700G-28TQ"},
		".1.3.6.1.4.1.11863.5.28":   {"platform": stringPlatform, "model": "T1600G-28TS"},
		".1.3.6.1.4.1.11863.5.29":   {"platform": stringPlatform, "model": "T1600G-52TS"},
		".1.3.6.1.4.1.11863.5.30":   {"platform": stringPlatform, "model": "T3700G-54TQ"},
		".1.3.6.1.4.1.11863.5.31":   {"platform": stringPlatform, "model": "T1700G-28TQ"},
		".1.3.6.1.4.1.11863.5.32":   {"platform": stringPlatform, "model": "T1700G-52TQ"},
		".1.3.6.1.4.1.11863.5.33":   {"platform": stringPlatform, "model": "T2600G-28TS"},
		".1.3.6.1.4.1.11863.5.34":   {"platform": stringPlatform, "model": "T2600G-52TS"},
		".1.3.6.1.4.1.11863.5.35":   {"platform": stringPlatform, "model": "T1500-28PCT"},
		".1.3.6.1.4.1.11863.5.37":   {"platform": stringPlatform, "model": "T1600G-28PS"},
		".1.3.6.1.4.1.11863.5.38":   {"platform": stringPlatform, "model": "T1600G-52PS"},
		".1.3.6.1.4.1.11863.5.39":   {"platform": stringPlatform, "model": "TL-SG2224P"},
		".1.3.6.1.4.1.11863.5.40":   {"platform": stringPlatform, "model": "TL-SG3428P"},
		".1.3.6.1.4.1.11863.5.41":   {"platform": stringPlatform, "model": "T1700X-16TS"},
		".1.3.6.1.4.1.11863.5.72":   {"platform": stringPlatform, "model": "T2600G-18TS"},
		".1.3.6.1.4.1.11863.5.122":  {"platform": stringPlatform, "model": "TL-SG3428"},
		".1.3.6.1.4.1.11863.5.186":  {"platform": stringPlatform, "model": "TL-SG3428MP"},
		".1.3.6.1.4.1.11863.3.2.10": {"platform": stringPlatform, "model": "EAP245"},
	}

	data, ok := oidMap[sysObjId]
	if !ok {
		return &devicemodel.DeviceModel{
			Platform:     platform.TPLink,
			Manufacturer: manufacturer.TPLink,
			DeviceModel:  devicemodel.UnknownDeviceModel,
		}
	}

	return &devicemodel.DeviceModel{
		Platform:     platform.Platform(data["platform"]),
		Manufacturer: manufacturer.TPLink,
		DeviceModel:  data["model"],
	}

}
