package config

import (
	"os"
	"time"
)

// MongoDBConfig almacena la configuración específica para MongoDB
type MongoDBConfig struct {
	MongoURL                      string
	DatabaseName                  string
	ReservationsCollection        string
	RoutesCollection              string
	BaggageReservationsCollection string
	BaggageTypesCollection        string
	ServerPort                    string
	MongoTimeout                  time.Duration
}

// MySQLConfig almacena la configuración específica para MySQL
type MySQLConfig struct {
	Username     string
	Password     string
	Host         string
	Port         string
	DatabaseName string
}

// Config almacena la configuración global del programa
type Config struct {
	MongoDB    MongoDBConfig
	MySQL      MySQLConfig
	ServerPort string
	UsingMongo bool
}

// NewConfig crea y retorna una nueva instancia de Config con los valores proporcionados
func NewConfig() *Config {
	return &Config{
		MongoDB: MongoDBConfig{
			MongoURL:                      getEnv("MONGO_URL", "mongodb://localhost:27017"),
			DatabaseName:                  getEnv("DATABASE_NAME", "venta-de-pasajes"),
			MongoTimeout:                  getEnvDuration("MONGO_TIMEOUT", 10*time.Second),
			ReservationsCollection:        getEnv("RESERVATIONS_COLLECTION", "reservations"),
			RoutesCollection:              getEnv("ROUTES_COLLECTION", "routes"),
			BaggageReservationsCollection: getEnv("BAGGAGES_COLLECTION", "baggageReservations"),
			BaggageTypesCollection:        getEnv("BAGGAGE_TYPES_COLLECTION", "baggageTypes"),
		},
		MySQL: MySQLConfig{
			Username:     getEnv("MYSQL_USERNAME", "root"),
			Password:     getEnv("MYSQL_PASSWORD", ""),
			Host:         getEnv("MYSQL_HOST", "localhost"),
			Port:         getEnv("MYSQL_PORT", "3306"),
			DatabaseName: getEnv("MYSQL_DATABASE", "venta_de_pasajes"),
		},
		ServerPort: getEnv("SERVER_PORT", "8080"),
		UsingMongo: true,
	}
}

// getEnv es una función de utilidad para obtener valores de variables de entorno con un valor predeterminado
func getEnv(key, fallbackValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallbackValue
}

// getEnvDuration es una función de utilidad para obtener valores de variables de entorno como duración
func getEnvDuration(key string, fallbackValue time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallbackValue
}
