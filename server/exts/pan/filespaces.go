package pan

import (
	"os"

	"github.com/mikzone/miknas/server/miknas"
)

type PrivFileSpace struct {
	miknas.SimpleFileSpace
}

func (fs *PrivFileSpace) GetSubRoot(ch *miknas.ContextHelper, fssubid string) string {
	uid := ch.GetUserAuth().MustGetUid()
	if len(uid) <= 0 {
		panic(miknas.NewFailRet("invalid folder"))
	}
	app := ch.GetApp()
	fullpath := app.WorkSpace.MustAbs("priv_files/" + uid)
	os.MkdirAll(fullpath, 0777)
	return fullpath
}
