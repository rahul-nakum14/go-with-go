package config

type Config struct {
	DBDriver string
	DBSource string
}

func Load() *Config {
	return &Config{
		DBDriver: "postgres", 
		DBSource: "postgres://user:pass@localhost:5432/demo?sslmode=disable",
	}
}
