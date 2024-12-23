package structure

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type OrderedMap struct {
	m map[string]interface{}
	l []string
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		m: make(map[string]interface{}),
		l: []string{},
	}
}

func (om *OrderedMap) Set(key string, value interface{}) {
	if _, exists := om.m[key]; !exists {
		om.l = append(om.l, key)
	}
	om.m[key] = value
}

func (om *OrderedMap) Get(key string) (interface{}, bool) {
	value, exists := om.m[key]
	return value, exists
}

func (om *OrderedMap) Keys() []string {
	return om.l
}

// ParseOrderedMap parses a JSON string into an OrderedMap
func ParseOrderedMap(jsonStr string) (*OrderedMap, error) {
	i := 0
	return parseObject(&i, jsonStr)
}

// parseObject parses a JSON object into an OrderedMap
func parseObject(i *int, jsonStr string) (*OrderedMap, error) {
	skipWhitespace(i, jsonStr)

	if *i >= len(jsonStr) || jsonStr[*i] != '{' {
		return nil, errors.New(fmt.Sprintf("expected '{' at start of object,  %s", jsonStr))
	}
	*i++ // Skip '{'

	om := NewOrderedMap()
	for {
		skipWhitespace(i, jsonStr)
		if *i < len(jsonStr) && jsonStr[*i] == '}' {
			*i++ // Skip '}'
			break
		}

		// Read key
		key, err := readString(i, jsonStr)
		if err != nil {
			return nil, err
		}

		skipWhitespace(i, jsonStr)
		if *i >= len(jsonStr) || jsonStr[*i] != ':' {
			return nil, errors.New("expected ':' after key")
		}
		*i++ // Skip ':'

		// Read value
		value, err := parseValue(i, jsonStr)
		if err != nil {
			return nil, err
		}
		om.Set(key, value)

		skipWhitespace(i, jsonStr)
		if *i < len(jsonStr) && jsonStr[*i] == ',' {
			*i++ // Skip ','
		} else if *i < len(jsonStr) && jsonStr[*i] == '}' {
			continue // Handle object end in the next iteration
		} else if *i < len(jsonStr) {
			return nil, errors.New("expected ',' or '}' in object")
		}
	}
	return om, nil
}

// parseValue parses a JSON value (string, number, object, boolean, null)
func parseValue(i *int, jsonStr string) (interface{}, error) {
	skipWhitespace(i, jsonStr)
	if *i >= len(jsonStr) {
		return nil, errors.New("unexpected end of JSON")
	}

	switch jsonStr[*i] {
	case '"':
		return readString(i, jsonStr)
	case '{':
		return parseObject(i, jsonStr)
	case '[':
        	return parseArray(i, jsonStr) // array
	case 't': // true
		if strings.HasPrefix(jsonStr[*i:], "true") {
			*i += 4
			return true, nil
		}
	case 'f': // false
		if strings.HasPrefix(jsonStr[*i:], "false") {
			*i += 5
			return false, nil
		}
	case 'n': // null
		if strings.HasPrefix(jsonStr[*i:], "null") {
			*i += 4
			return nil, nil
		}
	default:
		if jsonStr[*i] == '-' || unicode.IsDigit(rune(jsonStr[*i])) {
			return readNumber(i, jsonStr)
		}
	}

	return nil, errors.New(fmt.Sprintf("unexpected value type in JSON %s", jsonStr[*i]))
}

func readString(i *int, jsonStr string) (string, error) {
	if *i >= len(jsonStr) || jsonStr[*i] != '"' {
		return "", errors.New("expected '\"' at start of string")
	}
	*i++ // Skip the opening quote

	var result strings.Builder
	for *i < len(jsonStr) {
		ch := jsonStr[*i]
		*i++
		if ch == '"' {
			return result.String(), nil
		} else if ch == '\\' {
			// Handle escaped characters
			if *i >= len(jsonStr) {
				return "", errors.New("unexpected end of string escape")
			}
			escaped := jsonStr[*i]
			*i++
			switch escaped {
			case '"', '\\', '/':
				result.WriteByte(escaped)
			case 'b':
				result.WriteByte('\b')
			case 'f':
				result.WriteByte('\f')
			case 'n':
				result.WriteByte('\n')
			case 'r':
				result.WriteByte('\r')
			case 't':
				result.WriteByte('\t')
			default:
				return "", errors.New(fmt.Sprintf("invalid escape character in string %s", jsonStr))
			}
		} else {
			result.WriteByte(ch)
		}
	}
	return "", errors.New("unexpected end of string")
}

