package gorai

import (
	"encoding/json"
	"github.com/go51/log551"
	"github.com/go51/mysql551"
	"io/ioutil"
	"os"
	"time"
)

type Config struct {
	Framework ConfigFramework `json:"framework"`
}

type ConfigFramework struct {
	WebServer ConfigWebServer `json:"web_server"`
	SystemLog log551.Config   `json:"system_log"`
	Session   ConfigSession   `json:"session"`
	Database  mysql551.Config `json:"database"`
}

type ConfigWebServer struct {
	Host         string        `json:"host"`
	Port         string        `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
}

type ConfigSession struct {
	CookieKeyName string `json:"cookie_key_name"`
}

var configInstance *Config

func loadConfig() *Config {
	if configInstance != nil {
		return configInstance
	}

	path := getConfigFilePath()

	file := getConfigJson(path)

	configInstance = &Config{}
	json.Unmarshal(file, &configInstance)

	return configInstance
}

func getConfigJson(path string) []byte {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return file
}

func getConfigFilePath() string {
	env := getConfigEnv()
	return "./config/config_" + env + ".json"
}

func getConfigEnv() string {
	return os.Getenv("GORAI_ENV")
}
