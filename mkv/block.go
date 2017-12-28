package mkv

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"time"
)

func (b Block) String() string {
	return fmt.Sprintf("%#v %s\n", b, b.Flag.String())
}

type Flag byte

func (f Flag) String() string {
	s := "flags( "
	if f&FKeyFrame != 0 {
		s += "IDR "
	}
	if f&FNoDuration != 0 {
		s += "Invisible "
	}
	if f&FLaceEMBL != 0 {
		s += "EMBL "
	} else if f&FLaceFixed != 0 {
		s += "FIXED "
	} else if f&FLaceXIPH != 0 {
		s += "XIPH "
	} else {
		s += "NOLACE "
	}

	if f&FDiscardable != 0 {
		s += "Discardable "
	}
	return s + ")"
}

const (
	FKeyFrame    Flag = 0x80
	FNoDuration  Flag = 0x08
	FLaceEMBL         = FLaceXIPH | FLaceFixed
	FLaceXIPH    Flag = 0x02
	FLaceFixed   Flag = 0x04
	FDiscardable Flag = 0x01
)

type Block struct {
	Track    int
	TimeCode time.Duration
	Flag
	NFrames int

	// wow, now we're getting to the good stuff
	// ... maybe later
}

func ReadBlock(r io.Reader) *Block {
	b := NewBitReader(r)
	track, _ := ReadEBLM(b)
	tc := b.Advance(16)
	flag := Flag(b.Advance(8))
	nframe := 0
	if flag&FLaceEMBL != 0 {
		nframe = int(b.Advance(8))
	}
	block := &Block{
		Track:    track,
		TimeCode: time.Duration(tc),
		Flag:     flag,
		NFrames:  nframe,
	}
	log.Printf("pending block: %+v\n", block)
	ReadNAL(4, b)
	return block
}

func NewBitReader(r io.Reader) *BitReader {
	switch r := r.(type) {
	case *bufio.Reader:
		return &BitReader{r: r}
	}
	return &BitReader{r: bufio.NewReader(r)}
}

type BitReader struct {
	tmp uint64
	c   byte
	n   uint
	r   *bufio.Reader
}

func (b *BitReader) AdvanceZero() (nz int) {
	for b.Advance(1) == 0 {
		nz++
	}
	log.Printf("t@here are %d leading zeroes\n", nz)
	return nz
}

func (b *BitReader) Align() (state int64) {
	defer func() {
		if b.n%8 != 0 {
			panic("failed")
		}
	}()
	if b.n%8 == 0 {
		return
	}
	return b.Advance(int(8 - b.n%8))
}

func (b *BitReader) Advance(bits int) (state int64) {
	mask := ^((^uint64(0)) << uint(bits))
	//	for bits >= 8 {
	//		b.c, _ = b.r.ReadByte()
	//		b.tmp = b.tmp<<8 | uint64(b.c)
	//		b.n += 8
	//		bits -= 8
	//	}
	for bits > 0 {
		if b.n%8 == 0 {
			b.c, _ = b.r.ReadByte()
		}
		b.tmp = b.tmp<<1 | uint64(b.c>>7)
		b.c <<= 1
		b.n++
		bits--
	}
	return int64(b.tmp & mask)
}
func (b *BitReader) Value() uint64 {
	return b.tmp
}

func (b *BitReader) Next() (uint64, error) {
	var err error
	if b.n%8 == 0 {
		b.c, err = b.r.ReadByte()
	}
	b.tmp = b.tmp<<1 | uint64(b.c>>7)
	b.c <<= 1
	b.n++
	return b.tmp, err
}
