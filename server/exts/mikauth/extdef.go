package mikauth

import (
	"github.com/mikzone/miknas/server/miknas"
)

type MikNasExt struct {
	miknas.Extension
}

func New() *MikNasExt {
	return &MikNasExt{miknas.NewExtension("MikAuth")}
}

func (ext *MikNasExt) OnBind() {
	// register auth
	ext.RegAuth(ext.Res("manager"), "管理用户和权限", true)
	ext.RegAuth(ext.Res("vist"), "管理自己账号", false)
	// routes
	regUserRoutes(ext)
	regRoleRoutes(ext)
}

func (ext *MikNasExt) OnInit() {
	db := ext.App.Db
	db.AutoMigrate(&MikauthUser{})
	db.AutoMigrate(&MikauthRole{})
}
