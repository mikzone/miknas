package miknas

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"time"
)

type IFileSpace interface {
	GetFstype() string
	GetRelExt() IExtension
	SetRelExt(ext IExtension)
	// mode can be "r" or "w"
	Ensure(ch *ContextHelper, mode string)
	NewFsDriver(ch *ContextHelper, fssubid string) IFsDriver
	NewFsDriverByAddr(fsaddr string) IFsDriver
	GetAddr(self IFileSpace, ch *ContextHelper, fssubid string, fspath string) string
}

type fileStat struct {
	name    string
	size    int64
	isDir   bool
	mode    fs.FileMode
	modTime time.Time
}

func (fst *fileStat) Name() string       { return fst.name }
func (fst *fileStat) Size() int64        { return fst.size }
func (fst *fileStat) Mode() fs.FileMode  { return fst.mode }
func (fst *fileStat) ModTime() time.Time { return fst.modTime }
func (fst *fileStat) IsDir() bool        { return fst.isDir }
func (fst *fileStat) Sys() any           { return nil }

var _ fs.FileInfo = (*fileStat)(nil)

type BaseFileSpace struct {
	Fstype   string
	rootPath string
	ResId    AuthResId
	Ext      IExtension
}

func (fsp *BaseFileSpace) GetFstype() string {
	return fsp.Fstype
}

func (fsp *BaseFileSpace) GetRelExt() IExtension {
	return fsp.Ext
}

func (fsp *BaseFileSpace) SetRelExt(ext IExtension) {
	fsp.Ext = ext
}

func (fsp *BaseFileSpace) NewFsDriver(ch *ContextHelper, fssubid string) IFsDriver {
	return NewBaseFsDriver(fsp.rootPath, true)
}

func (fsp *BaseFileSpace) GetAddr(self IFileSpace, ch *ContextHelper, fssubid string, fspath string) string {
	self.Ensure(ch, "w")
	fsd := self.NewFsDriver(ch, fssubid)
	_, err := fsd.Stat(fsd, fspath, false)
	ch.EnsureNoErr(err)
	return fsp.Fstype + "|" + fssubid + "|" + fspath
}

func (fsp *BaseFileSpace) NewFsDriverByAddr(fsaddr string) IFsDriver {
	parts := strings.Split(fsaddr, "|")
	if parts[0] != fsp.Fstype {
		return nil
	}
	fspath := parts[2]
	fullpath := filepath.Join(fsp.rootPath, fspath)
	return NewCommFsDriver(fullpath)
}

func (fsp *BaseFileSpace) ModeRes(mode string) AuthResId {
	if len(fsp.ResId) <= 0 {
		panic("ResId is empty")
	}
	return AuthResId(fmt.Sprintf("%s:%s", fsp.ResId, mode))
}

func (fsp *BaseFileSpace) Ensure(ch *ContextHelper, mode string) {
	if len(fsp.ResId) <= 0 {
		return
	}
	resid := fsp.ModeRes(mode)
	ch.Ensure(resid)
}

func (fsp *BaseFileSpace) RegAuth(folderDesc string, sendClient bool) {
	ext := fsp.GetRelExt()
	resid1 := fsp.ModeRes("r")
	desc1 := fmt.Sprintf("读取 %s 下的文件和目录", folderDesc)
	resid2 := fsp.ModeRes("w")
	desc2 := fmt.Sprintf("写入 %s 下的文件和目录", folderDesc)
	ext.RegAuth(resid1, desc1, sendClient)
	ext.RegAuth(resid2, desc2, sendClient)
}

func NewBaseFileSpace(fstype, rootPath string, resid AuthResId) *BaseFileSpace {
	return &BaseFileSpace{
		Fstype:   fstype,
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

func (fsp *SimpleFileSpace) Ensure(ch *ContextHelper, mode string) {
	ch.Ensure(fsp.ResId)
}

func (fsp *SimpleFileSpace) RegAuth(folderDesc string, sendClient bool) {
	ext := fsp.GetRelExt()
	ext.RegAuth(fsp.ResId, folderDesc, sendClient)
}

func NewSimpleFileSpace(fstype, rootPath string, resid AuthResId) *SimpleFileSpace {
	return &SimpleFileSpace{
		BaseFileSpace: *NewBaseFileSpace(fstype, rootPath, resid),
	}
}

// ----------------- 其它函数 --------------

func DecomposeFsid(fsid string) (Fstype, fssubid string) {
	// decompose into Fstype and fssubid
	Fstype, fssubid, _ = strings.Cut(string(fsid), "_")
	return
}
