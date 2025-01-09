package program


import (
	"fmt"
	"hello/pkg/core/network"
	C "hello/pkg/core/config"
	"hello/pkg/crypto"
	"hello/pkg/util"
	"log"
	"hello/pkg/core/structure"
	"github.com/spf13/cobra"

	"bytes"
	"os/exec"
	"strings"
    	"os"
    	"syscall"
 
	"go.uber.org/zap"
)


func miningProcess() {
   	file1, _ := os.OpenFile("mining.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	miningLogger, _ := util.CreateLogger(file1)

	// Step 1: Get last height and blockhash
	lastHeightCmd := exec.Command("go", "run", "main.go", "api", "lastheight")
	var lastHeightOut bytes.Buffer
	lastHeightCmd.Stdout = &lastHeightOut
	if err := lastHeightCmd.Run(); err != nil {
		return
	}
	lastHeight := strings.TrimSpace(lastHeightOut.String()) // 결과값에서 공백 제거
	miningLogger.Info("last height:", zap.String("height", lastHeight))

	// Step 2: execute mine binary
	mineCmd := exec.Command("./mine/cmine", lastHeight)
	var mineOut bytes.Buffer
	mineCmd.Stdout = &mineOut
	if err := mineCmd.Run(); err != nil {
		// fmt.Printf("Error running mine: %s\n", err)
		return
	}
	mineResult := strings.TrimSpace(mineOut.String()) // 결과값에서 공백 제거
	// fmt.Printf("%s\n", mineResult)

	lines := strings.Split(mineResult, "\n")
	if len(lines) < 2 {
		// fmt.Println("Error: Invalid mine result format")
		return
	}
	nonce := lines[0]
	hash := lines[1]

	miningLogger.Info("calculate hash:", zap.String("hash", hash), zap.String("nonce", nonce))

	submitCmd := exec.Command("go", "run", "main.go", "mining", "submit",
		"-c", hash, "-e", fmt.Sprint("blk-", lastHeight), "-n", nonce)
	var submitOut bytes.Buffer
	submitCmd.Stdout = &submitOut
	if err := submitCmd.Run(); err != nil {
		// fmt.Printf("Error submitting result: %s\n", err)
		return
	}
	// fmt.Printf("Submit result: %s\n", submitOut.String())
}


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

	cmd.Flags().StringVarP(&peer, "peer", "p", C.CLI_DEFAULT_REQUEST, "peer to get balance")
	cmd.Flags().StringVarP(&address, "address", "a", _address, "mining address")
	cmd.Flags().StringVarP(&epoch, "epoch", "e", "", "epoch")
	cmd.Flags().StringVarP(&nonce, "nonce", "n", "", "nonce")
	cmd.Flags().StringVarP(&chash, "chash", "c", "", "nonce")

	cmd.MarkFlagRequired("epoch")
	cmd.MarkFlagRequired("nonce")
	cmd.MarkFlagRequired("chash")

	return cmd
}

func CreateStartMiningCmd() *cobra.Command {
	var peer string
	var child string

	cmd := &cobra.Command{
		Use:   "start",
		Short: "start mining",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {

			if child == "c" {
				for true {
					miningProcess()
				}
			} else {
				    cmd := exec.Command("./sg", "mining", "start", "-c", "c")

				    cmd.Stdout = os.Stdout
				    cmd.Stderr = os.Stderr
				    cmd.Stdin = os.Stdin

				    cmd.SysProcAttr = &syscall.SysProcAttr{
					 Setsid: true,
				    }

				    err := cmd.Start()
				    if err != nil {
					fmt.Printf("Failed to start command: %s\n", err)
					return
				    }

				    fmt.Printf("Started mining process with PID %d\n", cmd.Process.Pid)

			}
		},
	}

	cmd.Flags().StringVarP(&peer, "peer", "p", C.CLI_DEFAULT_REQUEST, "peer to get balance")
	cmd.Flags().StringVarP(&child, "child", "c", "", "")
	return cmd
}


func CreateMiningCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mining",
		Short: "mining cli tool",
	}

	cmd.AddCommand(
		CreateSubmitMiningCmd(),
		CreateStartMiningCmd(),
	)

	return cmd
}
