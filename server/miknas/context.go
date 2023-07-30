package miknas

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type ContextHelper struct {
	Ctx *gin.Context
}

//go:embed tpl
var tplFs embed.FS

// 需要注意的是template使用文件名(非完整路径)来作为模版名称,所以文件名相同的会有覆盖的情况
var tplTemplate = template.Must(template.ParseFS(tplFs, "tpl/*"))

// if err happen, stop curent request and return err massage to client
func (ch *ContextHelper) EnsureNoErr(err error) {
	if err != nil {
		p, ok := err.(IFailRet)
		if ok {
			panic(p)
		} else {
			panic(NewFailRet(err.Error()))
		}
	}
}

func (ch *ContextHelper) SucResp(ret any) {
	ch.Ctx.JSON(http.StatusOK, gin.H{
		"suc": true,
		"ret": ret,
	})
}

func (ch *ContextHelper) FailResp(format string, a ...any) {
	why := fmt.Sprintf(format, a...)
	ch.Ctx.JSON(http.StatusOK, gin.H{
		"suc": false,
		"why": why,
	})
}

func (ch *ContextHelper) FailRespWithStatus(status int, format string, a ...any) {
	why := fmt.Sprintf(format, a...)
	ch.Ctx.JSON(status, gin.H{
		"suc": false,
		"why": why,
	})
}

const ctxFailRetStatusKey = "miknasFailretStatus"

func (ch *ContextHelper) SetFailRetStatus(status int) {
	ch.Ctx.Set(ctxFailRetStatusKey, status)
}

func (ch *ContextHelper) Redirect(location string) {
	ch.Ctx.Redirect(http.StatusFound, location)
}

func (ch *ContextHelper) Jump(message, jumpurl string, cd int, flag bool) {
	obj := H{
		"message": message,
		"jumpurl": jumpurl,
		"cd":      cd,
		"flag":    flag,
	}
	// fmt.Printf("tplTemplate %v\n", tplTemplate.DefinedTemplates())
	err := tplTemplate.ExecuteTemplate(ch.Ctx.Writer, "jump.html", H{
		"SERVER_TEMPLATE_DATA": obj,
	})
	if err == nil {
		fmt.Printf("Jump Render err: %v", err)
	}
}

func (ch *ContextHelper) GetApp() *App {
	return ch.Ctx.MustGet(ctxAppKey).(*App)
}

func (ch *ContextHelper) GetRelExt() IExtension {
	return ch.Ctx.MustGet(ctxExtKey).(IExtension)
}

func (ch *ContextHelper) BindJSON(obj any) {
	err := ch.Ctx.ShouldBindJSON(obj)
	if err != nil {
		panic(NewFailRet(err.Error()))
	}
}

func (ch *ContextHelper) MustBind(obj any) {
	err := ch.Ctx.ShouldBind(obj)
	if err != nil {
		panic(NewFailRet(err.Error()))
	}
}

func (ch *ContextHelper) GetUserAuth() IUserAuth {
	app := ch.GetApp()
	extid := app.ConfMgr.Get("MIKNAS_AUTH_EXTS").(string)
	ext := app.GetExt(extid)
	return ext.GetUserAuth(ch)
}

func (ch *ContextHelper) GetSession() sessions.Session {
	return sessions.Default(ch.Ctx)
}

func (ch *ContextHelper) ClientUrl(suburl string) string {
	rootUrl := ch.GetApp().ConfMgr.Get("MIKNAS_CLIENT_PREFIX").(string)
	return path.Join(rootUrl, suburl)
}

func (ch *ContextHelper) Ensure(resid AuthResId) {
	app := ch.GetApp()
	if !app.AuthMgr.HasRes(resid) {
		errMsg := fmt.Sprintf("未定义权限(%s)", resid)
		panic(NewFailRet(errMsg))
	}
	ua := ch.GetUserAuth()
	if !ua.CanAccess(resid) {
		panic(NewFailRet("你的权限不足"))
	}
}

func (ch *ContextHelper) OpenFs(fsid string, mode string) IPathHelper {
	app := ch.GetApp()
	fsType, fsSubid := DecomposeFsid(fsid)
	filespace, exist := app.FileSpaces[fsType]
	if !exist {
		panic(NewFailRet("FsType(%s) not found", fsType))
	}
	filespace.Ensure(ch, mode)
	rootPath := filespace.GetSubRoot(ch, fsSubid)
	if rootPath == "" {
		panic(NewFailRet("rootPath of fsid(%s) cannot be empty", fsid))
	}
	return &BasePathHelper{rootPath: rootPath}
}

func (ch *ContextHelper) UsLog(action string, v H) {
	app := ch.GetApp()
	ext, exists := ch.Ctx.Get(ctxExtKey)
	vBytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	vStr := string(vBytes)
	if exists {
		ext = ext.(IExtension)
		extid := ch.GetRelExt().GetId()
		app.UsLogger.Printf("[miknas|%s]%s: %s", extid, action, vStr)
		return
	}
	app.UsLogger.Printf("[miknas]%s: %s", action, vStr)
}

func MakeCtxHelper(c *gin.Context) *ContextHelper {
	ch := &ContextHelper{c}
	return ch
}
