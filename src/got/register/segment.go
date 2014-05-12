package register

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/elivoa/got/config"
	"github.com/elivoa/got/logs"
	"github.com/elivoa/got/parser"
	"got/core"
	"log"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
)

// ----  Identity & Template path  ---------------------------------------------------------------

var pathMap = map[core.Kind]string{
	core.PAGE:      "pages",
	core.COMPONENT: "components",
	core.MIXIN:     "mixins",
}

var identityPrefixMap = map[core.Kind]string{
	core.PAGE:      "p/",
	core.COMPONENT: "c/",
	core.MIXIN:     "x/",
}

var conf = config.Config

// ProtonSegment is a tree like structure to hold path to Page/Component
// 1. Support quick lookup to locate a page or component. (TODO need improve performance)
// 2. Each kind of page has one ProtonSegment instance. (one path)
// TODO
//   - refactor this.
//
type ProtonSegment struct {
	// as a tree node
	Name     string                    // segment name
	Alias    []string                  // alias, e.g.(order/OrderEdit): edit, orderedit
	Parent   *ProtonSegment            //
	Children map[string]*ProtonSegment //
	Level    int                       // depth

	// template related
	IsTemplateLoaded  bool
	Blocks            map[string]*Block         // blocks
	ContentOrigin     string                    // template's html
	ContentTransfered string                    // template's transfered html
	EmbedComponents   map[string]*ProtonSegment // lowercased id

	// TODO: Test Performance: New Method
	//   - Test Perforance between `reflect new` and `native func call`
	//   ? Use Generated New function (e.g. NewSomePage) to create new Page? Is This Faster?
	// TODO: Chagne name
	Proton core.Protoner // The base proton segment. Create new one when installed.

	// associated external resources.
	ModulePackage string             // e.g. got/builtin, syd; used in init.
	StructInfo    *parser.StructInfo // from parser package
	module        *core.Module       // associated Module

	// caches
	identity     string // cache identity, default the same name with StructName
	templatePath string // cache template path.

	// TODO - try the method that use use channel to lock.
	L sync.RWMutex
}

type Block struct {
	ID                string // block's id
	ContentOrigin     string
	ContentTransfered string
}

func (s *ProtonSegment) AddChild(segname string, seg *ProtonSegment) {
	if s.Children == nil {
		s.Children = map[string]*ProtonSegment{}
	}
	s.Children[strings.ToLower(segname)] = seg
}

func (s *ProtonSegment) HasChild(seg string) bool {
	return s.Children != nil && s.Children[strings.ToLower(seg)] != nil
}

// used to update register
func (s *ProtonSegment) Remove() {
	panic("not implement!")
	// TODO implement this. used in auto reload.
}

// unique identity used as template key.
// TODO refactor all Identities of proton. with event call and event path call.
func (s *ProtonSegment) Identity() string {
	if s.identity == "" {
		s.identity = fmt.Sprintf("%v:%v",
			path.Join(identityPrefixMap[s.Proton.Kind()], s.StructInfo.ProtonPath()),
			s.StructInfo.StructName)
	}
	return s.identity
}

// TemplatePath returns the tempalte key and it's full path.
func (s *ProtonSegment) TemplatePath() (string, string) {
	if s.templatePath == "" {
		module := s.Module()
		if s.templatePath == "" {
			s.templatePath = filepath.Join(
				module.BasePath,
				s.StructInfo.ImportPath,
				fmt.Sprintf("%v%v", s.StructInfo.StructName, conf.TemplateFileExtension),
			) // TODO Configthis
		}
	}
	return s.Identity(), s.templatePath
}

// Find it's Module
func (s *ProtonSegment) Module() *core.Module {
	if s.module == nil {
		if s.StructInfo != nil {
			module := Modules.Get(s.StructInfo.ModulePackage)
			if module == nil {
				panic(fmt.Sprint("Can't find module for ", s.StructInfo.ModulePackage))
			}
			s.module = module
		}
	}
	return s.module
}

