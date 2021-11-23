package image

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"

	"github.com/nfnt/resize"
)

func Resize(file multipart.File, filename string) (io.Reader, error) {
	img, formatType, _ := image.Decode(file)
	fmt.Println("formatType", formatType)

	m := resize.Resize(200, 0, img, resize.Lanczos3)

	return encode(m, formatType, filename)
}

func encode(image image.Image, formatType string, filename string) (io.Reader, error) {
	writer, _ := os.Create(filename)

	if formatType == "png" || formatType == "ping" {
		encoder := png.Encoder{CompressionLevel: png.BestSpeed}
		encoder.Encode(writer, image)
	}
	if formatType == "jpeg" {
		jpeg.Encode(writer, image, nil)
	}
	reader, _ := os.Open(filename)
	os.Remove(filename)
	return reader, nil
}