// readNumber parses a JSON number (int64 or float64)
func readNumber(i *int, jsonStr string) (interface{}, error) {
	start := *i
	for *i < len(jsonStr) && (jsonStr[*i] == '-' || jsonStr[*i] == '.' || unicode.IsDigit(rune(jsonStr[*i]))) {
		*i++
	}
	numStr := jsonStr[start:*i]
	if strings.Contains(numStr, ".") {
		// Parse as float64
		floatValue, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return nil, err
		}
		return floatValue, nil
	}
	// Parse as int64
	intValue, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return nil, err
	}
	return intValue, nil
}

// skipWhitespace skips any whitespace characters
func skipWhitespace(i *int, jsonStr string) {
	for *i < len(jsonStr) && unicode.IsSpace(rune(jsonStr[*i])) {
		*i++
	}
}

// escapeString escapes special characters in a string for JSON
func escapeString(s string) string {
	var builder strings.Builder
	for _, ch := range s {
		switch ch {
		case '"':
			builder.WriteString(`\"`)
		case '\\':
			builder.WriteString(`\\`)
		case '\b':
			builder.WriteString(`\b`)
		case '\f':
			builder.WriteString(`\f`)
		case '\n':
			builder.WriteString(`\n`)
		case '\r':
			builder.WriteString(`\r`)
		case '\t':
			builder.WriteString(`\t`)
		default:
			builder.WriteRune(ch)
		}
	}
	return builder.String()
}

func (om *OrderedMap) Ser() string {
	var result strings.Builder
	err := serialize(&result, om)
	fmt.Println("err: ", err)
	return result.String()
}

func serialize(builder *strings.Builder, value interface{}) error {
	switch v := value.(type) {
	case *OrderedMap:
		if v == nil {
			builder.WriteString("null")
			return nil
		}
		return serialize(builder, *v) // Dereference and serialize
	case OrderedMap:
		builder.WriteString("{")
		first := true
		for _, key := range v.Keys() {
			if !first {
				builder.WriteString(",")
			}
			first = false
			// Serialize the key
			builder.WriteString(`"`)
			builder.WriteString(escapeString(key))
			builder.WriteString(`":`)
			// Serialize the value
			val := v.m[key]
			if err := serialize(builder, val); err != nil {
				return err
			}
		}
		builder.WriteString("}")
	case string:
		builder.WriteString(`"`)
		builder.WriteString(escapeString(v))
		builder.WriteString(`"`)
	case int64, float64, bool, int, float32, uint64, uint32, uint16, uint8, int16, int8, uint:
		builder.WriteString(fmt.Sprintf("%v", v))
	case *int64:
		builder.WriteString(fmt.Sprintf("%v", *v))
	case *float64:
		builder.WriteString(fmt.Sprintf("%v", *v))
	case *string:
		builder.WriteString(`"`)
		builder.WriteString(escapeString(*v))
		builder.WriteString(`"`)
	case *bool:
		builder.WriteString(fmt.Sprintf("%v", *v))
	case nil:
		builder.WriteString("null")
	default:
		panic(fmt.Sprintf("unsupported value type for serialization: %T", v))
		return fmt.Errorf("unsupported value type for serialization: %T", v)
	}
	return nil
}

func parseArray(i *int, jsonStr string) (interface{}, error) {
    if jsonStr[*i] != '[' {
        return nil, errors.New("expected '[' at the beginning of array")
    }
    *i++ // '[' skipping character.
    var array []interface{}

    for *i < len(jsonStr) {
        skipWhitespace(i, jsonStr)
        if jsonStr[*i] == ']' {
            *i++ // ']' skipping character.
            return array, nil
        }

        value, err := parseValue(i, jsonStr)
        if err != nil {
            return nil, err
        }
        array = append(array, value)

        skipWhitespace(i, jsonStr)
        if jsonStr[*i] == ',' {
            *i++ // ',' skipping character.
        } else if jsonStr[*i] != ']' {
            return nil, errors.New("expected ',' or ']' in array")
        }
    }
    return nil, errors.New("unexpected end of JSON while parsing array")
}
