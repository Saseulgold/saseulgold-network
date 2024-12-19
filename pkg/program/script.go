package program

import (
	"fmt"
	"hello/pkg/core/model"
	"hello/pkg/core/storage"
	"hello/pkg/core/structure"
	"hello/pkg/service"
	"hello/pkg/util"

	C "hello/pkg/core/config"

	"github.com/spf13/cobra"
)

func createGenesisCmd() *cobra.Command {
	var privateKey string
	var publicKey string

	cmd := &cobra.Command{
		Use:   "genesis",
		Short: "get genesis block information",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {

			isRunning := util.ServiceIsRunning(storage.DataRootDir(), "oracle")
			if isRunning {
				fmt.Println("oracle is already running. stop it first.")
				return
			}

			err := service.CommitGenesis(privateKey, publicKey)
			if err != nil {
				fmt.Println("failed to create genesis transaction: ", err)
				return
			}

		},
	}

	cmd.Flags().StringVarP(&privateKey, "privatekey", "k", "", "private key to sign genesis block")
	cmd.Flags().StringVarP(&publicKey, "publickey", "b", "", "public key to sign genesis block")
	cmd.Flags().BoolVarP(&C.CORE_TEST_MODE, "debug", "d", false, "Enable test mode")
	cmd.Flags().StringVarP(&C.DATA_TEST_ROOT_DIR, "rootdir", "r", "", "root dir")
	cmd.MarkFlagRequired("peer")

	return cmd
}

func createForceCommitCmd() *cobra.Command {
	var message string

	cmd := &cobra.Command{
		Use:   "forcecommit",
		Short: "force commit transaction to network",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {

			lastBlockHeight := storage.LastHeight()
			if lastBlockHeight == 0 {
				fmt.Println("no genesis block. create genesis block first.")
				return
			}

			isRunning := util.ServiceIsRunning(storage.DataRootDir(), "oracle")
			if isRunning {
				fmt.Println("oracle is already running. stop it first.")
				return
			}

			data, err := structure.ParseOrderedMap(message)

			if err != nil {
				fmt.Println("failed to parse transaction: ", err)
				return
			}

			fmt.Println("parsed transaction: ", data.Ser())

			tx, err := model.NewSignedTransaction(data)
			if err != nil {
				fmt.Println("failed to parse transaction: ", err)
				return
			}

			txs := make(map[string]*model.SignedTransaction)
			txs[tx.GetTxHash()] = &tx

			fmt.Println("txs: ", tx.GetTxData().Data.Ser())

			err = service.ForceCommit(txs)

			if err != nil {
				fmt.Println("failed to force commit: ", err)
				return
			}

		},
	}

	cmd.Flags().StringVarP(&message, "message", "m", "", "Transaction message to broadcast")
	cmd.Flags().BoolVarP(&C.CORE_TEST_MODE, "debug", "d", false, "Enable test mode")
	cmd.Flags().StringVarP(&C.DATA_TEST_ROOT_DIR, "rootdir", "r", "", "root dir")

	return cmd
}

func createGetStatusCmd() *cobra.Command {
	var key string

	cmd := &cobra.Command{
		Use:   "getstatus",
		Short: "check value for status key in network storage.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			err := service.GetStatus(key)

			if err != nil {
				fmt.Println("failed to read status: ", err)
				return
			}

		},
	}

	cmd.Flags().StringVarP(&key, "hashkey", "k", "", "status hash key")
	cmd.Flags().BoolVarP(&C.CORE_TEST_MODE, "debug", "d", false, "Enable test mode")
	cmd.Flags().StringVarP(&C.DATA_TEST_ROOT_DIR, "rootdir", "r", "", "root dir")

	return cmd
}


func createScriptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "script",
		Short: "script cli tool",
	}

	cmd.AddCommand(
		createGenesisCmd(),
		createForceCommitCmd(),
		createGetStatusCmd(),
	)

	return cmd
}
