package miknas

import (
	"fmt"
	"io/fs"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	cp "github.com/otiai10/copy"

	getFolderSize "github.com/markthree/go-get-folder-size/src"
)

type IFsDriver interface {
	Stat(self IFsDriver, fspath string, needDirSize bool) (fs.FileInfo, error)
	ListDir(self IFsDriver, fspath string, needDirSize bool) ([]fs.FileInfo, error)
	CalcDirSize(self IFsDriver, fspath string) (int64, error)
	ReadFile(self IFsDriver, fspath string) ([]byte, error)
	WriteFile(self IFsDriver, fspath string, data []byte) error
	MkdirAll(self IFsDriver, fspath string) error
	RemoveAll(self IFsDriver, fspath string) error
	MoveAll(self IFsDriver, fspath string, toFspath string) error
	CopyAll(self IFsDriver, fspath string, toFspath string) error

	// 请求相关的
	WebDownload(self IFsDriver, ch *ContextHelper, fspath string)
	WebView(self IFsDriver, ch *ContextHelper, fspath string)
	WebSaveUploadedFile(self IFsDriver, ch *ContextHelper, file *multipart.FileHeader, fspath string) error
}

type IDiskFsDriver interface {
	// 硬盘上的文件系统
	IFsDriver
	MustAbs(fspath string) string   // 根据fspath返回在真实文件系统中的路径
	MustRel(fullpath string) string // 根据真实文件系统的路径来计算在Fs中的路径
}

type BaseFsDriver struct {
	rootPath string // need to be a absolute path
}

func (fsh *BaseFsDriver) GetRoot() string {
	return fsh.rootPath
}

// 根据相对于 workspace 的相对路径，获得对应的绝对路径
func (fsh *BaseFsDriver) MustAbs(relpath string) string {
	fullpath := filepath.Join(fsh.rootPath, relpath)
	result, err := filepath.Abs(fullpath)
	if err != nil {
		panic(err)
	}
	if !strings.HasPrefix(result, fsh.rootPath) {
		panic(NewFailRet("超出目录范围"))
	}
	return result
}

// 根据相对于 workspace 的相对路径，获得对应的绝对路径
func (fsh *BaseFsDriver) MustRel(fullpath string) string {
	if !strings.HasPrefix(fullpath, fsh.rootPath) {
		panic(NewFailRet("超出目录范围"))
	}
	result, err := filepath.Rel(fsh.rootPath, fullpath)
	if err != nil {
		panic(err)
	}
	return result
}

func (*BaseFsDriver) CalcDirSize(self IFsDriver, fspath string) (int64, error) {
	fsd := self.(IDiskFsDriver)
	fullpath := fsd.MustAbs(fspath)
	return getFolderSize.Parallel(fullpath)
}

// 读取某个文件
func (*BaseFsDriver) Stat(self IFsDriver, fspath string, needDirSize bool) (fs.FileInfo, error) {
	fsd := self.(IDiskFsDriver)
	fullpath := fsd.MustAbs(fspath)
	fileInfo, err := os.Stat(fullpath)
	if err != nil {
		return nil, err
	}
	isDir := fileInfo.IsDir()
	fsize := fileInfo.Size()
	if isDir {
		fsize = 0
		if needDirSize {
			newSize, err := self.CalcDirSize(self, fspath)
			if err == nil {
				fsize = newSize
			}
		}
	}
	fst := fileStat{
		name:    fileInfo.Name(),
		isDir:   isDir,
		size:    fsize,
		mode:    fileInfo.Mode(),
		modTime: fileInfo.ModTime(),
	}
	return &fst, nil
}

// 读取文件夹内容
func (*BaseFsDriver) ListDir(self IFsDriver, fspath string, needDirSize bool) ([]fs.FileInfo, error) {
	fsd := self.(IDiskFsDriver)
	fullpath := fsd.MustAbs(fspath)
	files, err := os.ReadDir(fullpath)
	if err != nil {
		return nil, err
	}
	fileInfos := []fs.FileInfo{}
	for _, file := range files {
		newFspath := filepath.Join(fspath, file.Name())
		fCInfo, _ := self.Stat(self, newFspath, needDirSize)
		if fCInfo != nil {
			fileInfos = append(fileInfos, fCInfo)
		}
	}
	return fileInfos, nil
}

func NewBaseFsDriver(rootPath string, checkExist bool) *BaseFsDriver {
	root, absErr := filepath.Abs(rootPath)
	if absErr != nil {
		panic(NewFailRet(absErr.Error()))
	}
	if checkExist {
		if _, err := os.Stat(root); err != nil {
			panic(NewFailRet(err.Error()))
		}
	}
	return &BaseFsDriver{rootPath: root}
}

func (*BaseFsDriver) ReadFile(self IFsDriver, fspath string) ([]byte, error) {
	fsd := self.(IDiskFsDriver)
	fullpath := fsd.MustAbs(fspath)
	return os.ReadFile(fullpath)
}

