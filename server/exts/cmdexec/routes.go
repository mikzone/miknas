package cmdexec

import (
	"github.com/mikzone/miknas/server/miknas"
)

type inDataJobId struct {
	JobId string `json:"jobId" binding:"required"`
}

func queryJobResult(ch *miknas.ContextHelper) {
	var loc inDataJobId
	ch.BindJSON(&loc)
	jm := GetJobMgr(ch)
	jobItem, exist := jm.Jobs[loc.JobId]
	if !exist {
		ch.FailResp("jobid(%s)不存在", loc.JobId)
	}
	ch.SucResp(jobItem.PackClientDict(true))
}

func queryAllJobs(ch *miknas.ContextHelper) {
	jm := GetJobMgr(ch)
	ret := miknas.H{}
	for jobId, jobItem := range jm.Jobs {
		ret[jobId] = jobItem.PackClientDict(false)
	}
	ch.SucResp(ret)
}

type inDataCancelJob struct {
	JobId    string `json:"jobId" binding:"required"`
	KillType string `json:"killType" binding:"required"`
}

func reqCancelJob(ch *miknas.ContextHelper) {
	var loc inDataCancelJob
	ch.BindJSON(&loc)
	jm := GetJobMgr(ch)
	err := jm.TryStopJob(ch, loc.JobId)
	ch.EnsureNoErr(err)
	ch.SucResp("取消成功")
}

func regRoutes(ext *MikNasExt) {
	ext.POST("/queryJobResult", queryJobResult)
	ext.POST("/queryAllJobs", queryAllJobs)
	ext.POST("/reqCancelJob", reqCancelJob)
}
