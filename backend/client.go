package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	"github.com/lightningnetwork/lnd/macaroons"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/macaroon.v2"
)

func CreateLNDClient(macaroonPath, macaroonHex, tlsCertPath, ipAddr string) (lnrpc.LightningClient, []byte, error) {
	fmt.Println("Connecting to LND node...")

	tlsBytes, err := ioutil.ReadFile(tlsCertPath)
	if err != nil {
		return nil, nil, err
	}

	tlsCreds, err := credentials.NewClientTLSFromFile(tlsCertPath, "")
	if err != nil {
		fmt.Println("Cannot get node tls credentials", err)
		return nil, nil, err
	}

	var macaroonData []byte
	if macaroonHex != "" {
		macBytes, err := hex.DecodeString(macaroonHex)
		if err != nil {
			return nil, nil, err
		}
		macaroonData = macBytes
	} else if macaroonPath != "" {
		macBytes, err := ioutil.ReadFile(macaroonPath)
		if err != nil {
			return nil, nil, err
		}
		macaroonData = macBytes // make it available outside of the else if block
	} else {
		return nil, nil, fmt.Errorf("LND macaroon is missing")
	}

	mac := &macaroon.Macaroon{}
	if err := mac.UnmarshalBinary(macaroonData); err != nil {
		return nil, nil, err
	}

	macCred, err := macaroons.NewMacaroonCredential(mac)
	if err != nil {
		return nil, nil, err
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(tlsCreds),
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(macCred),
	}

	conn, err := grpc.Dial(ipAddr, opts...)
	if err != nil {
		fmt.Println("cannot dial to lnd", err)
		return nil, nil, err
	}
	client := lnrpc.NewLightningClient(conn)
	fmt.Println("Connected to LND node")
	return client, tlsBytes, nil
}

func CreateLNDRouterClient(macaroonPath, macaroonHex, tlsCertPath, ipAddr string) (routerrpc.RouterClient, []byte, error) {
	fmt.Println("Connecting to LND node...")

	tlsBytes, err := ioutil.ReadFile(tlsCertPath)
	if err != nil {
		return nil, nil, err
	}

	tlsCreds, err := credentials.NewClientTLSFromFile(tlsCertPath, "")
	if err != nil {
		fmt.Println("Cannot get node tls credentials", err)
		return nil, nil, err
	}

	var macaroonData []byte
	if macaroonHex != "" {
		macBytes, err := hex.DecodeString(macaroonHex)
		if err != nil {
			return nil, nil, err
		}
		macaroonData = macBytes
	} else if macaroonPath != "" {
		macBytes, err := ioutil.ReadFile(macaroonPath)
		if err != nil {
			return nil, nil, err
		}
		macaroonData = macBytes // make it available outside of the else if block
	} else {
		return nil, nil, fmt.Errorf("LND macaroon is missing")
	}

	mac := &macaroon.Macaroon{}
	if err := mac.UnmarshalBinary(macaroonData); err != nil {
		return nil, nil, err
	}

	macCred, err := macaroons.NewMacaroonCredential(mac)
	if err != nil {
		return nil, nil, err
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(tlsCreds),
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(macCred),
	}

	conn, err := grpc.Dial(ipAddr, opts...)
	if err != nil {
		fmt.Println("cannot dial to lnd", err)
		return nil, nil, err
	}
	client := routerrpc.NewRouterClient(conn)
	fmt.Println("Connected to LND node")
	return client, tlsBytes, nil
}
