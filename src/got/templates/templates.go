/*
   Time-stamp: <[templates.go] Elivoa @ Tuesday, 2013-08-20 19:17:28>
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
	"log"
	"os"
	"strings"
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
	Templates = template.New("-")
	registerBuiltinFuncs()   // Register built-in templates.
	registerComponentFuncs() // no use?
}

/*_______________________________________________________________________________
  Register components
*/
func registerComponentFuncs() {
	// init functions
	Templates.Funcs(template.FuncMap{})
}

// register components as template function call.
func RegisterComponent(name string, f interface{}) {
	funcName := fmt.Sprintf("t_%v", strings.Replace(name, "/", "_", -1))
	lowerFuncName := strings.ToLower(funcName)
	// debuglog("-108- [RegisterComponent] %v", funcName)
	Templates.Funcs(template.FuncMap{funcName: f, lowerFuncName: f})
}

/*
 @return
   template - when tempalte available and parse successful.
   nil      - when template not exists or error occurs.
*/
func parseTempaltes(key string, filename string) (map[string]*template.Template, error) {

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
	trans.Parse(r)

	templatesToParse := map[string]string{}
	templatesToParse[key] = trans.RenderToString()
	blocks := trans.RenderBlocks()
	if blocks != nil {
		for blockId, html := range blocks {
			templatesToParse[fmt.Sprintf("%v:%v", key, blockId)] = html
		}
	}

	fmt.Println("--------------------------------------------------------------------------------------")
	fmt.Println("------------------", filename, "--------------------------------------")
	fmt.Println(templatesToParse[key])
	fmt.Println("``````````````````````````````````````````````````````````````````````````````````````")

	// Old version uses filename as key, I make my own key. not
	// filepath.Base(filename) First template becomes return value if
	// not already defined, and we use that one for subsequent New
	// calls to associate all the templates together. Also, if this
	// file has the same name as t, this file becomes the contents of
	// t, so t, err := New(name).Funcs(xxx).ParseFiles(name)
	// works. Otherwise we create a new template associated with t.
	returns := map[string]*template.Template{}
	t := Templates
	for name, html := range templatesToParse {
		// Add to template
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
		returns[name] = tmpl
	}
	return returns, nil
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

// TODO Test & Improve Performance.
// . the first return value is not used.
// . the name shoud change
func (t *TemplateCache) Get(key string, templatePath string) (*template.Template, error) {
	t.l.Lock()
	_, ok := t.Templates[templatePath]
	t.l.Unlock()
	if !ok {
		tmpls, err := parseTempaltes(key, templatePath) // map[string]*template.Template, error
		if err != nil {
			return nil, err
		}
		if tmpls == nil {
			err = errors.New(fmt.Sprintf("Templates for '%v' not found!", key))
			return nil, err
		}

		t.l.Lock()
		t.Templates[templatePath] = true
		t.l.Unlock()

		return tmpls[key], nil
	}
	return nil, nil
}

// --------------------------------------------------------------------------------
// log
//
var debugLog = true

func debuglog(format string, params ...interface{}) {
	if debugLog {
		log.Printf(format, params...)
	}
}
