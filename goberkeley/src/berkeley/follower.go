package berkeley

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/pebbe/zmq4"
)

// Follower representa un nodo seguidor en el sistema distribuido.
type Follower struct {
	aAbstractNode *AbstractNode
	LeaderAddress string
}

// InitializeNode inicializa el nodo seguidor con su información específica.
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

	return follower, nil
}

// HandleProcess maneja y procesa los mensajes recibidos del líder.
func (f *Follower) HandleProcess(message string) (string, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(message), &data); err != nil {
		log.Printf("Error al deserializar el mensaje JSON: %v", err)
		return `{"error":"Error al procesar el mensaje JSON"}`, nil
	}

	leaderName := data["leaderName"].(string)
	operation := data["operation"].(string)
	leaderMessage := data["message"].(string)

	log.Printf("Procesando mensaje del líder: %s desde %s con operación: %s", leaderName, f.LeaderAddress, operation)

	switch operation {
	case "GET_TIME":
		T0 := int64(data["T0"].(float64))
		currentTime := f.getCurrentTime()
		TP := f.displayLeaderMessage(leaderName, leaderMessage, T0, currentTime)
		return fmt.Sprintf(`{"followerName":"%s", "localTime":"%d", "addressFollower":"%s"}`, f.aAbstractNode.Name, TP, f.aAbstractNode.Address), nil
	case "MOD_TIME":
		delta := int64(data["DELTA"].(float64))
		return f.modSystemTime(delta), nil
	case "CLOSE":
		return fmt.Sprintf(`{"followerName":"%s","operation":"CLOSE"}`, f.aAbstractNode.Name), nil
	default:
		log.Printf("Operación no reconocida en el mensaje del líder: %s", operation)
		return `{"error":"Operación no reconocida"}`, nil
	}
}

// displayLeaderMessage muestra el mensaje del líder y calcula un tiempo promedio (TP).
func (f *Follower) displayLeaderMessage(leaderName, leaderMessage string, T0, localTime int64) int64 {
	log.Printf("Mensaje del líder (%s) recibido: %s con T0 %s", leaderName, leaderMessage, time.UnixMilli(T0).String())
	TP := (T0 + localTime) / 2
	log.Printf("TP del seguidor %s es %s", f.aAbstractNode.Name, time.UnixMilli(TP).String())
	return TP
}

// modSystemTime modifica el tiempo local del sistema basado en un delta.
func (f *Follower) modSystemTime(delta int64) string {
	currentLocalTime := time.Now().UnixMilli()
	modSystemTime := currentLocalTime + delta
	log.Printf("Tiempo del seguidor %s modificado de %s a %s", f.aAbstractNode.Name, time.UnixMilli(currentLocalTime).String(), time.UnixMilli(modSystemTime).String())
	return fmt.Sprintf(`{"followerName":"%s", "localTime":"%d","operation":"OK_MOD_TIME"}`, f.aAbstractNode.Name, modSystemTime)
}

// getCurrentTime obtiene la hora actual del sistema en milisegundos desde la época Unix.
func (f *Follower) getCurrentTime() int64 {
	TP := time.Now().UnixMilli()
	log.Printf("Fecha y hora local del seguidor: TP: %s", time.UnixMilli(TP).String())
	return TP
}

// StartAlgorithm configura y enlaza el socket REP y comienza la escucha de mensajes.
func (f *Follower) StartAlgorithm() error {
	socket, err := zmq4.NewSocket(zmq4.REP)
	if err != nil {
		log.Printf("Error al crear el socket REP: %v", err)
		return fmt.Errorf("error al inicializar el socket REP: %v", err)
	}
	f.aAbstractNode.Socket = socket

	err = f.aAbstractNode.Socket.Bind("tcp://" + f.aAbstractNode.Address)
	if err != nil {
		log.Printf("Error al enlazar el socket REP para el seguidor %s: %v", f.aAbstractNode.Name, err)
		return fmt.Errorf("error al enlazar el socket REP: %v", err)
	}

	log.Printf("Iniciando escucha para el seguidor %s en %s", f.aAbstractNode.Name, f.aAbstractNode.Address)
	f.aAbstractNode.StartListening()
	return nil
}
