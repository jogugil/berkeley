package org.jogugi.berkeley;

import java.util.HashMap;
import java.util.concurrent.ExecutionException;

/**
 * @brief Interfaz que define las operaciones básicas para un nodo en el sistema.
 *
 * La interfaz `INode` establece los métodos que deben implementarse para inicializar un nodo,
 * enviar mensajes (sincrónicos y asincrónicos), iniciar el algoritmo del nodo y cerrar sus recursos.
 */
public interface INode {

    /**
     * @brief Inicializa el contexto y el socket para el nodo.
     *
     * @param name Nombre del nodo.
     * @param address Dirección del nodo en formato `host:puerto`.
     * @param timeout Tiempo de espera en milisegundos.
     */
    void initializeNode(String name, String address, int timeout);

    /**
     * @brief Inicializa el contexto y el socket para el nodo, incluyendo una lista de direcciones de otros nodos.
     *
     * @param name Nombre del nodo.
     * @param address Dirección del nodo en formato `host:puerto`.
     * @param timeout Tiempo de espera en milisegundos.
     * @param listNodeAddresses Mapa con las direcciones de otros nodos.
     */
    void initializeNode(String name, String address, int timeout, HashMap<String, String> listNodeAddresses);

    /**
     * @brief Envía un mensaje síncrono a un nodo y espera su respuesta.
     *
     * @param address Dirección del nodo receptor en formato `host:puerto`.
     * @param message Mensaje a enviar.
     * @return Respuesta recibida del nodo.
     * @throws SocketZeroMQException Si ocurre un error durante la comunicación.
     */
    String sendMessageSync(String address, String message) throws SocketZeroMQException;

    /**
     * @brief Envía un mensaje de manera asincrónica a un nodo.
     *
     * @param address Dirección del nodo receptor en formato `host:puerto`.
     * @param message Mensaje a enviar.
     * @throws SocketZeroMQException Si ocurre un error durante la comunicación.
     */
    void sendMessageAsync(String address, String message) throws SocketZeroMQException;

    /**
     * @brief Inicia el nodo en modo escucha, esperando recibir mensajes.
     *
     * @throws SocketZeroMQException Si ocurre un error al iniciar el socket.
     */
    void startListening() throws SocketZeroMQException;

    /**
     * @brief Inicia un algoritmo especial para el nodo (puede ser utilizado para procesos específicos como el cálculo del tiempo o sincronización).
     *
     * @throws SocketZeroMQException Si ocurre un error durante la ejecución del algoritmo.
     * @throws InterruptedException Si el hilo es interrumpido durante la ejecución.
     * @throws ExecutionException Si ocurre un error en la ejecución del algoritmo.
     */
    void startAlgorithm() throws SocketZeroMQException, InterruptedException, ExecutionException;

    /**
     * @brief Cierra los recursos del nodo, liberando cualquier conexión o socket abierto.
     *
     * @throws SocketZeroMQException Si ocurre un error al cerrar el nodo.
     */
    void close() throws SocketZeroMQException;
}