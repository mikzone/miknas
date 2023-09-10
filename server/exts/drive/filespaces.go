package drive

import (
	"io/fs"
	"os"

	"github.com/mikzone/miknas/server/miknas"
)

func EnsureExisted(fsd miknas.IFsDriver, fspath string) fs.FileInfo {
	var fInfo fs.FileInfo
	var err error
	if fsd == nil {
		fInfo, err = os.Stat(fspath)
	} else {
		fInfo, err = fsd.Stat(fsd, fspath, false)
	}
	if err != nil {
		panic(miknas.NewFailRet(err.Error()))
	}
	return fInfo
}

func EnsureIsExistedFile(fsd miknas.IFsDriver, fspath string) fs.FileInfo {
	fInfo := EnsureExisted(fsd, fspath)
	if fInfo.IsDir() {
		panic(miknas.NewFailRet("不是一个有效的文件"))
	}
	return fInfo
}

func EnsureIsExistedDir(fsd miknas.IFsDriver, fspath string) fs.FileInfo {
	fInfo := EnsureExisted(fsd, fspath)
	if !fInfo.IsDir() {
		panic(miknas.NewFailRet("不是一个有效的文件夹"))
	}
	return fInfo
}

func EnsureNotExisted(fsd miknas.IFsDriver, fspath string) {
	_, err := fsd.Stat(fsd, fspath, false)
	if err == nil {
		panic(miknas.NewFailRet("已存在同名文件名或目录"))
	} else if !os.IsNotExist(err) {
		panic(miknas.NewFailRet(err.Error()))
	}
}

type FileClientInfo struct {
	Name   string `json:"name"`
	IsFile bool   `json:"isFile"`
	Size   int64  `json:"size"`
	Modify int64  `json:"modify"`
}

func PackFileClientInfo(fInfo fs.FileInfo) *FileClientInfo {
	return &FileClientInfo{
		Name:   fInfo.Name(),
		IsFile: !fInfo.IsDir(),
		Size:   fInfo.Size(),
		Modify: fInfo.ModTime().Unix(),
	}
}
