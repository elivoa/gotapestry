package templates

import (
	"errors"
	"fmt"
	"got/debug"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"sync"
)

var Templates *template.Template

func init() {
	// init template
	Templates = template.Must(template.ParseFiles(
		LocatePath("nothing"),
	))

	registerHelperFuncs()
	registerComponentFuncs() // TEST

	// init functions
	// Templates.Funcs(template.FuncMap{
	// 	"formattime":     FormatTime,
	// 	"beautytime":     BeautyTime,
	// 	"formatcurrency": FormatCurrency,
	// })
}

// old
func Add(templateName string) {
	tml, err := Templates.ParseFiles(LocatePath(templateName))
	template.Must(tml, err)
	if err != nil {
		panic(err.Error())
	}
}

/*
 return
   template,  when tempalte available and parse successful.
   nil,  when template not exists or error occurs.
*/
func AddGOTTemplate(key string, filename string) (*template.Template, error) {

	debug.Log("-   - [ParseTempalte] %v, %v", key, filename)

	// copied from html/tempate

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// file not exists.
		return nil, nil
	}

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	s := string(b)
	name := key // old is use filename as key, I make my own key. not filepath.Base(filename)
	// First template becomes return value if not already defined,
	// and we use that one for subsequent New calls to associate
	// all the templates together. Also, if this file has the same name
	// as t, this file becomes the contents of t, so
	//  t, err := New(name).Funcs(xxx).ParseFiles(name)
	// works. Otherwise we create a new template associated with t.

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
	_, err = tmpl.Parse(s)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

/*
   Utils
*/
func LocatePath(name string) (templates string) {
	return fmt.Sprintf("../src/template/%v.html", name)
}

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
// TemplateCache value
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

/*_______________________________________________________________________________
  Old things
*/

// very old
func RenderTemplate(w io.Writer, tmpl string, p interface{}) error {
	err := Templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		panic(err)
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}
