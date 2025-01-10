package program

import (
	"encoding/json"
	"fmt"
	C "hello/pkg/core/config"
	"hello/pkg/core/model"
	"hello/pkg/core/network"
	"hello/pkg/core/storage"
	"hello/pkg/core/structure"
	"hello/pkg/core/vm"
	"hello/pkg/crypto"
	"hello/pkg/rpc"
	"hello/pkg/swift"
	"hello/pkg/util"
	"log"

	"github.com/spf13/cobra"
)

func CreateRequestCmd() *cobra.Command {
	var peer string
	var privateKey string
	var payload string

	_privateKey, _ := GetPrivateKey()

	cmd := &cobra.Command{
		Use:   "rawrequest",
		Short: "get raw request",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			data, err := structure.ParseOrderedMap(payload)
			address := crypto.GetAddress(crypto.GetXpub(privateKey))

			data.Set("timestamp", util.Utime())
			data.Set("from", address)

			if err != nil {
				log.Fatalf("Failed to parse payload: %v", err)
			}

			signedRequest := model.NewSignedRequestFromRawData(data, privateKey)
			payload, err = signedRequest.Ser()
			if err != nil {
				log.Fatalf("Failed to serialize signed request: %v", err)
			}

			req := rpc.CreateRequest(payload, peer)
			response, err := network.CallRawRequest(req)

			if err != nil {
				log.Fatalf("Failed to call raw request: %v", err)
			}

			rstr := FormatResponse(&response.Payload)
			fmt.Println(rstr)
		},
	}

	cmd.Flags().StringVarP(&peer, "peer", "p", "localhost:9001", "peer to get balance")
	cmd.Flags().StringVarP(&privateKey, "privatekey", "k", _privateKey, "private key to sign the request")
	cmd.Flags().StringVarP(&payload, "payload", "l", "", "payload to sign")
	cmd.MarkFlagRequired("payload")

	return cmd
}

func CreateListTransactionCmd() *cobra.Command {
	var peer string
	var count int

	cmd := &cobra.Command{
		Use:   "listtx",
		Short: "list transactions",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			payload := structure.NewOrderedMap()
			privateKey, err := GetPrivateKey()

			if err != nil {
				log.Fatalf("Failed to get private key: %v", err)
			}

			address := crypto.GetAddress(crypto.GetXpub(privateKey))

			payload.Set("type", "ListTransaction")
			payload.Set("address", address)
			payload.Set("timestamp", util.Utime())
			payload.Set("from", address)

			req := CreateWalletRequest(peer, payload.Ser())

			response, err := network.CallRawRequest(req)

			if err != nil {
				log.Fatalf("Failed to send request: %v", err)
			}

			rstr := FormatResponse(&response.Payload)

			fmt.Println(rstr)
		},
	}

	cmd.Flags().StringVarP(&peer, "peer", "p", "localhost:9001", "peer to get balance")
	cmd.Flags().IntVarP(&count, "count", "c", 1, "count to get balance")

	return cmd
}

func CreateListBlockCmd() *cobra.Command {
	var peer string
	var page int
	var count int

	cmd := &cobra.Command{
		Use:   "listblk",
		Short: "list transactions",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			payload := structure.NewOrderedMap()
			privateKey, err := GetPrivateKey()

			if err != nil {
				log.Fatalf("Failed to get private key: %v", err)
			}

			address := crypto.GetAddress(crypto.GetXpub(privateKey))

			payload.Set("type", "ListBlock")
			payload.Set("address", address)
			payload.Set("timestamp", util.Utime())
			payload.Set("from", address)

			req := CreateWalletRequest(peer, payload.Ser())

			response, err := network.CallRawRequest(req)

			if err != nil {
				log.Fatalf("Failed to send request: %v", err)
			}

			rstr := FormatResponse(&response.Payload)
			fmt.Println(rstr)
		},
	}

	cmd.Flags().StringVarP(&peer, "peer", "p", "localhost:9001", "peer to get balance")
	cmd.Flags().IntVarP(&page, "page", "a", 1, "page to get balance")
	cmd.Flags().IntVarP(&count, "count", "c", 1, "count to get balance")

	return cmd
}

func CreatePairInfoRequestCmd() *cobra.Command {
	var peer string
	var pair_address string
	var address_a string
	var address_b string

	cmd := &cobra.Command{
		Use:   "pairinfo",
		Short: "",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			payload := structure.NewOrderedMap()
			privateKey, err := GetPrivateKey()

			if err != nil {
				log.Fatalf("Failed to get private key: %v", err)
			}

			address := crypto.GetAddress(crypto.GetXpub(privateKey))

			payload.Set("type", "GetPairInfo")

			if pair_address != "" {
				payload.Set("pair_address", pair_address)
			} else {
				payload.Set("token_address_a", address_a)
				payload.Set("token_address_b", address_b)
			}

			payload.Set("timestamp", util.Utime())
			payload.Set("from", address)

			req := CreateWalletRequest(peer, payload.Ser())

			response, err := network.CallRawRequest(req)

			if err != nil {
				log.Fatalf("Failed to send request: %v", err)
			}

			rstr := FormatResponse(&response.Payload)
			fmt.Println(rstr)
		},
	}

	cmd.Flags().StringVarP(&peer, "peer", "p", "localhost:9001", "peer")
	cmd.Flags().StringVarP(&address_a, "address_a", "a", "", "")
	cmd.Flags().StringVarP(&address_b, "address_b", "b", "", "")
	cmd.Flags().StringVarP(&pair_address, "pair_address", "i", "", "pair address")

	return cmd
}

