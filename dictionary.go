package hu

import (
	"http"
	"log"
	"regexp"
	"io/ioutil"
	"path"
	"encoding/base64"
	"url"
)

type Term struct {
	Label        string
	PartOfSpeech string
	Definition   []Term
}

var posCache = map[string][]string{}

// http://en.wikipedia.org/wiki/Word
//
// Thinking we want to include punctuation? For example,
//   http://en.wiktionary.org/w/api.php?action=parse&page=,&prop=wikitext&format=json
//
type Word string

func (w *Word) Terms() (terms []Term) {
	for _, p := range w.PartsOfSpeech() {
		terms = append(terms, Term{Label: w.String(), PartOfSpeech: p})
	}
	return terms
}

func (w *Word) PartsOfSpeech() []string {
	v, ok := posCache[w.String()]
	if ok {
		return v
	}
	s := w.String()
	if s == "1" || s == "2" || s == "3" || s == "4" || s == "5" ||
		s == "6" || s == "7" || s == "8" || s == "9" || s == "0" {
		return []string{"Quantity"}
	}
	if s == " " {
		return []string{"Space"}
	}
	if s == "," {
		return []string{"Comma"}
	}
	pos := w.wiktionaryPartsOfSpeech()
	posCache[w.String()] = pos
	return pos
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
	word := w.String()

	if len(word) == 0 {
		return ""
	}
	// word_re := regexp.MustCompile("^[ a-zA-Z]+$")
	// if word_re.MatchString(word)==false {
	// 	return ""
	// }
	encoded_word := base64.StdEncoding.EncodeToString([]byte(word))
	cache_file := path.Join(".wiktionary", encoded_word)
	buf, err := ioutil.ReadFile(cache_file)
	if err != nil {
		log.Print("Getting wikitext from wiktionary for: "+word, []byte(word))
		var _URL = "http://en.wiktionary.org/w/api.php"
		URL := _URL + "?action=parse&page=" + url.QueryEscape(word) + "&prop=wikitext&format=json"
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
	english_level := 0
	in_english := false
	header := regexp.MustCompile("(==+)([^= ]+)(==+)")
	for _, t := range header.FindAllStringSubmatch(s, -1) {
		//fmt.Println(t)
		if t[1] == t[3] {
			if in_english && len(t[1]) > english_level {
				v := t[2]
				if v != "Quotations" &&
					v != "Etymology" &&
					v != "Pronunciation" &&
					v != "Translations" &&
					v != "Statistics" &&
					v != "References" &&
					v != "Synonyms" &&
					v != "Anagrams" &&
					v != "Antonyms" {
					contains := false
					for _, vv := range pos {
						if v == vv {
							contains = true
							break
						}
					}
					if contains == false {
						pos = append(pos, v)
					}
				}
			}
			if in_english && len(t[1]) == english_level {
				in_english = false
			}
			if t[2] == "English" {
				english_level = len(t[1])
				in_english = true
			}
		}
	}
	return pos
}

func (w Word) String() string {
	return string(w)
}
