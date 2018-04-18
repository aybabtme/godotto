package snapshots

import (
	"context"

	"github.com/aybabtme/godotto/pkg/extra/godoutil"
	"github.com/digitalocean/godo"
)

type Client interface {
	Get(ctx context.Context, id string) (Snapshot, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) (<-chan Snapshot, <-chan error)
	ListDroplet(ctx context.Context) (<-chan Snapshot, <-chan error)
	ListVolume(ctx context.Context) (<-chan Snapshot, <-chan error)
}

type Snapshot interface {
	Struct() *godo.Snapshot
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

func (svc *client) Get(ctx context.Context, id string) (Snapshot, error) {
	s, _, err := svc.g.Snapshots.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &snapshot{g: svc.g, s: s}, nil
}

func (svc *client) Delete(ctx context.Context, id string) error {
	_, err := svc.g.Snapshots.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (svc *client) List(ctx context.Context) (<-chan Snapshot, <-chan error) {
	outc := make(chan Snapshot, 1)
	errc := make(chan error, 1)

	go func() {
		defer close(outc)
		defer close(errc)
		err := godoutil.IterateList(ctx, func(ctx context.Context, opt *godo.ListOptions) (*godo.Response, error) {
			r, resp, err := svc.g.Snapshots.List(ctx, opt)
			for _, s := range r {
				ss := s // copy ranged over variable
				select {
				case outc <- &snapshot{g: svc.g, s: &ss}:
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

func (svc *client) ListDroplet(ctx context.Context) (<-chan Snapshot, <-chan error) {
	outc := make(chan Snapshot, 1)
	errc := make(chan error, 1)

	go func() {
		defer close(outc)
		defer close(errc)
		err := godoutil.IterateList(ctx, func(ctx context.Context, opt *godo.ListOptions) (*godo.Response, error) {
			r, resp, err := svc.g.Snapshots.ListDroplet(ctx, opt)
			for _, s := range r {
				ss := s // copy ranged over variable
				select {
				case outc <- &snapshot{g: svc.g, s: &ss}:
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

func (svc *client) ListVolume(ctx context.Context) (<-chan Snapshot, <-chan error) {
	outc := make(chan Snapshot, 1)
	errc := make(chan error, 1)

	go func() {
		defer close(outc)
		defer close(errc)
		err := godoutil.IterateList(ctx, func(ctx context.Context, opt *godo.ListOptions) (*godo.Response, error) {
			r, resp, err := svc.g.Snapshots.ListVolume(ctx, opt)
			for _, s := range r {
				ss := s // copy ranged over variable
				select {
				case outc <- &snapshot{g: svc.g, s: &ss}:
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

type snapshot struct {
	g *godo.Client
	s *godo.Snapshot
}

func (svc *snapshot) Struct() *godo.Snapshot { return svc.s }
