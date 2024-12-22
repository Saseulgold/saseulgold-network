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
	"hello/pkg/util"
	"os"
	"time"
)

func OracleLog(format string, args ...any) {
	fmt.Printf(format, args...)
	fmt.Println()
}

type Oracle struct {
	swift        *swift.Server
	machine      *vm.Machine
	storage      *storage.StatusFile
	storageIndex *storage.StatusIndex
	chain        *storage.ChainStorage
	mempool      *storage.MempoolStorage
}

var oracleInstance *Oracle

func GetOracleService() *Oracle {
	if oracleInstance == nil {
		oracleInstance = &Oracle{
			machine:      vm.GetMachineInstance(),
			mempool:      storage.GetMempoolInstance(),
			chain:        storage.GetChainStorageInstance(),
			swift:        nil,
			storage:      storage.GetStatusFileInstance(),
			storageIndex: storage.GetStatusIndexInstance(),
		}
	}
	return oracleInstance
}

func (o *Oracle) RemoveCommittedTransactions(txs map[string]*model.SignedTransaction) {
	for _, tx := range txs {
		o.mempool.RemoveTransaction(tx.GetTxHash())
	}
}

func (o *Oracle) Consensus(txs map[string]*model.SignedTransaction) {

}

func (o *Oracle) Commit(txs map[string]*model.SignedTransaction) error {
	var previousBlockhash string

	lastBlockHeight := storage.LastHeight()
	previousBlock, err := storage.GetChainStorageInstance().GetBlock(lastBlockHeight)

	if err != nil {
		return err
	}

	fmt.Println("previousBlock: ", previousBlock)

	if previousBlock == nil {
		previousBlockhash = ""
	} else {
		previousBlockhash = previousBlock.BlockHash()
	}

	o.machine.Init(previousBlock, int64(util.Utime()))
	o.machine.SetTransactions(txs)
	o.machine.PreCommit()

	universals := o.machine.GetInterpreter().GetUniversals()
	DebugLog("universals: %v", universals)

	block := model.NewBlock(storage.LastHeight()+1, previousBlockhash)
	block.SetTimestamp(int64(util.Utime()))

	expectedBlock := o.machine.ExpectedBlock()

	if expectedBlock.GetTransactionCount() == 0 {
		DebugLog("no transactions to commit. invalid block.")
		return fmt.Errorf("no transactions to commit. invalid block.")
	}

	DebugLog("unv: %v", expectedBlock.UniversalUpdates)
	DebugLog("loc: %v", expectedBlock.LocalUpdates)

	err = o.machine.Commit(expectedBlock)

	if err != nil {
		DebugLog("Commit error: %v", err)
		return fmt.Errorf("Commit error: %v", err)
	}

	DebugLog("Commit success")
	return nil

}

func (o *Oracle) Run() error {
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()
	var epoch string = o.machine.Epoch()

	for {
		select {
		case <-ticker.C:
			epoch = o.machine.Epoch()
			fmt.Println("epoch: ", epoch)
			switch epoch {

			case "txtime":
				OracleLog("Validating transactions in mempool during transaction time")

			case "blocktime":
				transactions := o.mempool.GetTransactionsHashMap()
				if len(transactions) == 0 {
					OracleLog("No transactions to commit, skipping")
					continue
				}

				err := o.Commit(transactions)

				if err != nil {
					o.RemoveCommittedTransactions(transactions)
					OracleLog("failed to commit transactions: %v", err)
				}

			}
		}
	}
}

func (o *Oracle) validatePendingTransactions() error {
	transactions := o.mempool.GetTransactions(20)

	for _, tx := range transactions {
		valid, err := o.machine.TxValidity(tx)
		if err != nil {
			OracleLog("Transaction validation failed (%s): %v\n", tx.GetTxHash(), err)
			o.mempool.RemoveTransaction(tx.GetTxHash())
			continue
		}

		if !valid {
			OracleLog("Removing invalid transaction: %s\n", tx.GetTxHash())
			o.mempool.RemoveTransaction(tx.GetTxHash())
		}
	}

	return nil
}

func (o *Oracle) createAndBroadcastBlock() error {
	if !o.machine.IsInBlockTime() {
		return nil
	}

	// Create new block
	expectedBlock := o.machine.ExpectedBlock()
	if expectedBlock == nil {
		return fmt.Errorf("failed to create block")
	}

	// Commit block
	if err := o.machine.Commit(expectedBlock); err != nil {
		return fmt.Errorf("failed to commit block: %v", err)
	}

	// Serialize block to JSON
	blockData, err := json.Marshal(expectedBlock)
	if err != nil {
		return fmt.Errorf("failed to serialize block: %v", err)
	}

	// Broadcast block to network
	packet := &swift.Packet{
		Type:    swift.PacketTypeBroadcastBlockRequest,
		Payload: blockData,
	}

	if err := o.swift.Broadcast(context.Background(), packet); err != nil {
		return fmt.Errorf("failed to broadcast block: %v", err)
	}

	OracleLog("Successfully created and broadcast new block (height: %d)\n", expectedBlock.Height)
	return nil
}

