package structure

import (
	"container/list"
	"encoding/json"
	"fmt"
)

type OrderedMap struct {
	m map[string]interface{} // value를 interface{}로 변경
	l *list.List
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		m: make(map[string]interface{}), // map의 value 타입을 interface{}로 변경
		l: list.New(),
	}
}

func (om *OrderedMap) Set(key string, value interface{}) { // value 파라미터의 타입을 interface{}로 변경
	if _, exists := om.m[key]; !exists {
		om.l.PushBack(key) // 새로운 키는 리스트에 추가
	}
	om.m[key] = value
}

func (om *OrderedMap) Get(key string) (interface{}, bool) { // 반환 타입을 interface{}로 변경
	value, exists := om.m[key]
	return value, exists
}

func (om *OrderedMap) Ser() string {
	result := "{"
	for e := om.l.Front(); e != nil; e = e.Next() {
		key := e.Value.(string)
		value := om.m[key]

		// If value is OrderedMap, recursively call Ser()
		if nestedMap, ok := value.(*OrderedMap); ok {
			result += fmt.Sprintf(`"%s":%s`, key, nestedMap.Ser())
		} else if strValue, ok := value.(string); ok {
			// If value is string, wrap with quotes
			result += fmt.Sprintf(`"%s":"%s"`, key, strValue)
		} else {
			// For other types, print as is
			result += fmt.Sprintf(`"%s":%v`, key, value)
		}

		if e.Next() != nil {
			result += "," // Add comma if there are more elements
		}
	}
	result += "}"
	return result
}

func ParseOrderedMap(jsonStr string) (*OrderedMap, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, err
	}

	om := NewOrderedMap()

	// Iterate through JSON decoded map and convert to OrderedMap
	for key, value := range data {
		switch v := value.(type) {
		case map[string]interface{}:
			// For nested objects, recursively call ParseOrderedMap
			nestedJsonStr, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			nestedMap, err := ParseOrderedMap(string(nestedJsonStr))
			if err != nil {
				return nil, err
			}
			om.Set(key, nestedMap)
		case float64:
			// Numbers in JSON are parsed as float64 by default
			// Convert to int64 if it can be represented as an integer
			if v == float64(int64(v)) {
				om.Set(key, int64(v))
			} else {
				om.Set(key, v)
			}
		default:
			// Set strings and other types as is
			om.Set(key, value)
		}
	}

	return om, nil
}

func (om *OrderedMap) Keys() []string {
	keys := make([]string, om.l.Len())
	i := 0
	for e := om.l.Front(); e != nil; e = e.Next() {
		keys[i] = e.Value.(string)
		i++
	}
	return keys
}
