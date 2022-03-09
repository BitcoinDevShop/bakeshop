package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/futurepaul/bakeshop/backend/bakedgood"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	"google.golang.org/protobuf/proto"
	"gopkg.in/macaroon.v2"
)

func (s *serverLnd) createGrpcInterceptor(client lnrpc.LightningClient) error {
	go func() {
		ctx := context.Background()
		rpcMiddlewareClient, err := client.RegisterRPCMiddleware(ctx)
		if err != nil {
			panic(err)
		}

		fmt.Println("Created lnd middleware stream")

		// Register interceptor immediately
		err = rpcMiddlewareClient.Send(&lnrpc.RPCMiddlewareResponse{
			MiddlewareMessage: &lnrpc.RPCMiddlewareResponse_Register{
				Register: &lnrpc.MiddlewareRegistration{
					MiddlewareName:           "subscribe",
					CustomMacaroonCaveatName: "subscribe",
					ReadOnlyMode:             false,
				},
			},
		})
		if err != nil {
			panic(err)
		}
		fmt.Println("Registered lnd middleware stream")

		// Listen for responses
		fmt.Println("Listening to lnd middleware stream")
		go func() {
			for {
				// TODO handle reconnections
				resp, err := rpcMiddlewareClient.Recv()
				if err != nil {
					// TODO return deny msg
					panic(err)
				}
				fmt.Println("Got lnd middleware response")
				fmt.Println(resp)

				// Craft deny msg in case it needs to be denied
				denyResp := &lnrpc.RPCMiddlewareResponse{
					RefMsgId: resp.GetMsgId(),
					MiddlewareMessage: &lnrpc.RPCMiddlewareResponse_Feedback{
						Feedback: &lnrpc.InterceptFeedback{
							Error:           "Invalid",
							ReplaceResponse: false,
						},
					},
				}

				macId, amount, ms, times, err := getMacaroonCaveats(resp.GetRawMacaroon())
				if err != nil {
					// Deny
					err = rpcMiddlewareClient.Send(denyResp)
					if err != nil {
						fmt.Println(err)
					}
					continue
				}

				// Get macaroon/caveat info info
				if amount == 0 {
					// Deny
					err = rpcMiddlewareClient.Send(denyResp)
					if err != nil {
						fmt.Println(err)
					}
					continue
				}
				if ms == 0 {
					// Deny
					err = rpcMiddlewareClient.Send(denyResp)
					if err != nil {
						fmt.Println(err)
					}
					continue
				}
				if times == 0 {
					// Deny
					err = rpcMiddlewareClient.Send(denyResp)
					if err != nil {
						fmt.Println(err)
					}
					continue
				}

				// Parse request details
				// Make sure it was a payment req
				switch resp.InterceptType.(type) {

				case *lnrpc.RPCMiddlewareRequest_StreamAuth:
					// Streams can go through without additional validation
					if resp.GetStreamAuth().GetMethodFullUri() != "/routerrpc.Router/SendPaymentV2" {
						// Deny
						err = rpcMiddlewareClient.Send(denyResp)
						if err != nil {
							fmt.Println(err)
						}
						continue
					}

				case *lnrpc.RPCMiddlewareRequest_Request:
					if resp.GetRequest().GetMethodFullUri() != "/routerrpc.Router/SendPaymentV2" {
						// Deny
						err = rpcMiddlewareClient.Send(denyResp)
						if err != nil {
							fmt.Println(err)
						}
						continue
					}

					if resp.GetRequest().GetTypeName() != "routerrpc.SendPaymentRequest" {
						// Deny
						err = rpcMiddlewareClient.Send(denyResp)
						if err != nil {
							fmt.Println(err)
						}
						continue
					}

					// Parse the actual request
					parsedRequest := &routerrpc.SendPaymentRequest{}
					reqBytes := resp.GetRequest().GetSerialized()
					err = proto.Unmarshal(reqBytes, parsedRequest)
					if err != nil {
						fmt.Println(err.Error())

						// Deny
						err = rpcMiddlewareClient.Send(denyResp)
						if err != nil {
							fmt.Println(err)
						}
						continue
					}

					// Verify the request matches details of macaroon
					if parsedRequest.GetAmt() != amount {
						fmt.Println("Subscriber tried to pull invalid amount")

						// Deny
						err = rpcMiddlewareClient.Send(denyResp)
						if err != nil {
							fmt.Println(err)
						}
						continue
					}

					// pull all payments that have been made
					payments, err := s.db.GetPaymentsByMacaroonId(string(macId))
					if err != nil {
						fmt.Println("Could not get payments by macaroon id")
						fmt.Println(err)
						// Deny
						err = rpcMiddlewareClient.Send(denyResp)
						if err != nil {
							fmt.Println(err)
						}
						continue
					}

					// validate times
					if int64(len(payments)) > times {
						fmt.Println("Payment times exceeded")

						// Deny
						err = rpcMiddlewareClient.Send(denyResp)
						if err != nil {
							fmt.Println(err)
						}
						continue
					}

					// TODO validate aggregated amount
					// Maybe if interval and times is missing, it's an "unlimited withdraws up to amount" scenario
					// TODO we need to handle fees into the payment

					// validate time interval
					// get the latest payment and make sure it's not within interval
					// TODO if productionalized, need a surefire way to lock requests while this one processes
					if len(payments) > 0 {
						latestPaymentTime := payments[len(payments)-1].CreatedAt
						if latestPaymentTime.Add(time.Duration(ms) * time.Millisecond).After(time.Now()) {
							fmt.Println("Too many payments in given interval")

							// Deny
							err = rpcMiddlewareClient.Send(denyResp)
							if err != nil {
								fmt.Println(err)
							}
							continue
						}
					}

					// Approved payment
					fmt.Println("Approving rpc middleware..")
					err = rpcMiddlewareClient.Send(&lnrpc.RPCMiddlewareResponse{
						RefMsgId: resp.GetMsgId(),
						MiddlewareMessage: &lnrpc.RPCMiddlewareResponse_Feedback{
							Feedback: &lnrpc.InterceptFeedback{
								Error:           "",
								ReplaceResponse: false,
							},
						},
					})
					if err != nil {
						// TODO process error, it may have actually approved
						fmt.Println("Could not approve payment")
						fmt.Println(err)
					}
					fmt.Println("Sent lnd middleware response")

					// save successful payment in db
					// TODO make sure this doesn't fail
					fmt.Println("Saving payment to db")
					_, err = s.db.CreatePayment(bakedgood.Payment{
						MacaroonId: string(macId),
						Amount:     amount,
					})
					if err != nil {
						fmt.Println("Could not save payment")
						fmt.Println(err)
					}
					fmt.Println("Saved payment to db")
					continue

				// Streams can go through without additional validation
				case *lnrpc.RPCMiddlewareRequest_Response:
					if resp.GetResponse().GetMethodFullUri() != "/routerrpc.Router/SendPaymentV2" {
						// Deny
						err = rpcMiddlewareClient.Send(denyResp)
						if err != nil {
							fmt.Println(err)
						}
						continue
					}

				}

				// Approved transaction
				fmt.Println("Approving rpc middleware..")
				err = rpcMiddlewareClient.Send(&lnrpc.RPCMiddlewareResponse{
					RefMsgId: resp.GetMsgId(),
					MiddlewareMessage: &lnrpc.RPCMiddlewareResponse_Feedback{
						Feedback: &lnrpc.InterceptFeedback{
							Error:           "",
							ReplaceResponse: false,
						},
					},
				})
				if err != nil {
					// TODO process error, it may have actually approved
					fmt.Println("Could not approve payment")
					fmt.Println(err)
				}
				fmt.Println("Sent lnd middleware response")
			}
		}()
	}()
	return nil
}

func getMacaroonCaveats(mac []byte) (id []byte, amt int64, ms int64, times int64, err error) {
	formattedMac := macaroon.Macaroon{}
	err = formattedMac.UnmarshalBinary(mac)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	cav := formattedMac.Caveats()[0].Id
	fmt.Println("cavs: " + string(cav))

	cavFields := strings.Fields(string(cav))
	for _, cavField := range cavFields {
		cavFieldSplit := strings.Split(cavField, ":")
		if len(cavFieldSplit) != 2 {
			continue
		}

		switch cavFieldSplit[0] {
		case "ms":
			ms, err = strconv.ParseInt(cavFieldSplit[1], 10, 64)
			if err != nil {
				continue
			}

		case "amount":
			amt, err = strconv.ParseInt(cavFieldSplit[1], 10, 64)
			if err != nil {
				continue
			}

		case "times":
			times, err = strconv.ParseInt(cavFieldSplit[1], 10, 64)
			if err != nil {
				continue
			}
		}
	}

	return formattedMac.Id(), amt, ms, times, nil
}
