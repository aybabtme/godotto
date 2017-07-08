package spycloud

import (
	"context"
	"sync"

	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/domains"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/droplets"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/floatingips"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/keys"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/loadbalancers"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/tags"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/volumes"
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

// Volumes lets you visit all the volumes that are
// still created by the spied upon client.
func Volumes(fn func(*godo.Volume)) Spy {
	return func(c *client) {
		for _, v := range c.volumes {
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

// Tags lets you visit all the tags that are still created by the spied upon
// client.
func Tags(fn func(*godo.Tag)) Spy {
	return func(c *client) {
		for _, v := range c.tags {
			fn(v)
		}
	}
}

// Load Balancers lets you visit all the load balancers that are still created
// by the spied upon client.
func LoadBalancers(fn func(*godo.LoadBalancer)) Spy {
	return func(c *client) {
		for _, v := range c.loadbalancers {
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

	mu            sync.Mutex
	droplets      map[int]*godo.Droplet
	volumes       map[string]*godo.Volume
	snapshots     map[string]*godo.Snapshot
	domains       map[string]*godo.Domain
	records       map[int]*godo.DomainRecord
	floatingips   map[string]*godo.FloatingIP
	keys          map[int]*godo.Key
	tags          map[string]*godo.Tag
	loadbalancers map[string]*godo.LoadBalancer
}

func newClient(cloud cloud.Client) (*client, *mockcloud.Mock) {
	mock := mockcloud.Client(cloud)

	c := &client{
		real:          cloud,
		droplets:      make(map[int]*godo.Droplet),
		volumes:       make(map[string]*godo.Volume),
		snapshots:     make(map[string]*godo.Snapshot),
		domains:       make(map[string]*godo.Domain),
		records:       make(map[int]*godo.DomainRecord),
		floatingips:   make(map[string]*godo.FloatingIP),
		keys:          make(map[int]*godo.Key),
		tags:          make(map[string]*godo.Tag),
		loadbalancers: make(map[string]*godo.LoadBalancer),
	}

	// capture all create/delete actions

	mock.MockDroplets.CreateFn = c.interceptDropletCreate
	mock.MockDroplets.DeleteFn = c.interceptDropletDelete
	mock.MockVolumes.CreateVolumeFn = c.interceptVolumeCreate
	mock.MockVolumes.DeleteVolumeFn = c.interceptVolumeDelete
	mock.MockVolumes.CreateSnapshotFn = c.interceptVolumeSnapshotCreate
	mock.MockVolumes.DeleteSnapshotFn = c.interceptVolumeSnapshotDelete
	mock.MockDomains.CreateFn = c.interceptDomainCreate
	mock.MockDomains.DeleteFn = c.interceptDomainDelete
	mock.MockDomains.CreateRecordFn = c.interceptDomainRecordCreate
	mock.MockDomains.DeleteRecordFn = c.interceptDomainRecordDelete
	mock.MockFloatingIPs.CreateFn = c.interceptFloatingIPCreate
	mock.MockFloatingIPs.DeleteFn = c.interceptFloatingIPDelete
	mock.MockKeys.CreateFn = c.interceptKeyCreate
	mock.MockKeys.DeleteByIDFn = c.interceptKeyDeleteByID
	mock.MockKeys.DeleteByFingerprintFn = c.interceptKeyDeleteByFingerprint
	mock.MockTags.CreateFn = c.interceptTagCreate
	mock.MockTags.DeleteFn = c.interceptTagDelete
	mock.MockLoadBalancers.CreateFn = c.interceptLoadBalancerCreate
	mock.MockLoadBalancers.DeleteFn = c.interceptLoadBalancerDelete
	mock.MockSnapshots.DeleteFn = c.interceptSnapshotDelete
	return c, mock
}

func (client *client) interceptDropletCreate(ctx context.Context, name, region, size, image string, opts ...droplets.CreateOpt) (droplets.Droplet, error) {
	d, err := client.real.Droplets().Create(ctx, name, region, size, image, opts...)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		client.droplets[d.Struct().ID] = d.Struct()
	}
	return d, err
}

func (client *client) interceptDropletDelete(ctx context.Context, id int) error {
	err := client.real.Droplets().Delete(ctx, id)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		delete(client.droplets, id)
	}
	return err
}

func (client *client) interceptVolumeCreate(ctx context.Context, name, region string, sizeGibiBytes int64, opts ...volumes.CreateOpt) (volumes.Volume, error) {
	d, err := client.real.Volumes().CreateVolume(ctx, name, region, sizeGibiBytes, opts...)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		client.volumes[d.Struct().ID] = d.Struct()
	}
	return d, err
}

func (client *client) interceptVolumeDelete(ctx context.Context, id string) error {
	err := client.real.Volumes().DeleteVolume(ctx, id)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		delete(client.volumes, id)
	}
	return err
}

func (client *client) interceptVolumeSnapshotCreate(ctx context.Context, volumeID, name string, opts ...volumes.SnapshotOpt) (volumes.Snapshot, error) {
	d, err := client.real.Volumes().CreateSnapshot(ctx, volumeID, name, opts...)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		client.snapshots[d.Struct().ID] = d.Struct()
	}
	return d, err
}

func (client *client) interceptVolumeSnapshotDelete(ctx context.Context, id string) error {
	err := client.real.Volumes().DeleteSnapshot(ctx, id)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		delete(client.snapshots, id)
	}
	return err
}

