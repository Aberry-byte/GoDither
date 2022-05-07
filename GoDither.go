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
	if gray != 0 {
		newGray = gray / (uint8(mean)) * 255
	} else {
		newGray = gray
	}
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
	for _, pixelVal := range pixelVals {
		sumOfPixels += (int)(pixelVal)
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

			ditherFactorBtmLeft := greyImg.At(x-1, y+1).(color.Gray)
			ditherFactorBtmLeft.Y += (uint8)((float32)(quantError) * (3.0 / 16.0))

			ditherFactorDown := greyImg.At(x, y+1).(color.Gray)
			ditherFactorDown.Y += (uint8)((float32)(quantError) * (5.0 / 16.0))

			ditherFactorBtmRight := greyImg.At(x+1, y+1).(color.Gray)
			ditherFactorBtmRight.Y += (uint8)((float32)(quantError) * (1.0 / 16.0))
			// * Set up new pixels
			ditheredImg.Set(x, y, newPixel)
			ditheredImg.Set(x+1, y, ditherFactorRight)
			ditheredImg.Set(x-1, y+1, ditherFactorBtmLeft)
			ditheredImg.Set(x, y+1, ditherFactorDown)
			ditheredImg.Set(x+1, y+1, ditherFactorBtmRight)
		}
	}
	// create dithered.png
	newImgFile, err := os.Create("Dithered.png")
	if err != nil {
		log.Fatal(err)
	}
	defer newImgFile.Close()

	// write out dithered image
	if err := png.Encode(newImgFile, ditheredImg); err != nil {
		log.Fatal(err)
	}

	// create black and white image
	newGrayFile, err := os.Create("Grey.png")
	if err != nil {
		log.Fatal(err)
	}
	defer newGrayFile.Close()

	// Write out black and white image
	if err := png.Encode(newGrayFile, greyImg); err != nil {
		log.Fatal(err)
	}
}
