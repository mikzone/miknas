package miknas

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type App struct {
	svrRouter      *gin.RouterGroup
	extids         []string
	exts           map[string]IExtension
	ConfMgr        *ConfigManager
	AuthMgr        *AuthResMgr
	sessionHandler gin.HandlerFunc
	loggers        map[string]*slog.Logger
	Db             *gorm.DB
	WorkSpace      IDiskFsDriver
	FileSpaces     map[string]IFileSpace
}

const ctxAppKey = "miknasApp"
const ctxExtKey = "miknasExt"

const ctxSessionNameUserKey = "u"
const ctxSessionNameFixKey = "f"

// register one extension
func (a *App) AddExt(ext IExtension) {
	extid := ext.GetId()
	if _, ok := a.exts[extid]; ok {
		fmt.Println("AddExt Failed: ExtId (", extid, ") is existed!!!")
		return
	}
	a.exts[extid] = ext
	a.extids = append(a.extids, extid)
	ext.bindApp(a)
	extRouter := ext.GetRouter()
	// 注册一个中间件用来获取extension
	extRouter.Use(func(c *gin.Context) {
		c.Set(ctxExtKey, ext)
		defer func() {
			c.Set(ctxExtKey, nil)
		}()
		c.Next()
	})
	// 注册一个中间件用来管理session
	extRouter.Use(func(c *gin.Context) {
		// 由于配置没有那么快初始化好，所以放在第一次访问时初始化
		if a.sessionHandler == nil {
			secretVal := a.ConfMgr.Get("MIKNAS_SECRET_KEY")
			if secretVal == nil {
				panic("config MIKNAS_SECRET_KEY must be set")
			}
			secret := secretVal.(string)
			store := cookie.NewStore([]byte(secret))
			sessionNames := []string{ctxSessionNameUserKey, ctxSessionNameFixKey}
			// 两个session，一个是跟随user登录增加删除的，一个是一直固定的登录也不会清理的
			a.sessionHandler = sessions.SessionsMany(sessionNames, store)
		}
		a.sessionHandler(c)
	})
	ext.OnBind()
}

func (a *App) GetExtids() []string {
	return a.extids
}

func (a *App) GetLogger(name string) *slog.Logger {
	if len(name) < 1 {
		name = "run_out"
	}
	l, ok := a.loggers[name]
	if !ok {
		filename, needStdOut := strings.CutSuffix(name, "_out")
		l = CreateSlogLogger(a, filename, needStdOut)
		a.loggers[name] = l
	}
	return l
}

func (a *App) GetExt(extid string) IExtension {
	return a.exts[extid]
}

// a middleware for inject cur app to gin context
func (a *App) InjectToCtxMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			c.Set(ctxAppKey, nil)
		}()
		c.Set(ctxAppKey, a)
		c.Next()
	}
}

func (a *App) StartInit() {
	wspath := a.ConfMgr.Get("MIKNAS_WORKSPACE").(string)
	a.WorkSpace = NewBaseFsDriver(wspath, true)
	// 初始化db
	dbpath := a.ConfMgr.Get("MIKNAS_DATABASE_PATH").(string)
	dbDebug := a.ConfMgr.Get("MIKNAS_DATABASE_DEBUG").(string) == "1"
	gormConf := gorm.Config{}
	if dbDebug {
		gormConf.Logger = logger.Default.LogMode(logger.Info)
	}
	db, err := openSqlite3(dbpath, &gormConf)
	if err != nil {
		panic(err)
	}
	a.Db = db
	for _, ext := range a.exts {
		ext.OnInit()
	}
}

// 注册文件空间
func (a *App) RegFileSpace(filespace IFileSpace) error {
	fstype := filespace.GetFstype()
	if fstype == "" {
		return NewFailRet("fstype 不能为空")
	}
	if filespace.GetRelExt() == nil {
		return NewFailRet("文件空间对应的扩展不能为空")
	}
	preFilespace, exist := a.FileSpaces[fstype]
	if exist {
		return NewFailRet("扩展(%s)中的Fstype(%s)在扩展(%s)中已定义", filespace.GetRelExt().GetId(), fstype, preFilespace.GetRelExt().GetId())
	}
	a.FileSpaces[fstype] = filespace
	return nil
}

func NewApp(r *gin.RouterGroup) *App {
	svrRouter := r.Group("/s")
	app := &App{
		svrRouter:  svrRouter,
		extids:     []string{},
		exts:       map[string]IExtension{},
		loggers:    map[string]*slog.Logger{},
		ConfMgr:    NewConfigManager(),
		AuthMgr:    NewAuthResMgr(),
		FileSpaces: map[string]IFileSpace{},
	}
	svrRouter.Use(app.InjectToCtxMiddleware())
	svrRouter.Use(HandleFailRetMiddleware)
	return app
}
