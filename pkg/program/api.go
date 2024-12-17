package program

import (
	"fmt"
	"hello/pkg/core/network"
	"hello/pkg/rpc"
	"log"

	"github.com/spf13/cobra"
)

func CreateRawRequestCmd() *cobra.Command {
	var requestType string
	var peer string

	payload := rpc.CreateRawRequest(requestType, peer)
	response, err := network.CallRawRequest(payload)
	if err != nil {
		log.Fatalf("Failed to call raw request: %v", err)
	}

	fmt.Println(response)

	cmd := &cobra.Command{
		Use:   "rawrequest",
		Short: "get raw request",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	cmd.Flags().StringVarP(&requestType, "type", "t", "", "type of raw request")
	cmd.Flags().StringVarP(&peer, "peer", "p", "", "peer to get balance")
	cmd.MarkFlagRequired("peer")

	return cmd
}

func CreateApiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "api cli tool",
	}

	cmd.AddCommand(
		CreateRawRequestCmd(),
	)

	return cmd
}