// ________________________________________________________________________________
// parse url and put it into segments.
// return what?
// TODO return [][]string:
//   order, list
//   order, orderlist
//   order/create/OrderCreateDetail
//
func (s *ProtonSegment) Add(si *parser.StructInfo, p core.Protoner) (selectors [][]string) {

	// TODO segment has structinfo
	src := si.ModulePackage
	segments := strings.Split(si.ProtonPath(), "/")
	if len(segments) > 0 && segments[0] == "" {
		segments = segments[1:]
	}
	segments = append(segments, si.StructName)

	dlog("-___- [Register %v] %v::%v url:%v", pathMap[p.Kind()], src, segments, si)

	// add to registerc
	var (
		currentSeg     = s             // always use root segment.
		prevSeg        = "//nothing//" // previous lowercase seg
		prevSegs       = []string{}    // previous lowercase seg[]
		selectorPrefix = []string{}    // tempvalue
	)

	// 1. process path segments to reach the end, without last node.
	for idx, seg := range segments[0:(len(segments) - 1)] {
		var lowerSeg = strings.ToLower(seg)

		var segment *ProtonSegment
		if currentSeg.HasChild(seg) {
			segment = currentSeg.Children[seg]
			// TODO detect conflict
		} else {
			// Add path to segment.
			segment = &ProtonSegment{
				Name:   seg,
				Parent: currentSeg,
				Level:  idx,
			}
			dlog("!!!! add path to structure: seg: %v\n", seg) // ------------------
			currentSeg.AddChild(segment.Name, segment)
		}
		currentSeg = segment
		selectorPrefix = append(selectorPrefix, seg)
		prevSeg = lowerSeg
		prevSegs = append(prevSegs, lowerSeg)
	}

	// 2. process last node
	// ~ 2 ~ overlapped keywords: i.e.: order/OrderEdit ==> order/edit
	var (
		seg           string   = segments[len(segments)-1]
		lowerSeg      string   = strings.ToLower(seg)
		shortLowerSeg string   = lowerSeg
		shortSeg      string   = seg           // shorted seg with case
		finalSegs     []string = []string{seg} // alias
	)

	// fmt.Printf("-www- [RegisterPage] enter last node %v; \n", seg)
	// fmt.Printf("-www- --------- %v; %v \n", lowerSeg, prevSeg)

	// Match origin paths: /order/create/OrderCreateIndex, in this example we can ignore
	// prefix 'Order' and the ignore 'Create', and 'Index' can be automatically ignored.
	for _, p := range prevSegs {
		// dlog("+++ strings.HasPrefix: %v, %v = %v\n", shortLowerSeg, p, strings.HasPrefix(shortLowerSeg, p))
		if strings.HasPrefix(shortLowerSeg, p) {
			shortSeg = shortSeg[len(p):]
			shortLowerSeg = shortLowerSeg[len(p):]
			if shortSeg != "" { // e.g. order/create/OrderCreateIndex
				finalSegs = append(finalSegs, shortSeg)
				dlog("!!! p:%v, add to final: %v\n", p, shortSeg)
			}
		}
	}

	// TODO remove suffix in the same way.
	// TODO kill this.
	// /order/Create[Order] - ignore [Order]
	if strings.HasSuffix(lowerSeg, prevSeg) {
		dlog("+++ Match suffix, \n") // ------------------------------------------
		s := seg[:len(seg)-len(prevSeg)]
		if s != "" {
			finalSegs = append(finalSegs, s)
		} else {
			// fallback TODO
			// currentSeg.Src = src
			currentSeg.Proton = p
		}
	}

	// judge empty/index

	// remove index if any
	// /order/Order[Index] - fall back to /order/
	// /api/suggest/Suggest - fall back to /api/suggest
	if p.Kind() == core.PAGE && strings.HasSuffix(shortLowerSeg, "index") {
		dlog("+++ Match Index, \n") // ------------------------------------------

		var trimlen = len(shortLowerSeg) - len("index")
		if trimlen >= 0 {
			shortSeg = shortSeg[:trimlen]
			shortLowerSeg = shortLowerSeg[:trimlen]
		}
	}

	// Fallback if needed.
	if shortSeg == "" {
		// e.g.: /api/suggest/Suggest --> /api/suggest
		dlog("+++++ Fallback.\n") // ------------------------------------------
		currentSeg.Proton = p
		currentSeg.StructInfo = si
	} else {
		// e.g.: /order/OrderDetailIndex --> /order/detail
		finalSegs = append(finalSegs, shortSeg)
	}

	dlog(">>>>> FinalSegs: %v\n", finalSegs) // ------------------------------------------

	// 4. finally add segment struct to chains.
	segment := &ProtonSegment{
		Name:       finalSegs[0], // Name is a bitch.
		Alias:      finalSegs,
		Parent:     currentSeg,
		Level:      len(segments) - 1,
		Proton:     p,
		StructInfo: si,
	}

	for _, s := range finalSegs {
		currentSeg.AddChild(s, segment) // link segment together.

		// add selector
		selector := []string{}
		selector = append(selector, selectorPrefix...)
		selector = append(selector, s)
		selectors = append(selectors, selector)
	}

	// add the first segment into typemap
	switch segment.Proton.Kind() {
	case core.PAGE:
		PageTypeMap[reflect.TypeOf(p).Elem()] = segment
	case core.COMPONENT:
		ComponentTypeMap[reflect.TypeOf(p).Elem()] = segment
	case core.MIXIN:
		MixinTypeMap[reflect.TypeOf(p).Elem()] = segment
	}
	// dlog(">>>>> Selectors: %v\n", selectors)
	return
}

