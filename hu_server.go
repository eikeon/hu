package main

import (
	"os"
	"path"
	"flag"
	"http"
	"log"
	"strings"
	"fmt"
	"time"
	"recipe"
)


func NotFoundHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Cache-Control", "max-age=10, must-revalidate")
	w.WriteHeader(http.StatusNotFound)
	page := newPage("Not Found")
	page.NotFound = true
	page.Write(w, req)
	return
}

func setCacheControl(w http.ResponseWriter, req *http.Request) {
	if req.Header["X-Draft"] != nil {
		w.Header().Set("Cache-Control", "max-age=1, must-revalidate")
	} else {
		now := time.UTC()
		d := time.Time{2011, 4, 11, 3, 0, 0, time.Monday, 0, "UTC"}
		ONE_WEEK := int64(604800)
		ttl := ONE_WEEK - (now.Seconds()-d.Seconds())%ONE_WEEK
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", ttl))
	}
	w.Header().Set("Vary", "Accept-Encoding")
}

func PageHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		setCacheControl(w, req)
		page := newPage("")
		page.Write(w, req)
		return
	}
	StaticHandler(w, req)
}

func RecipesHandler(w http.ResponseWriter, req *http.Request) {
	var r = recipe.Recipes[path.Base(req.URL.Path)]
	if r != nil {
		var p = strings.Replace(req.URL.Path, "recipes", "recipe", -1)
		w.Header().Set("Location", p)
		w.WriteHeader(http.StatusMovedPermanently)
		return
	}

	if req.URL.Path != "/recipes/" {
		NotFoundHandler(w, req)
		return
	}
	setCacheControl(w, req)
	page := newPage("Recipes")
	page.Recipes = recipe.Recipe_list
	page.Write(w, req)

}

func RecipeHandler(w http.ResponseWriter, req *http.Request) {
	var r = recipe.Recipes[path.Base(req.URL.Path)]
	if r == nil {
		NotFoundHandler(w, req)
		return
	}
	setCacheControl(w, req)
	page := newPage(r.Name + " Recipe")
	page.Recipe = r
	page.Write(w, req)
}

func CanonicalHostHandler(w http.ResponseWriter, req *http.Request) {
	var canonical = "www.eikeon.com"
	if req.Host != canonical {
		http.Redirect(w, req, "http://"+canonical+req.URL.Path, http.StatusMovedPermanently)
	} else {
		NotFoundHandler(w, req)
	}
	// TODO: set CacheControl
}

func StaticHandler(w http.ResponseWriter, req *http.Request) {
	var filename = path.Join(*StaticRoot, req.URL.Path)
	f, err := os.Open(filename)
	if err != nil {
		log.Print(err)
		NotFoundHandler(w, req)
		return
	}
	err = f.Close()

	if strings.Contains(req.URL.Path, "^") {
		w.Header().Set("Cache-Control", "max-age=3153600")
	} else {
		w.Header().Set("Cache-Control", "max-age=1, must-revalidate")
	}
	http.ServeFile(w, req, filename)
}

var Address *string
var StaticRoot *string

func main() {
	Address = flag.String("address", ":9999", "http service address")
	StaticRoot = flag.String("root", "static", "...")
	flag.Parse()

	http.Handle("eikeon.com/", http.HandlerFunc(CanonicalHostHandler))
	http.Handle("/", http.HandlerFunc(PageHandler))
	http.Handle("/recipes/", http.HandlerFunc(RecipesHandler))
	http.Handle("/recipe/", http.HandlerFunc(RecipeHandler))

	err := http.ListenAndServe(*Address, nil)
	if err != nil {
		log.Print("ListenAndServe:", err)
	}
}
