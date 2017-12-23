package mkv

import "log"

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
		copy(b, p[2:])
		p = p[2:]
		a.SeqParam = append(a.SeqParam, b)
	}
	for i := byte(0); i < a.NPicParam; i++ {
		b := make([]byte, uint16(p[0])<<8|uint16(p[1]))
		copy(b, p[2:])
		p = p[2:]
		a.SeqParam = append(a.SeqParam, b)
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
