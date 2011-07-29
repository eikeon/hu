package main

import (
	"testing"
	"recipe"
	//"log"
	"strings"
	"dictionary"
	"fmt"
)


/*

 1 teaspoon grapeseed oil

 quantifier
 |- unit

 quantifier -> unit -> noun

 qualified ingredient: amount ingredient
 qualified ingredient: ingredient

 ingredient: Noun
   
 */

func pos(s string) []string {
	if (s=="1" || s=="2" || s=="3" || s=="4" || s=="5" ||
		s=="6" || s=="7" || s=="8" || s=="9" || s=="0") {
		return 	[]string{"quantity"}
	}
	return dictionary.PartsOfSpeech(s)
}

// find two word nouns
func find_two_word_nouns(in chan string) chan string {
	out := make(chan string, 100)
	go func() {
		var last string
		for i:= range in {
			if last!="" {
				var pair = last + " " + i
				if len(pos(pair))>0 {
					//fmt.Println(pair, "->", pos(pair))
					out <- pair
				}
			}
			last = i
		}
		close(out)
	}()
	return out
} 

func parse_qualified_ingredient(s string) {
        ch := make(chan string)
	out := find_two_word_nouns(ch)
	for _, word := range strings.Split(s, " ", -1) {
		//fmt.Println(pos(word))
		ch <- word;
	}
	close(ch)

	fmt.Print(s, " -> ")
	for r:= range out {
		fmt.Print("'", r, "'", " ")
	}
	fmt.Println()
}

func TestIngredientParse(t *testing.T) {
	for _, recipe := range recipe.Recipe_list {
		//fmt.Println(recipe)
		for _, ingredient_line := range recipe.Ingredients {
			//fmt.Println(ingredient_line)
			parse_qualified_ingredient(ingredient_line)
		}
	}
}

