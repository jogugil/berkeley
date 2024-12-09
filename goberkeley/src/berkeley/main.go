package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"encoding/json"
	"berkeley"
)

type Config struct {
	Leader  berkeley.LeaderConfig    `json:"leader"`
	Followers []berkeley.FollowerConfig `json:"followers"`
	Timeout int `json:"timeout"`
}

func main() {
	// Cargar la configuración desde el archivo JSON
	configFile := "/path/to/config.json"
	config, err := loadConfig(configFile)
	if err != nil {
		log.Fatalf("Error al cargar la configuración: %v", err)
	}

	// Mostrar las direcciones de los seguidores
	followerAddresses := make(map[string]string)
	for _, follower := range config.Followers {
		followerAddresses[follower.Name] = follower.Address
		log.Printf("Dirección: %s, Nombre: %s", follower.Address, follower.Name)
	}

	// Crear el líder
	leader := berkeley.Leader{}
	err = leader.InitializeNode(config.Leader.Address, config.Leader.Name, config.Timeout, followerAddresses)
	if err != nil {
		log.Fatalf("Error al inicializar el líder: %v", err)
	}
	log.Printf("Líder %s inicializado en dirección %s", config.Leader.Name, config.Leader.Address)

	// Crear los seguidores
	for _, followerConfig := range config.Followers {
		follower := berkeley.Follower{}
		err = follower.InitializeNode(followerConfig.Name, followerConfig.Address, config.Leader.Address, config.Timeout)
		if err != nil {
			log.Fatalf("Error al inicializar el seguidor %s: %v", followerConfig.Name, err)
		}
		log.Printf("Seguidor %s inicializado en dirección %s", followerConfig.Name, followerConfig.Address)

		go func(followerName string) {
			log.Printf("Iniciando algoritmo para el seguidor %s", followerName)
			err := follower.StartAlgorithm()
			if err != nil {
				log.Printf("Error al iniciar el algoritmo del seguidor %s: %v", followerName, err)
			}
		}(followerConfig.Name)
	}

	// Esperar 2 segundos para asegurar que los seguidores estén listos
	time.Sleep(2 * time.Second)
	log.Println("Esperando 2 segundos para asegurar que los seguidores estén listos.")

	// Iniciar el algoritmo del líder
	log.Println("Iniciando algoritmo del líder.")
	err = leader.StartAlgorithm()
	if err != nil {
		log.Fatalf("Error al iniciar el algoritmo del líder: %v", err)
	}

	// Esperar que el líder reciba respuestas de los seguidores
	time.Sleep(time.Duration(config.Timeout+200) * time.Millisecond)
	log.Println("Esperando respuestas de los seguidores.")

	// Cerrar el nodo del líder
	err = leader.Close()
	if err != nil {
		log.Fatalf("Error al cerrar el nodo líder: %v", err)
	}
	log.Println("Nodo líder cerrado correctamente.")
}

// Cargar la configuración desde el archivo JSON
func loadConfig(filePath string) (Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("no se puede abrir el archivo de configuración: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, fmt.Errorf("error al decodificar el archivo de configuración: %w", err)
	}

	return config, nil
}
