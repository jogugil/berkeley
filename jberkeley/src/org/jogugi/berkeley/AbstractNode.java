package org.jogugi.berkeley;
 

import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ExecutionException;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.zeromq.SocketType;
import org.zeromq.ZContext;
import org.zeromq.ZMQ;

/**
 * @brief Clase abstracta que implementa la interfaz `INode`.
 *
 * La clase `AbstractNode` proporciona la implementación base para un nodo en el sistema.
 * Maneja el contexto de ZeroMQ, los sockets de comunicación y contiene las operaciones comunes
 * para el envío y recepción de mensajes, además de almacenar la información del nodo como
 * su nombre, dirección y tiempos de espera.
 */
public abstract class AbstractNode implements INode {

    /**
     * @brief Contexto de ZeroMQ utilizado para la creación de sockets.
     */
    protected ZContext context = null;

    /**
     * @brief Socket REP del nodo utilizado para recibir mensajes.
     */
    protected ZMQ.Socket socket = null;

    /**
     * @brief Nombre del nodo.
     */
    protected String name;

    /**
     * @brief Dirección del nodo en formato `host:puerto`.
     */
    protected String address;

    /**
     * @brief Tiempo de espera en milisegundos para operaciones de socket.
     */
    protected int timeout;

    /**
     * @brief Mapa que contiene las direcciones de otros nodos en el sistema.
     */
    protected Map<String, String> nodeAddresses = null;

    /**
     * @brief Lista de sockets REQ utilizados para enviar solicitudes a otros nodos.
     */
    protected List<ZMQ.Socket> listSocketsREQ;

    /**
     * @brief Logger utilizado para registrar información, advertencias y errores.
     */
    private static final Logger logger = LoggerFactory.getLogger(AbstractNode.class);
    
 // Constructor común
    @Override
    /**
     * @brief Inicializa el nodo con el nombre, dirección y tiempo de espera.
     *
     * Este constructor inicializa un nodo con los parámetros básicos: nombre del nodo,
     * dirección del nodo y el tiempo de espera. También crea el contexto de ZeroMQ para
     * la comunicación.
     * 
     * @param name Nombre del nodo.
     * @param address Dirección del nodo en formato `host:puerto`.
     * @param timeout Tiempo de espera para las operaciones de socket.
     */
    public void initializeNode(String name, String address, int timeout) {
        this.name    = name;
        this.address = address;
        this.timeout = timeout;
        this.context = new ZContext(); // Crear el contexto ZeroMQ
        logger.info("Nodo {} inicializado en dirección {}", name, address);
    }

