package main

import (
	"fmt"
	"image-trans/primitive"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		html := ` <html> <body> <form action = "/upload" method = "post" enctype = "multipart/form-data">
				<input type = "file" name = "image">
				<button type = "submit">upload image</button>
				</form>
				</body></html>`
		fmt.Fprint(w,html)
	})
	mux.HandleFunc("/upload", func (w http.ResponseWriter, r *http.Request){
		file, header , err := r.FormFile("image")
		if err != nil {
			http.Error(w,err.Error(),http.StatusBadRequest)
			return
		}
		defer file.Close()
		ext := filepath.Ext(header.Filename)[1:]

		out , err := primitive.Transform(file,50)
		if err != nil {
			http.Error(w,err.Error(), http.StatusInternalServerError)
		}
		switch ext {
			case "jpg":
				fallthrough
			case "jpeg":
				w.Header().Set("Content-Type", "image/jpeg")
			case "png":
				w.Header().Set("Content-Type", "image/png")
			default:
				http.Error(w, "invalid type ", http.StatusBadRequest)
		}

		_, err = io.Copy(w,out)

		if err != nil {
			http.Error(w,err.Error(), http.StatusInternalServerError)
		}
	})

	log.Fatal(http.ListenAndServe(":3000",mux))
}

