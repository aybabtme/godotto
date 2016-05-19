package droplets_test

import (
	"errors"
	"testing"

	"golang.org/x/net/context"

	"github.com/aybabtme/godotto/internal/vmtest"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/droplets"
	"github.com/aybabtme/godotto/pkg/extra/do/mockcloud"
	"github.com/digitalocean/godo"
)

func TestApply(t *testing.T) {
	cloud := mockcloud.Client(nil)
	vmtest.Run(t, cloud, `
var pkg = cloud.droplets;

assert(pkg != null, "package should be loaded");
assert(pkg.list != null, "list function should be defined");
assert(pkg.get != null, "get function should be defined");
assert(pkg.create != null, "create function should be defined");
assert(pkg.delete != null, "delete function should be defined");
    `)
}

func TestThrows(t *testing.T) {
	cloud := mockcloud.Client(nil)

	cloud.MockDroplets.ListFn = func(_ context.Context) (<-chan droplets.Droplet, <-chan error) {
		lc := make(chan droplets.Droplet)
		close(lc)
		ec := make(chan error, 1)
		ec <- errors.New("throw me")
		close(ec)
		return lc, ec
	}
	cloud.MockDroplets.GetFn = func(_ context.Context, _ int) (droplets.Droplet, error) {
		return nil, errors.New("throw me")
	}
	cloud.MockDroplets.CreateFn = func(_ context.Context, _, _, _, _ string, _ ...droplets.CreateOpt) (droplets.Droplet, error) {
		return nil, errors.New("throw me")
	}
	cloud.MockDroplets.DeleteFn = func(_ context.Context, _ int) error {
		return errors.New("throw me")
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.droplets;

[
	{ name: "list",          fn: function() { pkg.list() } },
	{ name: "get",           fn: function() { pkg.get(42) } },
	{ name: "create",        fn: function() { pkg.create({"image":{}}) } },
	{ name: "delete",        fn: function() { pkg.delete({}) } },

].forEach(function(kv) {
	var name = kv.name;
	var fn = kv.fn;
	try {
		fn(); throw "dont catch me";
	} catch (e) {
		equals("throw me", e.message, name +" should send the right exception");
	};
})`)
}

type droplet struct {
	*godo.Droplet
}

func (k *droplet) Struct() *godo.Droplet { return k.Droplet }

func TestList(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockDroplets.ListFn = func(_ context.Context) (<-chan droplets.Droplet, <-chan error) {
		lc := make(chan droplets.Droplet, 1)
		lc <- &droplet{&godo.Droplet{
			ID:          42,
			Name:        "my_name",
			Memory:      20,
			Vcpus:       21,
			Disk:        22,
			Region:      &godo.Region{Slug: "nyc1"},
			Image:       &godo.Image{ID: 43, Slug: "coreos-stable"},
			Size:        &godo.Size{Slug: "4gb"},
			SnapshotIDs: []int{42},
			BackupIDs:   []int{43},
			Networks: &godo.Networks{
				V4: []godo.NetworkV4{
					{IPAddress: "127.0.0.1", Type: "public"},
				},
			},
			Status: "loling",
		}}
		close(lc)
		ec := make(chan error)
		close(ec)
		return lc, ec
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.droplets;

var list = pkg.list();
assert(list != null, "should have received a list");
assert(list.length > 0, "should have received some elements")

var d = list[0];
var want = {
	id:           42,
	name:         "my_name",
	memory:       20,
	vcpus:        21,
	disk:         22,
	region_slug:  "nyc1",
	image_id:     43,
	image_slug:   "coreos-stable",
	size_slug:    "4gb"
	snapshot_ids: [42],
	backup_ids:   [43],
	locked:       false,
	public_ipv4:  "127.0.0.1",
	status:       "loling"
};
equals(d, want, "should have proper object");
`)
}

func TestGet(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockDroplets.GetFn = func(_ context.Context, id int) (droplets.Droplet, error) {
		return &droplet{&godo.Droplet{
			ID:          42,
			Name:        "my_name",
			Memory:      20,
			Vcpus:       21,
			Disk:        22,
			Region:      &godo.Region{Slug: "nyc1"},
			Image:       &godo.Image{ID: 43, Slug: "coreos-stable"},
			Size:        &godo.Size{Slug: "4gb"},
			SnapshotIDs: []int{42},
			BackupIDs:   []int{43},
			Networks: &godo.Networks{
				V4: []godo.NetworkV4{
					{IPAddress: "127.0.0.1", Type: "public"},
				},
			},
			Status: "loling",
		}}, nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.droplets;

var d = pkg.get(42)
var want = {
	id:           42,
	name:         "my_name",
	memory:       20,
	vcpus:        21,
	disk:         22,
	region_slug:  "nyc1",
	image_id:     43,
	image_slug:   "coreos-stable",
	size_slug:    "4gb"
	snapshot_ids: [42],
	backup_ids:   [43],
	locked:       false,
	public_ipv4:  "127.0.0.1",
	status:       "loling"
};
equals(d, want, "should have proper object");
`)
}

func TestCreate(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockDroplets.CreateFn = func(_ context.Context, name, region, size, image string, _ ...droplets.CreateOpt) (droplets.Droplet, error) {
		return &droplet{&godo.Droplet{
			ID:          42,
			Name:        "my_name",
			Memory:      20,
			Vcpus:       21,
			Disk:        22,
			Region:      &godo.Region{Slug: "nyc1"},
			Image:       &godo.Image{ID: 43, Slug: "coreos-stable"},
			Size:        &godo.Size{Slug: "4gb"},
			SnapshotIDs: []int{42},
			BackupIDs:   []int{43},
			Networks: &godo.Networks{
				V4: []godo.NetworkV4{
					{IPAddress: "127.0.0.1", Type: "public"},
				},
			},
			Status: "loling",
		}}, nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.droplets;


var d = pkg.create({
	name:         "my_name",
	memory:       20,
	vcpus:        21,
	disk:         22,
	region:       "nyc1",
	image:        {slug:"coreos-stable"},
	size:         "4gb",
	snapshot_ids: [42],
	backups:      true,
	ipv6:         true,
	private_networking: true,
	user_data:    "lolll"
});

var want = {
	id:           42,
	name:         "my_name",
	memory:       20,
	vcpus:        21,
	disk:         22,
	region_slug:  "nyc1",
	image_id:     43,
	image_slug:   "coreos-stable",
	size_slug:    "4gb"
	snapshot_ids: [42],
	backup_ids:   [43],
	locked:       false,
	public_ipv4:  "127.0.0.1",
	status:       "loling"
};
equals(d, want, "should have proper object");
`)
}

func TestDelete(t *testing.T) {
	wantID := 42
	cloud := mockcloud.Client(nil)
	cloud.MockDroplets.DeleteFn = func(_ context.Context, gotID int) error {
		if gotID != wantID {
			t.Fatalf("want %v got %v", wantID, gotID)
		}
		return nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.droplets;

pkg.delete(42);
`)
}
