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
	"sync"

	_ "image/jpeg"

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

var QCbCr = []int{
	17, 18, 24, 47, 99, 99, 99, 99,
	18, 21, 26, 66, 99, 99, 99, 99,
	24, 26, 56, 99, 99, 99, 99, 99,
	47, 66, 99, 99, 99, 99, 99, 99,
	77, 99, 99, 99, 99, 99, 99, 99,
	92, 99, 99, 99, 99, 99, 99, 99,
	101, 99, 99, 99, 99, 99, 99, 99,
	99, 99, 99, 99, 99, 99, 99, 99,
}

var QQQ = []int{
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
}

var iQY = func() (q []int) {
	for _, v := range QY {
		q = append([]int{v}, q...)
	}
	return q
}()

var QY = []int{
	16, 11, 10, 16, 24, 40, 51, 61,
	12, 12, 14, 19, 26, 58, 60, 55,
	14, 13, 16, 24, 40, 57, 69, 56,
	14, 17, 22, 29, 51, 87, 80, 62,
	18, 22, 37, 56, 68, 109, 103, 77,
	24, 35, 55, 64, 81, 104, 113, 92,
	49, 64, 78, 87, 103, 121, 120, 101,
	72, 92, 95, 98, 112, 100, 103, 99,
}

var QY16 = []int{
	16, 11, 10, 16, 24, 40, 51, 61, 16, 11, 10, 16, 24, 40, 51, 61,
	16, 11, 10, 16, 24, 40, 51, 61, 16, 11, 10, 16, 24, 40, 51, 61,
	12, 12, 14, 19, 26, 58, 60, 55, 12, 12, 14, 19, 26, 58, 60, 55,
	12, 12, 14, 19, 26, 58, 60, 55, 12, 12, 14, 19, 26, 58, 60, 55,
	14, 13, 16, 24, 40, 57, 69, 56, 14, 13, 16, 24, 40, 57, 69, 56,
	14, 13, 16, 24, 40, 57, 69, 56, 14, 13, 16, 24, 40, 57, 69, 56,
	14, 17, 22, 29, 51, 87, 80, 62, 14, 17, 22, 29, 51, 87, 80, 62,
	14, 17, 22, 29, 51, 87, 80, 62, 14, 17, 22, 29, 51, 87, 80, 62,
	18, 22, 37, 56, 68, 109, 103, 77, 18, 22, 37, 56, 68, 109, 103, 77,
	18, 22, 37, 56, 68, 109, 103, 77, 18, 22, 37, 56, 68, 109, 103, 77,
	24, 35, 55, 64, 81, 104, 113, 92, 24, 35, 55, 64, 81, 104, 113, 92,
	24, 35, 55, 64, 81, 104, 113, 92, 24, 35, 55, 64, 81, 104, 113, 92,
	49, 64, 78, 87, 103, 121, 120, 101, 49, 64, 78, 87, 103, 121, 120, 101,
	49, 64, 78, 87, 103, 121, 120, 101, 49, 64, 78, 87, 103, 121, 120, 101,
	72, 92, 95, 98, 112, 100, 103, 99, 72, 92, 95, 98, 112, 100, 103, 99,
	72, 92, 95, 98, 112, 100, 103, 99, 72, 92, 95, 98, 112, 100, 103, 99,
}

