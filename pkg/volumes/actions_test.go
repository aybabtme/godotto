package volumes_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aybabtme/godotto/internal/vmtest"
	"github.com/aybabtme/godotto/pkg/extra/do/mockcloud"
)

func TestActionsApply(t *testing.T) {
	cloud := mockcloud.Client(nil)
	vmtest.Run(t, cloud, `
var pkg = cloud.volumes.actions;

assert(pkg != null, "package should be loaded");
assert(pkg.attach != null, "attach function should be defined");
assert(pkg.detach_by_droplet_id != null, "detach function should be defined");
    `)
}

func TestActionsThrows(t *testing.T) {
	cloud := mockcloud.Client(nil)

	mock := cloud.MockVolumes.MockVolumeActions
	mock.AttachFn = func(ctx context.Context, ip string, did int) error {
		return errors.New("throw me")
	}
	mock.DetachByDropletIDFn = func(ctx context.Context, ip string, dropletID int) error {
		return errors.New("throw me")
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.volumes.actions;

[

	{ name: "attach",	fn: function() { pkg.attach("127.0.0.1", 42) } },
	{ name: "detach_by_droplet_id",	fn: function() { pkg.detach_by_droplet_id("127.0.0.1", 42) } },

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

func TestActionattach(t *testing.T) {
	cloud := mockcloud.Client(nil)
	mock := cloud.MockVolumes.MockVolumeActions
	mock.AttachFn = func(ctx context.Context, ip string, did int) error {
		if want, got := "127.0.0.1", ip; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		if want, got := 42, did; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		return nil
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.volumes.actions;
pkg.attach("127.0.0.1", 42);
	`)

}

func TestActiondetachByDropletID(t *testing.T) {
	cloud := mockcloud.Client(nil)
	mock := cloud.MockVolumes.MockVolumeActions

	mock.DetachByDropletIDFn = func(ctx context.Context, ip string, dropletID int) error {
		if want, got := "127.0.0.1", ip; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		if want, got := 42, dropletID; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		return nil
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.volumes.actions;
pkg.detach_by_droplet_id("127.0.0.1", 42);
	`)
}
