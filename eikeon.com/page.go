package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

var site_template = template.Must(template.ParseFiles("site.html"))
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
	Recipes    []*Recipe
	Recipe     *Recipe
	// TODO: add baseURL and use URL.ParseURL to resolve relative URLs such as the photo URLs.
}

func newPage(title string) *page {
	return &page{Title: title, Stylesheet: site_style}
}

func (p *page) Write(w http.ResponseWriter, req *http.Request) (err error) {
	var bw bytes.Buffer
	h := md5.New()
	mw := io.MultiWriter(&bw, h)
	err = site_template.Execute(mw, p)
	if err == nil {
		w.Header().Set("ETag", fmt.Sprintf("\"%x\"", h.Sum(nil)))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", bw.Len()))
		w.Write(bw.Bytes())
	} else {
		log.Println("template error:", err)
	}

	return err
}
