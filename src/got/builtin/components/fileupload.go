/*
 Component: fileupload
  Should include in a form.
  Support multiple file upload in the same time.
*/

package components

import (
	"got/core"
	"path/filepath"
)

// ________________________________________________________________________________
// Must used within a form. Will generate one or more Input:hidden.
//
type FileUpload struct {
	core.Component
	Tid     string   // Component Id / can be used as client id
	Name    string   // name property for INPUT, used to submit key with form.
	Folder  string   // subfolder in under the root picture folder.
	Restore []string // files to restore, generate some
	Style   string   // style
	Class   string   // class

	// TODO
	MaxFiles int // >0 means max files can upload. else many
}

func (p *FileUpload) Activate() {
	p.MaxFiles = 0 // can upload many files.
}

func (p *FileUpload) PictureLink(filekey string) string {
	return filepath.Join("/pictures", filekey)
}

func (p *FileUpload) FileLink(folder string, filekey string) string {
	return filepath.Join(folder, filekey)
}
