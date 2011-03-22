package main

import (
	"os"
	"sort"
	"path"
	"flag"
	"http"
	"io"
	"bufio"
	"log"
	"template"
	"io/ioutil"
)


func UrlHtmlFormatter(w io.Writer, fmt string, v ...interface{}) {
	template.HTMLEscape(w, []byte(http.URLEscape(v[0].(string))))
}

var fmap = template.FormatterMap{
	"html": template.HTMLFormatter,
	"url+html": UrlHtmlFormatter,
}

var site_template = template.MustParseFile("site.html", fmap)


type recipe struct {
	Original string
	Description string
	Ingredients []string
	Directions []string
}

func RecipeFromFile(filename string) *recipe {
	var result, err = ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("ReadFile", err)
	}

	f, err := os.Open(filename, os.O_RDONLY, 0)	
	if err != nil {
		log.Fatal("open", err)
	}
	var input = bufio.NewReader(f)

	var description, _ = input.ReadString('\n')

	var ingredients = [...]string{"TODO", "TODO"}[:]
	var directions = [...]string{"TODO", "TODO"}[:]

	return &recipe{Original: string(result), Description: description, Ingredients: ingredients, Directions: directions}
}

type page struct {
	Title string
	Names []string
	Recipe *recipe
}

func HomePageHandler(w http.ResponseWriter, req *http.Request) {
	f, err := os.Open("recipes", os.O_RDONLY, 0)
	if err != nil {
		log.Fatal("open", err)
	}
	dirs, err := f.Readdir(-1)
	if err != nil {
		log.Fatal("readdir", err)
	}
	names := make([]string, len(dirs))
	for i, d := range dirs {
		names[i] = d.Name
	}
	sort.SortStrings(names)

	site_template.Execute(w, &page{Title: "Our Recipes", Names: names})
}

func RecipeHandler(w http.ResponseWriter, req *http.Request) {
	var title = path.Base(req.URL.Path)
	var recipe = RecipeFromFile(path.Join(".", req.URL.Path))
	site_template.Execute(w, &page{Title: title, Recipe: recipe})
}


var addr = flag.String("addr", ":9999", "http service address")

func main() {
	flag.Parse()
	http.Handle("/", http.HandlerFunc(HomePageHandler))
	http.Handle("/recipes/", http.HandlerFunc(RecipeHandler))
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
