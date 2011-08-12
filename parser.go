package hu

import (
	"io"
)

// Words can be put together to build larger elements of language,
// such as phrases (a red rock), clauses (I threw a rock), and
// sentences (he threw one too but he missed).

// A parser for taking words back apart.
type parser struct {
	wordList []Word
}

func NewParser(reader io.RuneScanner) *parser {
	p := &parser{}
	for {
		var s Word = ""
		rune, _, err := reader.ReadRune()
		if err != nil {
			//p.wordList = append(p.wordList, s)
			break
		}
		s = Word(string(rune))
		if rune==' ' {
			p.wordList = append(p.wordList, s)
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

func (p *parser) parseAs(pos PartOfSpeech) (list []Term) {
	words := p.wordList
    for i := 0; i < len(pos.neighbors); i++ {
		w := words[i]
		match := false
		for _, t:= range(w.Terms()) {
			if t.PartOfSpeech==pos.neighbors[i] {
				list = append(list, t)
				match = true
				break
			}
		}
		if match==false {
			break
		}
	}
	if len(words)!=len(list) {
		list = list[0:0]
	}
	return
}
