package rolectrl

import (
	"github.com/mikzone/miknas/server/miknas"
)

type MikNasExt struct {
	miknas.Extension
}

func New() *MikNasExt {
	return &MikNasExt{miknas.NewExtension("RoleCtrl")}
}

func (ext *MikNasExt) OnBind() {
	ext.RegAuth(ext.Res("vist"), "角色添加、权限管理", true)
	// routes
	regRoleRoutes(ext)
}

func (ext *MikNasExt) OnInit() {
	db := ext.App.Db
	db.AutoMigrate(&RolectrlRole{})
	tourRole := GetRoleById(db, "tour")
	if tourRole == nil {
		AddOneRole(db, "tour")
	}
}
