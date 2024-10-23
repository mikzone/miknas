package customworks

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type WorkActionDef struct {
	Id    string
	Name  string
	Form  [](map[string]any)
	Shell map[string]any
}

type WorkDef struct {
	DefId   string
	Desc    string
	Actions []WorkActionDef
}

func ReadWorkDef(file string) (*WorkDef, error) {
	fileBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var ret WorkDef
	err = toml.Unmarshal(fileBytes, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

type SpaceExtDef struct {
	// ExtDef 定义在目录对应的文件里
	WorkDefId string
}

type SpaceDef struct {
	Id   string
	Name string
	Path string
}

func ReadSpaceExt(sp *SpaceDef) error {
	path := filepath.Join(sp.Path, ".custom_work_space.toml")
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var ret SpaceExtDef
	err = toml.Unmarshal(fileBytes, &ret)
	if err != nil {
		return err
	}
	// sp.Ext = ret
	return nil
}
