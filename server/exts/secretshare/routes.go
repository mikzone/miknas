package secretshare

import (
	"time"

	"github.com/mikzone/miknas/server/miknas"
)

type inDataAddSs struct {
	Name string `json:"name" binding:"required"`
	Txt  string `json:"txt" binding:"required"`
	Hint string `json:"hint"`
	Intv int    `json:"intv" binding:"required"`
}
type inDataModifySs struct {
	inDataAddSs
	Mid string `json:"mid" binding:"required"`
}

type inDataMid struct {
	Mid string `json:"mid" binding:"required"`
}

func addSecret(ch *miknas.ContextHelper) {
	var loc inDataAddSs
	ch.BindJSON(&loc)
	userauth := ch.GetUserAuth()
	uid := userauth.MustGetUid()
	db := ch.GetApp().Db
	mid, err := GenMid(db)
	if err != nil {
		ch.FailResp(err.Error())
		return
	}
	bm := &SecretShareItem{
		Uid:      uid,
		Mid:      mid,
		Name:     loc.Name,
		Txt:      loc.Txt,
		Hint:     loc.Hint,
		ExpireAt: time.Now().Add(time.Duration(loc.Intv) * time.Second),
	}
	result := db.Create(bm)
	ch.EnsureNoErr(result.Error)
	info := PackSecretInfo(bm)
	ch.UsLog("addSecret", info)
	ch.SucResp(info)
}

func modifySecret(ch *miknas.ContextHelper) {
	var loc inDataModifySs
	ch.BindJSON(&loc)
	db := ch.GetApp().Db
	bm := GetSecretByMid(db, loc.Mid)
	if bm == nil {
		ch.FailResp("密文不存在了")
		return
	}
	userauth := ch.GetUserAuth()
	uid := userauth.MustGetUid()
	if bm.Uid != uid {
		ch.FailResp("不是你的密文")
		return
	}
	modifyInfo := map[string]interface{}{
		"name":      loc.Name,
		"txt":       loc.Txt,
		"hint":      loc.Hint,
		"expire_at": time.Now().Add(time.Duration(loc.Intv) * time.Second),
	}
	db.Model(bm).Updates(modifyInfo)
	info := PackSecretInfo(bm)
	ch.UsLog("modifySecret", info)
	ch.SucResp(info)
}

func viewOne(ch *miknas.ContextHelper) {
	var loc inDataMid
	ch.BindJSON(&loc)
	db := ch.GetApp().Db
	bm := GetSecretByMid(db, loc.Mid)
	if bm == nil {
		ch.FailResp("密文不存在了")
		return
	}
	userauth := ch.GetUserAuth()
	uid := userauth.MustGetUid()
	if bm.Uid != uid {
		now := time.Now()
		if bm.ExpireAt.Before(now) {
			ch.FailResp("密文已过期")
			return
		}
		return
	}
	info := PackSecretInfo(bm)
	ch.SucResp(info)
}

func removeSecret(ch *miknas.ContextHelper) {
	var loc inDataMid
	ch.BindJSON(&loc)
	db := ch.GetApp().Db
	bm := GetSecretByMid(db, loc.Mid)
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
	info := PackSecretInfo(bm)
	result := db.Delete(&bm)
	if result.Error != nil {
		ch.FailResp(result.Error.Error())
		return
	}
	ch.UsLog("removeSecret", info)
	ch.SucResp(result.RowsAffected)
}

func querySecrets(ch *miknas.ContextHelper) {
	userauth := ch.GetUserAuth()
	uid := userauth.MustGetUid()
	db := ch.GetApp().Db
	var bms []SecretShareItem
	db.Where("uid = ?", uid).Find(&bms)
	infos := []miknas.H{}
	for _, bm := range bms {
		infos = append(infos, PackSecretInfo(&bm))
	}
	ch.SucResp(infos)
}

func regRoutes(ext *MikNasExt) {
	ext.POST("/addSecret", addSecret)
	ext.POST("/modifySecret", modifySecret)
	ext.POST("/viewOne", viewOne)
	ext.POST("/querySecrets", querySecrets)
	ext.POST("/removeSecret", removeSecret)
}
