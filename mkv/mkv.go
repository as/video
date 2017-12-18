//
package mkv

import (
	"bufio"
	"bytes"
	"encoding"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/bits"
	"strings"
	"time"
)

const maxSpins = 1e10

type Decoder interface {
	Decode(r io.Reader) error
}

type Config struct {
	Title         string
	Width, Height int
	//	ColorModel    color.Model	// This will not be straightforward
	FourCC   string
	Tags     map[string]string
	Duration time.Duration
	Codec    []*CodecInfo
	Muxer    string
	Writer   string
}

func (c CodecInfo) String() string {
	return fmt.Sprintf("%#v\n", c)
}

// DecodeConfig decodes the MKV config. Unlike image.DecodeConfig,
// this one reads the entire MKV (not just the headers). Metadata
// can reside in the trailer or in arbitrary locations of the file
func DecodeConfig(r io.Reader) (*Config, error) {
	var err error

	s := NewScanner(r)
	tags := make(map[string]string)
	key := ""
	readstring := func(len int64) string {
		str, _ := s.ReadString(int(len))
		return str
	}
	var (
		dx, dy               int64
		tcs                  int64
		dur                  float64
		fcc                  = make([]byte, 4)
		muxer, writer, title string
		codec                CodecInfo
		codecs               []*CodecInfo
		ncodecs              int
	)
	// Scan through the MKV, looking for elements we want
	// including their predecessors.
Scan:
	for i := 0; i < maxSpins; i++ {
		e, a, err := s.Next()
		if err != nil {
			break
		}
		nm := Name(e)
		log.Println(nm, a)
		switch nm {
		case "PixelWidth":
			s.Decode(&dx)
		case "PixelHeight":
			s.Decode(&dy)
		case "CodecID":
			if ncodecs > len(codecs) {
				codecs = append(codecs, &codec)
			}
			id := readstring(a)
			codec = codectab[id]
			ncodecs++
		case "CodecName":
			codec.Name = readstring(a)
		case "CodecInfoURL":
			codec.URL = readstring(a)
		case "CodecDownloadURL":
			codec.Download = readstring(a)
		case "CodecSettings":
			codec.Settings = readstring(a)
		case "Muxer":
			muxer = readstring(a)
		case "Writer":
			writer = readstring(a)
		case "Title":
			title = readstring(a)
		case "ColourModel":
			s.Decode(fcc)
		case "TimecodeScale":
			s.Decode(&tcs)
		case "Duration":
			s.Decode(&dur)
		case "Void":
			s.ParseVoid()
		case "SimpleTag", "Tags", "Tag":
			// nothing, just descend without advancing
		case "TagName":
			key = readstring(a)
		case "TagString":
			tags[key] = readstring(a)
		default:
			if err = s.Advance(a); err != nil {
				break Scan
			}
		}
	}

	if ncodecs > len(codecs) {
		codecs = append(codecs, &codec)
	}
	log.Printf("dur %v tcs %v\n", dur, tcs)
	return &Config{
		Tags:     tags,
		Muxer:    muxer,
		Writer:   writer,
		Title:    title,
		FourCC:   string(fcc),
		Width:    int(dx),
		Height:   int(dy),
		Duration: time.Duration(float64(dur) * float64(tcs)),
		Codec:    codecs,
	}, err
}

// NewScanner
func NewScanner(r io.Reader) *Scanner {
	switch r := r.(type) {
	case *bufio.Reader:
		return &Scanner{r: r}
	}
	return &Scanner{r: bufio.NewReader(r)}
}

type Scanner struct {
	r       *bufio.Reader
	elem    int
	advance int64
	err     error
}

// Decode decodes the next element under the scanner's current offset.
// If v implements the Decoder or encoding.BinaryUnmarshaler interface,
// it uses those corresponding method sets to complete the decoding.
// Decode ensures the scanner always advances the correct number of
// bytes in the underlying stream, regardless of the behavior of v.
//
// TODO(as): document the encoding of integers
func (s *Scanner) Decode(v interface{}) (err error) {
	l := &io.LimitedReader{
		N: s.advance,
		R: s.r,
	}
	switch v := v.(type) {
	case Decoder:
		err = v.Decode(l)
	case encoding.BinaryUnmarshaler:
		b, _ := ioutil.ReadAll(l)
		err = v.UnmarshalBinary(b)
	case *float32, *float64:
		if l.N == 3 {
			l.N = 4
		}
		err = binary.Read(l, binary.BigEndian, v)
	case *int64:
		switch s.advance {
		case 1:
			var b byte
			err = binary.Read(l, binary.BigEndian, &b)
			*v = int64(b)
		case 2:
			var b uint16
			err = binary.Read(l, binary.BigEndian, &b)
			*v = int64(b)
		case 3:
			data, _ := ioutil.ReadAll(l)
			var b uint32
			if s.advance == 3 {
				data = append([]byte{0}, data...)
			}
			err = binary.Read(bytes.NewReader(data), binary.BigEndian, &b)
			*v = int64(b)
		case 4:
			var b uint32
			err = binary.Read(l, binary.BigEndian, &b)
			*v = int64(b)
		}
	case []byte:
		_, err = io.ReadAtLeast(l, v, len(v))
	}
	if l.N > 0 {
		// Ensure that the parser is always element-aligned by advancing
		// through the slop the decoder functions missed
		ioutil.ReadAll(l)
	}
	return err
}

func (s *Scanner) extract() (nz int, advance int64, v uint32) {
	p, _ := s.r.Peek(4)
	err := binary.Read(bytes.NewReader(p), binary.BigEndian, &v)
	if err != nil {
		return -1, -1, 0
	}
	nz = int(bits.LeadingZeros32(v))
	advance = int64(nz + 1)
	return int(nz), advance, v
}
func (s *Scanner) Read(p []byte) (n int, err error) {
	return s.r.Read(p)
}
func (s *Scanner) ReadString(n int) (string, error) {
	m := make([]byte, n)
	_, err := io.ReadAtLeast(s.r, m, n)
	return strings.Trim(string(m), "\x00"), err
}
func (s *Scanner) ReadByte() (byte, error) {
	return s.r.ReadByte()
}
func (s *Scanner) UnreadByte() error {
	return s.r.UnreadByte()
}

func (s *Scanner) ParseVoid() {
	if Name(s.elem) != "Void" {
		return
	}
	for {
		c, _ := s.ReadByte()
		if c != 0 {
			s.UnreadByte()
			break
		}
	}
}

func (s *Scanner) Advance(n int64) (err error) {
	for i := int64(0); i < n; i++ {
		_, err = s.r.ReadByte()
		if err != nil {
			break
		}
	}
	return err
}
func (s *Scanner) Next() (elem int, advance int64, err error) {
	_, advance, v := s.extract()
	if advance == -1 {
		return 0, 0, io.EOF
	}
	if err = s.Advance(advance); err != nil {
		return 0, 0, err
	}

	elem = int(v)
	elem = elem >> ((4 - uint(advance)) * 8)
	_, advance, v = s.extract()
	err = s.Advance(advance)
	L := uint(advance)
	v = v << L >> L
	v >>= (4 - L) * 8
	s.elem, s.advance, s.err = elem, int64(v), err
	return s.elem, s.advance, s.err
}

func at(r *bufio.Reader) int64 {
	return int64(2)
}
