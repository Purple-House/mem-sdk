package maps

import (
	"context"
	"errors"

	"github.com/odio4u/mem-sdk/memsdk/pkg"
	pb "github.com/odio4u/mem-sdk/memsdk/protobuf"
)

type Client struct {
	grpc *pkg.Client
}

type AddRouterRequest struct {
	Region     string
	RouterIp   string
	RouterPort int32
	Identity   string
	RpcPort    int32
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

func (c *Client) Addgateway(ctx context.Context, router AddRouterRequest) (Gateway, error) {
	gateway := &pb.GatewayPutRequest{
		Region:             router.Region,
		GatewayIp:          router.RouterIp,
		GatewayPort:        router.RouterPort,
		VerifiableCredHash: router.Identity,
		WssPort:            router.RpcPort,
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

func (c *Client) ConnectAgent(ctx context.Context, agent_domain string, gateway_id string, creadhash string, region string) (Agent, error) {

	agentReq := &pb.AgentConnectionRequest{
		AgentDomain:        agent_domain,
		GatewayId:          gateway_id,
		VerifiableCredHash: creadhash,
		Region:             region,
	}

	resp, err := c.grpc.RegisterAgent(ctx, agentReq)
	if err != nil {
		return Agent{}, err
	}
	return *agentFromProto(resp), nil
}

func (c *Client) GetGatewayInfo(ctx context.Context, region string) ([]Gateway, error) {
	agentReq := &pb.GatewayHandshake{
		Region: region,
	}

	resp, err := c.grpc.ResolveGatewayForAgent(ctx, agentReq)
	if err != nil {
		return nil, err
	}

	var out []Gateway
	for _, g := range resp.Gateways {
		out = append(out, *gatewayFromProto(g))
	}
	return out, nil
}

func (c *Client) GetAgentProxyMapping(ctx context.Context, region string, domain string) (Agent, error) {

	agentReq := &pb.ProxyMapping{
		AgentDomain: domain,
		Region:      region,
	}

	resp, err := c.grpc.ResolveGatewayForProxy(ctx, agentReq)
	if err != nil {
		return Agent{}, err
	}
	return *agentFromProto(resp), nil
}
