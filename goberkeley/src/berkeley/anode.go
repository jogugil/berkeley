package berkeley

import (
	"errors"
	"log"
	"sync"
	"time"

	zmq "github.com/pebbe/zmq4" // Librería para trabajar con ZeroMQ
)

// AbstractNode proporciona la funcionalidad base para nodos en el sistema.
type AbstractNode struct {
	Name          string
	Address       string
	Timeout       time.Duration
	NodeAddresses map[string]string
	Context       *zmq.Context
	Socket        *zmq.Socket
	Logger        *log.Logger
	Mutex         sync.Mutex // Para proteger recursos compartidos
}

// NewAbstractNode crea e inicializa un nuevo nodo base.
func NewAbstractNode(name, address string, timeout time.Duration) (*AbstractNode, error) {
	context, err := zmq.NewContext()
	if err != nil {
		return nil, errors.New("error al crear el contexto ZeroMQ")
	}

	return &AbstractNode{
		Name:    name,
		Address: address,
		Timeout: timeout,
		Context: context,
		Logger:  log.Default(),
	}, nil
}

// InitializeNodeWithAddresses inicializa un nodo con direcciones de otros nodos.
func InitializeNodeWithAddresses(name, address string, timeout time.Duration, addresses map[string]string) (*AbstractNode, error) {
	var n *AbstractNode
	n, err := NewAbstractNode(name, address, timeout)
	if err != nil {
		return nil, err
	}
	n.NodeAddresses = addresses
	n.Logger.Printf("Nodo %s inicializado con direcciones %v", n.Name, addresses)
	return n, nil
}

// SendMessageSync envía un mensaje de forma síncrona y espera una respuesta.
func (n *AbstractNode) SendMessageSync(address, message string) (string, error) {
	n.Mutex.Lock()
	defer n.Mutex.Unlock()
	var socket *zmq.Socket
	socket, err := n.Context.NewSocket(zmq.REQ)
	if err != nil {
		return "", errors.New("error al crear el socket REQ")
	}
	defer socket.Close()

	socket.SetRcvtimeo(n.Timeout)
	err = socket.Connect("tcp://" + address)
	if err != nil {
		return "", errors.New("error al conectar con " + address)
	}

	_, err = socket.Send(message, 0)
	if err != nil {
		return "", errors.New("error al enviar el mensaje")
	}

	reply, err := socket.Recv(0)
	if err != nil {
		return "", errors.New("no se recibió respuesta del socket")
	}

	n.Logger.Printf("Mensaje enviado a %s: %s", address, message)
	n.Logger.Printf("Respuesta recibida: %s", reply)
	return reply, nil
}

// SendMessageAsync envía un mensaje de manera asíncrona.
func (n *AbstractNode) SendMessageAsync(address, message string) error {
	n.Mutex.Lock()
	defer n.Mutex.Unlock()

	socket, err := n.Context.NewSocket(zmq.PUSH)
	if err != nil {
		return errors.New("error al crear el socket PUSH")
	}
	defer socket.Close()

	err = socket.Connect("tcp://" + address)
	if err != nil {
		return errors.New("error al conectar con " + address)
	}

	_, err = socket.Send(message, 0)
	if err != nil {
		return errors.New("error al enviar el mensaje")
	}

	n.Logger.Printf("Mensaje asincrónico enviado a %s: %s", address, message)
	return nil
}

// StartListening inicia el proceso de escucha para mensajes entrantes.
func (n *AbstractNode) StartListening() error {
	n.Mutex.Lock()
	defer n.Mutex.Unlock()

	socket, err := n.Context.NewSocket(zmq.REP)
	if err != nil {
		return errors.New("error al crear el socket REP")
	}
	n.Socket = socket

	err = socket.Bind("tcp://" + n.Address)
	if err != nil {
		return errors.New("error al enlazar el socket en " + n.Address)
	}

	go func() {
		n.Logger.Printf("Nodo %s escuchando en %s", n.Name, n.Address)
		for {
			message, err := socket.Recv(0)
			if err != nil {
				n.Logger.Printf("Error al recibir mensaje: %v", err)
				break
			}
			response := n.handleProcess(message)
			_, err = socket.Send(response, 0)
			if err != nil {
				n.Logger.Printf("Error al enviar respuesta: %v", err)
				break
			}
		}
	}()
	return nil
}

// Close cierra los recursos del nodo.
func (n *AbstractNode) Close() error {
	n.Mutex.Lock()
	defer n.Mutex.Unlock()

	if n.Socket != nil {
		n.Socket.Close()
		n.Logger.Printf("Socket cerrado para el nodo %s", n.Name)
	}
	if n.Context != nil {
		n.Context.Term()
		n.Logger.Printf("Contexto cerrado para el nodo %s", n.Name)
	}
	return nil
}

// handleProcess debe ser sobrescrito por las subclases para procesar mensajes.
func (n *AbstractNode) handleProcess(message string) string {
	n.Logger.Printf("Procesando mensaje: %s", message)
	return `{"status":"unhandled"}`
}
