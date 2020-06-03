package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
)

const (
	ServiceName = "cosmoscan-api"
	configPath  = "./config.json"
	Currency    = "atom"
)

type (
	Config struct {
		API                   API        `json:"api"`
		Mysql                 Mysql      `json:"mysql"`
		Clickhouse            Clickhouse `json:"clickhouse"`
		Parser                Parser     `json:"parser"`
		CMCKey                string     `json:"cmc_key"`
		ECDSAPublicKeyBase64  string     `json:"ecdsa_public_key_base64"`
		ECDSAPrivateKeyBase64 string     `json:"ecdsa_private_key_base64"`
	}
	Parser struct {
		Node     string `json:"node"`
		Batch    uint64 `json:"batch"`
		Fetchers uint64 `json:"fetchers"`
	}
	API struct {
		Port         string   `json:"port"`
		AllowedHosts []string `json:"allowed_hosts"`
	}
	Mysql struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		DB       string `json:"db"`
		User     string `json:"user"`
		Password string `json:"password"`
	}
	Clickhouse struct {
		Protocol string `json:"protocol"`
		Host     string `json:"host"`
		Port     uint   `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"database"`
	}
)

func GetConfig() Config {
	path, _ := filepath.Abs(configPath)
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("Invalid config path : "+configPath, err)
	}
	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatalln("Failed unmarshal config ", err)
	}
	return config
}
