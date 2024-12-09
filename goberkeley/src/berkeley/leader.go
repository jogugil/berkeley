package berkeley

import (
<<<<<<< HEAD
	"encoding/json"
	"fmt"

	"log"
	"strconv"
=======
	"log"
>>>>>>> 4df5b3de8e63625e208a4e5d62139f2f2b1cb612
	"sync"
	"time"
)

<<<<<<< HEAD
// Mensaje JSON que se envía al seguidor.
type TimeRequest struct {
	Message    string `json:"message"`
	Operation  string `json:"operation"`
	Time       int64  `json:"time"`
	LeaderAddr string `json:"leader_address"`
}
type DeltaRequest struct {
	Message    string `json:"message"`
	Operation  string `json:"operation"`
	Delta      int64  `json:"delta"`
	LeaderAddr string `json:"leader_address"`
}
type closeRequest struct {
	Message    string `json:"message"`
	Operation  string `json:"operation"`
	LeaderAddr string `json:"leader_address"`
}

// Leader representa el nodo líder en el sistema Berkeley.
type Leader struct {
	aAbstractNode          *AbstractNode
	UnreachableFollowers   map[string]*FollowerInfo
	SuccessfulFollowers    map[string]*FollowerInfo
	NonRespondingFollowers map[string]*FollowerInfo
	TimeUpdatedFollowers   map[string]*FollowerInfo
	FailedFollowers        map[string]*FollowerInfo
	//mu                     sync.Mutex // Mutex para proteger los mapas en accesos concurrentes
	Logger *log.Logger
}

// HandleProcess implements Handler.
func (l *Leader) HandleProcess(message string) (string, error) {
	panic("unimplemented")
}

