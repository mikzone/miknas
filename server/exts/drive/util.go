package drive

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"
)

var Exts2Type = map[string]string{
	"apng":  "img",
	"avif":  "img",
	"gif":   "img",
	"jpg":   "img",
	"jpeg":  "img",
	"jfif":  "img",
	"pjpeg": "img",
	"pjp":   "img",
	"png":   "img",
	"svg":   "img",
	"webp":  "img",
	"bmp":   "img",
	"ico":   "img",
	"cur":   "img",
	"tif":   "img",
	"tiff":  "img",

	"wmv":   "video",
	"avi":   "video",
	"mpeg":  "video",
	"mpg":   "video",
	"rm":    "video",
	"rmvb":  "video",
	"flv":   "video",
	"mp4":   "video",
	"3gp":   "video",
	"mov":   "video",
	"divx":  "video",
	"vob":   "video",
	"mkv":   "video",
	"fli":   "video",
	"flc":   "video",
	"f4v":   "video",
	"m4v":   "video",
	"mod":   "video",
	"m2t":   "video",
	"webm":  "video",
	"mts":   "video",
	"m2ts":  "video",
	"3g2":   "video",
	"mpe":   "video",
	"ts":    "video",
	"div":   "video",
	"lavf":  "video",
	"dirac": "video",

	"txt":  "text",
	"md":   "text",
	"json": "text",
	"c":    "text",
	"cpp":  "text",
	"h":    "text",
	"hpp":  "text",
	"js":   "text",
	"vue":  "text",
	"html": "text",
	"htm":  "text",
	"css":  "text",
	"sql":  "text",
	"log":  "text",
	"ini":  "text",
	"yml":  "text",
	"yaml": "text",
	"py":   "text",
	"go":   "text",
	"java": "text",
	"php":  "text",
	"xml":  "text",
}

func CalcExt(filename string) string {
	ext := filepath.Ext(filename)
	ext = strings.TrimPrefix(ext, ".")
	if len(ext) <= 0 {
		return ""
	}
	ext = strings.ToLower(ext)
	return ext
}

func GetFileType(ext string) string {
	ftype, ok := Exts2Type[ext]
	if !ok {
		return ""
	}
	return ftype
}

func GetImageSize(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0
	}

	defer file.Close()

	config, format, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0
	}

	fmt.Println("Width:", config.Width, "Height:", config.Height, "Format:", format)

	return config.Width, config.Height
}
