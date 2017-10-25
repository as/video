// Package mjpeg provides a fast splitter for a Motion-JPEG A dual-field
// stream. This package conforms to the standard located here:
//
// https://developer.apple.com/standards/qtff-2001.pdf
//
// The Frame does not decode the jpeg image data until Decode is called
// and this optimization improves performance 100-fold on systems where
// the jpeg data is forwarded.
//
// ffmpeg  -i video.mp4 -f mjpeg -bsf:v mjpegadump -bsf:v mjpeg2jpeg -> out3.mjpg
package mjpg

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	_ "image/jpeg"
	"io"
	"time"
)

var (
	ErrBadSOI   = errors.New("bad start of image marker")
	ErrBadAPP1  = errors.New("bad APP1 marker")
	ErrBadCLen  = errors.New("bad content length")
	ErrBadMagic = errors.New("bad magic number")
)

func Open(r io.Reader) *Scanner {
	return &Scanner{
		br: bufio.NewReader(r),
	}
}

type Scanner struct {
	hdr   HdrA
	br    *bufio.Reader
	frame Frame
	err   error
	first *time.Time
}

func (s *Scanner) readhdr() error {
	hdr := HdrA{}
	err := binary.Read(s.br, binary.BigEndian, &hdr)
	if err != nil {
		return err
	}
	s.hdr = hdr
	return hdr.Check()
}

func (s *Scanner) Scan() bool {
	if s.first == nil {
		t := time.Now()
		s.first = &t
	}
	s.err = s.readhdr()
	if s.err != nil {
		return false
	}
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s.hdr)
	_, s.err = io.CopyN(buf, s.br, int64(s.hdr.PaddedSize-22))
	if s.err != nil {
		return false
	}
	s.frame = &frame{data: buf}
	return true
}

func (s *Scanner) Frame() Frame {
	return s.frame
}

func (s *Scanner) Err() error {
	return s.err
}

func (s *Scanner) TimeStamp(f Frame) time.Duration {
	if s.first == nil {
		return 0
	}
	return f.Time().Sub(*s.first)
}

// Frame supports returning the raw MJPEG data or decoding it
// into an image.Image.
type Frame interface {
	Time() time.Time
	Decode() (img image.Image, fm string, err error)
	Raw() *bytes.Buffer
}

type frame struct {
	ts   time.Time
	data *bytes.Buffer
}

func (f *frame) Decode() (img image.Image, fm string, err error) {
	return image.Decode(bytes.NewReader(f.data.Bytes()))
}
func (f *frame) Raw() *bytes.Buffer {
	return f.data
}
func (f *frame) Time() time.Time {
	return f.ts
}

type HdrA struct {
	SOI        uint16
	APP1       uint16
	CLen       uint16
	R1         uint32
	Magic      uint32
	Size       uint32
	PaddedSize uint32
}

func (h HdrA) Check() error {
	if h.SOI != 0xffd8 {
		return ErrBadSOI
	}
	if h.APP1 != 0xffe1 {
		return ErrBadAPP1
	}
	if h.CLen != 0x002a {
		return ErrBadCLen
	}
	if h.Magic != 0x6d6a7067 {
		return ErrBadMagic
	}
	return nil
}
