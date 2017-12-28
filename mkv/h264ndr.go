package mkv

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
)

type NKind byte

const (
	NUnspec NKind = iota
	NCodedNonIDR
	NCodedLayerA
	NCodedLayerB
	NCodedLayerC
	NCodedIDR
	NSupplemental
	NSeqParam
	NPicParam
	NAccessDelim
	NSeqEnd
	NStreamEnd
	NFiller
	NSeqParamExt
	NPrefixNAL
	NSeqParamSubset
	_
	_
	_
	NCodedAux
	NCodedExt
)

type SPS struct {
	Forbidden   byte // 1
	NalRefIDC   byte // 2
	NalUnitType byte // 5

	Profile byte // 8

	Constraints byte // first 4 bits

	LevelIDC byte

	ID          int
	MaxFrames   int // originally encoded in log 2 - 4
	PicOrder    int
	NReferences int

	GapsAllowed bool // one bit

	MBWidth  int // -1
	MBHeight int // -1

	MBNoField          byte // 1
	DirectInference8x8 byte //1
	Cropping           byte //1
	HasVUI             byte //1
	Stop               byte //1
}

func ReadExp(b *BitReader) (int, error) {
	N := b.AdvanceZero()
	log.Println("advanced", N, "zeroes")
	if N >= 32 {
		panic("hachi machi")
		return 0, fmt.Errorf("variable length exp golomb code overflow (%d >=32 bits)", N)
	}
	return int(1<<uint(N) | b.Advance(N+1)), nil
}

func ReadEBLM(b *BitReader) (int, error) {
	N := b.AdvanceZero()
	log.Println("advanced", N, "zeroes")
	if N >= 32 {
		panic("hachi machi")
		return 0, fmt.Errorf("variable length exp golomb code overflow (%d >=32 bits)", N)
	}
	N++
	return int(b.Advance(N * 7)), nil
}

var (
	ErrForbiddenBitSet = errors.New("forbidden bit set to 1")
)

func ReadSPS(br *BitReader) (*SPS, error) {
	sps := new(SPS)
	if br.Advance(1) != 0 {
		panic("ErrForbiddenBitSet")
		return nil, ErrForbiddenBitSet
	}
	sps.NalRefIDC = byte(br.Advance(2))
	sps.NalUnitType = byte(br.Advance(5))
	sps.Profile = byte(br.Advance(8))
	sps.Constraints = byte(br.Advance(4))
	br.Advance(4)
	sps.LevelIDC = byte(br.Advance(8))
	log.Printf("the SPS is %#v\n", sps)

	br.Advance(8)
	sps.ID, _ = ReadExp(br)
	x, _ := ReadExp(br)
	log.Printf("x SPS is %#v\n", sps)
	sps.MaxFrames = (1 << uint(x)) + 4
	sps.PicOrder, _ = ReadExp(br)
	sps.NReferences, _ = ReadExp(br)

	if br.Advance(1) == 1 {
		sps.GapsAllowed = true
	}

	sps.MBWidth, _ = ReadExp(br)
	sps.MBHeight, _ = ReadExp(br)
	sps.MBWidth++
	sps.MBHeight++

	sps.MBNoField = byte(br.Advance(1))
	sps.DirectInference8x8 = byte(br.Advance(1))
	sps.Cropping = byte(br.Advance(1))
	sps.HasVUI = byte(br.Advance(1))
	if byte(br.Advance(1)) != 1 {
		panic("bad stop bit")
		return nil, fmt.Errorf("stop bit set to 0")
	}
	return sps, nil
}

type AACPrivate struct {
	_         byte
	Profile   byte
	_         byte
	Level     byte
	NALUSize  byte // only the first two bits on the right
	NSeqParam byte // only the first 5 bits
	SeqParam  [][]byte
	NPicParam byte
	PicParam  [][]byte
}

func (a *AACPrivate) MarshalBinary(p []byte) error {
	a.Profile = p[1]
	a.Level = p[3]
	a.NALUSize = p[4]&3 + 1
	a.NSeqParam = p[5] & 0x1f
	p = p[6:]
	for i := byte(0); i < a.NSeqParam; i++ {
		b := make([]byte, uint16(p[0])<<8|uint16(p[1]))
		p = p[2:]
		n := copy(b, p)
		p = p[n:]
		a.SeqParam = append(a.SeqParam, b)
	}
	if len(a.SeqParam) > 0 {
		sps, err := ReadSPS(&BitReader{r: bufio.NewReader(bytes.NewReader(a.SeqParam[0]))})
		if err != nil {
			return err
		}
		log.Printf("SPS-Decoded: %#v\n", sps)
	}
	if len(p) > 3 {
		a.NPicParam = p[0]
		log.Println(a.NPicParam)
		p = p[1:]
		for i := byte(0); i < a.NPicParam; i++ {
			b := make([]byte, uint16(p[0])<<8|uint16(p[1]))
			p = p[2:]
			n := copy(b, p)
			p = p[n:]
			a.PicParam = append(a.PicParam, b)
		}
	}
	return nil
}

type NAL struct {
	Forbid bool
	Ref    byte
	Kind   NKind
}

func parseNAL(b byte) (NAL, error) {
	log.Printf("input is %x %08b\n", b, b)
	var n NAL
	n.Forbid = b&0x80 != 0
	log.Printf("ref is %x %08b\n", b>>6, b>>6)
	n.Ref = byte(b >> 5)
	log.Printf("Kind is 0x%x %08b\n", (0x1f & b), (0x1f & b))
	n.Kind = NKind(0x1f & b)
	return n, nil
}
