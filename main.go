package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"golang.org/x/image/draw"
	"golang.org/x/term"
)

const approxFontHeightToWidthRatio = 0.5

func main() {
	if len(os.Args) != 2 {
		log.Fatal("invalid number of command line arguments")
	}

	reader, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	sourceImg, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	outputWidth, err := getTerminalWidth()
	if err != nil {
		log.Fatal(err)
	}

	ascii := []byte{'.', ',', '-', ':', ';', '=', '+', '*', '?', 'S', '$', '%', '&', 'M', '#', '@'}

	getOutputSymbol := func(grayScaleValue uint8) byte {
		return ascii[grayScaleValue/16]
	}

	var (
		srcBounds    = sourceImg.Bounds()
		imageWidth   = srcBounds.Max.X - srcBounds.Min.X
		imageHeight  = srcBounds.Max.Y - srcBounds.Min.Y
		outputHeight = int(float64(imageHeight) *
			(float64(outputWidth) / float64(imageWidth)) *
			approxFontHeightToWidthRatio)
		outputBounds = image.Rect(0, 0, outputWidth, outputHeight)
		output       = image.NewGray(outputBounds)
	)

	draw.NearestNeighbor.Scale(output, outputBounds, sourceImg, srcBounds, draw.Over, nil)

	for y := outputBounds.Min.Y; y < outputBounds.Max.Y; y++ {
		outputLine := make([]byte, outputWidth)

		for x := outputBounds.Min.X; x < outputBounds.Max.X; x++ {
			outputLine[x] = getOutputSymbol(output.GrayAt(x, y).Y)
		}

		fmt.Println(string(outputLine))
	}
}

func getTerminalWidth() (int, error) {
	terminalFileDescriptor := int(os.Stdin.Fd())
	width, _, err := term.GetSize(terminalFileDescriptor)

	return width, err
}
