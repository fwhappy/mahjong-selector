package util

import (
	"encoding/json"
	"fmt"
)

// JsonMap 存储任意数据的map
type JsonMap map[string]interface{}

// JsonEncode 将JsonMap转成[]byte
func JsonEncode(m JsonMap) ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		err = fmt.Errorf("JsonEncode err: %s", err.Error())
	}
	return b, err
}

// JsonDecode 将[]byte转成JsonMap
func JsonDecode(b []byte) (JsonMap, error) {
	m := make(JsonMap)
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		err = fmt.Errorf("JsonDecode err: %s", err.Error())
	} else {
		m = f.(map[string]interface{})
	}
	return m, err
}

// JsonGetInt 从JsonMap中解析出一个int值
func (m JsonMap) JsonGetInt(key string) (int, bool) {
	if val, exists := m[key]; exists {
		switch val.(type) {
		case float64:
			return int(val.(float64)), true
		case int:
			return val.(int), true
		}
	}
	return 0, false
}

// JsonGetUint16 从JsonMap中解析出一个ushort值
func (m JsonMap) JsonGetUint16(key string) (uint16, bool) {
	if val, exists := m[key]; exists {
		return uint16(val.(float64)), true
	}
	return uint16(0), false
}

// JsonGetString 从JsonMap中解析出一个string值
func (m JsonMap) JsonGetString(key string) (string, bool) {
	if val, exists := m[key]; exists {
		return val.(string), true
	}
	return "", false
}

// JsonGetJsonMap 从JsonMap中解析出一个JsonMap值
func (this JsonMap) JsonGetJsonMap(key string) JsonMap {
	if val, exists := this[key]; exists {
		switch val.(type) {
		case map[string]interface{}:
			return JsonMap(val.(map[string]interface{}))
		case interface{}:
			return val.(JsonMap)
		}
	}
	return JsonMap{}
}

// InterfaceToJsonString 将任意类型的数据，转成json格式的字符串
func InterfaceToJsonString(s interface{}) (string, error) {
	byt, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(byt), nil
}
