package customworks

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/mikzone/miknas/server/miknas"
	"github.com/pelletier/go-toml/v2"
)

type PluginActionCmdDef struct {
	Path string
	Args []string
}

type FormSelectOptionDef struct {
	Label string `toml:"label" json:"label"`
	Value string `toml:"value" json:"value"`
}

type FormConfDef struct {
	Id             string                 `toml:"id" json:"id"`
	Title          string                 `toml:"title" json:"title"`
	Component      string                 `toml:"component" json:"component"`
	Default        string                 `toml:"default" json:"default"`
	Desc           string                 `toml:"desc" json:"desc"`
	ComponentProps map[string]any         `toml:"componentProps" json:"componentProps"`
	SelectOptions  *[]FormSelectOptionDef `toml:"selectOptions" json:"selectOptions"`
	HintList       []string               `toml:"hintList" json:"hintList"` // 在MdcTextAutoComplete组件中显示的候选词
}

type PluginJobFormDef struct {
	ConfirmLabel string        `toml:"confirmLabel" json:"confirmLabel"`
	FormConfs    []FormConfDef `toml:"formConfs" json:"formConfs"`
}

type PluginJobDef struct {
	Id                   string
	Name                 string
	Icon                 string
	Form                 PluginJobFormDef
	Cmd                  PluginActionCmdDef
	Confirm              bool               // 是否在执行前给个二次确认
	NameSpaceTpl         string             // 同一个命名空间下的任务会排队执行
	NameSpaceTemplate    *template.Template // 命名空间模板，用于生成实际的命名空间
	FormAbstractTpl      string
	FormAbstractTemplate *template.Template
}

type PluginDef struct {
	Id            string
	Title         string
	Desc          string
	Anchor        string
	ShowInSubDirs bool
	Version       string
	Jobs          []*PluginJobDef
	JobMap        map[string]*PluginJobDef
	RootDir       string
}

func ReadPluginJobDef(file string) (*PluginJobDef, error) {
	fileBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var ret PluginJobDef
	err = toml.Unmarshal(fileBytes, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

const mLenTomlExt = len(".toml")

func ReadPluginDef(file string) (*PluginDef, error) {
	fileBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var ret PluginDef
	err = toml.Unmarshal(fileBytes, &ret)
	if err != nil {
		return nil, err
	}
	// 扫描jobs文件夹里的所有toml文件
	rootDir := filepath.Dir(file)
	jobsDir := filepath.Join(rootDir, "jobs")
	err = filepath.Walk(jobsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".toml" {
			return nil
		}
		job, err := ReadPluginJobDef(path)
		if err != nil {
			return fmt.Errorf("ReadPluginJobDefFail, file: %s, error: %v", path, err)
		}
		baseName := filepath.Base(path)
		job.Id = baseName[:len(baseName)-mLenTomlExt]
		if job.NameSpaceTpl != "" {
			tplName := fmt.Sprintf("Customworks_NameSpace_%s_%s", ret.Id, job.Id)
			job.NameSpaceTemplate = template.Must(template.New(tplName).Parse(job.NameSpaceTpl))
		}
		if job.FormAbstractTpl != "" {
			tplName := fmt.Sprintf("Customworks_FormAbstract_%s_%s", ret.Id, job.Id)
			job.FormAbstractTemplate = template.Must(template.New(tplName).Parse(job.FormAbstractTpl))
		}
		ret.Jobs = append(ret.Jobs, job)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk jobs fail: %v", err)
	}

	ret.JobMap = miknas.List2Map(ret.Jobs, "Id")
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, fmt.Errorf("get abs root dir fail: %v", err)
	}
	ret.RootDir = absRootDir
	return &ret, nil
}

type SpaceDef struct {
	Id              string
	Name            string
	Path            string
	Plugins         []string
	CanWalkDirRoles []string
}

func MustExecTemplate(tpl *template.Template, data any) string {
	if tpl == nil {
		return ""
	}
	buf := bytes.Buffer{}
	err := tpl.Execute(&buf, data)
	if err != nil {
		panic(fmt.Errorf("exec template fail: %v", err))
	}
	return buf.String()
}
