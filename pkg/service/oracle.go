package service

import (
	"context"
	"encoding/json"
	"fmt"
	C "hello/pkg/core/config"
	. "hello/pkg/core/debug"
	"hello/pkg/core/model"
	"hello/pkg/core/storage"
	"hello/pkg/core/structure"
	"hello/pkg/core/vm"
	"hello/pkg/swift"
)

func OracleLog(format string, args ...any) {
	fmt.Printf(format, args...)
}

type Oracle struct {
	swift        *swift.Server
	machine      *vm.Machine
	storage      *storage.StatusFile
	storageIndex *storage.StatusIndex
	mempool      *storage.MempoolStorage
}

var oracleInstance *Oracle

func GetOracleService() *Oracle {
	if oracleInstance == nil {
		oracleInstance = &Oracle{
			machine:      vm.GetMachineInstance(),
			mempool:      storage.GetMempoolInstance(),
			swift:        nil,
			storage:      storage.GetStatusFileInstance(),
			storageIndex: storage.GetStatusIndexInstance(),
		}
	}
	return oracleInstance
}

func (o *Oracle) OnStartUp(config swift.SecurityConfig) error {
	//TEMP TODO
	C.CORE_TEST_MODE = true
	C.DATA_TEST_ROOT_DIR = "genesis_test_2"

	o.storageIndex.Load()

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

	// broadcast transaction request
	o.swift.RegisterHandler(swift.PacketTypeBroadcastTransactionRequest, func(ctx context.Context, packet *swift.Packet) error {
		swift.SwiftInfoLog("broadcast transaction request: %s", string(packet.Payload))

		var tx model.SignedTransaction
		data, err := structure.ParseOrderedMap(string(packet.Payload))

		if err != nil {
			return err
		}
		t, ok := data.Get("transaction")
		if !ok {
			return fmt.Errorf("transaction data not found")
		}
		DebugLog("\n\nparsed data: %s\n\n", t)

		tx, err = model.NewSignedTransaction(data)

		if err != nil {
			return err
		}

		txHash, err := tx.GetTxHash()

		if err != nil {
			return err
		}

		OracleLog("successfully added transaction to mempool: %s", txHash)
		swift.SwiftInfoLog("successfully added transaction to mempool: %s", txHash)

		err = o.mempool.AddTransaction(&tx)
		if err != nil {
			return err
		}

		responseData, err := json.Marshal(map[string]bool{"ok": true})
		if err != nil {
			return err
		}

		response := &swift.Packet{
			Type:    swift.PacketTypeBroadcastTransactionResponse,
			Payload: json.RawMessage(responseData),
		}

		return o.swift.Send(ctx, response)
	})

	// list mempool transaction request
	o.swift.RegisterHandler(swift.PacketTypeListMempoolTransactionRequest, func(ctx context.Context, packet *swift.Packet) error {
		DebugLog("list mempool transaction request")

		data, err := json.Marshal(o.mempool.FormatTransactions())
		if err != nil {
			return err
		}

		response := &swift.Packet{
			Type:    swift.PacketTypeListMempoolTransactionResponse,
			Payload: json.RawMessage(data),
		}

		return o.swift.Send(ctx, response)
	})

	o.swift.RegisterHandler(swift.PacketTypeGetStatusBundleRequest, func(ctx context.Context, packet *swift.Packet) error {
		OracleLog("get status bundle request for key %s", string(packet.Payload))

		bundle, ok := o.storage.CachedUniversalIndexes[string(packet.Payload)]
		if !ok {
			return fmt.Errorf("status bundle not found")
		}

		data, err := json.Marshal(bundle)
		if err != nil {
			return err
		}
		response := &swift.Packet{
			Type:    swift.PacketTypeGetStatusBundleResponse,
			Payload: json.RawMessage(data),
		}
		return o.swift.Send(ctx, response)
	})

}

func GetEpoch() string {
	return vm.GetMachineInstance().Epoch()
}
