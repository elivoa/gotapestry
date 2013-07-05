package register

import (
	"errors"
	"fmt"
	"got/core"
	"log"
	"strings"
	"sync"
)

/* ________________________________________________________________________________
   Segment, like a trie.
*/
type ProtonSegment struct {
	Name     string                    // segment name
	Path     string                    // TODO: URL path
	Parent   *ProtonSegment            //
	Children map[string]*ProtonSegment //
	Src      string                    // source package, used to select app
	Level    int                       // depth
	Proton   core.IProton              // Proton

	l sync.Mutex // TODO used to synchronized, also in page?
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
	// TODO
}

var pathMap map[string]string = map[string]string{
	"page": "pages", "component": "components",
}

// ________________________________________________________________________________
// parse url and put it into segments.
// return what?
// TODO return [][]string:
//   order, list
//   order, orderlist
//   order/create/OrderCreateDetail
//
func (s *ProtonSegment) Add(baseUrl string, p core.IProton, protonType string) (selectors [][]string) {

	src, segments := trimPathSegments(baseUrl, pathMap[protonType])

	dlog("-___- [Register %v] %v::%v url:%v", protonType, src, segments, baseUrl)

	// add to registerc
	var (
		currentSeg = s
		prevSeg    = "//nothing//" // previous lowercase seg
		prevSegs   = []string{}    // previous lowercase seg
		isPage     = (protonType == "page")

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
	// > the last node.

	// log.Printf("-www- [RegisterPage] enter last node %v; \n", seg)
	// log.Printf("-www- --------- %v; %v \n", lowerSeg, prevSeg)

	// ~ 2 ~ overlap keywords: i.e.: order/OrderEdit ==> order/edit
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
	if isPage {
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

	// 4. finally add segment struct to chains.
	dlog(">>>>> FinalSegs: %v\n", finalSegs) // ------------------------------------------
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
		logLookup("- - - [Lookup] Step: Level %v Seg:[ %-10v ] segment:[ %-20v ]\n",
			level, seg, segment)

		// skip the first / segment.
		if level == 0 && seg == "" {
			continue
		}

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
			segment = segment.Children[seg]
		}
	}

	{
		// fmt.Printf("+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+\n")
		// fmt.Printf("+ segments: \n")
		// for idx, s := range segments {
		// 	fmt.Printf("%v:%v\n", idx, s)
		// }
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
