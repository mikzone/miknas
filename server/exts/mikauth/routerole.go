package mikauth

import "github.com/mikzone/miknas/server/miknas"

func allRoles(ch *miknas.ContextHelper) {
	ch.Ensure(ch.GetRelExt().Res("manager"))
	db := ch.GetApp().Db
	var roles []MikauthRole
	db.Find(&roles)
	infos := []miknas.H{}
	for _, role := range roles {
		infos = append(infos, PackRoleInfo(&role))
	}

	ch.SucResp(miknas.H{
		"auths": ch.GetApp().AuthMgr.PackAllAuthInfo(),
		"roles": infos,
	})
}

type inDataRole struct {
	Role string `json:"role" binding:"required"`
}

func oneRole(ch *miknas.ContextHelper) {
	ch.Ensure(ch.GetRelExt().Res("manager"))
	var modify inDataRole
	ch.BindJSON(&modify)
	db := ch.GetApp().Db
	roleRec := GetRoleById(db, modify.Role)
	if roleRec == nil {
		ch.FailResp("角色(%s)不存在", modify.Role)
		return
	}
	ch.SucResp(miknas.H{
		"auths":    ch.GetApp().AuthMgr.PackAllAuthInfo(),
		"roleInfo": PackRoleInfo(roleRec),
	})
}

func addRole(ch *miknas.ContextHelper) {
	ch.Ensure(ch.GetRelExt().Res("manager"))
	var modify inDataRole
	ch.BindJSON(&modify)
	db := ch.GetApp().Db
	roleRec, err := AddOneRole(db, modify.Role)
	if err != nil {
		ch.FailResp(err.Error())
		return
	}
	ch.SucResp(miknas.H{
		"auths":    ch.GetApp().AuthMgr.PackAllAuthInfo(),
		"roleInfo": PackRoleInfo(roleRec),
	})
}

func removeRole(ch *miknas.ContextHelper) {
	ch.Ensure(ch.GetRelExt().Res("manager"))
	var modify inDataRole
	ch.BindJSON(&modify)
	db := ch.GetApp().Db
	roleRec := GetRoleById(db, modify.Role)
	if roleRec == nil {
		ch.FailResp("角色(%s)不存在", modify.Role)
		return
	}
	result := db.Delete(&roleRec)
	if result.Error != nil {
		ch.FailResp(result.Error.Error())
		return
	}
	ch.SucResp("删除成功")
}

type inDataSaveRole struct {
	Role string                    `json:"role" binding:"required"`
	Cans map[miknas.AuthResId]bool `json:"cans" binding:"required"`
}

func saveRole(ch *miknas.ContextHelper) {
	ch.Ensure(ch.GetRelExt().Res("manager"))
	var modify inDataSaveRole
	ch.BindJSON(&modify)
	db := ch.GetApp().Db
	roleRec := GetRoleById(db, modify.Role)
	if roleRec == nil {
		ch.FailResp("角色(%s)不存在", modify.Role)
		return
	}
	saveCans := ch.GetApp().AuthMgr.FilterDict(modify.Cans)
	roleRec.Cans = saveCans
	result := db.Save(&roleRec)
	if result.Error != nil {
		ch.FailResp(result.Error.Error())
		return
	}
	ch.SucResp("删除成功")
}

func regRoleRoutes(ext *MikNasExt) {
	ext.POST("allRoles", allRoles)
	ext.POST("oneRole", oneRole)
	ext.POST("addRole", addRole)
	ext.POST("removeRole", removeRole)
	ext.POST("saveRole", saveRole)
}
