package spycloud

import (
	"sync"

	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/domains"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/drives"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/droplets"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/floatingips"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/keys"
	"github.com/aybabtme/godotto/pkg/extra/do/mockcloud"
	"github.com/digitalocean/godo"
)

// A Spy lets you see what's been created by a client.
type Spy func(*client)

// Droplets lets you visit all the droplets that are
// still created by the spied upon client.
func Droplets(fn func(*godo.Droplet)) Spy {
	return func(c *client) {
		for _, v := range c.droplets {
			fn(v)
		}
	}
}

// Drives lets you visit all the drives that are
// still created by the spied upon client.
func Drives(fn func(*godo.Drive)) Spy {
	return func(c *client) {
		for _, v := range c.drives {
			fn(v)
		}
	}
}

// Snapshots lets you visit all the snapshots that are
// still created by the spied upon client.
func Snapshots(fn func(*godo.Snapshot)) Spy {
	return func(c *client) {
		for _, v := range c.snapshots {
			fn(v)
		}
	}
}

// Domains lets you visit all the domains that are
// still created by the spied upon client.
func Domains(fn func(*godo.Domain)) Spy {
	return func(c *client) {
		for _, v := range c.domains {
			fn(v)
		}
	}
}

// Records lets you visit all the records that are
// still created by the spied upon client.
func Records(fn func(*godo.DomainRecord)) Spy {
	return func(c *client) {
		for _, v := range c.records {
			fn(v)
		}
	}
}

// FloatingIPs lets you visit all the floatingips that are
// still created by the spied upon client.
func FloatingIPs(fn func(*godo.FloatingIP)) Spy {
	return func(c *client) {
		for _, v := range c.floatingips {
			fn(v)
		}
	}
}

// Keys lets you visit all the keys that are
// still created by the spied upon client.
func Keys(fn func(*godo.Key)) Spy {
	return func(c *client) {
		for _, v := range c.keys {
			fn(v)
		}
	}
}

// Client wraps a client with a spy, which allows looking at
// the resources that currently exist in the client.
func Client(cloud cloud.Client) (cloud.Client, func(...Spy)) {
	c, mock := newClient(cloud)
	return mock, func(opts ...Spy) {
		c.mu.Lock()
		defer c.mu.Unlock()
		for _, opt := range opts {
			opt(c)
		}
	}
}

type client struct {
	real cloud.Client

	mu          sync.Mutex
	droplets    map[int]*godo.Droplet
	drives      map[string]*godo.Drive
	snapshots   map[string]*godo.Snapshot
	domains     map[string]*godo.Domain
	records     map[int]*godo.DomainRecord
	floatingips map[string]*godo.FloatingIP
	keys        map[int]*godo.Key
}

func newClient(cloud cloud.Client) (*client, *mockcloud.Mock) {
	mock := mockcloud.Client(cloud)

	c := &client{
		real:        cloud,
		droplets:    make(map[int]*godo.Droplet),
		drives:      make(map[string]*godo.Drive),
		snapshots:   make(map[string]*godo.Snapshot),
		domains:     make(map[string]*godo.Domain),
		records:     make(map[int]*godo.DomainRecord),
		floatingips: make(map[string]*godo.FloatingIP),
		keys:        make(map[int]*godo.Key),
	}

	// capture all create/delete actions

	mock.MockDroplets.CreateFn = c.interceptDropletCreate
	mock.MockDroplets.DeleteFn = c.interceptDropletDelete
	mock.MockDrives.CreateDriveFn = c.interceptDriveCreate
	mock.MockDrives.DeleteDriveFn = c.interceptDriveDelete
	mock.MockDrives.CreateSnapshotFn = c.interceptSnapshotCreate
	mock.MockDrives.DeleteSnapshotFn = c.interceptSnapshotDelete
	mock.MockDomains.CreateFn = c.interceptDomainCreate
	mock.MockDomains.DeleteFn = c.interceptDomainDelete
	mock.MockDomains.CreateRecordFn = c.interceptDomainRecordCreate
	mock.MockDomains.DeleteRecordFn = c.interceptDomainRecordDelete
	mock.MockFloatingIPs.CreateFn = c.interceptFloatingIPCreate
	mock.MockFloatingIPs.DeleteFn = c.interceptFloatingIPDelete
	mock.MockKeys.CreateFn = c.interceptKeyCreate
	mock.MockKeys.DeleteByIDFn = c.interceptKeyDeleteByID
	mock.MockKeys.DeleteByFingerprintFn = c.interceptKeyDeleteByFingerprint

	return c, mock
}

