package main

import (
	"os"
	"path"
	"log"
	"strings"
	"io/ioutil"
	"encoding/line"
)


type Recipe struct {
	Original string
	Name string
	Description string
	Ingredients []string
	Directions []string
	Photo string
}

func RecipeFromFile(filename string) *Recipe {
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

	var photo string
	for {
		line, isPrefix, err := input.ReadLine()
		if err != nil {
			break;
		}
		if isPrefix {
			log.Print("TODO")
		}
		var s = string(line)
		if strings.HasPrefix(s, "Photo:") {
			photo = strings.TrimSpace(strings.SplitAfter(s, "Photo:", 2)[1])
		}
	}

	return &Recipe{Name: path.Base(filename), Original: string(result), Description: description, Ingredients: ingredients, Directions: directions,
	Photo: photo}
}

func (r *Recipe) Id() string {
	return strings.ToLower(strings.Replace(r.Name, " ", "_", -1))
}


var Recipes = map[string]*Recipe{}


func init() {
	f, err := os.Open("recipes", os.O_RDONLY, 0)
	if err != nil {
		log.Print("open", err)
	}
	dirs, err := f.Readdir(-1)
	if err != nil {
		log.Print("readdir", err)
	}
	for _, d := range dirs {
		var recipe = RecipeFromFile(path.Join("./recipes/", d.Name))
		Recipes[recipe.Id()] = recipe
	}

	log.Print(Recipes)
}


type RecipeArray []*Recipe

func (p RecipeArray) Len() int           { return len(p) }
func (p RecipeArray) Less(i, j int) bool { return p[i].Name < p[j].Name }
func (p RecipeArray) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
