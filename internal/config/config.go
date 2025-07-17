package config

import (
	"os"

	"github.com/nerfthisdev/backend-test-task/internal/database"
	"gopkg.in/yaml.v3"
)

type Server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Config struct {
	Server   Server          `yaml:"server"`
	Database database.Config `yaml:"database"`
}

func GetConfiguration(configPath string, cfg any) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return err
	}

	return nil
}
