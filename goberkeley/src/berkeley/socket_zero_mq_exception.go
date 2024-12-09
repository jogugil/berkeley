package berkeley

import (
	"fmt"
)

// ErrorType enum para los diferentes tipos de errores posibles en ZeroMQ.
type ErrorType int

const (
	// CONNECTION_ERROR Error relacionado con la conexión de ZeroMQ
	CONNECTION_ERROR ErrorType = iota
	// SEND_ERROR Error al enviar un mensaje a través del socket
	SEND_ERROR
	// RECEIVE_ERROR Error al recibir un mensaje del socket
	RECEIVE_ERROR
)

// SocketZeroMQException estructura que encapsula un error personalizado para ZeroMQ.
type SocketZeroMQException struct {
	Message   string    // Mensaje de error
	ErrorType ErrorType // Tipo de error
}

// NewSocketZeroMQException crea una nueva instancia de SocketZeroMQException.
func NewSocketZeroMQException(message string, errorType ErrorType) *SocketZeroMQException {
	return &SocketZeroMQException{
		Message:   message,
		ErrorType: errorType,
	}
}

// Error implementación de la interfaz error para SocketZeroMQException.
func (e *SocketZeroMQException) Error() string {
	return fmt.Sprintf("Error: %s, Tipo de error: %v", e.Message, e.ErrorType)
}

// GetErrorType obtiene el tipo de error.
func (e *SocketZeroMQException) GetErrorType() ErrorType {
	return e.ErrorType
}
