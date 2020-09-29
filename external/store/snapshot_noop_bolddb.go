package store

import (
	"io"

	"golang.org/x/sync/errgroup"

	ap "github.com/alex60217101990/nietzsche/external/alloc-pool"
	"github.com/alex60217101990/nietzsche/external/configs"
	"github.com/alex60217101990/nietzsche/external/logger"

	"github.com/boltdb/bolt"
	"github.com/hashicorp/raft"
	"github.com/valyala/gozstd"
)

// snapshotNoop handle noop snapshot
type snapshotNoopBoltDB struct {
	db          *bolt.DB
	outStream   io.WriteCloser
	buffersPool ap.BufferPool
}

// Persist persist to disk. Return nil on success, otherwise return error.
func (s snapshotNoopBoltDB) Persist(sink raft.SnapshotSink) (err error) {
	pbuf := s.buffersPool.GetBuffer()
	defer func() {
		sink.Close()
		s.buffersPool.PutBuffer(pbuf)
	}()

	err = s.db.View(func(tx *bolt.Tx) (err error) {

		if configs.Conf.Store.UseStreamDataCompression {
			eg := new(errgroup.Group)

			eg.Go(func() (err error) {
				_, err = tx.WriteTo(pbuf)
				return err
			})

			eg.Go(func() (err error) {
				return gozstd.StreamCompressLevel(s.outStream, pbuf, 30)
			})

			err = eg.Wait()
		} else {
			_, err = tx.WriteTo(s.outStream)
		}

		return err
	})

	if err != nil {
		logger.AppLogger.Errorf(err.Error(),
			map[string]interface{}{
				"boltdb-shapshot-noop": "persist",
			})

		sink.Cancel()
	}

	return err
}

// Release release the lock after persist snapshot.
// Release is invoked when we are finished with the snapshot.
func (s snapshotNoopBoltDB) Release() {}

// newSnapshotNoop is returned by an FSM in response to a snapshotNoop
// It must be safe to invoke FSMSnapshot methods with concurrent
// calls to Apply.
func newSnapshotNoopBoltDB(db *bolt.DB, out io.WriteCloser, bp ap.BufferPool) (raft.FSMSnapshot, error) {
	return &snapshotNoopBoltDB{
		db:          db,
		outStream:   out,
		buffersPool: bp,
	}, nil
}
