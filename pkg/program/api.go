package program

import (
	"fmt"
	"hello/pkg/core/model"
	"hello/pkg/core/network"
	"hello/pkg/core/structure"
	"hello/pkg/crypto"
	"hello/pkg/rpc"
	"hello/pkg/swift"
	"hello/pkg/util"
	"log"

	"github.com/spf13/cobra"
)

func CreateRequestCmd() *cobra.Command {
	var requestType string
	var peer string
	var privateKey string
	var payload string
	var pubKey string

	cmd := &cobra.Command{
		Use:   "rawrequest",
		Short: "get raw request",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {

			data, err := structure.ParseOrderedMap(payload)

			if err != nil {
				log.Fatalf("Failed to parse payload: %v", err)
			}

			signedRequest := model.NewSignedRequestFromRawData(data, privateKey)
			payload, err = signedRequest.Ser()
			if err != nil {
				log.Fatalf("Failed to serialize signed request: %v", err)
			}
			fmt.Println(payload)

			req := rpc.CreateRequest(payload, peer)
			response, err := network.CallRawRequest(req)

			if err != nil {
				log.Fatalf("Failed to call raw request: %v", err)
			}

			fmt.Println(response)
		},
	}

	cmd.Flags().StringVarP(&requestType, "type", "t", "", "type of raw request")
	cmd.Flags().StringVarP(&peer, "peer", "p", "", "peer to get balance")
	cmd.Flags().StringVarP(&privateKey, "privatekey", "k", "", "private key to sign the request")
	cmd.Flags().StringVarP(&payload, "payload", "l", "", "payload to sign")
	cmd.Flags().StringVarP(&pubKey, "pubkey", "u", "", "public key to sign the request")
	cmd.MarkFlagRequired("peer")

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
			payload.Set("token_address_a", address_a)
			payload.Set("token_address_b", address_b)

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
	cmd.MarkFlagRequired("address_a")
	cmd.Flags().StringVarP(&address_b, "address_b", "b", "", "")
	cmd.MarkFlagRequired("address_b")

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

			fmt.Println(peer)

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
	)

	return cmd
}
