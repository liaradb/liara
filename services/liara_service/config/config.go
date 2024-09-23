package config

import (
	"os"

	"github.com/cardboardrobots/config"
)

type Config struct {
	Port          int    `yml:"port" config:"PORT"`
	PostgresDbUri string `yaml:"postgresDbUri" config:"POSTGRES_DB_URI"`
	SqliteDbUri   string `yaml:"sqliteDbUri" config:"SQLITE_DB_URI"`
}

func LoadConfig() (Config, error) {
	return config.ReadConfigFile[Config](os.DirFS("."), "config/config.yml")
}
