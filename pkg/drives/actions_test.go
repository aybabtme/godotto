package drives_test

import (
	"errors"
	"testing"

	"golang.org/x/net/context"

	"github.com/aybabtme/godotto/internal/vmtest"
	"github.com/aybabtme/godotto/pkg/extra/do/mockcloud"
)

func TestActionsApply(t *testing.T) {
	cloud := mockcloud.Client(nil)
	vmtest.Run(t, cloud, `
var pkg = cloud.drives.actions;

assert(pkg != null, "package should be loaded");
assert(pkg.attach != null, "attach function should be defined");
assert(pkg.detach != null, "detach function should be defined");
    `)
}

func TestActionsThrows(t *testing.T) {
	cloud := mockcloud.Client(nil)

	mock := cloud.MockDrives.MockDriveActions
	mock.AttachFn = func(ctx context.Context, ip string, did int) error {
		return errors.New("throw me")
	}
	mock.DetachFn = func(ctx context.Context, ip string) error {
		return errors.New("throw me")
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.drives.actions;

[

	{ name: "attach",	fn: function() { pkg.attach("127.0.0.1", 42) } },
	{ name: "detach",	fn: function() { pkg.detach("127.0.0.1") } },

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
	mock := cloud.MockDrives.MockDriveActions
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
var pkg = cloud.drives.actions;
pkg.attach("127.0.0.1", 42);
	`)

}

func TestActiondetach(t *testing.T) {
	cloud := mockcloud.Client(nil)
	mock := cloud.MockDrives.MockDriveActions

	mock.DetachFn = func(ctx context.Context, ip string) error {
		if want, got := "127.0.0.1", ip; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		return nil
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.drives.actions;
pkg.detach("127.0.0.1");
	`)
}
