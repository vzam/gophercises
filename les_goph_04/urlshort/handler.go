package urlshort

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-yaml/yaml"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if redirectUrl, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(rw, r, redirectUrl, http.StatusTemporaryRedirect)
		} else {
			fallback.ServeHTTP(rw, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yamlb []byte, fallback http.Handler) (http.HandlerFunc, error) {
	redirects, err := parseYaml(yamlb)
	if err != nil {
		return nil, err
	}

	return MapHandler(buildMap(redirects), fallback), nil
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
// [
//		{
// 			"path": "/some-path",
//       	"url": "https://www.some-url.com/demo"
//		}
// ]
//
// The only errors that can be returned all related to having
// invalid JSON data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func JSONHandler(jsonb []byte, fallback http.Handler) (http.HandlerFunc, error) {
	redirects, err := parseJson(jsonb)
	if err != nil {
		return nil, err
	}

	return MapHandler(buildMap(redirects), fallback), nil
}

func buildMap(redirects []pathRedirect) map[string]string {
	mapping := make(map[string]string, len(redirects))

	for _, redirect := range redirects {
		mapping[redirect.Path] = redirect.RedirectUrl
	}
	return mapping
}

func parseJson(jsonb []byte) ([]pathRedirect, error) {
	var redirects []pathRedirect
	if err := json.Unmarshal(jsonb, &redirects); err != nil {
		return nil, fmt.Errorf("JSON parser error: %v", err)
	}
	return redirects, nil
}

func parseYaml(yamlb []byte) ([]pathRedirect, error) {
	var redirects []pathRedirect
	if err := yaml.Unmarshal(yamlb, &redirects); err != nil {
		return nil, fmt.Errorf("YAML parser error: %v", err)
	}
	return redirects, nil
}

type pathRedirect struct {
	Path        string `yaml:"path" json:"path"`
	RedirectUrl string `yaml:"url" json:"url"`
}
