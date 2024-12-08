package service

import (
	"context"
	"encoding/json"
	C "hello/pkg/core/config"
	. "hello/pkg/core/debug"
	"hello/pkg/core/storage"
	"hello/pkg/core/vm"
	"hello/pkg/swift"
)

type Oracle struct {
	swift   *swift.Server
	machine *vm.Machine
	storage *storage.StatusFile
	mempool *storage.MempoolStorage
}

var oracleInstance *Oracle

func GetOracleService() *Oracle {
	if oracleInstance == nil {
		oracleInstance = &Oracle{
			machine: vm.GetMachineInstance(),
			mempool: storage.GetMempoolInstance(),
			swift:   nil,
			storage: storage.GetStatusFileInstance(),
		}
	}
	return oracleInstance
}

func (o *Oracle) OnStartUp(config swift.SecurityConfig) error {
	if err := o.storage.Cache(); err != nil {
		return err
	}

	o.swift = swift.NewServer(C.SWIFT_HOST+":"+C.SWIFT_PORT, config)
	o.registerPacketHandlers()

	return nil
}

func (o *Oracle) Start() error {
	err := o.swift.Start()
	if err != nil {
		return err
	}
	return nil
}

func (o *Oracle) registerPacketHandlers() {
	o.swift.RegisterHandler(swift.PacketTypeListMempoolTransactionRequest, func(ctx context.Context, packet *swift.Packet) error {
		DebugLog("list mempool transaction request")

		response := &swift.Packet{
			Type:    swift.PacketTypeListMempoolTransactionResponse,
			Payload: json.RawMessage(o.mempool.FormatTransactions()),
		}
		return o.swift.Send(ctx, response)
	})

}

func GetEpoch() string {
	return vm.GetMachineInstance().Epoch()
}
