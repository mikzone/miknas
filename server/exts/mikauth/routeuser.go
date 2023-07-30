package mikauth

import (
	"fmt"

	"github.com/mikzone/miknas/server/miknas"
	"golang.org/x/crypto/bcrypt"
)

func login(ch *miknas.ContextHelper) {
	session := ch.GetSession()
	uid := session.Get("uid")
	// 已经登录了
	if uid != nil {
		ch.Redirect(ch.ClientUrl(""))
		return
	}
	session.Clear()
	session.Save()
	ext := ch.GetRelExt()
	ch.Redirect(ext.ClientUrl("login"))
}

var myrsa miknas.MyRSA = miknas.MustGenMyRSA()
var pubkey string

func init() {
	pk, err := myrsa.DumpPublicKeyBase64()
	if err != nil {
		panic("can not gen public key, it will failed for login")
	}
	pubkey = pk
}

func PwdHash(pwd string) (string, error) {
	// 第二个参数是进行哈希的次数，这里采用了默认值10,数字越大生成的密码速度越慢，成本越大。但是更安全
	// bcrypt每次生成的编码是不同的，较于md5更安全
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

func PwdVerify(pwd, hash string) bool {
	// CompareHashAndPassword 比较用户输入的明文和和数据库取出的的密码解析后是否匹配
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err == nil
}

func querySecretToken(ch *miknas.ContextHelper) {
	ch.SucResp(miknas.H{
		"token": pubkey,
	})
}

func tryLoginWithUser(ch *miknas.ContextHelper, uid, pwd string, why string) error {
	if len(uid) <= 0 {
		return fmt.Errorf("uid can not be empty")
	}
	if len(pwd) <= 0 {
		return fmt.Errorf("pwd can not be empty")
	}
	session := ch.GetSession()
	session.Clear()

	// 查找数据库看看用户是否存在
	user := GetUserByUid(ch.GetApp().Db, uid)
	if user == nil {
		return fmt.Errorf("user not exist")
	}
	if !PwdVerify(pwd, user.Pwd) {
		return fmt.Errorf("password incorrect")
	}
	app := ch.GetApp()
	adminUid := app.ConfMgr.Get("MIKNAS_ADMIN_UID").(string)
	role := user.Role
	if uid == adminUid {
		role = "admin"
	}
	session.Set("uid", uid)
	session.Set("role", role)
	session.Save()
	ch.UsLog("login", miknas.H{
		"uid": uid,
		"why": why,
	})
	return nil
}

func jumpLoginWithMsg(ch *miknas.ContextHelper, msg string) {
	ext := ch.GetRelExt()
	ch.Jump(msg, ext.ClientUrl("/login"), 3, false)
}

func jumpRegisterWithMsg(ch *miknas.ContextHelper, msg string) {
	ext := ch.GetRelExt()
	ch.Jump(msg, ext.ClientUrl("/register"), 3, false)
}

func loginto(ch *miknas.ContextHelper) {
	app := ch.GetApp()
	ext := ch.GetRelExt()
	if app.ConfMgr.Get("MIKNAS_AUTH_EXTS") != ext.GetId() {
		jumpLoginWithMsg(ch, "当前的用户管理不归本插件管理")
		return
	}
	uid := ch.Ctx.PostForm("uid")
	if len(uid) <= 0 {
		jumpLoginWithMsg(ch, "uid can not be empty")
		return
	}
	pwd := ch.Ctx.PostForm("pwd")
	if len(pwd) <= 0 {
		jumpLoginWithMsg(ch, "pwd can not be empty")
		return
	}
	pwd, decErr := myrsa.Decrypt(pwd)
	if decErr != nil {
		jumpLoginWithMsg(ch, decErr.Error())
		return
	}
	err := tryLoginWithUser(ch, uid, pwd, "loginto")
	if err != nil {
		jumpLoginWithMsg(ch, err.Error())
		return
	}
	ch.Redirect(ch.ClientUrl(""))
}

func logout(ch *miknas.ContextHelper) {
	session := ch.GetSession()
	session.Clear()
	session.Save()
	ch.Redirect(ch.ClientUrl(""))
}

func register(ch *miknas.ContextHelper) {
	app := ch.GetApp()
	ext := ch.GetRelExt()
	if app.ConfMgr.Get("MIKNAS_AUTH_EXTS") != ext.GetId() {
		jumpRegisterWithMsg(ch, "当前的用户管理不归本插件管理")
		return
	}
	uid := ch.Ctx.PostForm("uid")
	if len(uid) < 3 {
		jumpRegisterWithMsg(ch, "uid不得小于3个字符")
		return
	}
	pwd := ch.Ctx.PostForm("pwd")
	if len(pwd) < 6 {
		jumpRegisterWithMsg(ch, "密码不得小于6个字符")
		return
	}
	pwd, decErr := myrsa.Decrypt(pwd)
	if decErr != nil {
		jumpRegisterWithMsg(ch, decErr.Error())
		return
	}
	name := ch.Ctx.PostForm("name")
	if len(name) < 3 {
		jumpRegisterWithMsg(ch, "名字不得小于3个字符")
		return
	}
	user := GetUserByUid(app.Db, uid)
	if user != nil {
		jumpRegisterWithMsg(ch, fmt.Sprintf("uid(%s)已被占用，请换一个", uid))
		return
	}
	hashedPwd, hashErr := PwdHash(pwd)
	if hashErr != nil {
		jumpRegisterWithMsg(ch, hashErr.Error())
		return
	}
	adminUid := app.ConfMgr.Get("MIKNAS_ADMIN_UID").(string)
	role := "tour"
	if uid == adminUid {
		role = "admin"
	}
	newUser := MikauthUser{
		Uid:  uid,
		Pwd:  hashedPwd,
		Name: name,
		Role: role,
	}
	app.Db.Create(&newUser)
	err := tryLoginWithUser(ch, uid, pwd, "register")
	if err != nil {
		jumpLoginWithMsg(ch, err.Error())
		return
	}
	ch.Redirect(ch.ClientUrl(""))
}

func queryAllUser(ch *miknas.ContextHelper) {
	ch.Ensure(ch.GetRelExt().Res("manager"))
	db := ch.GetApp().Db
	var users []MikauthUser
	db.Find(&users)
	infos := []miknas.H{}
	for _, user := range users {
		infos = append(infos, PackUserInfo(&user))
	}
	var roleRecs []MikauthRole
	db.Select("Id").Find(&roleRecs)
	roles := []string{}
	for _, roleRec := range roleRecs {
		roles = append(roles, roleRec.Id)
	}

	ch.SucResp(miknas.H{
		"roles": roles,
		"users": infos,
	})
}

type inDataModifyUserRole struct {
	Uid  string `json:"uid" binding:"required"`
	Role string `json:"role" binding:"required"`
}

func modifyUserRole(ch *miknas.ContextHelper) {
	ch.Ensure(ch.GetRelExt().Res("manager"))
	var modify inDataModifyUserRole
	ch.BindJSON(&modify)
	db := ch.GetApp().Db
	user := GetUserByUid(db, modify.Uid)
	if user == nil {
		ch.FailResp("用户不存在")
		return
	}
	roleRec := GetRoleById(db, modify.Role)
	if roleRec == nil {
		ch.FailResp("角色(%s)不存在", modify.Role)
		return
	}
	user.Role = modify.Role
	result := db.Save(&user)
	if result.Error != nil {
		ch.FailResp(result.Error.Error())
		return
	}
	ch.SucResp(result.RowsAffected)
}

type inDataUid struct {
	Uid string `json:"uid" binding:"required"`
}

func removeUser(ch *miknas.ContextHelper) {
	ch.Ensure(ch.GetRelExt().Res("manager"))
	var modify inDataUid
	ch.BindJSON(&modify)
	db := ch.GetApp().Db
	user := GetUserByUid(db, modify.Uid)
	if user == nil {
		ch.SucResp(1)
		return
	}
	result := db.Delete(&user)
	if result.Error != nil {
		ch.FailResp(result.Error.Error())
		return
	}
	ch.SucResp(result.RowsAffected)
}

func currentUserinfo(ch *miknas.ContextHelper) {
	session := ch.GetSession()
	uid := miknas.AnyToStr(session.Get("uid"))
	if uid == "" {
		ch.FailResp("当前未登录")
		return
	}
	user := GetUserByUid(ch.GetApp().Db, uid)
	if user == nil {
		ch.FailResp("用户不存在")
		return
	}
	info := PackUserInfo(user)
	realRole := miknas.AnyToStr(session.Get("role"))
	info["realRole"] = realRole
	ch.SucResp(info)
}

type inDataName struct {
	Name string `json:"name" binding:"required"`
}

func modifyNickname(ch *miknas.ContextHelper) {
	var modify inDataName
	ch.BindJSON(&modify)
	session := ch.GetSession()
	uid := miknas.AnyToStr(session.Get("uid"))
	if uid == "" {
		ch.FailResp("当前未登录")
		return
	}
	db := ch.GetApp().Db
	user := GetUserByUid(db, uid)
	if user == nil {
		ch.FailResp("用户不存在")
		return
	}
	user.Name = modify.Name
	result := db.Save(&user)
	if result.Error != nil {
		ch.FailResp(result.Error.Error())
		return
	}
	ch.SucResp(result.RowsAffected)
}

type inDataChangePwd struct {
	OldPwd string `json:"oldPwd" binding:"required"`
	NewPwd string `json:"newPwd" binding:"required"`
}

func modifyPassword(ch *miknas.ContextHelper) {
	var modify inDataChangePwd
	ch.BindJSON(&modify)
	session := ch.GetSession()
	uid := miknas.AnyToStr(session.Get("uid"))
	if uid == "" {
		ch.FailResp("当前未登录")
		return
	}
	db := ch.GetApp().Db
	user := GetUserByUid(db, uid)
	if user == nil {
		ch.FailResp("用户不存在")
		return
	}
	oldPwd, decErr := myrsa.Decrypt(modify.OldPwd)
	if decErr != nil {
		ch.FailResp(decErr.Error())
		return
	}
	if !PwdVerify(oldPwd, user.Pwd) {
		ch.FailResp("旧密码不匹配")
		return
	}
	newPwd, decErr2 := myrsa.Decrypt(modify.NewPwd)
	if decErr2 != nil {
		ch.FailResp(decErr2.Error())
		return
	}
	if len(newPwd) < 6 {
		ch.FailResp("密码不得小于6个字符")
		return
	}
	savePwd, hashErr := PwdHash(newPwd)
	if hashErr != nil {
		ch.FailResp(hashErr.Error())
		return
	}
	user.Pwd = savePwd
	result := db.Save(&user)
	if result.Error != nil {
		ch.FailResp(result.Error.Error())
		return
	}
	ch.SucResp(result.RowsAffected)
}

func regUserRoutes(ext *MikNasExt) {
	ext.GET("login", login)
	ext.GET("logout", logout)
	ext.POST("querySecretToken", querySecretToken)
	ext.POST("loginto", loginto)
	ext.POST("register", register)
	ext.POST("queryAllUser", queryAllUser)
	ext.POST("modifyUserRole", modifyUserRole)
	ext.POST("removeUser", removeUser)
	ext.POST("currentUserinfo", currentUserinfo)
	ext.POST("modifyNickname", modifyNickname)
	ext.POST("modifyPassword", modifyPassword)
}
