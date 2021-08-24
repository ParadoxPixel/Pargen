package main

import (
	"fmt"
	"github.com/ParadoxPixel/Pargen/parser"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Gen struct {
	in        string
	inAbs     string
	parsed	  string
	parsedAbs string
	update    bool
	Parser    *parser.Parser
}

func NewDefaultGen(in , parsed string, f func(string)string) (*Gen, error) {
	return NewGen(
		in,
		parsed,
		false,
		parser.DefaultRegex,
		f,
	)
}

func NewGen(in , parsed string, update bool, regex *regexp.Regexp, f func(string)string) (*Gen, error) {
	inAbs, err := filepath.Abs(in)
	if err != nil {
		return nil, err
	}

	parsedAbs, err := filepath.Abs(parsed)
	if err != nil {
		return nil, err
	}

	return &Gen{
		in:        in,
		inAbs:     inAbs,
		parsed:    parsed,
		parsedAbs: parsedAbs,
		update:    update,
		Parser:    parser.NewParser(regex, f),
	}, nil
}

func(g *Gen) Initialize() error {
	var err error
	if err = g.Prepare(); err != nil {
		return err
	}

	return g.ParseAll()
}

func(g *Gen) Prepare() error {
	var err error
	if _, err = os.Stat(g.inAbs); os.IsNotExist(err) {
		if err = os.MkdirAll(g.inAbs, 0700); err != nil {
			return err
		}
	}

	if _, err = os.Stat(g.parsedAbs); os.IsNotExist(err) {
		return os.MkdirAll(g.parsedAbs, 0700)
	}

	return nil
}

func(g *Gen) ParseAll() error {
	return filepath.Walk(g.inAbs, func(fileName string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return err
		}

		return g.ParseFile(fileName)
	})
}

func(g *Gen) ParseFile(fileName string) error {
	fileName, err := filepath.Abs(fileName)
	if err != nil {
		return err
	}

	fileName = strings.Replace(fileName, g.inAbs, "", 1)
	fileName = strings.Trim(fileName, string(os.PathSeparator))
	dir, fileName := filepath.Split(fileName)

	parsed := filepath.Join(g.parsedAbs, dir)
	if _, err = os.Stat(parsed); os.IsNotExist(err) {
		if err = os.MkdirAll(parsed, 0700); err != nil {
			return err
		}
	} else if !g.update {
		if _, err = os.Stat(filepath.Join(parsed, fileName)); !os.IsNotExist(err) {
			return err
		}
	}

	in := filepath.Join(g.inAbs, dir)
	bytes, err := ioutil.ReadFile(filepath.Join(in, fileName))
	if err != nil {
		return err
	}

	str := string(bytes)
	str = g.Parser.Parse(str)
	if strings.HasSuffix(fileName, ".temp.html") {
		fmt.Println(filepath.Join(dir, strings.TrimSuffix(fileName, ".temp.html")))
		str = "{{define \"" + filepath.Join(dir, strings.TrimSuffix(fileName, ".temp.html")) + "\"}}\n" + str + "\n{{end}}"
	}


	bytes = []byte(str)
	f, err := os.OpenFile(filepath.Join(parsed, fileName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(f, "%s", bytes)
	if err != nil {
		return err
	}

	return f.Close()
}

func(g *Gen) Load(fm template.FuncMap) (*template.Template, error) {
	templates := template.New("")
	if fm != nil {
		templates.Funcs(fm)
	}

	err := filepath.Walk(g.parsedAbs, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info == nil || info.IsDir() {
			return nil
		}

		if 5 >= len(path) || path[len(path) - 5:] != ".html" {
			return nil
		}

		rel, err := filepath.Rel(g.parsedAbs, path)
		if err != nil {
			return err
		}

		name := filepath.ToSlash(rel)
		name = strings.TrimSuffix(name, ".html")
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		_, err = templates.New(name).Parse(string(buf))
		return err
	})

	return templates, err
}