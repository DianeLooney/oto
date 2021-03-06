// Copyright 2020 The Oto Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mux

import (
	"bytes"
	"sync"
)

// ConcurrentBuffer provides a bytes.Buffer that is safe for concurrent use.
type ConcurrentBuffer struct {
	buf bytes.Buffer
	m   sync.Mutex
	ch  chan struct{}
}

// Len returns the number of bytes currently available to be read from the
// buffer.
func (b *ConcurrentBuffer) Len() int {
	b.m.Lock()
	defer b.m.Unlock()

	return b.buf.Len()
}

func (b *ConcurrentBuffer) Read(buf []byte) (int, error) {
	b.m.Lock()
	defer b.m.Unlock()

	n, err := b.buf.Read(buf)
	if b.ch != nil && b.buf.Len() == 0 {
		b.ch <- struct{}{}
		close(b.ch)
		b.ch = nil
	}
	return n, err
}

func (b *ConcurrentBuffer) Write(buf []byte) (int, error) {
	b.m.Lock()
	n, err := b.buf.Write(buf)

	ch := make(chan struct{})
	b.ch = ch
	b.m.Unlock()

	<-ch

	return n, err
}
