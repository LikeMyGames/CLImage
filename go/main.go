package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"golang.org/x/term"
)

func main() {
	imagePath := ""
	if len(os.Args) >= 2 {
		imagePath = os.Args[1]
	} else {
		imagePath = "image.png"
	}
	termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}
	// fmt.Printf("Width: %v, Height: %v\n", termWidth, termHeight)

	file, err := os.Open(imagePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var img image.Image

	if strings.HasSuffix(imagePath, ".png") {
		img, err = png.Decode(file)
		if err != nil {
			panic(err)
		}
	} else if strings.HasSuffix(imagePath, ".jpeg") || strings.HasSuffix(imagePath, ".jpg") {
		img, err = jpeg.Decode(file)
		if err != nil {
			panic(err)
		}
	}

	// fmt.Println(img.Bounds().Dx())

	scalFact := float64(termWidth) / float64(img.Bounds().Dx())
	// fmt.Println("s", scalFact)

	imgScaled := image.NewRGBA64(image.Rect(0, 0, termWidth, int(float64(img.Bounds().Dy())*scalFact/2)))
	imgScaledWidth := imgScaled.Bounds().Dx()
	imgScaledHeight := imgScaled.Bounds().Dy()

	widthAdjust := img.Bounds().Dx() / termWidth
	heightAdjust := img.Bounds().Dx() / (termWidth) * 2
	// fmt.Println(widthAdjust, "\t", uint32(widthAdjust))

	for x := range imgScaledWidth {
		for y := range imgScaledHeight {
			imgScaled.Set(x, y, img.At(x*widthAdjust, y*heightAdjust))
			// var (
			// 	R int64 = 0
			// 	G int64 = 0
			// 	B int64 = 0
			// 	A int64 = 0
			// )
			// for i := 0; i < int(widthAdjust); i++ {
			// 	for j := 0; j < int(heightAdjust); j++ {
			// 		rgba := img.At(x+i, y+j)
			// 		r, g, b, a := rgba.RGBA()
			// 		R += int64(mapRange(float64(r), 0, 65535, 0, 255))
			// 		G += int64(mapRange(float64(g), 0, 65535, 0, 255))
			// 		B += int64(mapRange(float64(b), 0, 65535, 0, 255))
			// 		A += int64(mapRange(float64(a), 0, 65535, 0, 255))
			// 	}
			// }
			// R = int64(float64(R) / float64(widthAdjust*heightAdjust))
			// G = int64(float64(G) / float64(widthAdjust*heightAdjust))
			// B = int64(float64(B) / float64(widthAdjust*heightAdjust))
			// A = int64(float64(A) / float64(widthAdjust*heightAdjust))
			// imgScaled.Set(x, y, color.RGBA{R: uint8(R), G: uint8(G), B: uint8(B), A: uint8(A)})
		}
	}
	// for x := range imgScaledWidth {
	// 	for y := range imgScaledHeight {
	// 		fmt.Println("X:", x, "\tY:", y, "\t", imgScaled.At(x, y))
	// 	}
	// }

	outFile, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}

	png.Encode(outFile, imgScaled)

	for y := range imgScaledHeight {
		for x := range imgScaledWidth {
			r, g, b, a := imgScaled.At(x, y).RGBA()
			r = uint32(Int16toInt8(uint16(r)))
			g = uint32(Int16toInt8(uint16(g)))
			b = uint32(Int16toInt8(uint16(b)))
			a = uint32(Int16toInt8(uint16(a)))
			if a == 0 {
				fmt.Print(" ")
				continue
			}
			SetTermColor(uint8(r), uint8(g), uint8(b))
			fmt.Print("\u2588")
		}
		fmt.Print("\n")
	}

	ResetTermColor()
}

func mapRange(in, in_min, in_max, out_min, out_max float64) float64 {
	return (out_min + (((in - in_min) / (in_max - in_min)) * (out_max - out_min)))
}

func Int16toInt8(num uint16) uint8 {
	return uint8(mapRange(float64(num), 0, 65535, 0, 255))
}

func SetTermColor(r, g, b uint8) {
	fmt.Printf("\x1b[38;2;%d;%d;%dm", r, g, b)
}

func ResetTermColor() {
	fmt.Printf("\x1b[0;0m")
}

// func RGBtoHSV(rgba color.RGBA) (H, S, V, A uint8) {
// 	r, g, b, a := rgba.RGBA()
// 	cmax := math.Max(float64(r), math.Max(float64(g), float64(b)))
// 	cmin := math.Min(float64(r), math.Min(float64(g), float64(b)))
// 	diff := cmax - cmin
// 	if diff == 0 {
// 		H = 0
// 	}
// 	if cmax == r {
// 		H = uint8(60*(float64(g-b)/(diff+6))) % 6
// 	}
// }
