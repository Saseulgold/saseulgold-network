package structure

import (
	"container/list"
	"fmt"
	"strconv"
	"strings"
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
	if len(jsonStr) < 2 {
		return nil, fmt.Errorf("jsonStr is not a valid JSON string")
	}

	if jsonStr[0] != '{' || jsonStr[len(jsonStr)-1] != '}' {
		return nil, fmt.Errorf("jsonStr is not a valid JSON object")
	}

	om := NewOrderedMap()
	jsonStr = jsonStr[1 : len(jsonStr)-1] // 중괄호 제거

	if len(jsonStr) == 0 {
		return om, nil
	}

	var key string
	var value string
	var isInQuote bool
	var isInValue bool
	var nestCount int
	var start int

	for i := 0; i < len(jsonStr); i++ {
		c := jsonStr[i]

		if c == '"' && (i == 0 || jsonStr[i-1] != '\\') {
			isInQuote = !isInQuote
			if !isInValue {
				if isInQuote {
					start = i + 1
				} else {
					key = jsonStr[start:i]
				}
			}
			continue
		}

		if !isInQuote {
			if c == ':' {
				isInValue = true
				start = i + 1
				continue
			}

			if c == '{' {
				nestCount++
			} else if c == '}' {
				nestCount--
			}

			if (c == ',' && nestCount == 0) || i == len(jsonStr)-1 {
				if i == len(jsonStr)-1 {
					value = jsonStr[start : i+1]
				} else {
					value = jsonStr[start:i]
				}

				value = strings.TrimSpace(value)
				if len(value) > 0 {
					if value[0] == '"' && value[len(value)-1] == '"' {
						// 문자열 값
						om.Set(key, value[1:len(value)-1])
					} else if value[0] == '{' {
						// 중첩된 OrderedMap
						nestedMap, err := ParseOrderedMap(value)
						if err != nil {
							return nil, err
						}
						om.Set(key, nestedMap)
					} else {
						// 숫자나 다른 값들
						if i, err := strconv.ParseInt(value, 10, 64); err == nil {
							om.Set(key, i)
						} else if f, err := strconv.ParseFloat(value, 64); err == nil {
							om.Set(key, f)
						} else if value == "true" {
							om.Set(key, true)
						} else if value == "false" {
							om.Set(key, false)
						} else if value == "null" {
							om.Set(key, nil)
						} else {
							om.Set(key, value)
						}
					}
				}

				isInValue = false
				start = i + 1
			}
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
