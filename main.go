package main

import (
	"image"
	"os"
	"github.com/nfnt/resize"
	"image/png"
	"log"
	"fmt"
)

func main() {
	file, err := os.Open("image.png")
	if err != nil {
		log.Fatalf("%s", err)
	}
	grayImage := Grayscale(file, "decomposition.min")

	data, _ := os.Create("ishan.resized.png")
	//normal, _ := os.Create("ishan.png")

	resizedGrayImage := resize.Resize(9, 9, grayImage, resize.Bicubic)

	png.Encode(data, resizedGrayImage)

	rowHash, columnHash := calcHash(resizedGrayImage)

	for i, v := range rowHash {
		if i%8 == 0 {
			fmt.Printf("\n")
		}
		fmt.Printf("%d ", v)
	}
	fmt.Printf("\n\n")
	for i, v := range columnHash {
		if i%8 == 0 {
			fmt.Printf("\n")
		}
		fmt.Printf("%d ", v)
	}
}

func calcHash(img image.Image) ([]int, []int) {
	x := img.Bounds().Max.X
	y := img.Bounds().Max.Y
	X := x - 1
	Y := y - 1
	rowHash := []int{}
	columnHash := []int{}
	for i := 0; i < X; i++ {
		for j := 0; j < Y; j++ {

			currentGray, _, _, _ := img.At(i, j).RGBA()
			nextRowGray, _, _, _ := img.At(i, j+1).RGBA()    //(i + x*j) - 1
			nextColumnGray, _, _, _ := img.At(j, i+1).RGBA() //(j + y*i) - 1

			if currentGray >= nextRowGray {
				rowHash = append(rowHash, 1)
			}
			if currentGray < nextRowGray {
				rowHash = append(rowHash, 0)
			}

			if currentGray >= nextColumnGray {
				columnHash = append(columnHash, 1)
			}
			if currentGray < nextColumnGray {
				columnHash = append(columnHash, 0)
			}
		}
	}
	return rowHash, columnHash
}
