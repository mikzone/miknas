package rolectrl

import "github.com/mikzone/miknas/server/miknas"

const ctxRoleRecKey = "miknasRoleCtrlRec"

func CanAccess(ua miknas.IUserAuth, resid miknas.AuthResId) bool {
	extra := ua.GetExtra()
	role := extra["role"].(string)
	ch := extra["ch"].(*miknas.ContextHelper)
	app := ch.GetApp()
	item := app.AuthMgr.GetItem(resid)
	if item == nil {
		return false
	}
	if role == "admin" {
		return true
	}
	db := app.Db
	roleRecAny := ch.Ctx.Value(ctxRoleRecKey)
	var roleRec *RolectrlRole
	if roleRecAny == nil {
		// 记录一次缓存
		roleRec = GetRoleById(db, role)
		if roleRec == nil {
			roleRec = &RolectrlRole{}
		}
		ch.Ctx.Set(ctxRoleRecKey, roleRec)
	} else {
		roleRec = roleRecAny.(*RolectrlRole)
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
