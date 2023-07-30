package miknas

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*
Config Item need to be declare before use.
We assume that all config value are jsonable,
because sometime we need to send to client!
*/
type ConfItem struct {
	Key        string
	Default    any
	Desc       string
	SendClient bool
	// indicate who register it
	ExtId string
	// validate input value and convert to the right value.
	CheckConv func(any) (any, error)
}

type ConfigManager struct {
	items    map[string]ConfItem
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

func (m *ConfigManager) RegConfItem(item ConfItem) error {
	k := item.Key
	if preItem, ok := m.items[k]; ok {
		return fmt.Errorf("key %s existed, registed by extid(%s)", k, preItem.ExtId)
	}
	m.items[k] = item
	err := m.UpdateValue(k, item.Default)
	if err != nil {
		delete(m.items, k)
		return err
	}
	if item.SendClient {
		m.sendlist = append(m.sendlist, k)
	}
	return nil
}

func (m *ConfigManager) UpdateValue(k string, v any) error {
	item, ok := m.items[k]
	if !ok {
		return fmt.Errorf("config key(%s) not register", k)
	}
	if item.CheckConv != nil {
		real, err := item.CheckConv(v)
		if err != nil {
			return fmt.Errorf("value(%v) of key(%s) cannot pass CheckConv, err: %v", v, k, err)
		}
		v = real
	}
	m.values[item.Key] = v
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
				fmt.Printf("UpdateFromMap Fail: err: %v", err)
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
				fmt.Printf("UpdateFromEnv Fail: err: %v", err)
			}
		}
	}
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
		map[string]ConfItem{},
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
