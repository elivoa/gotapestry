package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

func main2() {
	const (
		master  = `Names:{{block "list" .}}{{"\n"}}{{range .}}{{println "-" .}}{{end}}{{end}}`
		overlay = `{{define "list"}} {{join . ", "}}{{end}} `
	)
	var (
		funcs     = template.FuncMap{"join": strings.Join}
		guardians = []string{"Gamora", "Groot", "Nebula", "Rocket", "Star-Lord"}
	)
	masterTmpl, err := template.New("master").Funcs(funcs).Parse(master)
	if err != nil {
		log.Fatal(err)
	}
	overlayTmpl, err := template.Must(masterTmpl.Clone()).Parse(overlay)
	if err != nil {
		log.Fatal(err)
	}
	if err := masterTmpl.Execute(os.Stdout, guardians); err != nil {
		log.Fatal(err)
	}
	if err := overlayTmpl.Execute(os.Stdout, guardians); err != nil {
		log.Fatal(err)
	}
}

// Engine instance. Unique.
var Engine = NewTemplateEngine()

type TemplateEngine struct {
	template *template.Template
}

func NewTemplateEngine() *TemplateEngine {
	e := &TemplateEngine{
		// init template TODO remove this, change another init method.
		// TODO: use better way to init.
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

	fmt.Printf("[parse tempalte] parseTempalte(%s,<<\n%s\n>>);\n", key, "--") //content) // REMOVE

	var tmpl *template.Template
	if Engine.template == nil {
		Engine.template = template.New(key)
	}

	if key == Engine.template.Name() {
		fmt.Println("[-vvvvvvvvvvvvvvvvvvvvvv") // REMOVE
		tmpl = Engine.template
	} else {
		fmt.Println("[-xxxxxxxxxxxxxxxxxxxx") // REMOVE
		tmpl = Engine.template.New(key)
		Engine.template = tmpl
	}
	fmt.Println("\tKEY    is", key)
	fmt.Println("\tt.name is", Engine.template.Name())
	fmt.Println("[-equals] template is: ", key == Engine.template.Name()) // REMOVE

	_, err := tmpl.Parse(content)
	fmt.Printf("[parse tempalte] End parseTempalte(%s, << ignored >>);\n", key) // REMOVE
	if err != nil {
		fmt.Println("[ERROR] \n", err) // REMOVE
		return err
	}
	return nil
}

func main() {
	fmt.Println("start test it")

	const (
		master  = `Names:{{block "list" .}}{{"\n"}}{{range .}}{{println "-" .}}{{end}}{{end}}`
		overlay = `{{define "list"}} {{join . ", "}}{{end}} `
	)
	var (
		funcs     = template.FuncMap{"join": strings.Join}
		guardians = []string{"Gamora", "Groot", "Nebula", "Rocket", "Star-Lord"}
	)
	fmt.Println("===========================")
	parseTemplate("ddd", master)
	fmt.Println("===========================")

	if err := Engine.template.Execute(os.Stdout, guardians); err != nil {
		log.Fatal(err)
	}
	parseTemplate("ddd", master)


	
	masterTmpl, err := template.New("master").Funcs(funcs).Parse(master)
	if err != nil {
		log.Fatal(err)
	}

	overlayTmpl, err := template.Must(masterTmpl.Clone()).Parse(overlay)
	if err != nil {
		log.Fatal(err)
	}
	if err := masterTmpl.Execute(os.Stdout, guardians); err != nil {
		log.Fatal(err)
	}
	if err := overlayTmpl.Execute(os.Stdout, guardians); err != nil {
		log.Fatal(err)
	}
}
