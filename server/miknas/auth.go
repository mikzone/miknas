package miknas

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// 权限资源Id的类型，建议用Extension.Res生成
type AuthResId string

// 所有需要替代默认权限管理的都需要实现该接口
type IUserAuth interface {
	GetUid() string
	MustGetUid() string
	CanAccess(AuthResId) bool
	Refresh()
}

/*
AuthResItem need to be declare before use.
We assume that all config value are jsonable,
because sometime we need to send to client!
*/
type AuthResItem struct {
	// 资源id,统一格式为"extid/resouce[@w][:action]",
	// 中括号里面是可选的,带有@w表示是白名单权限,默认值为true
	ResId      AuthResId
	Desc       string
	SendClient bool
	// indicate who register it
	ExtId string
	// default的值其实是看资源里面是否带有@w，这里是缓存起来避免再次计算
	Default bool
}

// 全局的 权限资源 管理器，只管资源的注册，不管权限认证
type AuthResMgr struct {
	items    map[AuthResId]AuthResItem
	sendlist []AuthResId
}

func (m *AuthResMgr) PackClientDict(ua IUserAuth) gin.H {
	ret := gin.H{}
	for _, resid := range m.sendlist {
		item := m.items[resid]
		flag := ua.CanAccess(resid)
		// 只传输不是默认值的
		if flag != item.Default {
			ret[string(resid)] = flag
		}
	}
	return ret
}

func (m *AuthResMgr) FilterDict(authDict map[AuthResId]bool) map[AuthResId]bool {
	ret := map[AuthResId]bool{}
	for resid, flag := range authDict {
		item := m.items[resid]
		// 只保留不是默认值的
		if flag != item.Default {
			ret[resid] = flag
		}
	}
	return ret
}

func (m *AuthResMgr) RegAuthItem(item AuthResItem) error {
	resid := item.ResId
	extid := item.ExtId
	if len(extid) <= 0 {
		return fmt.Errorf("extid cannot be empty, when register item(%v)", item)
	}
	if len(resid) <= 0 {
		return fmt.Errorf("resid cannot be empty, registed by extid(%s)", item.ExtId)
	}
	if preItem, ok := m.items[resid]; ok {
		return fmt.Errorf("authres %s existed, registed by extid(%s)", resid, preItem.ExtId)
	}

	// 限定一个权限资源id必须以其扩展名开头
	resIdStr := string(resid)
	needPrefix := extid + "/"
	if !strings.HasPrefix(resIdStr, needPrefix) {
		return fmt.Errorf("resid(%s) must start with '{its extid}/'", resid)
	}

	// 获取Default值
	strl := strings.SplitN(resIdStr, ":", 2)
	defv := false
	if strings.HasSuffix(strl[0], "@w") {
		defv = true
	}
	item.Default = defv

	m.items[resid] = item
	if item.SendClient {
		m.sendlist = append(m.sendlist, resid)
	}
	return nil
}

func (m *AuthResMgr) HasRes(resid AuthResId) bool {
	_, ok := m.items[resid]
	return ok
}

func (m *AuthResMgr) GetItem(resid AuthResId) *AuthResItem {
	item, ok := m.items[resid]
	if ok {
		return &item
	}
	return nil
}

func (m *AuthResMgr) PackAllAuthInfo() gin.H {
	ret := gin.H{}
	for resid, item := range m.items {
		ret[string(resid)] = gin.H{
			"resource":   item.ResId,
			"desc":       item.Desc,
			"sendClient": item.SendClient,
		}
	}
	return ret
}

func (m *AuthResMgr) PrintAllAuthInfo() {
	str, err := json.MarshalIndent(m.items, "", "  ")
	if err != nil {
		fmt.Printf("PrintAllAuthInfo Error in json marshal: %v", err)
		return
	}
	fmt.Println("[MikNas]AuthInfos:", string(str))
}

func NewAuthResMgr() *AuthResMgr {
	return &AuthResMgr{
		map[AuthResId]AuthResItem{},
		[]AuthResId{},
	}
}
