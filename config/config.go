package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"

	"github.com/tidwall/sjson"
)

// ...
var (
	Conf Config
)

const (
	configPath = "./config.json"
)

// Config ...
type Config struct {
	DatabaseURL string `json:"databaseURL" env:"DATABASE_URL"`
	AdminToken  string `json:"adminToken" env:"ADMIN_TOKEN"`
}

func readConfig() Config {
	f, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Cannot run scripts without %s\n", configPath)
		} else {
			fmt.Printf("Failed to load config from %s: %s\n", configPath, err.Error())
		}
	}

	var conf Config
	err = json.NewDecoder(f).Decode(&conf)
	if err != nil {
		fmt.Printf("%v is in an incorrect format.\n", configPath)
	}

	if err := env.Parse(&conf); err != nil {
		fmt.Printf("%+v\n", err)
	}

	return conf
}

// UpdateConfig writes a value into the config file
func UpdateConfig(key string, value interface{}) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Println("Error reading config file")
		panic(err)
	}
	newData, _ := sjson.Set(string(data), key, value)

	err = ioutil.WriteFile(configPath, []byte(newData), 0644)
	if err != nil {
		fmt.Println("Error writing config file")
		panic(err)
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
	Conf = readConfig()
}
