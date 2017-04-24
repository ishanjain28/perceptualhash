package main

import (
	"image"
	"fmt"
	"os"
	"github.com/nfnt/resize"
	"image/png"
)

func main() {
	file, _ := os.Open("image.png")

	grayImage := Grayscale(file, "decomposition.min")

	data, _ := os.Create("ishan.resized.png")
	//normal, _ := os.Create("ishan.png")

	resizedGrayImage := resize.Resize(9, 9, grayImage, resize.Bicubic)

	png.Encode(data, resizedGrayImage)
	//png.Encode(normal, grayImage)
	//rowHash := rowHash(resizedGrayImage)
	//columnHash := columnHash(resizedGrayImage)
	//completeHash := [][]int{}
	//
	//for _, v := range rowHash {
	//	fmt.Printf("%d ", v)
	//}
	//fmt.Printf("\n\n")
	//for _, v := range columnHash {
	//	fmt.Printf("%d ", v)
	//}
	calcHash(resizedGrayImage)
	//fmt.Println(len(rowHash), len(columnHash))
}

func rowHash(img image.Image) []int {

	x := img.Bounds().Max.X
	y := img.Bounds().Max.Y

	xHash := make([]int, (x-1)*(y-1))

	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {

			currentGrayValue, _, _, _ := img.At(i, j).RGBA()
			previousGrayValue, _, _, _ := img.At(i, j-1).RGBA()

			if (currentGrayValue >= previousGrayValue) {
				xHash[i*x+j] = 1
			} else {
				xHash[i*x+j] = 0
			}

		}
	}
	return xHash
}

func columnHash(img image.Image) []int {

	x := img.Bounds().Max.X
	y := img.Bounds().Max.Y

	fmt.Println(x, y)

	hash := make([]int, x*y)
	count := 0
	for i := 1; i < y; i++ {
		for j := 1; j < x; j++ {

			currentGrayValue, _, _, _ := img.At(j-1, i).RGBA()
			previousGrayValue, _, _, _ := img.At(j-1, i).RGBA()
			count++
			if (currentGrayValue >= previousGrayValue) {
				hash[j*(y-1)+i] = 1
			} else {
				hash[j*(y-1)+i] = 0
			}
		}
	}
	fmt.Println(count)
	return hash
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
			previousRowGray, _, _, _ := img.At(i, j-1).RGBA()    //(i + x*j) - 1
			previousColumnGray, _, _, _ := img.At(j, i-1).RGBA() //(j + y*i) - 1

			if currentGray >= previousRowGray {
				rowHash = append(rowHash, 1)
			}
			if currentGray < previousRowGray {
				rowHash = append(rowHash, 1)
			}

			if currentGray >= previousColumnGray {
				columnHash = append(columnHash, 1)
			}
			if currentGray < previousColumnGray {
				columnHash = append(columnHash, 0)
			}
		}
	}
	return rowHash, columnHash
}
