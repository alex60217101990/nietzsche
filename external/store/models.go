package store

// CommandPayload is payload sent by system when calling raft.Apply(cmd []byte, timeout time.Duration)
type CommandPayload struct {
	Operation string
	Key       string
	Value     interface{}
}

// ApplyResult response from Apply raft
type ApplyResult struct {
	Error error
	Data  interface{}
}
