package program

import (
	"github.com/spf13/cobra"
)

func CreateGetBalanceCmd() *cobra.Command {
	var requestType string
	var peer string

	cmd := &cobra.Command{
		Use:   "balance",
		Short: "get wallet balance",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	cmd.Flags().StringVarP(&requestType, "type", "t", "", "type of wallet interface")
	cmd.Flags().StringVarP(&peer, "peer", "p", "", "peer to get balance")
	cmd.MarkFlagRequired("peer")

	return cmd
}

func CreateWalletCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wallet",
		Short: "api cli tool",
	}

	cmd.AddCommand(
		CreateGetBalanceCmd(),
	)

	return cmd
}
