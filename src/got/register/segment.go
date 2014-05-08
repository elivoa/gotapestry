package register

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/elivoa/got/config"
	"github.com/elivoa/got/parser"
	"got/core"
	"log"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
)

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
	Parent   *ProtonSegment            //
	Children map[string]*ProtonSegment //
	Level    int                       // depth

	// TODO: Test Performance: New Method
	//   - Test Perforance between `reflect new` and `native func call`
	//   ? Use Generated New function (e.g. NewSomePage) to create new Page? Is This Faster?
	Proton core.Protoner // The base proton segment. Create new one when installed.

	// associated external resources.
	ModulePackage string             // e.g. got/builtin, syd; used in init.
	StructInfo    *parser.StructInfo // from parser package
	module        *core.Module       // associated Module

	// caches
	identity     string // cache identity, default the same name with StructName
	templatePath string // cache template path.

	// TODO - try the method that use use channel to lock.
	// TODO - Use RWMutex lock
	l sync.RWMutex

	// TODO replace with typeinfo
	// Path          string // ? TODO: URL path; TODO use appconfig
	// Src           string // source package, used to select app
}

func (s *ProtonSegment) AddChild(seg *ProtonSegment) {
	if s.Children == nil {
		s.Children = map[string]*ProtonSegment{}
	}
	s.Children[strings.ToLower(seg.Name)] = seg
}

func (s *ProtonSegment) HasChild(seg string) bool {
	return s.Children != nil && s.Children[strings.ToLower(seg)] != nil
}

// used to update register
func (s *ProtonSegment) Remove() {
	// TODO implement this. used in auto reload.
}

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

func (s *ProtonSegment) Module() *core.Module {
	if s.module == nil {
		if s.StructInfo != nil {
			// for k, module := range Modules.Map() {
			// 	fmt.Println("--- ", k, module.String())
			// }
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
		currentSeg = s
		prevSeg    = "//nothing//" // previous lowercase seg
		prevSegs   = []string{}    // previous lowercase seg[]

		selectorPrefix = []string{} // tempvalue
	)

	// 1. process path segments, without last node
	for idx, seg := range segments[0:(len(segments) - 1)] {
		var lowerSeg = strings.ToLower(seg)
		var segment *ProtonSegment

		if currentSeg.HasChild(seg) {
			segment = currentSeg.Children[seg]
			// TODO detect conflict
		} else {
			segment = &ProtonSegment{
				Name:   seg,
				Parent: currentSeg,
				Level:  idx,
			}
			dlog("!!!! add path to structure: seg: %v\n", seg) // ------------------
			currentSeg.AddChild(segment)
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

	// match origin paths: /order/create/OrderCreateIndex
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
	for _, s := range finalSegs {
		// link segment together.
		segment := &ProtonSegment{
			Name:       s,
			Parent:     currentSeg,
			Level:      len(segments) - 1,
			Proton:     p,
			StructInfo: si,
		}
		currentSeg.AddChild(segment)

		// add the first segment into typemap
		if segment.Proton.Kind() == core.PAGE {
			PageTypeMap[reflect.TypeOf(p).Elem()] = segment
		} else if segment.Proton.Kind() == core.COMPONENT {
			ComponentTypeMap[reflect.TypeOf(p).Elem()] = segment
		} else if segment.Proton.Kind() == core.MIXIN {
			MixinTypeMap[reflect.TypeOf(p).Elem()] = segment
		}

		// add selector
		selector := []string{}
		selector = append(selector, selectorPrefix...)
		selector = append(selector, s)
		selectors = append(selectors, selector)
	}
	// dlog(">>>>> Selectors: %v\n", selectors)
	return
}

// ----  Lookup & Results  ------------------------------------------------------------------------------

type LookupResult struct {
	Segment   *ProtonSegment
	PageUrl   string // value of request.URL.Path
	EventName string
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

var average_lookup_time int

// Lookup the structure, find the right page or component.
// Can detect event calls, event calls on embed components.
// TODO performance
func (s *ProtonSegment) Lookup(url string) (result *LookupResult, err error) {
	logLookup("- - - [Lookup] '%v'\n", url)

	// 1. pre-process url
	trimedUrl := strings.Trim(url, " ")
	if !strings.HasSuffix(trimedUrl, "/") {
		trimedUrl += "/"
	}
	segments := strings.Split(trimedUrl, "/")

	var (
		level int
		seg   string
		event string
	)
	// BUG: param segment not used?
	segment := s // loop channel object
	for level, seg = range segments {
		// skip the first / segment.
		if level == 0 && seg == "" {
			continue
		}
		logLookup("- - - [Lookup] Step: Level %v Seg:[ %-10v ] segment:[ %-20v ]\n",
			level, seg, segment)

		// If contains ".", this is an event call.
		// and match stops here, others are parameters of event.
		index := strings.Index(seg, ".")
		if index > 0 {
			event = seg[index+1:]
			seg = strings.ToLower(seg[0:index])
			segment = segment.Children[seg]
			level = level + 1
			break
		}

		// NEXT LEVEL LOOP
		if segment.Children == nil || len(segment.Children) == 0 || !segment.HasChild(seg) {
			logLookup("- - - [Lookup] match finished.")
			break
		} else {
			segment = segment.Children[strings.ToLower(seg)]
		}
	}

	// get page url
	pageUrl := strings.Join(segments[:level], "/")
	if event != "" {
		index := strings.LastIndex(pageUrl, ".")
		pageUrl = pageUrl[:index]
	}
	// log.Printf("- - - [Lookup] 'pageurl is' %v  (including event)\n", pageUrl)

	if nil == segment {
		err = errors.New("Lookup Failed.")
	}
	result = &LookupResult{
		Segment:   segment,
		PageUrl:   pageUrl,
		EventName: event,
	}
	logLookup("- - - [Lookup] Result is %v", result)
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
