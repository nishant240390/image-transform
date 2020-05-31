package main

import (
	"image-trans/primitive"
	"io"
	"log"
	"os"
)

func main() {

	f, err := os.Open("/Users/nishantagarwal/Desktop/profile.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

    out, err := primitive.Transform(f, 50)

    if err != nil {
    	panic(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	os.Remove("out.png")
     outfile , err := os.Create("out.png")
     if err != nil {
     	panic(err)
	 }
    io.Copy(outfile,out)
}

