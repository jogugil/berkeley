package org.jogugi.berkeley;

import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.ExecutorCompletionService;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.zeromq.SocketType;
import org.zeromq.ZContext;
import org.zeromq.ZMQ;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;

/**
 * @brief Representa un nodo líder en el sistema.
 *
 * La clase Leader gestiona el estado de los seguidores (followers) en el sistema.
 * Administra los seguidores que están alcanzables, exitosos, no responden, actualizados o fallidos,
 * y se comunica con ellos a través de mensajes.
 */
public class Leader extends AbstractNode {

	/**
	 * @brief Logger para registrar la actividad de la clase Leader.
	 */
    private static final Logger logger = LoggerFactory.getLogger(Leader.class);
    
    /**
     * @brief Mapa de seguidores no alcanzables.
     */
    private ConcurrentHashMap<String, FollowerInfo> unreachableFollowers = new ConcurrentHashMap<>();

    /**
     * @brief Mapa de seguidores exitosos.
     */
    private ConcurrentHashMap<String, FollowerInfo> successfulFollowers = new ConcurrentHashMap<>();

    /**
     * @brief Mapa de seguidores que no respondieron.
     */
    private ConcurrentHashMap<String, FollowerInfo> nonRespondingFollowers = new ConcurrentHashMap<>();

    /**
     * @brief Mapa de seguidores con tiempo actualizado.
     */
    private ConcurrentHashMap<String, FollowerInfo> timeUpdatedFollowers = new ConcurrentHashMap<>();

    /**
     * @brief Mapa de seguidores fallidos.
     */
    private ConcurrentHashMap<String, FollowerInfo> failedFollowers = new ConcurrentHashMap<>();
    
    

    /**
     * @brief Inicializa el nodo líder con la dirección, nombre, timeout y lista de seguidores.
     * 
     * Este método configura el nodo líder. Llama al método de inicialización de la clase base
     * para configurar el contexto y los parámetros comunes del nodo.
     * 
     * @param address Dirección del nodo líder.
     * @param name Nombre del nodo líder.
     * @param timeout Tiempo de espera en milisegundos para las operaciones de comunicación.
     * @param followers Lista de direcciones de los seguidores que están conectados al líder.
     */
    public void initializeNode(String address, String name, int timeout, HashMap<String, String> followers) {
        super.initializeNode(name, address, timeout, followers);
        logger.info("Líder inicializado en {} con timeout de {} ms", address, timeout);
    }
    
    /**
     * Inicia el algoritmo de sincronización de relojes entre el líder y los seguidores.
     * El algoritmo se ejecuta en varias fases:
     * 
     * - **Fase 1**: Solicitar la hora a todos los seguidores en paralelo.
     * - **Fase 2**: Calcular el delta entre los tiempos de los seguidores obteniendo la media de esos tiempos (delta).
     * - **Fase 3**: Actualizar los relojes de los seguidores con el delta calculado.
     * - **Fase 4**: Enviar un mensaje de cierre a los seguidores.
     * - **Fase 5**: Mostrar los resultados finales de la sincronización.
     * 
     * El método utiliza un pool de hilos para ejecutar la solicitud de hora a cada seguidor en paralelo.
     * Se gestionan las respuestas de cada seguidor y se realiza el cálculo necesario para sincronizar los relojes.
     * En esta versión, para simplificar el código, generramos tantos hilos en el pool como seguidores teenmos.
     * 
     * @throws InterruptedException Si el hilo es interrumpido durante la ejecución.
     * @throws ExecutionException Si ocurre un error durante la ejecución de las tareas en paralelo.
     */
    @Override
    public void startAlgorithm () throws InterruptedException, ExecutionException {
        logger.info ("Iniciando algoritmo de sincronización...");
        
        // Fase 1: Solicitar hora a todos los seguidores en paralelo y procesar sus respuestas
        logger.debug ("** Fase 1 **: Solicitar hora a todos los seguidores y obtener sus diferencias");
        processFollowers ();  // Llama a la función que interactúa con los seguidores para obtener sus tiempos

        // Fase 2: Calcular el delta con la media de los tiempos
        logger.debug("** Fase 2 **: Calcular el delta con la media de los tiempos");
        long delta = calculateDeltaTimeDifference();
        if (delta != 0) {
            // Paso 3: Actualizar relojes de los seguidores
            logger.debug("** Paso 3 **: Llamar a los seguidores para actualizar sus relojes");
            callFollowersWithUpdatedTime(delta);

            // Fase 4: Enviar mensaje de cierre
            logger.debug("** Fase 4 **: Enviar mensaje de cierre a los seguidores");
            sendCloseMessagesToFollowers();

            // Fase 5: Mostrar los resultados finales
            //logger.debug("** Fase 5: Mostrar los resultados de la sincronización");
            printResults();       	
        } else {
        	//Se registra esta situación para comrobarlo mñas adelante en lso logs
        	logger.warn("El resultado del delta es cero por lqo eu no se envia ctuakizaciones a ningun seguidor.");
        }

    }

