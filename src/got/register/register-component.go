package register

import ()

/* ________________________________________________________________________________
   ComponentRegister
*/
var Components = ProtonSegment{Name: "/"}

// // Comopnent register a component to got resgitry.
// func Component(f func(), components ...core.Componenter) int {
// 	for _, c := range components {
// 		url := makeUrl(f, c)
// 		// TODO has space to improve.
// 		selectors := Components.Add(url, c)
// 		for _, selector := range selectors {
// 			lowerKey := strings.ToLower(strings.Join(selector, "/"))
// 			templates.RegisterComponent(lowerKey, componentLifeCircle(lowerKey))
// 		}
// 	}
// 	return len(components)
// }
