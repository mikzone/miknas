package official

import (
	"github.com/gin-gonic/gin"
	"github.com/mikzone/miknas/server/miknas"
)

func getClientInitInfo(ch *miknas.ContextHelper) {
	app := ch.GetApp()
	userAuth := ch.GetUserAuth()
	uid := userAuth.GetUid()
	userAuth.Refresh()
	needids := []string{}
	extids := app.GetExtids()
	for _, tmpid := range extids {
		ext := app.GetExt(tmpid)
		resid := ext.Res("vist")
		if userAuth.CanAccess(resid) {
			needids = append(needids, tmpid)
		}
	}
	ch.SucResp(gin.H{
		"uid":           uid,
		"serverConfigs": app.ConfMgr.PackClientDict(),
		"userAuths":     app.AuthMgr.PackClientDict(userAuth),
		"extids":        needids,
	})
}

func regRoutes(ext *MikNasExt) {
	ext.POST("/getClientInitInfo", getClientInitInfo)
}