 // FASE 1: Enviar la petición de hora a cada seguior y obtne su diferncia respecto al lider
    
    /**
     * Procesa las respuestas de todos los seguidores en paralelo.
     * 
     * Este método solicita la hora a todos los seguidores y luego procesa las respuestas a medida que llegan. 
     * Los seguidores son procesados en paralelo utilizando un `ExecutorService` y `ExecutorCompletionService` 
     * para manejar la ejecución de múltiples hilos.
     * 
     * Fases:
     * - **Fase 1.1**: Lanzar tareas concurrentes para solicitar la hora a todos los seguidores.
     * - **Fase 1.2**: Procesar los resultados conforme se completan.
     * 
     * Los resultados de las respuestas de los seguidores se almacenan en dos mapas:
     * - **successfulFollowers**: Los seguidores que respondieron correctamente.
     * - **nonRespondingFollowers**: Los seguidores que no respondieron o respondieron incorrectamente.
     * 
     * @throws InterruptedException Si el hilo principal es interrumpido mientras espera los resultados.
     */
    private void processFollowers () throws InterruptedException {
        
    	// Inicializamos el ExecutorService y CompletionService
        ExecutorService executor = Executors.newFixedThreadPool (Math.min (nodeAddresses.size (),
                Runtime.getRuntime().availableProcessors ()));
        ExecutorCompletionService<FollowerInfo> completionService = new ExecutorCompletionService<> (executor);

        // Mapa para almacenar los futuros y los nombres de los seguidores
        Map<String, Future<FollowerInfo>> futureToFollowerMap = new HashMap<> ();

        // Fase 1.1: Lanzar las tareas concurrentes
        logger.debug ("** Fase 1.1 **: Solicitar hora a todos los seguidores");
        long T0 = System.currentTimeMillis (); // Pasamosa todos lso seguidores la hora local del líder

        // Se lanza una tarea por cada seguidor en la lista de direcciones
        nodeAddresses.forEach ((name, address) -> {
            // Enviamos el trabajo a ejecutar
            Future<FollowerInfo> future = completionService.submit (() -> processFollowerTime (name, address, T0));
            futureToFollowerMap.put (name,  future);  // Guardamos el futuro asociado con el nombre del seguidor
        });

     // Fase 1.2: Procesar los resultados conforme se completen
        logger.debug ("** Fase 1.2 **: Procesar los resultados de cada seguidor según llegan");

        int completedTasks     = 0;                            // Contamos los trabajos que van llegando
        successfulFollowers    = new ConcurrentHashMap<> ();   // Seguidores que han contestado
        nonRespondingFollowers = new ConcurrentHashMap<> ();   // Seguidores que no han contestado

        // Lista de futuros que todavía están pendientes de completar
        HashMap<String, Future<FollowerInfo>> pendingFutures = new HashMap<>(futureToFollowerMap);

        // Esperamos todas las peticiones de hora, bien porque llegan las respuestas o se alcanza el timeout global.
        while (completedTasks < nodeAddresses.size()) {
            try {
                // Esperamos la siguiente respuesta o un timeout
                Future<FollowerInfo> future = completionService.poll (timeout, TimeUnit.MILLISECONDS);

                if (future == null) {
                    // Timeout global alcanzado: Registramos los seguidores pendientes y salimos del loop
                    logger.warn ("Timeout global alcanzado. Algunos seguidores no respondieron.");
                    for (Map.Entry<String, Future<FollowerInfo>> entry : pendingFutures.entrySet()) {
                        Future<FollowerInfo> pendingFuture = entry.getValue();
                        if (pendingFuture != null) {
                            String followerName = getFollowerName(pendingFuture, futureToFollowerMap);
                            
                            // Solo procesamos aquellos seguidores que no han respondido aún
                            if (!successfulFollowers.containsKey(followerName)) {
                                // Crear una instancia de FollowerInfo con estado "no respondido"
                                FollowerInfo nonRespondingInfo = new FollowerInfo (
                                    nodeAddresses.get(followerName), // Dirección del seguidor
                                    followerName,                    // Nombre del seguidor
                                    0,                               // Hora local (desconocida)
                                    0,                               // Tiempo de comunicación (0 indica no respondido)
                                    0       						 // Tiempo de registro de no respuesta, se calcula cuando se obtiene el delta
                                );
                                nonRespondingInfo.setState(FollowerState.NO_RESPONSE); // Estado inicial por defecto
                                nonRespondingFollowers.put(followerName, nonRespondingInfo);
                            }
                        }
                    }
                    break; // Salimos del loop porque ya no esperamos más respuestas al alcanar el timeout  (único para todos)
                }

                // Procesamos la respuesta del seguidor
                FollowerInfo followerInfo = future.get ();
                String followerName = getFollowerName (future, futureToFollowerMap);
                logger.warn ("El seguidor {}  con estado {}.", followerName, followerInfo.getState());
                if (followerInfo.getState() == FollowerState.RESPONDED) {
                	logger.warn ("El seguidor {} respondió correctamente.", followerName);
                    successfulFollowers.put (followerName, followerInfo);
                 }else if (followerInfo.getState() == FollowerState.REQUEST_NOT_SENT) {
                	logger.warn ("No se pudo enviar la peticion de tiempo al seguidor {}.", followerName);
                	unreachableFollowers.put (followerName, followerInfo);
                } else {
                   logger.warn ("El seguidor {} no respondió correctamente.", followerName);
                   nonRespondingFollowers.put (followerName, followerInfo);
                }

                 
                // Eliminar este futuro de la lista pendiente utilizando su referencia en el mapa
                pendingFutures.put(followerName, null);
                completedTasks++;

            } catch (Exception e) {
                logger.error ("Error en la ejecución del hilo principal.", e);
                completedTasks++; // Contamos como tarea completada
            }
        }

        // Al final, loggeamos los resultados
        logger.info ("Tareas completadas: {}", completedTasks);
        logger.info ("Seguidores exitosos: {}", successfulFollowers.keySet ());
        logger.info ("Seguidores que no respondieron: {}", nonRespondingFollowers.keySet ());
        // Cerramos el ExecutorService después de procesar todas las tareas
        executor.shutdown();
    }

