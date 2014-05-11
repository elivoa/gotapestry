package main

import (
	"fmt"
	"html/template"
	"os"
	"time"
)

func main() {
	print("start")
	t := template.New("test template")

	t.Funcs(template.FuncMap{
		"testf": TestFunc,
		"now":   func() time.Time { return time.Now() },
	})

	t.Parse(`Test Templates[  {{testf (. | urlquery) }} ]`)

	err := t.Execute(os.Stdout, "Boutique Swe/ater Factory")
	if err != nil {
		panic(err)
	}
	print("end")
}

func TestFunc(a string) string {
	return fmt.Sprint("(((", a, ")))")
}
