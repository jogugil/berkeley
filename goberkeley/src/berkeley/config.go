package berkeley

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type LeaderConfig struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type FollowerConfig struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Config struct {
	Leader    LeaderConfig     `json:"leader"`
	Followers []FollowerConfig `json:"followers"`
	Timeout   time.Duration    `json:"timeout"`
}

func LoadConfig(filepath string) *Config {
	// Abrir el archivo JSON
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error abriendo el archivo:", err)
		return nil
	}
	defer file.Close()

	// Decodificar el JSON en la estructura Config
	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		fmt.Println("Error decodificando JSON:", err)
		return nil
	}

	// Imprimir la configuración cargada
	fmt.Printf("Líder: %s (%s)\n", config.Leader.Name, config.Leader.Address)
	for _, follower := range config.Followers {
		fmt.Printf("Seguidor: %s (%s)\n", follower.Name, follower.Address)
	}
	fmt.Printf("Timeout: %d ms\n", config.Timeout)

	return &config
}
