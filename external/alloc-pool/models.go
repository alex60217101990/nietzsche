package pool

import "bytes"

type WriteCloser struct {
	*bytes.Buffer
}

func (wc *WriteCloser) Close() error {
	// Noop
	return nil
}
