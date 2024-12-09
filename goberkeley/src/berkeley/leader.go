package berkeley

import (
	"encoding/json"
	"fmt"

	"log"
	"strconv"
	"sync"
	"time"
)

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
	leader.aAbstractNode.Handler = leader

	return leader, nil
}
func (l *Leader) initializeStructs() {
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
}

// StartAlgorithm implementa el algoritmo de sincronización Berkeley para el líder.
func (l *Leader) StartAlgorithm() {
	l.Logger.Println("\n\n\t*************** Iniciando algoritmo de sincronización Berkeley... *****************")
	log.Println(" ")
	log.Println(" ")

	l.initializeStructs()

	// Simula el envío de solicitudes de tiempo a los seguidores
	log.Println("\n\n\t** Fase 1 **:  Petición de tiempos a los seguidores y calculo de sus diferncias.")
	log.Println(" ")

	l.processFollowers()

	// Fase 2: Calcular el delta con la media de los tiempos
	log.Println("\n\n\t** Fase 2 **: Calcular el delta con la media de los tiempos")
	log.Println(" ")

	delta := l.calculateDeltaTimeDifference()

	if delta != 0 {
		// Paso 3: Actualizar relojes de los seguidores
		log.Println("\n\n\t** Paso 3 **: Llamar a los seguidores para actualizar sus relojes")
		log.Println(" ")

		l.callFollowersWithUpdatedTime(delta)

		// Fase 4: Enviar mensaje de cierre
		log.Println("\n\n\t** Fase 4 **: Enviar mensaje de cierre a los seguidores")
		log.Println(" ")

		l.sendCloseMessagesToFollowers() //Realmente no es del algoritmo pero para evitar problemas con los contextos de Go y sockets sincronizo el cierre!!!

		// Fase 5: Mostrar los resultados finales
		log.Println("\n\n\t** Fase 5: Mostrar los resultados de la sincronización")
		log.Println(" ")

		l.printResults()
	} else {
		// Se registra esta situación para comprobarlo más adelante en los logs
		log.Println("El resultado del delta es cero por lo que no se envían actualizaciones a ningún seguidor.")
	}
}

///////// FASE 1:

