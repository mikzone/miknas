package devtool

import (
	"github.com/mikzone/miknas/server/miknas"
)

type MikNasExt struct {
	miknas.Extension
}

func New() *MikNasExt {
	return &MikNasExt{miknas.NewExtension("DevTool")}
}

func (ext *MikNasExt) OnBind() {
	// you can register config, auth, routes in here
	ext.RegAuth(ext.Res("vist"), "使用开发管理工具，一般只有管理员可用", false)
	regRoutes(ext)
}

func (ext *MikNasExt) OnInit() {
	// only in init, you can access db, workspace, loaded configs
	// you can register your filespace, init your db here
	wsRoot := ext.App.WorkSpace.MustAbs("")
	wsFileSpace := miknas.NewSimpleFileSpace("Ws", wsRoot, ext.Res("vist"))
	ext.RegFileSpace(wsFileSpace)
}
