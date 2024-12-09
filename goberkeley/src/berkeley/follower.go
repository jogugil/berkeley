package berkeley

import (
	"encoding/json"
<<<<<<< HEAD
	"errors"
	"fmt"
	"log"
	"time"
=======
	"fmt"
	"log"
	"time"

	"github.com/pebbe/zmq4"
>>>>>>> 4df5b3de8e63625e208a4e5d62139f2f2b1cb612
)

// Follower representa un nodo seguidor en el sistema distribuido.
type Follower struct {
<<<<<<< HEAD
	aAbstractNode *AbstractNode
=======
	AbstractNode
>>>>>>> 4df5b3de8e63625e208a4e5d62139f2f2b1cb612
	LeaderAddress string
}

// InitializeNode inicializa el nodo seguidor con su información específica.
<<<<<<< HEAD
func InitializeFollowerNode(name, address, leaderAddress string, timeout time.Duration) (*Follower, error) {
	// Asegúrate de llamar al método correcto de AbstractNode.

	f, err := NewFollower(name, address, leaderAddress, timeout)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// NewFollower crea e inicializa un nuevo nodo seguidor.
func NewFollower(name, address, leaderAddress string, timeout time.Duration) (*Follower, error) {
	// Inicializar el nodo abstracto usando NewAbstractNode
	abstractNode, err := NewAbstractNode(name, address, timeout)
	if err != nil {
		return nil, err
	}

	// Crear e inicializar la estructura Follower
	follower := &Follower{
		aAbstractNode: abstractNode,
		LeaderAddress: leaderAddress,
	}
	follower.aAbstractNode.Handler = follower
	return follower, nil
}

// HandleProcess maneja y procesa los mensajes recibidos del líder.
// Implementación de HandleProcess para Follower
func (f *Follower) HandleProcess(message string) (string, error) {
	log.Printf("Recibiendo mensaje del líder: %s", message) // Traza para ver el mensaje recibido

	// Deserialización del mensaje JSON
=======
func (f *Follower) InitializeNode(name, address, leaderAddress string, timeout int) error {
	f.AbstractNode.InitializeNode(name, address, timeout)
	f.LeaderAddress = leaderAddress
	return nil
}

// HandleProcess maneja y procesa los mensajes recibidos del líder.
func (f *Follower) HandleProcess(message string) (string, error) {
>>>>>>> 4df5b3de8e63625e208a4e5d62139f2f2b1cb612
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(message), &data); err != nil {
		log.Printf("Error al deserializar el mensaje JSON: %v", err)
		return `{"error":"Error al procesar el mensaje JSON"}`, nil
	}
<<<<<<< HEAD
	log.Printf("Mensaje deserializado correctamente: %+v", data) // Traza para mostrar el mensaje deserializado

	// Extraer los campos del mensaje deserializado
	leaderAddr := data["leader_address"].(string)
	operation := data["operation"].(string)
	leaderMessage := data["message"].(string)

	log.Printf("Procesando mensaje del líder: %s desde %s con operación: %s", leaderAddr, f.LeaderAddress, operation)

	// Procesar según la operación especificada
	switch operation {
	case "GET_TIME":
		// Acceder a 'time' y asegurar su tipo como float64
		T0Float, ok := data["time"].(float64)
		if !ok {
			log.Printf("Error al convertir data[\"time\"] a float64")
			return "", errors.New("Error al procesar el campo Time")
		}

		// Convertir de float64 a int64
		T0 := int64(T0Float)

		log.Printf("T0 recibido: %d", T0)
		currentTime := f.getCurrentTime()

		log.Printf("Operación GET_TIME: T0 recibido %d, Hora local calculada: %d", T0, currentTime) // Traza para tiempo

		// Llamada a la función que maneja el mensaje del líder y muestra su mensaje
		TP := f.displayLeaderMessage(leaderAddr, leaderMessage, T0, currentTime)
		log.Printf("Hora procesada TP: %d", TP) // Traza para TP

		// Responder con el formato esperado
		return fmt.Sprintf(`{"followerName":"%s", "localTime":"%d", "addressFollower":"%s"}`, f.aAbstractNode.Name, TP, f.aAbstractNode.Address), nil
	case "UPDATE_TIME":
		delta := int64(data["delta"].(float64))
		log.Printf("Operación UPDATE_TIME: Delta recibido: %d", delta) // Traza para delta

		// Modificar el sistema según el delta
		return f.modSystemTime(delta), nil
	case "CLOSE":
		log.Printf("Operación CLOSE: Cerrando seguidor %s", f.aAbstractNode.Name) // Traza para CLOSE
		return fmt.Sprintf(`{"followerName":"%s","operation":"CLOSE"}`, f.aAbstractNode.Name), nil
	default:
		log.Printf("Operación no reconocida en el mensaje del líder: %s", operation) // Traza para operación no reconocida
=======

	leaderName := data["leaderName"].(string)
	operation := data["operation"].(string)
	leaderMessage := data["message"].(string)

	log.Printf("Procesando mensaje del líder: %s desde %s con operación: %s", leaderName, f.LeaderAddress, operation)

	switch operation {
	case "GET_TIME":
		T0 := int64(data["T0"].(float64))
		currentTime := f.getCurrentTime()
		TP := f.displayLeaderMessage(leaderName, leaderMessage, T0, currentTime)
		return fmt.Sprintf(`{"followerName":"%s", "localTime":"%d", "addressFollower":"%s"}`, f.Name, TP, f.Address), nil
	case "MOD_TIME":
		delta := int64(data["DELTA"].(float64))
		return f.modSystemTime(delta), nil
	case "CLOSE":
		return fmt.Sprintf(`{"followerName":"%s","operation":"CLOSE"}`, f.Name), nil
	default:
		log.Printf("Operación no reconocida en el mensaje del líder: %s", operation)
>>>>>>> 4df5b3de8e63625e208a4e5d62139f2f2b1cb612
		return `{"error":"Operación no reconocida"}`, nil
	}
}

// displayLeaderMessage muestra el mensaje del líder y calcula un tiempo promedio (TP).
<<<<<<< HEAD
func (f *Follower) displayLeaderMessage(leaderAddr, leaderMessage string, T0, localTime int64) int64 {
	log.Printf("Mensaje del líder con dirección (%s) recibido: %s con T0 %s", leaderAddr, leaderMessage, time.UnixMilli(T0).String())
	TP := (T0 + localTime) / 2
	log.Printf("TP del seguidor %s es %s", f.aAbstractNode.Name, time.UnixMilli(TP).String())
=======
func (f *Follower) displayLeaderMessage(leaderName, leaderMessage string, T0, localTime int64) int64 {
	log.Printf("Mensaje del líder (%s) recibido: %s con T0 %s", leaderName, leaderMessage, time.UnixMilli(T0).String())
	TP := (T0 + localTime) / 2
	log.Printf("TP del seguidor %s es %s", f.Name, time.UnixMilli(TP).String())
>>>>>>> 4df5b3de8e63625e208a4e5d62139f2f2b1cb612
	return TP
}

// modSystemTime modifica el tiempo local del sistema basado en un delta.
func (f *Follower) modSystemTime(delta int64) string {
	currentLocalTime := time.Now().UnixMilli()
	modSystemTime := currentLocalTime + delta
<<<<<<< HEAD
	log.Printf("Tiempo del seguidor %s modificado de %s a %s", f.aAbstractNode.Name, time.UnixMilli(currentLocalTime).String(), time.UnixMilli(modSystemTime).String())
	return fmt.Sprintf(`{"followerName":"%s", "localTime":"%d","operation":"OK_MOD_TIME"}`, f.aAbstractNode.Name, modSystemTime)
=======
	log.Printf("Tiempo del seguidor %s modificado de %s a %s", f.Name, time.UnixMilli(currentLocalTime).String(), time.UnixMilli(modSystemTime).String())
	return fmt.Sprintf(`{"followerName":"%s", "localTime":"%d","operation":"OK_MOD_TIME"}`, f.Name, modSystemTime)
>>>>>>> 4df5b3de8e63625e208a4e5d62139f2f2b1cb612
}

// getCurrentTime obtiene la hora actual del sistema en milisegundos desde la época Unix.
func (f *Follower) getCurrentTime() int64 {
	TP := time.Now().UnixMilli()
	log.Printf("Fecha y hora local del seguidor: TP: %s", time.UnixMilli(TP).String())
	return TP
}

<<<<<<< HEAD
// StartAlgorithm configura e inicia el socket REP para escuchar mensajes entrantes
// y delega la responsabilidad de iniciar la escucha al nodo abstracto.
// StartAlgorithm configura e inicia la escucha en el seguidor
func (f *Follower) StartAlgorithm() error {
	// Llama a StartListening, que se encargará de la creación y enlace del socket.
	log.Printf("Iniciando algoritmo para el seguidor %s en %s", f.aAbstractNode.Name, f.aAbstractNode.Address)
	return f.aAbstractNode.StartListening()
=======
// StartAlgorithm configura y enlaza el socket REP y comienza la escucha de mensajes.
func (f *Follower) StartAlgorithm() error {
	socket, err := zmq4.NewSocket(zmq4.REP)
	if err != nil {
		log.Printf("Error al crear el socket REP: %v", err)
		return fmt.Errorf("error al inicializar el socket REP: %v", err)
	}
	f.Socket = socket

	err = f.Socket.Bind("tcp://" + f.Address)
	if err != nil {
		log.Printf("Error al enlazar el socket REP para el seguidor %s: %v", f.Name, err)
		return fmt.Errorf("error al enlazar el socket REP: %v", err)
	}

	log.Printf("Iniciando escucha para el seguidor %s en %s", f.Name, f.Address)
	f.StartListening()
	return nil
>>>>>>> 4df5b3de8e63625e208a4e5d62139f2f2b1cb612
}
