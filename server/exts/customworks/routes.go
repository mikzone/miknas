package customworks

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mikzone/miknas/server/exts/cmdexec"
	"github.com/mikzone/miknas/server/miknas"
)

type inDataName struct {
	Name string `json:"name" form:"name" binding:"required"`
}

type inDataSpaceLocate struct {
	SpaceId string `json:"spaceId" binding:"required"`
	Fspath  string `json:"fspath"`
}

type inDataPlugin struct {
	Id string `json:"id" binding:"required"`
}

func hello(ch *miknas.ContextHelper) {
	var loc inDataName
	ch.BindJSON(&loc)
	ch.SucResp(fmt.Sprintf("hello, %s !", loc.Name))
}

func maininfo(ch *miknas.ContextHelper) {
	// spaces
	ext := ch.GetRelExt().(*MikNasExt)
	spaces := map[string]any{}
	for _, space := range ext.SpaceDefMap {
		spaceInfo := map[string]any{}
		spaceInfo["Id"] = space.Id
		spaceInfo["Name"] = space.Name
		spaceInfo["Path"] = space.Path
		spaceInfo["Plugins"] = space.Plugins
		spaces[space.Id] = spaceInfo
	}

	plugins := map[string]any{}
	for _, plugin := range ext.PluginDefMap {
		pluginInfo := map[string]any{}
		pluginInfo["Id"] = plugin.Id
		pluginInfo["Desc"] = plugin.Desc
		pluginInfo["Anchor"] = plugin.Anchor
		pluginInfo["Version"] = plugin.Version
		plugins[plugin.Id] = pluginInfo
	}

	ret := map[string]any{
		"spaces":  spaces,
		"plugins": plugins,
	}
	ch.SucResp(ret)
}

func FindPathHolder(absRootPath, curPath, targetName string) (string, error) {
	// 遍历curPath以及它的父节点，找出直接有targetName的文件
	curPath = filepath.Join(absRootPath, curPath)
	checkPath, err := filepath.Abs(curPath)
	if err != nil {
		return "", err
	}
	for {
		fp := filepath.Join(checkPath, targetName)
		_, err := os.Stat(fp)
		if err != nil {
			return checkPath, nil
		}
		if checkPath == absRootPath || checkPath == "/" || checkPath == "" {
			return "", fmt.Errorf("向上搜索找不到%s", targetName)
		}
		checkPath = filepath.Dir(checkPath)
	}
}

func queryFolderDetail(ch *miknas.ContextHelper) {
	var loc inDataSpaceLocate
	ch.BindJSON(&loc)
	ext := ch.GetRelExt().(*MikNasExt)
	spaceDef, ok := ext.SpaceDefMap[loc.SpaceId]
	if !ok {
		ch.FailResp("%s工作区不存在", loc.SpaceId)
	}
	absRootPath, err := filepath.Abs(spaceDef.Path)
	if err != nil {
		ch.FailResp("工作区路径非法")
	}
	pluginStats := []any{}
	for _, pluginId := range spaceDef.Plugins {
		pluginDef, ok := ext.PluginDefMap[pluginId]
		if !ok {
			continue
		}
		anchor := pluginDef.Anchor
		if anchor == "" {
			continue
		}
		pluginRootPath, err := FindPathHolder(absRootPath, loc.Fspath, anchor)
		if err != nil || pluginRootPath == "" {
			continue
		}
		pluginStats = append(pluginStats, map[string]any{
			"id":       pluginId,
			"rootPath": pluginRootPath,
		})
	}
	ret := map[string]any{
		"pluginStats": pluginStats,
	}
	ch.SucResp(ret)
}

func queryPluginDef(ch *miknas.ContextHelper) {
	var loc inDataPlugin
	ch.BindJSON(&loc)
	ext := ch.GetRelExt().(*MikNasExt)
	pluginDef, ok := ext.PluginDefMap[loc.Id]
	if !ok {
		ch.FailResp("%s插件不存在", loc.Id)
		return
	}
	ch.SucResp(pluginDef)
}

type inDataExecPluginJob struct {
	SpaceId  string          `json:"spaceId" binding:"required"`
	PluginId string          `json:"pluginId" binding:"required"`
	JobId    string          `json:"jobId" binding:"required"`
	Fspath   string          `json:"fspath"`
	FormData *map[string]any `json:"formData"`
}

func execPluginJob(ch *miknas.ContextHelper) {
	var loc inDataExecPluginJob
	ch.BindJSON(&loc)
	ext := ch.GetRelExt().(*MikNasExt)
	pluginDef, ok := ext.PluginDefMap[loc.PluginId]
	if !ok {
		ch.FailResp("%s插件不存在", loc.PluginId)
		return
	}
	spaceDef, ok := ext.SpaceDefMap[loc.SpaceId]
	if !ok {
		ch.FailResp("%s工作区不存在", loc.SpaceId)
		return
	}
	jobDef, ok := pluginDef.JobMap[loc.JobId]
	if !ok {
		ch.FailResp("%sJob不存在", loc.JobId)
		return
	}
	absRootPath, err := filepath.Abs(spaceDef.Path)
	if err != nil {
		ch.FailResp("工作区路径非法")
		return
	}
	pluginCurPath := filepath.Join(absRootPath, loc.Fspath)
	pluginRootPath, err := FindPathHolder(absRootPath, loc.Fspath, pluginDef.Anchor)
	if err != nil || pluginRootPath == "" {
		ch.FailResp("当前不在插件可管辖的目录下")
		return
	}
	needEnv := []string{}
	if len(jobDef.Form.FormConfs) > 0 {
		if loc.FormData == nil {
			ch.SucResp(map[string]any{
				"form":       jobDef.Form,
				"nextAction": "FillForm",
			})
			return
		} else {
			formData := *loc.FormData
			for _, conf := range jobDef.Form.FormConfs {
				key := conf["id"].(string)
				val, ok := formData[key]
				if !ok {
					ch.FailResp("表单数据不完整")
					return
				}
				needEnv = append(needEnv, fmt.Sprintf("CW_PARAM_%s=%v", key, val))
			}
		}
	}
	title := fmt.Sprintf("%s-%s", pluginDef.Title, jobDef.Name)
	jobItem := cmdexec.NewJob(title, jobDef.Cmd.Path, jobDef.Cmd.Args...)
	jobItem.Cmd.Dir = pluginCurPath
	jobItem.Cmd.Env = append(os.Environ(), needEnv...)
	cmdexec.SubmitJob(ch, jobItem)
	ch.SucResp(map[string]any{
		"jobInfo":    jobItem.PackClientDict(false),
		"nextAction": "ShowExec",
		"needEnv":    needEnv,
	})
}

func regRoutes(ext *MikNasExt) {
	ext.POST("/hello", hello)
	ext.POST("/maininfo", maininfo)
	ext.POST("/queryFolderDetail", queryFolderDetail)
	ext.POST("/queryPluginDef", queryPluginDef)
	ext.POST("/execPluginJob", execPluginJob)
}
