package droplets

import (
	"github.com/aybabtme/godotto/internal/godoutil"
	"github.com/digitalocean/godo"
	"golang.org/x/net/context"
)

// A Client can interact with the DigitalOcean Droplets service.
type Client interface {
	Create(name, region, size, image string, opts ...CreateOpt) (Droplet, error)
	Get(id int) (Droplet, error)
	Delete(id int) error
	List(ctx context.Context) (<-chan Droplet, <-chan error)
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

func (svc *client) Create(name, region, size, image string, opts ...CreateOpt) (Droplet, error) {
	opt := svc.defaultCreateOpts()
	opt.req.Name = name
	opt.req.Size = size
	opt.req.Region = region
	opt.req.Image.Slug = image

	for _, fn := range opts {
		fn(opt)
	}
	d, _, err := svc.g.Droplets.Create(opt.req)
	if err != nil {
		return nil, err
	}
	return &droplet{g: svc.g, d: d}, nil
}

func (svc *client) Get(id int) (Droplet, error) {
	d, _, err := svc.g.Droplets.Get(id)
	if err != nil {
		return nil, err
	}
	return &droplet{g: svc.g, d: d}, nil
}

func (svc *client) Delete(id int) error {
	_, err := svc.g.Droplets.Delete(id)
	return err
}

func (svc *client) List(ctx context.Context) (<-chan Droplet, <-chan error) {
	outc := make(chan Droplet, 1)
	errc := make(chan error, 1)

	go func() {
		defer close(outc)
		defer close(errc)
		err := godoutil.IterateList(ctx, func(opt *godo.ListOptions) (*godo.Response, error) {
			r, resp, err := svc.g.Droplets.List(opt)
			for _, d := range r {
				select {
				case outc <- &droplet{g: svc.g, d: &d}:
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

type droplet struct {
	g *godo.Client
	d *godo.Droplet
}

func (svc *droplet) Struct() *godo.Droplet { return svc.d }
