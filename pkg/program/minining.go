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

	var epoch string
	var nonce string
	var chash string

	privateKey, _ := GetPrivateKey()
	_address := crypto.GetAddress(crypto.GetXpub(privateKey))

	cmd := &cobra.Command{
		Use:   "submit",
		Short: "submit hash for mining",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			payload := structure.NewOrderedMap()
			privateKey, err := GetPrivateKey()

			if err != nil {
				log.Fatalf("Failed to get private key: %v", err)
			}

			address := crypto.GetAddress(crypto.GetXpub(privateKey))

			payload.Set("type", "Mining")
			payload.Set("epoch", epoch)
			payload.Set("nonce", nonce)
			payload.Set("calculated_hash", chash)
			payload.Set("timestamp", util.Utime())
			payload.Set("from", address)

			req := CreateWalletTransaction(peer, payload.Ser())
			fmt.Println(req.Payload)

			response, err := network.CallTransactionRequest(req)
			if err != nil {
				log.Fatalf("Failed to send request: %v", err)
			}

			rstr := FormatResponse(&response.Payload)
			fmt.Println(rstr)
		},
	}

	cmd.Flags().StringVarP(&peer, "peer", "p", "localhost:9001", "peer to get balance")
	cmd.Flags().StringVarP(&address, "address", "a", _address, "peer to get balance")
	cmd.Flags().StringVarP(&epoch, "epoch", "e", "", "epoch")
	cmd.Flags().StringVarP(&nonce, "nonce", "n", "", "nonce")
	cmd.Flags().StringVarP(&chash, "chash", "c", "", "nonce")

	cmd.MarkFlagRequired("epoch")
	cmd.MarkFlagRequired("nonce")
	cmd.MarkFlagRequired("chash")

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