// InitializeLeaderNode crea e inicializa un nuevo nodo líder.
func InitializeLeaderNode(name, address string, timeout time.Duration, nodeAddresses map[string]string) (*Leader, error) {
	// Suponiendo que InitializeNodeWithAddresses crea un nodo base y devuelve un puntero a AbstractNode
	baseNode, err := InitializeNodeWithAddresses(name, address, timeout, nodeAddresses)
=======
// Leader representa el nodo líder en el sistema Berkeley.
type Leader struct {
	AbstractNode
	UnreachableFollowers   sync.Map // Seguidores inalcanzables
	SuccessfulFollowers    sync.Map // Seguidores que respondieron con éxito
	NonRespondingFollowers sync.Map // Seguidores que no respondieron a tiempo
	TimeUpdatedFollowers   sync.Map // Seguidores que actualizaron su tiempo correctamente
	FailedFollowers        sync.Map // Seguidores que no pudieron actualizar su tiempo
	Logger                 *log.Logger
}

// NewLeader crea e inicializa un nuevo nodo líder.
func NewLeader(name, address string, timeout time.Duration, nodeAddresses map[string]string) (*Leader, error) {
	baseNode, err := NewAbstractNode(name, address, timeout)
>>>>>>> 4df5b3de8e63625e208a4e5d62139f2f2b1cb612
	if err != nil {
		return nil, err
	}

<<<<<<< HEAD
	// Inicializamos el líder, asignando el nodo base y un logger
	leader := &Leader{
		aAbstractNode: baseNode, // Asignamos el puntero a AbstractNode
		Logger:        log.Default(),
	}
	leader.aAbstractNode.Handler = leader

	return leader, nil
}
func (l *Leader) initStructs() {
	if l.UnreachableFollowers == nil {
		l.UnreachableFollowers = make(map[string]*FollowerInfo)
	}
	if l.SuccessfulFollowers == nil {
		l.SuccessfulFollowers = make(map[string]*FollowerInfo)
	}
	if l.NonRespondingFollowers == nil {
		l.NonRespondingFollowers = make(map[string]*FollowerInfo)
	}
	if l.TimeUpdatedFollowers == nil {
		l.TimeUpdatedFollowers = make(map[string]*FollowerInfo)
	}
	if l.FailedFollowers == nil {
		l.FailedFollowers = make(map[string]*FollowerInfo)
	}
=======
	leader := &Leader{
		AbstractNode: *baseNode,
		Logger:       log.Default(),
	}

	leader.InitializeNodeWithAddresses(nodeAddresses)
	return leader, nil
}

// SendCloseMessage envía un mensaje de cierre a un seguidor.
func (l *Leader) SendCloseMessage(followerAddress string) {
	message := `{"operation": "CLOSE", "message": "Cerrar conexión", "leaderName": "` + l.Name + `"}`

	response, err := l.SendMessageSync(followerAddress, message)
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
>>>>>>> 4df5b3de8e63625e208a4e5d62139f2f2b1cb612
}

// StartAlgorithm implementa el algoritmo de sincronización Berkeley para el líder.
func (l *Leader) StartAlgorithm() {
	l.Logger.Println("Iniciando algoritmo de sincronización Berkeley...")
<<<<<<< HEAD
	l.initStructs()
	// Simula el envío de solicitudes de tiempo a los seguidores
	l.processFollowers()

	// Fase 2: Calcular el delta con la media de los tiempos
	log.Println("** Fase 2 **: Calcular el delta con la media de los tiempos")
	delta := l.calculateDeltaTimeDifference()
	if delta != 0 {
		// Paso 3: Actualizar relojes de los seguidores
		log.Println("** Paso 3 **: Llamar a los seguidores para actualizar sus relojes")
		l.callFollowersWithUpdatedTime(delta)

		// Fase 4: Enviar mensaje de cierre
		log.Println("** Fase 4 **: Enviar mensaje de cierre a los seguidores")
		l.sendCloseMessagesToFollowers()

		// Fase 5: Mostrar los resultados finales
		log.Println("** Fase 5: Mostrar los resultados de la sincronización")
		l.printResults()
	} else {
		// Se registra esta situación para comprobarlo más adelante en los logs
		log.Println("El resultado del delta es cero por lo que no se envían actualizaciones a ningún seguidor.")
	}
}
func (l *Leader) processFollowers() {
	leaderTime := time.Now().UnixMilli() // T0
	leaderAddr := l.aAbstractNode.Address

	log.Printf("processFollowers: leaderAddr %s leaderTime %d.", leaderAddr, leaderTime)

	// Lista de seguidores
	followers := l.aAbstractNode.NodeAddresses

	// Canal para los resultados
	results := make(chan *FollowerInfo, len(followers))
	var wg sync.WaitGroup

	// Enviar solicitudes a los seguidores concurrentemente
	for followerName, followerAddr := range followers {
		wg.Add(1)
		go func(name, addr string) {
			defer wg.Done()
			log.Printf("Enviando solicitud de tiempo a %s (%s).", name, addr)
			l.sendTimeRequestToFollower(name, addr, leaderTime, leaderAddr, results)
		}(followerName, followerAddr)
	}

	// Cerrar el canal una vez que todos los goroutines hayan terminado
	go func() {
		wg.Wait()
		close(results)
	}()

	// Recoger y procesar los resultados
	for res := range results {
		if res.GetState() == "RESPONDED" {
			// Seguidor que respondió correctamente
			l.SuccessfulFollowers[res.GetName()] = res
			log.Printf("Seguidor %s respondió correctamente con tiempo local: %d ms", res.GetName(), res.GetLocalTime())
		} else if res.GetState() == "TIMEOUT" {
			// Seguidor que no respondió a tiempo
			l.NonRespondingFollowers[res.GetName()] = res
			log.Printf("Seguidor %s no respondió a tiempo.", res.GetName())
		} else {
			// Seguidor que tuvo algún otro problema
			l.FailedFollowers[res.GetName()] = res
			log.Printf("Seguidor %s falló con estado: %s.", res.GetName(), res.GetState())
		}
	}

	// Log final para indicar que el procesamiento de seguidores ha terminado
	log.Printf("Proceso de seguidores completado. Respuestas procesadas: %d, Fallos: %d.", len(l.SuccessfulFollowers), len(l.NonRespondingFollowers))
}
func (l *Leader) sendTimeRequestToFollower(followerName, followerAddr string, leaderTime int64, leaderAddr string, results chan<- *FollowerInfo) {
	// Crear el mensaje JSON
	request := TimeRequest{
		Message:    "Requesting time sync",
		Operation:  "GET_TIME",
		Time:       leaderTime, //T0
		LeaderAddr: leaderAddr,
	}

	// Serializar el mensaje en formato JSON
	requestData, err := json.Marshal(request)
	if err != nil {
		log.Printf("Error al serializar la solicitud a JSON: %v", err)
		results <- NewFollowerInfo(followerAddr, followerName, 0, 0, 0)
		return
	}

	// Convertir el JSON a string
	requestString := string(requestData)
	log.Printf("Solicitud enviada a %s: %s", followerAddr, requestString)

	// Enviar el mensaje al seguidor y recibir la respuesta
	reply, err := l.aAbstractNode.SendMessageSync(followerAddr, requestString)
	if err != nil {
		log.Printf("Error al recibir respuesta de %s: %v", followerAddr, err)
		followError := NewFollowerInfo(followerAddr, followerName, 0, 0, 0)
		results <- followError
		return
	}

	log.Printf("Respuesta recibida de %s: %s", followerAddr, reply)

	// Deserializar la respuesta JSON
	var response map[string]string
	if err := json.Unmarshal([]byte(reply), &response); err != nil {
		log.Printf("Error al procesar la respuesta de %s: %v", followerAddr, err)
		results <- NewFollowerInfo(followerAddr, followerName, 0, 0, 0)
		return
	}

	// Calcular la diferencia de tiempo
	endCommTime := time.Now().UnixMilli() // T0
	timeComm := endCommTime - leaderTime
	log.Printf("Tiempo de comunicación: %d ms", timeComm)

	// Obtener el tiempo local del seguidor
	localTimeStr, ok := response["localTime"]
	if !ok {
		log.Printf("No se encontró el campo 'localTime' en la respuesta de %s", followerAddr)
		results <- NewFollowerInfo(followerAddr, followerName, 0, 0, 0)
		return
	}

	// Convertir el valor de localTime de string a int64
	localTimeInt64, err := strconv.ParseInt(localTimeStr, 10, 64)
	if err != nil {
		log.Printf("Error al convertir 'localTime' a int64 de %s: %v", followerAddr, err)
		results <- NewFollowerInfo(followerAddr, followerName, 0, 0, 0)
		return
	}

	// Enviar los resultados al canal
	log.Printf("Tiempo local recibido de %s: %d", followerAddr, localTimeInt64)
	foll := NewFollowerInfo(followerAddr, followerName, localTimeInt64, timeComm, 0)
	foll.SetState(Responded)
	results <- foll
}

// calculateDeltaTimeDifference calcula la diferencia de tiempo (delta) entre el tiempo local y los tiempos de los seguidores válidos.
func (l *Leader) calculateDeltaTimeDifference() int64 {
	// Log de inicio de la operación de cálculo de la diferencia de tiempo
	log.Println("Calculando la diferencia de tiempo (delta)")

	// Obtener el tiempo actual del líder
	now := time.Now().UnixMilli()
	var sumTime int64 = 0           // Suma de los tiempos ajustados de los seguidores
	var validFollowersCount int = 0 // Contador de seguidores con respuestas válidas

	// Recorrer los seguidores exitosos y calcular la diferencia de tiempo para cada uno
	for _, follower := range l.SuccessfulFollowers {
		if follower.DiffTime != int64(^uint64(0)>>1) { // Verificar si el seguidor tiene un tiempo válido (Long.MAX_VALUE en Java)
			followerTime := now + follower.DiffTime // Calcular el tiempo ajustado del seguidor
			sumTime += followerTime                 // Sumar el tiempo ajustado
			validFollowersCount++                   // Incrementar el contador de seguidores válidos
		}
	}

	// Si se tienen seguidores válidos, calcular la diferencia de tiempo promedio
	if validFollowersCount > 0 {
		newNow := sumTime / int64(validFollowersCount) // Calcular el promedio de los tiempos de los seguidores
		delta := newNow - now                          // Calcular la diferencia de tiempo (delta)

		// Log de los resultados calculados
		log.Printf("Nuevo tiempo calculado (new_now): %d\n", newNow)
		log.Printf("Diferencia global (δ): %d\n", delta)

		// Retornar la diferencia de tiempo (delta)
		return delta
	} else {
		// Si no hay seguidores válidos, registrar advertencia y retornar 0
		log.Println("No se han recibido respuestas válidas de los seguidores.")
		return 0
	}
}

func (l *Leader) callFollowersWithUpdatedTime(delta int64) error {
	log.Printf("Enviando actualización de tiempo a los seguidores con delta: %d", delta)

	// Crear un canal para gestionar las tareas concurrentes
	ch := make(chan *FollowerInfo, len(l.SuccessfulFollowers))

	// Usamos un WaitGroup para esperar que todas las goroutines terminen
	var wg sync.WaitGroup

	// Enviar las tareas en paralelo para cada seguidor
	for _, follower := range l.SuccessfulFollowers {
		wg.Add(1) // Incrementamos el contador del WaitGroup
		go func(follower FollowerInfo) {
			defer wg.Done() // Decrementamos el contador del WaitGroup cuando la goroutine termina
			followerInfo := l.sendTimeUpdateToFollower(follower, delta)
			ch <- followerInfo
		}(*follower)
	}

	// Iniciar una goroutine para cerrar el canal una vez que todas las goroutines hayan terminado
	go func() {
		wg.Wait() // Esperamos a que todas las goroutines terminen
		close(ch) // Cerramos el canal después de que todas las goroutines hayan enviado sus respuestas
	}()

	// Procesar las respuestas conforme vayan llegando
	for followerInfo := range ch {
		if followerInfo.State == "TIME_UPDATED" {
			log.Printf("El seguidor %s respondió correctamente al cambio del timer.", followerInfo.Name)
			l.TimeUpdatedFollowers[followerInfo.Name] = followerInfo
		} else {
			log.Printf("El seguidor %s no respondió al cambio de su timer.", followerInfo.Name)
			l.FailedFollowers[followerInfo.Name] = followerInfo
		}
	}

	return nil
}

func (l *Leader) sendTimeUpdateToFollower(follower FollowerInfo, delta int64) *FollowerInfo {

	// Crear el mapa con la solicitud
	request := DeltaRequest{
		Message:    "Modifica el tiempo del sistema de tu servidor con el diferencial.",
		Operation:  "UPDATE_TIME",
		Delta:      delta,
		LeaderAddr: l.aAbstractNode.Address,
	}

	// Serializar la solicitud a JSON
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		log.Printf("Error al serializar la solicitud para el seguidor %s: %v", follower.Name, err)
		follower.State = TimeErrorSentUpdate
		return &follower
	}

	requestString := string(jsonRequest)
	log.Printf("Solicitud enviada a %s: %s", follower.GetAddress(), requestString)
	// Enviar la solicitud
	reply, err := l.aAbstractNode.SendMessageSync(follower.GetAddress(), requestString)
	if err != nil {
		log.Printf("Error al recibir respuesta de %s: %v", follower.GetAddress(), err)
		followError := NewFollowerInfo(follower.GetAddress(), follower.GetName(), 0, 0, 0)
		followError.State = TimeErrorSentUpdate
		return followError
	}

	var response map[string]string
	err = json.Unmarshal([]byte(reply), &response)
	if err != nil {
		log.Printf("Error al deserializar la respuesta del seguidor %s: %v", follower.Name, err)
		follower.State = TimeErrorSentUpdate
		return &follower
	}

	// Registrar la respuesta
	followerName := response["followerName"]
	operation := response["operation"]
	log.Printf("Respuesta de %s: Operación %s exitosa", followerName, operation)

	follower.State = TimeUpdated
	return &follower
}
func (l *Leader) sendCloseMessagesToFollowers() {
	// Crear un WaitGroup para esperar a que todas las goroutines terminen
	var wg sync.WaitGroup

	// Canal para recibir los resultados de las goroutines
	resultCh := make(chan string, len(l.TimeUpdatedFollowers))

	// Enviar mensaje de cierre a cada seguidor de manera concurrente
	for _, follower := range l.TimeUpdatedFollowers {
		wg.Add(1)
		go func(follower FollowerInfo) {
			defer wg.Done()

			// Enviar el mensaje de cierre
			err := l.sendCloseMessage(follower.Address, follower.Name)
			if err != nil {
				resultCh <- fmt.Sprintf("Error al enviar mensaje de cierre a %s: %s", follower.Name, err)
				return
			}

			// Simulamos que el mensaje fue enviado exitosamente
			resultCh <- fmt.Sprintf("Mensaje de cierre enviado con éxito a %s", follower.Name)
		}(*follower)
	}

	// Esperar a que todas las goroutines terminen
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Procesar los resultados a medida que vayan llegando
	for result := range resultCh {
		log.Println(result)
	}

}

