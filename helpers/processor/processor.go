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

// processor string defines standardized ways for string/bytes processing.

package processor

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
)

func String2Md5(s string) string {
	hash := md5.New()
	_, _ = hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}

func ByteToMd5(b []byte) string {
	hash := md5.New()
	_, _ = hash.Write(b)
	return hex.EncodeToString(hash.Sum(nil))
}

// 将字符串转换为指针类型的字符串，当字符串值为""零值时，转换为nil
// background：基础的sdk统一定义了结构体基础的字段没有指针类型
func StringToPtrString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func PtrStringToString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func RandomHexString(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}
