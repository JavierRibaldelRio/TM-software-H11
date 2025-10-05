package config

import (
	"encoding/json"
	"os"
	"ribal-backend-receiver/logger"
)

type Config struct {
	RingBufferSize int
	BufferSize     int
}

// default settings

var Default = Config{RingBufferSize: 10, BufferSize: 4096}

// loads configuration
func Load(path string) (Config, error) {
	cfg := Default

	file, err := os.Open(path)
	if err != nil {
		// Si el archivo no existe, devolvemos la configuración por defecto
		if os.IsNotExist(err) {
			logger.Info("Config file not found, using defaults")
			return cfg, nil
		}
		return cfg, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil

}
