package main

import "github.com/ParadoxPixel/Pargen/parser"

func main() {
	gen, err := NewGen(
		"./files/in",
		"./files/parsed",
		true,
		parser.DefaultRegex,
		func(s string) string {
			return s
		},
	)
	if err != nil {
		panic(err)
	}

	if err = gen.Initialize(); err != nil {
		panic(err)
	}

	templates, err := gen.Load(nil)
	if err != nil {
		panic(err)
	}

	err = templates.Lookup("box").Execute(nil, nil)
	if err != nil {
		panic(err)
	}
}
