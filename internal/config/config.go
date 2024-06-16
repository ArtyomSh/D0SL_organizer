package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Listen struct {
		Port string `yaml:"port" env-default:":8001"`
	} `yaml:"listen"`
	Milvus struct {
		Host     string `yaml:"host" env-default:"localhost"`
		Port     string `yaml:"port" env-default:"19530"`
		Database string `yaml:"database" env-default:"milvus"`
		// Username string `json:"username"`
		// Password string `json:"password"`
	} `yaml:"milvus"`
	// Repository struct {
	// 	Type string `json:"type"`
	// } `yaml:"repository"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		// _ = os.Chdir("../../")
		if err := cleanenv.ReadConfig("configs/config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Println(help)
		}
		log.Println(instance)
	})
	return instance
}
