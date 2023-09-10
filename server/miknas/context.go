package miknas

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"path"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
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
	status := ch.Ctx.GetInt(ctxFailStatusKey)
	if status == 0 {
		status = http.StatusOK
	}
	ch.Ctx.JSON(status, gin.H{
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

const ctxFailStatusKey = "miknasFailStatus"

func (ch *ContextHelper) SetFailStatus(status int) {
	ch.Ctx.Set(ctxFailStatusKey, status)
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
	if err != nil {
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
	return sessions.DefaultMany(ch.Ctx, ctxSessionNameUserKey)
}

func (ch *ContextHelper) GetFixSession() sessions.Session {
	return sessions.DefaultMany(ch.Ctx, ctxSessionNameFixKey)
}

func (ch *ContextHelper) GetSessionId() string {
	// 获取sessionid
	session := ch.GetFixSession()
	sid := session.Get("Sid")
	if sid == nil {
		guid := xid.New()
		newSid := guid.String()
		session.Set("Sid", newSid)
		sid = session.Get("Sid")
		session.Save()
	}
	return AnyToStr(sid)
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

func (ch *ContextHelper) OpenFs(fsid string, mode string) IFsDriver {
	app := ch.GetApp()
	fstype, fssubid := DecomposeFsid(fsid)
	filespace, exist := app.FileSpaces[fstype]
	if !exist {
		panic(NewFailRet("Fstype(%s) not found", fstype))
	}
	filespace.Ensure(ch, mode)
	return filespace.NewFsDriver(ch, fssubid)
}

func (ch *ContextHelper) UsLog(msg string, args ...any) {
	ch.GetRelExt().GetLogger("us").Info(msg, args...)
}

func (ch *ContextHelper) Logger() *slog.Logger {
	return ch.GetRelExt().Logger()
}

func MakeCtxHelper(c *gin.Context) *ContextHelper {
	ch := &ContextHelper{c}
	return ch
}
