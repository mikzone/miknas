package note

import (
	"errors"

	"github.com/mikzone/miknas/server/miknas"
	"gorm.io/gorm"
)

// -------------------- DB Models --------------------

type NoteFolder struct {
	miknas.PrimaryUintModel
	Uid    string `gorm:"index;size:64"`
	Name   string `gorm:"size:255"`
	Parent uint   `gorm:"index"`
}

func GetFolderById(db *gorm.DB, id uint) *NoteFolder {
	var folder NoteFolder
	err := db.First(&folder, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		panic(miknas.NewFailRet(err.Error()))
	}
	return &folder
}

func GetSubFolderByName(db *gorm.DB, parent uint, name string) *NoteFolder {
	var folder NoteFolder
	err := db.First(&folder, "parent = ? AND name = ?", parent, name).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		panic(miknas.NewFailRet(err.Error()))
	}
	return &folder
}

// -------------------- Routes --------------------

type inDataAddFolder struct {
	Name   string `json:"name" binding:"required"`
	Parent uint   `json:"parent"`
}

type inDataModifyFolder struct {
	Id      uint    `json:"id" binding:"required"`
	PName   *string `json:"name"`
	PParent *uint   `json:"parent"`
}

type inDataDelFolder struct {
	Id uint `json:"id" binding:"required"`
}

func PackNoteFolder(folder *NoteFolder) miknas.H {
	return miknas.H{
		"id":     folder.ID,
		"name":   folder.Name,
		"parent": folder.Parent,
	}
}

func addFolder(ch *miknas.ContextHelper) {
	var loc inDataAddFolder
	ch.BindJSON(&loc)
	userauth := ch.GetUserAuth()
	uid := userauth.MustGetUid()
	db := ch.GetApp().Db
	parent := loc.Parent
	if parent != 0 {
		pFolder := GetFolderById(db, parent)
		if pFolder == nil {
			ch.FailResp("创建文件夹失败,父文件夹可能已被删除")
			return
		}
		if pFolder.Uid != uid {
			ch.FailResp("这不是你的文件夹")
			return
		}
	}
	if GetSubFolderByName(db, parent, loc.Name) != nil {
		ch.FailResp("已有同名文件夹")
		return
	}
	folder := &NoteFolder{
		Uid:    uid,
		Name:   loc.Name,
		Parent: loc.Parent,
	}
	db.Create(folder)
	info := PackNoteFolder(folder)
	ch.UsLog("addFolder", info)
	ch.SucResp(info)
}

func modifyFolder(ch *miknas.ContextHelper) {
	var loc inDataModifyFolder
	ch.BindJSON(&loc)
	db := ch.GetApp().Db
	folder := GetFolderById(db, loc.Id)
	if folder == nil {
		ch.FailResp("文件夹不存在了")
		return
	}
	userauth := ch.GetUserAuth()
	uid := userauth.MustGetUid()
	if folder.Uid != uid {
		ch.FailResp("不是你的文件夹")
		return
	}
	newName := folder.Name
	modifyInfo := map[string]interface{}{}
	if loc.PName != nil {
		newName = *loc.PName
		if len(newName) <= 0 {
			ch.FailResp("重命名失败，新的名称不能为空")
			return
		}
		modifyInfo["name"] = newName
	}
	if loc.PParent != nil {
		parent := *loc.PParent
		if parent != folder.Parent && parent != 0 {
			pFolder := GetFolderById(db, parent)
			if pFolder == nil {
				ch.FailResp("移动失败,目标文件夹可能已被删除")
				return
			}
			if pFolder.Uid != uid {
				ch.FailResp("移动失败,目标文件夹不是你的")
				return
			}
		}
		if GetSubFolderByName(db, parent, newName) != nil {
			ch.FailResp("移动失败,已有同名文件夹")
			return
		}
		modifyInfo["parent"] = parent
	}
	result := db.Model(folder).Updates(modifyInfo)
	ch.EnsureNoErr(result.Error)
	info := PackNoteFolder(folder)
	ch.UsLog("modifyFolder", info)
	ch.SucResp(info)
}

func deleteFolder(ch *miknas.ContextHelper) {
	var loc inDataDelFolder
	ch.BindJSON(&loc)
	db := ch.GetApp().Db
	folder := GetFolderById(db, loc.Id)
	if folder == nil {
		ch.FailResp("文件夹不存在了")
		return
	}
	userauth := ch.GetUserAuth()
	uid := userauth.MustGetUid()
	if folder.Uid != uid {
		ch.FailResp("不是你的文件夹")
		return
	}

	var subFolders []NoteFolder
	if err := db.Where("parent = ?", loc.Id).Find(&subFolders).Error; err != nil {
		ch.FailResp("删除检测子目录失败: %s", err.Error())
		return
	}
	if len(subFolders) > 0 {
		ch.FailResp("不能删除非空文件夹，该文件夹下含有子文件夹")
		return
	}

	var items []NoteItem
	if err := db.Where("uid = ? And folder=?", uid, loc.Id).Find(&items).Error; err != nil {
		ch.FailResp("删除检测文件夹下的笔记失败: %s", err.Error())
		return
	}
	if len(items) > 0 {
		ch.FailResp("不能删除非空文件夹，该文件夹下含有笔记")
		return
	}

	info := PackNoteFolder(folder)
	result := db.Delete(&folder)
	ch.EnsureNoErr(result.Error)
	ch.UsLog("deleteFolder", info)
	ch.SucResp(result.RowsAffected)
}

func getUserFolder(ch *miknas.ContextHelper) {
	userauth := ch.GetUserAuth()
	uid := userauth.MustGetUid()
	db := ch.GetApp().Db
	var folders []NoteFolder
	db.Where("uid = ?", uid).Find(&folders)
	infos := []miknas.H{}
	for _, folder := range folders {
		infos = append(infos, PackNoteFolder(&folder))
	}
	ch.SucResp(infos)
}

func regFolderRoutes(ext *MikNasExt) {
	ext.POST("/addFolder", addFolder)
	ext.POST("/modifyFolder", modifyFolder)
	ext.POST("/getUserFolder", getUserFolder)
	ext.POST("/deleteFolder", deleteFolder)
}
