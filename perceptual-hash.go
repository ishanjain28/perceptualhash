package perceptualhash

import (
	"image/jpeg"
	"image"
	"image/png"
	"os"
	"image/color"
	"math"
	"github.com/nfnt/resize"
	"strconv"
	"strings"
	"path"
)

func decodeImage(img *os.File, imageExt string) (image.Image, error) {

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

func convertToGrayscale(img image.Image) *image.Gray16 {
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

func downsizeImage(img image.Image, hashLength int) image.Image {
	widthAndLength := uint(math.Ceil(math.Sqrt(float64(hashLength)/2.0)) + 1)
	return resize.Resize(widthAndLength, widthAndLength, img, resize.Bicubic)
}

func calculateRowHash(img image.Image) string {

	imageX, imageY := img.Bounds().Max.X, img.Bounds().Max.Y
	completeHash := [][]string{}

	for j := 0; j < imageY-1; j++ {
		rowHash := []string{}
		for i := 0; i < imageX-1; i++ {
			r0, _, _, _ := img.At(i, j).RGBA()
			r1, _, _, _ := img.At(i+1, j).RGBA()
			if r1 >= r0 {
				rowHash = append(rowHash, strconv.Itoa(1))
			}
			if r1 < r0 {
				rowHash = append(rowHash, strconv.Itoa(0))
			}
		}
		completeHash = append(completeHash, rowHash)
	}
	completeHashString := []string{}

	for i := 0; i < len(completeHash); i++ {
		rowHashString := strings.Join(completeHash[i], "")
		completeHashString = append(completeHashString, rowHashString)
	}

	return strings.Join(completeHashString, "")
}

func calculateColumnHash(img image.Image) string {
	imageX, imageY := img.Bounds().Max.X, img.Bounds().Max.Y
	completeHash := [][]string{}

	for i := 0; i < imageX-1; i++ {
		columnHash := []string{}
		for j := 0; j < imageY-1; j++ {
			r0, _, _, _ := img.At(i, j).RGBA()
			r1, _, _, _ := img.At(i, j+1).RGBA()
			if r1 >= r0 {
				columnHash = append(columnHash, strconv.Itoa(1))
			}
			if r1 < r0 {
				columnHash = append(columnHash, strconv.Itoa(0))
			}
		}
		completeHash = append(completeHash, columnHash)
	}

	completeHashString := []string{}
	for i := 0; i < len(completeHash); i++ {
		columnHashString := strings.Join(completeHash[i], "")
		completeHashString = append(completeHashString, columnHashString)
	}

	return strings.Join(completeHashString, "")
}

func CalculateHash(img image.Image) string {
	rowHash := calculateRowHash(img)
	columnHash := calculateColumnHash(img)
	hash := ""
	hash += rowHash
	hash += columnHash

	return hash
}

func diffHashes(hash1, hash2 string, hashLength float32) float32 {
	var similarityInTwoHashes float32
	for i := 0; i < len(hash1); i++ {
		if (byte(hash1[i]) ^ byte(hash2[i])) == 0 {
			similarityInTwoHashes++
		}
	}

	return (similarityInTwoHashes / hashLength) * 100
}

func loadFiles(firstImagePath, secondImagePath string) (*os.File, *os.File, error) {
	firstImage, err := os.Open(firstImagePath)
	if err != nil {
		return nil, nil, err
	}
	secondImage, err := os.Open(secondImagePath)
	if err != nil {
		return nil, nil, err
	}
	return firstImage, secondImage, nil
}

func CalculateSimilarity(firstImagePath, secondImagePath string, hashLength int) (float32, error) {
	firstImage, secondImage, err := loadFiles(firstImagePath, secondImagePath)
	if err != nil {
		return 0.0, err
	}
	defer firstImage.Close()
	defer secondImage.Close()
	firstDecodedImage, err := decodeImage(firstImage, path.Ext(firstImagePath))
	if err != nil {
		return 0.0, err
	}
	secondDecodedImage, err := decodeImage(secondImage, path.Ext(secondImagePath))
	if err != nil {
		return 0.0, err
	}
	firstGrayImage, secondGrayImage := convertToGrayscale(firstDecodedImage), convertToGrayscale(secondDecodedImage)

	firstDownsizedImage, secondDownsizedImage := downsizeImage(firstGrayImage, hashLength), downsizeImage(secondGrayImage, hashLength)

	return diffHashes(CalculateHash(firstDownsizedImage), CalculateHash(secondDownsizedImage), float32(hashLength)), nil

}
