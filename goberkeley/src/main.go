package main

import (
	"fmt"
	"goberkeley/berkeley"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Cargar la configuración desde el archivo JSON
	configFile := "config.json"
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
	var leader *berkeley.Leader
	leader, err_leader := berkeley.InitializeLeaderNode(config.Leader.Address, config.Leader.Name, config.Timeout, followerAddresses)
	if err_leader != nil {
		// Maneja el error adecuadamente
		fmt.Println("Error al inicializar el nodo líder:", err_leader)
		return // O maneja el error de otra manera
	}
	log.Printf("Líder %s inicializado en dirección %s", config.Leader.Name, config.Leader.Address)

	// Crear los seguidores
	for _, followerConfig := range config.Followers {
		var follower *berkeley.Follower
		follower, err := berkeley.InitializeFollowerNode(followerConfig.Name, followerConfig.Address, config.Leader.Address, config.Timeout)

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
	leader.StartAlgorithm()

	// Esperar que el líder reciba respuestas de los seguidores
	time.Sleep(time.Duration(config.Timeout+200) * time.Millisecond)
	log.Println("Esperando respuestas de los seguidores.")

	// Cerrar el nodo del líder
	leader.Close()

	log.Println("Nodo líder cerrado correctamente.")
}

// Cargar la configuración desde el archivo JSON
func loadConfig(filePath string) (config *berkeley.Config, error *string) {

	config = berkeley.LoadConfig(filePath)

	return config, nil
}
