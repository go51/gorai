package gorai

import (
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/go51/log551"
	"net/http"
	"time"
)

type gorai struct {
	config *Config
	logger *log551.Log551
}

var goraiInstance *gorai = nil

func Load() *gorai {
	if goraiInstance != nil {
		return goraiInstance
	}

	goraiInstance = &gorai{}

	goraiInstance.initialize()

	return goraiInstance
}

func (g *gorai) initialize() {
	// load config
	g.config = loadConfig()

	g.logger = log551.New(&g.config.Framework.SystemLog)
	g.logger.Open()
	defer g.logger.Close()

	g.logger.Information("--[ initialize gorai - START ]--")
	g.logger.Information("Success! [Log551]")
	g.logger.Information("--[ initialize gorai - END   ]--")
}

func (g *gorai) Run() {
	server := &http.Server{
		Addr:         g.config.Framework.WebServer.Host + ":" + g.config.Framework.WebServer.Port,
		Handler:      webHandler(),
		ReadTimeout:  g.config.Framework.WebServer.ReadTimeout * time.Second,
		WriteTimeout: g.config.Framework.WebServer.WriteTimeout * time.Second,
	}
	gracehttp.Serve(server)

	g.logger.Close()
}

func webHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", rootFunc)

	return mux
}

func rootFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, Gorai!")
}

func (g *gorai) Config() *Config {
	return g.config
}

func (g *gorai) Logger() *log551.Log551 {
	return g.logger
}
