package org.jogugi.berkeley;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.List;

/**
 * @brief Clase de configuración para los nodos en el sistema.
 *
 * Esta clase proporciona la configuración necesaria para inicializar los nodos en el sistema,
 * incluyendo la configuración del líder y los seguidores, así como el tiempo de espera global.
 */
public class Config {

    /**
     * @brief Configuración del nodo líder.
     *
     * Esta clase define las propiedades del nodo líder, incluyendo su nombre y dirección.
     */
    public static class LeaderConfig {
        
        /**
         * @brief Nombre del nodo líder.
         *
         * Este campo almacena el nombre del nodo líder en el sistema.
         */
        public String name;
        
        /**
         * @brief Dirección del nodo líder.
         *
         * Este campo almacena la dirección del nodo líder, que se usará para establecer conexiones.
         */
        public String address;
    }

    /**
     * @brief Configuración del nodo seguidor.
     *
     * Esta clase define las propiedades del nodo seguidor, incluyendo su nombre y dirección.
     */
    public static class FollowerConfig {
        
        /**
         * @brief Nombre del nodo seguidor.
         *
         * Este campo almacena el nombre de cada nodo seguidor en el sistema.
         */
        public String name;
        
        /**
         * @brief Dirección del nodo seguidor.
         *
         * Este campo almacena la dirección de cada nodo seguidor, que se usará para establecer conexiones.
         */
        public String address;
    }

    /**
     * @brief Nodo líder del sistema.
     *
     * Esta propiedad contiene la configuración del nodo líder, que incluye su nombre y dirección.
     */
    @JsonProperty("leader")
    public LeaderConfig leader;

    /**
     * @brief Lista de nodos seguidores.
     *
     * Esta propiedad contiene una lista de configuraciones de nodos seguidores, cada una con su nombre
     * y dirección. Los seguidores se conectan al líder.
     */
    @JsonProperty("followers")
    public List<FollowerConfig> followers;

    /**
     * @brief Tiempo de espera global para todos los nodos.
     *
     * Este campo define el tiempo máximo que los nodos esperarán por una respuesta antes de considerar
     * que ocurrió un error de tiempo de espera.
     */
    @JsonProperty("timeout")
    public int timeout;
}