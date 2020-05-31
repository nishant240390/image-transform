package main

import (
	"errors"
	"fmt"
	"image-trans/primitive"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

		out , err := primitive.Transform(file,50, primitive.WithMode(primitive.ModeRect))

		outfile , err := newTempFile("",ext)

		if err != nil {
			http.Error(w,err.Error(), http.StatusInternalServerError)
		}

		defer outfile.Close()

		io.Copy(outfile,out)

		redirectUrl := fmt.Sprintf("/%s",outfile.Name())

		http.Redirect(w,r,redirectUrl, http.StatusFound)

	})

	mux.Handle("/img/", http.StripPrefix("/img",http.FileServer(http.Dir("./img/"))))
	log.Fatal(http.ListenAndServe(":3000",mux))
}
func newTempFile(prefix , ext string ) (*os.File, error) {

	file ,err := ioutil.TempFile("./img", prefix)
	if  err != nil {
		return nil, errors.New("primitive: failed to create temporary file")
	}
	defer os.Remove(file.Name())
	return os.Create(fmt.Sprintf("%s.%s",file.Name(),ext))

}

