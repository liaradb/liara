package application

import (
	"os"

	"github.com/cardboardrobots/config"
)

type configuration struct {
	Port       int    `yaml:"port" config:"PORT"`
	Buffers    int    `yaml:"buffers" config:"BUFFERS"`
	BlockSize  int    `yaml:"blockSize" config:"BLOCK_SIZE"`
	RecordSize int    `yaml:"blockSize" config:"RECORD_SIZE"`
	Directory  string `yaml:"directory" config:"DIRECTORY"`
}

func LoadConfig() (configuration, error) {
	return config.ReadConfigFile[configuration](os.DirFS("."), "config.yml")
}