    // Método común que incluye la lista de nodos. Posible para el líder
    @Override
    /**
     * @brief Inicializa el nodo con el nombre, dirección, tiempo de espera y la lista de nodos.
     *
     * Este constructor inicializa un nodo con los parámetros básicos: nombre del nodo,
     * dirección del nodo, tiempo de espera y la lista de direcciones de otros nodos en el
     * sistema. También crea el contexto de ZeroMQ para la comunicación.
     * 
     * @param name Nombre del nodo.
     * @param address Dirección del nodo en formato `host:puerto`.
     * @param timeout Tiempo de espera para las operaciones de socket.
     * @param listNodeAddresses Lista de direcciones de otros nodos en el sistema.
     */
    public void initializeNode(String name, String address, int timeout, HashMap<String, String> listNodeAddresses) {
        this.name          = name;
        this.address       = address;
        this.timeout       = timeout;
        this.context       = new ZContext(); // Crear el contexto ZeroMQ
        this.nodeAddresses = listNodeAddresses;
        logger.info("Nodo {} con lista de nodos {} inicializado en dirección {}", name, listNodeAddresses, address);
    }
    @Override
    /**
     * @brief Envía un mensaje de manera síncrona y espera una respuesta.
     *
     * Este método crea un socket REQ para enviar un mensaje a una dirección específica y espera
     * una respuesta del servidor. Si no se recibe respuesta dentro del tiempo de espera o si ocurre
     * un error durante el envío/recepción, se lanza una excepción.
     * 
     * @param address La dirección a la que se enviará el mensaje en formato `host:puerto`.
     * @param message El mensaje que se va a enviar.
     * 
     * @return La respuesta recibida del servidor.
     * 
     * @throws SocketZeroMQException Si ocurre un error al enviar o recibir el mensaje.
     */
    public String sendMessageSync(String address, String message) throws SocketZeroMQException {
        try (ZMQ.Socket clientSocket = context.createSocket(SocketType.REQ)) {
            // Conectar el socket al servidor en la dirección indicada
            clientSocket.connect(address);

            // Enviar el mensaje
            clientSocket.send(message);
            logger.debug("Mensaje enviado: {}", message);

            // Esperar la respuesta del servidor
            String reply = clientSocket.recvStr();
            if (reply == null) {
                // Si no se recibe respuesta, lanzar una excepción
                throw new SocketZeroMQException("No se recibió respuesta del socket", SocketZeroMQException.ErrorType.RECEIVE_ERROR);
            }

            // Mostrar la respuesta recibida
            logger.debug("Respuesta recibida: {}", reply);
            return reply;
        } catch (Exception e) {
            // Si ocurre un error, lanzar una excepción con el mensaje de error
            throw new SocketZeroMQException("Error al enviar/recibir mensaje: " + e.getMessage(), SocketZeroMQException.ErrorType.SEND_ERROR);
        }
    }
    @Override
    /**
     * @brief Envía un mensaje de manera asíncrona.
     *
     * Este método crea un socket PUSH para enviar un mensaje de manera asíncrona a la dirección
     * especificada. No espera una respuesta, sino que simplemente envía el mensaje y termina la operación.
     * 
     * @param address La dirección a la que se enviará el mensaje en formato `host:puerto`.
     * @param message El mensaje que se va a enviar.
     * 
     * @throws SocketZeroMQException Si ocurre un error al enviar el mensaje.
     */
    public void sendMessageAsync(String address, String message) throws SocketZeroMQException {
        try (ZMQ.Socket clientSocket = context.createSocket(SocketType.PUSH)) {
            // Conectar el socket al servidor en la dirección indicada
            clientSocket.connect(address);

            // Enviar el mensaje de manera asíncrona
            clientSocket.send(message);
            logger.debug("Mensaje enviado de forma asíncrona: {}", message);
        } catch (Exception e) {
            // Si ocurre un error, lanzar una excepción con el mensaje de error
            throw new SocketZeroMQException("Error al enviar mensaje asíncrono: " + e.getMessage(), SocketZeroMQException.ErrorType.SEND_ERROR);
        }
    }

