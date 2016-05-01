package sizes

import (
	"github.com/aybabtme/godotto/internal/godoutil"
	"github.com/digitalocean/godo"
	"golang.org/x/net/context"
)

// A Client can interact with the DigitalOcean sizes service.
type Client interface {
	List(ctx context.Context) (<-chan Size, <-chan error)
}

// A Size in the DigitalOcean cloud.
type Size interface {
	Struct() *godo.Size
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

func (svc *client) List(ctx context.Context) (<-chan Size, <-chan error) {
	outc := make(chan Size, 1)
	errc := make(chan error, 1)

	go func() {
		defer close(outc)
		defer close(errc)
		err := godoutil.IterateList(ctx, func(opt *godo.ListOptions) (*godo.Response, error) {
			r, resp, err := svc.g.Sizes.List(opt)
			for _, d := range r {
				select {
				case outc <- &size{g: svc.g, d: &d}:
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

type size struct {
	g *godo.Client
	d *godo.Size
}

func (svc *size) Struct() *godo.Size { return svc.d }
