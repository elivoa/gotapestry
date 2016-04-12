package main

import (
	"fmt"
	"html/template"
	"os"
)

const (
	tempalte1 = "Template1 {{func1 $ `arg1`}}"
	template2 = "Template2 {{func1 $ `arg2`}}"
)

func func1(i interface{}, c string) string {
	// if err := parseTemplate("a1", overlay); err != nil {
	// 	panic(err)
	// }
	// Engine.template.ExecuteTemplate(os.Stderr, "a2", []string{"a", "b"})
	return "[[" + c + "]]"
}

func main() {
	rootTemplate := template.New("root")

	t, err := rootTemplate.Clone()
	if err != nil {
		panic(err)
	}

	rootTemplate.Funcs(template.FuncMap{"func1": func1})
	if _, err := t.Parse(tempalte1); err != nil {
		panic(err)
	}
	fmt.Println("============ start ===============")
	if err := t.ExecuteTemplate(os.Stderr, "root", []string{"a", "b"}); err != nil {
		panic(err)
	}
	fmt.Println("\n============================== DONE")

	// clone and use template 2

	// clone! parse 2

}
