package oram

import (
	"crypto/rand"
	"io"
)

//RandomBytes : Helper function: Returns a byte slice of the specificed
// size filled with random data
func RandomBytes(bytes int) (data []byte) {
	data = make([]byte, bytes)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		panic(err)
	}
	return
}
