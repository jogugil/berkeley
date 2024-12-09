package org.jogugi.berkeley;

/**
 * @brief Excepción personalizada para errores en la operación de sockets de ZeroMQ.
 *
 * Esta clase extiende de `Exception` y se utiliza para manejar errores específicos que ocurren
 * durante las operaciones de sockets utilizando ZeroMQ. Los tipos de errores cubiertos son
 * relacionados con la conexión, el envío y la recepción de mensajes.
 */
public class SocketZeroMQException extends Exception {

    /**
	 * 
	 */
	private static final long serialVersionUID = 1L;

	/**
     * @brief Enumeración de los tipos de errores posibles.
     *
     * Esta enumeración se utiliza para clasificar los diferentes tipos de errores que pueden ocurrir
     * durante las operaciones con sockets ZeroMQ.
     */
    public enum ErrorType {
        CONNECTION_ERROR, /**< Error relacionado con la conexión de ZeroMQ. */
        SEND_ERROR,       /**< Error al enviar un mensaje a través del socket. */
        RECEIVE_ERROR     /**< Error al recibir un mensaje del socket. */
    }

    private ErrorType errorType; /**< Tipo de error que ocurrió. */

    /**
     * @brief Constructor de la excepción.
     *
     * @param message El mensaje de error que describe el problema ocurrido.
     * @param errorType El tipo de error que ocurrió, basado en la enumeración `ErrorType`.
     */
    public SocketZeroMQException(String message, ErrorType errorType) {
        super(message); // Llamada al constructor de la clase base (Exception)
        this.errorType = errorType;
    }

    /**
     * @brief Obtiene el tipo de error asociado a esta excepción.
     *
     * @return El tipo de error que se produjo durante la operación de ZeroMQ.
     */
    public ErrorType getErrorType() {
        return errorType;
    }
}