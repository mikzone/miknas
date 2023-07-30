package devtool

import (
	"github.com/mikzone/miknas/server/miknas"
)

func downDb(ch *miknas.ContextHelper) {
	ext := ch.GetRelExt()
	ch.Ensure(ext.Res("vist"))
	dbpath := ch.GetApp().ConfMgr.Get("MIKNAS_DATABASE_PATH").(string)
	ch.Ctx.Writer.Header().Set("Cache-Control", "no-store")
	ch.Ctx.File(dbpath)
}

func viewTable(ch *miknas.ContextHelper) {
	ext := ch.GetRelExt()
	ch.Ensure(ext.Res("vist"))
	c := ch.Ctx
	tablename := c.Param("tablename")
	if len(tablename) <= 0 {
		ch.FailResp("非法table名: %s", tablename)
		return
	}
	db := ch.GetApp().Db
	var results []map[string]interface{}
	err := db.Table(tablename).Find(&results).Error
	ch.EnsureNoErr(err)
	ch.SucResp(results)
}

func descTable(ch *miknas.ContextHelper) {
	ext := ch.GetRelExt()
	ch.Ensure(ext.Res("vist"))
	c := ch.Ctx
	tablename := c.Param("tablename")
	if len(tablename) <= 0 {
		ch.FailResp("非法table名: %s", tablename)
		return
	}
	db := ch.GetApp().Db
	tableNames := []string{}
	err1 := db.Raw("SELECT name FROM sqlite_master WHERE type='table'").Scan(&tableNames).Error
	ch.EnsureNoErr(err1)
	hasTable := false
	for _, name := range tableNames {
		if name == tablename {
			hasTable = true
			break
		}
	}
	if !hasTable {
		ch.FailResp("不存在的表名: %s", tablename)
		return
	}
	var results []map[string]interface{}
	err := db.Raw("PRAGMA table_info(" + tablename + ")").Scan(&results).Error
	ch.EnsureNoErr(err)
	ch.SucResp(results)
}

func regRoutes(ext *MikNasExt) {
	ext.GET("/miknas.db", downDb)
	ext.GET("/viewTable/:tablename", viewTable)
	ext.GET("/descTable/:tablename", descTable)
}