// SendCloseMessage envía un mensaje de cierre a un seguidor.
// Simulación de la función para enviar mensaje de cierre a un seguidor
func (l *Leader) sendCloseMessage(followerAddress, followerName string) *FollowerInfo {
	log.Printf("Enviando mensaje de cierre a %s...", followerAddress)

	request := closeRequest{
		Message:    "Cerrar conexión",
		Operation:  "CLOSE",
		LeaderAddr: l.aAbstractNode.Address,
	}

	// Convertir el mensaje a JSON
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		log.Printf("error al crear mensaje de cierre: %w", err)
		followError := NewFollowerInfo(followerAddress, followerName, 0, 0, 0)
		followError.State = ErrorClose
		return followError
	}

	// Simulamos el envío del mensaje al seguidor (usando un canal)
	log.Printf("Mensaje de cierre enviado a %s: %s", followerAddress, jsonRequest)
	// Serializar la solicitud a JSON

	requestString := string(jsonRequest)
	log.Printf("Solicitud enviada a %s: %s", followerAddress, requestString)
	// Enviar la solicitud
	reply, err := l.aAbstractNode.SendMessageSync(followerAddress, requestString)
	if err != nil {
		log.Printf("Error al recibir respuesta de %s: %v", followerAddress, err)
		followError := NewFollowerInfo(followerAddress, followerName, 0, 0, 0)
		followError.State = ErrorClose
		return followError
	}

	var response map[string]string
	err = json.Unmarshal([]byte(reply), &response)
	if err != nil {
		log.Printf("Error al deserializar la respuesta del seguidor %s: %v", followerName, err)
		followError := NewFollowerInfo(followerAddress, followerName, 0, 0, 0)
		followError.State = ErrorClose
		return followError
	}
	// Simulamos una respuesta (en un caso real recibirías la respuesta desde el socket)
	// Aquí se simula un éxito.
	response = map[string]string{
		"followerName": followerAddress,
		"operation":    "CLOSE",
	}

	// Convertimos la respuesta a JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("error al procesar la respuesta: %w", err)
		followError := NewFollowerInfo(followerAddress, followerName, 0, 0, 0)
		followError.State = ErrorClose
		return followError
	}

	// Simulamos la recepción de la respuesta (en un caso real, recibirías el mensaje desde el socket)
	log.Printf("Respuesta recibida de %s: %s", followerAddress, jsonResponse)
	follow := NewFollowerInfo(followerAddress, followerName, 0, 0, 0)
	follow.State = OkClose
	return follow
}

