package main

import (
	"fmt"
	"bufio"
	"os"
	"log"
	"github.com/ishanjain28/phash"
)

func main() {
	firstImagePath := ""
	secondImagePath := ""

	fmt.Printf("Enter First Image Path: ")
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		s := scanner.Text()
		firstImagePath = s
		break
	}
	if err := scanner.Err(); err != nil {
		os.Exit(1)
	}
	fmt.Printf("Enter Second Image Path: ")

	for scanner.Scan() {
		s := scanner.Text()
		secondImagePath = s
		break
	}
	if err := scanner.Err(); err != nil {
		os.Exit(1)
	}
	percentSimilar, err := perceptualhash.CalculateSimilarity(firstImagePath, secondImagePath, 512)
	if err != nil {
		log.Fatalf("%s", err)
	}
	fmt.Printf("Images are about %.3f similar\n", percentSimilar)
}