    /**
     * Procesa la solicitud de hora a un seguidor y calcula el tiempo de comunicación.
     * 
     * Este método establece una conexión con un seguidor a través de un socket de tipo REQ/REP utilizando ZMQ.
     * Envía una solicitud para obtener la hora local del seguidor y calcula el tiempo que tarda en recibir la respuesta.
     * Si se recibe una respuesta válida, devuelve un objeto `FollowerInfo` con la información del seguidor.
     * Si no se recibe respuesta o ocurre un error, se registra la incidencia y devuelve una respuesta con valores por defecto.
     * 
     * El proceso de solicitud de tiempo se realiza de la siguiente manera:
     * - Se crea un socket de tipo REQ para realizar la solicitud de hora.
     * - Se envía una solicitud en formato JSON al seguidor con la hora inicial (`T0`).
     * - Se espera la respuesta con la hora local del seguidor.
     * - Si la respuesta es válida, se calcula el tiempo de comunicación y se devuelve el resultado.
     * - Si no se recibe respuesta o ocurre un error, se gestionan estos casos mediante valores predeterminados.
     * 
     * @param followerName Nombre del seguidor al que se le solicita la hora.
     * @param followerAddress Dirección del seguidor al que se conecta.
     * @param T0 Tiempo de referencia desde el cual se mide el tiempo de comunicación.
     * @return Un objeto `FollowerInfo` con la hora del seguidor, el tiempo de comunicación, y otros datos relevantes.
     * @throws JsonProcessingException Si ocurre un error al procesar el JSON.
     */
    private FollowerInfo processFollowerTime (String followerName, String followerAddress, long T0) {
        if (context == null) {
            context = new ZContext(); // Creamos el contexto de ZMQ si no existe
        }
        ZMQ.Socket requester = null;
        try {
            requester = context.createSocket (SocketType.REQ);
            logger.info ("Conectando al seguidor: {} con timeout {} en {}", followerAddress, timeout, followerAddress);
            requester.setReceiveTimeOut (timeout);               // Establecer el tiempo de espera
            requester.connect ("tcp://" + followerAddress);      // Conectar al seguidor

            // Crear y enviar la solicitud
            ObjectMapper mapper = new ObjectMapper();
            Map<String, String> request = new HashMap<> ();
            request.put ("leaderName", name);               // Nombre del líder
            request.put ("operation", "GET_TIME");          // Tipo de operación
            request.put ("message", "Dame tu hora local");  // Mensaje a enviar
            request.put ("T0", String.valueOf(T0));         // Tiempo de referencia

            String jsonRequest = mapper.writeValueAsString (request); // Convertir solicitud a JSON
            requester.send (jsonRequest.getBytes(ZMQ.CHARSET), 0);    // Enviar la solicitud

            // Esperar la respuesta
            byte[] reply = requester.recv (0); // Recibir la respuesta
            if (reply != null) {
                String jsonResponse = new String (reply, ZMQ.CHARSET); // Convertir la respuesta a string
                // PAra luego parsear la respuesta
                Map<String, String> response = mapper.readValue(
                	    jsonResponse, new TypeReference<Map<String, String>>() {}
                	);
                String localTime       = response.get ("localTime");                       // Obtener la hora local del seguidor
                long endTime           = System.currentTimeMillis ();                      // Calcular el tiempo final
                long communicationTime = endTime - T0;                                     // Calcular el tiempo de comunicación

                if (localTime == null || localTime.isEmpty()) {
                    logger.warn ("El seguidor {} devolvió un tiempo nulo o vacío", followerName);
                    return new FollowerInfo (followerAddress, followerName, 0, communicationTime, endTime); // Manejo básico de error
                }

                logger.info("Respuesta de {}: Hora local={}, Tiempo comunicación={}ms", 
                            followerName, localTime, communicationTime);
                FollowerInfo respondingInfo = new FollowerInfo (
                		followerAddress,  					    // Dirección del seguidor
                        followerName,                           // Nombre del seguidor
                        Long.valueOf (localTime).longValue (),  // Hora local (desconocida)
                        communicationTime,                      // Tiempo de comunicación (0 indica no respondido)
                        0       								// Se inicializa cuando se calcula el delta
                    );
                respondingInfo.setState (FollowerState.RESPONDED);
                return respondingInfo;
            } else {
                logger.warn("No hubo respuesta del seguidor: {}", followerAddress);
                FollowerInfo nonRespondingInfo = new FollowerInfo (
                        nodeAddresses.get (followerName),  // Dirección del seguidor
                        followerName,                      // Nombre del seguidor
                        0,                                 // Hora local (desconocida)
                        0,                                 // Tiempo de comunicación (0 indica no respondido)
                        0						           // Se inicializa cuando se calcula el delta
                    );
               nonRespondingInfo.setState (FollowerState.NO_RESPONSE);  
               return nonRespondingInfo; // Si no se recibe respuesta
            }
        } catch (Exception e) {
            logger.error ("Error en la comunicación con el seguidor: {}", e.getMessage ());
            logger.warn("No hubo respuesta del seguidor: {}", followerAddress);
            FollowerInfo unreachableInfo = new FollowerInfo (
                    nodeAddresses.get (followerName),  // Dirección del seguidor
                    followerName,                      // Nombre del seguidor
                    0,                                 // Hora local (desconocida)
                    0,                                 // Tiempo de comunicación (0 indica no respondido)
                    0						           // Se inicializa cuando se calcula el delta
                );
            unreachableInfo.setState (FollowerState.REQUEST_NOT_SENT);
            return unreachableInfo; // Devolver un valor por defecto
        } finally {
            if (requester != null) {
                requester.close (); // Cerrar el socket después de usarlo
            }
        }
    }
    
// FASE 2: Calcular el delta 
 
