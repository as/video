package main

import (
	"io"
	"log"
	"os"

	"github.com/as/video/mjpeg"
	"image/png"
)

func main() {
	fd, err := os.Open(os.Args[1])
	no(err)
	sc := mjpg.Open(fd)
	for sc.Scan() {
		img, _, err := sc.Frame().Decode()
		if err != nil {
			log.Fatalln(err)
		}
		png.Encode(os.Stdout, img)
	}
	if sc.Err() != nil || sc.Err() != io.EOF {
		log.Fatalln(sc.Err())
	}
}
func no(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
