package pan

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

func regRoutes(ext *MikNasExt) {
	ext.POST("/hello", hello)
}
