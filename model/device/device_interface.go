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

package device

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/netty-community/netdisco/helpers/processor"
)

type DeviceInterface struct {
	IfIndex       uint64  `json:"ifIndex"`
	IfName        string  `json:"ifName"`
	IfDescr       string  `json:"ifDescr"`
	IfType        string  `json:"ifType"`
	IfMtu         uint64  `json:"ifMtu"`
	IfSpeed       uint64  `json:"ifSpeed"`
	IfPhysAddr    *string `json:"ifPhysAddr"`
	IfAdminStatus string  `json:"ifAdminStatus"`
	IfOperStatus  string  `json:"ifOperStatus"`
	IfLastChange  uint64  `json:"ifLastChange"`
	IfHighSpeed   uint64  `json:"ifHighSpeed"`
	IfIpAddress   *string `json:"ifIpAddress"`
	HashValue     string  `json:"hashValue"`
}

func (d *DeviceInterface) CalHashValue() string {

	hashString := fmt.Sprintf(
		"%s-%s-%s-%d-%d-%d-%s-%s-%s-%d-%s",
		d.IfName,
		d.IfDescr,
		d.IfType,
		d.IfHighSpeed,
		d.IfMtu,
		d.IfSpeed,
		processor.PtrStringToString(d.IfPhysAddr),
		d.IfAdminStatus,
		d.IfOperStatus,
		d.IfLastChange,
		processor.PtrStringToString(d.IfIpAddress))
	hash := md5.New()
	_, _ = hash.Write([]byte(hashString))
	return hex.EncodeToString(hash.Sum(nil))
}
