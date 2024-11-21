package rpc

import (
	"encoding/json"
	"hello/pkg/core/abi"
	. "hello/pkg/core/config"
	. "hello/pkg/core/debug"
	"hello/pkg/core/model"
	. "hello/pkg/core/model"
	. "hello/pkg/core/storage"
	. "hello/pkg/util"
	"math"
	"strings"
)

type Code struct{}

var SystemMethods = []string{"Genesis", "Register", "Grant", "Revoke", "Oracle", "Faucet", "Publish", "Send", "Submit"}

func (c *Code) Contracts() map[string]interface{} {
	storage := GetStatusFileInstance()

	methods := make(map[string]map[string]*Method)
	var postProcess *Method

	systemContracts := make([]string, len(SystemMethods))
	for i, method := range SystemMethods {
		systemContracts[i] = strings.ToLower(method)
	}
	DebugLog("system contracts member::", systemContracts)

	customContracts := []string{}
	customCount := storage.CountLocalStatus(config.ContractPrefix())

	count := 50
	page := int(math.Ceil(float64(customCount) / float64(count)))

	for i := 0; i < page; i++ {
		codes := status.GetInstance().ListLocalStatus(config.ContractPrefix(), i, count)
		customContracts = append(customContracts, codes...)
	}

	logger.Log("Custom contracts:", customContracts)

	// Custom contracts
	for _, code := range customContracts {
		var codeMap map[string]interface{}
		if err := json.Unmarshal([]byte(code), &codeMap); err != nil {
			continue
		}
		logger.Log("parse code:", codeMap)

		method := c.ContractToMethod(codeMap)
		logger.Log("parsed method:", method)

		if method != nil {
			spaceID := hasher.SpaceID(method.Writer(), method.Space())

			name := method.Name()
			w := method.Writer()
			s := method.Space()
			cc := method.CID()
			logger.Log("load custom contract " + name + ": writer=" + w + "; space=" + s + "; cid=" + cc)

			if spaceID == config.RootSpaceID() && name == "Fee" {
				postProcess = method
				continue
			}

			cid := method.CID()
			if methods[cid] == nil {
				methods[cid] = make(map[string]*model.Method)
			}
			methods[cid][name] = method
		}
	}

	logger.Log("Load system contracts:", SystemMethods)
	// System contracts
	for _, code := range systemContracts {
		method := systemcontract.GetMethod(code)
		cid := method.CID()
		name := method.Name()

		w := method.Writer()
		ci := method.Space()
		logger.Log("load system contract " + name + ": writer=" + w + "; space=" + ci + "; cid=" + cid)

		if methods[cid] == nil {
			methods[cid] = make(map[string]*model.Method)
		}
		methods[cid][name] = method
	}

	return map[string]interface{}{
		"methods":     methods,
		"postProcess": postProcess,
	}
}

func (c *Code) Requests() map[string]map[string]*model.Method {
	methods := make(map[string]map[string]*model.Method)

	systemRequests := systemrequest.GetMethods()
	customRequests := []string{}
	customCount := status.GetInstance().CountLocalStatus(config.RequestPrefix())

	count := 50
	page := int(math.Ceil(float64(customCount) / float64(count)))

	for i := 0; i < page; i++ {
		codes := status.GetInstance().ListLocalStatus(config.RequestPrefix(), i, count)
		customRequests = append(customRequests, codes...)
	}

	// Custom requests
	for _, code := range customRequests {
		var codeMap map[string]interface{}
		if err := json.Unmarshal([]byte(code), &codeMap); err != nil {
			continue
		}

		method := c.RequestToMethod(codeMap)
		if method != nil {
			cid := method.CID()
			name := method.Name()

			if methods[cid] == nil {
				methods[cid] = make(map[string]*model.Method)
			}
			methods[cid][name] = method
		}
	}

	// System requests
	for _, code := range systemRequests {
		method := systemrequest.GetMethod(code)
		cid := method.CID()
		name := method.Name()

		if methods[cid] == nil {
			methods[cid] = make(map[string]*model.Method)
		}
		methods[cid][name] = method
	}

	return methods
}

func (c *Code) Contract(name string, cid string) *model.Method {
	if cid == "" {
		cid = config.RootSpaceID()
	}

	codes := c.Contracts()
	contracts := codes["methods"].(map[string]map[string]*model.Method)
	return contracts[cid][name]
}

func (c *Code) Request(name string, cid string) *model.Method {
	if cid == "" {
		cid = config.RootSpaceID()
	}

	requests := c.Requests()
	return requests[cid][name]
}

func (c *Code) ContractToMethod(code map[string]interface{}) *model.Method {
	if code == nil {
		return nil
	}

	// Check if code uses new format
	if _, ok := code["t"]; ok {
		return model.NewMethod(code)
	}

	// Convert legacy format
	newCode := map[string]interface{}{
		"t": code["type"].(string),
		"m": "0.2.0",
		"n": code["name"].(string),
		"v": code["version"].(string),
		"s": code["nonce"].(string),
		"w": code["writer"].(string),
		"p": code["parameters"].([]interface{}),
		"e": []interface{}{},
	}

	logger.Log("parsed new code header:", newCode)

	conditions := code["conditions"].([]interface{})
	updates := code["updates"].([]interface{})

	for _, condition := range conditions {
		cond := condition.([]interface{})
		logic := cond[0]
		errMsg := "Conditional error"
		if len(cond) > 1 {
			errMsg = cond[1].(string)
		}
		newCode["e"] = append(newCode["e"].([]interface{}), abi.Condition(logic, errMsg))
	}

	for _, update := range updates {
		newCode["e"] = append(newCode["e"].([]interface{}), update)
	}

	return model.NewMethod(newCode)
}

func (c *Code) RequestToMethod(code map[string]interface{}) *model.Method {
	if code == nil {
		return nil
	}

	// Check if code uses new format
	if _, ok := code["t"]; ok {
		return model.NewMethod(code)
	}

	// Convert legacy format
	newCode := map[string]interface{}{
		"t": code["type"].(string),
		"m": "0.2.0",
		"n": code["name"].(string),
		"v": code["version"].(string),
		"s": code["nonce"].(string),
		"w": code["writer"].(string),
		"p": code["parameters"].([]interface{}),
		"e": []interface{}{},
	}

	conditions := code["conditions"].([]interface{})
	response := code["response"].(interface{})

	for _, condition := range conditions {
		cond := condition.([]interface{})
		logic := cond[0]
		errMsg := "Conditional error"
		if len(cond) > 1 {
			errMsg = cond[1].(string)
		}
		newCode["e"] = append(newCode["e"].([]interface{}), abi.Condition(logic, errMsg))
	}

	newCode["e"] = append(newCode["e"].([]interface{}), abi.Response(response))

	return model.NewMethod(newCode)
}
