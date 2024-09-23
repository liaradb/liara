package config

import (
	"os"

	"github.com/cardboardrobots/config"
)

type Config struct {
	Port int `yml:"port" config:"PORT"`
}

func LoadConfig() (Config, error) {
	return config.ReadConfigFile[Config](os.DirFS("."), "config/config.yml")
}
