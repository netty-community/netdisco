package sysobjectid

import (
	"github.com/netty-community/netdisco/devicemodel"
	manufacturer "github.com/netty-community/netdisco/model/manufacturer"
	platform "github.com/netty-community/netdisco/model/platform"
)

func F5DeviceModel(sysObjId string) *devicemodel.DeviceModel {
	// stringPlatform := string(platform.F5)
	oidMap := map[string]map[string]string{}

	data, ok := oidMap[sysObjId]
	if !ok {
		return &devicemodel.DeviceModel{
			Platform:     platform.F5,
			Manufacturer: manufacturer.F5,
			DeviceModel:  devicemodel.UnknownDeviceModel,
		}
	}

	return &devicemodel.DeviceModel{
		Platform:     platform.Platform(data["platform"]),
		Manufacturer: manufacturer.F5,
		DeviceModel:  data["model"],
	}

}
