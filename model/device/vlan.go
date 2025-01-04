package device

type VlanItem struct {
	VlanId   uint32 `json:"vlanId"`
	VlanName string `json:"vlanName"`
	IfIndex  uint64 `json:"ifIndex"` // vlanIf index
	Network  string `json:"network"` // vlanIf Ip address to CIDR
	Gateway  string `json:"gateway"` //
}
