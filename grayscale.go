package main

import (
	"os"
	"image/png"
	"fmt"
	"image"
	"image/color"
	"math"
	"log"
)

var grayscalingOptions = []string{
	"basic.improved",
	"basic",
	"desaturated",
	"decomposition.max",
	"decomposition.min",
	"single.channel.red",
	"single.channel.green",
	"single.channel.blue",
}

func Grayscale(img *os.File, GrayscalingMethod string) *image.Gray16 {
	pngimg, err := png.Decode(img)
	if err != nil {
		fmt.Printf("%s", err)
	}
	// w and h are number of rows and columns of pixels in an image
	w, h := pngimg.Bounds().Max.X, pngimg.Bounds().Max.Y

	return CreateGrayImage(pngimg, GrayscalingMethod, w, h)
}

func CreateGrayImage(pngimg image.Image, Method string, w, h int) *image.Gray16 {
	grayImage := image.NewGray16(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})
	for x := 0; x <= w; x++ {
		for y := 0; y <= h; y++ {
			point := pngimg.At(x, y)
			r, g, b, _ := point.RGBA()
			var avg float64
			switch Method {
			case "basic.improved":
				avg = humanEyesImprovedGrayscaling(r, g, b)
			case "basic":
				avg = simpleGrayScaling(r, g, b)
			case "desaturated":
				avg = desaturation(r, g, b)
			case "decomposition.max":
				avg = float64(Max(Max(r, g), b))
			case "decomposition.min":
				avg = float64(Min(Min(r, g), b))
			case "single.channel.red":
				avg = float64(r)
			case "single.channel.green":
				avg = float64(g)
			case "single.channel.blue":
				avg = float64(b)
			default:
				log.Fatalln("Please Provide a valid Option")
			}
			grayColor := color.Gray16{uint16(math.Ceil(avg))}
			grayImage.Set(x, y, grayColor)
		}
	}

	return grayImage
}

func simpleGrayScaling(r, g, b uint32) float64 {
	return float64((r + g + b) / 3)
}

func humanEyesImprovedGrayscaling(r, g, b uint32) float64 {
	return float64(0.3)*float64(r) + float64(0.59)*float64(g) + float64(0.11)*float64(b)
}

func desaturation(r, g, b uint32) float64 {
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
