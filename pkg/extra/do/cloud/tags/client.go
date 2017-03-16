package tags

import (
	"context"

	"github.com/digitalocean/godo"
)

// A Client can interact with the DigitalOcean Tags service
type Client interface {
	Create(ctx context.Context, name string, opt ...CreateOpt) (Tag, error)
	//Get(ctx context.Context, name string) (Tag, error)
	//Delete(ctx context.Context, name string) error
	//List(ctx context.Context) (<-chan Tag, <-chan error)
}

type Tag interface {
	Struct() *godo.Tag
}

func New(g *godo.Client) Client {
	c := &client{
		g: g,
	}
	return c
}

type client struct {
	g *godo.Client
}

type CreateOpt func(*createOpt)

func (svc *client) defaultCreateOpts() *createOpt {
	return &createOpt{
		req: &godo.TagCreateRequest{
			Name: "",
		},
	}
}

func UseGodoCreate(req *godo.TagCreateRequest) CreateOpt {
	return func(opt *createOpt) { opt.req = req }
}

type createOpt struct {
	req *godo.TagCreateRequest
}

type tag struct {
	g *godo.Client
	t *godo.Tag
}

func (svc *tag) Struct() *godo.Tag { return svc.t }

func (svc *client) Create(ctx context.Context, name string, opts ...CreateOpt) (Tag, error) {
	opt := svc.defaultCreateOpts()
	for _, fn := range opts {
		fn(opt)
	}

	opt.req.Name = name

	t, _, err := svc.g.Tags.Create(ctx, opt.req)
	if err != nil {
		return nil, err
	}

	return &tag{g: svc.g, t: t}, nil
}
