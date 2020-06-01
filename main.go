package main

import (
	"errors"
	"fmt"
	"html/template"
	"image-trans/primitive"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
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
		a, err := gen(file, ext,10,primitive.ModePolygon)
		if err != nil {
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
		file.Seek(0,0)
		b, err := gen(file, ext,10,primitive.ModeCircle)
		if err != nil {
			panic(err)
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
		file.Seek(0,0)
		c, err := gen(file, ext,10,primitive.ModeTriangle)
		if err != nil {
			panic(err)
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
		file.Seek(0,0)
		d, err := gen(file, ext,10,primitive.ModeBeziers)
		if err != nil {
			panic(err)
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
		html := `<html><body>
				{{range .}}
				<img src = "/{{.}}">
				{{end}}
				</body>	</html>`

		tpl := template.Must(template.New("").Parse(html))
		tpl.Execute(w, []string{a,b,c,d})

		//redirectUrl := fmt.Sprintf("/%s", a)
		//http.Redirect(w, r, redirectUrl, http.StatusFound)
	})
	mux.Handle("/img/", http.StripPrefix("/img", http.FileServer(http.Dir("./img/"))))
	log.Fatal(http.ListenAndServe(":3000", mux))
}

func gen(file multipart.File, ext string, num int, shape primitive.PrimitiveMode) (string,error) {
	out, err := primitive.Transform(file, num, primitive.WithMode(shape))
	outfile, err := newTempFile("", ext)
	if err != nil {
		return "", err
	}
	defer outfile.Close()
	io.Copy(outfile, out)
	return outfile.Name(),nil
}

func newTempFile(prefix , ext string ) (*os.File, error) {
	file ,err := ioutil.TempFile("./img", prefix)
	if  err != nil {
		return nil, errors.New("primitive: failed to create temporary file")
	}
	defer os.Remove(file.Name())
	return os.Create(fmt.Sprintf("%s.%s",file.Name(),ext))
}

