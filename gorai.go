package gorai

import (
	"github.com/facebookgo/grace/gracehttp"
	"github.com/go51/container551"
	"github.com/go51/log551"
	"github.com/go51/response551"
	"github.com/go51/router551"
	"net/http"
	"time"
)

type gorai struct {
	config *Config
	logger *log551.Log551
	router *router551.Router
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

	// Logger
	g.logger = log551.New(&g.config.Framework.SystemLog)
	g.logger.Open()
	defer g.logger.Close()

	g.logger.Information("--[ initialize gorai - START ]--")
	g.logger.Information("Success! [Log551]")

	// Router
	g.router = router551.Load()
	g.logger.Information("Success! [Router551]")

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

	g := Load()

	l := g.logger
	l.Open()
	defer l.Close()

	route := g.router.FindRouteByPathMatch(r.Method, r.URL.Path)

	var data interface{} = nil
	if route != nil {
		l.Debug("--[ Routing ]--")
		l.Debugf("Path: %s", r.URL.Path)
		l.Debugf("Neme: %s", route.Name())
		action := route.Action()
		c := container551.New()
		data = action(c)
		response551.Response(w, r, data, route.PackageName(), route.Name())
	} else {
		l.Debug("--[ Routing ]--")
		l.Debugf("Path: %s", r.URL.Path)
		l.Debugf("Neme: Route not found.")
		data = response551.Error(404, "Page, Action not found.")
		response551.Response(w, r, data, "", "")
	}

}

func (g *gorai) Config() *Config {
	return g.config
}

func (g *gorai) Logger() *log551.Log551 {
	return g.logger
}

func (g *gorai) Router() *router551.Router {
	return g.router
}
