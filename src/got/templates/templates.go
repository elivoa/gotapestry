/*
   Time-stamp: <[templates.go] Elivoa @ Tuesday, 2013-07-30 13:01:54>
*/
package templates

import (
	"bufio"
	"errors"
	"fmt"
	"got/debug"
	"got/templates/transform"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"sync"
)

// all templates are in one template.Tempalte.
var Templates *template.Template

func init() {
	// init template TODO remove this, change another init method.
	// TODO: use better way to init.
	// Templates = template.Must(template.ParseFiles(
	// 	LocatePath("nothing"),
	// ))
	Templates = template.New("Got-Template")
	registerBuiltinFuncs()   // Register built-in templates.
	registerComponentFuncs() // no use?
}

/*
 @return
   template - when tempalte available and parse successful.
   nil      - when template not exists or error occurs.
*/
func AddGOTTemplate(key string, filename string) (*template.Template, error) {

	debug.Log("-   - [ParseTempalte] %v, %v", key, filename)

	// borrowed from html/tempate
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, nil // file not exist, don't panic.
	}

	// open input file
	fi, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	// close fi on exit and check for its returned error
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()
	// make a read buffer
	r := bufio.NewReader(fi)

	// transform
	trans := transform.NewTransformer()
	html := trans.Parse(r).Render()
	// fmt.Println("--------------------------------------------------------------------------------------")
	// fmt.Println("------------------", filename, "--------------------------------------")
	// fmt.Println(html)
	// fmt.Println("``````````````````````````````````````````````````````````````````````````````````````")

	// Old version uses filename as key, I make my own key. not
	// filepath.Base(filename) First template becomes return value if
	// not already defined, and we use that one for subsequent New
	// calls to associate all the templates together. Also, if this
	// file has the same name as t, this file becomes the contents of
	// t, so t, err := New(name).Funcs(xxx).ParseFiles(name)
	// works. Otherwise we create a new template associated with t.
	name := key

	// Add to template
	t := Templates
	var tmpl *template.Template
	if t == nil {
		t = template.New(name)
	}
	if name == t.Name() {
		tmpl = t
	} else {
		tmpl = t.New(name)
	}

	_, err = tmpl.Parse(html)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func AddGOTTemplate__go_tempaltes_parser(key string, filename string) (*template.Template, error) {

	debug.Log("-   - [ParseTempalte] %v, %v", key, filename)

	// borrowed from html/tempate
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, nil // file not exist, don't panic.
	}

	// read source from file.
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	html := string(b)
	name := key
	// Old version uses filename as key, I make my own key. not
	// filepath.Base(filename) First template becomes return value if
	// not already defined, and we use that one for subsequent New
	// calls to associate all the templates together. Also, if this
	// file has the same name as t, this file becomes the contents of
	// t, so t, err := New(name).Funcs(xxx).ParseFiles(name)
	// works. Otherwise we create a new template associated with t.

	// Add to template
	t := Templates
	var tmpl *template.Template
	if t == nil {
		t = template.New(name)
	}
	if name == t.Name() {
		tmpl = t
	} else {
		tmpl = t.New(name)
	}

	// old
	_, err = tmpl.Parse(html)
	if err != nil {
		return nil, err
	}
	// fmt.Println("--------------------------------------------------------------------------------------")
	// fmt.Println("--------------------------------------------------------------------------------------")
	// fmt.Println(ttt)
	// fmt.Println(debug.PrintEntrails(ttt))
	// fmt.Println("``````````````````````````````````````````````````````````````````````````````````````")
	return tmpl, nil
}

/* ________________________________________________________________________________
   parse template
   TODO move to new package
*/

/*_______________________________________________________________________________
  Render Tempaltes
*/
func RenderGotTemplate(w io.Writer, tplkey string, p interface{}) error {
	err := Templates.ExecuteTemplate(w, tplkey, p)
	if err != nil {
		// panic(err)
		return err
	}
	return nil
}

/*_______________________________________________________________________________
  GOT Templates Caches
*/
// TemplateCache value, mark if template is parsed.
var GotTemplateCache TemplateCache = TemplateCache{Templates: map[string]bool{}}

type TemplateCache struct {
	l         sync.Mutex      // TODO lock this.
	Templates map[string]bool // is this template cached?
}

// TODO Performance issue
// . the first return value is not used.
// . the name shoud change
func (t *TemplateCache) Get(key string, templatePath string) (*template.Template, error) {
	t.l.Lock()
	_, ok := t.Templates[templatePath]
	t.l.Unlock()
	if !ok {
		tmpl, err := AddGOTTemplate(key, templatePath)
		if err != nil {
			// panic(err.Error())
			return nil, err
		}
		if tmpl == nil {
			err = errors.New(fmt.Sprintf("Templates for '%v' not found!", key))
			return nil, err
		}

		t.l.Lock()
		t.Templates[templatePath] = true
		t.l.Unlock()

		return tmpl, nil
	}
	return nil, nil
}

// // used by lifecircle-component
// func RenderTemplate(w io.Writer, tmpl string, p interface{}) error {
// 	err := Templates.ExecuteTemplate(w, tmpl+".html", p)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return nil
// }
