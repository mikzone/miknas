package mikauth

import (
	"github.com/mikzone/miknas/server/exts/rolectrl"
	"github.com/mikzone/miknas/server/miknas"
)

type MyUserAuth struct {
	role string
	ch   *miknas.ContextHelper
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

func (m *MyUserAuth) HasRole(role string) bool {
	return m.role == role
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
	m.role = role
	preRole := miknas.AnyToStr(session.Get("role"))
	if role != preRole {
		session.Set("role", role)
		session.Save()
	}
}

func (m *MyUserAuth) GetExtra() map[string]any {
	return map[string]any{
		"role": m.role,
		"ch":   m.ch,
	}
}

func (m *MyUserAuth) CanAccess(resid miknas.AuthResId) bool {
	return rolectrl.CanAccess(m, resid)
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
