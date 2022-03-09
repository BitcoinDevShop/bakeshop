package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	"github.com/lightningnetwork/lnd/lntypes"
	"github.com/lightningnetwork/lnd/record"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	subscriberCmd := &cobra.Command{
		Use:   "subscriber",
		Short: "Subscriber that pulls funds from the user",
		Run: func(cmd *cobra.Command, args []string) {
			if err := subscribeAction(); err != nil {
				panic(err)
			}
		},
	}

	rootCmd.AddCommand(subscriberCmd)
}

func subscribeAction() error {
	adminMacaroonHex := viper.GetString("lnd.macaroonHex")
	tlsCertPath := viper.GetString("lnd.tls")
	ipAddr := viper.GetString("lnd.url")
	subscriberPubkey := viper.GetString("subscriber.pubkey")

	fmt.Printf("Starting subscriber with parameters: %s - %s - %s - %s \n", adminMacaroonHex, tlsCertPath, ipAddr, subscriberPubkey)

	client, _, err := CreateLNDRouterClient("", adminMacaroonHex, tlsCertPath, ipAddr)
	if err != nil {
		panic(err)
	}

	// Parse macaroon details
	macBytes, err := hex.DecodeString(adminMacaroonHex)
	if err != nil {
		return err
	}
	_, amt, _, _, err := getMacaroonCaveats(macBytes)
	if err != nil {
		return err
	}

	// Send a keysend payment to the specified pubkey
	pubkey, err := hex.DecodeString(subscriberPubkey)
	if err != nil {
		return err
	}

	var preimage lntypes.Preimage
	if _, err := rand.Read(preimage[:]); err != nil {
		return err
	}
	preimageHash := preimage.Hash()

	destRecords := make(map[uint64][]byte)
	destRecords[record.KeySendType] = preimage[:]
	req := &routerrpc.SendPaymentRequest{
		Dest:              pubkey,
		Amt:               amt,
		DestCustomRecords: destRecords,
		FeeLimitSat:       9999,
		PaymentHash:       preimageHash[:],
		TimeoutSeconds:    100,
	}

	stream, err := client.SendPaymentV2(context.Background(), req)
	if err != nil {
		return err
	}

	for {
		payment, err := stream.Recv()
		if err != nil {
			fmt.Println(err.Error())
			break
		}

		switch payment.Status {
		case lnrpc.Payment_FAILED:
			fmt.Println(payment.FailureReason.String())
			break

		case lnrpc.Payment_SUCCEEDED:
			// TODO keep trying payment per subscription
			fmt.Println("payment succeeded")
			break
		}
	}

	return nil
}
