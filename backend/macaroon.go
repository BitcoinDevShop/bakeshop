package main

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/lightningnetwork/lnd/lnrpc"
	"gopkg.in/macaroon.v2"
)

func bakeMacaroon(client lnrpc.LightningClient, bakeReq BakeReq) (string, string, error) {
	ctx := context.Background()
	//secret := uint64(123)

	// real custom caveat
	resp, err := client.BakeMacaroon(ctx, &lnrpc.BakeMacaroonRequest{
		//RootKeyId:                secret,
		AllowExternalPermissions: true,
		Permissions: []*lnrpc.MacaroonPermission{
			{
				Entity: "offchain",
				Action: "read",
			},
			{
				Entity: "offchain",
				Action: "write",
			},
		},
	})
	if err != nil {
		return "", "", err
	}

	// Parse the mac to add a caveat to it
	freshMac := resp.GetMacaroon()
	freshMacBytes, err := hex.DecodeString(freshMac)
	if err != nil {
		return "", "", err
	}

	formattedMac := macaroon.Macaroon{}
	err = formattedMac.UnmarshalBinary(freshMacBytes)
	if err != nil {
		fmt.Println("Could not unmarshall mac")
		return "", "", err
	}

	formattedMac.AddFirstPartyCaveat([]byte(fmt.Sprintf("lnd-custom subscribe ms:%d amount:%d times:%d", bakeReq.Interval, bakeReq.Amount, bakeReq.Times)))
	reformattedMac, err := formattedMac.MarshalBinary()
	if err != nil {
		return "", "", err
	}
	return hex.EncodeToString(formattedMac.Id()), hex.EncodeToString(reformattedMac), nil
}
