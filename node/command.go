package node

// Command defines the interface for communication commands.
type Command interface {
	Execute() error
	Name() string
}
