package berkeley

import (
	"time"
)

// INode define las operaciones básicas para un nodo en el sistema.
type INode interface {
	InitializeNode(name string, address string, timeout time.Duration) error
	InitializeNodeWithAddresses(name string, address string, timeout time.Duration, addresses map[string]string) error
	SendMessageSync(address string, message string) (string, error)
	SendMessageAsync(address string, message string) error
	StartListening() error
	StartAlgorithm() error
	Close() error
}
