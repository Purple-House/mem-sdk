package main

import (
	"context"
	"time"

	"github.com/Purple-House/mem-sdk/memsdk/maps"
	"github.com/Purple-House/mem-sdk/memsdk/pkg"
)

func main() {

	config := pkg.Config{
		Address:     "localhost:8080",
		Fingerprint: "86f7b7b55c1591c0aafbb9470baff92f1021791ca8f6ee9e372d0986a886be00",
		Timeout:     5 * time.Second,
	}

	client, err := maps.NewSdkOperation(config)
	if err != nil {
		panic(err)
	}

	defer client.Close()
	region := "global"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	router := maps.AddRouterRequest{
		RouterIp:   "127.0.0.1",
		RouterPort: 9090,
		RpcPort:    9092,
		Region:     region,
		Identity:   "jknscdmklm",
	}

	gateways, err := client.Addgateway(ctx, router)
	if err != nil {
		panic(err)
	}
	println("Added Gateway:", gateways.ID, gateways.GatewayPort, gateways.WssPort, gateways.IP, gateways.Identity)
}
