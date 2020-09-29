package pool

import (
	"bytes"
	"io"
	"sync"

	"github.com/alex60217101990/nietzsche/external/consts"
)

type UnlimitPoolBuffer struct {
	pool sync.Pool
}

func (up *UnlimitPoolBuffer) InitPool() BufferPool {
	return &UnlimitPoolBuffer{
		pool: sync.Pool{
			New: func() interface{} { return new(bytes.Buffer) },
		},
	}
}

func (up *UnlimitPoolBuffer) GetBuffer() (b *bytes.Buffer) {
	ifc := up.pool.Get()
	if ifc != nil {
		b = ifc.(*bytes.Buffer)
	}
	return
}

func (up *UnlimitPoolBuffer) PutBuffer(b *bytes.Buffer) {
	if b.Cap() <= consts.PoolMaxCap {
		b.Reset()
		up.pool.Put(b)
	}
}

func (up *UnlimitPoolBuffer) ConvertToWriteCloser(b *bytes.Buffer) io.WriteCloser {
	return &WriteCloser{b}
}

type UnlimitPool struct {
	pool sync.Pool
}

func (up *UnlimitPool) InitPool() Pool {
	return &UnlimitPool{
		pool: sync.Pool{
			New: func() interface{} { return []byte{} },
		},
	}
}

// get bytes
func (up *UnlimitPool) GetBytes() (b []byte) {
	ifc := up.pool.Get()
	if ifc != nil {
		b = ifc.([]byte)
	}
	return
}

// put bytes
func (up *UnlimitPool) PutBytes(b []byte) {
	if cap(b) <= consts.PoolMaxCap {
		b = b[:0] // reset
		up.pool.Put(b)
	}
}
