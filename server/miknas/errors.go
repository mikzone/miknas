package miknas

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IFailRet interface {
	error
	MakeRespond(c *gin.Context)
}

type FailRet struct {
	why string
}

func (p *FailRet) MakeRespond(c *gin.Context) {
	status := c.GetInt(ctxFailStatusKey)
	if status == 0 {
		status = http.StatusOK
	}
	c.JSON(status, gin.H{
		"suc": false,
		"why": p.why,
	})
}

func (p *FailRet) Error() string {
	return p.why
}

func NewFailRet(format string, a ...any) IFailRet {
	return &FailRet{
		why: fmt.Sprintf(format, a...),
	}
}

var _ IFailRet = (*FailRet)(nil)

func HandleFailRetMiddleware(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			p, ok := err.(IFailRet)
			if ok {
				p.MakeRespond(c)
				return
			}
			panic(err)
		}
	}()
	c.Next()
}
