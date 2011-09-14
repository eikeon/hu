package main

import (
	"http"
	"io"
	"io/ioutil"
	"template"
	"os"
	"crypto/md5"
	"bytes"
	"fmt"
	"recipe"
	"log"
)

var site_template = template.Must(template.ParseFile("site.html"))
var site_style string

func init() {
	site_style_bytes, err := ioutil.ReadFile("site.css")
	if err != nil {
		log.Fatal("could not read site.css")
	}
	site_style = minify(string(site_style_bytes))
}

type page struct {
	Title      string
	Stylesheet string
	NotFound   bool
	Recipes    []*recipe.Recipe
	Recipe     *recipe.Recipe
	// TODO: add baseURL and use URL.ParseURL to resolve relative URLs such as the photo URLs.
}

func newPage(title string) *page {
	return &page{Title: title, Stylesheet: site_style}
}

func (p *page) Write(w http.ResponseWriter, req *http.Request) (err os.Error) {
	var bw bytes.Buffer
	h := md5.New()
	mw := io.MultiWriter(&bw, h)
	err = site_template.Execute(mw, p)
	if err == nil {
		w.Header().Set("ETag", fmt.Sprintf("\"%x\"", h.Sum()))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", bw.Len()))
		w.Write(bw.Bytes())
	}
	return err
}
