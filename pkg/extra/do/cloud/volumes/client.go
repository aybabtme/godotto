package volumes

import (
	"context"

	"github.com/aybabtme/godotto/internal/godoutil"
	"github.com/digitalocean/godo"
)

// A Client can interact with the DigitalOcean Volumes service.
type Client interface {
	CreateVolume(ctx context.Context, name, region string, sizeGibiBytes int64, opts ...CreateOpt) (Volume, error)
	GetVolume(context.Context, string) (Volume, error)
	DeleteVolume(context.Context, string) error
	ListVolumes(context.Context) (<-chan Volume, <-chan error)

	CreateSnapshot(ctx context.Context, volumeID, name string, opts ...SnapshotOpt) (Snapshot, error)
	GetSnapshot(context.Context, string) (Snapshot, error)
	DeleteSnapshot(context.Context, string) error
	ListSnapshots(ctx context.Context, volumeID string) (<-chan Snapshot, <-chan error)

	Actions() ActionClient
}

// A Volume in the DigitalOcean cloud.
type Volume interface {
	Struct() *godo.Volume
}

// A Snapshot in the DigitalOcean cloud.
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

// CreateOpt is an optional argument to Volumes.Create.
type CreateOpt func(*createOpt)

// SetVolumeDescription does what it says on the tin.
func SetVolumeDescription(desc string) CreateOpt {
	return func(opt *createOpt) { opt.req.Description = desc }
}

type createOpt struct {
	req *godo.VolumeCreateRequest
}

func (svc *client) defaultCreateOpts() *createOpt {
	return &createOpt{
		req: &godo.VolumeCreateRequest{},
	}
}

func (svc *client) CreateVolume(ctx context.Context, name, region string, sizeGibiBytes int64, opts ...CreateOpt) (Volume, error) {

	opt := svc.defaultCreateOpts()
	opt.req.Name = name
	opt.req.Region = region
	opt.req.SizeGigaBytes = sizeGibiBytes

	for _, fn := range opts {
		fn(opt)
	}
	d, _, err := svc.g.Storage.CreateVolume(ctx, opt.req)
	if err != nil {
		return nil, err
	}
	return &volume{g: svc.g, d: d}, nil
}

func (svc *client) GetVolume(ctx context.Context, name string) (Volume, error) {
	d, _, err := svc.g.Storage.GetVolume(ctx, name)
	if err != nil {
		return nil, err
	}
	return &volume{g: svc.g, d: d}, nil
}

func (svc *client) DeleteVolume(ctx context.Context, name string) error {
	_, err := svc.g.Storage.DeleteVolume(ctx, name)
	return err
}

func (svc *client) ListVolumes(ctx context.Context) (<-chan Volume, <-chan error) {
	outc := make(chan Volume, 1)
	errc := make(chan error, 1)

	go func() {
		defer close(outc)
		defer close(errc)
		err := godoutil.IterateList(ctx, func(ctx context.Context, opt *godo.ListOptions) (*godo.Response, error) {
			lvp := &godo.ListVolumeParams{ListOptions: opt}
			r, resp, err := svc.g.Storage.ListVolumes(ctx, lvp)
			for _, d := range r {
				dd := d // copy ranged over variable
				select {
				case outc <- &volume{g: svc.g, d: &dd}:
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

func (svc *client) Actions() ActionClient {
	return &actionClient{g: svc.g}
}

type volume struct {
	g *godo.Client
	d *godo.Volume
}

func (svc *volume) Struct() *godo.Volume { return svc.d }

// SnapshotOpt is an optional argument to Volumes.Edit.
type SnapshotOpt func(*snapshotOpt)

// SetSnapshotDescription does what it says on the tin.
func SetSnapshotDescription(desc string) SnapshotOpt {
	return func(opt *snapshotOpt) { opt.req.Description = desc }
}

type snapshotOpt struct {
	req *godo.SnapshotCreateRequest
}

func (svc *client) defaultSnapshotOpts() *snapshotOpt {
	return &snapshotOpt{
		req: &godo.SnapshotCreateRequest{},
	}
}

func (svc *client) CreateSnapshot(ctx context.Context, volumeID, name string, opts ...SnapshotOpt) (Snapshot, error) {

	opt := svc.defaultSnapshotOpts()
	for _, fn := range opts {
		fn(opt)
	}
	opt.req.VolumeID = volumeID
	opt.req.Name = name
	d, _, err := svc.g.Storage.CreateSnapshot(ctx, opt.req)
	if err != nil {
		return nil, err
	}
	return &snapshot{g: svc.g, d: d}, nil
}

func (svc *client) GetSnapshot(ctx context.Context, id string) (Snapshot, error) {
	d, _, err := svc.g.Storage.GetSnapshot(ctx, id)
	if err != nil {
		return nil, err
	}
	return &snapshot{g: svc.g, d: d}, nil
}

func (svc *client) DeleteSnapshot(ctx context.Context, id string) error {
	_, err := svc.g.Storage.DeleteSnapshot(ctx, id)
	return err
}

func (svc *client) ListSnapshots(ctx context.Context, volumeID string) (<-chan Snapshot, <-chan error) {
	outc := make(chan Snapshot, 1)
	errc := make(chan error, 1)

	go func() {
		defer close(outc)
		defer close(errc)
		err := godoutil.IterateList(ctx, func(ctx context.Context, opt *godo.ListOptions) (*godo.Response, error) {
			r, resp, err := svc.g.Storage.ListSnapshots(ctx, volumeID, opt)
			for _, d := range r {
				dd := d // copy ranged over variable
				select {
				case outc <- &snapshot{g: svc.g, d: &dd}:
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
	d *godo.Snapshot
}

func (svc *snapshot) Struct() *godo.Snapshot { return svc.d }
