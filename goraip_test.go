package gorai

import (
	"github.com/go51/container551"
	"github.com/go51/response551"
	"github.com/go51/router551"
	"github.com/go51/string551"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func indexAction(c *container551.Container) interface{} {
	return map[string]interface{}{
		"name": "Gorai!",
		"children": map[string]interface{}{
			"child_001": "Yuzu",
			"child_002": "Misaki",
		},
	}
}

func redirectAction(c *container551.Container) interface{} {
	uri := "https://golang.org/"
	return response551.Redirect(uri, 301)
}

func errAction(c *container551.Container) interface{} {
	return response551.Error(501, http.StatusText(501))
}

func getBody(resp *http.Response) string {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return string551.BytesToString(b)

}

func TestRoot(t *testing.T) {
	g := Load()
	g.router.Add(router551.GET, "index", "/", indexAction)
	g.router.Add(router551.GET, "redirect", "/redirect", redirectAction)
	g.router.Add(router551.GET, "err", "/err", errAction)

	ts := httptest.NewServer(http.HandlerFunc(rootFunc))
	defer ts.Close()

	r, err := http.Get(ts.URL + "/test")
	if err != nil {
		t.Logf("%#v\n", err)
		t.Error("unexpected")
		return
	}
	if r.StatusCode != 404 {
		t.Error("Status code error")
		return
	}

	r, err = http.Get(ts.URL)
	if err != nil {
		t.Error("unexpected")
		return
	}
	if r.StatusCode != 200 {
		t.Error("Status code error")
		return
	}
	body := getBody(r)
	if body == "" {
		t.Error("Response body error")
	}

	r, err = http.Get(ts.URL + "?format=json")
	if err != nil {
		t.Error("unexpected")
		return
	}
	if r.StatusCode != 200 {
		t.Error("Status code error")
		return
	}
	body = getBody(r)
	if body == "" {
		t.Error("Response body error")
	}

	// TODO: http.GET() で redirect していない状態の Response が取得する方法がわかるまで保留
	//	r, err = http.Get(ts.URL + "/redirect")
	//	if err != nil {
	//		t.Error("unexpected")
	//		return
	//	}
	//	if r.StatusCode != 301 {
	//		t.Error("Status code error")
	//		return
	//	}

	r, err = http.Get(ts.URL + "/err")
	if err != nil {
		t.Error("unexpected")
		return
	}
	if r.StatusCode != 501 {
		t.Error("Status code error")
		return
	}

}
