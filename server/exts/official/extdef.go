package official

import (
	"crypto/rand"
	"math/big"

	"github.com/mikzone/miknas/server/miknas"
)

type MikNasExt struct {
	miknas.Extension
}

func New() *MikNasExt {
	return &MikNasExt{miknas.NewExtension("Official")}
}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

func (ext *MikNasExt) OnBind() {
	// register config
	initSecretKey, err := GenerateRandomString(8)
	if err != nil {
		initSecretKey = ";@123~NZR"
	}
	ext.RegStrConf("MIKNAS_SECRET_KEY", initSecretKey, "用作Session加密的Key", false)
	ext.RegStrConf("MIKNAS_WORKSPACE", "./workspace", "工作空间目录", false)
	ext.RegStrConf("MIKNAS_CONFIG_DIR", "./config", "配置存储目录", false)

	ext.RegStrConf("MIKNAS_AUTH_EXTS", "MikAuth", "登录认证使用的扩展", true)
	ext.RegStrConf("MIKNAS_CLIENT_PREFIX", "/", "客户端地址前缀,用来给服务端拼凑客户端url的", true)
	ext.RegMapConf("MIKNAS_CLIENT_URL_MAP",
		map[string]any{
			"Official": "",
		},
		"客户端的URL映射, 把某些模块改名或者把某个扩展作为首页可以使用这个",
		true)
	ext.RegStrConf("MIKNAS_DATABASE_PATH", "config/miknas.sqlite", "sqlite数据库路径", false)
	ext.RegStrConf("MIKNAS_DATABASE_DEBUG", "0", "是否debug数据库", false)
	ext.RegStrConf("MIKNAS_ADMIN_UID", "admin", "管理员用户uid", false)

	ext.RegAuth(ext.Res("vist"), "使用miknas主页，一般所有人都需要", false)

	// routes
	regRoutes(ext)
}
