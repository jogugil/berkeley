package org.jogugi.berkeley;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * La clase {@code LoggerManager} es responsable de gestionar los loggers
 * para las clases que lo invocan. Utiliza la biblioteca {@link slf4j} para
 * crear instancias de loggers.
 * <p>
 * Esta clase proporciona un único método estático que devuelve el logger
 * correspondiente a la clase que lo invoca.
 * </p>
 * 
 * @author Tu Nombre
 * @version 1.0
 * @since 2024-12-07
 */
public class LoggerManager {

    /**
     * Obtiene el logger correspondiente a la clase que invoca este método.
     * 
     * @param clazz La clase para la cual se desea obtener el logger.
     * @return Un objeto {@link Logger} correspondiente a la clase {@code clazz}.
     * @see Logger
     */
    public static Logger getLogger(Class<?> clazz) {
        return LoggerFactory.getLogger(clazz);
    }
}
