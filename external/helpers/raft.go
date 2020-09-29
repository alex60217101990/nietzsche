package helpers

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/alex60217101990/nietzsche/external/logger"

	"github.com/alex60217101990/nietzsche/external/configs"
	"github.com/alex60217101990/nietzsche/external/consts"
	rft "github.com/alex60217101990/nietzsche/external/raft-udp-transport"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

var (
	errInvalidTransport = errors.New("invalid raft transport type in config file")
)

func initRaftCacheStore(store *raftboltdb.BoltStore) (cacheStore *raft.LogCache, err error) {
	if configs.Conf.Raft.LogCacheSize == 0 {
		configs.Conf.Raft.LogCacheSize = consts.RaftLogCacheSize
	}

	// Wrap the store in a LogCache to improve performance.
	return raft.NewLogCache(int(configs.Conf.Raft.LogCacheSize), store)
}

func initRaftSnapshotStore() (snapshotStore *raft.FileSnapshotStore, err error) {
	if configs.Conf.Raft.SnapShotRetain == 0 {
		configs.Conf.Raft.SnapShotRetain = consts.RaftSnapShotRetain
	}

	return raft.NewFileSnapshotStore(configs.Conf.Raft.VolumeDir, int(configs.Conf.Raft.SnapShotRetain), os.Stdout)
}

func initRaftTransport() (transport *raft.NetworkTransport, err error) {
	raftBinAddr := fmt.Sprintf(":%d", configs.Conf.Raft.Port)
	switch configs.Conf.Raft.Transport {
	case configs.TCP:
		var tcpAddr *net.TCPAddr
		tcpAddr, err = net.ResolveTCPAddr("tcp", raftBinAddr)
		if err != nil {
			return transport, err
		}

		return raft.NewTCPTransport(raftBinAddr, tcpAddr, int(configs.Conf.Raft.MaxPool), TimeoutSecond(configs.Conf.Timeouts.DefaultTimeout), os.Stdout)
	case configs.UDP:
		var udpAddr *net.UDPAddr
		udpAddr, err := net.ResolveUDPAddr("udp", raftBinAddr)
		if err != nil {
			return transport, err
		}
		return rft.NewUDPTransport(raftBinAddr, udpAddr, int(configs.Conf.Raft.MaxPool), TimeoutSecond(configs.Conf.Timeouts.DefaultTimeout), os.Stdout)
	default:
		return transport, errInvalidTransport
	}
}

func InitRaftNode() {
	// Init default configs for raft cluster
	raftConf := raft.DefaultConfig()
	raftConf.LocalID = raft.ServerID(configs.Conf.Raft.NodeID)
	raftConf.SnapshotThreshold = 2 << 10

	// Init stable store
	store, err := raftboltdb.NewBoltStore(filepath.Join(configs.Conf.Raft.VolumeDir, consts.RaftPathPreffix))
	if err != nil {
		logger.AppLogger.Fatal(err)
	}

	// Init cache store
	var cacheStore *raft.LogCache
	cacheStore, err = initRaftCacheStore(store)
	if err != nil {
		logger.AppLogger.Fatal(err)
	}

	// Init snapshot store
	var snapshotStore *raft.FileSnapshotStore
	snapshotStore, err = initRaftSnapshotStore()
	if err != nil {
		logger.AppLogger.Fatal(err)
	}

	fmt.Println(snapshotStore, cacheStore)

	// Init transport
	var transport *raft.NetworkTransport
	transport, err = initRaftTransport()
	if err != nil {
		logger.AppLogger.Fatal(err)
	}

	raftServer, err := raft.NewRaft(raftConf /*fsmStore*/, nil, cacheStore, store, snapshotStore, transport)
	if err != nil {
		logger.AppLogger.Fatal(err)
	}

	// always start single server as a leader
	configuration := raft.Configuration{
		Servers: []raft.Server{
			{
				ID:      raft.ServerID(configs.Conf.Raft.NodeID),
				Address: transport.LocalAddr(),
			},
		},
	}

	raftServer.BootstrapCluster(configuration)
}
