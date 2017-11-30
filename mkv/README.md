# mkv
Mkv implements operations on the MKV container format. 

```
// This program reads an MKV from standard input and emits
// all of the tags present in the MKV to standard output
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/as/video/mkv"
)

func main() {
	br := bufio.NewReader(os.Stdin)
	conf, err := mkv.DecodeConfig(br)
	if err != nil {
		log.Printf("decode: %s\n", err)
	}
	for k, v := range conf.Tags {
		fmt.Printf("%q=%q\n", k, v)
	}
}
```

# TODO

- Split
- Join
- Recursive Descent

# Trivia
- MKV is based on EBML, a human-readable simplification of XML
- XML is a deprecated markup language based on slashes, angle brackets, and squigglies
