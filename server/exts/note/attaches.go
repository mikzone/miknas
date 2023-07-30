package note

import (
	"errors"

	"github.com/mikzone/miknas/server/miknas"
	"gorm.io/gorm"
)

// -------------------- DB Models --------------------

type NoteAttach struct {
	miknas.PrimaryUintModel
	NoteItemID uint
	Uid        string `gorm:"index;size:64"`
	Type       string
	Data       string //json
}

func GetAttachById(db *gorm.DB, id uint) *NoteAttach {
	var attach NoteAttach
	err := db.First(&attach, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		panic(miknas.NewFailRet(err.Error()))
	}
	return &attach
}

// -------------------- Routes --------------------

type inDataAddAttach struct {
	NoteItemID uint   `json:"itemId" binding:"required"`
	Type       string `json:"type" binding:"required"`
	Data       string `json:"data"`
}

type inDataModifyAttach struct {
	inDataAddAttach
	Id uint `json:"id" binding:"required"`
}

type inDataDelAttach struct {
	Id uint `json:"id" binding:"required"`
}

func PackNoteAttach(attach *NoteAttach) miknas.H {
	return miknas.H{
		"id":     attach.ID,
		"itemId": attach.NoteItemID,
		"uid":    attach.Uid,
		"type":   attach.Type,
		"data":   attach.Data,
	}
}

func addAttach(ch *miknas.ContextHelper) {
	var loc inDataAddAttach
	ch.BindJSON(&loc)
	userauth := ch.GetUserAuth()
	uid := userauth.MustGetUid()
	db := ch.GetApp().Db
	itemId := loc.NoteItemID
	item := GetItemById(db, itemId, false)
	if item == nil {
		ch.FailResp("笔记不存在，可能已被删除")
		return
	}
	if item.Uid != uid {
		ch.FailResp("这不是你的笔记")
		return
	}
	attach := &NoteAttach{
		NoteItemID: itemId,
		Uid:        uid,
		Type:       loc.Type,
		Data:       loc.Data,
	}
	db.Create(attach)
	info := PackNoteAttach(attach)
	ch.UsLog("addAttach", info)
	ch.SucResp(info)
}

func modifyAttach(ch *miknas.ContextHelper) {
	var loc inDataModifyAttach
	ch.BindJSON(&loc)
	db := ch.GetApp().Db
	attach := GetAttachById(db, loc.Id)
	if attach == nil {
		ch.FailResp("附件不存在了")
		return
	}
	userauth := ch.GetUserAuth()
	uid := userauth.MustGetUid()
	if attach.Uid != uid {
		ch.FailResp("不是你的附件")
		return
	}
	itemId := loc.NoteItemID
	if itemId != attach.NoteItemID {
		item := GetItemById(db, itemId, false)
		if item == nil {
			ch.FailResp("笔记不存在，可能已被删除")
			return
		}
		if item.Uid != uid {
			ch.FailResp("这不是你的笔记")
			return
		}
	}

	result := db.Model(attach).Updates(NoteAttach{
		Type: loc.Type,
		Data: loc.Data,
	})
	ch.EnsureNoErr(result.Error)
	info := PackNoteAttach(attach)
	ch.UsLog("modifyAttach", info)
	ch.SucResp(info)
}

func deleteAttach(ch *miknas.ContextHelper) {
	var loc inDataDelAttach
	ch.BindJSON(&loc)
	db := ch.GetApp().Db
	attach := GetAttachById(db, loc.Id)
	if attach == nil {
		ch.FailResp("附件不存在了")
		return
	}
	userauth := ch.GetUserAuth()
	uid := userauth.MustGetUid()
	if attach.Uid != uid {
		ch.FailResp("不是你的附件")
		return
	}

	info := PackNoteAttach(attach)
	result := db.Delete(&attach)
	ch.EnsureNoErr(result.Error)
	ch.UsLog("deleteAttach", info)
	ch.SucResp(result.RowsAffected)
}

func regAttachRoutes(ext *MikNasExt) {
	ext.POST("/addAttach", addAttach)
	ext.POST("/modifyAttach", modifyAttach)
	ext.POST("/deleteAttach", deleteAttach)
}
