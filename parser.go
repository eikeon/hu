package hu

import (
	"io"
	"fmt"
)

type production struct {
	left string
	right []string
}

func (p production) String() string {
	s := fmt.Sprintf("%v -> ", p.left)
	for _, right:= range(p.right) {
		s += fmt.Sprintf("%v ", right)
	}
	return s
}

type Item struct {
	left string
	right []string
	dot int
	start int
	end int
	parents [][]Item
}

func (a Item) equal(b Item) bool {
	if a.start!=b.start || a.end!=b.end || a.dot!=b.dot {
		return false
	}
	if len(a.left)!=len(b.left) {
		return false
	}
	for i:=0; i<len(a.left); i++ {
		if a.left[i]!=b.left[i] {
			return false
		}
	}
	for i:=0; i<len(a.right); i++ {
		if a.right[i]!=b.right[i] {
			return false
		}
	}
	return true
}

func (p Item) String() string {
	s := fmt.Sprintf("%v -> ", p.left)
	for i, right:= range(p.right) {
		if i==p.dot {
			s += "·"
		}
		s += fmt.Sprintf(" %v", right)
	}
	s += fmt.Sprintf(", %d, %d", p.start, p.end)
	return s
}

func (p Item) ParseTree(words []Word) string {
	s := ""
	s += p.left
	if p.end - p.start == 1 {
		s += fmt.Sprintf(" %v", words[p.start])
	}
	for _, parent:= range(p.parents) {
		s += "["
		for i, item := range(parent) {
			s += fmt.Sprintf("%v", item.ParseTree(words))
			if i<len(parent)-1 {
				s += ","
			}
		}
		s += "]"
	}
	return s
}

// Words can be put together to build larger elements of language,
// such as phrases (a red rock), clauses (I threw a rock), and
// sentences (he threw one too but he missed).

// A parser for taking words back apart.
type parser struct {
	wordList []Word
	productions []production
	items []Item
}

func NewParser(reader io.RuneScanner) *parser {
	p := &parser{}
	for {
		var s Word = Word("")
		rune, _, err := reader.ReadRune()
		if err != nil {
			//p.wordList = append(p.wordList, s)
			break
		}
		s = Word(string(rune))
		if rune==' ' {
			//p.wordList = append(p.wordList, s)
			continue
		}
		for {
			rune, _, err := reader.ReadRune()
			//p.position += size
			if err != nil {
				p.wordList = append(p.wordList, s)
				break
			}
			if rune == ',' || rune == ' ' {
				p.wordList = append(p.wordList, s)
				reader.UnreadRune()
				break
			}
			s = s + Word(string(rune))
		}
	}
	return p
}

func (p *parser) contains(b Item) bool {
	for _, a := range(p.items) {
		if a.equal(b) {
			return true
		}
	}
	return false
}

func (p *parser) predict(i Item) (items []Item) {
	for _, production := range(p.productions) {
		if i.dot<len(i.right) && production.left==i.right[i.dot] {
			item := Item{left: production.left, right: production.right, dot: 0, start: i.end, end: i.end}
			if p.contains(item)==false {
				//fmt.Println("predicted: ", item)
				p.items = append(p.items, item)
			}
		}
	}
	for _, word := range(p.wordList) {
		for _, partOfSpeech := range(word.PartsOfSpeech()) {
			pos := partOfSpeech
			//if partOfSpeech=="en-noun" {
			//	pos = "NP"
			//}
			if i.dot<len(i.right) && pos==i.right[i.dot] {
				item := Item{left: pos, right: []string{word.String()}, dot: 0, start: i.end, end: i.end}
				if p.contains(item)==false {
					//fmt.Println("predicted (pos): ", item)
					p.items = append(p.items, item)
				}
			}
		}
	}
	return
}

func (p *parser) scan(i Item) (items []Item) {
	if i.dot<len(i.right) && i.start<len(p.wordList) && i.right[i.dot]==p.wordList[i.start].String() {
		item := Item{left: i.left, right: i.right, dot: i.dot+1, start: i.start, end: i.end+1}
		if p.contains(item)==false {
			//fmt.Println("scanned: ", item)
			p.items = append(p.items, item)
		}
	}
	return
}

func (p *parser) complete(b Item) (items []Item) {
	for _, a := range(p.items) {
		if a.end==b.start && a.dot<len(a.right) && a.right[a.dot]==b.left {
			item := Item{left: a.left, right: a.right, dot: a.dot+1, start: a.start, end: b.end}
			item.parents = append(item.parents, []Item{a, b})
			if p.contains(item)==false {
				//fmt.Println("completed: ", item)
				p.items = append(p.items, item)
			}
		}
	}
	return
}

func (p *parser) parse() (result [][]Term) {

	// print POS for debugging
	for _, words := range(p.wordList) {
		fmt.Print(words.PartsOfSpeech(), " ")
	}
	fmt.Println()

	for _, production := range(p.productions) {
		if production.left=="S" {
			item := Item{left: production.left, right: production.right, dot: 0, start: 0, end: 0}
			p.items = append(p.items, item)
		}
	}
	for i:=0; i<len(p.items); i++ {
		item := p.items[i]
		p.predict(item)
		p.scan(item)
		p.complete(item)
	}
	for _, item := range(p.items) {
		if item.start==0 && item.end==len(p.wordList) {
			fmt.Println(":", item)
			fmt.Println(item.ParseTree(p.wordList))
		}
	}
	return
}
