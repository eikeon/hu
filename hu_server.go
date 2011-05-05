package main

import (
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
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d, must-revalidate", ttl))
	}
}

func HomePageHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		NotFoundHandler(w, req)
		return
	}
	setCacheControl(w, req)
	page := newPage("")
	page.Write(w, req)
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

var addr = flag.String("addr", ":9999", "http service address")

func main() {
	flag.Parse()
	http.Handle("www.eikeon.com/", http.RedirectHandler("http://eikeon.com/", http.StatusMovedPermanently))
	http.Handle("/", http.RedirectHandler("http://eikeon.com/", http.StatusMovedPermanently))
	http.Handle("eikeon.com/", http.HandlerFunc(HomePageHandler))
	http.Handle("eikeon.com/recipes/", http.HandlerFunc(RecipesHandler))
	http.Handle("eikeon.com/recipe/", http.HandlerFunc(RecipeHandler))

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Print("ListenAndServe:", err)
	}
}
