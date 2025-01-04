package sysobjectid

import (
	"github.com/netty-community/netdisco/devicemodel"
	manufacturer "github.com/netty-community/netdisco/model/manufacturer"
	platform "github.com/netty-community/netdisco/model/platform"
)

func PaloAltoDeviceModel(sysObjId string) *devicemodel.DeviceModel {
	stringPlatform := string(platform.PaloAlto)
	oidMap := map[string]map[string]string{
		".1.3.6.1.4.1.25461.2.3.1":  {"platform": stringPlatform, "model": "PA-4050"},
		".1.3.6.1.4.1.25461.2.3.11": {"platform": stringPlatform, "model": "PA-5020"},
		".1.3.6.1.4.1.25461.2.3.12": {"platform": stringPlatform, "model": "PA-200"},
		".1.3.6.1.4.1.25461.2.3.17": {"platform": stringPlatform, "model": "PA-3050"},
		".1.3.6.1.4.1.25461.2.3.18": {"platform": stringPlatform, "model": "PA-3020"},
		".1.3.6.1.4.1.25461.2.3.19": {"platform": stringPlatform, "model": "PA-3060"},
		".1.3.6.1.4.1.25461.2.3.2":  {"platform": stringPlatform, "model": "PA-4020"},
		".1.3.6.1.4.1.25461.2.3.21": {"platform": stringPlatform, "model": "PA-3250"},
		".1.3.6.1.4.1.25461.2.3.22": {"platform": stringPlatform, "model": "PA-5260"},
		".1.3.6.1.4.1.25461.2.3.23": {"platform": stringPlatform, "model": "PA-5250"},
		".1.3.6.1.4.1.25461.2.3.24": {"platform": stringPlatform, "model": "PA-5220"},
		".1.3.6.1.4.1.25461.2.3.29": {"platform": stringPlatform, "model": "PA-VM"},
		".1.3.6.1.4.1.25461.2.3.3":  {"platform": stringPlatform, "model": "PA-2050"},
		".1.3.6.1.4.1.25461.2.3.30": {"platform": stringPlatform, "model": "M-100"},
		".1.3.6.1.4.1.25461.2.3.31": {"platform": stringPlatform, "model": "PA-7050"},
		".1.3.6.1.4.1.25461.2.3.32": {"platform": stringPlatform, "model": "PALO ALTO GP-100"},
		".1.3.6.1.4.1.25461.2.3.33": {"platform": stringPlatform, "model": "WF-500"},
		".1.3.6.1.4.1.25461.2.3.34": {"platform": stringPlatform, "model": "PA-7080"},
		".1.3.6.1.4.1.25461.2.3.35": {"platform": stringPlatform, "model": "M-500"},
		".1.3.6.1.4.1.25461.2.3.36": {"platform": stringPlatform, "model": "PA-820"},
		".1.3.6.1.4.1.25461.2.3.37": {"platform": stringPlatform, "model": "PA-850"},
		".1.3.6.1.4.1.25461.2.3.38": {"platform": stringPlatform, "model": "PA-220"},
		".1.3.6.1.4.1.25461.2.3.39": {"platform": stringPlatform, "model": "M-600"},
		".1.3.6.1.4.1.25461.2.3.4":  {"platform": stringPlatform, "model": "PA-2020"},
		".1.3.6.1.4.1.25461.2.3.40": {"platform": stringPlatform, "model": "M-200"},
		".1.3.6.1.4.1.25461.2.3.41": {"platform": stringPlatform, "model": "PA-220R"},
		".1.3.6.1.4.1.25461.2.3.42": {"platform": stringPlatform, "model": "PA-5280"},
		".1.3.6.1.4.1.25461.2.3.43": {"platform": stringPlatform, "model": "PA-3220"},
		".1.3.6.1.4.1.25461.2.3.44": {"platform": stringPlatform, "model": "PA-3260"},
		".1.3.6.1.4.1.25461.2.3.5":  {"platform": stringPlatform, "model": "PA-4060"},
		".1.3.6.1.4.1.25461.2.3.6":  {"platform": stringPlatform, "model": "PA-500"},
		".1.3.6.1.4.1.25461.2.3.7":  {"platform": stringPlatform, "model": "PANORAMA"},
		".1.3.6.1.4.1.25461.2.3.8":  {"platform": stringPlatform, "model": "PA-5060"},
		".1.3.6.1.4.1.25461.2.3.9":  {"platform": stringPlatform, "model": "PA-5050"},
		".1.3.6.1.4.1.25461.2.3.69": {"platform": stringPlatform, "model": "PA-1400"},
		".1.3.6.1.4.1.25461.2.3.70": {"platform": stringPlatform, "model": "PA-1420"},
	}

	data, ok := oidMap[sysObjId]
	if !ok {
		return &devicemodel.DeviceModel{
			Platform:     platform.PaloAlto,
			Manufacturer: manufacturer.PaloAlto,
			DeviceModel:  devicemodel.UnknownDeviceModel,
		}
	}

	return &devicemodel.DeviceModel{
		Platform:     platform.Platform(data["platform"]),
		Manufacturer: manufacturer.PaloAlto,
		DeviceModel:  data["model"],
	}

}
