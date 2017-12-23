// This program reads an MKV from standard input and emits
// all of the tags present in the MKV to standard output
package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"

	"github.com/as/video/mkv"
)

func main() {
	video := new(bytes.Buffer)
	r := io.TeeReader(os.Stdin, video)
	br := bufio.NewReader(r)
	conf, err := mkv.DecodeConfig(br)
	if err != nil {
		log.Printf("decode: %s\n", err)
	}
	log.Printf("config: %+v\n", conf)

	for k, v := range conf.Tags {
		log.Printf("%q=%q\n", k, v)
	}
	io.Copy(os.Stdout, &io.LimitedReader{
		video,
		int64(video.Len() - br.Buffered()),
		// Substract the length of the video from
		// what is buffered in the reader. These are
		// data that haven't been sent to the decoder
		// but read out of the response.
	})
}