var QY1 = []int{
	1, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

const N, M = 32, 32

var (
	oneSQtwo = 1 / math.Sqrt(2)
	cos2x1pi = func() []float64 {
		var costab []float64
		for i := 0; i < N; i++ {
			costab = append(costab, float64(2*i+1)*math.Pi)
		}
		return costab
	}()
)

func ntsc(v uint32) byte {
	v >>= 8
	return byte(clamp(v, 16, 235))
}
func clamp(v, l, h uint32) uint32 {
	if v > h {
		return h
	}
	if v < l {
		return l
	}
	return v
}

func main() {
	prant(a)
	prant(dct(a))
	prant(quant(dct(a), QY))
	prant(scale(quant(dct(a), QY), QY))
	prant(idct(scale(quant(dct(a), QY), QY)))

	x, _, _ := image.Decode(os.Stdin)
	img := image.NewRGBA(x.Bounds())
	draw.Draw(img, x.Bounds(), x, x.Bounds().Min, draw.Src)
	var wg sync.WaitGroup
	for dy := 0; dy+M < img.Bounds().Size().Y; dy += M {
		for dx := 0; dx+N < img.Bounds().Size().X; dx += N {
			dx, dy := dx, dy
			wg.Add(1)
			go func() {
				defer wg.Done()
				blockY := make([]int, 0, N*M)
				blockCb := make([]int, 0, N*M)
				blockCr := make([]int, 0, N*M)
				for y := 0; y < M; y++ {
					for x := 0; x < N; x++ {
						r, g, b, _ := img.At(x+dx, y+dy).RGBA()
						R, G, B := ntsc(r), ntsc(g), ntsc(b)
						Y, Cb, Cr := color.RGBToYCbCr(R, G, B)
						blockY = append(blockY, int(Y))
						blockCb = append(blockCb, int(Cb))
						blockCr = append(blockCr, int(Cr))
					}
				}
				blockY = idct(scale(quant(dct(blockY), QY), QY))
				//blockCr = idct(scale(quant(dct(blockCr), QCbCr), QCbCr))
				//blockCb = idct(scale(quant(dct(blockCb), QCbCr), QCbCr))

				sp := 0
				for y := 0; y < M; y++ {
					for x := 0; x < N; x++ {
						Y := byte(blockY[sp])
						cb := byte(blockCb[sp])
						cr := byte(blockCr[sp])
						img.Set(x+dx, y+dy, color.YCbCr{Y, cb, cr})
						sp++
					}
				}
			}()
		}
	}
	wg.Wait()
	pn := &png.Encoder{CompressionLevel: png.NoCompression}
	pn.Encode(os.Stdout, img)

}

func prant(A []int) {
	for i := range A {
		if i%N == 0 {
			log.Println()
		}
		log.Printf("%4d\t", A[i])
	}
	log.Println()
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func quant(DCT, QTAB []int) []int {
	DCT = append([]int{}, DCT...)
	for i := 0; i < min(len(DCT), len(QTAB)); i++ {
		DCT[i] = int(float64(DCT[i]) / float64(QTAB[i]))
	}
	return DCT
}
func scale(DCT, QTAB []int) []int {
	DCT = append([]int{}, DCT...)
	for i := 0; i < min(len(DCT), len(QTAB)); i++ {
		DCT[i] *= QTAB[i]
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
	q := 0
Loop:
	for x := 0; x < N; x++ {
		for y := 0; y < M; y++ {
			if q >= len(DCT) {
				break Loop
			}
			tmp := 0.0
			p := 0
			for i := 0; i < N; i++ {
				for j := 0; j < M; j++ {
					if p >= len(A) {
						continue
					}
					cosi := math.Cos(cos2x1pi[x] * float64(i) / (2.0 * N))
					cosj := math.Cos(cos2x1pi[y] * float64(j) / (2.0 * M))
					tmp += cosi * cosj * float64(A[p]) * (C(i) * C(j))
					p++
				}
			}
			tmp *= SQ2MN
			DCT[q] = int(tmp)
			q++
		}
	}
	for i := range DCT {
		DCT[i] += 128
	}
	return DCT
}

var SQ2MN = math.Sqrt(2.0/M) * math.Sqrt(2.0/N)

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

	q := 0
Loop:
	for i := 0; i < N; i++ {
		for j := 0; j < M; j++ {
			if q >= len(DCT) {
				break Loop
			}
			tmp := 0.0
			p := 0
			for x := 0; x < N; x++ {
				for y := 0; y < M; y++ {
					if p >= len(A) {
						continue
					}
					cosi := math.Cos(cos2x1pi[x] * float64(i) / (2.0 * N))
					cosj := math.Cos(cos2x1pi[y] * float64(j) / (2.0 * M))
					tmp += cosi * cosj * float64(A[p])
					p++
				}
			}
			tmp *= SQ2MN * (C(i) * C(j))
			DCT[q] = int(tmp)
			q++
		}
	}
	return append([]int{}, DCT...)
}
