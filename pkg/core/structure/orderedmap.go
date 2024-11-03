package structure

import (
	"container/list"
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

		// Check if value is of type string, and wrap in quotes if so
		if strValue, ok := value.(string); ok {
			result += fmt.Sprintf(`"%s":"%s"`, key, strValue)
		} else {
			result += fmt.Sprintf(`"%s":%v`, key, value)
		}

		if e.Next() != nil {
			result += "," // Add a comma if there are more elements
		}
	}
	result += "}"
	return result
}
