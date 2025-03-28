package program

import (
	"encoding/json"
	"fmt"
	. "hello/pkg/core/debug"
	"hello/pkg/core/network"
	"hello/pkg/service"
	"hello/pkg/swift"
	"log"
	"os"
	"os/signal"
	"syscall"

	C "hello/pkg/core/config"

	"github.com/spf13/cobra"
)

func createRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "program",
		Short: "program cli tool",
	}
}

func createNetworkCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "network",
		Short: "network cli tool",
	}
}

func createNetworkStartCmd(useTLS *bool) *cobra.Command {
	var foreground bool

	cmd := &cobra.Command{
		Use:   "start",
		Short: "network start",
		Run: func(cmd *cobra.Command, args []string) {
			security := swift.SecurityConfig{
				UseTLS: *useTLS,
			}
			port := cmd.Flag("port").Value.String()

			oracle := service.GetOracleService()
			err := oracle.OnStartUp(security, port)
			if err != nil {
				log.Fatalf("Failed to start oracle: %v", err)
			}

			if !foreground {
				// Fork process
				if pid := os.Getpid(); pid != 1 {
					// Parent process
					if err := oracle.Start(); err != nil {
						log.Fatalf("Failed to start server: %v", err)
					}
					fmt.Printf("Server started in background mode. PID: %d\n", pid)
					os.Exit(0)
				}
			}

			// Child process or foreground mode
			DebugLog("server starting")
			if err := oracle.Start(); err != nil {
				log.Fatalf("Failed to start server: %v", err)
			}

			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

			<-sigChan

			fmt.Println("Shutting down server...")
		},
	}

	var port string
	cmd.Flags().BoolVarP(useTLS, "tls", "t", false, "Use TLS for security")
	cmd.Flags().BoolVarP(&foreground, "foreground", "f", false, "Run server in foreground mode")
	cmd.Flags().StringVarP(&port, "port", "p", "8080", "server port")

	cmd.Flags().BoolVarP(&C.CORE_TEST_MODE, "debug", "d", false, "Enable test mode")
	cmd.Flags().StringVarP(&C.DATA_TEST_ROOT_DIR, "rootdir", "r", "", "root dir")

	return cmd
}

func createConnectReplicaCmd() *cobra.Command {
	var targetAddr string
	var peerAddr string

	cmd := &cobra.Command{
		Use:   "replica",
		Short: "network replica",
		Run: func(cmd *cobra.Command, args []string) {

			if targetAddr == peerAddr {
				fmt.Println("target address and peer address cannot be the same")
				return
			}

			response, err := network.CallRPC(peerAddr, swift.Packet{
				Type:    swift.PacketTypeRegisterReplicaRequest,
				Payload: json.RawMessage(`{"targetAddr": "` + targetAddr + `"}`),
			})

			if err != nil {
				log.Fatalf("Failed to send request: %v", err)
			}

			rstr := FormatResponse(&response.Payload)
			fmt.Println(rstr)
		},
	}

	cmd.Flags().StringVarP(&targetAddr, "target", "t", "", "target address")
	cmd.MarkFlagRequired("target")
	cmd.Flags().StringVarP(&peerAddr, "peer", "p", "localhost:9001", "peer address")
	cmd.Flags().BoolVarP(&C.CORE_TEST_MODE, "debug", "d", false, "Enable test mode")
	cmd.Flags().StringVarP(&C.DATA_TEST_ROOT_DIR, "rootdir", "r", "", "root dir")

	return cmd
}

func createNetworkStopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "network stop",
		Run: func(cmd *cobra.Command, args []string) {
			oracle := service.GetOracleService()
			oracle.Shutdown()
		},
	}
	cmd.Flags().BoolVarP(&C.CORE_TEST_MODE, "debug", "d", false, "Enable test mode")
	cmd.Flags().StringVarP(&C.DATA_TEST_ROOT_DIR, "rootdir", "r", "", "root dir")

	return cmd
}

func createNetworkMetricsCmd() *cobra.Command {
	var peer string

	cmd := &cobra.Command{
		Use:   "metrics",
		Short: "network metrics",
		Run: func(cmd *cobra.Command, args []string) {
			packet := swift.Packet{
				Type:    swift.PacketTypeMetricsRequest,
				Payload: json.RawMessage(`{}`),
			}

			response, err := network.CallRPC(peer, packet)
			if err != nil {
				log.Fatalf("Failed to send request: %v", err)
			}

			rstr := FormatResponse(&response.Payload)
			fmt.Println(rstr)
		},
	}

	cmd.Flags().StringVarP(&peer, "peer", "p", "", "peer address")
	cmd.MarkFlagRequired("peer")

	return cmd
}

func RunNetworkCMD() *cobra.Command {
	var useTLS bool

	rootCmd := createRootCmd()
	networkCmd := createNetworkCmd()
	networkStartCmd := createNetworkStartCmd(&useTLS)
	networkStopCmd := createNetworkStopCmd()
	networkConnectReplicaCmd := createConnectReplicaCmd()
	networkMetricsCmd := createNetworkMetricsCmd()

	networkCmd.AddCommand(networkStartCmd)
	networkCmd.AddCommand(networkStopCmd)
	networkCmd.AddCommand(networkConnectReplicaCmd)
	networkCmd.AddCommand(networkMetricsCmd)

	rootCmd.AddCommand(networkCmd)

	walletCmd := CreateWalletCmd()
	dexCmd := CreateDexCmd()
	apiCmd := CreateApiCmd()
	miningCmd := CreateMiningCmd()

	//TODO remove dev commands in production
	nodeCmd := createNodeCmd()
	scriptCmd := createScriptCmd()

	rootCmd.AddCommand(nodeCmd)
	rootCmd.AddCommand(scriptCmd)
	rootCmd.AddCommand(walletCmd)
	rootCmd.AddCommand(apiCmd)
	rootCmd.AddCommand(dexCmd)
	rootCmd.AddCommand(miningCmd)

	return rootCmd
}
