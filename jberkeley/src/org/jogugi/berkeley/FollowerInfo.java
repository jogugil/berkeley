package org.jogugi.berkeley;

import java.util.Date;

/**
 * @brief Clase que representa información detallada de un nodo seguidor en el sistema.
 *
 * La clase `FollowerInfo` encapsula información sobre un seguidor, incluyendo su nombre, dirección, 
 * tiempos relacionados con la comunicación, diferencias horarias y estado actual.
 */
public class FollowerInfo {
    
    /** @brief Nombre del seguidor. */
    String name;

    /** @brief Dirección del seguidor en formato `host:puerto`. */
    String address;

    /** @brief Estado actual del seguidor. */
    FollowerState state;

    /** @brief Fecha y hora local del seguidor en formato `Date`. */
    Date dateFollower;

    /** @brief Hora local del seguidor en milisegundos desde la época UNIX. */
    long localTime;

    /** @brief Tiempo de comunicación entre el líder y el seguidor en milisegundos. */
    long communicationTime;

    /** @brief Tiempo de ida y vuelta (triptime) estimado para el seguidor. */
    long tripTime;

    /** @brief Diferencia de tiempo calculada entre el líder y el seguidor. */
    long diffTime;

    /** @brief Diferencia global de tiempo (delta) aplicada al seguidor. */
    long delta;

    /**
     * @brief Constructor de la clase `FollowerInfo`.
     *
     * Inicializa un objeto `FollowerInfo` con los valores proporcionados.
     * 
     * @param address Dirección del seguidor.
     * @param name Nombre del seguidor.
     * @param localTime Hora local del seguidor en milisegundos desde la época UNIX.
     * @param communicationTime Tiempo de comunicación en milisegundos.
     * @param nowTime Hora actual del líder en milisegundos desde la época UNIX.
     */
    public FollowerInfo(String address, String name, long localTime, long communicationTime, long nowTime) {
        this.name = name;
        this.address = address;
        this.localTime = localTime;
        this.dateFollower = new Date(localTime);
        this.communicationTime = communicationTime;
        this.tripTime = communicationTime / 2;
        this.diffTime = (localTime + tripTime) - nowTime;
        this.state = FollowerState.REQUEST_NOT_SENT; // Estado inicial por defecto
    }

    /** @brief Obtiene la dirección del seguidor. */
    public String getAdrressFollower() {
        return address;
    }

    /** @brief Obtiene el nombre del seguidor. */
    public String getName() {
        return name;
    }

    /** @brief Obtiene la fecha y hora local del seguidor. */
    public Date getDateFollower() {
        return dateFollower;
    }

    /** @brief Obtiene la diferencia de tiempo entre el líder y el seguidor. */
    public long getDiffTime() {
        return diffTime;
    }

    /** @brief Obtiene la hora local del seguidor. */
    public long getLocalTime() {
        return localTime;
    }

    /** @brief Obtiene el tiempo de comunicación entre el líder y el seguidor. */
    public long getCommunicationTime() {
        return communicationTime;
    }

    /** @brief Obtiene la diferencia global de tiempo (delta). */
    public long getDelta() {
        return delta;
    }

    /**
     * @brief Establece la diferencia global de tiempo (delta) para el seguidor.
     * @param delta Valor de la diferencia global de tiempo en milisegundos.
     */
    public void setDelta(long delta) {
        this.delta = delta;
    }

    /** @brief Obtiene el estado actual del seguidor. */
    public FollowerState getState() {
        return state;
    }

    /**
     * @brief Establece el estado del seguidor.
     * @param state Nuevo estado del seguidor.
     */
    public void setState(FollowerState state) {
        this.state = state;
    }

    /**
     * @brief Representación en cadena del objeto `FollowerInfo`.
     * 
     * @return Una descripción detallada de los atributos del seguidor.
     */
    @Override
    public String toString() {
        return "Nombre: " + name + ", " + "State: " + state + ", " +
               "Hora local: " + (localTime != -1 ? localTime : "No respondió") + ", " +
               "Hora del seguidor: " + (dateFollower != null ? dateFollower.toString() : "No disponible") + ", " +
               "Tiempo de comunicación: " + communicationTime + " ms, " +
               "Triptime con el seguidor: " + tripTime + " ms, " +
               "Diferencia de tiempo para el seguidor: " + diffTime + " ms, " +
               "Diferencia de tiempo global (delta): " + delta + " ms, " +
               "Dirección del seguidor: " + (address != null ? address : "No disponible");
    }
}