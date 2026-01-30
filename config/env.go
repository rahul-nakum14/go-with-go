package config

type Config struct{
	PublicHost string
	Port string
	DBHost string
	DBPort string
	DBUser string
	DBPassword string
	DBName string
	DBSSLMode string
}

var Envs = initConfig()


func initConfig() Config {
	godotenv.Load()
	return Config{
		PublicHost: getEnv("PUBLIC_HOST", "localhost"),
		Port: getEnv("PORT", "8080"),
		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "27017"),
		DBUser: getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "root"),
		DBName: getEnv("DB_NAME", "prodgo"),
		DBSSLMode: getEnv("DB_SSL_MODE", "disable"),
	}
}

func getEnv(key string, defaultValue string) string {
	if value,ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}