func (*BaseFsDriver) WriteFile(self IFsDriver, fspath string, data []byte) error {
	fsd := self.(IDiskFsDriver)
	fullpath := fsd.MustAbs(fspath)
	return os.WriteFile(fullpath, data, 0644)
}

func (*BaseFsDriver) MkdirAll(self IFsDriver, fspath string) error {
	fsd := self.(IDiskFsDriver)
	fullpath := fsd.MustAbs(fspath)
	return os.MkdirAll(fullpath, 0750)
}

func (*BaseFsDriver) RemoveAll(self IFsDriver, fspath string) error {
	fsd := self.(IDiskFsDriver)
	fullpath := fsd.MustAbs(fspath)
	return os.RemoveAll(fullpath)
}

func (*BaseFsDriver) MoveAll(self IFsDriver, fspath string, toFspath string) error {
	fsd := self.(IDiskFsDriver)
	fullpath := fsd.MustAbs(fspath)
	toFullpath := fsd.MustAbs(toFspath)
	err := os.Rename(fullpath, toFullpath)
	if err == nil {
		return nil
	}
	// 重命名失败的话,尝试复制和删除
	err2 := fsd.CopyAll(self, fspath, toFspath)
	if err2 != nil {
		return fmt.Errorf("ErrInCopy: %v", err2)
	}
	err3 := fsd.RemoveAll(self, fspath)
	if err3 != nil {
		return fmt.Errorf("ErrInRemove: %v", err3)
	}
	return nil
}

func (*BaseFsDriver) CopyAll(self IFsDriver, fspath string, toFspath string) error {
	fsd := self.(IDiskFsDriver)
	fullpath := fsd.MustAbs(fspath)
	toFullpath := fsd.MustAbs(toFspath)
	return cp.Copy(fullpath, toFullpath)
}

func (*BaseFsDriver) WebDownload(self IFsDriver, ch *ContextHelper, fspath string) {
	fsd := self.(IDiskFsDriver)
	fullpath := fsd.MustAbs(fspath)
	filename := filepath.Base(fullpath)
	ch.Ctx.FileAttachment(fullpath, filename)
}

func (*BaseFsDriver) WebView(self IFsDriver, ch *ContextHelper, fspath string) {
	fsd := self.(IDiskFsDriver)
	fullpath := fsd.MustAbs(fspath)
	ch.Ctx.File(fullpath)
}

func (*BaseFsDriver) WebSaveUploadedFile(self IFsDriver, ch *ContextHelper, file *multipart.FileHeader, fspath string) error {
	fsd := self.(IDiskFsDriver)
	fullpath := fsd.MustAbs(fspath)
	return ch.Ctx.SaveUploadedFile(file, fullpath)
}

var _ IDiskFsDriver = (*BaseFsDriver)(nil)

// ----------------- SingleFileFsDriver ---------------

// 单文件的一个Fstype
type SingleFileFsDriver struct {
	BaseFsDriver
}

func (fsd *SingleFileFsDriver) MustAbs(relpath string) string {
	rootPath := fsd.GetRoot()
	bN := filepath.Base(rootPath)
	if relpath != "" && relpath != bN && relpath != "/"+bN {
		panic(NewFailRet("非法路径"))
	}
	return rootPath
}

func (fsd *SingleFileFsDriver) MustRel(fullpath string) string {
	rootPath := fsd.GetRoot()
	if fullpath != rootPath {
		panic(NewFailRet("超出目录范围"))
	}
	return filepath.Base(rootPath)
}

// 读取文件夹内容
func (fsd *SingleFileFsDriver) ListDir(self IFsDriver, fspath string, needDirSize bool) ([]fs.FileInfo, error) {
	if fspath != "" && fspath != "/" {
		return nil, NewFailRet("目录不存在")
	}
	info, err := self.Stat(fsd, fspath, needDirSize)
	if err != nil {
		return nil, err
	}
	fileInfos := []fs.FileInfo{info}
	return fileInfos, nil
}

func NewSingleFileFsDriver(rootPath string, checkExist bool) IFsDriver {
	root, absErr := filepath.Abs(rootPath)
	if absErr != nil {
		panic(NewFailRet(absErr.Error()))
	}
	if checkExist {
		if _, err := os.Stat(root); err != nil {
			panic(NewFailRet(err.Error()))
		}
	}
	return &SingleFileFsDriver{
		*NewBaseFsDriver(rootPath, false),
	}
}

func NewCommFsDriver(rootPath string) IFsDriver {
	root, absErr := filepath.Abs(rootPath)
	if absErr != nil {
		panic(NewFailRet(absErr.Error()))
	}
	fInfo, err := os.Stat(root)
	if err != nil {
		panic(NewFailRet(err.Error()))
	}
	if fInfo.IsDir() {
		return NewBaseFsDriver(root, false)
	} else {
		return NewSingleFileFsDriver(root, false)
	}
}
