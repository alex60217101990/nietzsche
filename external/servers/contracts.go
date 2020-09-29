package servers

type Node interface {
	Start() error
	Close() error
}
