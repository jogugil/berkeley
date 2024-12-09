package org.jogugi.berkeley;

/**
 * Enumeración que define los posibles estados de un seguidor durante el proceso de actualización de la hora del sistema.
 * Cada estado representa una etapa en la interacción entre el líder y el seguidor.
 */
public enum FollowerState {
    
    /**
     * El seguidor ha respondido correctamente a la solicitud del líder.
     */
    RESPONDED,

    /**
     * El seguidor no ha respondido a la solicitud dentro del tiempo esperado.
     */
    NO_RESPONSE,

    /**
     * Hubo un error de conexión al intentar comunicarse con el seguidor.
     */
    CONNECTION_ERROR,

    /**
     * La solicitud de actualización de la hora no fue enviada debido a algún error.
     */
    REQUEST_NOT_SENT,

    /**
     * La solicitud de actualización de la hora fue enviada al seguidor para actualizar su hora.
     */
    REQUEST_DELTA_SENT,
    /**
     * La solicitud de actualización de la hora fue enviada pero hubo un error en ese envío.
     */
    TIME_ERROR_SENT_UPDATE,
    /**
     * La hora del sistema del seguidor se actualizó correctamente.
     */
    TIME_UPDATED
}