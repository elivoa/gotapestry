package main

import (
	"fmt"
	"html/template"
	"os"
)

// Engine instance. Unique.
var Engine = NewTemplateEngine()

type TemplateEngine struct {
	template *template.Template
}

func NewTemplateEngine() *TemplateEngine {
	e := &TemplateEngine{
		template: template.New("-"),
	}
	return e
}

func parseTemplate(key string, content string) error {
	// Old version uses filename as key, I make my own key. not
	// filepath.Base(filename) First template becomes return value if
	// not already defined, we use that one for subsequent New
	// calls to associate all the templates together. Also, if this
	// file has the same name as t, this file becomes the contents of
	// t, so t, err := New(name).Funcs(xxx).ParseFiles(name)
	// works. Otherwise we create a new template associated with t.

	fmt.Println("\n\n--------------------------------------------------------------------------------")
	fmt.Printf("[parse tempalte] parseTempalte(%s,<<%s>>);\n", key, content) //content) // REMOVE

	engine := Engine

	if engine.template == nil {
		engine.template = template.New(key)
	}

	var tmpl *template.Template
	// fmt.Println("0000000000000 > ", key, engine.template.Name())
	if key == engine.template.Name() {
		// fmt.Println(".... User old name, ", engine.template.Name())
		tmpl = engine.template
	} else {
		// fmt.Println(".... New name, ", key)
		tmpl = engine.template.New(key)
	}

	if true { // -------------------------- debug print templates.
		fmt.Println("\ndebug info { // templates loop ; tmpl.name is : ", tmpl.Name())
		for _, t := range engine.template.Templates() {
			fmt.Println("  ", t.Name())
		}
		fmt.Println("}")
	}

	// newt, e := tmpl.Clone()
	// if e != nil {
	// 	panic(e)
	// }
	// _, err := newt.Parse(content)
	_, err := tmpl.Parse(content)

	// fmt.Printf("[parse tempalte] End parseTempalte(%s, << ignored >>);\n", key) // REMOVE
	if err != nil {
		// fmt.Println("[ERROR] : \t", err) // REMOVE
		return err
	}
	engine.template = tmpl
	return nil
}

const (
	master = `<head><title>Title</title>` + "{{t_PageHeadBootstrap $ `span`}}</head>" + `<body></html>`

	overlay = `(____PageHeadBootstrap_replace_to_html____) `
)

func nestedFunc(i interface{}, c string) string {
	if err := parseTemplate("a1", overlay); err != nil {
		panic(err)
	}
	Engine.template.ExecuteTemplate(os.Stderr, "a2", []string{"a", "b"})
	return "********" + c + "<<"
}

func main() {

	t := Engine.template
	t.Funcs(template.FuncMap{"t_PageHeadBootstrap": nestedFunc})

	if err := parseTemplate("p/test.Test", master); err != nil {
		panic(err)
	}

	fmt.Println("============ start ===============")
	if err := t.ExecuteTemplate(os.Stderr, "p/test.Test", []string{"a", "b"}); err != nil {
		panic(err)
	}
	fmt.Println("\n========== DONE")
}

// func main2() {

// 	t := templates2.Engine
// 	t.Template().Funcs(template.FuncMap{"t_PageHeadBootstrap": nestedFunc})
// 	if err := templates2.ParseTemplate("p/test.Test", master); err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("============ Execute ===============")
// 	if err := t.Template().ExecuteTemplate(os.Stderr, "p/test.Test", []string{"a", "b"}); err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("\n========== DONE")
// }
