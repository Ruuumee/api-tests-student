package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
)

const testImageSize = 10

// GeneratePNG возвращает байты небольшого валидного PNG-изображения для загрузки фотографии в объявление
func GeneratePNG() []byte {
	img := image.NewRGBA(image.Rect(0, 0, testImageSize, testImageSize))
	for x := 0; x < testImageSize; x++ {
		for y := 0; y < testImageSize; y++ {
			r := uint8((float64(x) / float64(testImageSize)) * 255)
			g := uint8((float64(y) / float64(testImageSize)) * 255)
			b := uint8(150)

			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}

	buf := &bytes.Buffer{}
	if err := png.Encode(buf, img); err != nil {
		panic(fmt.Sprintf("Не удалось создать PNG: %s", err.Error()))
	}

	return buf.Bytes()
}
