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
	"crypto/md5"
	"fmt"
	"bufio"
	//"compress/gzip"
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
	Name string
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

	return &recipe{Name: path.Base(filename), Original: string(result), Description: description, Ingredients: ingredients, Directions: directions}
}

func (r *recipe) Id() string {
	return strings.ToLower(strings.Replace(r.Name, " ", "_", -1))
}


var recipes = map[string]*recipe{}

func init() {
	f, err := os.Open("recipes", os.O_RDONLY, 0)
	if err != nil {
		log.Print("open", err)
	}
	dirs, err := f.Readdir(-1)
	if err != nil {
		log.Print("readdir", err)
	}
	//names := make([]string, len(dirs))
	for _, d := range dirs {
		var recipe = RecipeFromFile(path.Join("./recipes/", d.Name))
		recipes[recipe.Id()] = recipe
	}
	log.Print(recipes)
}


type page struct {
	Title string
	NotFound bool
	Recipes []*recipe
	Recipe *recipe
}

type RecipeArray []*recipe

func (p RecipeArray) Len() int           { return len(p) }
func (p RecipeArray) Less(i, j int) bool { return p[i].Name < p[j].Name }
func (p RecipeArray) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }


func HomePageHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		w.SetHeader("Cache-Control", "max-age=10, must-revalidate")
		w.WriteHeader(http.StatusNotFound)
		site_template.Execute(w, &page{NotFound: true})
		return
	}

	recipe_list := make(RecipeArray, len(recipes))
	var i int
	for _, recipe := range recipes {
		recipe_list[i] = recipe
		i += 1
	}
	sort.Sort(recipe_list)

	bw := bufio.NewWriter(nil)
	h := md5.New()
	mw := io.MultiWriter(bw, h)
	site_template.Execute(mw, &page{Title: "Our Recipes", Recipes: recipe_list})

	//w.SetHeader("Vary", "Accept-Encoding")
	w.SetHeader("Cache-Control", "max-age=1, must-revalidate")
	w.SetHeader("ETag", fmt.Sprintf("\"%x\"", h.Sum()))
	site_template.Execute(w, &page{Title: "Our Recipes", Recipes: recipe_list})
}

func RecipeHandler(w http.ResponseWriter, req *http.Request) {
	var recipe = recipes[path.Base(req.URL.Path)]
	if recipe == nil {
		w.WriteHeader(http.StatusNotFound)
		site_template.Execute(w, &page{NotFound: true})
		return
	}
	var title = recipe.Name
	// w.SetHeader("Content-Encoding", "gzip")
	// ww, err := gzip.NewWriter(w)
	// if err != nil {
	// 	log.Print("gzip", err)
	// }

	site_template.Execute(w, &page{Title: title, Recipe: recipe})
}

func StaticHandler(w http.ResponseWriter, req *http.Request) {
	w.SetHeader("Content-Type", "text/css")
	w.SetHeader("Cache-Control", "max-age=3153600")

	var filename = path.Join(".", "/static/page.css")

	f, err := os.Open(filename, os.O_RDONLY, 0)
	if err != nil {
		log.Print("open", err)
	}

	// w.SetHeader("Content-Encoding", "gzip")
	// ww, err := gzip.NewWriter(w)
	// if err != nil {
	// 	log.Print("gzip", err)
	// }

	var input = line.NewReader(f, 1024)
	io.Copy(w, input)
}

var addr = flag.String("addr", ":9999", "http service address")

func main() {
	flag.Parse()
	http.Handle("www.eikeon.com/", http.RedirectHandler("http://eikeon.com/", http.StatusMovedPermanently))
	http.Handle("/", http.HandlerFunc(HomePageHandler))
	http.Handle("/recipes/", http.HandlerFunc(RecipeHandler))
	http.Handle("/37c58cfd8b54050a0f34d8fe5ebc2ff0_page.css", http.HandlerFunc(StaticHandler))
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Print("ListenAndServe:", err)
	}
}
