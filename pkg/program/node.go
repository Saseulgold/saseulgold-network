package program

import (
	"fmt"
	"hello/pkg/core/network"
	"hello/pkg/swift"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

func createNodeStartCmd(port *string, useTLS *bool) *cobra.Command {
	var foreground bool
	var rpcPort string

	cmd := &cobra.Command{
		Use:   "start",
		Short: "node start",
		Run: func(cmd *cobra.Command, args []string) {

			security := swift.SecurityConfig{
				UseTLS: *useTLS,
			}

			server := swift.NewServer("localhost:"+*port, security)
			node := network.NewNodeService(server, "localhost:"+rpcPort)

			if !foreground {
				// Fork process
				if pid := os.Getpid(); pid != 1 {
					// Parent process
					if err := node.Start(); err != nil {
						log.Fatalf("Failed to start node: %v", err)
					}
					fmt.Printf("Node started in background mode. PID: %d\n", pid)
					os.Exit(0)
				}
			}

			// Child process or foreground mode
			if err := node.Start(); err != nil {
				log.Fatalf("Failed to start node: %v", err)
			}

			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

			<-sigChan

			fmt.Println("Shutting down node...")
		},
	}

	cmd.Flags().StringVarP(port, "port", "p", "9090", "Port for Swift server")
	cmd.Flags().StringVarP(&rpcPort, "rpc-port", "r", "9091", "Port for RPC server")
	cmd.Flags().BoolVarP(useTLS, "tls", "t", false, "Use TLS for security")
	cmd.Flags().BoolVarP(&foreground, "foreground", "f", false, "Run node in foreground mode")

	return cmd
}

func createNodeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "node",
		Short: "node cli tool",
	}
}
