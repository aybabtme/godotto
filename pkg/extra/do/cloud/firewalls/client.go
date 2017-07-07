package firewalls

import (
	"context"

	"github.com/aybabtme/godotto/internal/godoutil"
	"github.com/digitalocean/godo"
)

type Client interface {
	Create(ctx context.Context, name string, inboundRules []godo.InboundRule, outboundRules []godo.OutboundRule, opts ...CreateOpt) (Firewall, error)
	List(ctx context.Context) (<-chan Firewall, <-chan error)
	Get(ctx context.Context, id string) (Firewall, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, opts ...UpdateOpt) (Firewall, error)
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

func (svc *client) Delete(ctx context.Context, id string) error {
	_, err := svc.g.Firewalls.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (svc *client) List(ctx context.Context) (<-chan Firewall, <-chan error) {
	outc := make(chan Firewall, 1)
	errc := make(chan error, 1)

	go func() {
		defer close(outc)
		defer close(errc)
		err := godoutil.IterateList(ctx, func(ctx context.Context, opt *godo.ListOptions) (*godo.Response, error) {
			r, resp, err := svc.g.Firewalls.List(ctx, opt)
			for _, f := range r {
				ff := f
				select {
				case outc <- &firewall{g: svc.g, f: &ff}:
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

// UpdateOpt is an optional argument to firewalls.Update
type UpdateOpt func(*updateOpt)

type updateOpt struct {
	req *godo.FirewallRequest
}

func (svc *client) defaultUpdateOpts() *updateOpt {
	return &updateOpt{
		req: &godo.FirewallRequest{},
	}
}

func UseGodoFirewall(req *godo.FirewallRequest) UpdateOpt {
	return func(opt *updateOpt) { opt.req = req }
}

func (svc *client) Update(ctx context.Context, id string, opts ...UpdateOpt) (Firewall, error) {
	opt := svc.defaultUpdateOpts()
	for _, fn := range opts {
		fn(opt)
	}

	f, _, err := svc.g.Firewalls.Update(ctx, id, opt.req)
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
