package urlshort

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var pathsToUrls = map[string]string{
	"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
	"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
}
var yamls = `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`

var jsons = `
[
	{
		"path": "/urlshort",
		"url": "https://github.com/gophercises/urlshort"
	},
	{
		"path": "/urlshort-final",
		"url": "https://github.com/gophercises/urlshort/tree/solution"
	}
]`

func TestBuildMap(t *testing.T) {
	pathRedirects := []pathRedirect{
		{
			Path:        "/urlshort",
			RedirectUrl: "https://github.com/gophercises/urlshort",
		},
		{
			Path:        "/urlshort-final",
			RedirectUrl: "https://github.com/gophercises/urlshort/tree/solution",
		},
	}

	mapping := buildMap(pathRedirects)

	if len(mapping) != 2 {
		t.Errorf("got length %d, wanted 2", len(mapping))
	}
}

func TestParseYaml(t *testing.T) {
	pathRedirects, err := parseYaml([]byte(yamls))

	if err != nil {
		t.Errorf("got error %v", err)
	}

	if len(pathRedirects) != 2 {
		t.Errorf("got length %d, wanted 2", len(pathRedirects))
	}
}

func TestParseJson(t *testing.T) {
	pathRedirects, err := parseJson([]byte(jsons))

	if err != nil {
		t.Errorf("got error %v", err)
	}

	if len(pathRedirects) != 2 {
		t.Errorf("got length %d, wanted 2", len(pathRedirects))
	}
}

func TestMapHandler(t *testing.T) {
	handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if _, ok := pathsToUrls[r.URL.Path]; ok {
			t.Errorf("fallback called for %v, expected redirect", r.URL.Path)
		} else {
			t.Logf("Fallback called for path %v", r.URL.Path)
		}
	})

	mappedHandler := MapHandler(pathsToUrls, handler)

	for k, v := range pathsToUrls {
		req := httptest.NewRequest("GET", k, nil)
		rw := httptest.NewRecorder()

		mappedHandler(rw, req)

		if rw.Result().StatusCode == http.StatusTemporaryRedirect {
			t.Logf("mapped %v to %v", k, v)
		} else {
			t.Errorf("expected mapping from %v to %v", k, v)
		}
	}
}

type httpHandler struct {
	handleFunc http.HandlerFunc
}

func (h *httpHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.handleFunc(rw, r)
}

var _ http.Handler = (*httpHandler)(nil)
