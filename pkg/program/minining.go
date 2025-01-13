package program

import (
	"fmt"
	C "hello/pkg/core/config"
	"hello/pkg/core/network"
	"hello/pkg/core/storage"
	"hello/pkg/core/structure"
	"hello/pkg/crypto"
	"hello/pkg/util"
	"log"

	"github.com/spf13/cobra"

	"bytes"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"go.uber.org/zap"
	"net/http"
	"io"
	"encoding/json"
	"io/ioutil"
)

func getNodeInfo() string {
	url := "http://ipinfo.io/json"

	// Perform the GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching URL: %v\n", err)
		os.Exit(1)
	}
	// Ensure the response body is closed after function finishes
	defer resp.Body.Close()

	// Read the response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading response body: %v\n", err)
		os.Exit(1)
	}

	// Print the response data as a string
	return string(data)
}

func IncNodeCount(address string, info interface{}) (*http.Response, error) {
    url := fmt.Sprintf("https://api.saseulgold.org/mining/v2/%s", address)

    jsonData, err := json.Marshal(info)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal info: %v", err)
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %v", err)
    }

    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("request failed: %v", err)
    }

    if resp.StatusCode != http.StatusOK {
        bodyBytes, readErr := ioutil.ReadAll(resp.Body)
        if readErr != nil {
            return nil, fmt.Errorf("inc failed with status %d and failed to read response body", resp.StatusCode)
        }
        fmt.Println(string(bodyBytes))
        resp.Body.Close()
        return nil, fmt.Errorf("inc failed with status code: %d", resp.StatusCode)
    }

    return resp, nil
}


func beforeMiningStart(address string) error {
	info := getNodeInfo()
	fmt.Println(info)
	_, err := IncNodeCount(address, info)
	return err
}

func miningProcess() error {
	file1, _ := os.OpenFile("mining.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	miningLogger, _ := util.CreateLogger(file1)

	// Step 1: Get last height and blockhash
	lastHeightCmd := exec.Command("./sg", "api", "lastheight")
	var lastHeightOut bytes.Buffer
	lastHeightCmd.Stdout = &lastHeightOut
	if err := lastHeightCmd.Run(); err != nil {
		return err
	}

	lastHeight := strings.TrimSpace(lastHeightOut.String()) // 결과값에서 공백 제거
	/**
	lastHeightInt, _ := strconv.Atoi(lastHeight)
	nodeLastHeight := storage.LastHeight()
	
	if nodeLastHeight > lastHeightInt {
		return fmt.Errorf("node last height is greater than network last height")
	}

	if lastHeightInt-10 > nodeLastHeight {
		return fmt.Errorf("block is not synced")
	}
	**/

	miningLogger.Info("last height:", zap.String("height", lastHeight))

	ts := strconv.Itoa(int(util.Utime()))
	seed := "blk-" + lastHeight + "," + ts

	mineCmd := exec.Command("./mine/cmine", lastHeight, ts)
	var mineOut bytes.Buffer

	mineCmd.Stdout = &mineOut
	if err := mineCmd.Run(); err != nil {
		fmt.Printf("Error running mine: %s\n", err)
		return err
	}
	mineResult := strings.TrimSpace(mineOut.String()) // 결과값에서 공백 제거

	lines := strings.Split(mineResult, "\n")
	if len(lines) < 2 {
		fmt.Println("Error: Invalid mine result format")
		return fmt.Errorf("invalid mine result format")
	}
	nonce := lines[0]
	hash := lines[1]

	miningLogger.Info("calculate hash:", zap.String("hash", hash), zap.String("nonce", nonce))

	submitCmd := exec.Command("./sg", "mining", "submit",
		"-c", hash, "-e", seed, "-n", nonce)
	var submitOut bytes.Buffer
	submitCmd.Stdout = &submitOut
	if err := submitCmd.Run(); err != nil {
		return fmt.Errorf("error submitting result: %s", err)
	}
	return nil
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

			isRunning := util.ServiceIsRunning(storage.DataRootDir(), "mining")

			if isRunning {
				fmt.Println("mining is already running")
			}

			if child == "c" {
				err := util.ProcessStart(storage.DataRootDir(), "mining", os.Getpid())

				for true {
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					err = miningProcess()
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
				}
			} else {

				privateKey, _ := GetPrivateKey()
				address := crypto.GetAddress(crypto.GetXpub(privateKey))

				err := beforeMiningStart(address)

				if err != nil {
					fmt.Println(err)
				}
				
				cmd := exec.Command("./sg", "mining", "start", "-c", "c")

				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Stdin = os.Stdin

				cmd.SysProcAttr = &syscall.SysProcAttr{
					Setsid: true,
				}

				err = cmd.Start()
				if err != nil {
					fmt.Printf("Failed to start command: %s\n", err)
					os.Exit(1)
				}

				fmt.Printf("Started mining process with PID %d\n", cmd.Process.Pid)
			}
		},
	}

	cmd.Flags().StringVarP(&peer, "peer", "p", C.CLI_DEFAULT_REQUEST, "peer to get balance")
	cmd.Flags().StringVarP(&child, "child", "c", "", "")
	return cmd
}

func CreateStopMiningCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "stop mining",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			util.TerminateProcess(storage.DataRootDir(), "mining")
		},
	}

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
		CreateStopMiningCMD(),
	)

	return cmd
}