    @Override
    /**
     * @brief Inicia el proceso de escucha de mensajes en el nodo.
     *
     * Este método pone al nodo en un estado en el que puede escuchar mensajes entrantes en su socket.
     * El nodo recibe mensajes de forma continua en un hilo separado, los procesa y, opcionalmente, responde
     * con un mensaje. Si recibe un mensaje que indica que el nodo debe cerrarse (por ejemplo, un mensaje con
     * la operación `"CLOSE"`), el nodo cierra su socket y termina el proceso de escucha.
     * 
     * @throws SocketZeroMQException Si el socket no está inicializado o ocurre un error en el proceso de escucha.
     */
    public void startListening() throws SocketZeroMQException {
        // Verificar si el socket está inicializado
        if (this.socket == null) {
            throw new SocketZeroMQException("Socket no inicializado en el nodo: " + this.name, SocketZeroMQException.ErrorType.CONNECTION_ERROR);
        }

        // Iniciar un nuevo hilo para escuchar mensajes
        new Thread(() -> {
            logger.info("Iniciando escucha en {}...", name);

            // Bucle principal para recibir mensajes mientras el hilo no sea interrumpido
            while (!Thread.currentThread().isInterrupted()) {
                try {
                    // Recibir un mensaje del socket
                    String receivedMessage = socket.recvStr();
                    
                    // Verificar si el mensaje recibido es nulo
                    if (receivedMessage == null) {
                        throw new SocketZeroMQException("Error al recibir mensaje", SocketZeroMQException.ErrorType.RECEIVE_ERROR);
                    }

                    // Registrar el mensaje recibido
                    logger.debug("Mensaje recibido en {}: {}", name, receivedMessage);

                    // Procesar el mensaje recibido y generar una respuesta
                    String response = handleProcess(receivedMessage); 

                    // Enviar la respuesta al remitente
                    socket.send(response);
                    logger.debug("Mensaje enviado por {}: {}", name, response);

                    // Si el mensaje contiene la operación "CLOSE", cerrar el socket
                    if (response.contains("\"operation\":\"CLOSE\"")) {
                        logger.info("Nodo {} llamando a close.", name);
                        close(); // Cerrar el socket
                        break;   // Salir del bucle de escucha
                    } 
                } catch (SocketZeroMQException e) {
                    // Manejar errores en la escucha de mensajes
                    logger.error("Error en la escucha de mensajes: {}", e.getMessage());
                    break; // Terminar el hilo en caso de un error grave
                }
            }
        }).start(); // Iniciar el hilo para escuchar los mensajes
    }

    @Override
    /**
     * @brief Cierra los recursos utilizados por el nodo, incluyendo el socket y el contexto.
     *
     * Este método cierra el socket asociado al nodo y también el contexto ZeroMQ, liberando así los
     * recursos utilizados por el nodo. Una vez cerrados estos recursos, el nodo deja de poder enviar
     * o recibir mensajes. El proceso de cierre es registrado en los logs del sistema.
     *
     * @throws SocketZeroMQException Si ocurre un error al intentar cerrar el socket o el contexto.
     */
    public void close() throws SocketZeroMQException {
        // Verificar si el socket está inicializado y cerrarlo
        if (socket != null) {
            socket.close();
            logger.info("Socket cerrado para el nodo {}", name);
        }
        
        // Verificar si el contexto está inicializado y cerrarlo
        if (context != null) {
            context.close();
            logger.info("Contexto cerrado para el nodo {}", name);
        }
        
        // Registrar el cierre del nodo
        logger.info("Nodo {} cerrado.", name);
    }
    @Override
    /**
     * @brief Método para iniciar el algoritmo del nodo.
     *
     * Este método se sobrecarga por el tipo de nodo (por ejemplo, nodo líder o seguidor) y debe ser
     * implementado en las subclases específicas. Si no se implementa en una subclase, lanza un error.
     *
     * @throws SocketZeroMQException Si ocurre un error en la comunicación con ZeroMQ.
     * @throws InterruptedException Si el hilo es interrumpido mientras se ejecuta el algoritmo.
     * @throws ExecutionException Si ocurre un error durante la ejecución del algoritmo.
     */
    public void startAlgorithm() throws SocketZeroMQException, InterruptedException, ExecutionException {
        // Para sobrecargar por el tipo de nodo
        logger.error("Función no implementada. Debe ser sobrecargada e implementada en la subclase.");
    }

    /**
     * @brief Método abstracto para manejar el procesamiento de los mensajes recibidos.
     *
     * Este método debe ser implementado en las subclases específicas para definir el comportamiento
     * que debe seguir el nodo al recibir un mensaje. El mensaje es procesado y se devuelve una respuesta
     * basada en la lógica del nodo.
     *
     * @param message El mensaje recibido que debe ser procesado.
     * @return La respuesta procesada al mensaje.
     * @throws SocketZeroMQException Si ocurre un error al procesar el mensaje.
     */
    protected abstract String handleProcess(String message) throws SocketZeroMQException;
}
