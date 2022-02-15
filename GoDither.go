package main

import (
	"fmt"
	"image/color"
	"image/png"
	"log"
	"os"
)

func findClosestPaletteColorGrey(gray uint8) color.Gray {
	newGray := gray / 255
	Gray := new(color.Gray)
	Gray.Y = newGray
	// leave alpha channel as is
	return *Gray
}

func main() {

	image, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer image.Close()

	// Maybe change this to image.Decode for more formats
	pngImage, err := png.Decode(image)
	if err != nil {
		log.Fatal(err)
	}

	// Bruh how do I do this GoLang
	ditheredImg := image.newGray(image.Rectangle{})

	// fmt.Println(pngImage)

	// levels := []string{" ", "░", "▒", "▓", "█"}

	for y := pngImage.Bounds().Min.Y; y < pngImage.Bounds().Max.Y; y++ {
		for x := pngImage.Bounds().Min.X; x < pngImage.Bounds().Max.X; x++ {
			oldPixel := color.GrayModel.Convert(pngImage.At(x, y)).(color.Gray)
			var newPixel color.Gray = findClosestPaletteColorGrey(oldPixel.Y)
			/*
				level := c.Y / 51 // 51 * 5 = 255
				if level == 5 {
					level--
				}
				fmt.Print(levels[level])
			*/
		}
		fmt.Print("\n")
	}
}
