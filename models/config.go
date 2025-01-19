package models

type DatabaseConfig struct {
	Host     string `mapstructure:"host"` //"localhost"
	Port     string `mapstructure:"port"` //= 5432
	User     string //= "postgres"
	Password string
	Dbname   string
}

type Config struct {
	DBConfig DatabaseConfig
}
