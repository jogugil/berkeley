package berkeley

import (
	"log"
	"sync"
	"time"
)

// Leader representa el nodo líder en el sistema Berkeley.
type Leader struct {
	aAbstractNode          *AbstractNode
	UnreachableFollowers   sync.Map // Seguidores inalcanzables
	SuccessfulFollowers    sync.Map // Seguidores que respondieron con éxito
	NonRespondingFollowers sync.Map // Seguidores que no respondieron a tiempo
	TimeUpdatedFollowers   sync.Map // Seguidores que actualizaron su tiempo correctamente
	FailedFollowers        sync.Map // Seguidores que no pudieron actualizar su tiempo
	Logger                 *log.Logger
}

// InitializeLeaderNode crea e inicializa un nuevo nodo líder.
func InitializeLeaderNode(name, address string, timeout time.Duration, nodeAddresses map[string]string) (*Leader, error) {
	// Suponiendo que InitializeNodeWithAddresses crea un nodo base y devuelve un puntero a AbstractNode
	baseNode, err := InitializeNodeWithAddresses(name, address, timeout, nodeAddresses)
	if err != nil {
		return nil, err
	}

	// Inicializamos el líder, asignando el nodo base y un logger
	leader := &Leader{
		aAbstractNode: baseNode, // Asignamos el puntero a AbstractNode
		Logger:        log.Default(),
	}

	return leader, nil
}

// StartAlgorithm implementa el algoritmo de sincronización Berkeley para el líder.
func (l *Leader) StartAlgorithm() {
	l.Logger.Println("Iniciando algoritmo de sincronización Berkeley...")

	// Simula el envío de solicitudes de tiempo a los seguidores
	for followerName, followerAddress := range l.aAbstractNode.NodeAddresses {
		go l.RequestTimeFromFollower(followerName, followerAddress)
	}

	// FALTA EL RESTOP--123-- vamos poco a poco
}

// SendCloseMessage envía un mensaje de cierre a un seguidor.
func (l *Leader) SendCloseMessage(followerAddress string) {
	message := `{"operation": "CLOSE", "message": "Cerrar conexión", "leaderName": "` + l.aAbstractNode.Name + `"}`

	response, err := l.aAbstractNode.SendMessageSync(followerAddress, message)
	if err != nil {
		l.Logger.Printf("Error al enviar mensaje de cierre a %s: %v", followerAddress, err)
		return
	}

	l.Logger.Printf("Respuesta de cierre del seguidor %s: %s", followerAddress, response)
}

// SendCloseMessagesToFollowers envía mensajes de cierre a todos los seguidores que respondieron con éxito.
func (l *Leader) SendCloseMessagesToFollowers() {
	var wg sync.WaitGroup

	l.SuccessfulFollowers.Range(func(key, value any) bool {
		wg.Add(1)
		go func(address string) {
			defer wg.Done()
			l.SendCloseMessage(address)
		}(key.(string))
		return true
	})

	wg.Wait()
	l.Logger.Println("Todos los mensajes de cierre han sido enviados.")
}

// RequestTimeFromFollower solicita la hora a un seguidor y registra su respuesta.
func (l *Leader) RequestTimeFromFollower(followerName, followerAddress string) {
	message := `{"operation": "GET_TIME", "leaderName": "` + l.aAbstractNode.Name + `"}`

	response, err := l.aAbstractNode.SendMessageSync(followerAddress, message)
	if err != nil {
		l.NonRespondingFollowers.Store(followerName, followerAddress)
		l.Logger.Printf("Error al solicitar tiempo al seguidor %s: %v", followerName, err)
		return
	}

	l.SuccessfulFollowers.Store(followerName, response)
	l.Logger.Printf("Hora recibida del seguidor %s: %s", followerName, response)
}

// PrintResults muestra los resultados del algoritmo.
func (l *Leader) PrintResults() {
	l.Logger.Println("Resultados del algoritmo Berkeley:")

	l.Logger.Println("Seguidores inalcanzables:")
	l.UnreachableFollowers.Range(func(key, value any) bool {
		l.Logger.Printf("- %s: %v", key, value)
		return true
	})

	l.Logger.Println("Seguidores que no respondieron a tiempo:")
	l.NonRespondingFollowers.Range(func(key, value any) bool {
		l.Logger.Printf("- %s: %v", key, value)
		return true
	})

	l.Logger.Println("Seguidores que respondieron correctamente:")
	l.SuccessfulFollowers.Range(func(key, value any) bool {
		l.Logger.Printf("- %s: %v", key, value)
		return true
	})

	l.Logger.Println("Seguidores que actualizaron su tiempo correctamente:")
	l.TimeUpdatedFollowers.Range(func(key, value any) bool {
		l.Logger.Printf("- %s: %v", key, value)
		return true
	})

	l.Logger.Println("Seguidores que no pudieron actualizar su tiempo:")
	l.FailedFollowers.Range(func(key, value any) bool {
		l.Logger.Printf("- %s: %v", key, value)
		return true
	})
}
func (l *Leader) Close() {
	l.aAbstractNode.Close()
}
