package cloud

import (
	"github.com/aybabtme/godotto/internal/do/cloud/droplets"
	"github.com/digitalocean/godo"
)

// ClientOpt is an option to configure a client.
type ClientOpt func(*client)

// UseGodoOpts allows setting options on the underlying *godo.Client.
func UseGodoOpts(fn func(*godo.Client)) ClientOpt {
	return func(svc *client) {
		fn(svc.g)
	}
}

// A Client knows how to interact with the DigitalOcean cloud.
type Client interface {
	Droplets() droplets.Client
}

// New creates a Client to the DigitalOcean cloud. Options are applied in order,
// such that of multiple options modifying the same value, only the last one's
// effect will be observed.
func New(opts ...ClientOpt) Client {

	g := godo.NewClient(nil)

	c := &client{
		g:        g,
		droplets: droplets.New(g),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

type client struct {
	g *godo.Client

	droplets droplets.Client
}

func (svc *client) Droplets() droplets.Client { return svc.droplets }