// PrintResults muestra los resultados del algoritmo.
func (l *Leader) printResults() {
	l.Logger.Println("Resultados del algoritmo Berkeley:")

	// Mostrar seguidores inalcanzables
	l.Logger.Println("Seguidores inalcanzables:")
	for key, value := range l.UnreachableFollowers {
		l.Logger.Printf("- %s: %v", key, value)
	}

	// Mostrar seguidores que no respondieron a tiempo
	l.Logger.Println("Seguidores que no respondieron a tiempo:")
	for key, value := range l.NonRespondingFollowers {
		l.Logger.Printf("- %s: %v", key, value)
	}

	// Mostrar seguidores que respondieron correctamente
	l.Logger.Println("Seguidores que respondieron correctamente:")
	for key, value := range l.SuccessfulFollowers {
		l.Logger.Printf("- %s: %v", key, value)
	}

	// Mostrar seguidores que actualizaron su tiempo correctamente
	l.Logger.Println("Seguidores que actualizaron su tiempo correctamente:")
	for key, value := range l.TimeUpdatedFollowers {
		l.Logger.Printf("- %s: %v", key, value)
	}

	// Mostrar seguidores que no pudieron actualizar su tiempo
	l.Logger.Println("Seguidores que no pudieron actualizar su tiempo:")
	for key, value := range l.FailedFollowers {
		l.Logger.Printf("- %s: %v", key, value)
	}
}
func (l *Leader) Close() {
	l.aAbstractNode.Close()
}
=======

	// Simula el envío de solicitudes de tiempo a los seguidores
	l.NodeAddresses.Range(func(followerName, followerAddress any) bool {
		go l.RequestTimeFromFollower(followerName.(string), followerAddress.(string))
		return true
	})

	// Aquí puedes implementar la lógica de cálculo del delta y la sincronización.
}

// RequestTimeFromFollower solicita la hora a un seguidor y registra su respuesta.
func (l *Leader) RequestTimeFromFollower(followerName, followerAddress string) {
	message := `{"operation": "GET_TIME", "leaderName": "` + l.Name + `"}`

	response, err := l.SendMessageSync(followerAddress, message)
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

>>>>>>> 4df5b3de8e63625e208a4e5d62139f2f2b1cb612
