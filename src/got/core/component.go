/*
Tapestry like Components Design

Use in templates like this:

  {{t_select `data:"value" ` }}

TODO:
  . How to
  . Can I modify html node in tempaltes? linke this:

    <select type="{{c}}" data="jiujiujiujiujiujiujiujiu" />

*/

package core

/*_______________________________________________________________________________
  Component
*/
type IComponent interface {
	Protoner
}

type Component struct {
	Proton
	Tid string // component id
}

func (c *Component) Kind() Kind {
	return COMPONENT
}
