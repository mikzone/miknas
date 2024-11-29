package customworks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

func FindPathHolder(absRootPath, curPath, targetName string, travelParent bool) (string, error) {
	// 遍历curPath以及它的父节点，找出直接有targetName的文件
	curPath = filepath.Join(absRootPath, curPath)
	checkPath, err := filepath.Abs(curPath)
	if err != nil {
		return "", err
	}
	for {
		fp := filepath.Join(checkPath, targetName)
		_, err := os.Stat(fp)
		if err == nil {
			return checkPath, nil
		}
		if !travelParent {
			return "", fmt.Errorf("当前目录没有%s", targetName)
		}
		if checkPath == absRootPath || checkPath == "/" || checkPath == "" {
			return "", fmt.Errorf("向上搜索也找不到%s", targetName)
		}
		checkPath = filepath.Dir(checkPath)
	}
}

type spaceSubDirInfo struct {
	DirName   string   `json:"dirName"` // 文件夹名
	PluginIds []string `json:"pluginIds"`
}

func calcSubDirs(absRootPath, curPath string, pluginDefList []*PluginDef) ([]spaceSubDirInfo, error) {
	ret := []spaceSubDirInfo{}
	checkPath := filepath.Join(absRootPath, curPath)
	checkPath, err := filepath.Abs(checkPath)
	if err != nil {
		return ret, err
	}
	files, err := os.ReadDir(checkPath)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			dirName := file.Name()
			newCurPath := filepath.Join(curPath, dirName)
			pluginIds := []string{}
			for _, pluginDef := range pluginDefList {
				anchor := pluginDef.Anchor
				pluginRootPath, err := FindPathHolder(absRootPath, newCurPath, anchor, false)
				if err != nil || pluginRootPath == "" {
					continue
				}
				pluginIds = append(pluginIds, pluginDef.Id)
			}
			if len(pluginIds) > 0 {
				ret = append(ret, spaceSubDirInfo{
					DirName:   dirName,
					PluginIds: pluginIds,
				})
			}
		}
	}
	return ret, nil
}

func queryFolderDetail(ch *miknas.ContextHelper) {
	var loc inDataSpaceLocate
	ch.BindJSON(&loc)
	ext := ch.GetRelExt().(*MikNasExt)
	spaceDef, ok := ext.SpaceDefMap[loc.SpaceId]
	if !ok {
		ch.FailResp("%s工作区不存在", loc.SpaceId)
		return
	}
	absRootPath, err := filepath.Abs(spaceDef.Path)
	if err != nil {
		ch.FailResp("工作区路径非法")
		return
	}
	canWalkDir := false
	ua := ch.GetUserAuth()
	for _, roleId := range spaceDef.CanWalkDirRoles {
		if ua.HasRole(roleId) {
			canWalkDir = true
			break
		}
	}
	pluginStats := []any{}
	pluginDefList := []*PluginDef{}
	for _, pluginId := range spaceDef.Plugins {
		pluginDef, ok := ext.PluginDefMap[pluginId]
		if !ok {
			continue
		}
		anchor := pluginDef.Anchor
		if anchor == "" {
			continue
		}
		pluginDefList = append(pluginDefList, pluginDef)
		pluginRootPath, err := FindPathHolder(absRootPath, loc.Fspath, anchor, pluginDef.ShowInSubDirs)
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
		"canWalkDir":  canWalkDir,
	}
	if !canWalkDir {
		// 不能浏览文件夹的话，构造一个可用列表
		subdirs, err := calcSubDirs(absRootPath, loc.Fspath, pluginDefList)
		if err != nil {
			ch.FailResp("遍历子目录失败: %v", err.Error())
			return
		}
		ret["subdirs"] = subdirs
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
	pluginRootPath, err := FindPathHolder(absRootPath, loc.Fspath, pluginDef.Anchor, pluginDef.ShowInSubDirs)
	if err != nil || pluginRootPath == "" {
		ch.FailResp("当前不在插件可管辖的目录下")
		return
	}
	needEnv := []string{
		fmt.Sprintf("CW_PLUGIN_WORK_ROOT=%s", pluginRootPath),
		fmt.Sprintf("CW_PLUGIN_WORK_DIR=%s", pluginCurPath),
		fmt.Sprintf("CW_PLUGIN_DEF_ROOT=%s", pluginDef.RootDir),
	}
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
				key := conf.Id
				val, ok := formData[key]
				if !ok {
					ch.FailResp("表单数据不完整")
					return
				}
				if conf.SelectOptions != nil {
					valStr, ok := val.(string)
					if !ok {
						ch.FailResp("表单数据不完整")
						return
					}
					isInValues := false
					for _, option := range *conf.SelectOptions {
						if valStr == option.Value {
							isInValues = true
							break
						}
					}
					if !isInValues {
						ch.FailResp("表单数据项(%s)的值不在可选访问内", conf.Title)
						return
					}
				}

				if strings.HasPrefix(key, "CW_PARAM_") {
					needEnv = append(needEnv, fmt.Sprintf("%s=%v", key, val))
				}
			}
		}
	}
	title := fmt.Sprintf("%s-%s", pluginDef.Title, jobDef.Name)
	jobItem := cmdexec.NewJob(title, jobDef.Cmd.Path, jobDef.Cmd.Args...)
	jobItem.Cmd.Dir = pluginCurPath
	jobItem.Cmd.Env = append(os.Environ(), needEnv...)
	jobItem.NameSpace = jobDef.NameSpace
	cmdexec.SubmitJob(ch, jobItem)
	ch.SucResp(map[string]any{
		"jobInfo":    jobItem.PackClientDict(),
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
