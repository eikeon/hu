package hu

import (
	"io"
	"fmt"
)

type production struct {
	left  string
	right []string
}

func (p production) String() string {
	s := fmt.Sprintf("%v -> ", p.left)
	for _, right := range p.right {
		s += fmt.Sprintf("%v ", right)
	}
	return s
}

type Item struct {
	left    string
	right   []string
	dot     int
	start   int
	end     int
	parents [][]Item
}

func (a Item) equal(b Item) bool {
	if a.start != b.start || a.end != b.end || a.dot != b.dot {
		return false
	}
	if len(a.left) != len(b.left) {
		return false
	}
	if len(a.right) != len(b.right) {
		return false
	}
	for i := 0; i < len(a.left); i++ {
		if a.left[i] != b.left[i] {
			return false
		}
	}
	for i := 0; i < len(a.right); i++ {
		if a.right[i] != b.right[i] {
			return false
		}
	}
	return true
}

func (p Item) String() string {
	s := fmt.Sprintf("[%v -> ", p.left)
	for i, right := range p.right {
		if i == p.dot {
			s += "·"
		}
		s += fmt.Sprintf(" %v", right)
	}
	s += fmt.Sprintf(", %d, %d %v(%v)", p.start, p.end, len(p.parents), p.parents)
	//s += fmt.Sprintf(", %d, %d %v", p.start, p.end, len(p.parents))
	s += "]"
	return s
}

func (p Item) ParseTree(words []Word) string {
	s := ""
	s += p.left
	if p.end-p.start == 1 {
		s += fmt.Sprintf(" %v", words[p.start])
	}
	for _, parent := range p.parents {
		s += "["
		for i, item := range parent {
			s += fmt.Sprintf("%v", item.ParseTree(words))
			if i < len(parent)-1 {
				s += ","
			}
		}
		s += "]"
	}
	return s
}

func (p Item) Result(words []Word) string {
	s := fmt.Sprintf("%v:", p.left)
	for i := p.start; i < p.end; i++ {
		if i > p.start {
			s += fmt.Sprintf(" ")
		}
		s += fmt.Sprintf("%v", words[i])
	}
	//s += "]"
	return s
}

// Words can be put together to build larger elements of language,
// such as phrases (a red rock), clauses (I threw a rock), and
// sentences (he threw one too but he missed).

