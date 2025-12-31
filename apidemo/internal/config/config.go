package config

import (
	"log"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
    Env         string           `mapstructure:"env" validate:"required,oneof=development production"`
    StoragePath string           `mapstructure:"storage_path" validate:"required"`
    HTTPServer  HTTPServerConfig `mapstructure:"http_server" validate:"required"`
}

type HTTPServerConfig struct {
    Address string `mapstructure:"address" validate:"required,hostname_port"`
}
func LoadConfig() (*Config, error) {
    
    viper.SetConfigName("local")  
    viper.SetConfigType("yaml")   
    viper.AddConfigPath("./config") 

    if err := viper.ReadInConfig(); err != nil {
        log.Printf("Error reading config file: %s", err)
        return nil, err
    }
    
    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        log.Printf("Error unmarshaling config: %s", err)
        return nil, err
    }

    log.Println("Config File ->", config)
    validate := validator.New()
    if err := validate.Struct(config); err != nil {
        log.Printf("Config validation failed: %s", err)
        return nil, err
    }

    return &config, nil
}