package mkv

import (
	"fmt"
	"log"
)

type Slice struct {
	First bool // exp
	Kind  int  // exp
	// colour_plane_id (uint16)
	PPSID int // exp

	FrameNum int
	// field_pic_flag (byte)
	// bottom_field_flag (byte)
	IDRPicID int // IDR slices only
	// pic_order_cnt_lsb
	// delta_pic_order_cnt_bottom
	// delta_pic_order_cnt[ 0 ]
	// delta_pic_order_cnt[ 1 ]

	QPdx int // sexp
}

func sliceHeader(b *BitReader) *Slice {
	s := new(Slice)
	first, _ := ReadExp(b)
	s.First = first == 1
	s.Kind, _ = ReadExp(b)
	s.PPSID, _ = ReadExp(b)
	s.FrameNum = int(b.Advance(1))
	s.IDRPicID, _ = ReadExp(b)
	s.QPdx, _ = ReadExp(b)
	return s
}

func ReadNAL(lenlen int, br *BitReader) {
	i := 0
	for ; ; i++ {
		if i > 1024 {
			return
		}
		size := br.Advance(32) - 2
		//size = size<<24 | size>>24 | (size>>8)<<24 | (size<<8)>>24
		//size-=8
		log.Printf("found a nal of size %d\n", size)
		if br.Advance(1) == 1 {
			panic("forbidden")
		}
		cat, kind := br.Advance(2), br.Advance(5)
		log.Printf("	its cat is %d and type is %x\n", cat, kind)
		switch kind {
		case 5:
			log.Println("-[5]")
			// huuhhuuhuhuh
			log.Printf("hhurhuhu %#v\n", sliceHeader(br))
		default:
			for size >= 0 {
				fmt.Printf("%x", br.Advance(8))
				size--
			}
		}
		fmt.Println()

	}
	/*
		binary.Read(br, binary.LittleEndian, &v)
		data, err := ioutil.ReadAll(&io.LimitedReader{R: br, N: int64(v)})
		if err != nil {
			panic(err)
		}
		for len(data) > 0 {
			log.Printf("\t%x\n", data)
			br2 := NewBitReader(bytes.NewReader(data))

			switch k := br2.Advance(8); k {
			case 6:
				log.Println("-[6]")
				if br2.Advance(8) != 5 {
					panic("crappy sei message")
				}
				X := int64(2)

				x := int64(0)
				n := int64(0xff)
				for n == 0xff {
					n = br2.Advance(8)
					x += n
					X++
				}
				log.Println("computed length of sei is", x)
				x += X + 5
				data = data[x:]
				log.Printf("\t%x\n", data)
			case 5, 101:
				log.Println("-[5]")
				// huuhhuuhuhuh
				log.Printf("hhurhuhu %#v\n", sliceHeader(br2))
			default:
				log.Println("dont know what %x is", k)
				log.Fatalf("\t%x\n", data)
			}
		}
	*/
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
