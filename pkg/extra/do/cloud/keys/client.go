package keys

import (
	"context"

	"github.com/aybabtme/godotto/pkg/extra/godoutil"
	"github.com/digitalocean/godo"
)

// A Client can interact with the DigitalOcean Domains service.
type Client interface {
	Create(ctx context.Context, name, publicKey string, opts ...CreateOpt) (Key, error)
	GetByID(context.Context, int) (Key, error)
	GetByFingerprint(context.Context, string) (Key, error)
	UpdateByID(context.Context, int, ...UpdateOpt) (Key, error)
	UpdateByFingerprint(context.Context, string, ...UpdateOpt) (Key, error)
	DeleteByID(context.Context, int) error
	DeleteByFingerprint(context.Context, string) error
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

func (svc *client) Create(ctx context.Context, name, publicKey string, opts ...CreateOpt) (Key, error) {

	opt := svc.defaultCreateOpts()
	opt.req.Name = name
	opt.req.PublicKey = publicKey

	for _, fn := range opts {
		fn(opt)
	}
	d, _, err := svc.g.Keys.Create(ctx, opt.req)
	if err != nil {
		return nil, err
	}
	return &key{g: svc.g, d: d}, nil
}

func (svc *client) GetByID(ctx context.Context, id int) (Key, error) {
	d, _, err := svc.g.Keys.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &key{g: svc.g, d: d}, nil
}

func (svc *client) GetByFingerprint(ctx context.Context, fp string) (Key, error) {
	d, _, err := svc.g.Keys.GetByFingerprint(ctx, fp)
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

func (svc *client) UpdateByID(ctx context.Context, id int, opts ...UpdateOpt) (Key, error) {
	opt := svc.defaultUpdateOpts()
	for _, fn := range opts {
		fn(opt)
	}
	d, _, err := svc.g.Keys.UpdateByID(ctx, id, opt.req)
	if err != nil {
		return nil, err
	}
	return &key{g: svc.g, d: d}, nil
}

func (svc *client) UpdateByFingerprint(ctx context.Context, fp string, opts ...UpdateOpt) (Key, error) {
	opt := svc.defaultUpdateOpts()
	for _, fn := range opts {
		fn(opt)
	}
	d, _, err := svc.g.Keys.UpdateByFingerprint(ctx, fp, opt.req)
	if err != nil {
		return nil, err
	}
	return &key{g: svc.g, d: d}, nil
}

func (svc *client) DeleteByID(ctx context.Context, id int) error {
	_, err := svc.g.Keys.DeleteByID(ctx, id)
	return err
}

func (svc *client) DeleteByFingerprint(ctx context.Context, fp string) error {
	_, err := svc.g.Keys.DeleteByFingerprint(ctx, fp)
	return err
}

func (svc *client) List(ctx context.Context) (<-chan Key, <-chan error) {
	outc := make(chan Key, 1)
	errc := make(chan error, 1)

	go func() {
		defer close(outc)
		defer close(errc)
		err := godoutil.IterateList(ctx, func(ctx context.Context, opt *godo.ListOptions) (*godo.Response, error) {
			r, resp, err := svc.g.Keys.List(ctx, opt)
			for _, d := range r {
				dd := d // copy ranged over variable
				select {
				case outc <- &key{g: svc.g, d: &dd}:
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
