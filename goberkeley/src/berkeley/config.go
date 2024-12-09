package berkeley

import (
	"encoding/json"
	"fmt"
	"os"
)

func Config () {
	// Abrir el archivo JSON
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error abriendo el archivo:", err)
		return
	}
	defer file.Close()

	// Decodificar el JSON en la estructura Config
	var config berkeley.Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		fmt.Println("Error decodificando JSON:", err)
		return
	}

	// Imprimir la configuración cargada
	fmt.Printf("Líder: %s (%s)\n", config.Leader.Name, config.Leader.Address)
	for _, follower := range config.Followers {
		fmt.Printf("Seguidor: %s (%s)\n", follower.Name, follower.Address)
	}
	fmt.Printf("Timeout: %d ms\n", config.Timeout)
}

