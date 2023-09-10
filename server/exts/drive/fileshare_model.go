// 处理文件分享相关的事情

package drive

import (
	"errors"
	"time"

	"github.com/mikzone/miknas/server/miknas"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type FileShareItem struct {
	miknas.PrimaryStringModel
	Uid    string `gorm:"index;size:64"`
	Name   string
	Pwd    string
	Vist   uint
	Fstype string
	Fsaddr string
	Intv   uint
}

func GetShareItemByMid(db *gorm.DB, mid string) *FileShareItem {
	var item FileShareItem
	err := db.First(&item, "sid = ?", mid).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		panic(miknas.NewFailRet(err.Error()))
	}
	return &item
}

func PackShareItemInfo(item *FileShareItem) miknas.H {
	return miknas.H{
		"sid":  item.Sid,
		"uid":  item.Uid,
		"name": item.Name,
		"bts":  item.CreatedAt.Unix(),
		"intv": item.Intv,
	}
}

func GenMid(db *gorm.DB) (string, error) {
	for i := 0; i < 5; i++ {
		guid := xid.New()
		mid := guid.String()
		if GetShareItemByMid(db, mid) == nil {
			return mid, nil
		}
	}
	return "", miknas.NewFailRet("生成唯一id失败过多")
}

func IsShareItemExpired(item *FileShareItem) bool {
	now := time.Now()
	if item.Intv > 0 && item.CreatedAt.Unix()+int64(item.Intv) <= now.Unix() {
		return true
	}
	return false
}
