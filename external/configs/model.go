package configs

import (
	"encoding/json"

	"github.com/alex60217101990/nietzsche/external/logger"

	"github.com/pkg/errors"
)

var Conf *Configs

type Configs struct {
	Ver         *string       `yaml:"ver"`
	ClusterName string        `yaml:"service-name" json:"service_name"`
	IsDebug     bool          `yaml:"-" json:"-"`
	Logger      logger.Logger `yaml:"logger" json:"logger"`
	Server      *Server       `yaml:"http-server" json:"http_server"`
	DB          *DB           `yaml:"db"`
	Timeouts    *Timeouts     `yaml:"timeouts"`

	Raft  *Raft  `yaml:"consensus"`
	Store *Store `yaml:"store" json:"store"`
}

type Server struct {
	Host string `yaml:"server-host" json:"server_host"`
	Port uint16 `yaml:"server-port" json:"server_port"`
}

type Store struct {
	StoreType                StoreType `yaml:"store-type" json:"store_type"`
	DbName                   string    `yaml:"db-name" json:"db_name"`
	BucketName               string    `yaml:"bucket-name" json:"bucket_name"`
	UseStreamDataCompression bool      `yaml:"use-compression" json:"use_compression"`
}

func (s *Store) MarshalJSON() ([]byte, error) {
	type alias struct {
		StoreType  string `json:"store_type"`
		DbName     string `json:"db_name"`
		BucketName string `json:"bucket_name"`
	}
	if s == nil {
		s = &Store{}
	}
	return json.Marshal(alias{
		StoreType:  s.StoreType.String(),
		DbName:     s.DbName,
		BucketName: s.BucketName,
	})
}

func (s *Store) UnmarshalJSON(data []byte) (err error) {
	type alias struct {
		StoreType  string `json:"store_type"`
		DbName     string `json:"db_name"`
		BucketName string `json:"bucket_name"`
	}
	var tmp alias
	if err = json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	if s == nil {
		s = &Store{}
	}

	err = s.StoreType.Set(tmp.StoreType)
	if err != nil {
		return errors.WithMessagef(err, "failed to parse '%s'", tmp.StoreType)
	}

	s.BucketName = tmp.BucketName
	s.DbName = tmp.DbName

	return nil
}

func (s *Store) MarshalYAML() (interface{}, error) {
	type alias struct {
		StoreType  string `json:"store-type"`
		DbName     string `json:"db-name"`
		BucketName string `json:"bucket-name"`
	}
	if s == nil {
		s = &Store{}
	}
	return alias{
		StoreType:  s.StoreType.String(),
		BucketName: s.BucketName,
		DbName:     s.DbName,
	}, nil
}

func (s *Store) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type alias struct {
		StoreType  string `json:"store-type"`
		DbName     string `json:"db-name"`
		BucketName string `json:"bucket-name"`
	}
	var tmp alias
	if err := unmarshal(&tmp); err != nil {
		return err
	}

	if s == nil {
		s = &Store{}
	}

	err := s.StoreType.Set(tmp.StoreType)
	if err != nil {
		return errors.WithMessagef(err, "failed to parse '%s'", tmp.StoreType)
	}

	s.BucketName = tmp.BucketName
	s.DbName = tmp.DbName

	return nil
}

type DB struct {
	RepoType RepoType `yaml:"db-type" json:"db_type"`
	Host     string   `yaml:"db-host" json:"db_host"`
	UserName string   `yaml:"db-user" json:"db_user"`
	Password string   `yaml:"db-password" json:"db_password"`
	Port     uint16   `yaml:"db-port" json:"db_port"`
	DbName   string   `yaml:"db-name" json:"db_name"`
}

func (db *DB) MarshalJSON() ([]byte, error) {
	type alias struct {
		RepoType string `json:"db_type"`
		Host     string `json:"db_host"`
		UserName string `json:"db_user"`
		Password string `json:"db_password"`
		Port     uint16 `json:"db_port"`
		DbName   string `json:"db_name"`
	}
	if db == nil {
		db = &DB{}
	}
	return json.Marshal(alias{
		RepoType: db.RepoType.String(),
		Host:     db.Host,
		UserName: db.UserName,
		Password: db.Password,
		Port:     db.Port,
		DbName:   db.DbName,
	})
}

func (db *DB) UnmarshalJSON(data []byte) (err error) {
	type alias struct {
		RepoType string `json:"db_type"`
		Host     string `json:"db_host"`
		UserName string `json:"db_user"`
		Password string `json:"db_password"`
		Port     uint16 `json:"db_port"`
		DbName   string `json:"db_name"`
	}
	var tmp alias
	if err = json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	if db == nil {
		db = &DB{}
	}

	err = db.RepoType.Set(tmp.RepoType)
	if err != nil {
		return errors.WithMessagef(err, "failed to parse '%s'", tmp.RepoType)
	}

	db.Host = tmp.Host
	db.UserName = tmp.UserName
	db.Password = tmp.Password
	db.Port = tmp.Port
	db.DbName = tmp.DbName

	return nil
}

func (db *DB) MarshalYAML() (interface{}, error) {
	type alias struct {
		RepoType string `yaml:"db-type"`
		Host     string `yaml:"db-host"`
		UserName string `yaml:"db-user"`
		Password string `yaml:"db-password"`
		Port     uint16 `yaml:"db-port"`
		DbName   string `yaml:"db-name"`
	}
	if db == nil {
		db = &DB{}
	}
	return alias{
		RepoType: db.RepoType.String(),
		Host:     db.Host,
		UserName: db.UserName,
		Password: db.Password,
		Port:     db.Port,
		DbName:   db.DbName,
	}, nil
}

