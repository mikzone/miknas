package drive

import (
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/mikzone/miknas/server/miknas"
	"github.com/rs/xid"
)

type ShareFileSpace struct {
	miknas.SimpleFileSpace
}

func (fsp *ShareFileSpace) Ensure(ch *miknas.ContextHelper, mode string) {
	if mode != "r" {
		panic(miknas.NewFailRet("分享的文件没有%s权限", mode))
	}
}

func (fsp *ShareFileSpace) NewFsDriver(ch *miknas.ContextHelper, fssubid string) miknas.IFsDriver {
	app := ch.GetApp()
	db := app.Db
	item := GetShareItemByMid(db, fssubid)
	if item == nil {
		panic(miknas.NewFailRet("分享不存在或已被删除"))
	}
	if IsShareItemExpired(item) {
		panic(miknas.NewFailRet("分享已过期"))
	}
	if len(item.Pwd) > 0 && !hadVerifyShare(ch, item) {
		panic(miknas.NewFailRet("请验证提取码再操作"))
	}
	filespace, exist := app.FileSpaces[item.Fstype]
	if !exist {
		panic(miknas.NewFailRet("无法识别的分享源"))
	}
	return filespace.NewFsDriverByAddr(item.Fsaddr)
}

func RegShareFileSpace(ext *MikNasExt) {
	ShareFileSpace := &ShareFileSpace{
		*miknas.NewSimpleFileSpace("S", "", ext.Res("nouse")),
	}
	ext.RegFileSpace(ShareFileSpace)
}

func markVerifyShare(ch *miknas.ContextHelper, item *FileShareItem) {
	sessionid := ch.GetSessionId()
	key := sessionid + "_" + item.Sid
	mVistShareCache.Add(key, true)
}

func hadVerifyShare(ch *miknas.ContextHelper, item *FileShareItem) bool {
	if item == nil {
		return false
	}
	sessionid := ch.GetSessionId()
	key := sessionid + "_" + item.Sid
	canVist, ok := mVistShareCache.Get(key)
	if ok {
		return canVist
	}
	return false
}

type tmpDownInfo struct {
	Uid    string
	Fstype string
	Fsaddr string
	Name   string
	Ts     time.Time
}

func GenTmpDownId() (string, error) {
	for i := 0; i < 5; i++ {
		guid := xid.New()
		mid := guid.String()
		if !mTmpDownRecCache.Contains(mid) {
			return mid, nil
		}
	}
	return "", miknas.NewFailRet("生成下载码失败过多")
}

var mVistShareCache *expirable.LRU[string, bool]
var mTmpDownRecCache *expirable.LRU[string, tmpDownInfo]

func init() {
	mVistShareCache = expirable.NewLRU[string, bool](50000, nil, time.Hour*1)
	mTmpDownRecCache = expirable.NewLRU[string, tmpDownInfo](50000, nil, time.Minute*2)
}
