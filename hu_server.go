package main

import (
	"os"
	"sort"
	"path"
	"flag"
	"http"
	"io"
	"log"
	"strings"
	"template"
	"io/ioutil"
	"encoding/line"
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
		log.Print("ReadFile: ", err)
		return nil
	}

	f, err := os.Open(filename, os.O_RDONLY, 0)	
	if err != nil {
		log.Print("open", err)
	}
	var ingredients = [...]string{}[:]
	var directions = [...]string{}[:]

	var input = line.NewReader(f, 1024)
	line, isPrefix, err := input.ReadLine()
	if err != nil {
		log.Print("reading description")
	}
	if isPrefix {
		log.Print("TODO")
	}
	var description = string(line)

	line, isPrefix, err = input.ReadLine()
	if err != nil {
		log.Print("reading blank line")
	}
	if isPrefix {
		log.Print("TODO")
	}

	for {
		line, isPrefix, err := input.ReadLine()
		if err != nil {
			break;
		}
		if isPrefix {
			log.Print("TODO")
		}
		var ingredient = string(line)
		if len(strings.TrimSpace(ingredient))==0 {
			break
		}
		ingredients = append(ingredients, ingredient)
	}

	for {
		line, isPrefix, err := input.ReadLine()
		if err != nil {
			break;
		}
		if isPrefix {
			log.Print("TODO")
		}
		var direction = string(line)
		if len(strings.TrimSpace(direction))==0 {
			break
		}

		line, isPrefix, err = input.ReadLine()
		if err != nil {
			//log.Print("reading blank line")
		}
		if isPrefix {
			log.Print("TODO")
		}

		directions = append(directions, direction)
	}

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
		log.Print("open", err)
	}
	dirs, err := f.Readdir(-1)
	if err != nil {
		log.Print("readdir", err)
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
	if recipe == nil {
		w.WriteHeader(http.StatusNotFound)
	}
	site_template.Execute(w, &page{Title: title, Recipe: recipe})
}


var addr = flag.String("addr", ":9999", "http service address")

func main() {
	flag.Parse()
	http.Handle("/", http.HandlerFunc(HomePageHandler))
	http.Handle("/recipes/", http.HandlerFunc(RecipeHandler))
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Print("ListenAndServe:", err)
	}
}
