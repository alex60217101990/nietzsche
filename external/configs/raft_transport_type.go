package configs

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type RaftTransportType uint8

const (
	TCP RaftTransportType = iota
	UDP
)

var (
	_RaftTransportTypeNameToValue = map[string]RaftTransportType{
		"tcp": TCP,
		"udp": UDP,
	}

	_RaftTransportTypeValueToName = map[RaftTransportType]string{
		TCP: "tcp",
		UDP: "udp",
	}
)

func (rt RaftTransportType) MarshalYAML() (interface{}, error) {
	s, ok := _RaftTransportTypeValueToName[rt]
	if !ok {
		return nil, fmt.Errorf("invalid RaftTransportType: %d", rt)
	}
	return s, nil
}

func (rt *RaftTransportType) UnmarshalYAML(value *yaml.Node) error {
	v, ok := _RaftTransportTypeNameToValue[value.Value]
	if !ok {
		return fmt.Errorf("invalid RaftTransportType %q", value.Value)
	}
	*rt = v
	return nil
}

func (rt RaftTransportType) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(rt).(fmt.Stringer); ok {
		return json.Marshal(s.String())
	}
	s, ok := _RaftTransportTypeValueToName[rt]
	if !ok {
		return nil, fmt.Errorf("invalid RaftTransportType: %d", rt)
	}
	return json.Marshal(s)
}

func (rt *RaftTransportType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("RaftTransportType should be a string, got %s", data)
	}
	v, ok := _RaftTransportTypeNameToValue[s]
	if !ok {
		return fmt.Errorf("invalid RaftTransportType %q", s)
	}
	*rt = v
	return nil
}

func (rt RaftTransportType) Val() uint8 {
	return uint8(rt)
}

// it's for using with flag package
func (rt *RaftTransportType) Set(val string) error {
	if at, ok := _RaftTransportTypeNameToValue[val]; ok {
		*rt = at
		return nil
	}
	return fmt.Errorf("invalid repository type: %v", val)
}

func (rt RaftTransportType) String() string {
	return _RaftTransportTypeValueToName[rt]
}
