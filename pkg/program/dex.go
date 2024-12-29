package program

import (
	"fmt"
	"hello/pkg/core/network"
	"hello/pkg/crypto"
	"hello/pkg/util"
	"log"

	"hello/pkg/core/structure"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var logger = util.GetLogger()

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

			token_address := util.HashMany("qrc_20", address, symbol)
			fmt.Println("token_address: ", token_address)

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
	var token_address string
	var to string
	var amount string

	cmd := &cobra.Command{
		Use:   "transfer",
		Short: "transfer token to another address",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("Transfer token to another address", zap.String("token_address", token_address), zap.String("to", to), zap.String("amount", amount))

			payload := structure.NewOrderedMap()
			privateKey, err := GetPrivateKey()

			if err != nil {
				log.Fatalf("Failed to get private key: %v", err)
			}

			address := crypto.GetAddress(crypto.GetXpub(privateKey))

			payload.Set("type", "Transfer")
			payload.Set("token_address", token_address)
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
	cmd.Flags().StringVarP(&token_address, "token_address", "t", "", "token address")
	cmd.MarkFlagRequired("token_address")
	cmd.Flags().StringVarP(&to, "to", "o", "", "recipient address")
	cmd.MarkFlagRequired("to")
	cmd.Flags().StringVarP(&amount, "amount", "a", "", "amount to transfer")
	cmd.MarkFlagRequired("amount")

	return cmd
}

func CreateBalanceOfCmd() *cobra.Command {
	var peer string
	var token_address string

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
	cmd.Flags().StringVarP(&token_address, "token_address", "t", "", "token address")
	cmd.MarkFlagRequired("token_address")

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

func CreateProvideLiquidityCmd() *cobra.Command {
	var peer string
	var token_address_a string
	var token_address_b string
	var amount_a string
	var amount_b string

	cmd := &cobra.Command{
		Use:   "provide",
		Short: "provide liquidity to token pool",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("Provide liquidity to pool",
				zap.String("token_address_a", token_address_a),
				zap.String("token_address_b", token_address_b),
				zap.String("amount_a", amount_a),
				zap.String("amount_b", amount_b))

			payload := structure.NewOrderedMap()
			privateKey, err := GetPrivateKey()

			if err != nil {
				log.Fatalf("Failed to get private key: %v", err)
			}

			address := crypto.GetAddress(crypto.GetXpub(privateKey))

			payload.Set("type", "LiquidityProvide")
			payload.Set("token_address_a", token_address_a)
			payload.Set("token_address_b", token_address_b)
			payload.Set("amount_a", amount_a)
			payload.Set("amount_b", amount_b)
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
	cmd.Flags().StringVarP(&token_address_a, "token_a", "a", "", "first token address")
	cmd.MarkFlagRequired("token_a")
	cmd.Flags().StringVarP(&token_address_b, "token_b", "b", "", "second token address")
	cmd.MarkFlagRequired("token_b")
	cmd.Flags().StringVarP(&amount_a, "amount_a", "x", "", "amount of first token")
	cmd.MarkFlagRequired("amount_a")
	cmd.Flags().StringVarP(&amount_b, "amount_b", "y", "", "amount of second token")
	cmd.MarkFlagRequired("amount_b")

	return cmd
}

func CreateWithdrawLiquidityCmd() *cobra.Command {
	var peer string
	var token_address_a string
	var token_address_b string
	var lp_amount string

	cmd := &cobra.Command{
		Use:   "withdraw",
		Short: "withdraw liquidity from token pool",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("Withdraw liquidity from pool",
				zap.String("token_address_a", token_address_a),
				zap.String("token_address_b", token_address_b),
				zap.String("lp_amount", lp_amount))

			payload := structure.NewOrderedMap()
			privateKey, err := GetPrivateKey()

			if err != nil {
				log.Fatalf("Failed to get private key: %v", err)
			}

			address := crypto.GetAddress(crypto.GetXpub(privateKey))

			payload.Set("type", "WithdrawLiquidity")
			payload.Set("token_address_a", token_address_a)
			payload.Set("token_address_b", token_address_b)
			payload.Set("lp_amount", lp_amount)
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
	cmd.Flags().StringVarP(&token_address_a, "token_a", "a", "", "first token address")
	cmd.MarkFlagRequired("token_a")
	cmd.Flags().StringVarP(&token_address_b, "token_b", "b", "", "second token address")
	cmd.MarkFlagRequired("token_b")
	cmd.Flags().StringVarP(&lp_amount, "lp_amount", "l", "", "LP token amount")
	cmd.MarkFlagRequired("lp_amount")

	return cmd
}

func CreateSwapCmd() *cobra.Command {
	var peer string
	var token_address_a string
	var token_address_b string
	var amount_a string

	cmd := &cobra.Command{
		Use:   "swap",
		Short: "swap token",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("Swap token",
				zap.String("token_address_a", token_address_a),
				zap.String("token_address_b", token_address_b),
				zap.String("amount_a", amount_a))

			payload := structure.NewOrderedMap()
			privateKey, err := GetPrivateKey()

			if err != nil {
				log.Fatalf("Failed to get private key: %v", err)
			}

			address := crypto.GetAddress(crypto.GetXpub(privateKey))

			payload.Set("type", "Swap")
			payload.Set("token_address_a", token_address_a)
			payload.Set("token_address_b", token_address_b)
			payload.Set("amount_a", amount_a)
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
	cmd.Flags().StringVarP(&token_address_a, "token_a", "a", "", "first token address")
	cmd.MarkFlagRequired("token_a")
	cmd.Flags().StringVarP(&token_address_b, "token_b", "b", "", "second token address")
	cmd.MarkFlagRequired("token_b")
	cmd.Flags().StringVarP(&amount_a, "amount_a", "x", "", "amount of first token")
	cmd.MarkFlagRequired("amount_a")

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
		CreateProvideLiquidityCmd(),
		CreateWithdrawLiquidityCmd(),
		CreateSwapCmd(),
	)

	return cmd
}
