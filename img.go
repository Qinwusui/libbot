package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

func GetPic(imgUrl string) {
	// os.Remove("qr.png")
	client := http.Client{}
	request, _ := http.NewRequest("GET", imgUrl, nil)
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	bts, _ := io.ReadAll(resp.Body)
	os.WriteFile("qr.jpg", bts, 0777)
	file, _ := os.Open("qr.jpg")
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Println(err)
		return
	}
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			oldPixel := img.At(x, y)
			r, g, b, a := oldPixel.RGBA()
			newPixel := color.RGBA{
				uint8(clamp(float64(r)/0xfff*1.2) * 0xff),
				uint8(clamp(float64(g)/0xfff*1.2) * 0xff),
				uint8(clamp(float64(b)/0xfff*1.2) * 0xff),
				uint8(a / 0xff),
			}
			newImg.Set(x, y, newPixel)
		}
	}
	out, err := os.Create("newImg.jpg")
	if err != nil {
		log.Println(err)
		return
	}
	defer out.Close() 
	if err := jpeg.Encode(out, newImg, nil); err != nil {
		log.Println(err)
		return
	}

	bmp, _ := gozxing.NewBinaryBitmapFromImage(newImg)
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		log.Println(err)
		return
	}
	bts = []byte(result.String())
	log.Println(string(bts))
}
func clamp(x float64) float64 {
	if x < 0 {
		return 0
	} else if x > 1 {
		return 1
	} else {
		return x
	}
}
