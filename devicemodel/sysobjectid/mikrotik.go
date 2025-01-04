package sysobjectid

import (
	"github.com/netty-community/netdisco/devicemodel"
	manufacturer "github.com/netty-community/netdisco/model/manufacturer"
	platform "github.com/netty-community/netdisco/model/platform"
)

func MikroTikDeviceModel(sysObjId string) *devicemodel.DeviceModel {
	stringPlatform := string(platform.MikroTik)
	oidMap := map[string]map[string]string{
		"1.3.6.1.4.1.14988.1": {"platform": stringPlatform, "model": "RB1200"},
	}

	data, ok := oidMap[sysObjId]
	if !ok {
		return &devicemodel.DeviceModel{
			Platform:     platform.MikroTik,
			Manufacturer: manufacturer.MikroTik,
			DeviceModel:  devicemodel.UnknownDeviceModel,
		}
	}

	return &devicemodel.DeviceModel{
		Platform:     platform.Platform(data["platform"]),
		Manufacturer: manufacturer.MikroTik,
		DeviceModel:  data["model"],
	}

}
