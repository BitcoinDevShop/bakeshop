package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "serve the program",
		Run: func(cmd *cobra.Command, args []string) {
			if err := serveAction(); err != nil {
				panic(err)
			}
		},
	}

	rootCmd.AddCommand(serveCmd)
}

func serveAction() error {
	// config parameters
	adminMacaroonPath := viper.GetString("lnd.macaroon")
	tlsCertPath := viper.GetString("lnd.tls")
	ipAddr := viper.GetString("lnd.url")

	client, tls, err := CreateLNDClient(adminMacaroonPath, "", tlsCertPath, ipAddr)
	if err != nil {
		panic(err)
	}

	// Create website server
	fmt.Println("Starting server")
	httpServer, err := createServer(&ServerConfig{
		LndClient: client,
		LndURL:    ipAddr,
		LndTLS:    tls,
	})
	if err != nil {
		panic(err)
	}

	// Daemon wait until cancel
	c := make(chan os.Signal)
	done := make(chan bool, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		httpServer.Shutdown(context.Background())
		fmt.Println("Shut down server")
		done <- true
	}()
	fmt.Println("Started...")
	<-done
	fmt.Println("Stopping...")

	return nil
}
