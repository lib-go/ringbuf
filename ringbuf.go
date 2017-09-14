package ringbuf

import (
	"sync"
)

type RingBuf struct {
	sync.Mutex
	buf    []byte
	index  int
	cycled bool
}

func New(maxSize int) *RingBuf {
	return &RingBuf{buf: make([]byte, maxSize)}
}

func (b *RingBuf) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	b.Lock()
	capacity := len(b.buf)
	i := b.index
	length := len(p)
	if length > capacity {
		b.index = 0
		copy(b.buf, p[length-capacity:])
		b.cycled = true
	} else if length > capacity-i {
		copy(b.buf[i:], p)
		copy(b.buf, p[capacity-i:])
		b.cycled = true
		b.index = (i + length) % capacity
	} else {
		copy(b.buf[i:], p)
		b.index = i + length
	}
	b.Unlock()
	return length, nil
}

func (b *RingBuf) Bytes() []byte {
	var result []byte
	b.Lock()
	if b.cycled {
		capacity := len(b.buf)
		index := b.index
		result = make([]byte, capacity)
		copy(result, b.buf[index:])
		copy(result[capacity-index:], b.buf[:index])
	} else {
		result = make([]byte, b.index)
		copy(result, b.buf)
	}
	b.Unlock()
	return result
}
