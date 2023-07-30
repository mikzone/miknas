package miknas

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type IPathHelper interface {
	Root() string
	MustAbs(relpath string) string
	MustRel(fullpath string) string
}

type IFileSpace interface {
	GetFsType() string
	GetRelExt() IExtension
	SetRelExt(ext IExtension)
	// mode can be "r" or "w"
	Ensure(ch *ContextHelper, mode string)
	// return root abs path of fssubid
	GetSubRoot(ch *ContextHelper, fssubid string) string
}

type BasePathHelper struct {
	rootPath string // need to be a absolute path
}

func (ph *BasePathHelper) Root() string {
	return ph.rootPath
}

// 根据相对于 workspace 的相对路径，获得对应的绝对路径
func (ph *BasePathHelper) MustAbs(relpath string) string {
	fullpath := filepath.Join(ph.rootPath, relpath)
	result, err := filepath.Abs(fullpath)
	if err != nil {
		panic(err)
	}
	if !strings.HasPrefix(result, ph.rootPath) {
		panic(NewFailRet("超出目录范围"))
	}
	return result
}

// 根据相对于 workspace 的相对路径，获得对应的绝对路径
func (ph *BasePathHelper) MustRel(fullpath string) string {
	if !strings.HasPrefix(fullpath, ph.rootPath) {
		panic(NewFailRet("超出目录范围"))
	}
	result, err := filepath.Rel(fullpath, ph.rootPath)
	if err != nil {
		panic(err)
	}
	return result
}

func NewBasePathHelper(rootPath string, checkExist bool) IPathHelper {
	root, absErr := filepath.Abs(rootPath)
	if absErr != nil {
		panic(NewFailRet(absErr.Error()))
	}
	if checkExist {
		if _, err := os.Stat(root); os.IsNotExist(err) {
			panic(NewFailRet(err.Error()))
		}
	}
	return &BasePathHelper{rootPath: root}
}

var _ IPathHelper = (*BasePathHelper)(nil)

type BaseFileSpace struct {
	FsType   string
	rootPath string
	ResId    AuthResId
	Ext      IExtension
}

func (fs *BaseFileSpace) GetFsType() string {
	return fs.FsType
}

func (fs *BaseFileSpace) GetRelExt() IExtension {
	return fs.Ext
}

func (fs *BaseFileSpace) SetRelExt(ext IExtension) {
	fs.Ext = ext
}

func (fs *BaseFileSpace) GetSubRoot(ch *ContextHelper, fssubid string) string {
	return fs.rootPath
}

func (fs *BaseFileSpace) ModeRes(mode string) AuthResId {
	if len(fs.ResId) <= 0 {
		panic("ResId is empty")
	}
	return AuthResId(fmt.Sprintf("%s:%s", fs.ResId, mode))
}

func (fs *BaseFileSpace) Ensure(ch *ContextHelper, mode string) {
	if len(fs.ResId) <= 0 {
		return
	}
	resid := fs.ModeRes(mode)
	ch.Ensure(resid)
}

func (fs *BaseFileSpace) RegAuth(folderDesc string, sendClient bool) {
	ext := fs.GetRelExt()
	resid1 := fs.ModeRes("r")
	desc1 := fmt.Sprintf("读取 %s 下的文件和目录", folderDesc)
	resid2 := fs.ModeRes("w")
	desc2 := fmt.Sprintf("写入 %s 下的文件和目录", folderDesc)
	ext.RegAuth(resid1, desc1, sendClient)
	ext.RegAuth(resid2, desc2, sendClient)
}

func NewBaseFileSpace(fsType, rootPath string, resid AuthResId) *BaseFileSpace {
	return &BaseFileSpace{
		FsType:   fsType,
		rootPath: rootPath,
		ResId:    resid,
	}
}

var _ IFileSpace = (*BaseFileSpace)(nil)

// ----------------- 单一控制权限的FileSpace --------------

// 只需要单个权限的工作空间
type SimpleFileSpace struct {
	BaseFileSpace
}

func (fs *SimpleFileSpace) Ensure(ch *ContextHelper, mode string) {
	ch.Ensure(fs.ResId)
}

func (fs *SimpleFileSpace) RegAuth(folderDesc string, sendClient bool) {
	ext := fs.GetRelExt()
	ext.RegAuth(fs.ResId, folderDesc, sendClient)
}

func NewSimpleFileSpace(fsType, rootPath string, resid AuthResId) *SimpleFileSpace {
	return &SimpleFileSpace{
		BaseFileSpace: *NewBaseFileSpace(fsType, rootPath, resid),
	}
}

// ----------------- 其它函数 --------------

func DecomposeFsid(fsid string) (FsType, FsSubType string) {
	// decompose into FsType and FsSubType
	FsType, FsSubType, _ = strings.Cut(string(fsid), "_")
	return
}
