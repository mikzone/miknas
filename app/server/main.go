package main

import (
	"embed"
	"io/fs"
	"net/http"

	// _ "net/http/pprof"

	"github.com/mikzone/miknas/server/exts/bookmarks"
	"github.com/mikzone/miknas/server/exts/cmdexec"
	"github.com/mikzone/miknas/server/exts/devtool"

	// "github.com/mikzone/miknas/server/exts/extensionexample"
	"github.com/mikzone/miknas/server/exts/drive"
	"github.com/mikzone/miknas/server/exts/mikauth"
	"github.com/mikzone/miknas/server/exts/note"
	"github.com/mikzone/miknas/server/exts/official"
	"github.com/mikzone/miknas/server/exts/pan"
	"github.com/mikzone/miknas/server/exts/secretshare"
	"github.com/mikzone/miknas/server/miknas"
)

//go:embed client
var clientf embed.FS

func main() {
	// go func() {
	// 	http.ListenAndServe("0.0.0.0:8899", nil)
	// }()

	// 初始化app
	app := miknas.NewApp()
	app.AddExt(official.New())
	app.AddExt(mikauth.New())

	// 自行添加扩展
	app.AddExt(drive.New())
	app.AddExt(pan.New())
	// app.AddExt(extensionexample.New())
	app.AddExt(bookmarks.New())
	app.AddExt(secretshare.New())
	app.AddExt(cmdexec.New())
	app.AddExt(devtool.New())
	app.AddExt(note.New())
	app.ConfMgr.UpdateFromEnv()
	app.ConfMgr.PrintConfigs()
	app.AuthMgr.PrintAllAuthInfo()
	app.StartInit()

	// 使用embed来打包client，client目录需要提前编译好复制过来
	webfs, err := fs.Sub(clientf, "client")
	if err != nil {
		panic(err)
	}

	engine := app.Engine
	// 指定未识别的路由都交由前端处理
	engine.NoRoute(miknas.VueHandler(http.FS(webfs)))

	// 在指定端口上运行
	engine.Run(":2020")
}
