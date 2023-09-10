package cmdexec

import (
	"time"

	"github.com/mikzone/miknas/server/miknas"
)

const ExtId = "CmdExec"

type MikNasExt struct {
	miknas.Extension
	JmInst *JobMgr
}

func New() *MikNasExt {
	return &MikNasExt{
		miknas.NewExtension(ExtId),
		NewJobMgr(),
	}
}

func (ext *MikNasExt) OnBind() {
	// you can register config, auth, routes in here
	if ext.GetId() != ExtId {
		panic("不能修改 CmdExec 扩展的id")
	}
	ext.RegAuth(ext.Res("vist"), "浏览命令执行结果", false)
	regRoutes(ext)
}

func (ext *MikNasExt) OnInit() {
	// only in init, you can access db, workspace, loaded configs
	// you can register your filespace, init your db here
	ext.Logger().Info("CreatedJobMgr", "Cap", ext.JmInst.Pool.Cap())
	go func() {
		for {
			time.Sleep(1 * time.Second)
			ext.JmInst.TryMaintainLineUps()
		}
	}()
}
