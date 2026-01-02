package maps

import (
	"context"
	"errors"

	"github.com/Purple-House/mem-sdk/memsdk/pkg"
	pb "github.com/Purple-House/mem-sdk/memsdk/protobuf"
)

type Client struct {
	grpc *pkg.Client
}

func NewSdkOperation(cfg pkg.Config) (*Client, error) {
	if cfg.Address == "" {
		return nil, errors.New("address must be set")
	}

	cli, err := pkg.New(pkg.Config{
		Address:     cfg.Address,
		Fingerprint: cfg.Fingerprint,
		Timeout:     cfg.Timeout,
	})
	if err != nil {
		return nil, err
	}

	return &Client{grpc: cli}, nil
}

func (c *Client) Close() error {
	return c.grpc.Close()
}

func (c *Client) ResolveGatewayForAgent(ctx context.Context, region string) ([]Gateway, error) {
	res, err := c.grpc.ResolveGatewayForAgent(ctx, &pb.GatewayHandshake{Region: region})
	if err != nil {
		return nil, err
	}

	var out []Gateway
	for _, g := range res.Gateways {
		out = append(out, *gatewayFromProto(g))
	}
	return out, nil
}

func (c *Client) Addgateway(ctx context.Context, gateway_ip string, gateway_port int32, gateway_domain string) (Gateway, error) {
	gateway := &pb.GatewayPutRequest{
		GatewayIp:     gateway_ip,
		GatewayDomain: gateway_domain,
		GatewayPort:   gateway_port,
		Region:        "global",
		Capacity: &pb.Capacity{
			Cpu:       1,
			Bandwidth: 20,
			Memory:    4096,
			Storage:   40960,
		},
	}

	resp, err := c.grpc.RegisterGateway(ctx, gateway)
	if err != nil {
		return Gateway{}, err
	}

	return *gatewayFromProto(resp), nil
}