// ----  Lookup & Results  ------------------------------------------------------------------------------

type LookupResult struct {
	Segment        *ProtonSegment
	PageUrl        string   // value of request.URL.Path
	ComponentPaths []string // component path ids, for calling event.
	EventName      string
	Parameters     []string // parameters reconized.
}

func (lr *LookupResult) IsEventCall() bool {
	return lr.EventName != ""
}

func (lr *LookupResult) IsValid() bool {
	if lr.Segment != nil && lr.Segment.Proton != nil {
		return true
	}
	return false
}

func (lr *LookupResult) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf(">> [LookupResult]{\n"))
	buffer.WriteString(fmt.Sprintf("\tSegment:%v,\n", lr.Segment))
	buffer.WriteString(fmt.Sprintf("\tPageUrl:%v,\n", lr.PageUrl))
	buffer.WriteString(fmt.Sprintf("\tComponentPaths:%v,\n", lr.ComponentPaths))
	buffer.WriteString(fmt.Sprintf("\tEventName:%v,\n", lr.EventName))
	buffer.WriteString(fmt.Sprintf("\tParameters:%v,\n", lr.Parameters))
	buffer.WriteString("  }\n")
	return buffer.String()
}

var average_lookup_time int

var lookupLogger = logs.Get("URL Lookup")

