package mkv

import (
	"bufio"
	"fmt"
	"log"
	"time"
)

type Flag byte

const (
	FKeyFrame    Flag = 0x80
	FNoLace      Flag = 0x06
	FNoDuration  Flag = 0x08
	FDiscardable Flag = 0x01
)

type Block struct {
	Track    Vint
	TimeCode time.Duration
	Flag

	// wow, now we're getting to the good stuff
	// ... maybe later
}

type BitReader struct {
	tmp uint32
	c   byte
	n   uint
	r   *bufio.Reader
}

type X struct {
	br *bufio.Reader
}

func (x *X) Next() (err error) {
	var tmp [4]byte
	for err == nil {
		log.Println(err)
		data, err := x.br.Peek(4)
		if err != nil {
			return nil
		}
		n := 1
		switch {
		case string(data) == "\x00\x00\x00\x01":
			n = 4
		case string(data[0:3]) == "\x00\x00\x01":
			n = 3
		}
		n, err = x.br.Read(tmp[:n])
		if n == 1 {
			continue
		}
		b, err := x.br.ReadByte()
		if err != nil {
			return err
		}
		nal, err := parseNAL(b)
		log.Printf("nal: %#v\n", nal)
		if err != nil {
			return err
		}
	}
	return err
}

func StartCodePrefix(r *bufio.Reader) (code uint32, err error) {
	bitr := &BitReader{
		tmp: 0,
		r:   r,
	}
	var v uint32
	for err == nil {
		v, err = bitr.Next()
		//fmt.Printf("@%032b %d\n", v, bitr.n)
		if v == 1 {
			for i := 0; i < 8; i++ {
				v, err = bitr.Next()
			}
			fmt.Printf("	%032b (%x) %d\n", v, v, bitr.n)
			fmt.Println(parseNAL(byte(v)))
		}
	}
	return v, err
}
func (b *BitReader) yNext() uint32 {
	var err error
	if b.n%8 == 0 {
		b.c, err = b.r.ReadByte()
		if err != nil {
			log.Fatalln(err)
		}
	}
	b.tmp = b.tmp>>1 | uint32(b.c>>7)<<31
	b.c <<= 1
	b.n++
	return b.tmp
}
func (b *BitReader) Next() (uint32, error) {
	var err error
	if b.n%8 == 0 {
		b.c, err = b.r.ReadByte()
	}
	b.tmp = b.tmp<<1 | uint32(b.c>>7)
	b.c <<= 1
	b.n++
	return b.tmp, err
}
