package floatingips_test

import (
	"errors"
	"testing"
	"context"

	"github.com/aybabtme/godotto/pkg/extra/vmtest"
	"github.com/aybabtme/godotto/pkg/extra/do/mockcloud"
)

func TestActionsApply(t *testing.T) {
	cloud := mockcloud.Client(nil)
	vmtest.Run(t, cloud, `
var pkg = cloud.floating_ips.actions;

assert(pkg != null, "package should be loaded");
assert(pkg.assign != null, "assign function should be defined");
assert(pkg.unassign != null, "unassign function should be defined");
    `)
}

func TestActionsThrows(t *testing.T) {
	cloud := mockcloud.Client(nil)

	mock := cloud.MockFloatingIPs.MockFloatingIPActions
	mock.AssignFn = func(ctx context.Context, ip string, did int) error {
		return errors.New("throw me")
	}
	mock.UnassignFn = func(ctx context.Context, ip string) error {
		return errors.New("throw me")
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.floating_ips.actions;

[

	{ name: "assign",	fn: function() { pkg.assign("127.0.0.1", 42) } },
	{ name: "unassign",	fn: function() { pkg.unassign("127.0.0.1") } },

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

func TestActionAssign(t *testing.T) {
	cloud := mockcloud.Client(nil)
	mock := cloud.MockFloatingIPs.MockFloatingIPActions
	mock.AssignFn = func(ctx context.Context, ip string, did int) error {
		if want, got := "127.0.0.1", ip; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		if want, got := 42, did; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		return nil
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.floating_ips.actions;
pkg.assign("127.0.0.1", 42);
	`)

}

func TestActionUnassign(t *testing.T) {
	cloud := mockcloud.Client(nil)
	mock := cloud.MockFloatingIPs.MockFloatingIPActions

	mock.UnassignFn = func(ctx context.Context, ip string) error {
		if want, got := "127.0.0.1", ip; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		return nil
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.floating_ips.actions;
pkg.unassign("127.0.0.1");
	`)
}
