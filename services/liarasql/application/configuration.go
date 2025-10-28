package application

import (
	"os"

	"github.com/cardboardrobots/config"
)

type configuration struct {
	Port          int    `yml:"port" config:"PORT"`
	PostgresDbUri string `yaml:"postgresDbUri" config:"POSTGRES_DB_URI"`
	SqliteDbUri   string `yaml:"sqliteDbUri" config:"SQLITE_DB_URI"`
}

func LoadConfig() (configuration, error) {
	return config.ReadConfigFile[configuration](os.DirFS("."), "config.yml")
}
