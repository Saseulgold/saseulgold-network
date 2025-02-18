package service

import (
	. "hello/pkg/core/model"
	"hello/pkg/core/storage"
	. "hello/pkg/core/vm"
)

// static request is an interface for request when service is not running, from storage file.
// every time staticrequest called, status storage is inialized and re-computed.
func StaticRequest(request *SignedRequest) (interface{}, error) {
	oracle := GetOracleService()
	si := storage.GetStatusIndexInstance()

	sf := oracle.storage
	ci := oracle.chain

	sf.Reset()
	sf.Touch()
	ci.Touch()
	si.Load()

	machine := GetMachineInstance()
	response, err := machine.Response(*request)

	if err != nil {
		return nil, err
	}

	return response, nil
}
