package keys

import (
	"github.com/aybabtme/godotto/internal/godoutil"
	"github.com/digitalocean/godo"
	"golang.org/x/net/context"
)

// A Client can interact with the DigitalOcean Domains service.
type Client interface {
	Create(name, publicKey string, opts ...CreateOpt) (Key, error)
	GetByID(int) (Key, error)
	GetByFingerprint(string) (Key, error)
	UpdateByID(int, ...UpdateOpt) (Key, error)
	UpdateByFingerprint(string, ...UpdateOpt) (Key, error)
	DeleteByID(int) error
	DeleteByFingerprint(string) error
	List(context.Context) (<-chan Key, <-chan error)
}

// Key in the DigitalOcean cloud.
type Key interface {
	Struct() *godo.Key
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

// CreateOpt is an optional argument to domains.Create.
type CreateOpt func(*createOpt)

type createOpt struct {
	req *godo.KeyCreateRequest
}

func (svc *client) defaultCreateOpts() *createOpt {
	return &createOpt{
		req: &godo.KeyCreateRequest{},
	}
}

func (svc *client) Create(name, publicKey string, opts ...CreateOpt) (Key, error) {

	opt := svc.defaultCreateOpts()
	opt.req.Name = name
	opt.req.PublicKey = publicKey

	for _, fn := range opts {
		fn(opt)
	}
	d, _, err := svc.g.Keys.Create(opt.req)
	if err != nil {
		return nil, err
	}
	return &key{g: svc.g, d: d}, nil
}

func (svc *client) GetByID(id int) (Key, error) {
	d, _, err := svc.g.Keys.GetByID(id)
	if err != nil {
		return nil, err
	}
	return &key{g: svc.g, d: d}, nil
}

func (svc *client) GetByFingerprint(fp string) (Key, error) {
	d, _, err := svc.g.Keys.GetByFingerprint(fp)
	if err != nil {
		return nil, err
	}
	return &key{g: svc.g, d: d}, nil
}

// UpdateOpt is an optional argument to keys.Update.
type UpdateOpt func(*updateOpt)

func UseGodoKey(req *godo.KeyUpdateRequest) UpdateOpt {
	return func(opt *updateOpt) { opt.req = req }
}

type updateOpt struct {
	req *godo.KeyUpdateRequest
}

func (svc *client) defaultUpdateOpts() *updateOpt {
	return &updateOpt{
		req: &godo.KeyUpdateRequest{},
	}
}

func (svc *client) UpdateByID(id int, opts ...UpdateOpt) (Key, error) {
	opt := svc.defaultUpdateOpts()
	for _, fn := range opts {
		fn(opt)
	}
	d, _, err := svc.g.Keys.UpdateByID(id, opt.req)
	if err != nil {
		return nil, err
	}
	return &key{g: svc.g, d: d}, nil
}

func (svc *client) UpdateByFingerprint(fp string, opts ...UpdateOpt) (Key, error) {
	opt := svc.defaultUpdateOpts()
	for _, fn := range opts {
		fn(opt)
	}
	d, _, err := svc.g.Keys.UpdateByFingerprint(fp, opt.req)
	if err != nil {
		return nil, err
	}
	return &key{g: svc.g, d: d}, nil
}

func (svc *client) DeleteByID(id int) error {
	_, err := svc.g.Keys.DeleteByID(id)
	return err
}

func (svc *client) DeleteByFingerprint(fp string) error {
	_, err := svc.g.Keys.DeleteByFingerprint(fp)
	return err
}

func (svc *client) List(ctx context.Context) (<-chan Key, <-chan error) {
	outc := make(chan Key, 1)
	errc := make(chan error, 1)

	go func() {
		defer close(outc)
		defer close(errc)
		err := godoutil.IterateList(ctx, func(opt *godo.ListOptions) (*godo.Response, error) {
			r, resp, err := svc.g.Keys.List(opt)
			for _, d := range r {
				select {
				case outc <- &key{g: svc.g, d: &d}:
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

type key struct {
	g *godo.Client
	d *godo.Key
}

func (svc *key) Struct() *godo.Key { return svc.d }
