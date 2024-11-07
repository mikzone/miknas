package customworks

import "github.com/mikzone/miknas/server/miknas"

type CwFileSpace struct {
	miknas.SimpleFileSpace
}

const MFsType = "Cw"

func (fsp *CwFileSpace) Ensure(ch *miknas.ContextHelper, mode string) {
	if mode != "r" {
		panic(miknas.NewFailRet("工作区没有%s权限", mode))
	}
}

func (fsp *CwFileSpace) NewFsDriver(ch *miknas.ContextHelper, fssubid string) miknas.IFsDriver {
	ext := ch.GetApp().GetExt(MExtId).(*MikNasExt)
	spaceDef, ok := ext.SpaceDefMap[fssubid]
	if !ok {
		panic(miknas.NewFailRet("工作区不存在"))
	}
	rootDir := spaceDef.Path
	return miknas.NewBaseFsDriver(rootDir, true)
}

func regCwFileSpace(ext *MikNasExt) {
	CwFileSpace := &CwFileSpace{
		*miknas.NewSimpleFileSpace(MFsType, "", ext.Res("nouse")),
	}
	ext.RegFileSpace(CwFileSpace)
}
