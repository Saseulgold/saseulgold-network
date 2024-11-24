package program

import (
	"fmt"
	"hello/pkg/swift"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

func createRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "program",
		Short: "program cli tool",
	}
}

func createServerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "server cli tool",
	}
}

func createServerStartCmd(port *string, useTLS *bool) *cobra.Command {
	var foreground bool

	cmd := &cobra.Command{
		Use:   "start",
		Short: "server start",
		Run: func(cmd *cobra.Command, args []string) {
			security := swift.SecurityConfig{
				UseTLS: *useTLS,
			}

			server := swift.NewServer(":"+*port, security)

			if !foreground {
				// Fork process
				if pid := os.Getpid(); pid != 1 {
					// Parent process
					if err := server.Start(); err != nil {
						log.Fatalf("Failed to start server: %v", err)
					}
					fmt.Printf("Server started in background mode. PID: %d\n", pid)
					os.Exit(0)
				}
			}

			// Child process or foreground mode
			if err := server.Start(); err != nil {
				log.Fatalf("Failed to start server: %v", err)
			}

			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

			<-sigChan

			fmt.Println("Shutting down server...")
		},
	}

	// 플래그 설정
	cmd.Flags().StringVarP(port, "port", "p", "8080", "Port to run the server on")
	cmd.Flags().BoolVarP(useTLS, "tls", "t", false, "Use TLS for security")
	cmd.Flags().BoolVarP(&foreground, "foreground", "f", false, "Run server in foreground mode")

	return cmd
}

func RunServerCMD() *cobra.Command {
	var port string
	var useTLS bool

	rootCmd := createRootCmd()
	serverCmd := createServerCmd()
	serverStartCmd := createServerStartCmd(&port, &useTLS)

	serverCmd.AddCommand(serverStartCmd)
	rootCmd.AddCommand(serverCmd)
	nodeCmd := createNodeCmd()
	nodeStartCmd := createNodeStartCmd(&port, &useTLS)
	nodeCmd.AddCommand(nodeStartCmd)
	rootCmd.AddCommand(nodeCmd)

	return rootCmd
}