func (client *client) interceptDomainCreate(ctx context.Context, name, ip string, opts ...domains.CreateOpt) (domains.Domain, error) {
	d, err := client.real.Domains().Create(ctx, name, ip, opts...)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		client.domains[d.Struct().Name] = d.Struct()
	}
	return d, err
}

func (client *client) interceptDomainDelete(ctx context.Context, id string) error {
	err := client.real.Domains().Delete(ctx, id)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		delete(client.domains, id)
	}
	return err
}

func (client *client) interceptDomainRecordCreate(ctx context.Context, id string, opts ...domains.RecordOpt) (domains.Record, error) {
	d, err := client.real.Domains().CreateRecord(ctx, id, opts...)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		client.records[d.Struct().ID] = d.Struct()
	}
	return d, err
}

func (client *client) interceptDomainRecordDelete(ctx context.Context, name string, id int) error {
	err := client.real.Domains().DeleteRecord(ctx, name, id)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		delete(client.records, id)
	}
	return err
}

func (client *client) interceptFloatingIPCreate(ctx context.Context, region string, opts ...floatingips.CreateOpt) (floatingips.FloatingIP, error) {
	d, err := client.real.FloatingIPs().Create(ctx, region, opts...)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		client.floatingips[d.Struct().IP] = d.Struct()
	}
	return d, err
}

func (client *client) interceptFloatingIPDelete(ctx context.Context, ip string) error {
	err := client.real.FloatingIPs().Delete(ctx, ip)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		delete(client.floatingips, ip)
	}
	return err
}

func (client *client) interceptKeyCreate(ctx context.Context, name, publicKey string, opts ...keys.CreateOpt) (keys.Key, error) {
	d, err := client.real.Keys().Create(ctx, name, publicKey, opts...)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		client.keys[d.Struct().ID] = d.Struct()
	}
	return d, err
}

func (client *client) interceptKeyDeleteByID(ctx context.Context, id int) error {
	err := client.real.Keys().DeleteByID(ctx, id)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		delete(client.keys, id)
	}
	return err
}

func (client *client) interceptKeyDeleteByFingerprint(ctx context.Context, fp string) error {
	err := client.real.Keys().DeleteByFingerprint(ctx, fp)
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

func (client *client) interceptTagCreate(ctx context.Context, name string, opt ...tags.CreateOpt) (tags.Tag, error) {
	t, err := client.real.Tags().Create(ctx, name, opt...)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		client.tags[t.Struct().Name] = t.Struct()
	}

	return t, err
}

func (client *client) interceptTagDelete(ctx context.Context, name string) error {
	err := client.real.Tags().Delete(ctx, name)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		delete(client.tags, name)
	}

	return err
}

func (client *client) interceptLoadBalancerCreate(ctx context.Context, name, region string, rules []godo.ForwardingRule, opts ...loadbalancers.CreateOpt) (loadbalancers.LoadBalancer, error) {
	l, err := client.real.LoadBalancers().Create(ctx, name, region, rules, opts...)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		client.loadbalancers[l.Struct().ID] = l.Struct()
	}

	return l, err
}

func (client *client) interceptLoadBalancerDelete(ctx context.Context, id string) error {
	err := client.real.LoadBalancers().Delete(ctx, id)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		delete(client.loadbalancers, id)
	}

	return err
}

func (client *client) interceptSnapshotDelete(ctx context.Context, sId string) error {
	err := client.real.Snapshots().Delete(ctx, sId)
	if err == nil {
		client.mu.Lock()
		defer client.mu.Unlock()
		delete(client.snapshots, sId)
	}

	return err
}
