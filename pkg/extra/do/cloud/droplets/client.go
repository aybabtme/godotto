package droplets

import (
	"context"

	"github.com/aybabtme/godotto/pkg/extra/godoutil"
	"github.com/digitalocean/godo"
)

// A Client can interact with the DigitalOcean Droplets service.
type Client interface {
	Create(ctx context.Context, name, region, size, image string, opts ...CreateOpt) (Droplet, error)
	CreateMultiple(ctx context.Context, names []string, region, size, image string, opts ...CreateMultipleOpt) ([]Droplet, error)
	Get(ctx context.Context, id int) (Droplet, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) (<-chan Droplet, <-chan error)
	Actions() ActionClient
}

// A Droplet in the DigitalOcean cloud.
type Droplet interface {
	Struct() *godo.Droplet
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

// CreateOpt is an optional argument to droplets.Create.
type CreateOpt func(*createOpt)

func UseGodoCreate(req *godo.DropletCreateRequest) CreateOpt {
	return func(opt *createOpt) { opt.req = req }
}

type createOpt struct {
	req *godo.DropletCreateRequest
}

func (svc *client) defaultCreateOpts() *createOpt {
	return &createOpt{
		req: &godo.DropletCreateRequest{
			Image: godo.DropletCreateImage{},
		},
	}
}

type CreateMultipleOpt func(*createMultipleOpt)

func UseGodoMultiCreate(req *godo.DropletMultiCreateRequest) CreateMultipleOpt {
	return func(opt *createMultipleOpt) { opt.req = req }
}

type createMultipleOpt struct {
	req *godo.DropletMultiCreateRequest
}

func (svc *client) defaultCreateMultipleOpts() *createMultipleOpt {
	return &createMultipleOpt{
		req: &godo.DropletMultiCreateRequest{
			Image: godo.DropletCreateImage{},
		},
	}
}

func (svc *client) Create(ctx context.Context, name, region, size, image string, opts ...CreateOpt) (Droplet, error) {
	opt := svc.defaultCreateOpts()
	for _, fn := range opts {
		fn(opt)
	}
	opt.req.Name = name
	opt.req.Size = size
	opt.req.Region = region
	opt.req.Image.Slug = image

	d, resp, err := svc.g.Droplets.Create(ctx, opt.req)
	if err != nil {
		return nil, err
	}

	return &droplet{g: svc.g, d: d}, godoutil.WaitForActions(ctx, svc.g, resp.Links)
}

func (svc *client) CreateMultiple(ctx context.Context, names []string, region, size, image string, opts ...CreateMultipleOpt) ([]Droplet, error) {
	opt := svc.defaultCreateMultipleOpts()
	for _, fn := range opts {
		fn(opt)
	}
	opt.req.Names = names
	opt.req.Size = size
	opt.req.Region = region
	opt.req.Image.Slug = image

	r, _, err := svc.g.Droplets.CreateMultiple(ctx, opt.req)
	if err != nil {
		return nil, err
	}

	droplets := make([]Droplet, 0, len(r))
	for _, d := range r {
		dd := d
		droplets = append(droplets, &droplet{g: svc.g, d: &dd})
	}

	return droplets, nil
}

func (svc *client) Get(ctx context.Context, id int) (Droplet, error) {
	d, _, err := svc.g.Droplets.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &droplet{g: svc.g, d: d}, nil
}

func (svc *client) Delete(ctx context.Context, id int) error {
	resp, err := svc.g.Droplets.Delete(ctx, id)
	if err != nil {
		return err
	}
	return godoutil.WaitForActions(ctx, svc.g, resp.Links)
}

func (svc *client) List(ctx context.Context) (<-chan Droplet, <-chan error) {
	outc := make(chan Droplet, 1)
	errc := make(chan error, 1)

	go func() {
		defer close(outc)
		defer close(errc)
		err := godoutil.IterateList(ctx, func(ctx context.Context, opt *godo.ListOptions) (*godo.Response, error) {
			r, resp, err := svc.g.Droplets.List(ctx, opt)
			for _, d := range r {
				dd := d // copy ranged over variable
				select {
				case outc <- &droplet{g: svc.g, d: &dd}:
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

type droplet struct {
	g *godo.Client
	d *godo.Droplet
}

func (svc *droplet) Struct() *godo.Droplet { return svc.d }
