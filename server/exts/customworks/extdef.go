package customworks

import (
	"github.com/mikzone/miknas/server/miknas"
)

type MikNasExt struct {
	miknas.Extension
	WorkDefMap  map[string]*WorkDef
	SpaceDefMap map[string]*SpaceDef
}

const MExtId = "CustomWorks"

func New() *MikNasExt {
	return &MikNasExt{
		miknas.NewExtension(MExtId),
		make(map[string]*WorkDef),
		make(map[string]*SpaceDef),
	}
}

func (ext *MikNasExt) OnBind() {
	// you can register config, auth, routes in here
	ext.RegAuth(ext.Res("vist"), "使用CustomWorks", false)
	ext.RegConfs(
		miknas.NewConfItem("CUSTOM_WORKS_PLUGINS", []string{}, "CustomWorks相关定义文件列表", false),
		miknas.NewConfItem("CUSTOM_WORKS_SPACES", []SpaceDef{}, "CustomWorks相关实例目录列表", false),
	)
	regRoutes(ext)
}

func (ext *MikNasExt) scanDefs() {
	ConfMgr := ext.App.ConfMgr
	defFiles := ConfMgr.Get("CUSTOM_WORKS_PLUGINS").([]string)
	for _, defFile := range defFiles {
		workDef, err := ReadWorkDef(defFile)
		if err != nil {
			ext.Logger().Warn("ReadWorkDefFail", "file", defFile, "err", err)
			continue
		}
		ext.Logger().Info("ReadWorkDefSuccess", "file", defFile, "DefId", workDef.DefId)
		ext.WorkDefMap[workDef.DefId] = workDef
	}
}

func (ext *MikNasExt) scanSpaces() {
	ConfMgr := ext.App.ConfMgr
	spaceDefs := ConfMgr.Get("CUSTOM_WORKS_SPACES").([]SpaceDef)
	for _, spaceDef := range spaceDefs {
		// err := ReadSpaceExt(&spaceDef)
		// if err != nil {
		// 	ext.Logger().Warn("ReadSpaceExtFail", "SpaceDef", spaceDef, "err", err)
		// 	continue
		// }
		// workDefId := spaceDef.Ext.WorkDefId
		// _, ok := ext.WorkDefMap[workDefId]
		// if !ok {
		// 	ext.Logger().Warn("ReadSpaceExtFail", "SpaceDef", spaceDef, "err", "WorkDefId not found", "WorkDefId", workDefId)
		// 	continue
		// }
		ext.Logger().Info("ReadSpaceExtSuccess", "SpaceDef", spaceDef)
		ext.SpaceDefMap[spaceDef.Id] = &spaceDef
	}
}

func (ext *MikNasExt) OnInit() {
	// only in init, you can access db, workspace, loaded configs
	// you can register your filespace, init your db here
	ext.scanDefs()
	ext.scanSpaces()
	regCwFileSpace(ext)
}
