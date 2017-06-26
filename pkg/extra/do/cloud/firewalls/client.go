package firewalls

import (
	"context"

	"github.com/digitalocean/godo"
)

type Client interface {
	Create(ctx context.Context, name string, inboundRules []godo.InboundRule, outboundRules []godo.OutboundRule, opts ...CreateOpt) (Firewall, error)
	Get(ctx context.Context, id string) (Firewall, error)
}

type Firewall interface {
	Struct() *godo.Firewall
}

// New creates a Client.
func New(g *godo.Client) Client {
	c := &client{
		g: g,
	}
	return c
}

type client struct {
	g *godo.Client
}

// CreateOpt is an optional argument to firewalls.Create
type CreateOpt func(*createOpt)

type createOpt struct {
	req *godo.FirewallRequest
}

func UseGodoCreate(req *godo.FirewallRequest) CreateOpt {
	return func(opt *createOpt) { opt.req = req }
}

func (svc *client) defaultCreateOpts() *createOpt {
	return &createOpt{
		req: &godo.FirewallRequest{},
	}
}

func (svc *client) Create(ctx context.Context, name string, inboundRules []godo.InboundRule, outboundRules []godo.OutboundRule, opts ...CreateOpt) (Firewall, error) {
	opt := svc.defaultCreateOpts()

	for _, fn := range opts {
		fn(opt)
	}

	opt.req.Name = name
	opt.req.InboundRules = inboundRules
	opt.req.OutboundRules = outboundRules

	f, _, err := svc.g.Firewalls.Create(ctx, opt.req)
	if err != nil {
		return nil, err
	}

	return &firewall{g: svc.g, f: f}, nil
}

func (svc *client) Get(ctx context.Context, id string) (Firewall, error) {
	f, _, err := svc.g.Firewalls.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &firewall{g: svc.g, f: f}, nil
}

type firewall struct {
	g *godo.Client
	f *godo.Firewall
}

func (svc *firewall) Struct() *godo.Firewall { return svc.f }
