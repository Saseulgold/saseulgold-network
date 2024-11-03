package structure

import (
	"container/list"
	"encoding/json"
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

func (om *OrderedMap) MarshalJSON() ([]byte, error) {
	kvPairs := make([]map[string]interface{}, 0) // kvPairs의 타입을 변경
	for e := om.l.Front(); e != nil; e = e.Next() {
		key := e.Value.(string)
		kvPairs = append(kvPairs, map[string]interface{}{key: om.m[key]}) // m의 value 타입에 맞게 변경
	}
	return json.Marshal(kvPairs)
}

/*
func main() {
	orderedMap := NewOrderedMap()
	orderedMap.Set("key2", "value2")
	orderedMap.Set("key1", 42) // 다양한 타입의 값이 허용됨

	jsonData, err := json.Marshal(orderedMap)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonData))
}
*/