func CreateSearchCmd() *cobra.Command {
	var peer string
	var prefix string
	var page int
	var count int

	cmd := &cobra.Command{
		Use:   "search",
		Short: "search cli tool",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {

			packet := swift.Packet{
				Type:    swift.PacketTypeSearchRequest,
				Payload: []byte(fmt.Sprintf(`{"prefix": "%s", "page": %d, "count": %d}`, prefix, page, count)),
			}

			response, err := network.CallRPC(peer, packet)
			if err != nil {
				log.Fatalf("Failed to send request: %v", err)
			}

			rstr := FormatResponse(&response.Payload)
			fmt.Println(rstr)
		},
	}

	cmd.Flags().StringVarP(&peer, "peer", "p", "localhost:9001", "peer")
	cmd.Flags().StringVarP(&prefix, "prefix", "k", "", "prefix")
	cmd.Flags().IntVarP(&page, "page", "a", 1, "page")
	cmd.Flags().IntVarP(&count, "count", "c", 10, "count")

	cmd.MarkFlagRequired("prefix")

	return cmd
}

func CreateGetLastHeightCmd() *cobra.Command {
	var targetNode string

	cmd := &cobra.Command{
		Use:   "lastheight",
		Short: "get last height of the network.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			req := swift.Packet{
				Type:    swift.PacketTypeLastHeightRequest,
				Payload: json.RawMessage("{}"),
			}

			response, err := network.CallRPC(targetNode, req)
			if err != nil {
				log.Fatalf("RPC call failed: %v", err)
			}

			rst := FormatResponse(&response.Payload)
			fmt.Println(rst)

		},
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.Flags().Parse(args)
		},
	}

	cmd.Flags().StringVarP(&targetNode, "peer", "p", C.CLI_DEFAULT_REQUEST, "Target node to ping")

	return cmd
}

func CreateSyncCmd() *cobra.Command {
	const SYNC_BATCH_SIZE = 100
	const SYNC_BATCH_LIMIT = SYNC_BATCH_SIZE - 1

	var targetNode string
	var startHeight int
	var endHeight int

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Synchronize blocks from target node",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if startHeight == -1 {
				heightReq := swift.Packet{
					Type:    swift.PacketTypeLastHeightRequest,
					Payload: json.RawMessage("{}"),
				}

				response, err := network.CallRPC(targetNode, heightReq)
				if err != nil {
					log.Fatalf("Failed to get last height: %v", err)
				}

				var lastHeight int
				if err := json.Unmarshal(response.Payload, &lastHeight); err != nil {
					log.Fatalf("Failed to parse height: %v", err)
				}

				startHeight = 1
				endHeight = lastHeight
			}

			for currentStart := startHeight; currentStart <= endHeight; currentStart += SYNC_BATCH_SIZE {
				currentEnd := currentStart + SYNC_BATCH_LIMIT
				if currentEnd > endHeight {
					currentEnd = endHeight
				}

				syncReq := swift.Packet{
					Type: swift.PacketTypeSyncBlockRequest,
					Payload: json.RawMessage(fmt.Sprintf(`{
						"start_height": %d,
						"end_height": %d
					}`, currentStart, currentEnd)),
				}

				response, err := network.CallRPC(targetNode, syncReq)
				if err != nil {
					log.Fatalf("Failed to sync blocks (height %d-%d): %v",
						currentStart, currentEnd, err)
				}

				var blocks []string

				if err := json.Unmarshal(response.Payload, &blocks); err != nil {
					log.Fatalf("Failed to parse block data: %v", err)
				}

				machine := vm.GetMachineInstance()
				for _, block := range blocks {
					parsed, err := storage.ParseBlock([]byte(block))
					if err != nil {
						log.Fatalf("Failed to parse block data: %v", err)
					}
					if err := machine.Commit(parsed); err != nil {
						log.Fatalf("Failed to commit block (height %d): %v",
							parsed.Height, err)
					}
				}
				fmt.Printf("Height %d-%d synchronized\n", currentStart, currentEnd)
			}

			fmt.Printf("Full synchronization completed (height %d-%d)\n", startHeight, endHeight)
		},
	}

	cmd.Flags().StringVarP(&targetNode, "peer", "p", C.CLI_DEFAULT_REQUEST, "Target node to sync from")
	cmd.Flags().IntVarP(&startHeight, "start", "s", -1, "Start height (default: 0)")
	cmd.Flags().IntVarP(&endHeight, "end", "e", -1, "End height (default: latest block)")

	return cmd
}

func CreateApiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "api cli tool",
	}

	cmd.AddCommand(
		CreateRequestCmd(),
		CreateListBlockCmd(),
		CreateListTransactionCmd(),
		CreatePairInfoRequestCmd(),
		CreateSearchCmd(),
		CreateGetLastHeightCmd(),
		CreateSyncCmd(),
	)

	return cmd
}
