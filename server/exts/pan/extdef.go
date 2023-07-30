package pan

import (
	"github.com/mikzone/miknas/server/miknas"
)

type MikNasExt struct {
	miknas.Extension
}

func New() *MikNasExt {
	return &MikNasExt{miknas.NewExtension("Pan")}
}

func (ext *MikNasExt) OnBind() {
	// you can register config, auth, routes in here
	ext.RegAuth(ext.Res("vist"), "使用个人网盘", false)
	regRoutes(ext)
}

func (ext *MikNasExt) OnInit() {
	// only in init, you can access db, workspace, loaded configs
	// you can register your filespace, init your db here
	privFileSpace := &PrivFileSpace{
		*miknas.NewSimpleFileSpace("Priv", "", ext.Res("vist")),
	}
	ext.RegFileSpace(privFileSpace)
}
