package accounts

import "github.com/digitalocean/godo"

// A Client can interact with the DigitalOcean Accounts service.
type Client interface {
	Get() (Account, error)
}

// A Account in the DigitalOcean cloud.
type Account interface {
	Struct() *godo.Account
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

func (svc *client) Get() (Account, error) {
	d, _, err := svc.g.Account.Get()
	if err != nil {
		return nil, err
	}
	return &account{g: svc.g, d: d}, nil
}

type account struct {
	g *godo.Client
	d *godo.Account
}

func (svc *account) Struct() *godo.Account { return svc.d }
