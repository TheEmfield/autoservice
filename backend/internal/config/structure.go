package config

type Config struct {
	Database Database `yaml:"database" env-prefix:"POSTGRES_"`
	App      App      `yaml:"app"`
}

type Database struct {
	Host     string `yaml:"host"     env:"HOST"     env-default:"localhost"`
	Port     string `yaml:"port"     env:"PORT"     env-default:"5432"`
	User     string `yaml:"user"     env:"USER"     env-default:"postgres"`
	Password string `yaml:"password" env:"PASSWORD" env-default:"postgres"`
	DB       string `yaml:"db"       env:"DB"       env-default:"coursework"`
	SSLMode  string `yaml:"ssl_mode" env:"SSL_MODE" env-default:"disable"`
}

type App struct {
	Port string `yaml:"port" env:"APP_PORT" env-default:"8080"`
}
