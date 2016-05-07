package regions

import (
	"github.com/aybabtme/godotto/internal/godoutil"
	"github.com/digitalocean/godo"
	"golang.org/x/net/context"
)

// A Client can interact with the DigitalOcean Regions service.
type Client interface {
	List(ctx context.Context) (<-chan Region, <-chan error)
}

// A Region in the DigitalOcean cloud.
type Region interface {
	Struct() *godo.Region
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

func (svc *client) List(ctx context.Context) (<-chan Region, <-chan error) {
	outc := make(chan Region, 1)
	errc := make(chan error, 1)

	go func() {
		defer close(outc)
		defer close(errc)
		err := godoutil.IterateList(ctx, func(opt *godo.ListOptions) (*godo.Response, error) {
			r, resp, err := svc.g.Regions.List(opt)
			for _, d := range r {
				dd := d // copy ranged over variable
				select {
				case outc <- &region{g: svc.g, d: &dd}:
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

type region struct {
	g *godo.Client
	d *godo.Region
}

func (svc *region) Struct() *godo.Region { return svc.d }
