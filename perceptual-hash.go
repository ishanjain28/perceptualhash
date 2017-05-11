package perceptual_hash

import (
	"image/jpeg"
	"image"
	"image/png"
	"os"
	"image/color"
	"math"
	"github.com/nfnt/resize"
)

func DecodeImage(img *os.File, imageExt string) (image.Image, error) {

	var decodedImage image.Image
	var err error
	switch imageExt {
	case ".jpg", ".JPG", ".JPEG", ".jpeg":
		decodedImage, err = jpeg.Decode(img)
	case ".png", ".PNG":
		decodedImage, err = png.Decode(img)
	}
	if err != nil {
		return nil, err
	}
	return decodedImage, nil
}

func ConvertToGrayscale(img image.Image) *image.Gray16 {
	imageX, imageY := img.Bounds().Max.X, img.Bounds().Max.Y
	grayImage := image.NewGray16(image.Rectangle{image.Point{0, 0}, image.Point{imageX, imageY}})

	for i := 0; i < imageX; i++ {
		for j := 0; j < imageY; j++ {
			r, g, b, _ := img.At(i, j).RGBA()
			pixelAvg := math.Ceil((float64(r) + float64(g) + float64(b)) / 3.0)
			grayColor := color.Gray16{uint16(pixelAvg)}

			grayImage.SetGray16(i, j, grayColor)
		}
	}

	return grayImage
}

func DownsizeImage(img image.Image) image.Image {
	return resize.Resize(4, 4, img, resize.Bicubic)
}

func CalculateRowHash(img image.Image) []int {

	imageX, imageY := img.Bounds().Max.X, img.Bounds().Max.Y
	hash := []int{}
	for j := 0; j < imageY; j++ {
		for i := 0; i < imageX; i++ {
			r0, _, _, _ := img.At(j, i).RGBA()
			r1, _, _, _ := img.At(j+1, i).RGBA()

			if r0 >= r1 {
				hash = append(hash, 1)
			} else {
				hash = append(hash, 0)
			}
		}
	}
	return hash
}

func CalculateColumnHash(img image.Image) []int {
	imageX, imageY := img.Bounds().Max.X, img.Bounds().Max.Y
	hash := []int{}

	for i := 0; i < imageX; i++ {
		for j := 0; j < imageY; j++ {
			r0, _, _, _ := img.At(j, i).RGBA()
			r1, _, _, _ := img.At(j, i+1).RGBA()

			if r0 > r1 {
				hash = append(hash, 1)
			} else {
				hash = append(hash, 0)
			}
		}
	}

	return hash
}
