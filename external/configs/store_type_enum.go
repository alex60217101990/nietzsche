package configs

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type StoreType uint8

const (
	StoreBoldDB StoreType = iota
	StoreBadgerDB
)

var (
	_StoreTypeNameToValue = map[string]StoreType{
		"boldbd":   StoreBoldDB,
		"BoldDB":   StoreBoldDB,
		"badgerdb": StoreBadgerDB,
		"BadgerDB": StoreBadgerDB,
	}

	_StoreTypeValueToName = map[StoreType]string{
		StoreBoldDB:   "boldbd",
		StoreBadgerDB: "badgerdb",
	}
)

func (st StoreType) MarshalYAML() (interface{}, error) {
	s, ok := _StoreTypeValueToName[st]
	if !ok {
		return nil, fmt.Errorf("invalid StoreType: %d", st)
	}
	return s, nil
}

func (st *StoreType) UnmarshalYAML(value *yaml.Node) error {
	v, ok := _StoreTypeNameToValue[value.Value]
	if !ok {
		return fmt.Errorf("invalid StoreType %q", value.Value)
	}
	*st = v
	return nil
}

func (st StoreType) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(st).(fmt.Stringer); ok {
		return json.Marshal(s.String())
	}
	s, ok := _StoreTypeValueToName[st]
	if !ok {
		return nil, fmt.Errorf("invalid StoreType: %d", st)
	}
	return json.Marshal(s)
}

func (st *StoreType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("StoreType should be a string, got %s", data)
	}
	v, ok := _StoreTypeNameToValue[s]
	if !ok {
		return fmt.Errorf("invalid StoreType %q", s)
	}
	*st = v
	return nil
}

func (st StoreType) Val() uint8 {
	return uint8(st)
}

// it's for using with flag package
func (st *StoreType) Set(val string) error {
	if at, ok := _StoreTypeNameToValue[val]; ok {
		*st = at
		return nil
	}
	return fmt.Errorf("invalid repository type: %v", val)
}

func (st StoreType) String() string {
	return _StoreTypeValueToName[st]
}
