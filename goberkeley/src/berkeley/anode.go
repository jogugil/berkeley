package berkeley

import (
	"errors"
	"log"
<<<<<<< HEAD
	"time"

	zmq "github.com/pebbe/zmq4" // Librería para trabajar con ZeroMQ
)

// Interfaz Handler con el método HandleProcess
type Handler interface {
	HandleProcess(message string) (string, error)
}

=======
	"sync"
	"time"

	"github.com/pebbe/zmq4" // Librería para trabajar con ZeroMQ
)

>>>>>>> 4df5b3de8e63625e208a4e5d62139f2f2b1cb612
// AbstractNode proporciona la funcionalidad base para nodos en el sistema.
type AbstractNode struct {
	Name          string
	Address       string
	Timeout       time.Duration
	NodeAddresses map[string]string
<<<<<<< HEAD
	Context       *zmq.Context
	Socket        *zmq.Socket
	Logger        *log.Logger
	Handler       // Composición de la interfaz Handler
=======
	Context       *zmq4.Context
	Socket        *zmq4.Socket
	Logger        *log.Logger
	Mutex         sync.Mutex // Para proteger recursos compartidos
>>>>>>> 4df5b3de8e63625e208a4e5d62139f2f2b1cb612
}

// NewAbstractNode crea e inicializa un nuevo nodo base.
func NewAbstractNode(name, address string, timeout time.Duration) (*AbstractNode, error) {
<<<<<<< HEAD
	context, err := zmq.NewContext()
=======
	context, err := zmq4.NewContext()
>>>>>>> 4df5b3de8e63625e208a4e5d62139f2f2b1cb612
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
<<<<<<< HEAD
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
func (n *AbstractNode) SendMessageSync(address string, message string) (string, error) {
	// Creación del socket REQ (Request) para enviar el mensaje
	var socket *zmq.Socket
	socket, err := n.Context.NewSocket(zmq.REQ)
	if err != nil {
		n.Logger.Printf("Error al crear el socket REQ: %v", err) // Traza adicional
=======
func (n *AbstractNode) InitializeNodeWithAddresses(addresses map[string]string) {
	n.NodeAddresses = addresses
	n.Logger.Printf("Nodo %s inicializado con direcciones %v", n.Name, addresses)
}

// SendMessageSync envía un mensaje de forma síncrona y espera una respuesta.
func (n *AbstractNode) SendMessageSync(address, message string) (string, error) {
	n.Mutex.Lock()
	defer n.Mutex.Unlock()

	socket, err := n.Context.NewSocket(zmq4.REQ)
	if err != nil {
>>>>>>> 4df5b3de8e63625e208a4e5d62139f2f2b1cb612
		return "", errors.New("error al crear el socket REQ")
	}
	defer socket.Close()

<<<<<<< HEAD
	// Establecer el timeout para recibir respuestas
	socket.SetRcvtimeo(n.Timeout * time.Millisecond) // Convertir el timeout a milisegundos

	// Traza: Mostrar el valor del timeout configurado
	n.Logger.Printf("Timeout de recepción configurado a: %v milisegundos", int64(n.Timeout)) // Traza para el timeout

	// Conectar al socket en la dirección proporcionada
	err = socket.Connect("tcp://" + address)
	if err != nil {
		n.Logger.Printf("Error al conectar con %s: %v", address, err) // Traza adicional
		return "", errors.New("error al conectar con " + address)
	}

	n.Logger.Printf("Conectado a %s", address) // Traza para verificar la conexión

	// Enviar el mensaje al servidor
	_, err = socket.Send(message, 0)
	if err != nil {
		n.Logger.Printf("Error al enviar el mensaje a %s: %v", address, err) // Traza adicional
		return "", errors.New("error al enviar el mensaje")
	}

	n.Logger.Printf("Mensaje enviado a %s: %s", address, message) // Traza de envío

	// Intentar recibir la respuesta
	reply, err := socket.Recv(0)
	if err != nil {
		n.Logger.Printf("Error al recibir respuesta de %s: %v", address, err) // Traza para el error de recepción
		return "", errors.New("no se recibió respuesta del socket")
	}

	n.Logger.Printf("Respuesta recibida de %s: %s", address, reply) // Traza de respuesta recibida
=======
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
>>>>>>> 4df5b3de8e63625e208a4e5d62139f2f2b1cb612
	return reply, nil
}

// SendMessageAsync envía un mensaje de manera asíncrona.
func (n *AbstractNode) SendMessageAsync(address, message string) error {
<<<<<<< HEAD
	socket, err := n.Context.NewSocket(zmq.PUSH)
=======
	n.Mutex.Lock()
	defer n.Mutex.Unlock()

	socket, err := n.Context.NewSocket(zmq4.PUSH)
>>>>>>> 4df5b3de8e63625e208a4e5d62139f2f2b1cb612
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

<<<<<<< HEAD
// StartListening inicia el proceso de escucha para mensajes entrantes en el nodo.
// Configura un socket de tipo REP (Response) para recibir y responder a mensajes.
func (n *AbstractNode) StartListening() error {
	// Log para indicar que estamos intentando crear un socket REP.
	n.Logger.Printf("Intentando crear un socket REP para el nodo %s en la dirección %s", n.Name, n.Address)
	socket, err := n.Context.NewSocket(zmq.REP)
	if err != nil {
		// Si ocurre un error al crear el socket, loguea el error y retorna un mensaje de error.
		n.Logger.Printf("Error al crear el socket REP: %v", err)
		return errors.New("error al crear el socket REP")
	}
	n.Socket = socket
	n.Logger.Printf("Socket REP creado exitosamente para el nodo %s", n.Name)

	// Enlaza el socket a la dirección TCP proporcionada en el nodo para esperar conexiones.
	n.Logger.Printf("Enlazando el socket REP en la dirección tcp://%s", n.Address)
	err = socket.Bind("tcp://" + n.Address)
	if err != nil {
		// Si ocurre un error al enlazar el socket, loguea el error y retorna un mensaje de error.
		n.Logger.Printf("Error al enlazar el socket en %s: %v", n.Address, err)
		return errors.New("error al enlazar el socket en " + n.Address)
	}
	n.Logger.Printf("Socket enlazado exitosamente en %s", n.Address)

	// Inicia una nueva goroutine para escuchar de manera concurrente sin bloquear el hilo principal.
	n.Logger.Printf("Iniciando goroutine para escuchar en el nodo %s", n.Name)
	go func() {
		// Logea que el nodo ha comenzado a escuchar en la dirección configurada.
		n.Logger.Printf("Nodo %s escuchando en %s", n.Name, n.Address)

		// Entra en un bucle infinito para recibir y procesar mensajes.
		for {
			// Recibe un mensaje del socket.
			n.Logger.Printf("Esperando mensaje en %s...", n.Name)
			message, err := socket.Recv(0)
			if err != nil {
				// Si ocurre un error al recibir el mensaje, loguea el error y termina el bucle.
				n.Logger.Printf("Error al recibir mensaje: %v", err)
				break
			}
			n.Logger.Printf("Mensaje recibido en %s: %v", n.Name, message)

			// Llama al método HandleProcess, que es implementado por el tipo real de nodo (Follower, Leader, etc.).
			// Este método procesará el mensaje recibido y generará una respuesta.
			n.Logger.Printf("Procesando mensaje en %s...", n.Name)
			response, err := n.Handler.HandleProcess(message)
			if err != nil {
				// Si hay un error al procesar el mensaje, loguea el error y termina el bucle.
				n.Logger.Printf("Error en HandleProcess en %s: %v", n.Name, err)
				break
			}

			// Envía la respuesta de vuelta al cliente a través del socket.
			n.Logger.Printf("Enviando respuesta en %s: %v", n.Name, response)
			_, err = socket.Send(response, 0)
			if err != nil {
				// Si ocurre un error al enviar la respuesta, loguea el error y termina el bucle.
				n.Logger.Printf("Error al enviar respuesta en %s: %v", n.Name, err)
=======
// StartListening inicia el proceso de escucha para mensajes entrantes.
func (n *AbstractNode) StartListening() error {
	n.Mutex.Lock()
	defer n.Mutex.Unlock()

	socket, err := n.Context.NewSocket(zmq4.REP)
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
>>>>>>> 4df5b3de8e63625e208a4e5d62139f2f2b1cb612
				break
			}
		}
	}()
<<<<<<< HEAD

	// La función regresa nil si todo se configura correctamente y la goroutine se inicia sin errores.
	n.Logger.Printf("Escucha iniciada exitosamente en el nodo %s", n.Name)
=======
>>>>>>> 4df5b3de8e63625e208a4e5d62139f2f2b1cb612
	return nil
}

// Close cierra los recursos del nodo.
func (n *AbstractNode) Close() error {
<<<<<<< HEAD
=======
	n.Mutex.Lock()
	defer n.Mutex.Unlock()

>>>>>>> 4df5b3de8e63625e208a4e5d62139f2f2b1cb612
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
<<<<<<< HEAD
=======


>>>>>>> 4df5b3de8e63625e208a4e5d62139f2f2b1cb612
