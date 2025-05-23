package program

import (
	"bytes"
	"fmt"
	"hello/pkg/core/config"
	C "hello/pkg/core/config"
	"hello/pkg/core/model"
	"hello/pkg/core/network"
	"hello/pkg/crypto"
	"hello/pkg/rpc"
	"hello/pkg/util"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"hello/pkg/core/structure"

	"encoding/json"

	"strings"

	"github.com/spf13/cobra"
)

func FormatResponse(payload *json.RawMessage) string {
	var prettyJSON bytes.Buffer

	if err := json.Indent(&prettyJSON, *payload, "", "    "); err != nil {
		log.Fatalf("JSON formatting failed: %v", err)
		return ""
	}
	return prettyJSON.String()
}

func CreateWalletTransaction(peer string, payload string) *rpc.TransactionRequest {
	privateKey, err := GetPrivateKey()

	if err != nil {
		log.Fatalf("Failed to get private key: %v", err)
	}

	data, err := structure.ParseOrderedMap(payload)

	if err != nil {
		log.Fatalf("Failed to parse payload: %v", err)
	}

	signedTx, err := model.FromRawData(data, privateKey, crypto.GetXpub(privateKey))

	if err != nil {
		log.Fatalf("Failed to create signed transaction: %v", err)
	}

	err = signedTx.Validate()
	if err != nil {
		log.Fatalf("Failed to validate signed transaction: %v", err)
	}

	payload, err = signedTx.Ser()

	if err != nil {
		log.Fatalf("Failed to serialize signed transaction: %v", err)
	}

	return rpc.CreateTransactionRequest(payload, peer)
}

func CreateWalletRequest(peer string, payload string) *rpc.RawRequest {
	privateKey, err := GetPrivateKey()

	if err != nil {
		log.Fatalf("Failed to get private key: %v", err)
	}

	data, err := structure.ParseOrderedMap(payload)

	if err != nil {
		log.Fatalf("Failed to parse payload: %v", err)
	}

	signedRequest := model.NewSignedRequestFromRawData(data, privateKey)
	err = signedRequest.Validate()
	if err != nil {
		log.Fatalf("Failed to validate signed request: %v", err)
	}
	payload, err = signedRequest.Ser()

	if err != nil {
		log.Fatalf("Failed to serialize signed request: %v", err)
	}

	return rpc.CreateRequest(payload, peer)
}

func GetPrivateKey() (string, error) {
	walletPath := filepath.Join(config.DATA_ROOT_DIR, ".wallet")

	// Check if wallet file exists
	if _, err := os.Stat(walletPath); os.IsNotExist(err) {
		return "", fmt.Errorf("wallet file not found: %v", err)
	}

	// Read private key from wallet file
	privateKey, err := os.ReadFile(walletPath)
	if err != nil {
		return "", fmt.Errorf("failed to read wallet file: %v", err)
	}

	return string(privateKey), nil
}

func CreateSetWalletCmd() *cobra.Command {
	var privateKey string

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set default wallet private key",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			walletPath := filepath.Join(config.DATA_ROOT_DIR, ".wallet")

			// Create data directory if it doesn't exist
			if err := os.MkdirAll(config.DATA_ROOT_DIR, 0755); err != nil {
				log.Fatalf("Failed to create data directory: %v", err)
			}

			// Write private key to .wallet file
			if err := os.WriteFile(walletPath, []byte(privateKey), 0600); err != nil {
				log.Fatalf("Failed to save wallet file: %v", err)
			}

			fmt.Println("Default wallet has been set successfully")
		},
	}

	cmd.Flags().StringVarP(&privateKey, "privatekey", "k", "", "private key for default wallet")
	cmd.MarkFlagRequired("privatekey")

	return cmd
}

func CreateGetBalanceCmd() *cobra.Command {
	var peer string
	var address string

	privateKey, err := GetPrivateKey()
	_address := crypto.GetAddress(crypto.GetXpub(privateKey))

	cmd := &cobra.Command{
		Use:   "balance",
		Short: "get wallet balance",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			payload := structure.NewOrderedMap()
			fmt.Println("address: ", address)

			if err != nil {
				log.Fatalf("Failed to get private key: %v", err)
			}

			payload.Set("type", "GetBalance")
			payload.Set("address", address)
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

	cmd.Flags().StringVarP(&peer, "peer", "p", C.CLI_DEFAULT_REQUEST, "peer to get balance")
	cmd.Flags().StringVarP(&address, "address", "a", _address, "peer to get balance")

	return cmd
}

func CreateSendTransactionCmd() *cobra.Command {
	var peer string
	var toaddress string
	var amount string

	cmd := &cobra.Command{
		Use:   "send",
		Short: "send transaction",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			payload := structure.NewOrderedMap()
			privateKey, err := GetPrivateKey()

			if err != nil {
				log.Fatalf("Failed to get private key: %v", err)
			}

			address := crypto.GetAddress(crypto.GetXpub(privateKey))
			amount = util.MulByE18(amount)

			payload.Set("type", "Send")
			payload.Set("amount", amount)
			payload.Set("to", toaddress)
			payload.Set("from", address)
			payload.Set("timestamp", util.Utime())

			req := CreateWalletTransaction(peer, payload.Ser())

			response, err := network.CallTransactionRequest(req)
			if err != nil {
				log.Fatalf("Failed to send request: %v", err)
			}

			rstr := FormatResponse(&response.Payload)
			fmt.Println(rstr)
		},
	}

	cmd.Flags().StringVarP(&peer, "peer", "p", C.CLI_DEFAULT_REQUEST, "peer to get balance")
	cmd.Flags().StringVarP(&toaddress, "toaddress", "t", "", "to address")
	cmd.MarkFlagRequired("toaddress")
	cmd.Flags().StringVarP(&amount, "amount", "a", "", "amount")
	cmd.MarkFlagRequired("amount")

	return cmd
}

