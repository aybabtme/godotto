package droplets_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aybabtme/godotto/internal/vmtest"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/droplets"
	"github.com/aybabtme/godotto/pkg/extra/do/mockcloud"
	"github.com/digitalocean/godo"
)

func TestDropletApply(t *testing.T) {
	cloud := mockcloud.Client(nil)
	vmtest.Run(t, cloud, `
var pkg = cloud.droplets;

assert(pkg != null, "package should be loaded");
assert(pkg.list != null, "list function should be defined");
assert(pkg.get != null, "get function should be defined");
assert(pkg.create != null, "create function should be defined");
assert(pkg.create_multiple != null, "create_multiple function should be defined");
assert(pkg.delete != null, "delete function should be defined");
    `)
}

func TestDropletThrows(t *testing.T) {
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

var region = { name: "", slug: "", sizes: [], available: true, features: [] };
var image = { id: 42, name: "", type: "", distribution: "", slug: "", public: true, regions: [], min_disk_size: 0, };
var size = { slug: "", memory: 1, vcpus: 2, disk: 2, price_monthly: 1.0, price_hourly: 0.1, regions: [], available: true, transfer: 1.0, };

var d = {
  backup_ids: [ 43 ],
  disk: 22,
  id: 42,
  image: image,
  locked: false,
  memory: 20,
  name: "my_name",
  public_ipv4: "127.0.0.1",
  public_ipv6: "::1",
	private_ipv4: "",
  networks: {"v4":[], "v6":[]},
  region: region,
  size: size,
  snapshot_ids: [ 42 ],
  status: "loling",
  vcpus: 21
};
[
	{ name: "list",          fn: function() { pkg.list() } },
	{ name: "get",           fn: function() { pkg.get(42) } },
	{ name: "create",        fn: function() { pkg.create(d) } },
	{ name: "delete",        fn: function() { pkg.delete(42) } },

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

var (
	region = &godo.Region{Name: "newyork3", Slug: "nyc3", Sizes: []string{"small"}, Available: true, Features: []string{"all"}}
	size   = &godo.Size{Slug: "lol", Memory: 1, Vcpus: 2, Disk: 2, PriceMonthly: 1.0, PriceHourly: 0.1, Regions: []string{"lol"}, Available: true, Transfer: 1.0}
	image  = &godo.Image{ID: 42, Name: "derp", Type: "herp", Distribution: "coreos", Slug: "coreos-stable", Public: true, Regions: []string{"atlantis"}}
	d      = &godo.Droplet{ID: 42, Name: "my_name", Memory: 20, Vcpus: 21, Disk: 22, Region: region, Image: image, Size: size, SnapshotIDs: []int{42}, BackupIDs: []int{43}, Networks: &godo.Networks{}, Status: "loling", Tags: []string{"test"}, VolumeIDs: []string{""}}
)

type droplet struct {
	*godo.Droplet
}

func (k *droplet) Struct() *godo.Droplet { return k.Droplet }

func TestDropletList(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockDroplets.ListFn = func(_ context.Context) (<-chan droplets.Droplet, <-chan error) {
		lc := make(chan droplets.Droplet, 1)
		lc <- &droplet{d}
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

var region = { name: "newyork3", slug: "nyc3", sizes: ["small"], available: true, features: ["all"] };
var image = { id: 42, name: "derp", type: "herp", distribution: "coreos", slug: "coreos-stable", public: true, regions: ["atlantis"], min_disk_size: 0, };
var size = { slug: "lol", memory: 1, vcpus: 2, disk: 2, price_monthly: 1.0, price_hourly: 0.1, regions: ["lol"], available: true, transfer: 1.0, };

var want = {
  backup_ids: [ 43 ],
  disk: 22,
  id: 42,
  image: image,
  locked: false,
  memory: 20,
  name: "my_name",
  region: region,
  size: size,
  snapshot_ids: [ 42 ],
  status: "loling",
  networks: {"v4":{}, "v6":{}},
  tags: ["test"],
  created_at: "",
  volumes: [""],
  size_slug: "",
  public_ipv4: "",
  public_ipv6: "",
	private_ipv4: "",
  vcpus: 21
};

var d = list[0];
equals(d, want, "should have proper object");
`)
}

func TestDropletGet(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockDroplets.GetFn = func(_ context.Context, id int) (droplets.Droplet, error) {
		return &droplet{d}, nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.droplets;

var region = { name: "newyork3", slug: "nyc3", sizes: ["small"], available: true, features: ["all"] };
var image = { id: 42, name: "derp", type: "herp", distribution: "coreos", slug: "coreos-stable", public: true, regions: ["atlantis"], min_disk_size: 0, };
var size = { slug: "lol", memory: 1, vcpus: 2, disk: 2, price_monthly: 1.0, price_hourly: 0.1, regions: ["lol"], available: true, transfer: 1.0, };

var want = {
  backup_ids: [ 43 ],
  disk: 22,
  id: 42,
  image: image,
  locked: false,
  memory: 20,
  name: "my_name",
  region: region,
  size: size,
  snapshot_ids: [ 42 ],
  status: "loling",
  vcpus: 21,
  networks: {"v4":{}, "v6":{}},
  tags: ["test"],
  created_at: "",
  volumes: [""],
  size_slug: "",
  public_ipv4: "",
  public_ipv6: "",
	private_ipv4: "",
};

var d = pkg.get(42)
equals(d, want, "should have proper object");
`)
}

func TestDropletCreate(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockDroplets.CreateFn = func(_ context.Context, name, region, size, image string, _ ...droplets.CreateOpt) (droplets.Droplet, error) {
		return &droplet{d}, nil
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
	user_data:    "lolll",
	monitoring:   true,
	tags:         ["test"],
});

var region = { name: "newyork3", slug: "nyc3", sizes: ["small"], available: true, features: ["all"] };
var image = { id: 42, name: "derp", type: "herp", distribution: "coreos", slug: "coreos-stable", public: true, regions: ["atlantis"], min_disk_size: 0, };
var size = { slug: "lol", memory: 1, vcpus: 2, disk: 2, price_monthly: 1.0, price_hourly: 0.1, regions: ["lol"], available: true, transfer: 1.0, };

var want = {
  backup_ids: [ 43 ],
  disk: 22,
  id: 42,
  image: image,
  locked: false,
  memory: 20,
  name: "my_name",
  region: region,
  size: size,
  snapshot_ids: [ 42 ],
  status: "loling",
  vcpus: 21,
  networks: {"v4":{}, "v6":{}},
  tags: ["test"],
  created_at: "",
  volumes: [""],
  size_slug: "",
  public_ipv4: "",
  public_ipv6: "",
	private_ipv4: ""
};

equals(d, want, "should have proper object");
`)
}

func TestDropletDelete(t *testing.T) {
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
