package gorai

import (
	"crypto/tls"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/go51/auth551"
	"github.com/go51/container551"
	"github.com/go51/cookie551"
	"github.com/go51/log551"
	"github.com/go51/memcache551"
	"github.com/go51/model551"
	"github.com/go51/mysql551"
	"github.com/go51/response551"
	"github.com/go51/router551"
	"github.com/go51/secure551"
	"net/http"
	"os"
	"time"
)

type gorai struct {
	config       *Config
	logger       *log551.Log551
	router       *router551.Router
	modelManager *model551.Model
	auth         *auth551.Auth
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

	// Logger
	if isConsole() {
		g.logger = log551.New(&g.config.Framework.CommandLog)
	} else {
		g.logger = log551.New(&g.config.Framework.SystemLog)
	}
	g.logger.Open()
	defer g.logger.Close()

	// Router
	g.router = router551.Load()

	// ModelManager
	g.modelManager = model551.Load()

	// Add Auth Model
	g.modelManager.Add(auth551.NewUserModel, auth551.NewUserModelPointer)
	g.modelManager.Add(auth551.NewUserTokenModel, auth551.NewUserTokenModelPointer)

	// Auth
	g.auth = auth551.Load(&g.config.Framework.Auth)
}

func (g *gorai) Run() {
	if isConsole() {
		consoleHandler()
		return
	}

	cer, err := tls.LoadX509KeyPair(g.config.Framework.WebServerSSL.CtrFile, g.config.Framework.WebServerSSL.KeyFile)
	if err != nil {
		panic(err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}

	server := &http.Server{
		Addr:         g.config.Framework.WebServer.Host + ":" + g.config.Framework.WebServer.Port,
		Handler:      webReWriteHandler(),
		ReadTimeout:  g.config.Framework.WebServer.ReadTimeout * time.Second,
		WriteTimeout: g.config.Framework.WebServer.WriteTimeout * time.Second,
	}
	serverSSL := &http.Server{
		Addr:         g.config.Framework.WebServerSSL.Host + ":" + g.config.Framework.WebServerSSL.Port,
		Handler:      webHandler(),
		ReadTimeout:  g.config.Framework.WebServerSSL.ReadTimeout * time.Second,
		WriteTimeout: g.config.Framework.WebServerSSL.WriteTimeout * time.Second,
		TLSConfig:    config,
	}

	gracehttp.Serve(server, serverSSL)
}

func isConsole() bool {
	return len(os.Args) > 1
}

func consoleHandler() {

	g := Load()

	l := log551.New(&g.config.Framework.CommandLog)
	l.Open()
	defer l.Close()

	sid := secure551.Hash()

	// Routing
	name := os.Args[1]
	l.Debugf("%s [ Command ] %s", sid[:10], name)
	route := g.router.FindRouteByName(router551.COMMAND.String(), name)
	if route == nil {
		l.Errorf("%s %s", sid[:10], "Action not found...")
		return
	}

	// Options
	optionArgs := os.Args[2:]
	if len(optionArgs)%2 == 1 {
		l.Errorf("%s %s", sid[:10], "Missing options.")
	}
	options := make(map[string]string, len(optionArgs)/2)
	for i := 0; i < len(optionArgs); i += 2 {
		options[optionArgs[i][1:]] = optionArgs[i+1]
	}

	mysql := mysql551.New(&g.config.Framework.Database)
	mysql.Open()
	defer mysql.Close()

	session := memcache551.New(&g.config.Framework.Session.Server, sid)

	c := container551.New()
	c.SetSID(sid)
	c.SetResponseWriter(nil)
	c.SetRequest(nil)
	c.SetLogger(l)
	c.SetCookie(nil)
	c.SetDb(mysql)
	c.SetSession(session)
	c.SetModel(g.modelManager)
	c.SetCommandOptions(options)
	if g.config.Framework.WebServerSSL.Port == "443" {
		c.SetBaseURL("https://" + g.config.Framework.WebServerSSL.Host)
	} else {
		c.SetBaseURL("https://" + g.config.Framework.WebServerSSL.Host + ":" + g.config.Framework.WebServerSSL.Port)
	}

	action := route.Action()
	action(c)

}

func webHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/static/", staticResource)
	mux.HandleFunc("/", rootFunc)

	return mux
}

func webReWriteHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		g := Load()

		http.Redirect(w, r, "https://"+g.config.Framework.WebServerSSL.Host+":"+g.config.Framework.WebServerSSL.Port+r.URL.Path, http.StatusMovedPermanently)
	})

	return mux
}

func staticResource(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control-Max-Age", "10")
	http.ServeFile(w, r, r.URL.Path[1:])

}

func rootFunc(w http.ResponseWriter, r *http.Request) {

	g := Load()

	l := log551.New(&g.config.Framework.SystemLog)
	l.Open()
	defer l.Close()

	cookie := cookie551.New(w, r)

	sid := g.sid(cookie)
	l.Debugf("%s SID: %s", sid[:10], sid)
	session := memcache551.New(&g.config.Framework.Session.Server, sid)

	route := g.router.FindRouteByPathMatch(r.Method, r.URL.Path)
	response551.UrlFunction = g.router.Url

	var data interface{} = nil
	if route != nil {
		mysql := mysql551.New(&g.config.Framework.Database)
		mysql.Open()
		defer mysql.Close()

		l.Debugf("%s --[ Routing ]--", sid[:10])
		l.Debugf("%s Path: %s", sid[:10], r.URL.Path)
		l.Debugf("%s Neme: %s", sid[:10], route.Name())

		c := container551.New()
		c.SetSID(sid)
		c.SetResponseWriter(w)
		c.SetRequest(r)
		c.SetLogger(l)
		c.SetCookie(cookie)
		c.SetDb(mysql)
		c.SetSession(session)
		c.SetModel(g.modelManager)
		c.SetAuth(g.auth)
		c.SetUrlFunc(g.router.Url)
		if g.config.Framework.WebServerSSL.Port == "443" {
			c.SetBaseURL("https://" + g.config.Framework.WebServerSSL.Host)
		} else {
			c.SetBaseURL("https://" + g.config.Framework.WebServerSSL.Host + ":" + g.config.Framework.WebServerSSL.Port)
		}

		action := route.Action()
		data = action(c)
		response551.Response(w, r, data, route.PackageName(), route.Name(), c.User(), g.config.Application)
	} else {
		l.Errorf("%s --[ Routing ]--", sid[:10])
		l.Errorf("%s Path: %s", sid[:10], r.URL.Path)
		l.Errorf("%s Neme: Route not found.", sid[:10])
		data = response551.Error(404, "Route not found.")
		response551.Response(w, r, data, "", "", nil, g.config.Application)
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

func (g *gorai) ModelManager() *model551.Model {
	return g.modelManager
}

func (g *gorai) sid(cookie *cookie551.Cookie) string {
	sid, err := cookie.Get(g.config.Framework.Session.CookieKeyName)
	if err == nil {
		return sid
	}

	sid = secure551.Hash()

	cookie.Set(g.config.Framework.Session.CookieKeyName, sid, 60*60*24*365)

	return sid

}
