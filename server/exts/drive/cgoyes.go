//go:build cgo

package drive

import (
	"fmt"

	"github.com/h2non/bimg"
)

func doImgThumb(srcBuf []byte, needSize int) ([]byte, error) {
	srcImg := bimg.NewImage(srcBuf)
	// size, err := srcImg.Size()
	// if err != nil {
	// 	return nil, fmt.Errorf("imagesize error")
	// }
	// width, height := fitSize(size.Width, size.Height, needSize, needSize)
	width, height := needSize, needSize
	options := bimg.Options{
		Width:   width,
		Height:  height,
		Crop:    true,
		Quality: 80,
	}
	return srcImg.Process(options)
}

func doGenImgThumb(srcBuf []byte, needSize int) ([]byte, string, error) {
	srcImg := bimg.NewImage(srcBuf)
	size, err := srcImg.Size()
	if err != nil {
		return nil, "", fmt.Errorf("imagesize error")
	}
	width, height := fitSize(size.Width, size.Height, needSize, needSize)
	options := bimg.Options{
		Width:   width,
		Height:  height,
		Type:    bimg.WEBP,
		Embed:   true,
		Quality: 95,
	}
	ret, err := srcImg.Process(options)
	return ret, ".webp", err
}

func init() {
	bimg.VipsCacheSetMax(0)
	bimg.VipsCacheSetMaxMem(0)
}
