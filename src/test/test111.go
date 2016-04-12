package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

func nestedFunc2(c string) string {

	// if err := parseTemplate("a2", overlay); err != nil {
	// 	panic(err)
	// }

	return "********" + c + "<<"
}

func main() {

	t := template.New("bbb")
	t.Funcs(template.FuncMap{"nested": nestedFunc})

	if err := parseTemplate("a1", master); err != nil {
		panic(err)
	}
	fmt.Println("============ start ===============")
	// if _, err := t.Parse(master); err != nil {
	// 	panic(err)
	// }

	t.ExecuteTemplate(os.Stderr, "a1", []string{"a", "b"})

	fmt.Println("\n========== DONE")
}

func main2() {
	fmt.Println("start test it")

	const (
		master  = `Names:{{block "list" .}}{{"\n"}}{{range .}}{{println "-" .}}{{end}}{{end}}`
		overlay = `{{define "list"}} sdfsdfsd {{end}} `
	)
	var (
		funcs     = template.FuncMap{"join": strings.Join}
		guardians = []string{"Gamora", "Groot", "Nebula", "Rocket", "Star-Lord"}
	)

	fmt.Println("===========================")
	Engine.template.Funcs(funcs)
	parseTemplate("ddd", master)

	fmt.Println(">>>>>")
	if err := Engine.template.ExecuteTemplate(os.Stdout, "ddd", guardians); err != nil {
		log.Fatal(err)
	}
	fmt.Println("<<<")

	parseTemplate("ddd", overlay)

	fmt.Println(">>>>>")
	if err := Engine.template.ExecuteTemplate(os.Stdout, "ddd", guardians); err != nil {
		log.Fatal(err)
	}
	fmt.Println("<<<")

	fmt.Println("\n\n\n===========================")

	// masterTmpl, err := template.New("master").Funcs(funcs).Parse(master)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// overlayTmpl, err := template.Must(masterTmpl.Clone()).Parse(overlay)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if err := masterTmpl.Execute(os.Stdout, guardians); err != nil {
	// 	log.Fatal(err)
	// }
	// if err := overlayTmpl.Execute(os.Stdout, guardians); err != nil {
	// 	log.Fatal(err)
	// }
}
