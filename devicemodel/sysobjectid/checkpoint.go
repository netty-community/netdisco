package sysobjectid

import (
	"github.com/netty-community/netdisco/devicemodel"
	manufacturer "github.com/netty-community/netdisco/model/manufacturer"
	platform "github.com/netty-community/netdisco/model/platform"
)

func CheckPointDeviceModel(sysObjId string) *devicemodel.DeviceModel {
	// stringPlatform := string(platform.CheckPoint)
	oidMap := map[string]map[string]string{}

	data, ok := oidMap[sysObjId]
	if !ok {
		return &devicemodel.DeviceModel{
			Platform:     platform.CheckPoint,
			Manufacturer: manufacturer.CheckPoint,
			DeviceModel:  devicemodel.UnknownDeviceModel,
		}
	}

	return &devicemodel.DeviceModel{
		Platform:     platform.Platform(data["platform"]),
		Manufacturer: manufacturer.CheckPoint,
		DeviceModel:  data["model"],
	}

}
