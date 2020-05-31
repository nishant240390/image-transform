package primitive

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type PrimitiveMode int

const
(
	combo PrimitiveMode = iota
	ModeTriangle
	ModeRect
	ModeEllipse
	ModeCircle
	ModeRotatedrect
	ModeBeziers
	ModeRotatedellipse
	ModePolygon
)

func Transform (image io.Reader , numShapes int, opts ...func() []string)(io.Reader, error) {

	in, err := newTempFile("in_","png")
	if err != nil {
		return nil, errors.New("Failed to create temporary input file")
	}
	defer os.Remove(in.Name())
	out , err := newTempFile("out_","png")
	if err != nil {
		return nil, errors.New("Failed to create temporary input file")
	}
	defer os.Remove(out.Name())

	_,err = io.Copy(in , image)
	if err != nil {
		return nil, err
	}

	std, err := primitive(in.Name(),out.Name(), numShapes, ModePolygon)
	if err != nil {
		return nil,err
	}
	fmt.Println(std)

	 b := bytes.NewBuffer(nil)
	 _, err = io.Copy(b , out)

	 if err != nil {
	 	return nil, err
	}
	return b, nil
}

func primitive (inputFile string , outputFile string ,numShapes int , mode PrimitiveMode )(string, error) {
	args := fmt.Sprintf("-i %s -o %s -n %d -m %d",inputFile,outputFile,numShapes,mode)
	byte , err := exec.Command("primitive",strings.Fields(args)...).CombinedOutput()
	return string(byte),err
}
func newTempFile(prefix , ext string ) (*os.File, error) {

	 file ,err := ioutil.TempFile("", prefix)
	 if  err != nil {
	 	return nil, errors.New("primitive: failed to create temporary file")
	 }
	 defer os.Remove(file.Name())
	 return os.Create(fmt.Sprintf("%s.%s",file.Name(),ext))

}
