package drive

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/mikzone/miknas/server/miknas"
)

type inDataAddShare struct {
	Fsid   string `json:"fsid" binding:"required"`
	Fspath string `json:"fspath" binding:"required"`
	Pwd    string `json:"pwd"`
	Intv   uint   `json:"intv" binding:"required"`
}

type inDataMid struct {
	Sid string `json:"sid" binding:"required"`
}

type inDataVerify struct {
	Sid string `json:"sid" binding:"required"`
	Pwd string `json:"pwd" binding:"required"`
}

func addShare(ch *miknas.ContextHelper) {
	var loc inDataAddShare
	ch.BindJSON(&loc)
	userauth := ch.GetUserAuth()
	uid := userauth.MustGetUid()
	fspath := loc.Fspath
	fstype, fssubid := miknas.DecomposeFsid(loc.Fsid)
	name := filepath.Base(fspath)
	app := ch.GetApp()
	filespace, exist := app.FileSpaces[fstype]
	if !exist {
		ch.FailResp("无法识别的分享源")
		return
	}
	fsaddr := filespace.GetAddr(filespace, ch, fssubid, fspath)

	db := app.Db
	mid, err := GenMid(db)
	if err != nil {
		ch.FailResp(err.Error())
		return
	}
	item := &FileShareItem{
		PrimaryStringModel: miknas.PrimaryStringModel{
			Sid: mid,
		},
		Uid:    uid,
		Name:   name,
		Pwd:    loc.Pwd,
		Vist:   0,
		Fstype: fstype,
		Fsaddr: fsaddr,
		Intv:   loc.Intv,
	}
	result := db.Create(item)
	ch.EnsureNoErr(result.Error)
	info := PackShareItemInfo(item)
	ch.UsLog("addShare", "item", item)
	ch.SucResp(info)
}

func viewShare(ch *miknas.ContextHelper) {
	var loc inDataMid
	ch.BindJSON(&loc)
	db := ch.GetApp().Db
	item := GetShareItemByMid(db, loc.Sid)
	if item == nil {
		ch.FailResp("该分享不存在了")
		return
	}
	if IsShareItemExpired(item) {
		ch.FailResp("你来晚了，该分享已过期")
		return
	}

	if item.Pwd != "" {
		// 有密码保护的
		if !hadVerifyShare(ch, item) {
			ch.SucResp(miknas.H{
				"needPwd": true,
			})
			return
		}
	}

	info := PackShareItemInfo(item)
	ch.SucResp(info)
}

func verifyShare(ch *miknas.ContextHelper) {
	var loc inDataVerify
	ch.BindJSON(&loc)
	db := ch.GetApp().Db
	item := GetShareItemByMid(db, loc.Sid)
	if item == nil {
		ch.FailResp("该分享不存在了")
		return
	}
	if IsShareItemExpired(item) {
		ch.FailResp("你来晚了，该分享已过期")
		return
	}

	if item.Pwd != loc.Pwd {
		ch.FailResp("提取码不正确")
		return
	}
	markVerifyShare(ch, item)
	info := PackShareItemInfo(item)
	ch.SucResp(info)
}

func removeShare(ch *miknas.ContextHelper) {
	var loc inDataMid
	ch.BindJSON(&loc)
	db := ch.GetApp().Db
	item := GetShareItemByMid(db, loc.Sid)
	if item == nil {
		ch.FailResp("分享不存在了")
		return
	}
	userauth := ch.GetUserAuth()
	uid := userauth.MustGetUid()
	if item.Uid != uid {
		ch.FailResp("不是你的分享")
		return
	}
	result := db.Delete(&item)
	if result.Error != nil {
		ch.FailResp(result.Error.Error())
		return
	}
	ch.UsLog("removeShare", "item", item)
	ch.SucResp(result.RowsAffected)
}

func queryShares(ch *miknas.ContextHelper) {
	userauth := ch.GetUserAuth()
	uid := userauth.MustGetUid()
	db := ch.GetApp().Db
	var bms []FileShareItem
	db.Where("uid = ?", uid).Find(&bms)
	infos := []miknas.H{}
	for _, item := range bms {
		infos = append(infos, PackShareItemInfo(&item))
	}
	ch.SucResp(infos)
}

func genTmpDownUrl(ch *miknas.ContextHelper) {
	var loc inDataFsLocate
	ch.BindJSON(&loc)
	fsid := loc.Fsid
	fspath := loc.Fspath
	fsd := ch.OpenFs(fsid, "w")
	EnsureIsExistedFile(fsd, fspath)
	fstype, fssubid := miknas.DecomposeFsid(loc.Fsid)
	name := filepath.Base(fspath)
	app := ch.GetApp()
	filespace, exist := app.FileSpaces[fstype]
	if !exist {
		ch.FailResp("无法识别的分享源")
		return
	}
	fsaddr := filespace.GetAddr(filespace, ch, fssubid, fspath)

	mid, err := GenTmpDownId()
	if err != nil {
		ch.FailResp(err.Error())
		return
	}

	item := tmpDownInfo{
		Fstype: fstype,
		Name:   name,
		Fsaddr: fsaddr,
		Ts:     time.Now(),
	}

	mTmpDownRecCache.Add(mid, item)

	urlstr := fmt.Sprintf("tmpDown/%s/%s", mid, name)
	ch.SucResp(urlstr)
}

func tmpDown(ch *miknas.ContextHelper) {
	ch.SetFailStatus(400)
	c := ch.Ctx
	downid := c.Param("downid")
	name := c.Param("name")
	item, ok := mTmpDownRecCache.Peek(downid)
	if !ok {
		ch.FailResp("链接不存在或已过期")
		return
	}
	if item.Name != name {
		ch.FailResp("链接错误")
		return
	}
	app := ch.GetApp()
	fstype := item.Fstype
	filespace, exist := app.FileSpaces[fstype]
	if !exist {
		ch.FailResp("无法识别的分享源")
		return
	}
	fsd := filespace.NewFsDriverByAddr(item.Fsaddr)
	fsd.WebDownload(fsd, ch, name)
}

func regFileShareRoutes(ext *MikNasExt) {
	ext.POST("/addShare", addShare)
	ext.POST("/viewShare", viewShare)
	ext.POST("/verifyShare", verifyShare)
	ext.POST("/queryShares", queryShares)
	ext.POST("/removeShare", removeShare)
	ext.POST("/genTmpDownUrl", genTmpDownUrl)
	ext.GET("/tmpDown/:downid/:name", tmpDown)
}