// processFollowers es responsable de gestionar las solicitudes de tiempo a los seguidores.
// Envía solicitudes de tiempo a los seguidores de manera concurrente y luego procesa las respuestas.
// Dependiendo del estado de la respuesta, clasifica a los seguidores en tres grupos:
// - Seguidores que respondieron correctamente.
// - Seguidores que no respondieron a tiempo.
// - Seguidores que tuvieron algún otro problema.
// La función también registra los resultados y la cantidad de respuestas procesadas.
func (l *Leader) processFollowers() {
	// Obtener el tiempo actual del líder en milisegundos (T0)
	leaderTime := time.Now().UnixMilli() // T0

	// Obtener la dirección del líder
	leaderAddr := l.aAbstractNode.Address

	// Registrar la información del líder y el tiempo de la solicitud
	log.Printf("processFollowers: leaderAddr %s leaderTime %d.", leaderAddr, leaderTime)

	// Lista de direcciones de los seguidores
	followers := l.aAbstractNode.NodeAddresses

	// Crear un canal para recibir los resultados de las respuestas de los seguidores
	results := make(chan *FollowerInfo, len(followers))

	// Usar un WaitGroup para esperar que todos los goroutines terminen
	var wg sync.WaitGroup

	// Enviar solicitudes de tiempo a los seguidores concurrentemente
	for followerName, followerAddr := range followers {
		// Incrementar el contador del WaitGroup para cada goroutine
		wg.Add(1)

		// Crear un goroutine para cada seguidor
		go func(name, addr string) {
			defer wg.Done() // Decrementar el contador del WaitGroup cuando termine este goroutine
			log.Printf("Enviando solicitud de tiempo a %s (%s).", name, addr)
			// Enviar la solicitud de tiempo al seguidor y recibir la respuesta en el canal
			l.sendTimeRequestToFollower(name, addr, leaderTime, leaderAddr, results)
		}(followerName, followerAddr)
	}

	// Cerrar el canal una vez que todos los goroutines hayan terminado
	go func() {
		// Esperar a que todos los goroutines terminen
		wg.Wait()
		// Cerrar el canal de resultados
		close(results)
	}()

	// Recoger y procesar los resultados de los seguidores
	for res := range results {
		// Clasificar la respuesta según el estado
		if res.GetState() == "RESPONDED" {
			// Seguidor que respondió correctamente
			l.SuccessfulFollowers[res.GetName()] = res
			log.Printf("Seguidor %s respondió correctamente con tiempo local: %d ms", res.GetName(), res.GetFollowerTime())
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

// sendTimeRequestToFollower envía una solicitud de sincronización de tiempo a un seguidor específico.
// Se crea un mensaje de solicitud en formato JSON que incluye la dirección del líder y el tiempo actual (T0).
// La función espera la respuesta del seguidor, la procesa y calcula el tiempo de comunicación.
// Si la respuesta es válida, se obtiene el tiempo local del seguidor y se calcula el tiempo de comunicación.
// Los resultados se envían al canal de resultados con la información relevante.
func (l *Leader) sendTimeRequestToFollower(followerName, followerAddr string, leaderTime int64, leaderAddr string, results chan<- *FollowerInfo) {

	// Crear el mensaje JSON con la solicitud de sincronización de tiempo
	request := TimeRequest{
		Message:    "Requesting time sync", // Mensaje de la solicitud
		Operation:  "GET_TIME",             // Operación que se está solicitando
		Time:       leaderTime,             // El tiempo actual del líder (T0)
		LeaderAddr: leaderAddr,             // Dirección del líder
	}

	// Serializar el mensaje en formato JSON
	requestData, err := json.Marshal(request)
	if err != nil {
		// Si ocurre un error al serializar, se registra y se envía una respuesta de error al canal
		log.Printf("Error al serializar la solicitud a JSON: %v", err)
		results <- NewFollowerInfo(followerAddr, followerName, 0, 0, 0, 0)
		return
	}

	// Convertir el JSON a un string para su envío
	requestString := string(requestData)
	log.Printf("Solicitud enviada a %s: %s", followerAddr, requestString)

	// Enviar el mensaje al seguidor y recibir la respuesta
	reply, err := l.aAbstractNode.SendMessageSync(followerAddr, requestString)
	if err != nil {
		// Si ocurre un error al recibir la respuesta, se registra y se envía un error al canal
		log.Printf("Error al recibir respuesta de %s: %v", followerAddr, err)
		followError := NewFollowerInfo(followerAddr, followerName, 0, 0, 0, 0)
		results <- followError
		return
	}

	// Registrar la respuesta recibida
	log.Printf("Respuesta recibida de %s: %s", followerAddr, reply)

	// Deserializar la respuesta JSON del seguidor
	var response map[string]string
	if err := json.Unmarshal([]byte(reply), &response); err != nil {
		// Si ocurre un error al procesar la respuesta, se registra y se envía un error al canal
		log.Printf("Error al procesar la respuesta de %s: %v", followerAddr, err)
		results <- NewFollowerInfo(followerAddr, followerName, 0, 0, 0, 0)
		return
	}

	// Calcular el tiempo de comunicación entre el líder y el seguidor
	endCommTime := time.Now().UnixMilli() // T0 al final de la comunicación
	timeComm := endCommTime - leaderTime  // Tiempo de comunicación en milisegundos
	log.Printf("Tiempo de comunicación: %d ms", timeComm)

	// Obtener el tiempo local del seguidor desde la respuesta
	localTimeStr, ok := response["localTime"]
	if !ok {
		// Si no se encuentra el campo 'localTime' en la respuesta, se registra y se envía un error
		log.Printf("No se encontró el campo 'localTime' en la respuesta de %s", followerAddr)
		results <- NewFollowerInfo(followerAddr, followerName, 0, 0, 0, 0)
		return
	}

	// Convertir el valor de 'localTime' de string a int64
	followerTime, err := strconv.ParseInt(localTimeStr, 10, 64)
	if err != nil {
		// Si ocurre un error al convertir 'localTime', se registra y se envía un error
		log.Printf("Error al convertir 'localTime' a int64 de %s: %v", followerAddr, err)
		results <- NewFollowerInfo(followerAddr, followerName, 0, 0, 0, 0)
		return
	}

	// Enviar los resultados al canal, incluyendo el tiempo local, el tiempo de comunicación y la diferencia entre el lider y el seguidor
	// La diferncia se calcula al crear el objeto FollowerInfo. diff = (TP + trip_Time) - Now. Trip_time = (Now - T0)/2
	log.Printf("Tiempo local recibido de %s: %d", followerAddr, followerTime)
	foll := NewFollowerInfo(followerAddr, followerName, followerTime, leaderTime, timeComm, endCommTime)
	foll.SetState(Responded) // Marcar la respuesta como "RESPONDED"
	results <- foll
}

////////// FASE 2:

// calculateDeltaTimeDifference calcula la diferencia de tiempo (delta) entre el tiempo local del líder
// y el tiempo ajustado promedio de los seguidores válidos. La función recorre los seguidores exitosos,
// ajusta sus tiempos en función de sus diferencias de tiempo y calcula una diferencia global promedio (δ).
// Retorna la diferencia de tiempo calculada (delta) o 0 si no hay seguidores válidos.
func (l *Leader) calculateDeltaTimeDifference() int64 {
	// Log de inicio de la operación de cálculo de la diferencia de tiempo
	log.Println("Calculando la diferencia de tiempo (delta)")

	// Obtener el tiempo actual del líder en milisegundos
	now := time.Now().UnixMilli()

	// Inicializar variables para la suma de los tiempos ajustados y el contador de seguidores válidos
	var sumTime int64 = 0           // Suma de los tiempos ajustados de los seguidores
	var validFollowersCount int = 0 // Contador de seguidores con respuestas válidas

	// Recorrer los seguidores exitosos y calcular la diferencia de tiempo para cada uno
	for _, follower := range l.SuccessfulFollowers {
		// Verificar si el seguidor tiene un tiempo válido, evitando el valor especial de "Long.MAX_VALUE"
		if follower.DiffTime != int64(^uint64(0)>>1) { // Long.MAX_VALUE en Java
			// Ajustar el tiempo del seguidor sumando su diferencia de tiempo al tiempo actual del líder
			followerTime := now + follower.DiffTime
			// Sumar el tiempo ajustado del seguidor al total
			sumTime += followerTime
			// Incrementar el contador de seguidores válidos
			validFollowersCount++
		}
	}

	// Si hay seguidores válidos, calcular la diferencia de tiempo promedio
	if validFollowersCount > 0 {
		// Calcular el tiempo ajustado promedio de los seguidores
		newNow := sumTime / int64(validFollowersCount)
		// Calcular la diferencia de tiempo (delta) entre el tiempo promedio de los seguidores y el tiempo local del líder
		delta := newNow - now

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

////////// FASE 3:

// callFollowersWithUpdatedTime envía una actualización de tiempo a todos los seguidores con el delta de tiempo especificado.
// La actualización se realiza en paralelo utilizando goroutines para cada seguidor, y las respuestas se procesan conforme
// van llegando. La función maneja la concurrencia mediante un canal y un WaitGroup para asegurarse de que todas las
// goroutines terminen antes de procesar los resultados.
func (l *Leader) callFollowersWithUpdatedTime(delta int64) error {
	// Log que muestra el inicio de la actualización de tiempo a los seguidores con el delta calculado
	log.Printf("Enviando actualización de tiempo a los seguidores con delta: %d", delta)

	// Crear un canal para gestionar las respuestas de los seguidores, con un buffer del tamaño del número de seguidores exitosos
	ch := make(chan *FollowerInfo, len(l.SuccessfulFollowers))

	// Crear un WaitGroup para esperar a que todas las goroutines terminen su ejecución
	var wg sync.WaitGroup

	// Enviar las tareas en paralelo para cada seguidor exitoso
	for _, follower := range l.SuccessfulFollowers {
		wg.Add(1) // Incrementamos el contador del WaitGroup antes de iniciar cada goroutine
		// Goroutine para enviar la actualización de tiempo a un seguidor específico
		go func(follower FollowerInfo) {
			defer wg.Done() // Decrementamos el contador del WaitGroup cuando la goroutine termina
			// Enviar la actualización de tiempo al seguidor y obtener la respuesta
			followerInfo := l.sendTimeUpdateToFollower(&follower, delta)
			// Enviar la respuesta al canal para su posterior procesamiento
			ch <- followerInfo
		}(*follower) // Llamamos a la goroutine pasando el valor de 'follower'
	}

	// Iniciar una goroutine para cerrar el canal una vez que todas las goroutines hayan terminado
	go func() {
		wg.Wait() // Esperamos a que todas las goroutines terminen su ejecución
		close(ch) // Cerramos el canal después de que todas las respuestas hayan sido enviadas
	}()

	// Procesar las respuestas conforme vayan llegando del canal
	for followerInfo := range ch {
		// Si el estado del seguidor es "TIME_UPDATED", se ha actualizado correctamente
		if followerInfo.State == "TIME_UPDATED" {
			log.Printf("El seguidor %s respondió correctamente al cambio del timer.", followerInfo.Name)
			// Guardamos el seguidor como actualizado correctamente en la lista de seguidores actualizados
			l.TimeUpdatedFollowers[followerInfo.Name] = followerInfo
		} else {
			// Si el seguidor no respondió correctamente, lo agregamos a la lista de seguidores fallidos
			log.Printf("El seguidor %s no respondió al cambio de su timer.", followerInfo.Name)
			l.FailedFollowers[followerInfo.Name] = followerInfo
		}
	}

	// Retornamos nil indicando que no hubo errores en la ejecución de la función
	return nil
}

// sendTimeUpdateToFollower envía una solicitud para actualizar el tiempo del seguidor con un delta especificado.
// La solicitud es serializada a formato JSON y enviada de forma sincrónica al seguidor. Luego, se procesa la respuesta
// y se actualiza el estado del seguidor según el resultado. Si hay algún error en el proceso, se registra y se devuelve
// un seguidor con un estado de error.
func (l *Leader) sendTimeUpdateToFollower(follower *FollowerInfo, delta int64) *FollowerInfo {
	// Crear el mapa con la solicitud para modificar el tiempo del sistema del seguidor
	request := DeltaRequest{
		Message:    "Modifica el tiempo del sistema de tu servidor con el diferencial.",
		Operation:  "UPDATE_TIME",
		Delta:      delta,
		LeaderAddr: l.aAbstractNode.Address,
	}

	// Serializar la solicitud a JSON
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		// Si hay un error al serializar la solicitud, se registra el error y se marca el estado del seguidor como de error
		log.Printf("Error al serializar la solicitud para el seguidor %s: %v", follower.Name, err)
		follower.State = TimeErrorSentUpdate
		// Devolver el seguidor con estado de error
		return follower
	}

	// Convertir la solicitud serializada en una cadena JSON
	requestString := string(jsonRequest)
	log.Printf("Solicitud enviada a %s: %s", follower.GetAddress(), requestString)

	// Enviar la solicitud de manera sincrónica y esperar la respuesta
	reply, err := l.aAbstractNode.SendMessageSync(follower.GetAddress(), requestString)
	if err != nil {
		// Si hay un error al recibir la respuesta, se registra el error y se crea un nuevo objeto de seguidor con estado de error
		log.Printf("Error al recibir respuesta de %s: %v", follower.GetAddress(), err)
		followError := follower
		followError.State = TimeErrorSentUpdate
		// Devolver el seguidor con error
		return followError
	}

	// Deserializar la respuesta JSON recibida del seguidor
	var response map[string]string
	err = json.Unmarshal([]byte(reply), &response)
	if err != nil {
		// Si hay un error al deserializar la respuesta, se registra el error y se marca el estado del seguidor como de error
		log.Printf("Error al deserializar la respuesta del seguidor %s: %v", follower.Name, err)
		follower.State = TimeErrorSentUpdate
		// Devolver el seguidor con error
		return follower
	}

	// Registrar la respuesta exitosa del seguidor
	followerName := response["followerName"]
	operation := response["operation"]
	log.Printf("Respuesta de %s: Operación %s exitosa", followerName, operation)

	// Modificamos el seguidor con la información del delta que se uso para actualizar la hora local del seguidor y su estado
	follwerUpdate := follower
	follwerUpdate.SetDelta(delta)
	follwerUpdate.State = TimeUpdated // Marcar el estado como "TimeUpdated" (actualizado)

	// Devolver el seguidor con los datos actualizados
	return follwerUpdate
}

////////// FASE 4:

// sendCloseMessagesToFollowers envía un mensaje de cierre a todos los seguidores que han sido actualizados correctamente.
// Utiliza goroutines para enviar los mensajes de forma concurrente y espera que todas las goroutines terminen antes de
// finalizar el proceso. Los resultados de las operaciones son procesados a medida que van llegando y se registran.
func (l *Leader) sendCloseMessagesToFollowers() {
	// Crear un WaitGroup para esperar a que todas las goroutines terminen
	var wg sync.WaitGroup

	// Canal para recibir los resultados de las goroutines
	resultCh := make(chan string, len(l.TimeUpdatedFollowers))

	// Enviar mensaje de cierre a cada seguidor de manera concurrente
	for _, follower := range l.TimeUpdatedFollowers {
		wg.Add(1) // Incrementamos el contador del WaitGroup para cada goroutine
		go func(follower FollowerInfo) {
			defer wg.Done() // Decrementamos el contador cuando la goroutine termine

			// Enviar el mensaje de cierre al seguidor
			err := l.sendCloseMessage(&follower)
			if err != nil {
				// Si hay un error al enviar el mensaje, se envía un resultado con el error al canal
				resultCh <- fmt.Sprintf("Error al enviar mensaje de cierre a %s: %s", follower.Name, err)
				return
			}

			// Si el mensaje fue enviado correctamente, se simula el éxito y se envía al canal
			resultCh <- fmt.Sprintf("Mensaje de cierre enviado con éxito a %s", follower.Name)
		}(*follower) // Llamamos a la goroutine con la información del seguidor
	}

	// Iniciar una goroutine para esperar que todas las goroutines terminen y cerrar el canal de resultados
	go func() {
		wg.Wait()       // Esperamos a que todas las goroutines terminen
		close(resultCh) // Cerramos el canal cuando se haya completado el procesamiento
	}()

	// Procesar los resultados a medida que vayan llegando
	for result := range resultCh {
		// Registrar cada resultado recibido del canal
		log.Println(result)
	}
}

// SendCloseMessage envía un mensaje de cierre a un seguidor.
// Simulación de la función para enviar mensaje de cierre a un seguidor
func (l *Leader) sendCloseMessage(follower *FollowerInfo) *FollowerInfo {

	followerAddress := follower.GetAddress()

	followerName := follower.GetName()

	request := closeRequest{
		Message:    "Cerrar conexión",
		Operation:  "CLOSE",
		LeaderAddr: l.aAbstractNode.Address,
	}

	// Convertir el mensaje a JSON
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		wrappedErr := fmt.Errorf("Error during some operation: %w", err)
		log.Printf("An error occurred: %v", wrappedErr) // Logging the wrapped error
		followError := follower
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
		followError := follower
		followError.State = ErrorClose
		return followError
	}

	var response map[string]string
	err = json.Unmarshal([]byte(reply), &response)
	if err != nil {
		log.Printf("Error al deserializar la respuesta del seguidor %s: %v", followerName, err)
		followError := follower
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
		wrappedErr := fmt.Errorf("Error during some operation: %w", err)
		log.Printf("An error occurred: %v", wrappedErr) // Logging the wrapped error
		followError := follower
		followError.State = ErrorClose
		return followError
	}

	// Simulamos la recepción de la respuesta (en un caso real, recibirías el mensaje desde el socket)
	log.Printf("Respuesta recibida de %s: %s", followerAddress, jsonResponse)
	follow := follower
	follow.State = OkClose
	return follow
}

////////// FASE 5:

// PrintResults muestra los resultados del algoritmo Berkeley de manera organizada y tabulada.
func (l *Leader) printResults() {
	// Guardar los flags originales
	originalFlags := log.Flags()
	// Desactivar los flags temporalmente antes de mostrar el resumen
	log.SetFlags(0)

	l.Logger.Println("\n\t===============================")
	l.Logger.Println("\t Resultados del algoritmo Berkeley")
	l.Logger.Println("\t===============================")

	// Mostrar seguidores inalcanzables
	l.Logger.Println("\n\t❌\tSeguidores inalcanzables:")
	if len(l.UnreachableFollowers) == 0 {
		l.Logger.Println("\t\t\tNo hay seguidores inalcanzables.")
	} else {
		for key, value := range l.UnreachableFollowers {
			l.Logger.Printf("\t\t\t- %s:\t%v\n", key, value)
		}
	}

	// Mostrar seguidores que no respondieron a tiempo
	l.Logger.Println("\n\t⏳\tSeguidores que no respondieron a tiempo:")
	if len(l.NonRespondingFollowers) == 0 {
		l.Logger.Println("\t\t\tTodos los seguidores respondieron a tiempo.")
	} else {
		for key, value := range l.NonRespondingFollowers {
			l.Logger.Printf("\t\t\t- %s:\t%v\n", key, value)
		}
	}

	// Mostrar seguidores que respondieron correctamente
	l.Logger.Println("\n\t✅\tSeguidores que respondieron correctamente:")
	if len(l.SuccessfulFollowers) == 0 {
		l.Logger.Println("\t\t\tNo hubo respuestas correctas.")
	} else {
		for key, value := range l.SuccessfulFollowers {
			l.Logger.Printf("\t\t\t- %s:\t%v\n", key, value)
		}
	}

	// Mostrar seguidores que actualizaron su tiempo correctamente
	l.Logger.Println("\n\t🕰️\tSeguidores que actualizaron su tiempo correctamente:")
	if len(l.TimeUpdatedFollowers) == 0 {
		l.Logger.Println("\t\t\t Ningún seguidor actualizó su tiempo.")
	} else {
		for key, value := range l.TimeUpdatedFollowers {
			l.Logger.Printf("\t\t\t- %s:\t%v\n", key, value)
		}
	}

	// Mostrar seguidores que no pudieron actualizar su tiempo
	l.Logger.Println("\n\t❌\tSeguidores que no pudieron actualizar su tiempo:")
	if len(l.FailedFollowers) == 0 {
		l.Logger.Println("\t\t\tTodos los seguidores pudieron actualizar su tiempo.")
	} else {
		for key, value := range l.FailedFollowers {
			l.Logger.Printf("\t\t\t- %s:\t%v\n", key, value)
		}
	}

	l.Logger.Println("\t\t===============================")

	// Restaurar los flags originales
	log.SetFlags(originalFlags)
}

func (l *Leader) Close() {
	l.aAbstractNode.Close()
}

// HandleProcess implements Handler.
func (l *Leader) HandleProcess(message string) (string, error) {
	panic("unimplemented")
}
