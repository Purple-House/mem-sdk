package maps

import (
	pb "github.com/odio4u/agni-schema/maps"
)

type Error struct {
	Code    pb.ErrorCode
	Message string
}

type Capacity struct {
	CPU       int32
	Memory    int32
	Storage   int32
	Bandwidth int32
}

type Gateway struct {
	ID          string
	IP          string
	Address     string
	GatewayPort int32
	WssPort     int32
	Capacity    Capacity
	Identity    string

	Error *Error
}

type Agent struct {
	ID             string
	Domain         string
	GatewayID      string
	GatewayAddress string
	GatewayPort    int32
	WssPort        int32
	GatewayIP      string
	Capacity       Capacity
	Identity       string
	Error          *Error
}

func gatewayFromProto(g *pb.GatewayResponse) *Gateway {
	if g == nil {
		return nil
	}
	return &Gateway{
		ID:          g.GatewayId,
		IP:          g.GatewayIp,
		Address:     g.GatewayAddress,
		WssPort:     g.WssPort,
		GatewayPort: g.GatewayPort,
		Identity:    g.Identity,
		Capacity: Capacity{
			CPU:       g.Capacity.Cpu,
			Memory:    g.Capacity.Memory,
			Storage:   g.Capacity.Storage,
			Bandwidth: g.Capacity.Bandwidth,
		},

		Error: errorFromProto(g.Error),
	}
}

func agentFromProto(a *pb.AgentResponse) *Agent {
	if a == nil {
		return nil
	}
	return &Agent{
		ID:             a.AgentId,
		Domain:         a.AgentDomain,
		GatewayID:      a.GatewayId,
		GatewayAddress: a.GatewayAddress,
		GatewayIP:      a.GatewayIp,
		GatewayPort:    a.GatewayPort,
		WssPort:        a.WssPort,
		Identity:       a.Identity,
		Error:          errorFromProto(a.Error),
	}
}

func errorFromProto(e *pb.Error) *Error {
	if e == nil {
		return nil
	}
	return &Error{
		Code:    e.Code,
		Message: e.Message,
	}
}