func CreateFaucetTransactionCmd() *cobra.Command {
	var peer string

	cmd := &cobra.Command{
		Use:   "faucet",
		Short: "send faucet transaction",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			payload := structure.NewOrderedMap()
			privateKey, err := GetPrivateKey()

			if err != nil {
				log.Fatalf("Failed to get private key: %v", err)
			}

			address := crypto.GetAddress(crypto.GetXpub(privateKey))

			payload.Set("type", "Faucet")
			payload.Set("address", address)
			payload.Set("timestamp", util.Utime())
			payload.Set("from", address)
			spayload := payload.Ser()

			req := CreateWalletTransaction(peer, spayload)
			fmt.Println(req.Payload)

			response, err := network.CallTransactionRequest(req)
			if err != nil {
				log.Fatalf("Failed to send request: %v", err)
			}

			rstr := FormatResponse(&response.Payload)
			fmt.Println(rstr)
		},
	}

	cmd.Flags().StringVarP(&peer, "peer", "p", C.CLI_DEFAULT_REQUEST, "peer to get balance")

	return cmd
}

func CreateWallet() *cobra.Command {
	var peer string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create new wallet",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			privateKey, publicKey := crypto.GenerateKeyPair()
			address := crypto.GetAddress(publicKey)

			fmt.Println("prviate key: ", privateKey)
			fmt.Println("public key: ", publicKey)
			fmt.Println("address: ", address)
		},
	}

	cmd.Flags().StringVarP(&peer, "peer", "p", C.CLI_DEFAULT_REQUEST, "peer to get balance")
	return cmd
}

func CreateCountTransactionCmd() *cobra.Command {
	var peer string

	cmd := &cobra.Command{
		Use:   "count",
		Short: "count",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			payload := structure.NewOrderedMap()
			privateKey, err := GetPrivateKey()

			if err != nil {
				log.Fatalf("Failed to get private key: %v", err)
			}

			address := crypto.GetAddress(crypto.GetXpub(privateKey))

			payload.Set("type", "Count")
			payload.Set("from", address)
			payload.Set("timestamp", util.Utime())

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

	cmd.Flags().StringVarP(&peer, "peer", "p", C.CLI_DEFAULT_REQUEST, "peer to get balance")
	return cmd
}

func CreateMultiSendCmd() *cobra.Command {
	var peer string
	var toaddresses []string
	var amount string

	cmd := &cobra.Command{
		Use:   "multisend",
		Short: "send transaction to multiple addresses",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {

			if len(toaddresses) > 5 {
				log.Fatalf("Number of addresses cannot exceed 5. Got: %d", len(toaddresses))
			}

			// Check for duplicate addresses
			addressMap := make(map[string]bool)
			for _, addr := range toaddresses {
				if addressMap[addr] {
					log.Fatalf("Duplicate address found: %s", addr)
				}
				addressMap[addr] = true
			}

			payload := structure.NewOrderedMap()
			privateKey, err := GetPrivateKey()

			if err != nil {
				log.Fatalf("Failed to get private key: %v", err)
			}

			address := crypto.GetAddress(crypto.GetXpub(privateKey))
			amount = util.MulByE18(amount)

			// Concatenate all addresses without separator
			combinedAddresses := strings.Join(toaddresses, "")

			payload.Set("type", "MultiSend")
			payload.Set("amount", amount)
			payload.Set("beneficiaries", combinedAddresses)
			payload.Set("txcount", strconv.Itoa(len(toaddresses)))
			payload.Set("from", address)
			payload.Set("timestamp", util.Utime())

			fmt.Println(payload.Ser())

			req := CreateWalletTransaction(peer, payload.Ser())

			response, err := network.CallTransactionRequest(req)
			if err != nil {
				log.Fatalf("Failed to send request: %v", err)
			}

			rstr := FormatResponse(&response.Payload)
			fmt.Println(rstr)
		},
	}

	cmd.Flags().StringVarP(&peer, "peer", "p", C.CLI_DEFAULT_REQUEST, "peer to send transaction")
	cmd.Flags().StringArrayVarP(&toaddresses, "toaddresses", "t", []string{}, "list of recipient addresses")
	cmd.MarkFlagRequired("toaddresses")
	cmd.Flags().StringVarP(&amount, "amount", "a", "", "amount to send to each address")
	cmd.MarkFlagRequired("amount")

	return cmd
}

func CreateWalletCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wallet",
		Short: "api cli tool",
	}

	cmd.AddCommand(
		CreateGetBalanceCmd(),
		CreateSetWalletCmd(),
		CreateFaucetTransactionCmd(),
		CreateSendTransactionCmd(),
		CreateCountTransactionCmd(),
		CreateWallet(),
		CreateMultiSendCmd(),
	)

	return cmd
}