    /**
     * Calcula la diferencia de tiempo (delta) entre el tiempo local y los tiempos de los seguidores válidos.
     * 
     * Esta función calcula un promedio de los tiempos ajustados de los seguidores que han enviado respuestas
     * válidas, y calcula la diferencia entre el tiempo local del líder y el promedio calculado. 
     * El valor de delta se usa para sincronizar el tiempo entre el líder y los seguidores. Si no se reciben 
     * respuestas válidas de los seguidores, se retorna 0.
     * 
     * @return La diferencia de tiempo (delta) calculada en milisegundos. Si no hay seguidores válidos,
     *         retorna 0.
     */
    private long calculateDeltaTimeDifference () {
        // Log de inicio de la operación de cálculo de la diferencia de tiempo
        logger.info ("Calculando la diferencia de tiempo (delta)");
        
        // Obtener el tiempo actual del líder
        long now                = System.currentTimeMillis();
        long sumTime            = 0;  // Suma de los tiempos ajustados de los seguidores
        int validFollowersCount = 0;  // Contador de seguidores con respuestas válidas

        // Recorrer los seguidores exitosos y calcular la diferencia de tiempo para cada uno
        for (FollowerInfo follower : successfulFollowers.values ()) {
            if (follower.getDiffTime () != Long.MAX_VALUE) {        // Verificar si el seguidor tiene un tiempo válido
                long followerTime = now + follower.getDiffTime ();  // Calcular el tiempo ajustado del seguidor
                sumTime += followerTime;   // Sumar el tiempo ajustado
                validFollowersCount++;     // Incrementar el contador de seguidores válidos
            }
        }

        // Si se tienen seguidores válidos, calcular la diferencia de tiempo promedio
        if (validFollowersCount > 0) {
            long new_now = sumTime / validFollowersCount;  // Calcular el promedio de los tiempos de los seguidores
            long delta   = new_now - now;                  // Calcular la diferencia de tiempo (delta)
            
            // Log de los resultados calculados
            logger.info("Nuevo tiempo calculado (new_now): {}", new_now);
            logger.info("Diferencia global (δ): {}", delta);
            
            // Retornar la diferencia de tiempo (delta)
            return delta;
        } else {
            // Si no hay seguidores válidos, registrar advertencia y retornar 0
            logger.warn("No se han recibido respuestas válidas de los seguidores.");
            return 0;
        }
    }
    
