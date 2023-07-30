package mikauth

import (
	"github.com/mikzone/miknas/server/miknas"
)

type MyUserAuth struct {
	role string
	ch   *miknas.ContextHelper
	rec  *MikauthRole
}

func (m *MyUserAuth) GetUid() string {
	session := m.ch.GetSession()
	return miknas.AnyToStr(session.Get("uid"))
}

func (m *MyUserAuth) MustGetUid() string {
	uid := m.GetUid()
	if uid == "" {
		panic(miknas.NewFailRet("you have not logined"))
	}
	return uid
}

func (m *MyUserAuth) Refresh() {
	session := m.ch.GetSession()
	uid := miknas.AnyToStr(session.Get("uid"))
	if uid == "" {
		return
	}
	user := GetUserByUid(m.ch.GetApp().Db, uid)
	if user == nil {
		return
	}
	role := user.Role
	adminUid := m.ch.GetApp().ConfMgr.Get("MIKNAS_ADMIN_UID").(string)
	if uid == adminUid {
		role = "admin"
	}
	preRole := miknas.AnyToStr(session.Get("role"))
	if role != preRole {
		session.Set("role", role)
		session.Save()
	}
}

func (m *MyUserAuth) CanAccess(resid miknas.AuthResId) bool {
	// 没有定义或者找不到角色的都用默认值
	app := m.ch.GetApp()
	item := app.AuthMgr.GetItem(resid)
	if item == nil {
		return false
	}
	if m.role == "admin" {
		return true
	}
	db := app.Db
	roleRec := m.rec
	if roleRec == nil {
		// 记录一次缓存
		roleRec = GetRoleById(db, m.role)
		if roleRec == nil {
			roleRec = &MikauthRole{}
		}
		m.rec = roleRec
	}
	if roleRec != nil {
		flag, exist := roleRec.Cans[resid]
		if exist {
			return flag
		}
	}
	// 没有定义或者找不到角色的都用默认值
	return item.Default || false
}

func (ext *MikNasExt) GetUserAuth(ch *miknas.ContextHelper) miknas.IUserAuth {
	session := ch.GetSession()
	role := miknas.AnyToStr(session.Get("role"))
	if role == "" {
		role = "tour"
	}
	return &MyUserAuth{
		role: role,
		ch:   ch,
	}
}
