package miknas

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml/v2"
)

/*
Config Item need to be declare before use.
We assume that all config value are jsonable,
because sometime we need to send to client!
*/

type IConfItem interface {
	// 配置文件
	GetKey() string
	GetDefault() any
	GetDesc() string
	GetExtId() string
	SetExtId(string)
	GetSendClient() bool
	CheckConv(any) (any, error)
}

type TConfItem[T any] struct {
	Key        string
	Default    T
	Desc       string
	SendClient bool
	// indicate who register it
	ExtId           string
	CustomCheckConv func(any) (any, error)
}

func (c *TConfItem[T]) GetKey() string {
	return c.Key
}

func (c *TConfItem[T]) GetDefault() any {
	return c.Default
}

func (c *TConfItem[T]) GetDesc() string {
	return c.Desc
}

func (c *TConfItem[T]) GetExtId() string {
	return c.ExtId
}

func (c *TConfItem[T]) SetExtId(extid string) {
	c.ExtId = extid
}

func (c *TConfItem[T]) GetSendClient() bool {
	return c.SendClient
}

func (c *TConfItem[T]) CheckConv(v any) (any, error) {
	if c.CustomCheckConv != nil {
		return c.CustomCheckConv(v)
	}
	return CheckConvTAny[T](v)
}

var _ IConfItem = (*TConfItem[int])(nil)

type ConfigManager struct {
	items    map[string]IConfItem
	values   map[string]any
	sendlist []string
}

func (m *ConfigManager) PackClientDict() gin.H {
	ret := gin.H{}
	for _, confkey := range m.sendlist {
		ret[confkey] = m.values[confkey]
	}
	return ret
}

func (m *ConfigManager) RegConfItem(item IConfItem) error {
	k := item.GetKey()
	if preItem, ok := m.items[k]; ok {
		return fmt.Errorf("key %s existed, registed by extid(%s)", k, preItem.GetExtId())
	}
	m.items[k] = item
	err := m.UpdateValue(k, item.GetDefault())
	if err != nil {
		delete(m.items, k)
		return err
	}
	if item.GetSendClient() {
		m.sendlist = append(m.sendlist, k)
	}
	return nil
}

func (m *ConfigManager) UpdateValue(k string, v any) error {
	item, ok := m.items[k]
	if !ok {
		return fmt.Errorf("config key(%s) not register", k)
	}
	real, err := item.CheckConv(v)
	if err != nil {
		return fmt.Errorf("value(%v) of key(%s) cannot pass CheckConv, err: %v", v, k, err)
	}
	m.values[item.GetKey()] = real
	return nil
}

func (m *ConfigManager) IsConfKey(k string) bool {
	_, ok := m.items[k]
	return ok
}

func (m *ConfigManager) UpdateFromMap(obj map[string]any) {
	for k, v := range obj {
		if m.IsConfKey(k) {
			err := m.UpdateValue(k, v)
			if err != nil {
				panic(fmt.Errorf("ConfigManager.UpdateFromMap Fail: %v", err))
			}
		}
	}
}

func (m *ConfigManager) UpdateFromEnv() {
	for k := range m.items {
		v := os.Getenv(k)
		if len(v) > 0 {
			err := m.UpdateValue(k, v)
			if err != nil {
				panic(fmt.Errorf("UpdateFromEnv Fail: err: %v", err))
			}
		}
	}
}

func (m *ConfigManager) UpdateFromTomlFile(filepath string) error {
	Content, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("open toml file(%s) Fail: %v", filepath, err)
	}
	anyMap := make(map[string]any, 0)
	err = toml.Unmarshal(Content, &anyMap)
	if err != nil {
		return fmt.Errorf("unmarshal toml file(%s) Fail: %v", filepath, err)
	}

	m.UpdateFromMap(anyMap)
	return nil
}

func (m *ConfigManager) PrintConfigs() {
	ret := gin.H{}
	for k := range m.items {
		v := m.values[k]
		ret[k] = v
	}
	str, err := json.MarshalIndent(ret, "", "  ")
	if err != nil {
		fmt.Printf("PrintConfig Error in json marshal: %v", err)
		return
	}
	fmt.Println("[MikNas]AppConfigs:", string(str))
}

func (m *ConfigManager) Get(k string) any {
	return m.values[k]
}

func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		map[string]IConfItem{},
		map[string]any{},
		[]string{},
	}
}

func CheckConvStr(value any) (any, error) {
	v, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("value(%v) not a string value", v)
	}
	return v, nil
}

func CheckConvInt(value any) (any, error) {
	str1, ok1 := value.(string)
	if ok1 {
		// string need to convert to int
		intV, err1 := strconv.Atoi(str1)
		if err1 != nil {
			return nil, err1
		}
		return intV, nil
	}
	v, ok := value.(int)
	if !ok {
		return nil, fmt.Errorf("value(%v) is not an integer", v)
	}
	return v, nil
}

func CheckConvMap(value any) (any, error) {
	str1, ok1 := value.(string)
	if ok1 {
		// string need to convert to map
		anyMap := make(map[string]any, 0)
		if err := json.Unmarshal([]byte(str1), &anyMap); err != nil {
			return nil, err
		}
		return anyMap, nil
	}
	v, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("value(%v) is not map[string]any type", v)
	}
	return v, nil
}

func CheckConvList(value any) (any, error) {
	str1, ok1 := value.(string)
	if ok1 {
		// string need to convert to list
		anyList := make([]any, 0)
		if err := json.Unmarshal([]byte(str1), &anyList); err != nil {
			return nil, err
		}
		return anyList, nil
	}
	v, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("value(%v) is not []any type", v)
	}
	return v, nil
}

func CheckConvTAny[T any](value any) (any, error) {
	vv0, ok0 := value.(T)
	if ok0 {
		return vv0, nil
	}
	str1, ok1 := value.(string)
	if !ok1 {
		// 如果不是字符串，比如用toml这种会先格式化好的，则尝试先转回json字符串再Unmarshal（低效一些的方法）
		byte2, err2 := json.Marshal(value)
		if err2 != nil {
			return nil, fmt.Errorf("json marshal err: %v", err2)
		}
		str1 = string(byte2)
	}
	vv1 := new(T)
	if err := json.Unmarshal([]byte(str1), &vv1); err != nil {
		return nil, fmt.Errorf("json unmarshal err: %v", err)
	}
	return *vv1, nil
}

func NewConfItem[T any](key string, defv T, desc string, sendClient bool) IConfItem {
	return &TConfItem[T]{Key: key, Default: defv, Desc: desc, SendClient: sendClient}
}