  //FASE 3: Enviarle la petición de actualizaión de fecha del sistema a cada seguidor válido

    /**
     * Envía actualizaciones de tiempo a los seguidores en paralelo utilizando un `ExecutorService`.
     * 
     * Esta función recorre los seguidores exitosos y envía una solicitud de actualización de tiempo
     * a cada uno de ellos. Las actualizaciones se realizan en paralelo utilizando un pool de hilos.
     * Se gestionan las respuestas de cada seguidor y se espera un resultado con un tiempo de espera
     * definido. Si no se recibe respuesta a tiempo, se registra un mensaje de advertencia.
     * 
     * @param delta El valor diferencial de tiempo que se debe aplicar a cada seguidor.
     * @throws ExecutionException Si ocurre un error durante la ejecución de alguna de las tareas en paralelo.
     */
    private void callFollowersWithUpdatedTime(long delta) throws ExecutionException {
        logger.info("Enviando actualización de tiempo a los seguidores con delta: {}", delta);

        // Crear un pool de hilos con un tamaño igual al número de seguidores
        int poolSize = successfulFollowers.size ();
        ExecutorService executorService = Executors.newFixedThreadPool (poolSize);
        ExecutorCompletionService<FollowerInfo> completionService = new ExecutorCompletionService<> (executorService);

        // Mapa para almacenar las futuras tareas asociadas a cada seguidor
        Map<String, Future<FollowerInfo>> followerFutures = new HashMap<>();

        // Enviar las tareas para cada seguidor en paralelo
        for (FollowerInfo follower : successfulFollowers.values()) {
             // Enviar tarea para cada seguidor
            completionService.submit(() -> {
                // Llamar al método para enviar el mensaje y obtener el FollowerInfo
                FollowerInfo follw = sendTimeUpdateToFollower (follower, delta); // Enviar el mensaje al seguidor
                // Aquí envolvemos el FollowerInfo en un CompletableFuture y lo agregamos al mapa
                // Guardamos el Future<FollowerInfo> en el mapa, usando el nombre del seguidor como clave
                followerFutures.put(follw.getName(), CompletableFuture.completedFuture(follw));
                
                return follw; // Retornamos el FollowerInfo
            });
        }

        // Esperar y procesar las respuestas conforme vayan llegando
        int completedTasks = 0;
        while (completedTasks < followerFutures.size ()) {
            try {
                // Tomar la tarea completada
                Future<FollowerInfo> completedFuture = completionService.take (); // Bloquea hasta que una tarea se complete
                // Obtener el FollowerInfo asociado al Future completado
                FollowerInfo followerInfo = completedFuture.get();
                String followerName = getFollowerName (completedFuture, followerFutures);

                // Comprobamos si la tarea terminó correctamente o con timeout
                if (completedFuture.isDone ()) {
	                logger.warn ("El seguidor {}  con estado {}.", followerName, followerInfo.getState());
	                followerInfo.setState(FollowerState.TIME_UPDATED);
	                logger.warn ("El seguidor {} respondió correctamente al cambio del timer.", followerName);
	                timeUpdatedFollowers.put (followerName, followerInfo);

                } else {
                	logger.warn ("El seguidor {}  con estado {} no contesto al cambio de su timer.", followerName, followerInfo.getState());
                	failedFollowers.put (followerName, followerInfo);
                }
                completedTasks++; // Aumentar el contador de tareas completadas
            } catch (InterruptedException e) {
                // Manejar la excepción de interrupción
                logger.error ("Error al procesar una tarea de actualización de tiempo.", e);
                completedTasks++; // Aumentar el contador de tareas completadas en caso de error
                // Capturamos la excepción lanzada por el Future
                Throwable cause = e.getCause ();  // La causa original de la excepción
                String followerName = "Desconocido";
                if (cause instanceof Exception) {
                    logger.error("El Future para el seguidor falló debido a: {}", cause.getMessage());
                }
                FollowerInfo unreachableInfo = new FollowerInfo (
                        nodeAddresses.get (followerName),  // Dirección del seguidor
                        followerName,                      // Nombre del seguidor
                        0,                                 // Hora local (desconocida)
                        0,                                 // Tiempo de comunicación (0 indica no respondido)
                        0						           // Se inicializa cuando se calcula el delta
                    );
                failedFollowers.put (followerName, unreachableInfo); // Registrar como fallido
            }
        }

        // Apagar el executor después de procesar todas las tareas
        executorService.shutdown ();
    }
    
