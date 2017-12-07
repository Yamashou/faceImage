package faceImage

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"strconv"
	"time"

	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
	"gocv.io/x/gocv"
)

func GetFaceImages(filename string, n uint, isGray bool) {
	file, err := imaging.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	rects := getImageToFacesImagePlace(filename)

	for i, r := range rects {
		img := imgResize(file, r, n)
		if isGray {
			img = getGray(img)
		}
		name := strconv.Itoa(int(time.Now().UnixNano()) + i)
		createImg(img, name+".jpg")
	}

}

func getImageToFacesImagePlace(filename string) []image.Rectangle {
	img := gocv.IMRead(filename, gocv.IMReadColor)
	defer img.Clone()

	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	classifier.Load(os.Getenv("GOPATH") + "/src/gocv.io/x/gocv/data/haarcascade_frontalface_default.xml")
	return classifier.DetectMultiScale(img)
}

func createImg(m image.Image, name string) {
	err := imaging.Save(m, name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func imgResize(m image.Image, r image.Rectangle, n uint) image.Image {
	rectcropimg := imaging.Crop(m, r)
	return resize.Resize(n, 0, rectcropimg, resize.Lanczos3)
}

func getGray(m image.Image) *image.Gray16 {
	bounds := m.Bounds()
	dest := image.NewGray16(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.Gray16Model.Convert(m.At(x, y))
			gray, _ := c.(color.Gray16)
			dest.Set(x, y, gray)
		}
	}
	return dest
}
