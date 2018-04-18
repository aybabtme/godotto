package loadbalancers

import (
	"context"

	"github.com/aybabtme/godotto/pkg/extra/godoutil"
	"github.com/digitalocean/godo"
)

type Client interface {
	Create(ctx context.Context, name, region string, forwardingRules []godo.ForwardingRule, opts ...CreateOpt) (LoadBalancer, error)
	Get(ctx context.Context, id string) (LoadBalancer, error)
	Update(ctx context.Context, id string, opts ...UpdateOpt) (LoadBalancer, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) (<-chan LoadBalancer, <-chan error)
	AddDroplets(ctx context.Context, lbId string, dropletIDs ...int) error
	RemoveDroplets(ctx context.Context, lbId string, dropletIDs ...int) error
	AddForwardingRules(ctx context.Context, lbID string, rules ...godo.ForwardingRule) error
	RemoveForwardingRules(ctx context.Context, lbID string, rules ...godo.ForwardingRule) error
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

// UpdateOpt is an optional argument to Update.
type UpdateOpt func(*updateOpt)

func UseGodoLoadBalancer(req *godo.LoadBalancerRequest) UpdateOpt {
	return func(opt *updateOpt) { opt.req = req }
}

type updateOpt struct {
	req *godo.LoadBalancerRequest
}

func (svc *client) defaultUpdateOpts() *updateOpt {
	return &updateOpt{
		req: &godo.LoadBalancerRequest{},
	}
}

func (svc *client) Update(ctx context.Context, id string, opts ...UpdateOpt) (LoadBalancer, error) {
	opt := svc.defaultUpdateOpts()
	for _, fn := range opts {
		fn(opt)
	}

	l, _, err := svc.g.LoadBalancers.Update(ctx, id, opt.req)
	if err != nil {
		return nil, err
	}

	return &loadBalancer{g: svc.g, l: l}, nil
}

func (svc *client) Get(ctx context.Context, id string) (LoadBalancer, error) {
	l, _, err := svc.g.LoadBalancers.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &loadBalancer{g: svc.g, l: l}, nil
}

func (svc *client) Delete(ctx context.Context, id string) error {
	_, err := svc.g.LoadBalancers.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (svc *client) List(ctx context.Context) (<-chan LoadBalancer, <-chan error) {
	outc := make(chan LoadBalancer, 1)
	errc := make(chan error, 1)

	go func() {
		defer close(outc)
		defer close(errc)
		err := godoutil.IterateList(ctx, func(ctx context.Context, opt *godo.ListOptions) (*godo.Response, error) {
			r, resp, err := svc.g.LoadBalancers.List(ctx, opt)
			for _, l := range r {
				ll := l // copy ranged over variable
				select {
				case outc <- &loadBalancer{g: svc.g, l: &ll}:
				case <-ctx.Done():
					return resp, err
				}
			}
			return resp, err
		})
		if err != nil {
			errc <- err
		}
	}()

	return outc, errc
}

func (svc *client) AddDroplets(ctx context.Context, lbId string, dropletIds ...int) error {
	_, err := svc.g.LoadBalancers.AddDroplets(ctx, lbId, dropletIds...)
	if err != nil {
		return err
	}

	return nil
}

func (svc *client) RemoveDroplets(ctx context.Context, lbId string, dropletIds ...int) error {
	_, err := svc.g.LoadBalancers.RemoveDroplets(ctx, lbId, dropletIds...)
	if err != nil {
		return err
	}

	return nil
}

func (svc *client) AddForwardingRules(ctx context.Context, lbId string, rules ...godo.ForwardingRule) error {
	_, err := svc.g.LoadBalancers.AddForwardingRules(ctx, lbId, rules...)
	if err != nil {
		return err
	}

	return nil
}

func (svc *client) RemoveForwardingRules(ctx context.Context, lbId string, rules ...godo.ForwardingRule) error {
	_, err := svc.g.LoadBalancers.RemoveForwardingRules(ctx, lbId, rules...)
	if err != nil {
		return err
	}

	return nil
}

type loadBalancer struct {
	g *godo.Client
	l *godo.LoadBalancer
}

func (svc *loadBalancer) Struct() *godo.LoadBalancer { return svc.l }
