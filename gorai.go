package gorai

import (
	"github.com/facebookgo/grace/gracehttp"
	"github.com/go51/container551"
	"github.com/go51/log551"
	"github.com/go51/response551"
	"github.com/go51/router551"
	"github.com/go51/secure551"
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

	l := log551.New(&g.config.Framework.SystemLog)
	l.Open()
	defer l.Close()

	sid := g.sid(w, r)
	sidShort := sid[:10]
	l.Debugf("%s SID: %s", sidShort, sid)

	route := g.router.FindRouteByPathMatch(r.Method, r.URL.Path)

	var data interface{} = nil
	if route != nil {
		l.Debugf("%s --[ Routing ]--", sidShort)
		l.Debugf("%s Path: %s", sidShort, r.URL.Path)
		l.Debugf("%s Neme: %s", sidShort, route.Name())
		c := container551.New()

		action := route.Action()
		data = action(c)
		response551.Response(w, r, data, route.PackageName(), route.Name())
	} else {
		l.Debugf("%s --[ Routing ]--", sidShort)
		l.Debugf("%s Path: %s", sidShort, r.URL.Path)
		l.Debugf("%s Neme: Route not found.", sidShort)
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

func (g *gorai) sid(w http.ResponseWriter, r *http.Request) string {
	cookie, _ := r.Cookie(g.config.Framework.Session.CookieKeyName)
	if cookie != nil {
		return cookie.String()
	}

	sid := secure551.Hash()
	expire := time.Now().AddDate(1, 0, 0)
	setCookie := http.Cookie{
		Name:     "GOSID",
		Value:    sid,
		Expires:  expire,
		HttpOnly: true,
		Raw:      sid,
	}

	http.SetCookie(w, &setCookie)

	return sid

}
