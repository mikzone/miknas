package miknas

import (
	"fmt"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(*ContextHelper)

type IExtension interface {
	bindApp(*App)
	GetId() string
	GetRouter() *gin.RouterGroup
	GetUserAuth(*ContextHelper) IUserAuth
	ClientUrl(string) string
	// 注册权限(资源名称，权限描述，是否发送给客户端)
	RegAuth(AuthResId, string, bool) error
	OnBind()
	OnInit()
	Res(string) AuthResId
}

type Extension struct {
	ExtId  string
	App    *App
	Router *gin.RouterGroup
}

// bind to app
func (r *Extension) bindApp(app *App) {
	if r.ExtId == "" {
		panic("Extension ExtId must be init")
	}
	r.App = app
	r.Router = app.svrRouter.Group("/" + r.ExtId)
}

func (r *Extension) GetId() string {
	return r.ExtId
}

func (r *Extension) GetRouter() *gin.RouterGroup {
	return r.Router
}

// 如果是开发替代默认的权限校验功能的扩展必须实现这个，用户修改 MIKNAS_AUTH_EXTS 的配置即可
func (r *Extension) GetUserAuth(ch *ContextHelper) IUserAuth {
	return nil
}

func (r *Extension) Res(resource string) AuthResId {
	return AuthResId(fmt.Sprintf("%s/%s", r.ExtId, resource))
}

func (r *Extension) RegConf(item ConfItem) error {
	item.ExtId = r.ExtId
	err := r.App.ConfMgr.RegConfItem(item)
	if err != nil {
		// fmt.Printf("[%s]RegConf Fail: %v", item.ExtId, err)
		panic(fmt.Errorf("[%s]RegConf Fail: %v", item.ExtId, err))
	}
	return err
}

func (r *Extension) RegStrConf(key string, defv any, desc string, sendClient bool) error {
	item := ConfItem{Key: key, Default: defv, Desc: desc, SendClient: sendClient, CheckConv: CheckConvStr}
	return r.RegConf(item)
}

func (r *Extension) RegIntConf(key string, defv any, desc string, sendClient bool) error {
	item := ConfItem{Key: key, Default: defv, Desc: desc, SendClient: sendClient, CheckConv: CheckConvInt}
	return r.RegConf(item)
}

func (r *Extension) RegMapConf(key string, defv any, desc string, sendClient bool) error {
	item := ConfItem{Key: key, Default: defv, Desc: desc, SendClient: sendClient, CheckConv: CheckConvMap}
	return r.RegConf(item)
}

func (r *Extension) RegAuth(resid AuthResId, desc string, sendClient bool) error {
	item := AuthResItem{
		ResId:      resid,
		Desc:       desc,
		SendClient: sendClient,
	}
	item.ExtId = r.ExtId
	err := r.App.AuthMgr.RegAuthItem(item)
	if err != nil {
		// fmt.Printf("[%s]RegAuth Fail: %v", item.ExtId, err)
		panic(fmt.Errorf("[%s]RegAuth Fail: %v", item.ExtId, err))
	}
	return err
}

func (r *Extension) RegFileSpace(filespace IFileSpace) {
	filespace.SetRelExt(r)
	err := r.App.RegFileSpace(filespace)
	if err != nil {
		panic(fmt.Sprintf("[%s]RegFileSpace Fail: %v", filespace, err))
	}
}

func wrapToGinHanlders(handlers []HandlerFunc) gin.HandlersChain {
	ret := gin.HandlersChain{}
	for _, handler := range handlers {
		handler2 := handler // range 作用域原因，这里做个副本
		tmp := func(c *gin.Context) {
			ch := MakeCtxHelper(c)
			handler2(ch)
		}
		ret = append(ret, tmp)
	}
	return ret
}

func (r *Extension) ClientUrl(suburl string) string {
	ConfMgr := r.App.ConfMgr
	rootUrl := ConfMgr.Get("MIKNAS_CLIENT_PREFIX").(string)
	clientMap := ConfMgr.Get("MIKNAS_CLIENT_URL_MAP").(map[string]any)
	extMapUrl, ok := clientMap[r.ExtId]
	if !ok {
		extMapUrl = r.ExtId
	}
	return path.Join(rootUrl, extMapUrl.(string), suburl)
}

func (r *Extension) ServerUrl(suburl string) string {
	return path.Join(r.Router.BasePath(), suburl)
}

func (r *Extension) POST(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	return r.Router.Handle(http.MethodPost, relativePath, wrapToGinHanlders(handlers)...)
}

func (r *Extension) GET(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	return r.Router.Handle(http.MethodGet, relativePath, wrapToGinHanlders(handlers)...)
}

// 定义一些事件接口

// 在App注册扩展的时候
func (r *Extension) OnBind() {}

// 在加载完配置等之后，App初始化完后
func (r *Extension) OnInit() {}

var _ IExtension = (*Extension)(nil)

func NewExtension(ExtId string) Extension {
	return Extension{
		ExtId: ExtId,
	}
}