// A parser for taking words back apart.
type parser struct {
	wordList    []Word
	productions []production
	items       []Item
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
		if rune == ' ' {
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

func (p *parser) intern(b Item) (index int, existing bool) {
	for i, a := range p.items {
		if a.equal(b) {
			index = i
			existing = true
			return
		}
	}
	index = len(p.items)
	p.items = append(p.items, b)
	return
}

func (p *parser) debug_from(operation string, item Item, from Item) {
	item_index, existing := p.intern(item)
	if existing == false {
		fmt.Println("ERROR")
	}
	from_index, existing := p.intern(from)
	if existing == false {
		fmt.Println("ERROR")
	}
	if item_index == from_index && (item_index != 71 && item_index != 70) {
		return
	}
	fmt.Println(item_index, operation, ":", item, " from ", from_index)
}

func (p *parser) debug_with(operation string, item Item, a Item, b Item) {
	item_index, existing := p.intern(item)
	if existing == false {
		fmt.Println("ERROR")
	}
	a_index, existing := p.intern(a)
	if existing == false {
		fmt.Println("ERROR")
	}
	b_index, existing := p.intern(b)
	if existing == false {
		fmt.Println("ERROR")
	}
	fmt.Println(item_index, operation, ":", item, a_index, "with", b_index)
}

//70 predicted : [NP -> · Adjective Noun, 2, 2 0]  from  66
//71 predicted : [NP -> · Noun Noun, 2, 2 0]  from  66
//75 predicted (pos) : [Adjective -> · chicken, 2, 2 0]  from  70

func (p *parser) predict(i Item) (changed bool) {
	for _, production := range p.productions {
		//p.debug_from("predicted?", i, i)
		//fmt.Print("   ", production, production.left)
		//fmt.Println("  -->", i.dot, len(i.right), "<--", i.dot<len(i.right))
		// if i.dot<len(i.right) {
		// 	fmt.Println(fmt.Sprintf("'%v' '%v'\n", production.left, i.right[i.dot]))
		// }
		if i.dot < len(i.right) && production.left == i.right[i.dot] {
			item := Item{left: production.left, right: production.right, dot: 0, start: i.end, end: i.end}
			_, existing := p.intern(item)
			//p.items[index].parents = append(p.items[index].parents, []Item{i})
			if existing == false {
				//p.debug_from("predicted", item, i)
				changed = true
			}
		}
	}
	for _, word := range p.wordList {
		for _, partOfSpeech := range word.PartsOfSpeech() {
			pos := partOfSpeech
			if i.dot < len(i.right) && pos == i.right[i.dot] {
				item := Item{left: pos, right: []string{word.String()}, dot: 0, start: i.end, end: i.end}
				_, existing := p.intern(item)
				//p.items[index].parents = append(p.items[index].parents, []Item{i})
				if existing == false {
					//p.debug_from("predicted (pos)", item, i)
					changed = true
				}
			}
		}
	}
	return
}

func (p *parser) scan(i Item) (changed bool) {
	if i.dot < len(i.right) && i.start < len(p.wordList) && i.right[i.dot] == p.wordList[i.start].String() {
		item := Item{left: i.left, right: i.right, dot: i.dot + 1, start: i.start, end: i.end + 1}
		_, existing := p.intern(item)
		//p.items[index].parents = append(p.items[index].parents, []Item{i})
		if existing == false {
			//p.debug_from("scanned", item, i)
			changed = true
		}
	}
	return
}

func (p *parser) complete(b Item) (changed bool) {
	//for _, a := range(p.items) {
	for i := 0; i < len(p.items); i++ {
		a := p.items[i]
		//fmt.Println("item: ", b)
		if a.end == b.start && a.dot < len(a.right) && a.right[a.dot] == b.left && b.dot >= len(b.right) {
			item := Item{left: a.left, right: a.right, dot: a.dot + 1, start: a.start, end: b.end}
			index, existing := p.intern(item)
			new_parent := []Item{a, b}
			p.items[index].parents = append(p.items[index].parents, new_parent)
			if existing == false {
				//i = 0
				//p.debug_with("completed", p.items[index], a, b)
				changed = true
			}
		}
		// if b.end==a.start && b.dot<len(b.right) && b.right[b.dot]==a.left && b.dot>=len(b.right) {
		// 	item := Item{left: b.left, right: b.right, dot: b.dot+1, start: b.start, end: a.end}
		// 	index, existing:=p.intern(item)
		// 	new_parent := []Item{b, a}
		// 	p.items[index].parents = append(p.items[index].parents, new_parent)
		// 	if existing==false {
		// 		//i = 0
		// 		//p.debug_with("completed", p.items[index], b, a)
		// 		changed = true
		// 	}
		// }
	}
	return
}

func (p *parser) parse() (result [][]Term) {

	// print POS for debugging
	for _, words := range p.wordList {
		fmt.Print(words.PartsOfSpeech(), " ")
	}
	fmt.Println()

	// for _, word := range(p.wordList) {
	// 	for _, partOfSpeech := range(word.PartsOfSpeech()) {
	// 		pos := partOfSpeech
	// 		prod := production{left: pos, right: []string{word.String()}}
	// 		p.productions = append(p.productions, prod)
	// 	}
	// }

	for _, production := range p.productions {
		if production.left == "S" {
			item := Item{left: production.left, right: production.right, dot: 0, start: 0, end: 0}
			p.intern(item)
		}
	}
	for i := 0; i < len(p.items); i++ {
		item := p.items[i]
		//fmt.Println(i)
		for {
			if (p.predict(item) || p.scan(item) || p.complete(item)) == false {
				break
			}
			//i = 0
		}
	}
	for _, item := range p.items {
		if item.start == 0 && item.end == len(p.wordList) && item.left == "S" {
			fmt.Println(":", item)
			//var items []Item
			//items = append(items, item)
			// out := make(chan string, 100)
			// go p.Result(items, out, 0)
			// for r:= range(out) {
			// 	fmt.Println(r)
			// }
		}
	}
	return
}

func (p *parser) Result(items []Item, out chan string, depth int) {
	leaf := true
	ss := "["
	for i, item := range items {
		//fmt.Println(len(item.parents))
		//ss += item.Result(p.wordList)
		ss += fmt.Sprintf("%v (%v,%v)", item, item.dot, len(item.parents))
		// if item.dot!=0 {
		// 	leaf = false
		// }
		//fmt.Printf("%v (%v)\n", item, item.dot)
		if item.dot == 1 && len(item.right) != 1 {
			leaf = false
		}
		if item.dot == 1 && len(item.right) == 1 && len(item.parents) > 0 {
			leaf = false
		}
		//if item.dot<=len(item.right) && len(item.parents)==0 {
		//    leaf = false
		//}
		if item.end-item.start != item.dot || item.dot > 1 {
			leaf = false
		}
		if i < len(items)-1 {
			ss += ", "
		}
	}
	ss += "]"
	if leaf == true {
		out <- ss
	}
	found := false
	for i, item := range items {
		if found == false && len(item.parents) > 0 {
			found = true
			for _, parent := range item.parents {
				var nitems []Item
				for j := 0; j < i; j++ {
					if items[j].end-items[j].start > 0 {
						nitems = append(nitems, items[j])
					}
				}
				for _, pitem := range parent {
					nitems = append(nitems, pitem)
				}
				for j := i + 1; j < len(items); j++ {
					nitems = append(nitems, items[j])
				}
				fmt.Println(items, "==>", nitems)
				if depth < 2 {
					p.Result(nitems, out, depth+1)
				}
			}
		}
	}

	if depth == 0 {
		close(out)
	}
}
