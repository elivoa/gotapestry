/*
  All Files uplaods here.
*/

package fileupload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"got/config"
	"got/core"
	"got/register"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func Register() {}
func init() {
	register.Page(Register,
		&FileUploadIndex{},
		&FileUploadTest{},
	)
}

// ________________________________________________________________________________
//
type FileUploadTest struct {
	core.Page
}

// ________________________________________________________________________________
//
type FileUploadIndex struct {
	core.Page
}

func (p *FileUploadIndex) Setup() {
	fmt.Println("... file upload ...")
}

func (p *FileUploadIndex) AfterRender() (string, string) {
	return "json", ""
}

func (p *FileUploadIndex) OnSuccess() (string, string) {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("onsuccess")
	FU(p.W, p.R)
	return "json", "{chedan: 'submit success'}"
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// ________________________________________________________________________________
// directly bind in got.go; can't change func name
//
func FU(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":

		var (
			pathkey = "__got_fileupload_path__"
			path    = ""
		)

		// ________________________________________________________________________________
		// 1. read files
		reader, err := r.MultipartReader()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// var path string = ""
		// copy each part to destination.
		result := make(map[string][]*FileInfo, 1)
		fileInfos := make([]*FileInfo, 0)

		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}

			// _______________________________
			// get path
			if part.FormName() == pathkey {
				bytes := make([]byte, 2048) // 2048 is enough to path
				i, err := part.Read(bytes)
				if err != nil {
					panic(err.Error())
				}
				path = string(bytes[0:i])
				fmt.Println("___________________________________________________________________")
				fmt.Println(path)
				continue
			}

			// ________________________________
			filename := part.FileName()
			//if part.FileName() is empty, skip this iteration.
			if filename == "" {
				continue
			}

			// make path
			abspath := filepath.Join(config.Config.ResourcePath, path) // /var/.../product-pic/12/
			err = os.MkdirAll(abspath, 0777)                           // security?
			if err != nil {
				panic(err.Error())
			}

			// create filename
			filename = strings.Replace(filename, ";", "_", -1)                 // file_name.ext
			filekey := fmt.Sprintf("%v%-3v", time.Millisecond, rand.Intn(999)) // 123456789

			fullFilename := strings.Join([]string{filekey, filename}, "-") // 123-filename.ext
			filePathKey := filepath.Join(path, fullFilename)               // /pic/12/123-filename.ext
			absFilename := filepath.Join(abspath, fullFilename)            // full

			fi := &FileInfo{
				Name: filePathKey,
				Type: part.Header.Get("Content-Type"),
			}
			fileInfos = append(fileInfos, fi)

			fmt.Println("Create file: " + absFilename)
			// write file here
			dst, err := os.Create(absFilename)
			defer dst.Close()

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if _, err := io.Copy(dst, part); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// fmt.Println("done 1")
		}
		result["files"] = fileInfos

		// handle returned json result
		b, err := json.Marshal(result)
		check(err)
		if redirect := r.FormValue("redirect"); redirect != "" {
			if strings.Contains(redirect, "%s") {
				redirect = fmt.Sprintf(
					redirect,
					escape(string(b)),
				)
			}
			http.Redirect(w, r, redirect, http.StatusFound)
			return
		}
		w.Header().Set("Cache-Control", "no-cache")
		jsonType := "application/json"
		if strings.Index(r.Header.Get("Accept"), jsonType) != -1 {
			w.Header().Set("Content-Type", jsonType)
		}
		fmt.Fprintln(w, string(b))

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// ________________________________________________________________________________
// Fileupload server side

const (
	WEBSITE           = "http://blueimp.github.io/jQuery-File-Upload/"
	MIN_FILE_SIZE     = 1       // bytes
	MAX_FILE_SIZE     = 5000000 // bytes
	IMAGE_TYPES       = "image/(gif|p?jpeg|(x-)?png)"
	ACCEPT_FILE_TYPES = IMAGE_TYPES
	EXPIRATION_TIME   = 300 // seconds
	THUMBNAIL_PARAM   = "=s80"
)

var (
	imageTypes      = regexp.MustCompile(IMAGE_TYPES)
	acceptFileTypes = regexp.MustCompile(ACCEPT_FILE_TYPES)
)

type FileInfo struct {
	Key    string `json:"-"`
	Name   string `json:"name"`
	Serial string `json:"serial"` // 12345678-filename.ext

	// not used
	Url          string `json:"url,omitempty"`
	ThumbnailUrl string `json:"thumbnail_url,omitempty"`
	Type         string `json:"type"`
	Size         int64  `json:"size"`
	Error        string `json:"error,omitempty"`
	DeleteUrl    string `json:"delete_url,omitempty"`
	DeleteType   string `json:"delete_type,omitempty"`
}

func (fi *FileInfo) ValidateType() (valid bool) {
	if acceptFileTypes.MatchString(fi.Type) {
		return true
	}
	fi.Error = "Filetype not allowed"
	return false
}

func (fi *FileInfo) ValidateSize() (valid bool) {
	if fi.Size < MIN_FILE_SIZE {
		fi.Error = "File is too small"
	} else if fi.Size > MAX_FILE_SIZE {
		fi.Error = "File is too big"
	} else {
		return true
	}
	return false
}

func (fi *FileInfo) CreateUrls(r *http.Request) {
	u := &url.URL{
		Scheme: r.URL.Scheme,
		Host:   "", //appengine.DefaultVersionHostname(c),
		Path:   "/",
	}
	uString := u.String()
	fi.Url = uString + escape(string(fi.Key)) + "/" +
		escape(string(fi.Name))
	fi.DeleteUrl = fi.Url + "?delete=true"
	fi.DeleteType = "DELETE"
	// if imageTypes.MatchString(fi.Type) {
	// 	servingUrl, err := image.ServingURL(
	// 		c,
	// 		fi.Key,
	// 		&image.ServingURLOptions{
	// 			Secure: strings.HasSuffix(u.Scheme, "s"),
	// 			Size:   0,
	// 			Crop:   false,
	// 		},
	// 	)
	// 	check(err)
	// 	fi.ThumbnailUrl = servingUrl.String() + THUMBNAIL_PARAM
	// }
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func escape(s string) string {
	return strings.Replace(url.QueryEscape(s), "+", "%20", -1)
}

func getFormValue(p *multipart.Part) string {
	var b bytes.Buffer
	io.CopyN(&b, p, int64(1<<20)) // Copy max: 1 MiB
	return b.String()
}
