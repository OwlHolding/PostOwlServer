package main

import (
	"encoding/json"
	"log"
	"os"
)

type ServerConfig struct {
	Token           string
	Url             string
	Port            string
	MaxBotConns     int
	ChansPerUser    int
	CertFile        string
	KeyFile         string
	RedisUrl        string
	MlServers       []string
	SqlUser         string
	SqlPass         string
	MaxSqlConns     int
	MaxSqlIdleConns int
	AdminChatIDs    []int64
	BanList         []int64
	WhiteList       []int64
	AccessKeys      []string
}

func GenerateConfig(path string) {
	config, _ := json.Marshal(ServerConfig{})
	err := os.WriteFile(path, config, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func LoadConfig(path string) ServerConfig {
	byte_config, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	var config ServerConfig
	err = json.Unmarshal(byte_config, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func LoadConfigFromEnv(variable string) ServerConfig {
	value, exists := os.LookupEnv(variable)
	if !exists {
		log.Fatalf("Variable %s does not exist", variable)
	}
	return LoadConfig(value)
}
