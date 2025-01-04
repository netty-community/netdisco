package devicemodel

import (
	manufacturer "github.com/netty-community/netdisco/model/manufacturer"
	platform "github.com/netty-community/netdisco/model/platform"
)

var UnknownDeviceModel = "Unknown"

type DeviceModel struct {
	Platform     platform.Platform
	Manufacturer manufacturer.Manufacturer
	DeviceModel  string
}
