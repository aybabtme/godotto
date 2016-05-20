package drives

import (
	"github.com/aybabtme/godotto/internal/godoutil"
	"github.com/digitalocean/godo"
	"golang.org/x/net/context"
)

// A Client can interact with the DigitalOcean Drives service.
type Client interface {
	CreateDrive(ctx context.Context, name, region string, sizeGibiBytes int64, opts ...CreateOpt) (Drive, error)
	GetDrive(context.Context, string) (Drive, error)
	DeleteDrive(context.Context, string) error
	ListDrives(context.Context) (<-chan Drive, <-chan error)

	CreateSnapshot(ctx context.Context, driveID, name string, opts ...SnapshotOpt) (Snapshot, error)
	GetSnapshot(context.Context, string) (Snapshot, error)
	DeleteSnapshot(context.Context, string) error
	ListSnapshots(ctx context.Context, driveID string) (<-chan Snapshot, <-chan error)

	Actions() ActionClient
}

// A Drive in the DigitalOcean cloud.
type Drive interface {
	Struct() *godo.Drive
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

// CreateOpt is an optional argument to Drives.Create.
type CreateOpt func(*createOpt)

// SetDriveDescription does what it says on the tin.
func SetDriveDescription(desc string) CreateOpt {
	return func(opt *createOpt) { opt.req.Description = desc }
}

type createOpt struct {
	req *godo.DriveCreateRequest
}

func (svc *client) defaultCreateOpts() *createOpt {
	return &createOpt{
		req: &godo.DriveCreateRequest{},
	}
}

func (svc *client) CreateDrive(ctx context.Context, name, region string, sizeGibiBytes int64, opts ...CreateOpt) (Drive, error) {

	opt := svc.defaultCreateOpts()
	opt.req.Name = name
	opt.req.Region = region
	opt.req.SizeGibiBytes = sizeGibiBytes

	for _, fn := range opts {
		fn(opt)
	}
	d, _, err := svc.g.Storage.CreateDrive(opt.req)
	if err != nil {
		return nil, err
	}
	return &drive{g: svc.g, d: d}, nil
}

func (svc *client) GetDrive(ctx context.Context, name string) (Drive, error) {
	d, _, err := svc.g.Storage.GetDrive(name)
	if err != nil {
		return nil, err
	}
	return &drive{g: svc.g, d: d}, nil
}

func (svc *client) DeleteDrive(ctx context.Context, name string) error {
	_, err := svc.g.Storage.DeleteDrive(name)
	return err
}

func (svc *client) ListDrives(ctx context.Context) (<-chan Drive, <-chan error) {
	outc := make(chan Drive, 1)
	errc := make(chan error, 1)

	go func() {
		defer close(outc)
		defer close(errc)
		err := godoutil.IterateList(ctx, func(opt *godo.ListOptions) (*godo.Response, error) {
			r, resp, err := svc.g.Storage.ListDrives(opt)
			for _, d := range r {
				dd := d // copy ranged over variable
				select {
				case outc <- &drive{g: svc.g, d: &dd}:
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

type drive struct {
	g *godo.Client
	d *godo.Drive
}

func (svc *drive) Struct() *godo.Drive { return svc.d }

// SnapshotOpt is an optional argument to Drives.Edit.
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

func (svc *client) CreateSnapshot(ctx context.Context, driveID, name string, opts ...SnapshotOpt) (Snapshot, error) {

	opt := svc.defaultSnapshotOpts()
	for _, fn := range opts {
		fn(opt)
	}
	opt.req.DriveID = driveID
	opt.req.Name = name
	d, _, err := svc.g.Storage.CreateSnapshot(opt.req)
	if err != nil {
		return nil, err
	}
	return &snapshot{g: svc.g, d: d}, nil
}

func (svc *client) GetSnapshot(ctx context.Context, id string) (Snapshot, error) {
	d, _, err := svc.g.Storage.GetSnapshot(id)
	if err != nil {
		return nil, err
	}
	return &snapshot{g: svc.g, d: d}, nil
}

func (svc *client) DeleteSnapshot(ctx context.Context, id string) error {
	_, err := svc.g.Storage.DeleteSnapshot(id)
	return err
}

func (svc *client) ListSnapshots(ctx context.Context, driveID string) (<-chan Snapshot, <-chan error) {
	outc := make(chan Snapshot, 1)
	errc := make(chan error, 1)

	go func() {
		defer close(outc)
		defer close(errc)
		err := godoutil.IterateList(ctx, func(opt *godo.ListOptions) (*godo.Response, error) {
			r, resp, err := svc.g.Storage.ListSnapshots(driveID, opt)
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
