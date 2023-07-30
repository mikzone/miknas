package drive

import (
	"io/fs"
	"os"

	"github.com/mikzone/miknas/server/miknas"
)

func EnsureExisted(fullpath string) fs.FileInfo {
	fInfo, err := os.Stat(fullpath)
	if err != nil {
		panic(miknas.NewFailRet(err.Error()))
	}
	return fInfo
}

func EnsureIsExistedFile(fullpath string) fs.FileInfo {
	fInfo := EnsureExisted(fullpath)
	if fInfo.IsDir() {
		panic(miknas.NewFailRet("不是一个有效的文件"))
	}
	return fInfo
}

func EnsureIsExistedDir(fullpath string) fs.FileInfo {
	fInfo := EnsureExisted(fullpath)
	if !fInfo.IsDir() {
		panic(miknas.NewFailRet("不是一个有效的文件夹"))
	}
	return fInfo
}

func EnsureNotExisted(fullpath string) {
	_, err := os.Stat(fullpath)
	if err == nil {
		panic(miknas.NewFailRet("已存在同名文件名或目录"))
	} else if !os.IsNotExist(err) {
		panic(miknas.NewFailRet(err.Error()))
	}
}
