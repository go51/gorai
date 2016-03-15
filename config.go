package gorai

import (
	"encoding/json"
	"github.com/go51/auth551"
	"github.com/go51/log551"
	"github.com/go51/memcache551"
	"github.com/go51/mysql551"
	"io/ioutil"
	"os"
	"time"
)

type Config struct {
	Framework   ConfigFramework   `json:"framework"`
	Application ConfigApplication `json:"application"`
}

type ConfigFramework struct {
	WebServer    ConfigWebServer    `json:"web_server"`
	WebServerSSL ConfigWebServerSSL `json:"web_server_ssl"`
	SystemLog    log551.Config      `json:"system_log"`
	CommandLog   log551.Config      `json:"command_log"`
	Session      ConfigSession      `json:"session"`
	Database     mysql551.Config    `json:"database"`
	Auth         auth551.Config     `json:"auth"`
}

type ConfigWebServer struct {
	Host         string        `json:"host"`
	Port         string        `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
}

type ConfigWebServerSSL struct {
	Host         string        `json:"host"`
	Port         string        `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	CtrFile      string        `json:"crt_file"`
	KeyFile      string        `json:"key_file"`
}

type ConfigSession struct {
	CookieKeyName string             `json:"cookie_key_name"`
	Server        memcache551.Config `json:"server"`
}

type ConfigApplication struct {
	Name                      string `json:"name"`
	GoogleAnalyticsTrackingId string `json:"google_analytics_tracking_id"`
}

var configInstance *Config

func loadConfig(appConfig interface{}) *Config {
	if configInstance != nil {
		return configInstance
	}

	path := getConfigFilePath()

	file := getConfigJson(path)

	configInstance = &Config{}
	json.Unmarshal(file, &configInstance)

	configInstance.Application = appConfig

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