func (db *DB) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type alias struct {
		RepoType string `yaml:"db-type"`
		Host     string `yaml:"db-host"`
		UserName string `yaml:"db-user"`
		Password string `yaml:"db-password"`
		Port     uint16 `yaml:"db-port"`
		DbName   string `yaml:"db-name"`
	}
	var tmp alias
	if err := unmarshal(&tmp); err != nil {
		return err
	}

	if db == nil {
		db = &DB{}
	}

	err := db.RepoType.Set(tmp.RepoType)
	if err != nil {
		return errors.WithMessagef(err, "failed to parse '%s'", tmp.RepoType)
	}

	db.Host = tmp.Host
	db.UserName = tmp.UserName
	db.Password = tmp.Password
	db.Port = tmp.Port
	db.DbName = tmp.DbName

	return nil
}

type Timeouts struct {
	DefaultTimeout      uint8 `yaml:"default-timeout" json:"default_timeout"`
	DefaultStoreTimeout uint8 `yaml:"default-store-timeout" json:"default_store_timeout"`
}

type Raft struct {
	LogCacheSize   uint16            `yaml:"log-cache-size" json:"log_cache_size"`
	VolumeDir      string            `yaml:"volume-dir" json:"volume_dir"`
	NodeID         string            `yaml:"node-id"  json:"node_id"`
	Port           uint16            `yaml:"port"  json:"port"`
	MaxPool        uint16            `yaml:"max-pool"  json:"max_pool"`
	Transport      RaftTransportType `yaml:"transport-type"  json:"transport_type"`
	SnapShotRetain uint8             `yaml:"snap-shot-retain" json:"snap_shot_retain"`
}

func (r *Raft) MarshalJSON() ([]byte, error) {
	type alias struct {
		LogCacheSize   uint16 `json:"log_cache_size"`
		VolumeDir      string `json:"volume_dir"`
		NodeID         string `json:"node_id"`
		Port           uint16 `json:"port"`
		MaxPool        uint16 `json:"max_pool"`
		Transport      string `json:"transport_type"`
		SnapShotRetain uint8  `json:"snap_shot_retain"`
	}
	if r == nil {
		r = &Raft{}
	}
	return json.Marshal(alias{
		Transport:      r.Transport.String(),
		LogCacheSize:   r.LogCacheSize,
		VolumeDir:      r.VolumeDir,
		NodeID:         r.NodeID,
		Port:           r.Port,
		MaxPool:        r.MaxPool,
		SnapShotRetain: r.SnapShotRetain,
	})
}

func (r *Raft) UnmarshalJSON(data []byte) (err error) {
	type alias struct {
		LogCacheSize   uint16 `json:"log_cache_size"`
		VolumeDir      string `json:"volume_dir"`
		NodeID         string `json:"node_id"`
		Port           uint16 `json:"port"`
		MaxPool        uint16 `json:"max_pool"`
		Transport      string `json:"transport_type"`
		SnapShotRetain uint8  `json:"snap_shot_retain"`
	}
	var tmp alias
	if err = json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	if r == nil {
		r = &Raft{}
	}

	err = r.Transport.Set(tmp.Transport)
	if err != nil {
		return errors.WithMessagef(err, "failed to parse '%s'", tmp.Transport)
	}

	r.LogCacheSize = tmp.LogCacheSize
	r.VolumeDir = tmp.VolumeDir
	r.NodeID = tmp.NodeID
	r.Port = tmp.Port
	r.MaxPool = tmp.MaxPool
	r.SnapShotRetain = tmp.SnapShotRetain

	return nil
}

func (r *Raft) MarshalYAML() (interface{}, error) {
	type alias struct {
		LogCacheSize   uint16 `yaml:"log-cache-size"`
		VolumeDir      string `yaml:"volume-dir"`
		NodeID         string `yaml:"node-id"`
		Port           uint16 `yaml:"port"`
		MaxPool        uint16 `yaml:"max-pool"`
		Transport      string `yaml:"transport-type"`
		SnapShotRetain uint8  `yaml:"snap-shot-retain"`
	}
	if r == nil {
		r = &Raft{}
	}
	return alias{
		Transport:      r.Transport.String(),
		LogCacheSize:   r.LogCacheSize,
		VolumeDir:      r.VolumeDir,
		NodeID:         r.NodeID,
		Port:           r.Port,
		MaxPool:        r.MaxPool,
		SnapShotRetain: r.SnapShotRetain,
	}, nil
}

func (r *Raft) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type alias struct {
		LogCacheSize   uint16 `yaml:"log-cache-size"`
		VolumeDir      string `yaml:"volume-dir"`
		NodeID         string `yaml:"node-id"`
		Port           uint16 `yaml:"port"`
		MaxPool        uint16 `yaml:"max-pool"`
		Transport      string `yaml:"transport-type"`
		SnapShotRetain uint8  `yaml:"snap-shot-retain"`
	}
	var tmp alias
	if err := unmarshal(&tmp); err != nil {
		return err
	}

	if r == nil {
		r = &Raft{}
	}

	err := r.Transport.Set(tmp.Transport)
	if err != nil {
		return errors.WithMessagef(err, "failed to parse '%s'", tmp.Transport)
	}

	r.LogCacheSize = tmp.LogCacheSize
	r.VolumeDir = tmp.VolumeDir
	r.NodeID = tmp.NodeID
	r.Port = tmp.Port
	r.MaxPool = tmp.MaxPool
	r.SnapShotRetain = tmp.SnapShotRetain

	return nil
}