// Lookup the structure, find the right page or component.
// Can detect event calls, event calls on embed components.
// TODO performance
// 例如： /got/Status.TemplateStatus:TemplateDetail/c__got:TemplateStatus
// 当遇到第一个.的时候，后面的为components. 当再遇到：的时候，后面的是方法名，/截断作为参数。
func (s *ProtonSegment) Lookup(url string) (result *LookupResult, err error) {
	// pre-process url
	trimedUrl := strings.Trim(url, " ")
	if !strings.HasSuffix(trimedUrl, "/") {
		trimedUrl += "/"
	}

	if lookupLogger.Debug() {
		lookupLogger.Printf("[Lookup] '%v'", url)
	}

	var (
		level         int = -1
		buffer        bytes.Buffer
		segments           = []string{}
		parameterPart bool = false
	)
	result = &LookupResult{
		ComponentPaths: []string{}, // init component paths.
		Parameters:     []string{},
	}

	segment := s // loop channel object
	for _, c := range trimedUrl {
		switch c {
		default:
			buffer.WriteRune(c)
			continue
		case '/':
			level += 1
		}

		// arrive here means words finished. process segment
		seg := buffer.String()
		segments = append(segments, seg)
		buffer.Reset()

		// skip the first / segment.
		if level == 0 && seg == "" {
			continue
		}

		if lookupLogger.Debug() {
			lookupLogger.Printf("[Lookup] Step: Level %v Seg:[ %-10v ] segment:[ %-20v ]\n",
				level, seg, segment)
		}

		// parameter mode
		if parameterPart {
			result.Parameters = append(result.Parameters, seg)
			continue
		}

		// parth lookup mode

		// If contains ":", this is an event call. Or parameter.
		// and match stops here, others are parameters of event.
		if index := strings.Index(seg, ":"); index > 0 {
			result.EventName = seg[index+1:]
			array := strings.Split(seg[0:index], ".")
			seg = strings.ToLower(array[0])
			result.ComponentPaths = array[1:]

			level = level + 1
			parameterPart = true
		}

		if segment.Children == nil || len(segment.Children) == 0 || !segment.HasChild(seg) {
			if lookupLogger.Debug() {
				lookupLogger.Printf("- - - [Lookup] match finished.")
			}
			// Match finished, this must be the first paramete
			result.Parameters = append(result.Parameters, seg)
			parameterPart = true

			break
		} else {
			fmt.Println("going into next step: ", segment)
			segment = segment.Children[strings.ToLower(seg)]
		}

	}

	// get page url
	// pageUrl := strings.Join(segments[:level], "/")
	// if result.EventName != "" {
	// 	// TODO: bugs here. can', .
	// 	index := strings.LastIndex(pageUrl, ".")
	// 	fmt.Println("................", index, " >>> ", pageUrl)
	// 	pageUrl = pageUrl[:index]
	// }
	// log.Printf("- - - [Lookup] 'pageurl is' %v  (including event)\n", pageUrl)

	if nil == segment {
		err = errors.New("Lookup Failed.")
	}
	result.Segment = segment
	// result.PageUrl = pageUrl
	if lookupLogger.Debug() {
		lookupLogger.Printf("- - - [Lookup] Result is %v", result)
	}
	return
}

/* ________________________________________________________________________________
   Print Helper
*/

func (s *ProtonSegment) String() string {
	length, path := ".", "--"
	if len(s.Children) > 0 {
		length = fmt.Sprint(len(s.Children))
	}
	if s.StructInfo != nil {
		path = s.StructInfo.ImportPath
	}
	return fmt.Sprintf("%-20v (%v)[%v]", s.Name, length, path)
	// return fmt.Sprintf("%-20v (%v)[%v]", s.Name, length, path)
}

// print all details
func (s *ProtonSegment) PrintALL() string {
	s.print(s)
	return ""
}

func (s *ProtonSegment) StringTree(newline string) string {
	var out bytes.Buffer // = bytes.NewBuffer([]byte{})
	s.treeSegment(&out, s, newline)
	return out.String()
}

func (s *ProtonSegment) treeSegment(out *bytes.Buffer, segment *ProtonSegment, newline string) {
	out.WriteString(fmt.Sprintf("+ %v >> %v%s", segment, segment.StructInfo, newline))
	for _, seg := range s.Children {
		for i := 0; i <= seg.Level; i++ {
			out.WriteString("  ")
		}
		if seg != nil {
			seg.treeSegment(out, seg, newline)
		}
	}
}

// TODO: user treeSegment instead.
func (s *ProtonSegment) print(segment *ProtonSegment) string {
	fmt.Printf("+ %v >> %v\n", segment, segment.StructInfo)
	for _, seg := range s.Children {
		for i := 0; i <= seg.Level; i++ {
			fmt.Print("  ")
		}
		if seg != nil {
			seg.print(seg)
		}
	}
	return ""
}

/* ________________________________________________________________________________
   TrimPathSegments
   Return: src, segments
   Param:
     protonType - [page|component]

   e.g. f("/got/builtin/pages/order/list") = "order/list/"

*/
func trimPathSegments(url string, protonType string) (string, []string) {
	segments := strings.Split(url, "/")
	var index, seg = 0, ""
	for index, seg = range segments {
		if seg == protonType {
			break
		}
	}
	var src string = ""
	if index > 0 {
		src = strings.Join(segments[0:index], "/")
	}
	return src, segments[index+1:]
}

// --------------------------------------------------------------------------------
// Tools & helper methods
// --------------------------------------------------------------------------------

var debug_add = false

func dlog(format string, params ...interface{}) {
	if debug_add {
		log.Printf(format, params...)
	}
}
