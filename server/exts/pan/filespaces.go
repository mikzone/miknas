package pan

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mikzone/miknas/server/miknas"
)

type PrivFileSpace struct {
	miknas.SimpleFileSpace
}

func (fsp *PrivFileSpace) NewFsDriver(ch *miknas.ContextHelper, fssubid string) miknas.IFsDriver {
	uid := ch.GetUserAuth().MustGetUid()
	if len(uid) <= 0 {
		panic(miknas.NewFailRet("invalid folder"))
	}
	app := ch.GetApp()
	fullpath := app.WorkSpace.MustAbs("priv_files/" + uid)
	os.MkdirAll(fullpath, 0777)
	return miknas.NewBaseFsDriver(fullpath, true)
}

func (fsp *PrivFileSpace) GetAddr(self miknas.IFileSpace, ch *miknas.ContextHelper, fssubid string, fspath string) string {
	self.Ensure(ch, "w")
	fsd := self.NewFsDriver(ch, fssubid)
	_, err := fsd.Stat(fsd, fspath, false)
	ch.EnsureNoErr(err)
	uid := ch.GetUserAuth().MustGetUid()
	return fsp.Fstype + "|" + uid + "|" + fspath
}

func (fsp *PrivFileSpace) NewFsDriverByAddr(fsaddr string) miknas.IFsDriver {
	parts := strings.Split(fsaddr, "|")
	if parts[0] != fsp.Fstype {
		return nil
	}
	uid := parts[1]
	if len(uid) <= 0 {
		panic(miknas.NewFailRet("非法目录"))
	}
	fspath := parts[2]
	app := fsp.GetRelExt().GetApp()
	fullpath := filepath.Join("priv_files", uid, fspath)
	fullpath = app.WorkSpace.MustAbs(fullpath)
	return miknas.NewCommFsDriver(fullpath)
}
