package actions

import (
	"context"

	"github.com/aybabtme/godotto/internal/godoutil"
	"github.com/digitalocean/godo"
)

// A Client can interact with the DigitalOcean Actions service.
type Client interface {
	Get(context.Context, int) (Action, error)
	List(context.Context) (<-chan Action, <-chan error)
}

// A Action in the DigitalOcean cloud.
type Action interface {
	Struct() *godo.Action
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

func (svc *client) Get(ctx context.Context, id int) (Action, error) {
	d, _, err := svc.g.Actions.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &action{g: svc.g, d: d}, nil
}

func (svc *client) List(ctx context.Context) (<-chan Action, <-chan error) {
	outc := make(chan Action, 1)
	errc := make(chan error, 1)

	go func() {
		defer close(outc)
		defer close(errc)
		err := godoutil.IterateList(ctx, func(ctx context.Context, opt *godo.ListOptions) (*godo.Response, error) {
			r, resp, err := svc.g.Actions.List(ctx, opt)
			for _, d := range r {
				dd := d // copy ranged over variable
				select {
				case outc <- &action{g: svc.g, d: &dd}:
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

type action struct {
	g *godo.Client
	d *godo.Action
}

func (svc *action) Struct() *godo.Action { return svc.d }
