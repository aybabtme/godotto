package floatingips

import (
	"context"

	"github.com/aybabtme/godotto/pkg/extra/godoutil"
	"github.com/digitalocean/godo"
)

// A Client can interact with the DigitalOcean FloatingIPs service.
type Client interface {
	Create(ctx context.Context, region string, opts ...CreateOpt) (FloatingIP, error)
	Get(context.Context, string) (FloatingIP, error)
	Delete(context.Context, string) error
	List(context.Context) (<-chan FloatingIP, <-chan error)
	Actions() ActionClient
}

// FloatingIP in the DigitalOcean cloud.
type FloatingIP interface {
	Struct() *godo.FloatingIP
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

// CreateOpt is an optional argument to floatingips.Create.
type CreateOpt func(*createOpt)

type createOpt struct {
	req *godo.FloatingIPCreateRequest
}

func UseGodoFloatingIP(req *godo.FloatingIPCreateRequest) CreateOpt {
	return func(opt *createOpt) { opt.req = req }
}

func (svc *client) defaultCreateOpts() *createOpt {
	return &createOpt{
		req: &godo.FloatingIPCreateRequest{},
	}
}

func (svc *client) Create(ctx context.Context, region string, opts ...CreateOpt) (FloatingIP, error) {

	opt := svc.defaultCreateOpts()
	opt.req.Region = region
	for _, fn := range opts {
		fn(opt)
	}
	d, resp, err := svc.g.FloatingIPs.Create(ctx, opt.req)
	if err != nil {
		return nil, err
	}
	return &floatingIP{g: svc.g, d: d}, godoutil.WaitForActions(ctx, svc.g, resp.Links)
}

func (svc *client) Get(ctx context.Context, ip string) (FloatingIP, error) {
	d, _, err := svc.g.FloatingIPs.Get(ctx, ip)
	if err != nil {
		return nil, err
	}
	return &floatingIP{g: svc.g, d: d}, nil
}

func (svc *client) Delete(ctx context.Context, ip string) error {
	resp, err := svc.g.FloatingIPs.Delete(ctx, ip)
	if err != nil {
		return err
	}

	return godoutil.WaitForActions(ctx, svc.g, resp.Links)
}

func (svc *client) List(ctx context.Context) (<-chan FloatingIP, <-chan error) {
	outc := make(chan FloatingIP, 1)
	errc := make(chan error, 1)

	go func() {
		defer close(outc)
		defer close(errc)
		err := godoutil.IterateList(ctx, func(ctx context.Context, opt *godo.ListOptions) (*godo.Response, error) {
			r, resp, err := svc.g.FloatingIPs.List(ctx, opt)
			for _, d := range r {
				dd := d // copy ranged over variable
				select {
				case outc <- &floatingIP{g: svc.g, d: &dd}:
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

func (svc *client) Actions() ActionClient {
	return &actionClient{g: svc.g}
}

type floatingIP struct {
	g *godo.Client
	d *godo.FloatingIP
}

func (svc *floatingIP) Struct() *godo.FloatingIP { return svc.d }
