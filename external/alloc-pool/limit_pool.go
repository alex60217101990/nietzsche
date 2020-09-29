package pool

import (
	"fmt"

	"github.com/alex60217101990/nietzsche/external/consts"
	"github.com/alex60217101990/nietzsche/external/logger"
)

type LimitPool struct {
	pool chan []byte
}

func (lp *LimitPool) InitPool(poolCap uint16) Pool {
	if poolCap == 0 {
		poolCap = consts.PoolCap
	}
	return &LimitPool{
		pool: make(chan []byte, poolCap),
	}
}

func (lp *LimitPool) GetBytes() (b []byte) {
	select {
	case bt, ok := <-lp.pool:
		if ok {
			return bt
		}
		// non-normal behaivor, need fix!
		logger.AppLogger.Fatal(fmt.Errorf("closed the cahnnel of the LimitPool"))
	default:
	}
	return // nil<[]byte>
}

func (lp *LimitPool) PutBytes(b []byte) {
	if cap(b) > consts.PoolMaxCap {
		return
	}
	b = b[:0]
	select {
	case lp.pool <- b:
	default:
	}
	return
}
