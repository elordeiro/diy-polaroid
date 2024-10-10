package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

const (
	canvasWidth   = 1200 // 4x6 inches at 300 DPI
	canvasHeight  = 1800 // 4x6 inches at 300 DPI
	lineThickness = 5    // 5 pixels
	linePositionY = 1440 // 4.8 inches from the top
	squareSize    = 1012 // 3.37x3.37 inches at 300 DPI
	border        = 94   // 0.25-inch border on all sides
)

func main() {
	// Get the path of the images
	var path string
	fmt.Print("Enter the path of the images: ")
	fmt.Scanln(&path)

	// Create the polaroids directory
	polaroidDir := filepath.Join(path, "polaroids")
	err := os.MkdirAll(polaroidDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create polaroid directory: %v", err)
	}
	fmt.Println("Polaroid directory created")

	// Get the list of image files
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	for _, file := range files {
		// Filter image files
		if file.IsDir() {
			continue
		}
		if !(filepath.Ext(file.Name()) == ".jpg" ||
			filepath.Ext(file.Name()) == ".jpeg" ||
			filepath.Ext(file.Name()) == ".png") {
			fmt.Println("Skipping file:", file.Name())
			continue
		}

		imagePath := filepath.Join(path, file.Name())
		img, err := imaging.Open(imagePath, imaging.AutoOrientation(true))
		if err != nil {
			log.Printf("Error opening image %s: %v", file.Name(), err)
			continue
		}

		// Crop image to square if necessary
		img = cropToSquare(img)

		// Resize the image
		img = imaging.Resize(img, squareSize, squareSize, imaging.Lanczos)

		// Create a gg context for the canvas
		dc := gg.NewContext(canvasWidth, canvasHeight)

		// Fill the canvas with white
		dc.SetColor(color.White)
		dc.Clear()

		// Draw the image onto the canvas at the specified offset
		dc.DrawImage(img, border, border)

		// Draw the black cutting line
		drawLine(dc)

		// Export the final image
		outputImage := dc.Image()
		outputPath := filepath.Join(polaroidDir, file.Name())
		err = imaging.Save(outputImage, outputPath)
		if err != nil {
			log.Printf("Error saving polaroid for %s: %v", file.Name(), err)
			continue
		}
		fmt.Printf("Polaroid created for: %s\n", file.Name())
	}
	fmt.Println("Done. Press any key to exit.")
	fmt.Scanln()
}

// cropToSquare crops an image to be a square
func cropToSquare(img image.Image) image.Image {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	if width > height {
		offset := (width - height) / 2
		return imaging.Crop(img, image.Rect(offset, 0, offset+height, height))
	} else if height > width {
		offset := (height - width) / 2
		return imaging.Crop(img, image.Rect(0, offset, width, offset+width))
	}
	return img
}

// drawLine draws a horizontal line across the image canvas at a specific Y position
func drawLine(dc *gg.Context) {
	dc.SetRGB(0, 0, 0) // Set color to black
	dc.SetLineWidth(lineThickness)
	dc.MoveTo(0, linePositionY)           // Starting point of the line
	dc.LineTo(canvasWidth, linePositionY) // End point of the line
	dc.Stroke()                           // Apply the line drawing
}
