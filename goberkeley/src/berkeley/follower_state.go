package berkeley

// FollowerState representa los posibles estados de un seguidor durante el proceso de actualización de la hora del sistema.
type FollowerState string

const (
	// RESPONDED indica que el seguidor ha respondido correctamente a la solicitud del líder.
	Responded FollowerState = "RESPONDED"

	// NO_RESPONSE indica que el seguidor no ha respondido a la solicitud dentro del tiempo esperado.
	NoResponse FollowerState = "NO_RESPONSE"

	// CONNECTION_ERROR indica que hubo un error de conexión al intentar comunicarse con el seguidor.
	ConnectionError FollowerState = "CONNECTION_ERROR"

	// REQUEST_NOT_SENT indica que la solicitud de actualización de la hora no fue enviada debido a algún error.
	RequestNotSent FollowerState = "REQUEST_NOT_SENT"

	// REQUEST_DELTA_SENT indica que la solicitud de actualización de la hora fue enviada al seguidor para actualizar su hora.
	RequestDeltaSent FollowerState = "REQUEST_DELTA_SENT"

	// TIME_ERROR_SENT_UPDATE indica que la solicitud de actualización de la hora fue enviada pero hubo un error en ese envío.
	TimeErrorSentUpdate FollowerState = "TIME_ERROR_SENT_UPDATE"

	// TIME_UPDATED indica que la hora del sistema del seguidor se actualizó correctamente.
	TimeUpdated FollowerState = "TIME_UPDATED"
<<<<<<< HEAD

	// ERROR_CLOSE indica que  no pude enviar el cierre del socket en el seguidor.
	ErrorClose FollowerState = "ERROR_CLOSE"

	// Ok_CLOSE indica que   pude  cerrar el socket en el seguidor.
	OkClose FollowerState = "Ok_CLOSE"
)
=======
)
>>>>>>> 4df5b3de8e63625e208a4e5d62139f2f2b1cb612
