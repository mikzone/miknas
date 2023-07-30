package note

import (
	"github.com/mikzone/miknas/server/miknas"
)

type MikNasExt struct {
	miknas.Extension
}

func New() *MikNasExt {
	return &MikNasExt{miknas.NewExtension("Note")}
}

func (ext *MikNasExt) OnBind() {
	// you can register config, auth, routes in here
	ext.RegAuth(ext.Res("vist"), "使用笔记", false)
	regFolderRoutes(ext)
	regItemRoutes(ext)
	regAttachRoutes(ext)
}

func (ext *MikNasExt) OnInit() {
	// only in init, you can access db, workspace, loaded configs
	// you can register your filespace, init your db here
	db := ext.App.Db
	db.AutoMigrate(&NoteFolder{})
	db.AutoMigrate(&NoteItem{})
	db.AutoMigrate(&NoteAttach{})
}
