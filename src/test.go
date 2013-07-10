package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/upload", uploadHandler)

	//static file handler.
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("static"))))

	//Listen on port 8080
	http.ListenAndServe(":8080", nil)
}

//Compile templates on start
var templates = template.Must(template.ParseFiles("template/upload.html"))

//Display the named template
func display(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl+".html", data)
}

//This is where the action happens.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//GET displays the upload form.
	case "GET":
		display(w, "upload", nil)

	//POST takes the uploaded file(s) and saves it to disk.
	case "POST":
		//get the multipart reader for the request.
		debuglog := true
		if debuglog {
			fmt.Println("``````````````````````````````````````````````````````````````````````````````````````````")
			fmt.Printf("r.Method=%v\n", r.Method)
			err := r.ParseForm()
			if err != nil {
				panic(err.Error())
			}
			err = r.ParseMultipartForm(1024 * 1024 * 10)
			if err != nil {
				panic(err.Error())
			}

			fmt.Println(r.Form)
			fmt.Println(r.PostForm)
			fmt.Println(r.Body)
		}

		reader, err := r.MultipartReader()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("for parts")
		//copy each part to destination.
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}

			//if part.FileName() is empty, skip this iteration.
			if part.FileName() == "" {
				continue
			}
			fmt.Println("create folder???  " + part.FileName())
			dst, err := os.Create("/tmp/uploadfile-" + part.FileName())
			defer dst.Close()

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if _, err := io.Copy(dst, part); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Println("done 1")
		}
		display(w, "upload", "Upload successful.")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
