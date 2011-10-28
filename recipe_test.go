// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hu

import (
	"flag"
	"fmt"
	"testing"
)

var debug = flag.Bool("debug", false, "show the errors produced by the tests")

type parseTest struct {
	name   string
	input  string
	ok     bool
	result string
}

const (
	noError  = true
	hasError = false
)

var parseTests = []parseTest{
	//	{"empty", "", noError, `[]`},
	{"basic", "Red Lentils with Lemon\n\nQuick, easy, and delicious.\n\n3 tablespoons olive oil, plus more for drizzling\n1 large onion, chopped\n2 cloves garlic, minced\n1 tablespoon tomato paste\n1 teaspoon ground cumin\n1 teaspoon kosher salt, or more to taste\n1 teaspoon freshly ground black pepper\nlarge pinch of cayenne pepper, or more to taste\n4 cups chicken or vegetable broth\n2 cups water\n1 cup red lentils, picked over and rinsed\n1 large carrot, peeled and diced\njuice of 1 lemon, or more to taste\n3 tablespoons cilantro, chopped\n\nHeat the oil in a large pot over high heat until hot and shimmering. Add the onion and garlic and sauté until golden, about 4 minutes. Stir in the tomato paste, cumin, salt, black pepper, and chile powder and sauté for 2 minutes. Add the broth, water, lentils, and carrot and bring to a simmer, then partially cover the pot and turn the heat to medium low. Simmer until the lentils are soft, about 30 minutes. Taste and add salt if necessary.\n\nUsing an immersion or regular blender or a food processor, puree half the soup, then return it to the pot; do not over puree it, the soup should be somewhat chunky. Reheat the soup if necessary, then stir in the lemon juice and cilantro. Serve drizzled with olive oil and dusted lightly with chile powder if desired.\n\n", noError,
		`[]`},
	// {"comment", "{{/*\n\n\n*/}}", noError,
	// 	`[]`},
	// {"spaces", " \t\n", noError,
	// 	`[(text: " \t\n")]`},
	// {"text", "some text", noError,
	// 	`[(text: "some text")]`},
	// {"emptyAction", "{{}}", hasError,
	// 	`[(action: [])]`},
	// {"field", "{{.X}}", noError,
	// 	`[(action: [(command: [F=[X]])])]`},
	// {"simple command", "{{printf}}", noError,
	// 	`[(action: [(command: [I=printf])])]`},
	// {"$ invocation", "{{$}}", noError,
	// 	"[(action: [(command: [V=[$]])])]"},
	// {"variable invocation", "{{with $x := 3}}{{$x 23}}{{end}}", noError,
	// 	"[({{with [V=[$x]] := [(command: [N=3])]}} [(action: [(command: [V=[$x] N=23])])])]"},
	// {"variable with fields", "{{$.I}}", noError,
	// 	"[(action: [(command: [V=[$ I]])])]"},
	// {"multi-word command", "{{printf `%d` 23}}", noError,
	// 	"[(action: [(command: [I=printf S=`%d` N=23])])]"},
	// {"pipeline", "{{.X|.Y}}", noError,
	// 	`[(action: [(command: [F=[X]]) (command: [F=[Y]])])]`},
	// {"pipeline with decl", "{{$x := .X|.Y}}", noError,
	// 	`[(action: [V=[$x]] := [(command: [F=[X]]) (command: [F=[Y]])])]`},
	// {"declaration", "{{.X|.Y}}", noError,
	// 	`[(action: [(command: [F=[X]]) (command: [F=[Y]])])]`},
	// {"simple if", "{{if .X}}hello{{end}}", noError,
	// 	`[({{if [(command: [F=[X]])]}} [(text: "hello")])]`},
	// {"if with else", "{{if .X}}true{{else}}false{{end}}", noError,
	// 	`[({{if [(command: [F=[X]])]}} [(text: "true")] {{else}} [(text: "false")])]`},
	// {"simple range", "{{range .X}}hello{{end}}", noError,
	// 	`[({{range [(command: [F=[X]])]}} [(text: "hello")])]`},
	// {"chained field range", "{{range .X.Y.Z}}hello{{end}}", noError,
	// 	`[({{range [(command: [F=[X Y Z]])]}} [(text: "hello")])]`},
	// {"nested range", "{{range .X}}hello{{range .Y}}goodbye{{end}}{{end}}", noError,
	// 	`[({{range [(command: [F=[X]])]}} [(text: "hello")({{range [(command: [F=[Y]])]}} [(text: "goodbye")])])]`},
	// {"range with else", "{{range .X}}true{{else}}false{{end}}", noError,
	// 	`[({{range [(command: [F=[X]])]}} [(text: "true")] {{else}} [(text: "false")])]`},
	// {"range over pipeline", "{{range .X|.M}}true{{else}}false{{end}}", noError,
	// 	`[({{range [(command: [F=[X]]) (command: [F=[M]])]}} [(text: "true")] {{else}} [(text: "false")])]`},
	// {"range []int", "{{range .SI}}{{.}}{{end}}", noError,
	// 	`[({{range [(command: [F=[SI]])]}} [(action: [(command: [{{<.>}}])])])]`},
	// {"constants", "{{range .SI 1 -3.2i true false 'a'}}{{end}}", noError,
	// 	`[({{range [(command: [F=[SI] N=1 N=-3.2i B=true B=false N='a'])]}} [])]`},
	// {"template", "{{template `x`}}", noError,
	// 	`[{{template "x"}}]`},
	// {"template with arg", "{{template `x` .Y}}", noError,
	// 	`[{{template "x" [(command: [F=[Y]])]}}]`},
	// {"with", "{{with .X}}hello{{end}}", noError,
	// 	`[({{with [(command: [F=[X]])]}} [(text: "hello")])]`},
	// {"with with else", "{{with .X}}hello{{else}}goodbye{{end}}", noError,
	// 	`[({{with [(command: [F=[X]])]}} [(text: "hello")] {{else}} [(text: "goodbye")])]`},
	// // Errors.
	// {"unclosed action", "hello{{range", hasError, ""},
	// {"unmatched end", "{{end}}", hasError, ""},
	// {"missing end", "hello{{range .x}}", hasError, ""},
	// {"missing end after else", "hello{{range .x}}{{else}}", hasError, ""},
	// {"undefined function", "hello{{undefined}}", hasError, ""},
	// {"undefined variable", "{{$x}}", hasError, ""},
	// {"variable undefined after end", "{{with $x := 4}}{{end}}{{$x}}", hasError, ""},
	// {"variable undefined in template", "{{template $v}}", hasError, ""},
	// {"declare with field", "{{with $x.Y := 4}}{{end}}", hasError, ""},
	// {"template with field ref", "{{template .X}}", hasError, ""},
	// {"template with var", "{{template $v}}", hasError, ""},
	// {"invalid punctuation", "{{printf 3, 4}}", hasError, ""},
	// {"multidecl outside range", "{{with $v, $u := 3}}{{end}}", hasError, ""},
	// {"too many decls in range", "{{range $u, $v, $w := 3}}{{end}}", hasError, ""},
}

var builtins = map[string]interface{}{
	"printf": fmt.Sprintf,
}

func TestRecipeParse(t *testing.T) {
	for _, test := range parseTests {
		tmpl, err := New(test.name).Parse(test.input, builtins)
		switch {
		case err == nil && !test.ok:
			t.Errorf("%q: expected error; got none", test.name)
			continue
		case err != nil && test.ok:
			t.Errorf("%q: unexpected error: %v", test.name, err)
			continue
		case err != nil && !test.ok:
			// expected error, got one
			if *debug {
				fmt.Printf("%s: %s\n\t%s\n", test.name, test.input, err)
			}
			continue
		}
		result := tmpl.Recipe.String()
		if result != test.result {
			t.Errorf("%s=(%q): got\n\t%v\nexpected\n\t%v", test.name, test.input, result, test.result)
		}
	}
}