func (client *client) interceptDropletCreate(name, region, size, image string, opts ...droplets.CreateOpt) (droplets.Droplet, error) {
	d, err := client.real.Droplets().Create(name, region, size, image, opts...)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		client.droplets[d.Struct().ID] = d.Struct()
	}
	return d, err
}

func (client *client) interceptDropletDelete(id int) error {
	err := client.real.Droplets().Delete(id)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		delete(client.droplets, id)
	}
	return err
}

func (client *client) interceptDriveCreate(name, region string, sizeGibiBytes int64, opts ...drives.CreateOpt) (drives.Drive, error) {
	d, err := client.real.Drives().CreateDrive(name, region, sizeGibiBytes, opts...)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		client.drives[d.Struct().ID] = d.Struct()
	}
	return d, err
}

func (client *client) interceptDriveDelete(id string) error {
	err := client.real.Drives().DeleteDrive(id)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		delete(client.drives, id)
	}
	return err
}

func (client *client) interceptSnapshotCreate(driveID, name string, opts ...drives.SnapshotOpt) (drives.Snapshot, error) {
	d, err := client.real.Drives().CreateSnapshot(driveID, name, opts...)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		client.snapshots[d.Struct().ID] = d.Struct()
	}
	return d, err
}

func (client *client) interceptSnapshotDelete(id string) error {
	err := client.real.Drives().DeleteSnapshot(id)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		delete(client.snapshots, id)
	}
	return err
}

func (client *client) interceptDomainCreate(name, ip string, opts ...domains.CreateOpt) (domains.Domain, error) {
	d, err := client.real.Domains().Create(name, ip, opts...)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		client.domains[d.Struct().Name] = d.Struct()
	}
	return d, err
}

func (client *client) interceptDomainDelete(id string) error {
	err := client.real.Domains().Delete(id)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		delete(client.domains, id)
	}
	return err
}

func (client *client) interceptDomainRecordCreate(id string, opts ...domains.RecordOpt) (domains.Record, error) {
	d, err := client.real.Domains().CreateRecord(id, opts...)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		client.records[d.Struct().ID] = d.Struct()
	}
	return d, err
}

func (client *client) interceptDomainRecordDelete(name string, id int) error {
	err := client.real.Domains().DeleteRecord(name, id)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		delete(client.records, id)
	}
	return err
}

func (client *client) interceptFloatingIPCreate(region string, opts ...floatingips.CreateOpt) (floatingips.FloatingIP, error) {
	d, err := client.real.FloatingIPs().Create(region, opts...)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		client.floatingips[d.Struct().IP] = d.Struct()
	}
	return d, err
}

func (client *client) interceptFloatingIPDelete(ip string) error {
	err := client.real.FloatingIPs().Delete(ip)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		delete(client.floatingips, ip)
	}
	return err
}

func (client *client) interceptKeyCreate(name, publicKey string, opts ...keys.CreateOpt) (keys.Key, error) {
	d, err := client.real.Keys().Create(name, publicKey, opts...)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		client.keys[d.Struct().ID] = d.Struct()
	}
	return d, err
}

func (client *client) interceptKeyDeleteByID(id int) error {
	err := client.real.Keys().DeleteByID(id)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		delete(client.keys, id)
	}
	return err
}

func (client *client) interceptKeyDeleteByFingerprint(fp string) error {
	err := client.real.Keys().DeleteByFingerprint(fp)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		var id int
		for _, key := range client.keys {
			if key.Fingerprint == fp {
				id = key.ID
				break
			}
		}
		delete(client.keys, id)
	}
	return err
}
