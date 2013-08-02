package register

import (
	"errors"
	"fmt"
	"got/config"
	"got/core"
	"log"
	"path/filepath"
	"strings"
	"sync"
)

var conf = config.Config

/* ________________________________________________________________________________
   Segment, like a trie.
*/
type ProtonSegment struct {
	Name     string                    // segment name
	Path     string                    // TODO: URL path; TODO use appconfig
	Parent   *ProtonSegment            //
	Children map[string]*ProtonSegment //
	Level    int                       // depth
	Proton   core.Protoner             // Proton
	Src      string                    // source package, used to select app

	identity     string // cache identity
	templatePath string // cache template path.

	// TODO cache structure info.
	// TODO add pointer to proton?

	l sync.Mutex // TODO used to synchronized, use channel version to support multiwrite.
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
	core.PAGE:      "",
	core.COMPONENT: "c_",
	core.MIXIN:     "x_",
}

func (s *ProtonSegment) Identity() string {
	if s.identity == "" {
		appConfig := Apps.Get(s.Src)
		if appConfig == nil {
			panic(fmt.Sprintf("Can't find APP Config %v", s.Src))
		}
		s.identity = fmt.Sprintf("%v%v:%v", identityPrefixMap[s.Proton.Kind()], s.Src, s.Path)
	}
	return s.identity
}

func (s *ProtonSegment) TemplatePath() (string, string) {
	if s.identity == "" {
		appConfig := Apps.Get(s.Src)
		if appConfig == nil {
			panic(fmt.Sprintf("Can't find APP Config %v", s.Src))
		}
		s.identity = fmt.Sprintf("%v%v:%v", identityPrefixMap[s.Proton.Kind()], s.Src, s.Path)

		if s.templatePath == "" {
			s.templatePath = filepath.Join(
				appConfig.FilePath,
				pathMap[s.Proton.Kind()], // "pages","components"
				s.Path,
			) + conf.TemplateFileExtension // TODO Configthis
		}
	}
	return s.identity, s.templatePath
}

// ________________________________________________________________________________
// parse url and put it into segments.
// return what?
// TODO return [][]string:
//   order, list
//   order, orderlist
//   order/create/OrderCreateDetail
//
func (s *ProtonSegment) Add(baseUrl string, p core.Protoner) (selectors [][]string) {

	src, segments := trimPathSegments(baseUrl, pathMap[p.Kind()])

	dlog("-___- [Register %v] %v::%v url:%v", pathMap[p.Kind()], src, segments, baseUrl)

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
			dlog("!!!! cached path: currentSeg: %v, has seg: %v\n", currentSeg.Name, seg)
			dlog("!!!! children: %v\n", s.Children[seg])
			// detect conflict
			if segment.Src != "" && segment.Src != src {
				log.Fatalf("Conflict of Page defination %v.\n", baseUrl)
			}
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
	// log.Printf("-www- [RegisterPage] enter last node %v; \n", seg)
	// log.Printf("-www- --------- %v; %v \n", lowerSeg, prevSeg)

	// ~ 2 ~ overlapped keywords: i.e.: order/OrderEdit ==> order/edit
	var (
		seg           string   = segments[len(segments)-1]
		lowerSeg      string   = strings.ToLower(seg)
		shortLowerSeg string   = lowerSeg
		shortSeg      string   = seg           // shorted seg with case
		finalSegs     []string = []string{seg} // alias
	)

	// match origin paths: /order/create/OrderCreateIndex
	// fmt.Println(prevSegs)
	for _, p := range prevSegs {
		// dlog("+++ strings.HasPrefix: %v, %v = %v\n", shortLowerSeg, p, strings.HasPrefix(shortLowerSeg, p))
		if strings.HasPrefix(shortLowerSeg, p) {
			shortSeg = shortSeg[len(p):]
			shortLowerSeg = shortLowerSeg[len(p):]
			if shortSeg != "" {
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
			currentSeg.Src = src
			currentSeg.Path = strings.Join(segments, "/")
			currentSeg.Proton = p
		}
	}

	// judge empty/index

	// /order/Order[Index] - fall back to /order/
	if p.Kind() == core.PAGE {
		if strings.HasSuffix(shortLowerSeg, "index") {
			dlog("+++ Match Index, \n") // ------------------------------------------

			var trimlen = len(shortLowerSeg) - len("index")
			if trimlen >= 0 {
				shortSeg = shortSeg[:trimlen]
				shortLowerSeg = shortLowerSeg[:trimlen]
			}
			if shortSeg == "" {
				// fallback
				dlog("+++++ Fallback.\n") // ------------------------------------------
				currentSeg.Src = src
				currentSeg.Path = strings.Join(segments, "/")
				currentSeg.Proton = p
			} else {
				// eg: /order/OrderDetailIndex --> /order/detail
				finalSegs = append(finalSegs, shortSeg)
			}
		}
	}

	dlog(">>>>> FinalSegs: %v\n", finalSegs) // ------------------------------------------

	// 4. finally add segment struct to chains.
	for _, s := range finalSegs {
		// link segment together.
		segment := &ProtonSegment{
			Name:   s,
			Parent: currentSeg,
			Level:  len(segments) - 1,
			Src:    src,
			Path:   strings.Join(segments, "/"),
			Proton: p,
		}
		currentSeg.AddChild(segment)

		// add selector
		selector := []string{}
		selector = append(selector, selectorPrefix...)
		selector = append(selector, s)
		selectors = append(selectors, selector)
	}
	dlog(">>>>> Selectors: %v\n", selectors)
	return
}

// ----  Lookup & Results  ------------------------------------------------------------------------------

type LookupResult struct {
	Segment   *ProtonSegment
	PageUrl   string
	EventName string
}

// Lookup the structure, find the right page/component.
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

// string()
func (s *ProtonSegment) String() string {
	return fmt.Sprintf("%-14v (%v)[SRC='%v' PATH='%v']",
		s.Name, len(s.Children), s.Src, s.Path)
}

// print all details
func (s *ProtonSegment) PrintALL() string {
	s.print(s)
	return ""
}

func (s *ProtonSegment) print(segment *ProtonSegment) string {
	fmt.Printf("+ %v\n", segment)
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