    /**
     * Envía una solicitud para actualizar el tiempo de un seguidor en el sistema.
     * 
     * Esta función establece una conexión de tipo `REQ` con un seguidor a través de ZeroMQ,
     * enviando una solicitud de actualización de tiempo con un valor delta proporcionado.
     * Si la respuesta es recibida, se procesa y se muestra en el log.
     * Si no se recibe respuesta o ocurre un error, se registra un mensaje adecuado.
     *
     * @param followerAddress Dirección del seguidor a donde enviar la solicitud.
     * @param delta El valor diferencial de tiempo que se debe aplicar en el seguidor.
     */
    private FollowerInfo sendTimeUpdateToFollower (FollowerInfo follower, long delta) {
        // Verificamos si el contexto de ZeroMQ ya está creado. Si no, lo inicializamos.
        if (context == null) {
            context = new ZContext ();
        }
       
        // Bloque try-with-resources para asegurarse de que el socket se cierre automáticamente.
        try (ZMQ.Socket requester = context.createSocket (SocketType.REQ)) {
        	 String followerAddress = follower.getAdrressFollower();
        	 
            // Configuramos el tiempo de espera para la respuesta.
            requester.setReceiveTimeOut (timeout);

            // Establecemos la conexión con el seguidor.
            requester.connect ("tcp://" + followerAddress);

            // Creamos un mapa con la solicitud que se enviará.
            Map<String, String> request = new HashMap<> ();
            request.put ("leaderName", name);
            request.put ("operation", "UPDATE_TIME");
            request.put ("delta", String.valueOf (delta));
            request.put ("message", "Modifica el tiemp del sistema de tu servidor con el diferencial.");

            // Serializamos la solicitud a formato JSON utilizando ObjectMapper.
            ObjectMapper mapper = new ObjectMapper ();
            String jsonRequest  = mapper.writeValueAsString (request);

            // Enviamos la solicitud al seguidor.
            requester.send (jsonRequest.getBytes (ZMQ.CHARSET), 0);

            // Esperamos la respuesta del seguidor.
            byte[] reply = requester.recv (0);

            // Si recibimos una respuesta, la procesamos.
            if (reply != null) {
                // Convertimos la respuesta de bytes a String.
                String jsonResponse = new String (reply, ZMQ.CHARSET);

                // Deserializamos la respuesta en un mapa.
                Map<String, String> response = mapper.readValue (jsonResponse, new TypeReference<Map<String, String>>() {});

                // Extraemos los valores de la respuesta.
                String followerName = response.get ("followerName");
                String operation    = response.get ("operation");

                // Registramos que la operación fue exitosa.
                logger.info("Respuesta de {}: Operación: {} exitosa", followerName, operation);
                follower.setState (FollowerState.TIME_UPDATED);
            } else { // Si no se recibe respuesta, se registra una advertencia.
                logger.warn ("No hubo respuesta del seguidor: {}", followerAddress);
            }
            return follower;
        } catch (Exception e) {
            // Si ocurre un error, se registra el mensaje de error.
            logger.error ("Error al enviar actualización de tiempo a {}: {}", follower.getAdrressFollower(), e.getMessage ());
             follower.setState (FollowerState.TIME_ERROR_SENT_UPDATE);
             return follower;
        }
    }
  
 //Fas 4 Indicarles a lso seguidores que se teermino el proceso y si tienen algún socket con el lider que lo cierren
    
