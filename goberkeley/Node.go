package main

import (
	"fmt"
	"log"
	"time"

	"github.com/pebbe/zmq4"
)

type SocketZeroMQException struct {
	message string
	errType ErrorType
}

type ErrorType int

const (
	ReceiveError ErrorType = iota
	SendError
	ConnectionError
)

func (e *SocketZeroMQException) Error() string {
	return fmt.Sprintf("ErrorType: %v, Message: %s", e.errType, e.message)
}

type INode interface {
	InitializeNode(name, address string, timeout int)
	InitializeNodeWithAddresses(name, address string, timeout int, nodeAddresses map[string]string)
	SendMessageSync(address, message string) (string, error)
	SendMessageAsync(address, message string) error
	StartListening() error
	StartAlgorithm() error
	Close() error
}

type AbstractNode struct {
	context        *zmq4.Context
	socket         *zmq4.Socket
	name           string
	address        string
	timeout        int
	nodeAddresses  map[string]string
	listSocketsREQ []*zmq4.Socket
}

func (n *AbstractNode) InitializeNode(name, address string, timeout int) {
	n.name = name
	n.address = address
	n.timeout = timeout
	n.context, _ = zmq4.NewContext() // Initialize ZeroMQ context
	log.Printf("Node %s initialized at address %s", name, address)
}

func (n *AbstractNode) InitializeNodeWithAddresses(name, address string, timeout int, nodeAddresses map[string]string) {
	n.name = name
	n.address = address
	n.timeout = timeout
	n.context, _ = zmq4.NewContext()
	n.nodeAddresses = nodeAddresses
	log.Printf("Node %s initialized with node addresses %v at %s", name, nodeAddresses, address)
}

func (n *AbstractNode) SendMessageSync(address, message string) (string, error) {
	clientSocket, _ := n.context.NewSocket(zmq4.REQ)
	defer clientSocket.Close()

	err := clientSocket.Connect(address)
	if err != nil {
		return "", &SocketZeroMQException{"Failed to connect", ConnectionError}
	}

	clientSocket.Send(message, 0)
	log.Printf("Message sent: %s", message)

	reply, err := clientSocket.Recv(0)
	if err != nil {
		return "", &SocketZeroMQException{"Failed to receive reply", ReceiveError}
	}

	log.Printf("Received reply: %s", reply)
	return reply, nil
}

func (n *AbstractNode) SendMessageAsync(address, message string) error {
	clientSocket, _ := n.context.NewSocket(zmq4.PUSH)
	defer clientSocket.Close()

	err := clientSocket.Connect(address)
	if err != nil {
		return &SocketZeroMQException{"Failed to connect", ConnectionError}
	}

	clientSocket.Send(message, 0)
	log.Printf("Message sent asynchronously: %s", message)
	return nil
}

func (n *AbstractNode) StartListening() error {
	if n.socket == nil {
		return &SocketZeroMQException{"Socket not initialized", ConnectionError}
	}

	go func() {
		log.Printf("Listening started for node %s...", n.name)

		for {
			message, err := n.socket.Recv(0)
			if err != nil {
				log.Printf("Error receiving message: %v", err)
				break
			}

			log.Printf("Received message: %s", message)

			// Process the message and generate a response
			response := handleProcess(message)

			// Send response back to sender
			n.socket.Send(response, 0)
			log.Printf("Sent response: %s", response)

			// If the message contains "CLOSE", shut down the socket
			if response == `"operation":"CLOSE"` {
				log.Printf("Closing node %s.", n.name)
				n.Close()
				break
			}
		}
	}()
	return nil
}

func handleProcess(message string) string {
	// Handle the message processing here
	return "Processed " + message
}

func (n *AbstractNode) StartAlgorithm() error {
	// Start specific algorithm logic here
	log.Println("Starting algorithm...")
	time.Sleep(2 * time.Second) // Simulating algorithm work
	log.Println("Algorithm completed")
	return nil
}

func (n *AbstractNode) Close() error {
	if n.socket != nil {
		n.socket.Close()
	}
	if n.context != nil {
		n.context.Terminate()
	}
	log.Printf("Node %s resources closed", n.name)
	return nil
}

func main() {
	// Example usage
	node := &AbstractNode{}
	node.InitializeNode("Node1", "localhost:5555", 1000)
	node.StartListening()
	_, err := node.SendMessageSync("localhost:5555", "Hello, Node2")
	if err != nil {
		log.Printf("Error: %v", err)
	}
}