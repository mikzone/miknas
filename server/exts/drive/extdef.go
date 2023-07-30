package drive

import (
	"github.com/mikzone/miknas/server/miknas"
)

type MikNasExt struct {
	miknas.Extension
}

func New() *MikNasExt {
	return &MikNasExt{miknas.NewExtension("Drive")}
}

func (ext *MikNasExt) OnBind() {
	regRoutes(ext)
}

func (ext *MikNasExt) OnInit() {

}
