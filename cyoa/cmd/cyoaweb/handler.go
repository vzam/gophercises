package main

import (
	"cyoa"
	"html/template"
	"net/http"
)

type PathFunc func(*http.Request) string

type HandlerOption func(h *AdventureHandler)

func WithPathFunc(pathFunc PathFunc) HandlerOption {
	return func(h *AdventureHandler) {
		h.pathFunc = pathFunc
	}
}

func WithDefaultChapter(story string) HandlerOption {
	return func(h *AdventureHandler) {
		h.defaultStory = story
	}
}

type AdventureHandler struct {
	t *template.Template
	c map[string]cyoa.Chapter

	// pathFunc takes a http Request and returns the key of a story
	pathFunc     func(*http.Request) string
	defaultStory string
}

func NewHandler(adventure map[string]cyoa.Chapter, t *template.Template, opts ...HandlerOption) *AdventureHandler {
	h := AdventureHandler{
		t: t,
		c: adventure,
		pathFunc: func(r *http.Request) string {
			return r.URL.Path[1:]
		},
		defaultStory: "intro",
	}

	for _, opt := range opts {
		opt(&h)
	}

	return &h
}

func (h *AdventureHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	path := h.pathFunc(r)

	if path == "" {
		path = h.defaultStory
	}

	if c, ok := h.c[path]; ok {
		h.t.Execute(rw, c)
		return
	}
	rw.WriteHeader(http.StatusNotFound)
}
