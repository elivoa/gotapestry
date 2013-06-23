package got

/*___________________________________________
  GOT framework helper funcs
*/

const (
	ReturnTemplate   = "template"
	ReturnJson       = "json"
	ReturnPlain_Text = "plaintext"
	ReturnRedirect   = "redirect"
)

type ReturnMeta struct {
	Type string //
	Data string //
}

func ReturnText(returnType string, data string) *ReturnMeta {
	return &ReturnMeta{
		Type: returnType,
		Data: data,
	}
}
