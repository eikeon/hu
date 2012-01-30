// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hu

import (
	"reflect"
	"strings"
	"testing"
)

type lexTest struct {
	name  string
	input string
	items []item
}

var (
	tEOF   = item{itemEOF, ""}
	tQuote = item{itemString, `"abc \n\t\" "`}
)

var lexTests = []lexTest{
	{"empty", "", []item{tEOF}},
	{"words", "Red lentil soup",
		[]item{
			{itemWord, "Red"}, {itemSpace, " "},
			{itemWord, "lentil"}, {itemSpace, " "},
			{itemWord, "soup"},
			tEOF}},
	{"number and word", "1 onion",
		[]item{
			{itemNumber, "1"}, {itemSpace, " "},
			{itemWord, "onion"},
			tEOF}},
	{"colon", "Photo: apple",
		[]item{
			{itemWord, "Photo"}, {itemPunctuation, ":"}, {itemSpace, " "},
			{itemWord, "apple"},
			tEOF}},
	{"punctuation", "onion, chopped",
		[]item{
			{itemWord, "onion"}, {itemPunctuation, ","}, {itemSpace, " "},
			{itemWord, "chopped"},
			tEOF}},
}

// collect gathers the emitted items into a slice.
func collect(t *lexTest) (items []item) {
	l := lex(t.name, strings.NewReader(t.input))
	for {
		item := l.nextItem()
		items = append(items, item)
		if item.typ == itemEOF || item.typ == itemError {
			break
		}
	}
	return
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		items := collect(&test)
		if !reflect.DeepEqual(items, test.items) {
			t.Errorf("%s: got\n\t%v\nexpected\n\t%v", test.name, items, test.items)
		}
	}
}
