// Just testing. Do not use.
package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"

	_ "golang.org/x/image/bmp"
)

var a = []int{52, 55, 61, 66, 70, 61, 64, 73,
	63, 59, 55, 90, 109, 85, 69, 72,
	62, 59, 68, 113, 144, 104, 66, 73,
	63, 58, 71, 122, 154, 106, 70, 69,
	67, 61, 68, 104, 126, 88, 68, 70,
	79, 65, 60, 70, 77, 68, 58, 75,
	85, 71, 64, 59, 55, 61, 65, 83,
	87, 79, 69, 68, 65, 76, 78, 94,
}

var QCrCB = []int{
	17, 18, 24, 47, 99, 99, 99, 99,
	18, 21, 26, 66, 99, 99, 99, 99,
	24, 26, 56, 99, 99, 99, 99, 99,
	47, 66, 99, 99, 99, 99, 99, 99,
	77, 99, 99, 99, 99, 99, 99, 99, 99,
	92, 99, 99, 99, 99, 99, 99, 99, 99,
	101, 99, 99, 99, 99, 99, 99, 99, 99,
	99, 99, 99, 99, 99, 99, 99, 99, 99,
}

var QY = []int{
	16, 11, 10, 16, 24, 40, 51, 61,
	12, 12, 14, 19, 26, 58, 60, 55,
	14, 13, 16, 24, 40, 57, 69, 56,
	14, 17, 22, 29, 51, 87, 80, 62,
	18, 22, 37, 56, 68, 109, 103,
	24, 35, 55, 64, 81, 104, 113,
	49, 64, 78, 87, 103, 121, 120,
	72, 92, 95, 98, 112, 100, 103,
}

var (
	N        = 8
	oneSQtwo = 1 / math.Sqrt(2)
	cos2x1pi = func() []float64 {
		var costab []float64
		for i := 0; i < N; i++ {
			costab = append(costab, float64(2*i+1)*math.Pi)
		}
		return costab
	}()
)

func main() {
	prant(a)
	prant(dct(a))
	prant(quant(dct(a), QY))
	prant(scale(quant(dct(a), QY), QY))
	prant(idct(scale(quant(dct(a), QY), QY)))

	img, _, _ := image.Decode(os.Stdin)
	var block []int
	for dy := 0; dy+8 < img.Bounds().Size().Y; dy += 8 {
	for dx := 0; dx+8 < img.Bounds().Size().X; dx += 8 {
		block=block[:0]
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				r, g, b, _ := img.At(x+dx, y+dy).RGBA()
				Y, _, _ := color.RGBToYCbCr(byte(r/256), byte(g/256), byte(b/256))
				block = append(block, int(Y))
			}
		}
		compressed := idct(scale(quant(dct(a), QY), QY))
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				r, g, b, _ := img.At(x+dx, y+dy).RGBA()
				Y, cb, cr := color.RGBToYCbCr(byte(r/256), byte(g/256), byte(b/256))
				Y = byte(compressed[x+8*y])
				img.(draw.Image).Set(x+dx, y+dy, color.YCbCr{Y, cb, cr})
			}
		}
	}}
	pn := &png.Encoder{CompressionLevel: png.NoCompression}
	pn.Encode(os.Stdout, img)

}

func prant(A []int) {
	for i := range A {
		if i%8 == 0 {
			log.Println()
		}
		log.Printf("%4d\t", A[i])
	}
	log.Println()
}
func quant(DCT, QTAB []int) []int {
	DCT = append([]int{}, DCT...)
	for i, q := range QTAB {
		DCT[i] = int(float64(DCT[i]) / float64(q))
	}
	return DCT
}
func scale(DCT, QTAB []int) []int {
	DCT = append([]int{}, DCT...)
	for i, q := range QTAB {
		DCT[i] = DCT[i] * q
	}
	return DCT
}

func idct(A []int) []int {
	A = append([]int{}, A...)
	C := func(u int) float64 {
		if u == 0 {
			return oneSQtwo
		}
		if u > 0 {
			return 1
		}
		return 0
	}
	DCT := make([]int, len(A))
	for x := 0; x < N; x++ {
		for y := 0; y < N; y++ {
			tmp := 0.0
			for i := 0; i < N; i++ {
				for j := 0; j < N; j++ {
					cosi := math.Cos(cos2x1pi[x] * float64(i) / 16.0)
					cosj := math.Cos(cos2x1pi[y] * float64(j) / 16.0)
					tmp += cosi * cosj * float64(A[i+j*8]) * (C(i) * C(j))
				}
			}
			tmp *= 1 / math.Sqrt(16)
			DCT[x+y*8] = int(tmp)
		}
	}
	for i := range DCT {
		DCT[i] += 128
	}
	return DCT
}

func dct(A []int) []int {
	A = append([]int{}, A...)
	C := func(u int) float64 {
		if u == 0 {
			return oneSQtwo
		}
		if u > 0 {
			return 1
		}
		return 0
	}
	for i := range A {
		A[i] -= 128
	}
	DCT := make([]int, len(A))
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			tmp := 0.0
			for x := 0; x < N; x++ {
				for y := 0; y < N; y++ {
					cosi := math.Cos(cos2x1pi[x] * float64(i) / 16.0)
					cosj := math.Cos(cos2x1pi[y] * float64(j) / 16.0)
					tmp += cosi * cosj * float64(A[x+y*8])
				}
			}
			tmp *= 1 / math.Sqrt(16) * (C(i) * C(j))
			DCT[i+j*8] = int(tmp)
		}
	}
	return append([]int{}, DCT...)
}
