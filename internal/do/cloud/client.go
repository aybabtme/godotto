package cloud

import (
	"log"

	"github.com/aybabtme/godotto/internal/do/cloud/accounts"
	"github.com/aybabtme/godotto/internal/do/cloud/actions"
	"github.com/aybabtme/godotto/internal/do/cloud/domains"
	"github.com/aybabtme/godotto/internal/do/cloud/droplets"
	"github.com/aybabtme/godotto/internal/do/cloud/floatingips"
	"github.com/aybabtme/godotto/internal/do/cloud/images"
	"github.com/aybabtme/godotto/internal/do/cloud/keys"
	"github.com/aybabtme/godotto/internal/do/cloud/regions"
	"github.com/aybabtme/godotto/internal/do/cloud/sizes"
	"github.com/digitalocean/godo"
)

type clientOpts struct {
	g *godo.Client
}

// ClientOpt is an option to configure a client.
type ClientOpt func(*clientOpts)

// UseGodo allows setting options on the underlying *godo.Client.
func UseGodo(client *godo.Client) ClientOpt {
	return func(opt *clientOpts) {
		opt.g = client
	}
}

// A Client knows how to interact with the DigitalOcean cloud.
type Client interface {
	Droplets() droplets.Client
	Accounts() accounts.Client
	Actions() actions.Client
	Domains() domains.Client
	Images() images.Client
	Keys() keys.Client
	Regions() regions.Client
	Sizes() sizes.Client
	FloatingIPs() floatingips.Client
}

// New creates a Client to the DigitalOcean cloud. Options are applied in order,
// such that of multiple options modifying the same value, only the last one's
// effect will be observed.
func New(opts ...ClientOpt) Client {

	opt := &clientOpts{
		g: godo.NewClient(nil),
	}
	log.Printf("before %#v", opt)
	for _, fn := range opts {
		fn(opt)
	}
	log.Printf("after %#v", opt)
	c := &client{
		g:           opt.g,
		droplets:    droplets.New(opt.g),
		actions:     actions.New(opt.g),
		accounts:    accounts.New(opt.g),
		domains:     domains.New(opt.g),
		images:      images.New(opt.g),
		keys:        keys.New(opt.g),
		regions:     regions.New(opt.g),
		sizes:       sizes.New(opt.g),
		floatingips: floatingips.New(opt.g),
	}

	return c
}

type client struct {
	g *godo.Client

	droplets    droplets.Client
	accounts    accounts.Client
	actions     actions.Client
	domains     domains.Client
	images      images.Client
	keys        keys.Client
	regions     regions.Client
	sizes       sizes.Client
	floatingips floatingips.Client
}

func (svc *client) Droplets() droplets.Client       { return svc.droplets }
func (svc *client) Accounts() accounts.Client       { return svc.accounts }
func (svc *client) Actions() actions.Client         { return svc.actions }
func (svc *client) Domains() domains.Client         { return svc.domains }
func (svc *client) Images() images.Client           { return svc.images }
func (svc *client) Keys() keys.Client               { return svc.keys }
func (svc *client) Regions() regions.Client         { return svc.regions }
func (svc *client) Sizes() sizes.Client             { return svc.sizes }
func (svc *client) FloatingIPs() floatingips.Client { return svc.floatingips }
