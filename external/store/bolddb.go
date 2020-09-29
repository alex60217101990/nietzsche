package store

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	ap "github.com/alex60217101990/nietzsche/external/alloc-pool"
	"github.com/alex60217101990/nietzsche/external/configs"
	"github.com/alex60217101990/nietzsche/external/helpers"
	"github.com/alex60217101990/nietzsche/external/logger"

	"github.com/boltdb/bolt"
	"github.com/hashicorp/raft"
	"github.com/valyala/gozstd"
)

type BoldDBStore struct {
	db          *bolt.DB
	pool        ap.Pool
	buffersPool ap.BufferPool
}

func NewBoldDBStore() Store {
	// Open the [some name].db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open(fmt.Sprintf("%s.db", configs.Conf.Store.DbName), 0600,
		&bolt.Options{Timeout: helpers.TimeoutSecond(
			configs.Conf.Timeouts.DefaultStoreTimeout,
		)})
	if err != nil {
		logger.AppLogger.Fatal(err)
	}

	return &BoldDBStore{
		db:          db,
		pool:        new(ap.UnlimitPool).InitPool(),
		buffersPool: new(ap.UnlimitPoolBuffer).InitPool(),
	}
}

func (b *BoldDBStore) SetDumpWriter(w io.WriteCloser) {

}

func (b BoldDBStore) Close() error {
	return b.db.Close()
}

// get fetch data from boldDB
func (b *BoldDBStore) get(key string) (data interface{}, err error) {
	pbuf := b.buffersPool.GetBuffer()
	defer func() {
		b.buffersPool.PutBuffer(pbuf)

		if err != nil || (err == nil && data == nil) {
			data = map[string]interface{}{}
		}
	}()

	err = b.db.View(func(tx *bolt.Tx) (err error) {
		bucket := tx.Bucket([]byte(configs.Conf.Store.BucketName))

		if configs.Conf.Store.UseStreamDataCompression {
			err = gozstd.StreamDecompress(pbuf, bytes.NewReader(bucket.Get([]byte(key))))
		} else {
			_, err = pbuf.Read(bucket.Get([]byte(key)))
		}

		return err
	})

	if err != nil {
		return data, err
	}

	if pbuf.Len() > 0 {
		err = gob.NewDecoder(pbuf).Decode(&data)
	}

	return data, err
}

// set store data to boldDB
func (b *BoldDBStore) set(key string, value interface{}) (err error) {
	pbuf := b.buffersPool.GetBuffer()
	bbuf := b.pool.GetBytes()
	defer func() {
		b.pool.PutBytes(bbuf)
		b.buffersPool.PutBuffer(pbuf)
	}()

	err = gob.NewEncoder(pbuf).Encode(value)
	if err != nil {
		return err
	}

	if configs.Conf.Store.UseStreamDataCompression {
		bbuf = gozstd.CompressLevel(bbuf[:0], pbuf.Bytes(), 30)
	} else {
		bbuf = pbuf.Bytes()
	}

	return b.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(configs.Conf.Store.BucketName)).Put([]byte(key), bbuf)
	})
}

// delete remove data from badgerDB
func (b *BoldDBStore) delete(key string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(configs.Conf.Store.BucketName)).Delete([]byte(key))
	})
}

// Apply log is invoked once a log entry is committed.
// It returns a value which will be made available in the
// ApplyFuture returned by Raft.Apply method if that
// method was called on the same Raft node as the FSM.
func (b *BoldDBStore) Apply(log *raft.Log) interface{} {
	switch log.Type {
	case raft.LogCommand:
		var payload = CommandPayload{}
		if err := json.Unmarshal(log.Data, &payload); err != nil {
			logger.AppLogger.Errorf(
				fmt.Sprintf("error marshalling store payload %s\n", err.Error()),
				map[string]interface{}{
					"raft": "apply",
				})
			return &ApplyResult{
				Error: err,
				Data:  nil,
			}
		}

		op := strings.ToUpper(strings.TrimSpace(payload.Operation))
		switch op {
		case "SET":
			return &ApplyResult{
				Error: b.set(payload.Key, payload.Value),
				Data:  payload.Value,
			}
		case "GET":
			data, err := b.get(payload.Key)
			return &ApplyResult{
				Error: err,
				Data:  data,
			}

		case "DELETE":
			return &ApplyResult{
				Error: b.delete(payload.Key),
				Data:  nil,
			}
		}
	}

	if configs.Conf.IsDebug {
		logger.AppLogger.Warnf("not raft log command type",
			map[string]interface{}{
				"raft": "apply",
			})
	}

	return nil
}

// Snapshot will be called during make snapshot.
// Snapshot is used to support log compaction.
// No need to call snapshot since it already persisted in disk (using BadgerDB) when raft calling Apply function.
func (b *BoldDBStore) Snapshot() (raft.FSMSnapshot, error) {
	return newSnapshotNoopBoltDB(b.db)
}