    /**
     * @brief Envía un mensaje de cierre a un seguidor especificado.
     * 
     * Este método crea una solicitud de cierre para desconectar un seguidor del líder. El mensaje
     * se envía de manera síncrona a la dirección del seguidor proporcionada. Si el seguidor responde,
     * se registra la respuesta, indicando si la operación de cierre fue realizada correctamente.
     * 
     * @param followerAddress Dirección del seguidor al que se enviará el mensaje de cierre.
     * 
     * @note Si no se recibe respuesta del seguidor, se emite una advertencia en los logs.
     */
    private void sendCloseMessage(String followerAddress) {
        if (context == null) {
            context = new ZContext();
        }
        ZMQ.Socket requester = null;
        try {
            requester = context.createSocket(SocketType.REQ);
            requester.setReceiveTimeOut(timeout);
            requester.connect("tcp://" + followerAddress);

            // Crear y enviar el mensaje de cierre
            Map<String, String> request = new HashMap<>();
            request.put("leaderName", name);
            request.put("operation", "CLOSE");
            request.put("message", "Cerrar conexión");

            ObjectMapper mapper = new ObjectMapper();
            String jsonRequest = mapper.writeValueAsString(request);
            requester.send(jsonRequest.getBytes(ZMQ.CHARSET), 0);

            logger.info("Mensaje de cierre enviado a {}", followerAddress);

            // Esperar la respuesta
            byte[] reply = requester.recv(0);
            if (reply != null) {
                String jsonResponse = new String(reply, ZMQ.CHARSET);
                Map<String, String> response = mapper.readValue(
                    jsonResponse, new TypeReference<Map<String, String>>() {}
                );
                String followerName = response.get("followerName");
                String operation = response.get("operation");

                logger.info("Respuesta de {}: Operación : {} realizada", followerName, operation);
            } else {
                logger.warn("No hubo respuesta del seguidor: {}", followerAddress);
            }

        } catch (Exception e) {
            logger.error("Error al enviar mensaje de cierre a {}: {}", followerAddress, e.getMessage());
        } finally {
            if (requester != null) {
                requester.close();
            }
        }
    }
    /**
     * @brief Envía un mensaje de cierre a todos los seguidores que respondieron con éxito.
     * 
     * Este método utiliza un `ExecutorService` con un grupo de hilos para enviar un mensaje de cierre a
     * cada uno de los seguidores que han respondido correctamente. El mensaje de cierre es enviado de manera 
     * asíncrona, y el método espera las respuestas o el tiempo de espera para cada uno de los seguidores.
     * 
     * Para cada seguidor:
     * - Se envía un mensaje de cierre en un hilo separado.
     * - Se espera hasta que el mensaje se haya enviado correctamente o se alcance un tiempo de espera.
     * - Se registran las respuestas o se manejan los errores (por ejemplo, tiempo de espera alcanzado).
     * 
     * @throws ExecutionException Si ocurre un error durante la ejecución de las tareas asíncronas.
     */
    private void sendCloseMessagesToFollowers () throws ExecutionException {
        // Crear un ExecutorService para manejar las tareas de cierre
        ExecutorService executor = Executors.newFixedThreadPool (successfulFollowers.size());
        ExecutorCompletionService<FollowerInfo> completionService = new ExecutorCompletionService<> (executor);
        Map<String, Future<FollowerInfo>> closeFutures = new HashMap<> ();

        // Enviar el mensaje de cierre a los seguidores que respondieron con éxito
        for (Map.Entry<String, FollowerInfo> entry : successfulFollowers.entrySet ()) {
            //String followerName    = entry.getKey ();
            String followerAddress = entry.getValue ().getAdrressFollower ();  

            // Enviar mensaje de cierre en un hilo separado
            Future<FollowerInfo> closeFuture = executor.submit(() -> {
                try {
                    sendCloseMessage (followerAddress);  // Enviar mensaje de cierre
                    logger.info("Mensaje de cierre enviado a: {}", followerAddress);
                } catch (Exception e) {
                    logger.error("Error al enviar mensaje de cierre a: {}", followerAddress, e);
                }
                return null; // La tarea no devuelve nada útil, así que devolvemos null
            });

            closeFutures.put (followerAddress, closeFuture);  // Asociamos el Future con la dirección del seguidor
            completionService.submit (() -> {
                try {
                    closeFuture.get(timeout, TimeUnit.MILLISECONDS);  // Esperamos que el Future termine
                    logger.info ("Mensaje de cierre enviado con éxito al seguidor: {}", followerAddress);
                } catch (TimeoutException e) {
                    logger.warn ("Timeout alcanzado para el seguidor: {} en el envío de mensaje de cierre.", followerAddress);
                } catch (InterruptedException | ExecutionException e) {
                    logger.error ("Error al enviar mensaje de cierre al seguidor: {}", followerAddress, e);
                }
                return null;
            });
        }

        // Esperar y procesar las respuestas conforme vayan llegando
        int completedTasks = 0;
        while (completedTasks < successfulFollowers.size()) {
            try {
                // Esperar a que cualquier tarea se complete
                Future<FollowerInfo> completedFuture = completionService.take ();
                completedTasks++;

                // Aquí ya se procesó la respuesta correspondiente al seguidor
            } catch (InterruptedException e) {
                logger.error("Error al procesar la respuesta del seguidor.", e);
            }
        }

        // Apagar el executor una vez terminadas todas las tareas
        executor.shutdown();
    }
    /**
     * @brief Cierra el socket del líder y el contexto de ZeroMQ.
     * 
     * Este método se encarga de cerrar el socket del líder si está abierto y el contexto de ZeroMQ,
     * liberando los recursos asociados. Si el socket o el contexto no están inicializados, se
     * emite un mensaje de advertencia. Si ocurre un error durante el cierre, se registra un mensaje de
     * error correspondiente.
     * 
     * @throws SocketZeroMQException Si ocurre un error durante el cierre del socket.
     */
    public void close () throws SocketZeroMQException {
        logger.info ("Cerrando nodo del líder...");

        try {
            // Verifica si el socket es nulo y, en ese caso, lo cierra.
            if (socket != null) {
                super.close (); // Cerrar el socket del líder si está inicializado.
                logger.info ("Socket del líder cerrado correctamente.");
            }
        } catch (Exception e) {
            // Manejo de errores si ocurre una excepción al cerrar el socket.
            logger.error ("Error al cerrar el socket del líder.", e);
        }

        // Verifica si el contexto de ZeroMQ está inicializado antes de cerrarlo.
        if (context != null) {
            context.close (); // Cerrar el contexto de ZeroMQ si está presente.
            logger.info ("Contexto de ZeroMQ cerrado.");
        } else {
            // En caso de que el contexto ya haya sido cerrado previamente.
            logger.warn ("El contexto de ZeroMQ ya estaba cerrado.");
        }
    }
    /**
     * Obtiene el nombre del seguidor asociado al futuro pasado como argumento.
     * 
     * Este método busca el nombre de un seguidor en el mapa de futuros (`futureMap`) en función del objeto `Future<FollowerInfo>`.
     * Si el futuro corresponde a uno de los seguidores registrados en el mapa, se devuelve el nombre asociado a ese seguidor.
     * Si no se encuentra el seguidor, se devuelve un nombre por defecto: "Unknown".
     * 
     * @param future El objeto `Future<FollowerInfo>` que se está buscando en el mapa de futuros.
     * @param futureMap El mapa de futuros (`Map<String, Future<FollowerInfo>>`) que asocia nombres de seguidores con sus futuros.
     * @return El nombre del seguidor asociado al futuro, o "Unknown" si no se encuentra el seguidor.
     */
    private String getFollowerName (Future<FollowerInfo> future, Map<String, Future<FollowerInfo>> futureMap) {
        // Buscar el nombre asociado al futuro
        for (Map.Entry<String, Future<FollowerInfo>> entry : futureMap.entrySet ()) {
            if (entry.getValue ().equals (future)) {
                return entry.getKey ();  // Devolver el nombre del seguidor
            }
        }
        return "Unknown";  // En caso de que no se encuentre el seguidor (aunque no debería pasar)
    }
    /**
     * Imprime los resultados del algoritmo, mostrando los diferentes estados de los seguidores.
     * 
     * Este método imprime tres grupos de seguidores:
     * - Seguidores inalcanzables: Aquellos que no pudieron ser contactados.
     * - Seguidores que no respondieron a tiempo: Aquellos que no enviaron su respuesta dentro del plazo esperado.
     * - Seguidores que respondieron correctamente: Aquellos que respondieron con éxito y dentro del tiempo esperado.
     */
    private void printResults () {
        logger.info ("Resultados del algoritmo:");

        // Crear una lista de mapas con el nombre del grupo y su respectiva colección
        List<Map<String, FollowerInfo>> followerGroups = List.of (
												                unreachableFollowers,
												                nonRespondingFollowers,
												                successfulFollowers,
												                timeUpdatedFollowers,
												                failedFollowers
												        );
        
        // Crear una lista con los nombres de los grupos de seguidores
        List<String> groupNames = List.of (
						                "Seguidores inalcanzables",
						                "Seguidores que no respondieron a tiempo",
						                "Seguidores que respondieron correctamente",
						                "Seguidores que respondieron correctamente al cambio de timer",
						                "Seguidores que NO respondieron al cambio de timer"
						        );

        // Iterar sobre los grupos y sus nombres
        for (int i = 0; i < followerGroups.size (); i++) {
            Map<String, FollowerInfo> group = followerGroups.get (i);
            String groupName = groupNames.get (i);
            
            logger.info(groupName + ":");
            for (Map.Entry<String, FollowerInfo> entry : group.entrySet()) {
                logger.info ("- {}: Hora={}, Tiempo={}ms", 
                            entry.getKey (), entry.getValue ().dateFollower.toString (), entry.getValue ().communicationTime);
            }
        }
    }
    /**
     * Función que debe sobrecargarse e implementar la lógica para
     * procesar el mensaje recibido y realizar la lógica de negocio correspondiente.
     * @param message El mensaje que se recibirá para ser procesado. Este mensaje podría ser de cualquier
     *                tipo, dependiendo de la implementación futura.
     * 
     * @return El resultado del procesamiento del mensaje. Actualmente siempre retorna null, pero se espera
     *         que devuelva un resultado significativo en el futuro según la lógica de procesamiento.
     * 
     * @throws Exception Si ocurre algún error durante el procesamiento del mensaje, se lanzará una excepción.
     */
    @Override
    protected String handleProcess(String message) {
        logger.debug("Procesando mensaje: {}", message);
        // TODO: Implementar la lógica de procesamiento del mensaje
        return null;
    }
}

