package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
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
		now := time.Now().UTC()
		//d := time.Time{2011, 4, 11, 3, 0, 0, 0, time.Monday, 0, "UTC"}
		d := time.Date(2011, 4, 11, 3, 0, 0, 0, time.UTC)
		TTL := int64(86400)
		ttl := TTL - (now.Unix()-d.Unix())%TTL // shift
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", ttl))
	}
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
	var r = Recipes[path.Base(req.URL.Path)]
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
	page.Recipes = Recipe_list
	page.Write(w, req)

}

func RecipeHandler(w http.ResponseWriter, req *http.Request) {
	var r = Recipes[path.Base(req.URL.Path)]
	if r == nil {
		NotFoundHandler(w, req)
		return
	}
	if req.URL.Path[len(req.URL.Path)-1] != '/' {
		http.Redirect(w, req, req.URL.Path+"/", http.StatusMovedPermanently)
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
	} else {
		defer f.Close()
	}

	setCacheControl(w, req)

	http.ServeFile(w, req, filename)
}

var Address *string
var StaticRoot *string
var RecipesLocation *string

func main() {
	Address = flag.String("address", ":9999", "http service address")
	StaticRoot = flag.String("root", "static", "...")
	RecipesLocation = flag.String("recipes", "recipes", "location of recipes")
	flag.Parse()

	if strings.Contains(*RecipesLocation, "://") {
		r, err := http.Get(*RecipesLocation)
		if err == nil {
			initRecipes(r.Body)
		} else {
			log.Print(err)
		}
	} else {
		f, err := os.Open(*RecipesLocation)
		if err == nil {
			initRecipes(f)
		} else {
			log.Print(err)
		}
	}

	http.Handle("eikeon.com/", http.HandlerFunc(CanonicalHostHandler))
	http.Handle("/", http.HandlerFunc(PageHandler))
	http.Handle("/recipes/", http.HandlerFunc(RecipesHandler))
	http.Handle("/recipe/", http.HandlerFunc(RecipeHandler))
	http.Handle("/recipe/seaweed_and_cabbage_saute/",
		http.RedirectHandler("/recipe/seaweed_and_cabbage_sauté/",
			http.StatusMovedPermanently))

	err := http.ListenAndServe(*Address, nil)
	if err != nil {
		log.Print("ListenAndServe:", err)
	}
}
