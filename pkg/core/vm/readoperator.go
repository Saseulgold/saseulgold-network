package vm

func OpLoadParam(i *Interpreter, vars interface{}) interface{} {
	var result interface{}

	if arr, ok := vars.([]interface{}); ok {
		for _, v := range arr {
			str, ok := v.(string)
			if !ok {
				return nil
			}

			if result == nil {
				result = i.SignedData.GetAttribute(str)
			} else if arr, ok := result.(map[string]interface{}); ok {
				if val, exists := arr[str]; exists {
					result = val
				} else {
					break
				}
			} else {
				break
			}
		}
	}

	return result
}
