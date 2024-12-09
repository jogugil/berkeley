package berkeley

import (
	"encoding/json"
	"fmt"
	"os"
<<<<<<< HEAD
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
=======
)

func Config () {
	// Abrir el archivo JSON
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error abriendo el archivo:", err)
		return
>>>>>>> 4df5b3de8e63625e208a4e5d62139f2f2b1cb612
	}
	defer file.Close()

	// Decodificar el JSON en la estructura Config
<<<<<<< HEAD
	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		fmt.Println("Error decodificando JSON:", err)
		return nil
=======
	var config berkeley.Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		fmt.Println("Error decodificando JSON:", err)
		return
>>>>>>> 4df5b3de8e63625e208a4e5d62139f2f2b1cb612
	}

	// Imprimir la configuración cargada
	fmt.Printf("Líder: %s (%s)\n", config.Leader.Name, config.Leader.Address)
	for _, follower := range config.Followers {
		fmt.Printf("Seguidor: %s (%s)\n", follower.Name, follower.Address)
	}
	fmt.Printf("Timeout: %d ms\n", config.Timeout)
<<<<<<< HEAD

	return &config
}
=======
}

>>>>>>> 4df5b3de8e63625e208a4e5d62139f2f2b1cb612
