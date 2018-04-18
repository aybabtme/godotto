package domains

import (
	"context"

	"github.com/aybabtme/godotto/pkg/extra/godoutil"
	"github.com/digitalocean/godo"
)

// A Client can interact with the DigitalOcean Domains service.
type Client interface {
	Create(ctx context.Context, name, ip string, opts ...CreateOpt) (Domain, error)
	Get(context.Context, string) (Domain, error)
	Delete(context.Context, string) error
	List(context.Context) (<-chan Domain, <-chan error)

	CreateRecord(context.Context, string, ...RecordOpt) (Record, error)
	GetRecord(context.Context, string, int) (Record, error)
	UpdateRecord(context.Context, string, int, ...RecordOpt) (Record, error)
	DeleteRecord(context.Context, string, int) error
	ListRecord(ctx context.Context, name string) (<-chan Record, <-chan error)
}

// A Domain in the DigitalOcean cloud.
type Domain interface {
	Struct() *godo.Domain
}

// A Record in the DigitalOcean cloud.
type Record interface {
	Struct() *godo.DomainRecord
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
	req *godo.DomainCreateRequest
}

func (svc *client) defaultCreateOpts() *createOpt {
	return &createOpt{
		req: &godo.DomainCreateRequest{},
	}
}

func (svc *client) Create(ctx context.Context, name, ip string, opts ...CreateOpt) (Domain, error) {

	opt := svc.defaultCreateOpts()
	opt.req.Name = name

	for _, fn := range opts {
		fn(opt)
	}
	d, _, err := svc.g.Domains.Create(ctx, opt.req)
	if err != nil {
		return nil, err
	}
	return &domain{g: svc.g, d: d}, nil
}

func (svc *client) Get(ctx context.Context, name string) (Domain, error) {
	d, _, err := svc.g.Domains.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return &domain{g: svc.g, d: d}, nil
}

func (svc *client) Delete(ctx context.Context, name string) error {
	_, err := svc.g.Domains.Delete(ctx, name)
	return err
}

func (svc *client) List(ctx context.Context) (<-chan Domain, <-chan error) {
	outc := make(chan Domain, 1)
	errc := make(chan error, 1)

	go func() {
		defer close(outc)
		defer close(errc)
		err := godoutil.IterateList(ctx, func(ctx context.Context, opt *godo.ListOptions) (*godo.Response, error) {
			r, resp, err := svc.g.Domains.List(ctx, opt)
			for _, d := range r {
				dd := d // copy ranged over variable
				select {
				case outc <- &domain{g: svc.g, d: &dd}:
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

type domain struct {
	g *godo.Client
	d *godo.Domain
}

func (svc *domain) Struct() *godo.Domain { return svc.d }

// RecordOpt is an optional argument to domains.Edit.
type RecordOpt func(*recordOpt)

func UseGodoRecord(req *godo.DomainRecordEditRequest) RecordOpt {
	return func(opt *recordOpt) { opt.req = req }
}

type recordOpt struct {
	req *godo.DomainRecordEditRequest
}

func (svc *client) defaultRecordOpts() *recordOpt {
	return &recordOpt{
		req: &godo.DomainRecordEditRequest{},
	}
}

func (svc *client) CreateRecord(ctx context.Context, name string, opts ...RecordOpt) (Record, error) {

	opt := svc.defaultRecordOpts()
	for _, fn := range opts {
		fn(opt)
	}
	d, _, err := svc.g.Domains.CreateRecord(ctx, name, opt.req)
	if err != nil {
		return nil, err
	}
	return &record{g: svc.g, d: d}, nil
}

func (svc *client) GetRecord(ctx context.Context, name string, id int) (Record, error) {
	d, _, err := svc.g.Domains.Record(ctx, name, id)
	if err != nil {
		return nil, err
	}
	return &record{g: svc.g, d: d}, nil
}

func (svc *client) UpdateRecord(ctx context.Context, name string, id int, opts ...RecordOpt) (Record, error) {
	opt := svc.defaultRecordOpts()
	for _, fn := range opts {
		fn(opt)
	}
	d, _, err := svc.g.Domains.EditRecord(ctx, name, id, opt.req)
	if err != nil {
		return nil, err
	}
	return &record{g: svc.g, d: d}, nil
}

func (svc *client) DeleteRecord(ctx context.Context, name string, id int) error {
	_, err := svc.g.Domains.DeleteRecord(ctx, name, id)
	return err
}

func (svc *client) ListRecord(ctx context.Context, name string) (<-chan Record, <-chan error) {
	outc := make(chan Record, 1)
	errc := make(chan error, 1)

	go func() {
		defer close(outc)
		defer close(errc)
		err := godoutil.IterateList(ctx, func(ctx context.Context, opt *godo.ListOptions) (*godo.Response, error) {
			r, resp, err := svc.g.Domains.Records(ctx, name, opt)
			for _, d := range r {
				dd := d // copy ranged over variable
				select {
				case outc <- &record{g: svc.g, d: &dd}:
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

type record struct {
	g *godo.Client
	d *godo.DomainRecord
}

func (svc *record) Struct() *godo.DomainRecord { return svc.d }
