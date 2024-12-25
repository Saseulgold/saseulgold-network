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

func CreateMintTokenCmd() *cobra.Command {
	var peer string
	var name string
	var symbol string
	var supply string

	cmd := &cobra.Command{
		Use:   "mint",
		Short: "mint new token",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			payload := structure.NewOrderedMap()
			privateKey, err := GetPrivateKey()

			if err != nil {
				log.Fatalf("Failed to get private key: %v", err)
			}

			address := crypto.GetAddress(crypto.GetXpub(privateKey))

			payload.Set("type", "Mint")
			payload.Set("name", name)
			payload.Set("symbol", symbol)
			payload.Set("supply", supply)
			payload.Set("timestamp", util.Utime())
			payload.Set("from", address)
			fmt.Println(payload.Ser())

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

	cmd.Flags().StringVarP(&peer, "peer", "p", "localhost:9001", "peer to connect")
	cmd.Flags().StringVarP(&name, "name", "n", "", "token name")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(&symbol, "symbol", "s", "", "token symbol")
	cmd.MarkFlagRequired("symbol")
	cmd.Flags().StringVarP(&supply, "supply", "a", "", "token supply amount")
	cmd.MarkFlagRequired("supply")

	return cmd
}

func CreateTransferTokenCmd() *cobra.Command {
	var peer string
	var symbol string
	var to string
	var amount string

	cmd := &cobra.Command{
		Use:   "transfer",
		Short: "transfer token to another address",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			payload := structure.NewOrderedMap()
			privateKey, err := GetPrivateKey()

			if err != nil {
				log.Fatalf("Failed to get private key: %v", err)
			}

			address := crypto.GetAddress(crypto.GetXpub(privateKey))

			payload.Set("type", "Transfer")
			payload.Set("symbol", symbol)
			payload.Set("to", to)
			payload.Set("amount", amount)
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

	cmd.Flags().StringVarP(&peer, "peer", "p", "localhost:9001", "peer to connect")
	cmd.Flags().StringVarP(&symbol, "symbol", "s", "", "token symbol")
	cmd.MarkFlagRequired("symbol")
	cmd.Flags().StringVarP(&to, "to", "t", "", "recipient address")
	cmd.MarkFlagRequired("to")
	cmd.Flags().StringVarP(&amount, "amount", "a", "", "amount to transfer")
	cmd.MarkFlagRequired("amount")

	return cmd
}

func CreateBalanceOfCmd() *cobra.Command {
	var peer string
	var symbol string

	cmd := &cobra.Command{
		Use:   "balance",
		Short: "check token balance of an address",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			payload := structure.NewOrderedMap()
			privateKey, err := GetPrivateKey()

			if err != nil {
				log.Fatalf("Failed to get private key: %v", err)
			}

			address := crypto.GetAddress(crypto.GetXpub(privateKey))

			token_address := util.HashMany("qrc_20", address, symbol)
			fmt.Println("token_address: ", token_address)

			payload.Set("type", "BalanceOf")
			payload.Set("token_address", token_address)
			payload.Set("address", address)
			payload.Set("timestamp", util.Utime())
			payload.Set("from", address)

			req := CreateWalletRequest(peer, payload.Ser())

			response, err := network.CallRawRequest(req)
			if err != nil {
				log.Fatalf("Failed to send request: %v", err)
			}

			rstr := FormatResponse(&response.Payload)
			fmt.Println("balance: ", rstr)
		},
	}

	cmd.Flags().StringVarP(&peer, "peer", "p", "localhost:9001", "peer to connect")
	cmd.Flags().StringVarP(&symbol, "symbol", "s", "", "token symbol")
	cmd.MarkFlagRequired("symbol")

	return cmd
}

func CreateTokenInfoCmd() *cobra.Command {
	var peer string
	var token_address string

	cmd := &cobra.Command{
		Use:   "info",
		Short: "get token information",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			payload := structure.NewOrderedMap()

			payload.Set("type", "GetTokenInfo")
			payload.Set("token_address", token_address)
			payload.Set("timestamp", util.Utime())

			req := CreateWalletRequest(peer, payload.Ser())

			response, err := network.CallRawRequest(req)
			if err != nil {
				log.Fatalf("Failed to send request: %v", err)
			}

			rstr := FormatResponse(&response.Payload)
			fmt.Println(rstr)
		},
	}

	cmd.Flags().StringVarP(&peer, "peer", "p", "localhost:9001", "peer to connect")
	cmd.Flags().StringVarP(&token_address, "token_address", "t", "", "token address")
	cmd.MarkFlagRequired("token_address")

	return cmd
}

func CreateDexCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dex",
		Short: "dex cli tool",
	}

	cmd.AddCommand(
		CreateMintTokenCmd(),
		CreateTransferTokenCmd(),
		CreateBalanceOfCmd(),
		CreateTokenInfoCmd(),
	)

	return cmd
}
