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
	return &page{Title: title, Stylesheet: "/site^b8a9e95ed8b90c765a216ff17bf67510.css"}
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
