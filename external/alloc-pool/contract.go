package pool

import "bytes"

type Pool interface {
	GetBytes() (b []byte)
	PutBytes(b []byte)
}

type BufferPool interface {
	InitPool() BufferPool
	GetBuffer() (b *bytes.Buffer)
	PutBuffer(b *bytes.Buffer)
}
