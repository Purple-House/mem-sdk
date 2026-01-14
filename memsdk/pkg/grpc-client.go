package pkg

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	pb "github.com/odio4u/mem-sdk/memsdk/protobuf"
)

type Config struct {
	Address     string
	Fingerprint string
	Timeout     time.Duration
}

type Client struct {
	conn        *grpc.ClientConn
	rpc         pb.MapsClient
	fingerprint string
	timeout     time.Duration
}

func New(cfg Config) (*Client, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			cert, err := x509.ParseCertificate(rawCerts[0])
			if err != nil {
				return err
			}

			fp := sha256.Sum256(cert.Raw)
			expected := cfg.Fingerprint

			if hex.EncodeToString(fp[:]) != expected {
				return errors.New("server fingerprint mismatch")
			}
			return nil
		},
	}

	creds := credentials.NewTLS(tlsConfig)

	conn, err := grpc.Dial(
		cfg.Address,
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:        conn,
		rpc:         pb.NewMapsClient(conn),
		fingerprint: cfg.Fingerprint,
		timeout:     cfg.Timeout,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) ctx(ctx context.Context) (context.Context, context.CancelFunc) {
	if c.timeout > 0 {
		ctx, cancel := context.WithTimeout(ctx, c.timeout)
		ctx = metadata.AppendToOutgoingContext(ctx, "x-fingerprint", c.fingerprint)
		return ctx, cancel
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "x-fingerprint", c.fingerprint)
	return ctx, func() {}
}

// gRPC call wrappers
// ----------------------------

func (c *Client) RegisterGateway(ctx context.Context, req *pb.GatewayPutRequest) (*pb.GatewayResponse, error) {
	ctx, cancel := c.ctx(ctx)
	defer cancel()
	return c.rpc.RegisterGateway(ctx, req)
}

func (c *Client) RegisterAgent(ctx context.Context, req *pb.AgentConnectionRequest) (*pb.AgentResponse, error) {
	ctx, cancel := c.ctx(ctx)
	defer cancel()
	return c.rpc.RegisterAgent(ctx, req)
}

func (c *Client) ResolveGatewayForAgent(ctx context.Context, req *pb.GatewayHandshake) (*pb.MultipleGateways, error) {
	ctx, cancel := c.ctx(ctx)
	defer cancel()
	return c.rpc.ResolveGatewayForAgent(ctx, req)
}

func (c *Client) ResolveGatewayForProxy(ctx context.Context, req *pb.ProxyMapping) (*pb.AgentResponse, error) {
	ctx, cancel := c.ctx(ctx)
	defer cancel()
	return c.rpc.ResolveGatewayForProxy(ctx, req)
}
