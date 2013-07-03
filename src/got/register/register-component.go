package register

import (
	"bytes"
	"fmt"
	"got/core"
	"got/core/lifecircle"
	"got/debug"
	"got/templates"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

/* ________________________________________________________________________________
   ComponentRegister
*/
var Components = ProtonSegment{Name: "/"}

func Component(f func(), components ...core.IComponent) int {
	for _, c := range components {
		url := makeUrl(f, c)
		selectors := Components.Add(url, c, "component")

		for _, selector := range selectors {
			lowerKey := strings.ToLower(strings.Join(selector, "/"))
			templates.RegisterComponent(lowerKey, componentLifeCircle(lowerKey))
		}
	}
	return len(components)
}

/* ________________________________________________________________________________
   Execute Components
*/

/*
  Handler method
  Return: string or template.HTML
*/
func componentLifeCircle(name string) func(...interface{}) interface{} {

	return func(params ...interface{}) interface{} {

		log.Printf("-620- [flow] Render Component %v ....", name)

		// 1. find base component type
		result, err := Components.Lookup(name)
		if err != nil || result.Segment == nil {
			panic(fmt.Sprintf("Component %v not found!", name))
		}
		if len(params) < 1 {
			panic(fmt.Sprintf("First parameter of component must be '$' (container)"))
		}

		// 2. find container page/component
		container := params[0].(core.IProton)

		// 3. create lifecircle controler
		lcc := lifecircle.NewComponentFlow(container, result.Segment.Proton, params[1:])
		lcc.Flow()
		handleComponentReturn(lcc, result.Segment)
		return template.HTML(lcc.String)
	}
}

// handle component return
func handleComponentReturn(lcc *lifecircle.LifeCircleControl, seg *ProtonSegment) {
	// no error, no templates return or redirect.
	if seg != nil && lcc.Err == nil && lcc.ResultType == "" {
		// find default tempalte to return
		key, tplPath := LocateGOTComponentTemplate(seg.Src, seg.Path)
		debug.Log("-756- [ComponentTemplateSelect] %v -> %v", key, tplPath)
		_, err := templates.GotTemplateCache.Get(key, tplPath)
		if nil != err {
			lcc.Err = err
		} else {
			fmt.Println("render component tempalte " + key)

			var buffer bytes.Buffer
			err = templates.RenderGotTemplate(&buffer, key, lcc.Proton)
			if err != nil {
				lcc.Err = err
			}
			lcc.String = buffer.String()
		}
	}

	if lcc.Err != nil {
		debug.Error(lcc.Err)
		http.Error(lcc.W, fmt.Sprint(lcc.Err), http.StatusInternalServerError)
	}
}

// ________________________________________________________________________________
// Locate Templates
// return (template-key, template-file-path); TODO: performance issue
func LocateGOTComponentTemplate(src string, path string) (string, string) {
	appConfig := Apps.Get(src)
	if appConfig == nil {
		panic(fmt.Sprintf("Can't find APP Config %v", src))
	}
	key := fmt.Sprintf("c_%v:%v", src, path)
	templateFilePath := filepath.Join(appConfig.FilePath, "components", path) + ".html"
	return key, templateFilePath
}

/* ________________________________________________________________________________
   ComponentSegment
*/
// ________________________________________________________________________________
// parse url and put it into segments.
//
// func (s *ProtonSegment) Add(baseUrl string, c core.IComponent) {
// 	src, segments := trimPathSegments(baseUrl, "components")

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

// 					// component don't process index.
// 					// if strings.HasSuffix(lowerSeg, "index") {
// 					// 	shortSeg = shortSeg[:len(shortSeg)-len("index")]
// 					// 	if shortSeg == "" {
// 					// 		// fallback TODO
// 					// 		currentSeg.Src = src
// 					// 		currentSeg.Path = strings.Join(segments, "/")
// 					// 		currentSeg.Proton = c
// 					// 		// should return here.
// 					// 	} else {
// 					// 		// eg: /order/OrderDetailIndex --> /order/detail
// 					// 		finalSegs = append(finalSegs, shortSeg)
// 					// 	}
// 					// }
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
// 						Proton: c,
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

// // Lookup the structure, find the right page/component.
// // TODO performance
// func (s *ProtonSegment) Lookup(url string) (segment *ProtonSegment, pageUrl string, err error) {
// 	logLookup("- - - [Page Lookup] '%v'\n", url)

// 	var level int
// 	var seg string
// 	trimedUrl := strings.Trim(url, " ")
// 	if !strings.HasSuffix(trimedUrl, "/") {
// 		trimedUrl += "/"
// 	}
// 	segments := strings.Split(trimedUrl, "/")

// 	// fmt.Println("--------------------------------------------------------------------------")
// 	// fmt.Println(trimedUrl)
// 	// fmt.Println(segments)
// 	// fmt.Println(len(segments))
// 	for level, seg = range segments {
// 		logLookup("- - - [Page Lookup] Step: Level %v Seg:[ %-10v ] segment:[ %-20v ]\n",
// 			level, seg, segment)
// 		if level == 0 { // skip the first / segment.
// 			segment = s
// 			continue
// 		}

// 		seg = strings.ToLower(seg)

// 		// logLookup("--- segment[%-10v].NumChildren = %v :: HasChildren([ %v ]) is %v\n",
// 		// 	segment.Name, len(segment.Children),
// 		// 	seg, segment.HasChild(seg),
// 		// )
// 		// logLookup("-d- seg[%v], segment[%v]\n",
// 		// 	seg, segment,
// 		// )

// 		// try to go next level.
// 		if segment.Children == nil || len(segment.Children) == 0 || !segment.HasChild(seg) {
// 			logLookup("- - - [Page Lookup] match finished.")
// 			break // stop match here
// 		} else {
// 			// find, go next level
// 			segment = segment.Children[seg]
// 		}
// 	}

// 	// fmt.Printf("+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+\n")
// 	// fmt.Printf("+ segments: \n")
// 	// for idx, s := range segments {
// 	// 	fmt.Printf("%v:%v\n", idx, s)
// 	// }
// 	pageUrl = strings.Join(segments[:level], "/")
// 	log.Printf("- - - [Page Lookup] 'pageurl is' %v\n", pageUrl)

// 	if nil == segment {
// 		err = errors.New("Lookup Failed.")
// 	}
// 	return
// }
