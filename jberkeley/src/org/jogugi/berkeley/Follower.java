package org.jogugi.berkeley;

import java.io.IOException;
import java.util.Date;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.zeromq.SocketType;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;

/**
 * @brief Clase que representa un nodo seguidor en un sistema distribuido.
 *
 * La clase `Follower` extiende la funcionalidad de la clase abstracta `AbstractNode`. 
 * Representa un nodo seguidor que interactúa con un nodo líder en un sistema distribuido. 
 * El seguidor es capaz de recibir y procesar mensajes del líder, así como manejar operaciones 
 * específicas como obtener la hora, modificarla, y cerrar la conexión.
 */
public class Follower extends AbstractNode {

    /**
     * @brief Dirección del nodo líder en formato `host:puerto`.
     */
    private String leaderAddress;

    /**
     * @brief Logger para registrar eventos y mensajes del nodo seguidor.
     */
    private static final Logger logger = LoggerFactory.getLogger(Follower.class);


    
    /**
     * @brief Inicializa el nodo seguidor con la información específica.
     *
     * Este método configura el nodo seguidor asignándole su nombre, dirección, 
     * dirección del líder y tiempo de espera. Extiende la funcionalidad de la clase base 
     * `AbstractNode`.
     *
     * @param name Nombre del nodo seguidor.
     * @param address Dirección del nodo seguidor en formato `host:puerto`.
     * @param leaderAddress Dirección del líder en formato `host:puerto`.
     * @param timeout Tiempo de espera (en milisegundos) para las operaciones del nodo.
     * @throws IOException Si ocurre un error al inicializar el nodo.
     * @throws SocketZeroMQException Si ocurre un error relacionado con ZeroMQ.
     */
    public void initializeNode (String name, String address, String leaderAddress, int timeout) 
            throws IOException, SocketZeroMQException {
        // Llamar al método de inicialización de la clase AbstractNode
        super.initializeNode (name, address, timeout);

        // Asignar la dirección del líder
        this.leaderAddress = leaderAddress;
    }
    /**
     * @brief Maneja y procesa los mensajes recibidos del líder.
     *
     * Este método deserializa un mensaje JSON enviado por el líder, identifica la operación solicitada
     * (por ejemplo, obtener la hora, modificar el tiempo o cerrar la conexión) y responde de acuerdo 
     * con la operación especificada.
     *
     * @param message El mensaje JSON recibido del líder.
     * @return Una cadena JSON con la respuesta apropiada para la operación realizada.
     * @throws SocketZeroMQException En caso de errores relacionados con el socket.
     */
    @Override
    protected String handleProcess (String message) throws SocketZeroMQException {
        // Crear un objeto para manejar la deserialización de JSON
        ObjectMapper objectMapper = new ObjectMapper ();

        try {
            // Deserializar el mensaje JSON
            JsonNode rootNode    = objectMapper.readTree (message);
            String leaderName    = rootNode.get ("leaderName").asText ();
            String operation     = rootNode.get ("operation").asText ();
            String leaderMessage = rootNode.get ("message").asText ();

            // Registrar el inicio del procesamiento del mensaje
            logger.debug ("Procesando mensaje del líder: {} desde {} con operación: {}", leaderName, leaderAddress, operation);

            // Procesar según el tipo de operación
            if ("GET_TIME".equalsIgnoreCase(operation)) {
                // Obtener la marca de tiempo del líder
                long T0          = rootNode.get ("T0").asLong ();
                long currentTime = getCurrentTime ();

                // Mostrar el mensaje del líder y calcular el tiempo promedio TP
                long TP = displayLeaderMessage (leaderName, leaderMessage, T0, currentTime);

                // Crear y devolver la respuesta
                return String.format ("{\"followerName\":\"%s\", \"localTime\":\"%s\", \"addressFollower\":\"%s\"}", 
                                                name, TP, address);
            }

            if ("MOD_TIME".equalsIgnoreCase (operation)) {
                // Obtener el valor de delta y modificar el tiempo del sistema
                long delta = rootNode.get ("DELTA").asLong ();
                return modSystemTime (delta);
            }

            if ("CLOSE".equalsIgnoreCase (operation)) {
                // Responder al cierre solicitado por el líder
                return String.format ("{\"followerName\":\"%s\",\"operation\":\"CLOSE\"}", name);
            } else {
                // Operación no reconocida
                logger.warn ("Operación no reconocida en el mensaje del líder: {}", operation);
                return "{\"error\":\"Operación no reconocida\"}";
            }
        } catch (IOException e) {
            // Manejar errores de deserialización
            logger.error ("Error al procesar el mensaje JSON: {}", e.getMessage(), e);
            return "{\"error\":\"Error al procesar el mensaje JSON\"}";
        }
    }

