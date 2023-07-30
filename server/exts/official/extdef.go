package official

import (
	"github.com/mikzone/miknas/server/miknas"
)

type MikNasExt struct {
	miknas.Extension
}

func New() *MikNasExt {
	return &MikNasExt{miknas.NewExtension("Official")}
}

func (ext *MikNasExt) OnBind() {
	// register config
	ext.RegStrConf("MIKNAS_SECRET_KEY", ";@123~NZR", "用作Session加密的Key", false)
	ext.RegStrConf("MIKNAS_WORKSPACE", "./workspace", "工作空间目录", false)

	ext.RegStrConf("MIKNAS_AUTH_EXTS", "MikAuth", "登录认证使用的扩展", true)
	ext.RegStrConf("MIKNAS_CLIENT_PREFIX", "/", "客户端地址前缀,用来给服务端拼凑客户端url的", true)
	ext.RegMapConf("MIKNAS_CLIENT_URL_MAP",
		map[string]any{
			"Official": "",
		},
		"客户端的URL映射, 把某些模块改名或者把某个扩展作为首页可以使用这个",
		true)
	ext.RegStrConf("MIKNAS_DATABASE_PATH", "miknas.sqlite", "sqlite数据库路径", false)
	ext.RegStrConf("MIKNAS_DATABASE_DEBUG", "0", "是否debug数据库", false)
	ext.RegStrConf("MIKNAS_ADMIN_UID", "admin", "管理员用户uid", false)

	ext.RegAuth(ext.Res("vist"), "使用miknas主页，一般所有人都需要", false)

	// routes
	regRoutes(ext)
}
