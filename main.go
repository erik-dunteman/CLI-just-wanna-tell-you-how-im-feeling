package main

import (
	"fmt"
	"os"
	"time"

	"image"

	"gocv.io/x/gocv"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	vid, _ := gocv.VideoCaptureFile("video.mp4")

	ascii := []string{" ", ".", ":", "-", "=", "+", "*", "#", "%", "@"}

	fps := vid.Get(5)
	microSPF := int(1_000_000 / fps)

	ticker := time.NewTicker(time.Duration(microSPF) * time.Microsecond)
	img := gocv.NewMat()
	for {
		vid.Read(&img)
		if img.Empty() {
			break
		}

		gray := gocv.NewMat()
		gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

		maxCols, maxRows, err := terminal.GetSize(int(os.Stdin.Fd()))
		if err != nil {
			panic(err)
		}

		fx := float64(maxCols) / float64(gray.Cols())
		fy := float64(maxRows) / float64(gray.Rows())

		scaled := gocv.NewMat()
		pt := image.Point{
			X: 0,
			Y: 0,
		}

		gocv.Resize(gray, &scaled, pt, fx, fy, 0)

		str := ""
		for r := 0; r < scaled.Rows(); r++ {
			for c := 0; c < scaled.Cols(); c++ {
				b := float32(scaled.GetUCharAt(r, c)) / 255
				i := int(b * float32(len(ascii)))
				if i >= len(ascii) {
					i = len(ascii) - 1
				}
				str += ascii[i]
			}
			str += "\n"
		}

		// block
		<-ticker.C

		fmt.Println(str)
	}
}
