package ringbuf

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"crypto/rand"
)

func TestRingBuffer(t *testing.T) {
	buffer := New(10)
	buffer.Write([]byte("abcd"))

	assert.Equal(t, string(buffer.Bytes()), "abcd")
	assert.Equal(t, buffer.index, 4)

	buffer.Write([]byte("efghijklm"))
	assert.Equal(t, buffer.index, 3)
	assert.True(t, buffer.cycled)
	assert.Equal(t, buffer.Bytes(), []byte("defghijklm"))

	buffer.Write([]byte("01234567890"))
	assert.Equal(t, buffer.index, 0)
	assert.True(t, buffer.cycled)
	assert.Equal(t, buffer.Bytes(), []byte("1234567890"))

	bigBytes := make([]byte, 119)
	for i := range bigBytes {
		bigBytes[i] = byte(i)
	}
	buffer.Write(bigBytes)
	assert.Equal(t, buffer.index, 0)
	assert.True(t, buffer.cycled)
	assert.Equal(t, buffer.Bytes(), []byte{109, 110, 111, 112, 113, 114, 115, 116, 117, 118})
}

func BenchmarkRingBuffer_Write(b *testing.B) {
	buffer := New(1024)
	randomBuffer := make([]byte, 233)
	rand.Read(randomBuffer)
	for i := 0; i < b.N; i++ {
		buffer.Write(randomBuffer)
	}
}
