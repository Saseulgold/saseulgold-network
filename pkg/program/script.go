package program

import (
	"fmt"
	"hello/pkg/core/storage"
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

func createScriptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "script",
		Short: "script cli tool",
	}

	cmd.AddCommand(
		createGenesisCmd(),
	)

	return cmd
}
