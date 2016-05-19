package floatingips_test

import (
	"errors"
	"testing"

	"golang.org/x/net/context"

	"github.com/aybabtme/godotto/internal/vmtest"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/floatingips"
	"github.com/aybabtme/godotto/pkg/extra/do/mockcloud"
	"github.com/digitalocean/godo"
)

func TestApply(t *testing.T) {
	cloud := mockcloud.Client(nil)
	vmtest.Run(t, cloud, `
var pkg = cloud.floating_ips;

assert(pkg != null, "package should be loaded");
assert(pkg.list != null, "list function should be defined");
assert(pkg.get != null, "get function should be defined");
assert(pkg.create != null, "create function should be defined");
assert(pkg.delete != null, "delete function should be defined");
    `)
}

func TestThrows(t *testing.T) {
	cloud := mockcloud.Client(nil)

	cloud.MockFloatingIPs.ListFn = func(_ context.Context) (<-chan floatingips.FloatingIP, <-chan error) {
		lc := make(chan floatingips.FloatingIP)
		close(lc)
		ec := make(chan error, 1)
		ec <- errors.New("throw me")
		close(ec)
		return lc, ec
	}
	cloud.MockFloatingIPs.GetFn = func(_ context.Context, _ string) (floatingips.FloatingIP, error) {
		return nil, errors.New("throw me")
	}
	cloud.MockFloatingIPs.CreateFn = func(_ context.Context, _ string, _ ...floatingips.CreateOpt) (floatingips.FloatingIP, error) {
		return nil, errors.New("throw me")
	}
	cloud.MockFloatingIPs.DeleteFn = func(_ context.Context, _ string) error {
		return errors.New("throw me")
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.floating_ips;

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

type floatingip struct {
	*godo.FloatingIP
}

func (k *floatingip) Struct() *godo.FloatingIP { return k.FloatingIP }

func TestList(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockFloatingIPs.ListFn = func(_ context.Context) (<-chan floatingips.FloatingIP, <-chan error) {
		lc := make(chan floatingips.FloatingIP, 1)
		lc <- &floatingip{&godo.FloatingIP{
			Region:  &godo.Region{Slug: "nyc3"},
			Droplet: &godo.Droplet{ID: 42},
			IP:      "127.0.0.1",
		}}
		close(lc)
		ec := make(chan error)
		close(ec)
		return lc, ec
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.floating_ips;

var list = pkg.list();
assert(list != null, "should have received a list");
assert(list.length > 0, "should have received some elements")

var d = list[0];
var want = {
	region_slug: "nyc3",
	droplet_id:  42,
	ip:          "127.0.0.1"
};
equals(d, want, "should have proper object");
`)
}

func TestGet(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockFloatingIPs.GetFn = func(_ context.Context, ip string) (floatingips.FloatingIP, error) {
		return &floatingip{&godo.FloatingIP{
			Region:  &godo.Region{Slug: "nyc3"},
			Droplet: &godo.Droplet{ID: 42},
			IP:      ip,
		}}, nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.floating_ips;

var d = pkg.get("127.0.0.1")
var want = {
	region_slug: "nyc3",
	droplet_id:  42,
	ip:          "127.0.0.1"
};
equals(d, want, "should have proper object");
`)
}

func TestCreate(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockFloatingIPs.CreateFn = func(_ context.Context, region string, _ ...floatingips.CreateOpt) (floatingips.FloatingIP, error) {
		return &floatingip{&godo.FloatingIP{
			Region:  &godo.Region{Slug: region},
			Droplet: &godo.Droplet{ID: 42},
			IP:      "127.0.0.1",
		}}, nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.floating_ips;


var d = pkg.create({region:"nyc3", droplet_id: 42});

var want = {
	region_slug: "nyc3",
	droplet_id:  42,
	ip:          "127.0.0.1"
};
equals(d, want, "should have proper object");
`)
}

func TestDelete(t *testing.T) {
	wantIP := "127.0.0.1"
	cloud := mockcloud.Client(nil)
	cloud.MockFloatingIPs.DeleteFn = func(_ context.Context, gotIP string) error {
		if gotIP != wantIP {
			t.Fatalf("want %v got %v", wantIP, gotIP)
		}
		return nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.floating_ips;

pkg.delete("127.0.0.1");
`)
}
