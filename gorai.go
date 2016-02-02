package gorai

import (
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

func Load(appConfig ...interface{}) *gorai {
	if goraiInstance != nil {
		return goraiInstance
	}

	goraiInstance = &gorai{}

	goraiInstance.initialize(appConfig[0])

	return goraiInstance
}

func (g *gorai) initialize(appConfig interface{}) {
	// load config
	g.config = loadConfig(appConfig)

	// Logger
	g.logger = log551.New(&g.config.Framework.SystemLog)
	g.logger.Open()
	defer g.logger.Close()

	g.logger.Information("--[ initialize gorai - START ]--")
	g.logger.Information("Success! [Log551]")

	// Router
	g.router = router551.Load()
	g.logger.Information("Success! [Router551]")

	// ModelManager
	g.modelManager = model551.Load()
	g.logger.Information("Success! [Model551]")

	// Add Auth Model
	g.modelManager.Add(auth551.NewUserModel, auth551.NewUserModelPointer)
	g.modelManager.Add(auth551.NewUserTokenModel, auth551.NewUserTokenModelPointer)
	g.logger.Information("Success! [Add auth models]")

	// Auth
	g.auth = auth551.Load(g.config.Framework.Auth)
	g.logger.Information("Success! [Auth551]")

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

	mux.HandleFunc("/static/", staticResource)
	mux.HandleFunc("/", rootFunc)

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

	mysql := mysql551.New(&g.config.Framework.Database)

	cookie := cookie551.New(w, r)

	sid := g.sid(cookie)
	sidShort := sid[:10]
	l.Debugf("%s SID: %s", sidShort, sid)

	session := memcache551.New(&g.config.Framework.Session.Server, sid)

	route := g.router.FindRouteByPathMatch(r.Method, r.URL.Path)

	var data interface{} = nil
	if route != nil {
		l.Debugf("%s --[ Routing ]--", sidShort)
		l.Debugf("%s Path: %s", sidShort, r.URL.Path)
		l.Debugf("%s Neme: %s", sidShort, route.Name())
		c := container551.New()
		c.SetSID(sid)
		c.SetResponseWriter(w)
		c.SetRequest(r)
		c.SetLogger(l)
		c.SetLogger(l)
		c.SetCookie(cookie)
		c.SetDb(mysql)
		c.SetSession(session)
		c.SetModel(g.modelManager)
		c.SetAuth(g.auth)

		action := route.Action()
		data = action(c)
		response551.Response(w, r, data, route.PackageName(), route.Name())
	} else {
		l.Errorf("%s --[ Routing ]--", sidShort)
		l.Errorf("%s Path: %s", sidShort, r.URL.Path)
		l.Errorf("%s Neme: Route not found.", sidShort)
		data = response551.Error(404, "Route not found.")
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

func (g *gorai) sid(cookie *cookie551.Cookie) string {
	sid, err := cookie.Get(g.config.Framework.Session.CookieKeyName)
	if err == nil {
		return sid
	}

	sid = secure551.Hash()

	cookie.Set(g.config.Framework.Session.CookieKeyName, sid, 60*60*24*365)

	return sid

}
