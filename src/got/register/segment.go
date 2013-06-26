package register

import (
	"errors"
	"fmt"
	"got/core"
	"got/debug"
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
//
func (s *ProtonSegment) Add(baseUrl string, p core.IProton, protonType string) (selectors [][]string) {
	src, segments := trimPathSegments(baseUrl, pathMap[protonType])
	debug.Log("-___- [Register %v] %v::%v url:%v", protonType, src, segments, baseUrl)

	// add to register
	var (
		currentSeg = s
		prevSeg    = "//nothing//"
		isPage     = (protonType == "page")

		selectorPrefix = []string{} // tempvalue
	)
	for idx, seg := range segments {
		lowerSeg := strings.ToLower(seg)
		var segment *ProtonSegment
		if currentSeg.HasChild(seg) {
			// detect conflict
			existSeg := s.Children[seg]
			if existSeg.Src != "" && existSeg.Src != src {
				log.Fatalf("Conflict of Page defination %v.\n", baseUrl)
			}
			currentSeg = existSeg
			selectorPrefix = append(selectorPrefix, seg)
		} else {
			// the last node

			// log.Printf("-www- [RegisterPage] seg: %v; \n", seg)
			if idx == len(segments)-1 {
				// log.Printf("-www- [RegisterPage] enter last node %v; \n", seg)
				// log.Printf("-www- --------- %v; %v \n", lowerSeg, prevSeg)

				// ~ 2 ~ overlap keywords: i.e.: order/OrderEdit ==> order/edit
				var (
					shortSeg  string
					finalSegs []string = []string{seg}
				)
				if lowerSeg == prevSeg {
					// eg: /order/Order --> leave it
				} else if strings.HasPrefix(lowerSeg, prevSeg) {
					// eg: /order/OrderList
					shortSeg = seg[len(prevSeg):]
					if shortSeg != "" {
						finalSegs = append(finalSegs, shortSeg)
					}

					// if isPage { // only page removes suffix "Index"
					// 	// eg: /order/OrderIndex --> /order
					// 	if strings.HasSuffix(lowerSeg, "index") {
					// 		shortSeg = shortSeg[:len(shortSeg)-len("index")]
					// 		if shortSeg == "" {
					// 			fmt.Printf("[[[[[[[[[[[[[[[[[[[[[[[[[[[[[[[[[[---   %v \n", finalSegs)
					// 			// fmt.Println("fallback")
					// 			// fallback TODO
					// 			currentSeg.Src = src
					// 			currentSeg.Path = strings.Join(segments, "/")
					// 			currentSeg.Proton = p
					// 			// should return here.
					// 		} else {
					// 			// eg: /order/OrderDetailIndex --> /order/detail
					// 			finalSegs = append(finalSegs, shortSeg)
					// 		}
					// 	}
					// }
				}
				// fmt.Printf(">>>>>> seg: %v shortSeg: %v\n", seg, shortSeg)

				// speical: for segment Index, go back to parent.
				if isPage {
					if strings.HasSuffix(lowerSeg, "index") {
						var trimlen = len(shortSeg) - len("index")
						if trimlen >= 0 {
							shortSeg = shortSeg[:trimlen]
						}
						if shortSeg == "" {
							// fallback TODO
							currentSeg.Src = src
							currentSeg.Path = strings.Join(segments, "/")
							currentSeg.Proton = p
							// should return here.
						} else {
							// eg: /order/OrderDetailIndex --> /order/detail
							finalSegs = append(finalSegs, shortSeg)
						}

					}
				}

				if strings.HasSuffix(lowerSeg, prevSeg) {
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

				// finally add segment struct to chains.
				for _, s := range finalSegs {
					// link segment together.
					segment = &ProtonSegment{
						Name:   s,
						Parent: currentSeg,
						Level:  idx,
						Src:    src,
						Path:   strings.Join(segments, "/"),
						Proton: p,
					}
					currentSeg.AddChild(segment)
					selector := []string{}
					// fmt.Printf(">>>>>>>>>>>selectorPrefix %v\n", selectorPrefix)
					selector = append(selector, selectorPrefix...)
					selector = append(selector, s)
					selectors = append(selectors, selector)
				}
				currentSeg = segment
				return

			} else {
				// the middle part
				segment = &ProtonSegment{
					Name:   seg,
					Parent: currentSeg,
					Level:  idx,
				}
				currentSeg.AddChild(segment)
				currentSeg = segment
				selectorPrefix = append(selectorPrefix, seg)
			}

		}
		prevSeg = lowerSeg
	}
	return
}

// Lookup the structure, find the right page/component.
// TODO performance
func (s *ProtonSegment) Lookup(url string) (segment *ProtonSegment, pageUrl string, err error) {
	logLookup("- - - [Lookup] '%v'\n", url)

	var level int
	var seg string
	trimedUrl := strings.Trim(url, " ")
	if !strings.HasSuffix(trimedUrl, "/") {
		trimedUrl += "/"
	}
	segments := strings.Split(trimedUrl, "/")

	// fmt.Println("--------------------------------------------------------------------------")
	// fmt.Println(trimedUrl)
	// fmt.Println(segments)
	// fmt.Println(len(segments))
	segment = s
	for level, seg = range segments {
		logLookup("- - - [Lookup] Step: Level %v Seg:[ %-10v ] segment:[ %-20v ]\n",
			level, seg, segment)
		if level == 0 && seg == "" { // skip the first / segment.
			// fmt.Printf("----- %v \n", segment)
			// first segment must be /
			continue
		}

		seg = strings.ToLower(seg)

		// logLookup("--- segment[%-10v].NumChildren = %v :: HasChildren([ %v ]) is %v\n",
		// 	segment.Seg, len(segment.Children),
		// 	seg, segment.HasChild(seg),
		// )
		// logLookup("-d- seg[%v], segment[%v]\n",
		// 	seg, segment,
		// )

		// try to go next level.
		if segment.Children == nil || len(segment.Children) == 0 || !segment.HasChild(seg) {
			// fmt.Printf("----- %v \n", segment)
			logLookup("- - - [Lookup] match finished.")
			break // stop match here
		} else {
			// find, go next level
			segment = segment.Children[seg]
		}
	}

	// fmt.Printf("+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+\n")
	// fmt.Printf("+ segments: \n")
	// for idx, s := range segments {
	// 	fmt.Printf("%v:%v\n", idx, s)
	// }
	pageUrl = strings.Join(segments[:level], "/")
	log.Printf("- - - [Lookup] 'pageurl is' %v\n", pageUrl)

	if nil == segment {
		err = errors.New("Lookup Failed.")
	}
	return
}

/* ________________________________________________________________________________
   Print Helper
*/

// String()
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
   protonType=[page|component]
   e.g. f("/got/builtin/pages/order/list") = "order/list"
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
