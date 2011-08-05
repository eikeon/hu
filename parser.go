package hu

import (
	"io"
)

// Words can be put together to build larger elements of language,
// such as phrases (a red rock), clauses (I threw a rock), and
// sentences (he threw one too but he missed).

// A parser for taking words back apart.
type parser struct {
	reader io.RuneScanner
	word Word
	position int
}

//func NewParser(reader io.RuneReader) *parser { return &Parser{} }

// Advance to next word 
func (p *parser) next() {
	var s Word = ""
	rune, _, err := p.reader.ReadRune()
	if err != nil {
		p.word = s
		return
	}
	s = Word(string(rune))
	if rune==' ' {
		p.word = s
		return 
	}
	for {
		rune, size, err := p.reader.ReadRune()
		p.position += size
		if err != nil {
			p.word = s
			return
		}
		if rune == ',' || rune == ' ' {
			p.word = s
			p.reader.UnreadRune()
			return
		}
		s = s + Word(string(rune))
	}
}

func (p *parser) parseWordList() (list []Word) {
	for {
		p.next()
		w := p.word
		if w=="" {
			break
		}
		list = append(list, w)
	}
	return
}

type PartOfSpeech struct {
	label string
	neighbors []string
}

func (p *parser) parseAs(pos PartOfSpeech) (list []Word) {
    for i := 0; i < len(pos.neighbors); i++ {
		p.next()
		w := p.word
		if w.HasPartOfSpeech(pos.neighbors[i])==false {
			list = list[0:0]
			return
		} else {
			list = append(list, w)
		}
	}
	p.next()
	if p.word!="" { // we didn't reach the end
		list = list[0:0]
	}
	return
}
