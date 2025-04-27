package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Server struct {
	Port string `yaml:"port"`
}

type ServerConfig struct {
	Server *Server `yaml:"server"`
}

func LoadConfig(path string) (*ServerConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла: %v", err)
	}
	var config ServerConfig
	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("ошибка декодирования yaml: %v", err)
	}
	return &config, nil
}
