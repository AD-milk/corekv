// Copyright 2021 logicrec Project Authors
//
// Licensed under the Apache License, Version 2.0 (the "License")
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Copyright 2021 hardcore-os Project Authors
//
// Licensed under the Apache License, Version 2.0 (the "License")
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package codec

import (
	"encoding/binary"
	"time"
)

type Entry struct {
	Key       []byte
	Value     []byte
	ExpiresAt uint64
}

// EncodedSize 计算不包含key的部分的size大小
func (e *Entry) EncodedSize() uint32 {
	sz := len(e.Value)
	enc := sizeVarint(e.ExpiresAt)
	return uint32(sz + enc)
}

func NewEntry(key, value []byte) *Entry {
	return &Entry{
		Key:   key,
		Value: value,
	}
}

// DecodeEntry _
func (e *Entry) DecodeEntry(buf []byte) {
	var sz int
	e.ExpiresAt, sz = binary.Uvarint(buf)
	e.Value = buf[sz:]
}

func (e *Entry) WithTTL(dur time.Duration) *Entry {
	e.ExpiresAt = uint64(time.Now().Add(dur).Unix())
	return e
}

// Size 计算整体的大小
func (e *Entry) Size() int64 {
	return int64(len(e.Key) + len(e.Value))
}

// EncodeEntry 对整个entry的value部分解码
func (e *Entry) EncodeEntry(b []byte) uint32 {
	sz := binary.PutUvarint(b[:], e.ExpiresAt)
	n := copy(b[sz:], e.Value)
	return uint32(sz + n)
}

// 计算变长编码长度
func sizeVarint(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
