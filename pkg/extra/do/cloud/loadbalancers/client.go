package loadbalancers

import (
	"context"

	"github.com/digitalocean/godo"
)

type Client interface {
	Create(ctx context.Context, name, region string, forwardingRules []godo.ForwardingRule, opts ...CreateOpt) (LoadBalancer, error)
}

type LoadBalancer interface {
	Struct() *godo.LoadBalancer
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

// CreateOpt is an optional argument to loadbalancers.Create
type CreateOpt func(*createOpt)

type createOpt struct {
	req *godo.LoadBalancerRequest
}

func UseGodoCreate(req *godo.LoadBalancerRequest) CreateOpt {
	return func(opt *createOpt) { opt.req = req }
}

func (svc *client) defaultCreateOpts() *createOpt {
	return &createOpt{
		req: &godo.LoadBalancerRequest{},
	}
}

func (svc *client) Create(ctx context.Context, name, region string, forwardingRules []godo.ForwardingRule, opts ...CreateOpt) (LoadBalancer, error) {
	opt := svc.defaultCreateOpts()

	for _, fn := range opts {
		fn(opt)
	}

	opt.req.Name = name
	opt.req.Region = region
	opt.req.ForwardingRules = forwardingRules

	l, _, err := svc.g.LoadBalancers.Create(ctx, opt.req)
	if err != nil {
		return nil, err
	}

	return &loadBalancer{g: svc.g, l: l}, nil
}

type loadBalancer struct {
	g *godo.Client
	l *godo.LoadBalancer
}

func (svc *loadBalancer) Struct() *godo.LoadBalancer { return svc.l }
