package mikauth

import (
	"errors"
	"time"

	"github.com/mikzone/miknas/server/miknas"
	"gorm.io/gorm"
)

type MikauthUser struct {
	Uid       string `gorm:"primarykey;size:64"`
	Pwd       string
	Name      string `gorm:"size:255"`
	Role      string `gorm:"size:64"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GetUserByUid(db *gorm.DB, uid string) *MikauthUser {
	var user MikauthUser
	err := db.First(&user, "uid = ?", uid).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		panic(miknas.NewFailRet(err.Error()))
	}
	return &user
}

func PackUserInfo(user *MikauthUser) miknas.H {
	return miknas.H{
		"uid":  user.Uid,
		"name": user.Name,
		"role": user.Role,
		"cts":  user.CreatedAt.Unix(),
	}
}
