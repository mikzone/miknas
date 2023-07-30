package secretshare

import (
	"errors"
	"time"

	"github.com/mikzone/miknas/server/miknas"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type SecretShareItem struct {
	miknas.PrimaryUintModel
	Uid      string `gorm:"index;size:64"`
	Mid      string `gorm:"uniqueIndex"`
	Name     string
	Txt      string
	Hint     string
	ExpireAt time.Time
}

func GetSecretByMid(db *gorm.DB, mid string) *SecretShareItem {
	var item SecretShareItem
	err := db.First(&item, "mid = ?", mid).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		panic(miknas.NewFailRet(err.Error()))
	}
	return &item
}

func PackSecretInfo(bm *SecretShareItem) miknas.H {
	return miknas.H{
		"uid":  bm.Uid,
		"mid":  bm.Mid,
		"name": bm.Name,
		"txt":  bm.Txt,
		"hint": bm.Hint,
		"bts":  bm.CreatedAt.Unix(),
		"ts":   bm.ExpireAt.Unix(),
	}
}

func GenMid(db *gorm.DB) (string, error) {
	for i := 0; i < 5; i++ {
		guid := xid.New()
		mid := guid.String()
		if GetSecretByMid(db, mid) == nil {
			return mid, nil
		}
	}
	return "", miknas.NewFailRet("生成唯一id失败过多")
}
