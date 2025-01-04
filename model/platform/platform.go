package platform

type Platform string

const (
	CiscoIos      Platform = "ciscoIos"
	CiscoIosXE    Platform = "ciscoXe"
	CiscoIosXR    Platform = "ciscoXr"
	CiscoNexusOS  Platform = "ciscoNxos"
	CiscoASA      Platform = "ciscoAsa"
	Huawei        Platform = "huaweiVrp"
	HuaweiCE      Platform = "huaweiCE"
	HuaweiFM      Platform = "huaweiFM"
	Aruba         Platform = "arubaOs"
	ArubaOSSwitch Platform = "arubaOsswitch"
	Arista        Platform = "arista"
	RuiJie        Platform = "ruijie"
	H3C           Platform = "h3c"
	FortiNet      Platform = "fortinet"
	PaloAlto      Platform = "Panos"
	Juniper       Platform = "juniper"
	Netgear       Platform = "netgear"
	TPLink        Platform = "tp_link"
	Ruckus        Platform = "fastiron"
	F5            Platform = "bigip"
	CheckPoint    Platform = "checkpoint"
	Extreme       Platform = "extreme"
	MikroTik      Platform = "mikrotik"
	Unknown       Platform = "unknown"
)

func SupportedPlatform() []Platform {
	return []Platform{
		CiscoIos,
		CiscoIosXE,
		CiscoIosXR,
		CiscoNexusOS,
		Huawei,
		HuaweiCE,
		HuaweiFM,
		Aruba,
		ArubaOSSwitch,
		Arista,
		RuiJie,
		H3C,
		FortiNet,
		PaloAlto,
		Juniper,
		Netgear,
		TPLink,
		Ruckus,
		F5,
		CheckPoint,
		Extreme,
		MikroTik,
		Unknown,
	}
}
