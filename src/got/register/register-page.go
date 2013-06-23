package register

import (
	"got/core"
)

/* ________________________________________________________________________________
   PageRegister
*/
var Pages = ProtonSegment{Name: "/"}

/* ________________________________________________________________________________
   Register a Page,
   params:
     - f used to locate a page(reflect has no use).
     - pages are Pages located in that folder.
*/
func Page(f func(), pages ...core.IPage) int {
	for _, p := range pages {
		url := makeUrl(f, p)
		Pages.Add(url, p, "page")
	}
	return len(pages)
}

/* ________________________________________________________________________________
   Segment, like a trie.
*/
// TODO: Add Cache
// TODO: refactor

// type PageSegment struct {
// 	ProtonSegment
// }

// ________________________________________________________________________________
// parse url and put it into segments.
//
// func (s *PageSegment) Add(baseUrl string, p core.IPage) {
// 	src, segments := trimPathSegments(baseUrl, "pages")

// 	debug.Log("-___- [RegisterPage] register %v::%v url:%v", src, segments, baseUrl)

// 	// add to register
// 	var (
// 		currentSeg = &s.ProtonSegment
// 		prevSeg    = "//nothing//"
// 	)
// 	for idx, seg := range segments {
// 		lowerSeg := strings.ToLower(seg)
// 		var segment *ProtonSegment
// 		if currentSeg.HasChild(seg) {
// 			// detect conflict
// 			existSeg := s.Children[seg]
// 			if existSeg.Src != "" && existSeg.Src != src {
// 				log.Fatalf("Conflict of Page defination %v.\n", baseUrl)
// 			}
// 			currentSeg = existSeg
// 		} else {
// 			// log.Printf("-www- [RegisterPage] seg: %v; \n", seg)

// 			// create segment{} and add to chain
// 			if idx == len(segments)-1 {
// 				// log.Printf("-www- [RegisterPage] enter last node %v; \n", seg)

// 				// ~ 2 ~ overlap keywords: i.e.: order/OrderEdit ==> order/edit
// 				var (
// 					shortSeg  string
// 					finalSegs []string = []string{seg}
// 				)
// 				if lowerSeg == prevSeg {
// 					// eg: /order/Order --> leave it
// 				} else if strings.HasPrefix(lowerSeg, prevSeg) {
// 					// eg: /order/OrderList
// 					shortSeg = seg[len(prevSeg):]
// 					finalSegs = append(finalSegs, shortSeg)

// 					// eg: /order/OrderIndex --> /order
// 					if strings.HasSuffix(lowerSeg, "index") {
// 						shortSeg = shortSeg[:len(shortSeg)-len("index")]
// 						if shortSeg == "" {
// 							// fallback TODO
// 							currentSeg.Src = src
// 							currentSeg.Path = strings.Join(segments, "/")
// 							currentSeg.Proton = p
// 							// should return here.
// 						} else {
// 							// eg: /order/OrderDetailIndex --> /order/detail
// 							finalSegs = append(finalSegs, shortSeg)
// 						}
// 					}
// 				}

// 				if strings.HasSuffix(lowerSeg, prevSeg) {
// 					finalSegs = append(finalSegs, seg[:len(seg)-len(prevSeg)])
// 				}

// 				// finally add segment struct to chains.
// 				for _, s := range finalSegs {
// 					// link segment together.
// 					segment = &ProtonSegment{
// 						Name:   s,
// 						Parent: currentSeg,
// 						Level:  idx,
// 						Src:    src,
// 						Path:   strings.Join(segments, "/"),
// 						Proton: p,
// 					}
// 					currentSeg.AddChild(segment)
// 				}
// 				currentSeg = segment

// 			} else {
// 				// log.Printf("-www- [RegisterPage] normal path node %v; \n", seg)

// 				// the middle path
// 				segment = &ProtonSegment{
// 					Name:   seg,
// 					Parent: currentSeg,
// 					Level:  idx,
// 				}
// 				currentSeg.AddChild(segment)
// 				currentSeg = segment
// 			}

// 		}
// 		prevSeg = lowerSeg
// 	}
// }
