package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"net/http"
)

func splitImage(w http.ResponseWriter, r *http.Request) {
	// Decode the input image from the request body
	inputImage, _, err := image.Decode(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the bounds of the input image
	bounds := inputImage.Bounds()

	grid := [2]int{2, 2}

	var outImages []bytes.Buffer

	var outWriter bytes.Buffer

	var width, height int

	// Set the dimensions of the output images
	//const width, height = bounds.Max.Y/

	width = bounds.Max.X / grid[1]
	height = bounds.Max.Y / grid[0]

	// Split the input image into multiple smaller images
	for y := bounds.Min.Y; y < bounds.Max.Y; y += height {
		for x := bounds.Min.X; x < bounds.Max.X; x += width {
			// Create a new subimage for each smaller image
			subimage := inputImage.(interface {
				SubImage(r image.Rectangle) image.Image
			}).SubImage(image.Rect(x, y, x+width, y+height))

			// Encode the subimage as a JPEG
			var outputBuffer bytes.Buffer
			err = jpeg.Encode(&outputBuffer, subimage, nil)

			outWriter.Write(outputBuffer.Bytes())

			outImages = append(outImages, outputBuffer)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		}
	}
	//
	//outBytes := []byte{}
	//
	//oR := bytes.NewReader(outBytes)
	//
	//for _, img := range outImages {
	//	n, err := oR.Read(img.Bytes())
	//
	//	fmt.Println(n, err)
	//	break
	//}

	// Write the JPEG data to the response
	_, err = w.Write(outWriter.Bytes())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/split", splitImage)
	http.ListenAndServe(":128", nil)
}
