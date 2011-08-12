package hu

import (
	"http"
	"log"
	"regexp"
	"io/ioutil"
	"path"
	"strings"
	"encoding/base64"
)

type PartOfSpeech struct {
	label string
	neighbors []string
}

type Grammer struct {
	label string
	neighbors [][]string
}

type Term struct {
	Label string
	PartOfSpeech string
	Definition []Term
}

// http://en.wikipedia.org/wiki/Word
//
// Thinking we want to include punctuation? For example,
//   http://en.wiktionary.org/w/api.php?action=parse&page=,&prop=wikitext&format=json
//
type Word string

func (w *Word) Terms() (terms []Term) {
	for _, p := range w.PartsOfSpeech() {
		terms = append(terms, Term{Label: string(*w), PartOfSpeech: p})
	}
	return terms
}

func (w *Word) PartsOfSpeech() []string {
	s := string(*w)
	if s == "1" || s == "2" || s == "3" || s == "4" || s == "5" ||
		s == "6" || s == "7" || s == "8" || s == "9" || s == "0" {
		return []string{"quantity"}
	}
	if s == " " {
		return []string{"space"}
	}
	return w.wiktionaryPartsOfSpeech()
}

func (w *Word) HasPartOfSpeech(pos string) bool {
	for _, p := range w.PartsOfSpeech() {
		if p == pos {
			return true
		}
	}
	return false
}

func (w *Word) getWikitext() string {
	word := string(*w)

	if len(word)==0 {
		return ""
	}
	// word_re := regexp.MustCompile("^[ a-zA-Z]+$")
	// if word_re.MatchString(word)==false {
	// 	return ""
	// }
	encoded_word := base64.StdEncoding.EncodeToString([]byte(word))
	cache_file := path.Join(".wiktionary", encoded_word)
	buf, err := ioutil.ReadFile(cache_file)
	if err!=nil {
		log.Print("Getting wikitext from wiktionary for: " + word, []byte(word))
		var _URL = "http://en.wiktionary.org/w/api.php"		
		URL := _URL + "?action=parse&page="+http.URLEscape(word)+"&prop=wikitext&format=json"
		r, err := http.Get(URL)
		if err != nil {
			log.Print(err)
		}
		b, err := ioutil.ReadAll(r.Body)
		write_err := ioutil.WriteFile(cache_file, b, 0666)
		if write_err != nil {
			log.Print("Error writing cache file for: " + word)
		}
		return string(b)
	}
	return string(buf)
}

func (w *Word) wiktionaryPartsOfSpeech() (pos []string) {
	//pos := make([]string, 0, 5)
	s := w.getWikitext()
	template := regexp.MustCompile("{{([^}\\|]+)(\\|([^}]+))?}}")
	templates := template.FindAllStringSubmatch(s, -1)
	for _, t := range templates {
		if strings.HasPrefix(t[1], "en-") {
			pos = append(pos, t[1])
		}
	}
	return pos
}
