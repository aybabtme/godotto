package tags

import (
	"context"

	"github.com/digitalocean/godo"
)

// A Client can interact with the DigitalOcean Tags service
type Client interface {
	Create(ctx context.Context, name string, opt ...CreateOpt) (Tag, error)
	Get(ctx context.Context, name string) (Tag, error)
	//Delete(ctx context.Context, name string) error
	//List(ctx context.Context) (<-chan Tag, <-chan error)
	TagResources(ctx context.Context, name string, res []godo.Resource) error
	UntagResources(ctx context.Context, name string, res []godo.Resource) error
}

type Tag interface {
	Struct() *godo.Tag
}

type Resource interface {
	Struct() *godo.Resource
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
		req: &godo.TagCreateRequest{},
	}
}

func (svc *client) defaultTagResourcesOpts() *tagResourcesOpt {
	return &tagResourcesOpt{
		req: &godo.TagResourcesRequest{},
	}
}

func (svc *client) defaultUntagResourcesOpts() *untagResourcesOpt {
	return &untagResourcesOpt{
		req: &godo.UntagResourcesRequest{},
	}
}

func UseGodoCreate(req *godo.TagCreateRequest) CreateOpt {
	return func(opt *createOpt) { opt.req = req }
}

type createOpt struct {
	req *godo.TagCreateRequest
}

type tagResourcesOpt struct {
	req *godo.TagResourcesRequest
}

type untagResourcesOpt struct {
	req *godo.UntagResourcesRequest
}

type tag struct {
	g *godo.Client
	t *godo.Tag
}

type resource struct {
	g *godo.Client
	r *godo.Resource
}

func (svc *tag) Struct() *godo.Tag { return svc.t }

func (svc *resource) Struct() *godo.Resource { return svc.r }

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

func (svc *client) TagResources(ctx context.Context, name string, res []godo.Resource) error {
	opt := svc.defaultTagResourcesOpts()
	opt.req.Resources = res

	_, err := svc.g.Tags.TagResources(ctx, name, opt.req)
	if err != nil {
		return err
	}

	return nil
}

func (svc *client) UntagResources(ctx context.Context, name string, res []godo.Resource) error {
	opt := svc.defaultUntagResourcesOpts()
	opt.req.Resources = res

	_, err := svc.g.Tags.UntagResources(ctx, name, opt.req)
	if err != nil {
		return err
	}

	return nil
}

func (svc *client) Get(ctx context.Context, name string) (Tag, error) {
	t, _, err := svc.g.Tags.Get(ctx, name)
	if err != nil {
		return nil, err
	}

	return &tag{g: svc.g, t: t}, nil
}
