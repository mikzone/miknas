package note

import (
	"errors"

	"github.com/mikzone/miknas/server/miknas"
	"gorm.io/gorm"
)

// -------------------- DB Models --------------------

type NoteItem struct {
	miknas.PrimaryUintModel
	Folder      uint
	Uid         string `gorm:"index;size:64"`
	Title       string
	Content     string
	NoteAttachs []NoteAttach // has many
}

// 简约数据
type BriefNoteItem struct {
	ID     uint
	Folder uint
	Title  string
}

func GetItemById(db *gorm.DB, id uint, preloadAttach bool) *NoteItem {
	var item NoteItem
	_db := db
	if preloadAttach {
		_db = _db.Preload("NoteAttachs")
	}
	err := _db.First(&item, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		panic(miknas.NewFailRet(err.Error()))
	}
	return &item
}

// -------------------- Routes --------------------

type inDataAddItem struct {
	Folder  uint   `json:"folder" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}

type inDataModifyItem struct {
	Id       uint    `json:"id" binding:"required"`
	PFolder  *uint   `json:"folder"`
	PTitle   *string `json:"title"`
	PContent *string `json:"content"`
}

type inDataGetItem struct {
	Id uint `json:"id" binding:"required"`
}

type inDataDelItem struct {
	Id uint `json:"id" binding:"required"`
}

type inDataGetFolderItem struct {
	Folder uint `json:"folder" binding:"required"`
}

func PackNoteItem(item *NoteItem) miknas.H {
	attaches := []any{}
	for _, attach := range item.NoteAttachs {
		attaches = append(attaches, PackNoteAttach(&attach))
	}
	return miknas.H{
		"id":          item.ID,
		"uid":         item.Uid,
		"folder":      item.Folder,
		"title":       item.Title,
		"content":     item.Content,
		"noteAttachs": attaches,
		"modify":      item.UpdatedAt,
	}
}

func addItem(ch *miknas.ContextHelper) {
	var loc inDataAddItem
	ch.BindJSON(&loc)
	userauth := ch.GetUserAuth()
	uid := userauth.MustGetUid()
	db := ch.GetApp().Db
	folderId := loc.Folder
	folder := GetFolderById(db, folderId)
	if folder == nil {
		ch.FailResp("文件夹不存在，可能已被删除")
		return
	}
	if folder.Uid != uid {
		ch.FailResp("这不是你的文件夹")
		return
	}
	item := &NoteItem{
		Folder:  folderId,
		Uid:     uid,
		Title:   loc.Title,
		Content: loc.Content,
	}
	db.Create(item)
	info := PackNoteItem(item)
	ch.UsLog("addItem", info)
	ch.SucResp(info)
}

func getItem(ch *miknas.ContextHelper) {
	var loc inDataGetItem
	ch.BindJSON(&loc)
	db := ch.GetApp().Db
	item := GetItemById(db, loc.Id, true)
	if item == nil {
		ch.FailResp("笔记不存在了")
		return
	}
	userauth := ch.GetUserAuth()
	uid := userauth.MustGetUid()
	if item.Uid != uid {
		ch.FailResp("不是你的笔记")
		return
	}
	info := PackNoteItem(item)
	ch.SucResp(info)
}

func modifyItem(ch *miknas.ContextHelper) {
	var loc inDataModifyItem
	ch.BindJSON(&loc)
	db := ch.GetApp().Db.Debug()
	item := GetItemById(db, loc.Id, false)
	if item == nil {
		ch.FailResp("笔记不存在了")
		return
	}
	userauth := ch.GetUserAuth()
	uid := userauth.MustGetUid()
	if item.Uid != uid {
		ch.FailResp("不是你的笔记")
		return
	}

	modifyInfo := map[string]interface{}{}
	if loc.PFolder != nil {
		folderId := *loc.PFolder
		if folderId != item.Folder {
			folder := GetFolderById(db, folderId)
			if folder == nil {
				ch.FailResp("文件夹不存在，可能已被删除")
				return
			}
			if folder.Uid != uid {
				ch.FailResp("这不是你的文件夹")
				return
			}
			modifyInfo["folder"] = folderId
		}
	}

	if loc.PTitle != nil {
		modifyInfo["title"] = *loc.PTitle
	}

	if loc.PContent != nil {
		modifyInfo["content"] = *loc.PContent
	}

	result := db.Model(item).Updates(modifyInfo)
	ch.EnsureNoErr(result.Error)
	info := PackNoteItem(item)
	ch.UsLog("modifyItem", info)
	ch.SucResp(info)
}

func deleteItem(ch *miknas.ContextHelper) {
	var loc inDataDelItem
	ch.BindJSON(&loc)
	db := ch.GetApp().Db
	item := GetItemById(db, loc.Id, true)
	if item == nil {
		ch.FailResp("笔记不存在了")
		return
	}
	userauth := ch.GetUserAuth()
	uid := userauth.MustGetUid()
	if item.Uid != uid {
		ch.FailResp("不是你的笔记")
		return
	}

	if len(item.NoteAttachs) > 0 {
		ch.FailResp("需要删除该笔记下的所有附件才可以删除该笔记")
		return
	}

	info := PackNoteItem(item)
	result := db.Delete(&item)
	ch.EnsureNoErr(result.Error)
	ch.UsLog("deleteItem", info)
	ch.SucResp(result.RowsAffected)
}

func getFolderItem(ch *miknas.ContextHelper) {
	var loc inDataGetFolderItem
	ch.BindJSON(&loc)
	userauth := ch.GetUserAuth()
	uid := userauth.MustGetUid()
	db := ch.GetApp().Db
	var items []NoteItem
	db.Preload("NoteAttachs").Where("uid = ? And folder=?", uid, loc.Folder).Find(&items)
	infos := []miknas.H{}
	for _, item := range items {
		infos = append(infos, PackNoteItem(&item))
	}
	ch.SucResp(infos)
}

func PackBriefItem(item *BriefNoteItem) miknas.H {
	return miknas.H{
		"id":     item.ID,
		"folder": item.Folder,
		"title":  item.Title,
	}
}

func getUserItemBrief(ch *miknas.ContextHelper) {
	userauth := ch.GetUserAuth()
	uid := userauth.MustGetUid()
	db := ch.GetApp().Db
	var items []BriefNoteItem
	db.Model(&NoteItem{}).Where("uid = ?", uid).Find(&items)
	infos := []miknas.H{}
	for _, item := range items {
		infos = append(infos, PackBriefItem(&item))
	}
	ch.SucResp(infos)
}

type inDataListItem struct {
	Folder   uint   `json:"folder"`
	Search   string `json:"search"`
	PageNum  int    `json:"pageNum"`
	PageSize int    `json:"pageSize"`
}

func listItems(ch *miknas.ContextHelper) {
	var loc inDataListItem
	ch.BindJSON(&loc)
	userauth := ch.GetUserAuth()
	uid := userauth.MustGetUid()
	db := ch.GetApp().Db
	var items []NoteItem
	db = db.Where("uid = ?", uid)
	if loc.Folder > 0 {
		db = db.Where("folder = ?", loc.Folder)
	}
	if len(loc.Search) > 0 {
		search1 := "%" + loc.Search + "%"
		db = db.Where("title like ? OR content like ?", search1, search1)
	}

	pageSize := loc.PageSize
	switch {
	case pageSize > 30:
		pageSize = 30
	case pageSize <= 0:
		pageSize = 10
	}

	pageNum := loc.PageNum
	if pageNum <= 0 {
		pageNum = 1
	}
	offset := (pageNum - 1) * pageSize

	totalCnt := int64(0)
	countResult := db.Model(&NoteItem{}).Count(&totalCnt)
	ch.EnsureNoErr(countResult.Error)

	findResult := db.Order("id desc").Limit(pageSize).Offset(offset).Preload("NoteAttachs").Find(&items)
	ch.EnsureNoErr(findResult.Error)
	infos := []miknas.H{}
	for _, item := range items {
		infos = append(infos, PackNoteItem(&item))
	}
	ch.SucResp(miknas.H{
		"total": totalCnt,
		"notes": infos,
	})
}

func regItemRoutes(ext *MikNasExt) {
	ext.POST("/addItem", addItem)
	ext.POST("/getItem", getItem)
	ext.POST("/modifyItem", modifyItem)
	ext.POST("/getFolderItem", getFolderItem)
	ext.POST("/getUserItemBrief", getUserItemBrief)
	ext.POST("/deleteItem", deleteItem)
	ext.POST("/listItems", listItems)
}
