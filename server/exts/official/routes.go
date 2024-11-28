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
	ch.SucResp(gin.H{
		"uid":           uid,
		"serverConfigs": app.ConfMgr.PackClientDict(),
		"userAuths":     app.AuthMgr.PackClientDict(userAuth),
	})
}

func regRoutes(ext *MikNasExt) {
	ext.POST("/getClientInitInfo", getClientInitInfo)
}
