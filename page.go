package main

import (
	"http"
	"io"
	"template"
	"os"
	"crypto/md5"
	"bytes"
	"fmt"
	"recipe"
)

func UrlHtmlFormatter(w io.Writer, fmt string, v ...interface{}) {
	template.HTMLEscape(w, []byte(http.URLEscape(v[0].(string))))
}

var fmap = template.FormatterMap{
	"html":     template.HTMLFormatter,
	"url+html": UrlHtmlFormatter,
}

var site_template = template.MustParseFile("site.html", fmap)


type page struct {
	Title      string
	Stylesheet string
	NotFound   bool
	Recipes    []*recipe.Recipe
	Recipe     *recipe.Recipe
}

func newPage(title string) *page {
	return &page{Title: title, Stylesheet: "http://static.eikeon.com/site.css^39d51ea9dec24fd6ba35fa1fc033942e"}
}

func (p *page) Write(w http.ResponseWriter, req *http.Request) (err os.Error) {
	var bw bytes.Buffer
	h := md5.New()
	mw := io.MultiWriter(&bw, h)
	err = site_template.Execute(mw, p)
	if err == nil {
		w.SetHeader("ETag", fmt.Sprintf("\"%x\"", h.Sum()))
		w.Write(bw.Bytes())
	}
	return err
}
