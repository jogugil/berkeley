package org.jogugi.berkeley;

import java.io.IOException;
import java.io.InputStream;
import java.util.HashMap;

import org.jogugi.berkeley.Config.FollowerConfig;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.fasterxml.jackson.databind.ObjectMapper;

/**
 * Clase principal que gestiona la inicialización y ejecución del líder y los seguidores en el sistema.
 * 
 * Esta clase se encarga de cargar la configuración desde un archivo JSON, inicializar los nodos de los seguidores y el líder,
 * y coordinar la ejecución de los algoritmos en el sistema. Además, gestiona las excepciones y asegura que el flujo de trabajo
 * continúe o se detenga dependiendo de la correcta inicialización de los componentes.
 * 
 * Los pasos clave en este proceso son:
 * 1. Cargar la configuración desde un archivo JSON.
 * 2. Crear e inicializar los seguidores y el líder.
 * 3. Iniciar los algoritmos en los seguidores y el líder.
 * 4. Gestionar el tiempo de espera para la sincronización.
 * 5. Cerrar el nodo del líder de forma segura.
 * 
 * @author [Tu nombre]
 * @version 1.0
 * @since [Fecha de creación]
 */
public class Main {
    private static final Logger logger = LoggerFactory.getLogger(Main.class);

    /**
     * Método principal que inicia el proceso de inicialización y ejecución del sistema.
     * 
     * @param args Argumentos de línea de comandos. No se utilizan en esta implementación.
     * @throws IOException Si ocurre un error al leer el archivo de configuración o durante la inicialización de los nodos.
     * @throws InterruptedException Si la ejecución es interrumpida durante los períodos de espera.
     * @throws SocketZeroMQException Si ocurre un error al intentar establecer comunicación a través de ZeroMQ.
     */
    public static void main(String[] args) throws IOException, InterruptedException, SocketZeroMQException {
        // Declarar la variable config fuera del bloque try
        Config config;
 
        // Cargar la configuración desde el archivo JSON
        ObjectMapper mapper = new ObjectMapper();
        try (InputStream input = Main.class.getResourceAsStream("/org/jogugi/berkeley/config.json")) {
            if (input == null) {
                logger.error("Archivo config.json no encontrado en el classpath.");
                throw new IOException("Archivo config.json no encontrado en el classpath.");
            }
            config = mapper.readValue(input, Config.class);
            logger.info("Configuración cargada correctamente desde config.json.");
        }

 
        // Crear el HashMap de seguidores
        HashMap<String, String> followerAddresses = new HashMap<>();
        for (FollowerConfig follower : config.followers) {
            followerAddresses.put(follower.name, follower.address); // Dirección como clave y nombre como valor
        }
        logger.info("HashMap de seguidores: {}", followerAddresses);

        // Mostarmos todas las direcciones y nombres
        followerAddresses.forEach((name, address) -> {
        	logger.info ("Dirección: " + address + ", Nombre: " + name);
        });


        // Crear el líder
        Leader leader = new Leader();
        try {
            leader.initializeNode( config.leader.address, config.leader.name,config.timeout, followerAddresses);
            logger.info("Líder {} inicializado en dirección {}", config.leader.name, config.leader.address);
        } catch (Exception  e) {
            logger.error("Error al inicializar el líder: {}", e.getMessage());
            throw e; // Detener la ejecución en caso de error crítico
        }

        // Crear los seguidores
        for (Config.FollowerConfig followerConfig : config.followers) {
            Follower follower = new Follower();
            try {
                follower.initializeNode(followerConfig.name, followerConfig.address, config.leader.address, config.timeout);
                logger.info("Seguidor {} inicializado en dirección {}", followerConfig.name, followerConfig.address);

                new Thread(() -> {
                    try {
                        logger.info("Iniciando algoritmo para el seguidor {}", followerConfig.name);
                        follower.startAlgorithm();
                    } catch (SocketZeroMQException e) {
                        logger.error("Error al iniciar el algoritmo del seguidor {}: {}", followerConfig.name, e.getMessage());
                    }
                }).start();
            } catch (IOException | SocketZeroMQException e) {
                logger.error("Error al inicializar el seguidor {}: {}", followerConfig.name, e.getMessage());
                // Detener el programa si un seguidor no puede ser inicializado correctamente
                throw new RuntimeException("Error crítico al inicializar el seguidor, interrumpiendo el flujo.", e);
            }
        }

        // Espera 2 segundos para asegurar que los seguidores estén listos
        Thread.sleep(2000);
        logger.info("Esperando 2 segundos para asegurar que los seguidores estén listos.");

        // Iniciar el algoritmo del líder
        try {
            logger.info("Iniciando algoritmo del líder.");
            leader.startAlgorithm();
        } catch (Exception  e) {
            logger.error("Error al iniciar el algoritmo del líder: {}", e.getMessage());
            throw new RuntimeException("Error crítico al iniciar el algoritmo del líder.", e); // Detener si el algoritmo no puede comenzar
        }

        // Esperar que el líder reciba respuestas de los seguidores
        try {
            Thread.sleep(config.timeout + 200); // Dar tiempo adicional al líder
            logger.info("Esperando respuestas de los seguidores.");
        } catch (InterruptedException e) {
            logger.error("Error durante el tiempo de espera para las respuestas de los seguidores: {}", e.getMessage());
            throw new RuntimeException("Error durante el tiempo de espera, interrumpiendo ejecución.", e); // Detener si hay un error en el tiempo de espera
        }

        // Cerrar el nodo del líder
        try {
            leader.close();
            logger.info("Nodo líder cerrado correctamente.");
        } catch (SocketZeroMQException e) {
            logger.error("Error al cerrar el nodo líder: {}", e.getMessage());
            throw new RuntimeException("Error crítico al cerrar el nodo líder.", e); // Detener si no se puede cerrar el líder
        }
    }
}
