package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "main",
	Short: "bakeshop",
	Long:  "bakeshop program",
}

func main() {
	fmt.Println("Starting bakeshop")

	// Config loading
	viper.SetEnvPrefix("bakeshop") // Set the environment prefix to BAKESHOP_*
	viper.AutomaticEnv()           // Automatically search for environment variables

	// LND configs
	// lnd admin macaroon
	rootCmd.PersistentFlags().String("lnd.macaroon", "", "macaroon for lnd")
	viper.BindPFlag("lnd.macaroon", rootCmd.PersistentFlags().Lookup("lnd.macaroon"))
	viper.SetDefault("lnd.macaroon", "/Users/a/.polar/networks/1/volumes/lnd/dave/data/chain/bitcoin/regtest/admin.macaroon")

	rootCmd.PersistentFlags().String("lnd.macaroonHex", "", "macaroon hex for lnd")
	viper.BindPFlag("lnd.macaroonHex", rootCmd.PersistentFlags().Lookup("lnd.macaroonHex"))
	viper.SetDefault("lnd.macaroonHex", "")

	// lnd tls
	rootCmd.PersistentFlags().String("lnd.tls", "", "tls for lnd")
	viper.BindPFlag("lnd.tls", rootCmd.PersistentFlags().Lookup("lnd.tls"))
	viper.SetDefault("lnd.tls", "/Users/a/.polar/networks/1/volumes/lnd/dave/tls.cert")

	// lnd ip
	rootCmd.PersistentFlags().String("lnd.url", "", "url for lnd")
	viper.BindPFlag("lnd.url", rootCmd.PersistentFlags().Lookup("lnd.url"))
	viper.SetDefault("lnd.url", "127.0.0.1:10004")

	// subscriber destination pubkey
	rootCmd.PersistentFlags().String("subscriber.pubkey", "", "pubkey for the subscriber to send funds to")
	viper.BindPFlag("subscriber.pubkey", rootCmd.PersistentFlags().Lookup("subscriber.pubkey"))
	viper.SetDefault("subscriber.pubkey", "")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
