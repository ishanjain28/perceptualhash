package main

import (
	"os"
	"log"
	"image/png"
	"fmt"
	"image"
	"image/color"
	"math"
	"sync"
	"path"
	"strings"
	"image/jpeg"
)

var wg sync.WaitGroup

var filenames = []string{
	"basic.improved.for.human.eye",
	"basic",
	"desaturated",
	"decomposition.max",
	"decomposition.min",
	"single.channel.red",
	"single.channel.green",
	"single.channel.blue",
}

func main() {
	fmt.Printf("Enter File Path/Name: ")
	srcPath := ""
	fmt.Scanf("%s", &srcPath)
	img, err := os.Open(srcPath)
	defer img.Close()

	if err != nil {
		log.Fatalf("%s", err)
	}
	var pngimg image.Image
	fileExtension := path.Ext(srcPath)
	switch fileExtension {
	case ".png":
		pngimg, err = png.Decode(img)
	case ".jpg":
		pngimg, err = jpeg.Decode(img)
	case ".jpeg":
		pngimg, err = jpeg.Decode(img)
	}

	w, h := pngimg.Bounds().Max.X, pngimg.Bounds().Max.Y
	fmt.Printf("Image Resolution is %dx%d\n", w, h)

	for index, name := range filenames {

		//if index != 2 && index != 3 {
		wg.Add(1)
		go CreateGrayImage(pngimg, name, path.Base(srcPath), fileExtension, w, h, index)
		//}
	}
	wg.Wait()
}

func CreateGrayImage(pngimg image.Image, name, srcFilename, fileExtension string, w, h int, index int) {
	grayImage := image.NewGray16(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})
	for x := 0; x <= w; x++ {
		for y := 0; y <= h; y++ {
			point := pngimg.At(x, y)
			r, g, b, _ := point.RGBA()
			var avg float64
			switch index {
			case 0:
				avg = HumanEyesImprovedGrayscaling(r, g, b)
			case 1:
				avg = SimpleGrayScaling(r, g, b)
			case 2:
				avg = Desaturation(r, g, b)
			case 3:
				avg = float64(Max(Max(r, g), b))
			case 4:
				avg = float64(Min(Min(r, g), b))
			case 5:
				avg = float64(r)
			case 6:
				avg = float64(g)
			case 7:
				avg = float64(b)
			}
			grayColor := color.Gray16{uint16(math.Ceil(avg))}
			grayImage.Set(x, y, grayColor)
		}
	}

	srcFilename = strings.Replace(srcFilename, fileExtension, "", -1)

	os.Mkdir("grayscaled", 0777)

	os.Mkdir("grayscaled/"+srcFilename, 0777)

	outfile, err := os.Create("./grayscaled/" + srcFilename + "/" + name + ".png")
	if err != nil {
		log.Fatalln(err)
	}
	defer outfile.Close()
	png.Encode(outfile, grayImage)

	wg.Done()
}

func SimpleGrayScaling(r, g, b uint32) float64 {
	return float64((r + g + b) / 3)
}

func HumanEyesImprovedGrayscaling(r, g, b uint32) float64 {
	return float64(0.3)*float64(r) + float64(0.59)*float64(g) + float64(0.11)*float64(b)
}

func Desaturation(r, g, b uint32) float64 {
	return float64(MaxOfThree(r, g, b)+MinOfThree(r, g, b)) / 2
}

func MaxOfThree(r, g, b uint32) uint32 {
	return Max(Max(r, g), b)
}

func MinOfThree(r, g, b uint32) uint32 {
	return Min(Min(r, g), b)
}

func Max(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

func Min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}
