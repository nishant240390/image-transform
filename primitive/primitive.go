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

func WithMode(mode PrimitiveMode) func() []string {
	return func()[]string {
		return []string{"-m", fmt.Sprintf("%d",mode)}
	}
}

func Transform (image io.Reader , numShapes int, opts ...func() []string)(io.Reader, error) {

	var args []string
	for _,opt := range opts {
		args = append(args, opt()...)
	}

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

	std, err := primitive(in.Name(),out.Name(), numShapes, args...)
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

func primitive (inputFile string , outputFile string ,numShapes int , args ...string )(string, error) {
	argString := fmt.Sprintf("-i %s -o %s -n %d",inputFile,outputFile,numShapes)
	args = append(strings.Fields(argString), args...)
	byte , err := exec.Command("primitive",args...).CombinedOutput()
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
