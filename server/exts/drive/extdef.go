package drive

import (
	"github.com/mikzone/miknas/server/miknas"
)

type MikNasExt struct {
	miknas.Extension
}

func New() *MikNasExt {
	return &MikNasExt{miknas.NewExtension("Drive")}
}

func (ext *MikNasExt) OnBind() {
	regRoutes(ext)
	regFileShareRoutes(ext)
	ext.RegStrConf("MIKNAS_DRIVE_NEED_DIR_SIZE", "1", "是否计算文件夹大小", false)
}

func (ext *MikNasExt) OnInit() {
	db := ext.App.Db
	db.AutoMigrate(&FileShareItem{})
	RegShareFileSpace(ext)
}
