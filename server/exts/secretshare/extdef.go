package secretshare

import (
	"github.com/mikzone/miknas/server/miknas"
)

type MikNasExt struct {
	miknas.Extension
}

func New() *MikNasExt {
	return &MikNasExt{miknas.NewExtension("SecretShare")}
}

func (ext *MikNasExt) OnBind() {
	// you can register config, auth, routes in here
	ext.RegAuth(ext.Res("vist"), "使用密文分享", false)
	regRoutes(ext)
}

func (ext *MikNasExt) OnInit() {
	// only in init, you can access db, workspace, loaded configs
	// you can register your filespace, init your db here
	db := ext.App.Db
	db.AutoMigrate(&SecretShareItem{})
}
