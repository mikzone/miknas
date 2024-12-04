package customworks

import (
	"github.com/mikzone/miknas/server/miknas"
)

type inDataQueryJob struct {
	JobId           string `json:"jobId" binding:"required"`
	ReadStdoutStart int    `json:"readStdoutStart"`
}

func queryJobResult(ch *miknas.ContextHelper) {
	var loc inDataQueryJob
	ch.BindJSON(&loc)
	jm := GetJobMgr(ch)
	jobItem, exist := jm.Jobs[loc.JobId]
	if !exist {
		ch.FailResp("jobid(%s)不存在", loc.JobId)
	}
	ret := jobItem.PackClientDict()
	ret["stdoutInfo"] = jobItem.PackClientStdOut(loc.ReadStdoutStart)
	ch.SucResp(ret)
}

type inDataQueryAllJobs struct {
	SpaceId        string `json:"spaceId"`
	PluginId       string `json:"pluginId"`
	PluginWorkRoot string `json:"pluginWorkRoot"`
}

func queryAllJobs(ch *miknas.ContextHelper) {
	var loc inDataQueryAllJobs
	ch.BindJSON(&loc)
	jm := GetJobMgr(ch)
	ret := miknas.H{}
	for jobId, jobItem := range jm.Jobs {
		if loc.SpaceId != "" && jobItem.SpaceId != loc.SpaceId {
			continue
		}
		if loc.PluginId != "" && jobItem.PluginId != loc.PluginId {
			continue
		}
		if loc.PluginWorkRoot != "" && jobItem.PluginWorkRoot != loc.PluginWorkRoot {
			continue
		}
		ret[jobId] = jobItem.PackClientDict()
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
	err := jm.TryStopJob(ch, loc.JobId, loc.KillType)
	ch.EnsureNoErr(err)
	ch.SucResp("取消成功")
}

func regCmdRoutes(ext *MikNasExt) {
	ext.POST("/queryJobResult", queryJobResult)
	ext.POST("/queryAllJobs", queryAllJobs)
	ext.POST("/reqCancelJob", reqCancelJob)
}
