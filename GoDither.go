package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

func findClosestPaletteColorGrey(gray uint8, mean int) color.Gray {
	var newGray uint8
	/* if gray > 150 {
		newGray = 255
	} else {
		newGray = 0
	} */
	newGray = (gray / 255)
	Gray := new(color.Gray)
	Gray.Y = newGray
	// leave alpha channel as is
	return *Gray
}

func main() {

	imageFile, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer imageFile.Close()

	// Maybe change this to image.Decode for more formats
	decodedImage, err := png.Decode(imageFile)
	if err != nil {
		log.Fatal(err)
	}

	greyImg := image.NewGray(decodedImage.Bounds())
	ditheredImg := image.NewGray(decodedImage.Bounds())

	// * first we have to find median for a more accurate closet pixel
	var pixelVals []uint8
	for y := decodedImage.Bounds().Min.Y; y < decodedImage.Bounds().Max.Y; y++ {
		for x := decodedImage.Bounds().Min.X; x < decodedImage.Bounds().Max.X; x++ {
			oldPixel := color.GrayModel.Convert(decodedImage.At(x, y)).(color.Gray)
			greyImg.Set(x, y, oldPixel)
			pixelVals = append(pixelVals, oldPixel.Y)
		}
	}

	var sumOfPixels int = 0
	for pixelVal := range pixelVals {
		sumOfPixels += pixelVal
	}
	meanPixelVal := sumOfPixels / len(pixelVals)

	// important to use the grayscale image here instead of color
	for y := greyImg.Bounds().Min.Y; y < greyImg.Bounds().Max.Y; y++ {
		for x := greyImg.Bounds().Min.X; x < greyImg.Bounds().Max.X; x++ {
			oldPixel := greyImg.At(x, y).(color.Gray)
			var newPixel color.Gray = findClosestPaletteColorGrey(oldPixel.Y, meanPixelVal)
			quantError := oldPixel.Y - newPixel.Y
			// fmt.Printf("%v,", quantError)

			ditherFactorRight := greyImg.At(x+1, y).(color.Gray)
			ditherFactorRight.Y += (uint8)((float32)(quantError) * (7.0 / 16.0))

			ditherFactorUpLeft := greyImg.At(x-1, y+1).(color.Gray)
			ditherFactorUpLeft.Y += (uint8)((float32)(quantError) * (3.0 / 16.0))

			ditherFactorUp := greyImg.At(x, y+1).(color.Gray)
			ditherFactorUp.Y += (uint8)((float32)(quantError) * (5.0 / 16.0))

			ditherFactorUpRight := greyImg.At(x+1, y+1).(color.Gray)
			ditherFactorUpRight.Y += (uint8)((float32)(quantError) * (1.0 / 16.0))
			// * Set up new pixels
			ditheredImg.Set(x+1, y, ditherFactorRight)
			ditheredImg.Set(x-1, y+1, ditherFactorUpLeft)
			ditheredImg.Set(x, y+1, ditherFactorUp)
			ditheredImg.Set(x+1, y+1, ditherFactorUpRight)
		}
	}

	// Save new image
	newImgFile, err := os.Create("Dithered.png")
	if err != nil {
		log.Fatal(err)
	}
	defer newImgFile.Close()

	if err := png.Encode(newImgFile, ditheredImg); err != nil {
		log.Fatal(err)
	}
}
