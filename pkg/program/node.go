package program

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hello/pkg/core/network"
	"hello/pkg/swift"
	"log"

	"github.com/spf13/cobra"
)

func createPingCmd() *cobra.Command {
	var targetNode string

	cmd := &cobra.Command{
		Use:   "ping",
		Short: "ping to another node",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("ping to node: %s\n", targetNode)

			req := swift.Packet{
				Type:    swift.PacketTypePing,
				Payload: json.RawMessage(`{"message": "hello"}`),
			}

			response, err := network.CallRPC(targetNode, req)
			if err != nil {
				log.Fatalf("RPC call failed: %v", err)
			}
			fmt.Printf("RPC response: %v\n", response)

		},
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.Flags().Parse(args)
		},
	}

	cmd.Flags().StringVarP(&targetNode, "peer", "p", "", "Target node to ping")
	cmd.MarkFlagRequired("peer")

	return cmd
}

func createPeerCmd() *cobra.Command {
	var targetNode string

	cmd := &cobra.Command{
		Use:   "peer",
		Short: "get peer list",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("get peer list from node: %s\n", targetNode)

			req := swift.Packet{
				Type:    swift.PacketTypePeerRequest,
				Payload: json.RawMessage(`{"message": "hello"}`),
			}

			response, err := network.CallRPC(targetNode, req)
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, response.Payload, "", "    "); err != nil {
				log.Fatalf("JSON formatting failed: %v", err)
			}
			fmt.Printf("Peer list:\n%s\n", prettyJSON.String())
			if err != nil {
				log.Fatalf("RPC call failed: %v", err)
			}
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.Flags().Parse(args)
		},
	}

	cmd.Flags().StringVarP(&targetNode, "peer", "p", "", "Target node to get peer list")
	cmd.MarkFlagRequired("peer")

	return cmd
}

func createMempoolCmd() *cobra.Command {
	var targetNode string

	cmd := &cobra.Command{
		Use:   "listmemtx",
		Short: "get mempool transaction list",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("get mempool transaction list from node: %s\n", targetNode)

			req := swift.Packet{
				Type:    swift.PacketTypeListMempoolTransactionRequest,
				Payload: json.RawMessage(`{}`),
			}

			response, err := network.CallRPC(targetNode, req)
			if err != nil {
				log.Fatalf("RPC call failed: %v", err)
			}

			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, response.Payload, "", "    "); err != nil {
				log.Fatalf("JSON formatting failed: %v", err)
			}
			fmt.Printf("Mempool transaction list:\n%s\n", prettyJSON.String())
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.Flags().Parse(args)
		},
	}

	cmd.Flags().StringVarP(&targetNode, "peer", "p", "", "Target node to get mempool transactions")
	cmd.MarkFlagRequired("peer")

	return cmd
}

func createBroadcastTxCmd() *cobra.Command {
	var targetNode string
	var message string

	cmd := &cobra.Command{
		Use:   "broadcast",
		Short: "broadcast transaction to network",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("broadcasting transaction to node: %s\n", targetNode)

			req := swift.Packet{
				Type:    swift.PacketTypeBroadcastTransactionRequest,
				Payload: json.RawMessage(message),
			}

			response, err := network.CallRPC(targetNode, req)
			if err != nil {
				log.Fatalf("RPC call failed: %v", err)
			}
			fmt.Printf("Broadcast response: %s\n", string(response.Payload))
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.Flags().Parse(args)
		},
	}

	cmd.Flags().StringVarP(&targetNode, "peer", "p", "", "Target node to broadcast transaction")
	cmd.Flags().StringVarP(&message, "message", "m", "", "Transaction message to broadcast")
	cmd.MarkFlagRequired("peer")
	cmd.MarkFlagRequired("message")

	return cmd
}

func createStatusBundleCmd() *cobra.Command {
	var targetNode string

	cmd := &cobra.Command{
		Use:   "status",
		Short: "get node status bundle",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("get status bundle from node: %s\n", targetNode)

			req := swift.Packet{
				Type:    swift.PacketTypeGetStatusBundleRequest,
				Payload: json.RawMessage(`{}`),
			}

			response, err := network.CallRPC(targetNode, req)
			if err != nil {
				log.Fatalf("RPC call failed: %v", err)
			}

			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, response.Payload, "", "    "); err != nil {
				log.Fatalf("JSON formatting failed: %v", err)
			}
			fmt.Printf("Status bundle:\n%s\n", prettyJSON.String())
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.Flags().Parse(args)
		},
	}

	cmd.Flags().StringVarP(&targetNode, "peer", "p", "", "Target node to get status bundle")
	cmd.MarkFlagRequired("peer")

	return cmd
}

func createNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node",
		Short: "node cli tool",
	}

	cmd.AddCommand(
		createPingCmd(),
		createPeerCmd(),
		createMempoolCmd(),
		createBroadcastTxCmd(),
		createStatusBundleCmd(),
	)

	return cmd
}
