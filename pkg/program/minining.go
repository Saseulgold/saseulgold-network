package program

import (
	"fmt"
	"hello/pkg/core/network"
	"hello/pkg/crypto"
	"hello/pkg/util"
	"log"
	"hello/pkg/core/structure"
	"github.com/spf13/cobra"
)

// FormatResponse(payload *json.RawMessage) string {
// CreateWalletTransaction(peer string, payload string) *rpc.TransactionRequest {
// CreateWalletRequest(peer string, payload string) *rpc.RawRequest {
func CreateSubmitMiningCmd() *cobra.Command {
	var peer string
	var address string

	privateKey, err := GetPrivateKey()
	_address := crypto.GetAddress(crypto.GetXpub(privateKey))

	cmd := &cobra.Command{
		Use:   "submit",
		Short: "submit hash for mining",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			payload := structure.NewOrderedMap()
			fmt.Println("address: ", address)

			if err != nil {
				log.Fatalf("Failed to get private key: %v", err)
			}


			payload.Set("type", "Mining")
			payload.Set("timestamp", util.Utime())
			payload.Set("from", address)

			req := CreateWalletRequest(peer, payload.Ser())

			response, err := network.CallRawRequest(req)
			if err != nil {
				log.Fatalf("Failed to send request: %v", err)
			}

			rstr := FormatResponse(&response.Payload)
			frstr := util.DivideByE18(rstr)
			fmt.Println("balance: ", frstr)
		},
	}

	cmd.Flags().StringVarP(&peer, "peer", "p", "localhost:9001", "peer to get balance")
	cmd.Flags().StringVarP(&address, "address", "a", _address, "peer to get balance")

	return cmd
}


func CreateMiningCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mining",
		Short: "mining cli tool",
	}

	cmd.AddCommand(
		CreateSubmitMiningCmd(),
	)

	return cmd
}
