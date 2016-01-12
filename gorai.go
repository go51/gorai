package gorai
import (
	"net/http"
	"time"
	"github.com/facebookgo/grace/gracehttp"
	"fmt"
)

type gorai struct {
	config *Config
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
	g.config = loadConfig()
}

func (g *gorai) Run() {
	server := &http.Server{
		Addr: g.config.Framework.WebServer.Host + ":" + g.config.Framework.WebServer.Port,
		Handler: webHandler(),
		ReadTimeout: g.config.Framework.WebServer.ReadTimeout * time.Second,
		WriteTimeout: g.config.Framework.WebServer.WriteTimeout * time.Second,
	}
	gracehttp.Serve(server)
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