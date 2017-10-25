package main

import (
	"io"
	"log"
	"os"

	"github.com/as/video/mjpeg"
)

func main() {
	fd, err := os.Open(os.Args[1])
	no(err)
	sc := mjpg.Open(fd)
	for sc.Scan() {
		io.Copy(os.Stdout, sc.Frame().Raw())
	}
	if sc.Err() != nil || sc.Err() != io.EOF {
		log.Fatalln(sc.Err())
	}
}
func no(err error){
	if err !=nil{
		log.Fatalln(err)
	}
}