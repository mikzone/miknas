package customworks

import (
	"fmt"
	"os"
	"time"

	"github.com/mikzone/miknas/server/miknas"
)

type MikNasExt struct {
	miknas.Extension
	JmInst       *JobMgr
	PluginDefMap map[string]*PluginDef
	SpaceDefMap  map[string]*SpaceDef
}

const MExtId = "CustomWorks"

func New() *MikNasExt {
	return &MikNasExt{
		miknas.NewExtension(MExtId),
		NewJobMgr(),
		make(map[string]*PluginDef),
		make(map[string]*SpaceDef),
	}
}

func (ext *MikNasExt) OnBind() {
	// you can register config, auth, routes in here
	ext.RegAuth(ext.Res("vist"), "使用CustomWorks", true)
	ext.RegConfs(
		miknas.NewConfItem("CUSTOM_WORKS_PLUGINS", []string{}, "CustomWorks插件定义文件列表", false),
		miknas.NewConfItem("CUSTOM_WORKS_SPACES", []SpaceDef{}, "CustomWorks工作区列表", false),
	)
	regCmdRoutes(ext)
	regRoutes(ext)
}

func (ext *MikNasExt) scanPlugins() {
	ConfMgr := ext.App.ConfMgr
	defFiles := ConfMgr.Get("CUSTOM_WORKS_PLUGINS").([]string)
	for _, defFile := range defFiles {
		pluginDef, err := ReadPluginDef(defFile)
		if err != nil {
			panic(fmt.Errorf("ReadPluginDefFail, file: %s, error: %v", defFile, err))
		}
		ext.Logger().Info("RegPluginDef", "file", defFile, "PluginId", pluginDef.Id)
		ext.PluginDefMap[pluginDef.Id] = pluginDef
	}
}

func (ext *MikNasExt) scanSpaces() {
	ConfMgr := ext.App.ConfMgr
	spaceDefs := ConfMgr.Get("CUSTOM_WORKS_SPACES").([]SpaceDef)
	for _, spaceDef := range spaceDefs {
		_, err := os.Stat(spaceDef.Path)
		if err != nil {
			panic(fmt.Errorf("ScanSpaceFail, SpaceId: %v, Path(%v) is not exist", spaceDef.Id, spaceDef.Path))
		}
		_, ok := ext.SpaceDefMap[spaceDef.Id]
		if ok {
			panic(fmt.Errorf("ScanSpaceFail, SpaceId: %v is already registed", spaceDef.Id))
		}
		ext.Logger().Info("ReadSpaceSuccess", "SpaceDef", spaceDef)
		ext.SpaceDefMap[spaceDef.Id] = &spaceDef
	}
}

func (ext *MikNasExt) OnInit() {
	// only in init, you can access db, workspace, loaded configs
	// you can register your filespace, init your db here
	ext.scanPlugins()
	ext.scanSpaces()
	regCwFileSpace(ext)
	ext.Logger().Info("CreatedJobMgr", "Cap", ext.JmInst.Pool.Cap())
	go func() {
		for {
			time.Sleep(1 * time.Second)
			ext.JmInst.TryMaintainLineUps()
		}
	}()
}
