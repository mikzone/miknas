package rolectrl

import (
	"errors"
	"time"

	"github.com/mikzone/miknas/server/miknas"
	"gorm.io/gorm"
)

type RolectrlRole struct {
	Id        string                    `gorm:"primarykey;size:64"`
	Cans      map[miknas.AuthResId]bool `gorm:"serializer:json"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func PackRoleInfo(roleRec *RolectrlRole) miknas.H {
	return miknas.H{
		"role": roleRec.Id,
		"cans": roleRec.Cans,
	}
}

func GetRoleById(db *gorm.DB, role string) *RolectrlRole {
	var roleRec RolectrlRole
	err := db.First(&roleRec, "id = ?", role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		panic(miknas.NewFailRet(err.Error()))
	}
	return &roleRec
}

func AddOneRole(db *gorm.DB, role string) (*RolectrlRole, error) {
	roleRec := RolectrlRole{
		Id:   role,
		Cans: map[miknas.AuthResId]bool{},
	}
	if err := db.Create(&roleRec).Error; err != nil {
		return nil, err
	}
	return &roleRec, nil
}

func GetAllRoleIds(db *gorm.DB) []string {
	var roleRecs []RolectrlRole
	db.Select("Id").Find(&roleRecs)
	roles := []string{}
	for _, roleRec := range roleRecs {
		roles = append(roles, roleRec.Id)
	}
	return roles
}
