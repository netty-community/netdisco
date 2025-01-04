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
)

type LldpNeighbor struct {
	LocalChassisId  string `json:"localChassisId"`
	LocalHostname   string `json:"localHostname"`
	LocalIfName     string `json:"localIfName"`
	LocalIfDescr    string `json:"localIfDescr"`
	RemoteChassisId string `json:"remoteChassisId"`
	RemoteHostname  string `json:"remoteHostname"`
	RemoteIfName    string `json:"remoteIfName"`
	RemoteIfDescr   string `json:"remoteIfDescr"`
	HashValue       string `json:"hashValue"`
}

func (l *LldpNeighbor) CalHashValue() string {
	hashString := fmt.Sprintf(
		"%s-%s-%s-%s-%s-%s",
		l.LocalChassisId,
		l.LocalIfName,
		l.LocalIfDescr,
		l.RemoteChassisId,
		l.RemoteIfName,
		l.RemoteIfDescr)
	hash := md5.New()
	_, _ = hash.Write([]byte(hashString))
	return hex.EncodeToString(hash.Sum(nil))
}

func (l *LldpNeighbor) CalApHashValue() string {
	hashString := fmt.Sprintf(
		"%s-%s-%s-%s",
		l.LocalChassisId,
		l.LocalIfName,
		l.LocalIfDescr,
		l.RemoteChassisId)
	hash := md5.New()
	_, _ = hash.Write([]byte(hashString))
	return hex.EncodeToString(hash.Sum(nil))
}
