/*
   Time-stamp: <[templates.go] Elivoa @ Sunday, 2014-05-11 17:39:08>
*/
package templates

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/elivoa/got/config"
	"github.com/elivoa/got/templates/transform"
	"got/debug"
	"got/register"
	"html/template"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
)

// Templates stores all templates.
var Templates *template.Template

func init() {
	// init template TODO remove this, change another init method.
	// TODO: use better way to init.
	Templates = template.New("-")

	// Register built-in templates.
	registerBuiltinFuncs()
}

/*_______________________________________________________________________________
  Register components
*/

// register components as template function call.
func RegisterComponentAsFunc(name string, f interface{}) {
	funcName := fmt.Sprintf("t_%v", strings.Replace(name, "/", "_", -1))
	lowerFuncName := strings.ToLower(funcName)
	Templates.Funcs(template.FuncMap{funcName: f, lowerFuncName: f})
}

/*_______________________________________________________________________________
  Render Tempaltes
*/

// RenderTemplate render template into writer.
func RenderTemplate(w io.Writer, key string, p interface{}) error {
	err := Templates.ExecuteTemplate(w, key, p)
	if err != nil {
		return err
	}
	return nil
}

/*_______________________________________________________________________________
  GOT Templates Caches
*/

// TODO: Integrate to register.

// init template cache
var Cache TemplateCache = TemplateCache{
	Templates: map[reflect.Type]*TemplateUnit{},
	Keymap:    map[string]*TemplateUnit{},
}

// TemplateCache cache templates
type TemplateCache struct {
	l sync.RWMutex

	Templates map[reflect.Type]*TemplateUnit // type as key
	Keymap    map[string]*TemplateUnit       // template key as key
}

type TemplateUnit struct {
	Key               string // template key.
	FilePath          string
	ContentOrigin     string `json:"-"`
	ContentTransfered string `json:"-"`
	IsCached          bool   `json:"-"`

	Blocks map[string]*BlockUnit

	// todo components?
}

// Note: Component in blocks are directly belong to block's container.
type BlockUnit struct {
	Container         *TemplateUnit
	ID                string // block's id
	ContentOrigin     string
	ContentTransfered string
}

// Get cached TemplateUnit by proton type.
func (t *TemplateCache) Get(protonType reflect.Type) (*TemplateUnit, error) {
	t.l.RLock()
	defer t.l.RUnlock()
	if unit, ok := t.Templates[protonType]; ok {
		return unit, nil
	}
	return nil, errors.New("Template not loaded.")
}

// Get cached TemplateUnit by template key.
func (t *TemplateCache) GetByKey(key string) (*TemplateUnit, error) {
	t.l.RLock()
	defer t.l.RUnlock()
	if unit, ok := t.Keymap[key]; ok {
		return unit, nil
	}
	return nil, errors.New("Template not loaded.")
}

func (t *TemplateCache) GetBlock(protonType reflect.Type, blockName string) (*BlockUnit, error) {
	t.l.RLock()
	defer t.l.RUnlock()
	if unit, ok := t.Templates[protonType]; ok {
		if nil == unit {
			return nil, errors.New("Error: Templates are nil, can't has blocks.")
		}
		if bu, okb := unit.Blocks[blockName]; okb {
			return bu, nil
		}
		return nil, errors.New(fmt.Sprintf("Block '%v' not found.", blockName))
	}
	return nil, errors.New("Template not loaded.")
}

func (t *TemplateCache) GetBlockByKey(key string, blockName string) (*BlockUnit, error) {
	t.l.RLock()
	defer t.l.RUnlock()
	if unit, ok := t.Keymap[key]; ok {
		if nil == unit {
			return nil, errors.New("Error: Templates are nil, can't has blocks.")
		}
		if bu, okb := unit.Blocks[blockName]; okb {
			return bu, nil
		}
		return nil, errors.New(fmt.Sprintf("Block '%v' not found.", blockName))
	}
	return nil, errors.New("Template not loaded.")
}

// Parse and cache page or component's template, return the cached one.
func (t *TemplateCache) GetnParse(key string, templatePath string, protonType reflect.Type) (*TemplateUnit, error) {
	/* TODO 这里模板上锁的机制有问题。
	   1. 先上锁判断是否存在，然后初始化，设置的时候上第二道锁；
	      缺点：并发多的时候，会有多个进程同时初始化。
	   2. 解决方案, 用rw嗦，读取写入的时候上多到嗦。
	*/
	fmt.Println(register.Components)
	// if !ok {
	forceLoad := false
	tu, cached, err := t.LoadTemplates(protonType, key, templatePath, forceLoad)
	if err != nil { // error occured
		return nil, err
	}
	if tu == nil { // no error and no result.
		err = errors.New(fmt.Sprintf("Templates for '%v' not found!", key))
		return nil, err
	}
	if !cached { // return the cached one.
		// parse templates
		t.l.Lock() // write lock
		ParseTemplate(tu)
		t.l.Unlock()
	}
	return tu, nil
}

/*
  Load template and it's contents into memory.
  TODO: zip the source
  TODO: implement force reload
*/
func (t *TemplateCache) LoadTemplates(protonType reflect.Type, key string, filename string, forceReload bool) (tu *TemplateUnit, cached bool, err error) {
	debug.Log("-   - [ParseTemplate] %v", filename)

	// TODO: 这里的锁有问题，高并发时容易引起资源浪费。
	if !forceReload { // read cache.
		if tu, err = t.Get(protonType); err == nil {
			// Be Lazy, err is Tempalte not loaded yet!
			cached = true
			return // return cached version.
		}
	}

	// if file doesn't exist.
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		// set nil to cache
		t.Templates[protonType] = nil
		return
	} else if err != nil {
		return // other file error.
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

	tu = &TemplateUnit{
		Key:               key,
		FilePath:          filename,
		IsCached:          true,
		ContentTransfered: trans.RenderToString(),
	}

	blocks := trans.RenderBlocks() // blocks found in template.
	if blocks != nil {
		tu.Blocks = map[string]*BlockUnit{}
		for blockId, html := range blocks {
			tu.Blocks[blockId] = &BlockUnit{
				Container:         tu,
				ID:                blockId,
				ContentTransfered: html,
			}
		}
	}

	// add to cache
	t.Templates[protonType] = tu
	t.Keymap[tu.Key] = tu
	return
}

func ParseTemplate(tu *TemplateUnit) (err error) {
	// parse main
	if err = parseTemplate(tu.Key, tu.ContentTransfered); err != nil {
		return
	}

	if tu.Blocks != nil {
		for _, block := range tu.Blocks {
			key := fmt.Sprintf("%s%s%s", tu.Key, config.SPLITER_BLOCK, block.ID)
			if err = parseTemplate(key, block.ContentTransfered); err != nil {
				return
			}
		}
	}
	return nil
}

func parseTemplate(key string, content string) error {
	// Old version uses filename as key, I make my own key. not
	// filepath.Base(filename) First template becomes return value if
	// not already defined, and we use that one for subsequent New
	// calls to associate all the templates together. Also, if this
	// file has the same name as t, this file becomes the contents of
	// t, so t, err := New(name).Funcs(xxx).ParseFiles(name)
	// works. Otherwise we create a new template associated with t.

	t := Templates
	var tmpl *template.Template
	if t == nil {
		t = template.New(key)
	}
	if key == t.Name() {
		tmpl = t
	} else {
		tmpl = t.New(key)
	}

	_, err := tmpl.Parse(content)
	if err != nil {
		return err
	}
	return nil
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
