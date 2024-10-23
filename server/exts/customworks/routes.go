package customworks

import (
	"fmt"

	"github.com/mikzone/miknas/server/miknas"
)

type inDataName struct {
	Name string `json:"name" form:"name" binding:"required"`
}

func hello(ch *miknas.ContextHelper) {
	var loc inDataName
	ch.BindJSON(&loc)
	ch.SucResp(fmt.Sprintf("hello, %s !", loc.Name))
}

func maininfo(ch *miknas.ContextHelper) {
	// spaces
	ext := ch.GetRelExt().(*MikNasExt)
	ret := map[string]any{
		"spaces": ext.SpaceDefMap,
		"defs":   ext.WorkDefMap,
	}
	ch.SucResp(ret)
}

func regRoutes(ext *MikNasExt) {
	ext.POST("/hello", hello)
	ext.POST("/maininfo", maininfo)
}
