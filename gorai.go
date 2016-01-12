package gorai
import (
	"net/http"
	"time"
	"github.com/facebookgo/grace/gracehttp"
	"fmt"
)

type gorai struct {}

var goraiInstance *gorai = nil

func Load() *gorai {
	if goraiInstance != nil {
		return goraiInstance
	}

	goraiInstance = &gorai{}

	return goraiInstance
}

func (g *gorai) Run() {
	server := &http.Server{
		Addr: ":8080",
		Handler: webHandler(),
		ReadTimeout: 30 * time.Second,
		WriteTimeout: 60 * time.Second,
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
