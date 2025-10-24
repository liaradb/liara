package application

import (
	"os"

	"github.com/cardboardrobots/config"
)

// TODO: Fix camel case
type configuration struct {
	Port      int    `yml:"port" config:"PORT"`
	Buffers   int    `yml:"buffers" config:"BUFFERS"`
	BlockSize int    `yml:"blocksize" config:"BLOCK_SIZE"`
	Directory string `yml:"directory" config:"DIRECTORY"`
}

func LoadConfig() (configuration, error) {
	return config.ReadConfigFile[configuration](os.DirFS("."), "config.yml")
}
