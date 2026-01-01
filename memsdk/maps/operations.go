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

func New(cfg Config) (*Client, error) {
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
