package bookmarks

import (
	"github.com/mikzone/miknas/server/miknas"
)

type inDataAddBm struct {
	Kind string `json:"kind" binding:"required"`
	Name string `json:"name" binding:"required"`
	Url  string `json:"url" binding:"required"`
	Icon string `json:"icon"`
}

type inDataModifyBm struct {
	inDataAddBm
	Id int `json:"id" binding:"required"`
}

type inDataDelMb struct {
	Id int `json:"id" binding:"required"`
}

func addBookmark(ch *miknas.ContextHelper) {
	var loc inDataAddBm
	ch.BindJSON(&loc)
	userauth := ch.GetUserAuth()
	uid := userauth.MustGetUid()
	db := ch.GetApp().Db
	bm := &BookmarkItem{
		Uid:  uid,
		Kind: loc.Kind,
		Name: loc.Name,
		Url:  loc.Url,
		Icon: loc.Icon,
	}
	db.Create(bm)
	info := PackBookmarkInfo(bm)
	ch.UsLog("addBookmark", "bm", bm)
	ch.SucResp(info)
}

func modifyBookmark(ch *miknas.ContextHelper) {
	var loc inDataModifyBm
	ch.BindJSON(&loc)
	db := ch.GetApp().Db
	bm := GetBookmarkById(db, loc.Id)
	if bm == nil {
		ch.FailResp("书签不存在了")
		return
	}
	userauth := ch.GetUserAuth()
	uid := userauth.MustGetUid()
	if bm.Uid != uid {
		ch.FailResp("不是你的书签")
		return
	}
	db.Model(bm).Updates(BookmarkItem{
		Kind: loc.Kind,
		Name: loc.Name,
		Url:  loc.Url,
		Icon: loc.Icon,
	})
	info := PackBookmarkInfo(bm)
	ch.UsLog("modifyBookmarks", "bm", bm)
	ch.SucResp(info)
}

func deleteBookmark(ch *miknas.ContextHelper) {
	var loc inDataDelMb
	ch.BindJSON(&loc)
	db := ch.GetApp().Db
	bm := GetBookmarkById(db, loc.Id)
	if bm == nil {
		ch.FailResp("书签不存在了")
		return
	}
	userauth := ch.GetUserAuth()
	uid := userauth.MustGetUid()
	if bm.Uid != uid {
		ch.FailResp("不是你的书签")
		return
	}
	result := db.Delete(&bm)
	if result.Error != nil {
		ch.FailResp(result.Error.Error())
		return
	}
	ch.UsLog("deleteBookmarks", "bm", bm)
	ch.SucResp(result.RowsAffected)
}

func getUserBookmarks(ch *miknas.ContextHelper) {
	userauth := ch.GetUserAuth()
	uid := userauth.MustGetUid()
	db := ch.GetApp().Db
	var bms []BookmarkItem
	db.Where("uid = ?", uid).Find(&bms)
	infos := []miknas.H{}
	for _, bm := range bms {
		infos = append(infos, PackBookmarkInfo(&bm))
	}
	ch.SucResp(infos)
}

func regRoutes(ext *MikNasExt) {
	ext.POST("/add", addBookmark)
	ext.POST("/modify", modifyBookmark)
	ext.POST("/getall", getUserBookmarks)
	ext.POST("/delete", deleteBookmark)
}