    /**
     * @brief Muestra el mensaje del líder y calcula un valor intermedio de tiempo (TP).
     *
     * Este método registra el mensaje del líder recibido, la hora proporcionada por el líder (T0),
     * y la hora local del seguidor. Además, calcula un tiempo promedio (TP) basado en T0 y el tiempo
     * local del seguidor, y lo registra en el log.
     *
     * @param leaderName El nombre del líder que envía el mensaje.
     * @param leaderMessage El mensaje recibido del líder.
     * @param T0 La marca de tiempo enviada por el líder en milisegundos desde la época.
     * @param localTime La hora local actual del seguidor en milisegundos desde la época.
     * @return El tiempo promedio (TP) calculado a partir de T0 y la hora local del seguidor.
     */
    private long displayLeaderMessage (String leaderName, String leaderMessage, long T0, long localTime) {
        // Registrar el mensaje del líder recibido
        logger.info ("Mensaje del líder ({}) recibido: {} con T0 {}", leaderName, leaderMessage,  (new Date(T0)).toString());

        // Calcular el tiempo promedio (TP) entre T0 y la hora local
        long TP = (T0 + localTime) / 2;
        logger.info("TP del seguidor {} es {}", name,  (new Date(TP)).toString());

        // Devolver el tiempo promedio calculado
        return TP;
    }

    /**
     * @brief Modifica el tiempo local del sistema en función de un delta.
     *
     * Este método simula un cambio en el tiempo del sistema sumando un delta al tiempo local actual.
     * También registra el tiempo antes y después de la modificación en los logs y devuelve una
     * representación JSON de la operación.
     *
     * @param delta El desplazamiento en milisegundos a aplicar al tiempo local actual.
     * @return Una cadena JSON que contiene el nombre del seguidor, el nuevo tiempo local modificado 
     *         y el estado de la operación.
     */
    private String modSystemTime (long delta) {
       
        long currentLocalTime = System.currentTimeMillis ();   // Obtiene el tiempo local actual en milisegundos desde la época (1970-01-01T00:00:00Z)
        long modSystemTime    = currentLocalTime + delta;      // Modifica el tiempo local sumando el delta proporcionado

        // Registra los tiempos antes y después de la modificación
        logger.info("Tiempo del seguidor {} modificado de {} a {}", name, (new Date (currentLocalTime)).toString(), (new Date (modSystemTime)).toString());

        // Devuelve un JSON con el nombre del seguidor, el nuevo tiempo local y el estado de la operación
        return String.format("{\"followerName\":\"%s\", \"localTime\":\"%s\",\"operation\":\"OK_MOD_TIME\"}", name, modSystemTime);
    }
    
    
    /**
     * @brief Obtiene la hora actual del sistema en milisegundos desde la época Unix.
     * 
     * Este método registra el tiempo actual en milisegundos (TP) junto con la dirección 
     * del seguidor y también lo convierte en una fecha legible para su registro.
     * 
     * @return El tiempo actual del sistema en milisegundos desde el 1 de enero de 1970 (época Unix).
     */
    private long getCurrentTime () {
        // Obtener el tiempo actual del sistema en milisegundos desde la época Unix
        long TP = System.currentTimeMillis ();

        // Convertir el valor de TP en un objeto de fecha legible
        // Registrar la fecha y hora local del seguidor
        logger.info("Fecha y hora local del seguidor: TP: {}", (new Date (TP)).toString ());

        // Devolver el tiempo en milisegundos
        return TP;
    }
    
    /**
     * @brief Inicia el algoritmo configurando y enlazando el socket ZeroMQ.
     * 
     * Este método configura un socket de tipo REP (respuesta) para recibir mensajes y lo enlaza 
     * a la dirección especificada. Luego inicia el proceso de escucha de mensajes entrantes.
     * 
     * @throws SocketZeroMQException Si ocurre un error al configurar o enlazar el socket.
     */
    @Override
    public void startAlgorithm () throws SocketZeroMQException {
        try {
            // Crear el socket de tipo REP (respuesta) utilizando el contexto de ZeroMQ
            socket = context.createSocket (SocketType.REP);

            // Enlazar el socket a la dirección TCP configurada
            this.socket.bind ("tcp://" + address);

            // Iniciar la escucha de mensajes entrantes en el socket
            startListening ();
        } catch (Exception e) {
            // Registrar el error en los logs
            logger.error ("Error al inicializar el socket REP para el seguidor {}: {}", name, e.getMessage (), e);
            // Lanzar una excepción específica para errores de conexión de ZeroMQ
            throw new SocketZeroMQException("Error al inicializar el socket REP: " + e.getMessage (),
            		                                        SocketZeroMQException.ErrorType.CONNECTION_ERROR);
        }
    }
}

 