func (o *Oracle) OnStartUp(config swift.SecurityConfig, port string) error {
	isRunning := util.ServiceIsRunning(storage.DataRootDir(), "oracle")
	if isRunning {
		return fmt.Errorf("oracle is already running")
	}

	o.chain.Touch()
	o.storage.Touch()

	o.storage.Reset()
	o.storageIndex.Load()

	o.swift = swift.NewServer(C.SWIFT_HOST+":"+port, config)
	o.registerPacketHandlers()
	err := util.ProcessStart(storage.DataRootDir(), "oracle", os.Getpid())
	if err != nil {
		OracleLog("failed to start oracle: %v", err)
		return err
	}

	return nil
}

func (o *Oracle) Start() error {
	go o.Run()
	err := o.swift.Start()
	if err != nil {
		return err
	}

	return nil
}

func (o *Oracle) registerPacketHandlers() {

	// send trasaction a node
	o.swift.RegisterHandler(swift.PacketTypeSendTransactionRequest, func(ctx context.Context, packet *swift.Packet) error {
		swift.SwiftInfoLog("send transaction request: %s", string(packet.Payload))

		data, err := structure.ParseOrderedMap(string(packet.Payload))
		if err != nil {
			return err
		}

		tx, err := model.NewSignedTransaction(data)
		if err != nil {
			return o.swift.SendErrorResponse(ctx, err.Error())
		}

		txHash := tx.GetTxHash()
		if err != nil {
			return o.swift.SendErrorResponse(ctx, err.Error())
		}

		err = o.mempool.AddTransaction(&tx)
		if err != nil {
			return o.swift.SendErrorResponse(ctx, err.Error())
		}

		err = o.swift.Send(ctx, packet)
		if err != nil {
			return o.swift.SendErrorResponse(ctx, err.Error())
		}

		response := &swift.Packet{
			Type:    swift.PacketTypeSendTransactionResponse,
			Payload: json.RawMessage(fmt.Sprintf("%v", txHash)),
		}

		return o.swift.Send(ctx, response)
	})

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

		valid, err := o.machine.TxValidity(&tx)

		if err != nil {
			return o.swift.SendErrorResponse(ctx, err.Error())
		}

		if !valid {
			return fmt.Errorf("transaction is not valid")
		}

		err = o.mempool.AddTransaction(&tx)
		if err != nil {
			return err
		}

		txHash := tx.GetTxHash()
		OracleLog("successfully added transaction to mempool: %s", txHash)

		var responseData []byte

		swift.SwiftInfoLog("broadcasting transaction to peers: %s", o.swift.GetPeers())

		if err := o.swift.Broadcast(ctx, packet); err != nil {
			return o.swift.SendErrorResponse(ctx, err.Error())
		}

		responseData, err = json.Marshal(map[string]string{"ok": "true", "msg": "successfully broadcast transaction"})

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

		var payload map[string]string
		err := json.Unmarshal(packet.Payload, &payload)
		if err != nil {
			return err
		}

		cursor, ok := o.storage.CachedUniversalIndexes[payload["key"]]
		if !ok {
			return fmt.Errorf("status bundle not found")
		}

		data, err := o.storage.ReadUniversalStatus(cursor)
		if err != nil {
			return err
		}

		response := &swift.Packet{
			Type:    swift.PacketTypeGetStatusBundleResponse,
			Payload: json.RawMessage(fmt.Sprintf("%v", data)),
		}
		return o.swift.Send(ctx, response)
	})

	// handshake request handler
	o.swift.RegisterHandler(swift.PacketTypeHandshakeCMDRequest, func(ctx context.Context, packet *swift.Packet) error {
		var payload struct {
			Peer string `json:"peer"`
		}

		swift.SwiftInfoLog("handshake request: %s", string(packet.Payload))
		if err := json.Unmarshal(packet.Payload, &payload); err != nil {
			return fmt.Errorf("handshake request payload parsing failed: %v", err)
		}

		if err := o.swift.Connect(payload.Peer); err != nil {
			return fmt.Errorf("failed to connect to peer: %v", err)
		}

		responseData, err := json.Marshal(map[string]string{"status": "connected"})
		if err != nil {
			return fmt.Errorf("failed to create response: %v", err)
		}

		swift.SwiftInfoLog("handshake response: %s", string(responseData))

		response := &swift.Packet{
			Type:    swift.PacketTypeHandshakeCMDResponse,
			Payload: responseData,
		}

		return o.swift.Send(ctx, response)
	})
}

func (o *Oracle) Shutdown() error {
	return util.TerminateProcess(storage.DataRootDir(), "oracle")
}
