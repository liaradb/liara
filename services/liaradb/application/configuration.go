package application

import (
	"os"

	"github.com/cardboardrobots/config"
)

type configuration struct {
	Port int `yml:"port" config:"PORT"`
}

func LoadConfig() (configuration, error) {
	return config.ReadConfigFile[configuration](os.DirFS("."), "config.yml")
}
