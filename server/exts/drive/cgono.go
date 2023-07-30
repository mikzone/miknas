//go:build !cgo

package drive

import (
	"bytes"
	"fmt"

	"github.com/disintegration/imaging"
)

func doImgThumb(srcBuf []byte, needSize int) ([]byte, error) {
	imgBuf := bytes.NewBuffer(srcBuf)
	image, err := imaging.Decode(imgBuf)
	if err != nil {
		return nil, err
	}
	// srcW := image.Bounds().Dx()
	// srcH := image.Bounds().Dy()
	// width, height := fitSize(srcW, srcH, needSize, needSize)
	// thumbImg := imaging.Resize(image, width, height, imaging.Linear)
	thumbImg := imaging.Fill(image, needSize, needSize, imaging.Center, imaging.Linear)
	var buf bytes.Buffer
	err = imaging.Encode(&buf, thumbImg, imaging.PNG)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func doGenImgThumb(srcBuf []byte, needSize int) ([]byte, string, error) {
	imgBuf := bytes.NewBuffer(srcBuf)
	image, err := imaging.Decode(imgBuf)
	if err != nil {
		return nil, "", fmt.Errorf("DecodeImage发生错误: %s", err.Error())
	}
	srcW := image.Bounds().Dx()
	srcH := image.Bounds().Dy()
	width, height := fitSize(srcW, srcH, 200, 200)
	thumbImg := imaging.Resize(image, width, height, imaging.Lanczos)
	var buf bytes.Buffer
	err = imaging.Encode(&buf, thumbImg, imaging.PNG)
	if err != nil {
		return nil, "", fmt.Errorf("EncodeImage发生错误: %s", err.Error())
	}
	return buf.Bytes(), ".png", nil
}
