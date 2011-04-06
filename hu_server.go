package main

import (
	"sort"
	"path"
	"flag"
	"http"
	"io"
	"log"
	"strings"
	"template"
	"crypto/md5"
	"fmt"
	"bufio"
	"./recipe"
)


func UrlHtmlFormatter(w io.Writer, fmt string, v ...interface{}) {
	template.HTMLEscape(w, []byte(http.URLEscape(v[0].(string))))
}

var fmap = template.FormatterMap{
	"html": template.HTMLFormatter,
	"url+html": UrlHtmlFormatter,
}

var site_template = template.MustParseFile("site.html", fmap)


type page struct {
	Title string
	Stylesheet string
	NotFound bool
	Recipes []*recipe.Recipe
	Recipe *recipe.Recipe
}

func newPage(title string) *page {
        return &page{Title: title, Stylesheet: "http://h.eikeon.com/site.css^aa933dc876627b1a85509061aaad0bed"}
}

func NotFoundHandler(w http.ResponseWriter, req *http.Request) {
	w.SetHeader("Cache-Control", "max-age=10, must-revalidate")
	w.WriteHeader(http.StatusNotFound)
	page := newPage("Not Found")
	page.NotFound = true
	site_template.Execute(w, page)
	return
}

func HomePageHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		NotFoundHandler(w, req)
		return
	}

	page := newPage("")

	bw := bufio.NewWriter(nil)
	h := md5.New()
	mw := io.MultiWriter(bw, h)
	site_template.Execute(mw, page)

	//w.SetHeader("Vary", "Accept-Encoding")
	w.SetHeader("Cache-Control", "max-age=1, must-revalidate")
	w.SetHeader("ETag", fmt.Sprintf("\"%x\"", h.Sum()))
	site_template.Execute(w, page)
}

func RecipesHandler(w http.ResponseWriter, req *http.Request) {
	var r = recipe.Recipes[path.Base(req.URL.Path)]
	if r != nil {
		var p = strings.Replace(req.URL.Path, "recipes", "recipe", -1)
		w.SetHeader("Location", p)
		w.WriteHeader(http.StatusMovedPermanently)
		return
	}

	if req.URL.Path != "/recipes/" {
		NotFoundHandler(w, req)
		return
	}

	recipe_list := make(recipe.RecipeArray, len(recipe.Recipes))
	var i int
	for _, r := range recipe.Recipes {
		recipe_list[i] = r
		i += 1
	}
	sort.Sort(recipe_list)

	page := newPage("Recipes")
	page.Recipes = recipe_list

	bw := bufio.NewWriter(nil)
	h := md5.New()
	mw := io.MultiWriter(bw, h)
	site_template.Execute(mw, page)

	w.SetHeader("Cache-Control", "max-age=1, must-revalidate")
	w.SetHeader("ETag", fmt.Sprintf("\"%x\"", h.Sum()))
	site_template.Execute(w, page)
}

func RecipeHandler(w http.ResponseWriter, req *http.Request) {
	var r = recipe.Recipes[path.Base(req.URL.Path)]
	if r == nil {
		NotFoundHandler(w, req)
		return
	}
	page := newPage(r.Name)
	page.Recipe = r

	site_template.Execute(w, page)
}

var addr = flag.String("addr", ":9999", "http service address")

func main() {
	flag.Parse()
	http.Handle("www.eikeon.com/", http.RedirectHandler("http://eikeon.com/", http.StatusMovedPermanently))
	http.Handle("eikeon.com/", http.HandlerFunc(HomePageHandler))
	http.Handle("eikeon.com/recipes/", http.HandlerFunc(RecipesHandler))
	http.Handle("eikeon.com/recipe/", http.HandlerFunc(RecipeHandler))

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Print("ListenAndServe:", err)
	}
}
