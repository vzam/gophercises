package main

import (
	"cyoa"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	port := flag.Int("port", 3030, "The port where to listen at for story requests.")
	tplName := flag.String("template", "", "The template to use for rendering the adventure.")
	advName := flag.String("adventure", "", "The adventure to serve.")
	flag.Parse()

	if *tplName == "" {
		fmt.Println("template file is required")
		os.Exit(1)
	}
	if *advName == "" {
		fmt.Println("adventure file is required")
		os.Exit(1)
	}

	tpl, err := parseHtmlFileTemplate(*tplName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	adv, err := cyoa.ParseJsonFileAdventure(*advName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	handler := NewHandler(adv.Chapters, tpl, WithPathFunc(func(r *http.Request) string {
		fmt.Println(r.URL.Path)
		return r.URL.Path[len("/stories/"):]
	}), WithDefaultChapter(adv.InitialChapter))

	mux := http.NewServeMux()
	mux.HandleFunc("/stories/", func(rw http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(rw, r)
	})

	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}

func parseHtmlTemplate(r io.Reader, name string) (*template.Template, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("unable to read template: %v", err)
	}

	t := template.New(name)
	t.Parse(string(b))

	return t, nil
}

func parseHtmlFileTemplate(filename string) (*template.Template, error) {
	tpl, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open template file: %v", err)
	}
	defer tpl.Close()

	return parseHtmlTemplate(tpl, filename)
}
