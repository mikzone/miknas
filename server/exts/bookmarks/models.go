package bookmarks

import (
	"errors"

	"github.com/mikzone/miknas/server/miknas"
	"gorm.io/gorm"
)

type BookmarkItem struct {
	miknas.PrimaryUintModel
	Uid  string `gorm:"index;size:64"`
	Kind string
	Name string
	Url  string
	Icon string
}

func GetBookmarkById(db *gorm.DB, id int) *BookmarkItem {
	var item BookmarkItem
	err := db.First(&item, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		panic(miknas.NewFailRet(err.Error()))
	}
	return &item
}

func PackBookmarkInfo(bm *BookmarkItem) miknas.H {
	return miknas.H{
		"id":   bm.ID,
		"uid":  bm.Uid,
		"kind": bm.Kind,
		"name": bm.Name,
		"url":  bm.Url,
		"icon": bm.Icon,
	}
}
