package droplets_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aybabtme/godotto/pkg/extra/vmtest"
	"github.com/aybabtme/godotto/pkg/extra/do/mockcloud"
)

func TestDropletActionsApply(t *testing.T) {
	cloud := mockcloud.Client(nil)
	vmtest.Run(t, cloud, `
var pkg = cloud.droplets.actions;

assert(pkg != null, "package should be loaded");
assert(pkg.shutdown != null, "shutdown function should be defined");
assert(pkg.power_off != null, "power_off function should be defined");
assert(pkg.power_on != null, "power_on function should be defined");
assert(pkg.power_cycle != null, "power_cycle function should be defined");
assert(pkg.reboot != null, "reboot function should be defined");
assert(pkg.restore != null, "restore function should be defined");
assert(pkg.resize != null, "resize function should be defined");
assert(pkg.rename != null, "rename function should be defined");
assert(pkg.snapshot != null, "snapshot function should be defined");
assert(pkg.enable_backups != null, "enable_backups function should be defined");
assert(pkg.disable_backups != null, "disable_backups function should be defined");
assert(pkg.password_reset != null, "password_reset function should be defined");
assert(pkg.change_kernel != null, "change_kernel function should be defined");
assert(pkg.enable_ipv6 != null, "enable_ipv6 function should be defined");
assert(pkg.enable_private_networking != null, "enable_private_networking function should be defined");
assert(pkg.upgrade != null, "upgrade function should be defined");
    `)
}

func TestDropletActionsThrows(t *testing.T) {
	cloud := mockcloud.Client(nil)

	mock := cloud.MockDroplets.MockDropletActions
	mock.ShutdownFn = func(ctx context.Context, dropletID int) error {
		return errors.New("throw me")
	}
	mock.PowerOffFn = func(ctx context.Context, dropletID int) error {
		return errors.New("throw me")
	}
	mock.PowerOnFn = func(ctx context.Context, dropletID int) error {
		return errors.New("throw me")
	}
	mock.PowerCycleFn = func(ctx context.Context, dropletID int) error {
		return errors.New("throw me")
	}
	mock.RebootFn = func(ctx context.Context, dropletID int) error {
		return errors.New("throw me")
	}
	mock.RestoreFn = func(ctx context.Context, dropletID, imageID int) error {
		return errors.New("throw me")
	}
	mock.ResizeFn = func(ctx context.Context, dropletID int, sizeSlug string, resizeDisk bool) error {
		return errors.New("throw me")
	}
	mock.RenameFn = func(ctx context.Context, dropletID int, name string) error {
		return errors.New("throw me")
	}
	mock.SnapshotFn = func(ctx context.Context, dropletID int, name string) error {
		return errors.New("throw me")
	}
	mock.EnableBackupsFn = func(ctx context.Context, dropletID int) error {
		return errors.New("throw me")
	}
	mock.DisableBackupsFn = func(ctx context.Context, dropletID int) error {
		return errors.New("throw me")
	}
	mock.PasswordResetFn = func(ctx context.Context, dropletID int) error {
		return errors.New("throw me")
	}
	mock.RebuildByImageIDFn = func(ctx context.Context, dropletID int, imageID int) error {
		return errors.New("throw me")
	}
	mock.RebuildByImageSlugFn = func(ctx context.Context, dropletID int, imageSlug string) error {
		return errors.New("throw me")
	}
	mock.ChangeKernelFn = func(ctx context.Context, dropletID int, kernelID int) error {
		return errors.New("throw me")
	}
	mock.EnableIPv6Fn = func(ctx context.Context, dropletID int) error {
		return errors.New("throw me")
	}
	mock.EnablePrivateNetworkingFn = func(ctx context.Context, dropletID int) error {
		return errors.New("throw me")
	}
	mock.UpgradeFn = func(ctx context.Context, dropletID int) error {
		return errors.New("throw me")
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.droplets.actions;

[

	{ name: "shutdown",					  fn: function() { pkg.shutdown(42) } },
	{ name: "power_off",				  fn: function() { pkg.power_off(42) } },
	{ name: "power_on",					  fn: function() { pkg.power_on(42) } },
	{ name: "power_cycle",				  fn: function() { pkg.power_cycle(42) } },
	{ name: "reboot",					  fn: function() { pkg.reboot(42) } },
	{ name: "restore",					  fn: function() { pkg.restore(42, 43) } },
	{ name: "resize",					  fn: function() { pkg.resize(42, "4gb") } },
	{ name: "rename",					  fn: function() { pkg.rename(42, "lolzerg") } },
	{ name: "snapshot",					  fn: function() { pkg.snapshot(42, "lolol") } },
	{ name: "enable_backups",			  fn: function() { pkg.enable_backups(42) } },
	{ name: "disable_backups",			  fn: function() { pkg.disable_backups(42) } },
	{ name: "password_reset",			  fn: function() { pkg.password_reset(42) } },
	{ name: "change_kernel",			  fn: function() { pkg.change_kernel(42, 44) } },
	{ name: "enable_ipv6",				  fn: function() { pkg.enable_ipv6(42) } },
	{ name: "enable_private_networking",  fn: function() { pkg.enable_private_networking(42) } },
	{ name: "upgrade",					  fn: function() { pkg.upgrade(42) } },

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

func TestDropletActionShutdown(t *testing.T) {
	cloud := mockcloud.Client(nil)
	mock := cloud.MockDroplets.MockDropletActions
	mock.ShutdownFn = func(ctx context.Context, dropletID int) error {
		if want, got := 42, dropletID; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		return nil
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.droplets.actions;
pkg.shutdown(42);
	`)

}

func TestDropletActionPowerOff(t *testing.T) {
	cloud := mockcloud.Client(nil)
	mock := cloud.MockDroplets.MockDropletActions

	mock.PowerOffFn = func(ctx context.Context, dropletID int) error {
		if want, got := 42, dropletID; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		return nil
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.droplets.actions;
pkg.power_off(42);
	`)
}

func TestDropletActionPowerOn(t *testing.T) {
	cloud := mockcloud.Client(nil)
	mock := cloud.MockDroplets.MockDropletActions

	mock.PowerOnFn = func(ctx context.Context, dropletID int) error {
		if want, got := 42, dropletID; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		return nil
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.droplets.actions;
pkg.power_on(42);
	`)
}

func TestDropletActionPowerCycle(t *testing.T) {
	cloud := mockcloud.Client(nil)
	mock := cloud.MockDroplets.MockDropletActions

	mock.PowerCycleFn = func(ctx context.Context, dropletID int) error {
		if want, got := 42, dropletID; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		return nil
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.droplets.actions;
pkg.power_cycle(42);
	`)
}

func TestDropletActionReboot(t *testing.T) {
	cloud := mockcloud.Client(nil)
	mock := cloud.MockDroplets.MockDropletActions

	mock.RebootFn = func(ctx context.Context, dropletID int) error {
		if want, got := 42, dropletID; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		return nil
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.droplets.actions;
pkg.reboot(42);
	`)
}

func TestDropletActionRestore(t *testing.T) {
	cloud := mockcloud.Client(nil)
	mock := cloud.MockDroplets.MockDropletActions

	mock.RestoreFn = func(ctx context.Context, dropletID, imageID int) error {
		if want, got := 42, dropletID; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		if want, got := 43, imageID; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		return nil
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.droplets.actions;
pkg.restore(42, 43);
	`)
}

func TestDropletActionResize(t *testing.T) {
	cloud := mockcloud.Client(nil)
	mock := cloud.MockDroplets.MockDropletActions

	mock.ResizeFn = func(ctx context.Context, dropletID int, sizeSlug string, resizeDisk bool) error {
		if want, got := 42, dropletID; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		if want, got := "4gb", sizeSlug; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		if want, got := true, resizeDisk; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		return nil
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.droplets.actions;
pkg.resize(42, "4gb", true);
	`)
}

func TestDropletActionRename(t *testing.T) {
	cloud := mockcloud.Client(nil)
	mock := cloud.MockDroplets.MockDropletActions

	mock.RenameFn = func(ctx context.Context, dropletID int, name string) error {
		if want, got := 42, dropletID; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		if want, got := "hello", name; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		return nil
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.droplets.actions;
pkg.rename(42, "hello");
	`)
}

func TestDropletActionSnapshot(t *testing.T) {
	cloud := mockcloud.Client(nil)
	mock := cloud.MockDroplets.MockDropletActions

	mock.SnapshotFn = func(ctx context.Context, dropletID int, name string) error {
		if want, got := 42, dropletID; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		if want, got := "hello", name; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		return nil
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.droplets.actions;
pkg.snapshot(42, "hello");
	`)
}

func TestDropletActionEnableBackups(t *testing.T) {
	cloud := mockcloud.Client(nil)
	mock := cloud.MockDroplets.MockDropletActions

	mock.EnableBackupsFn = func(ctx context.Context, dropletID int) error {
		if want, got := 42, dropletID; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		return nil
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.droplets.actions;
pkg.enable_backups(42);
	`)
}

func TestDropletActionDisableBackups(t *testing.T) {
	cloud := mockcloud.Client(nil)
	mock := cloud.MockDroplets.MockDropletActions

	mock.DisableBackupsFn = func(ctx context.Context, dropletID int) error {
		if want, got := 42, dropletID; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		return nil
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.droplets.actions;
pkg.disable_backups(42);
	`)
}

func TestDropletActionPasswordReset(t *testing.T) {
	cloud := mockcloud.Client(nil)
	mock := cloud.MockDroplets.MockDropletActions

	mock.PasswordResetFn = func(ctx context.Context, dropletID int) error {
		if want, got := 42, dropletID; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		return nil
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.droplets.actions;
pkg.password_reset(42);
	`)
}

func TestDropletActionChangeKernel(t *testing.T) {
	cloud := mockcloud.Client(nil)
	mock := cloud.MockDroplets.MockDropletActions

	mock.ChangeKernelFn = func(ctx context.Context, dropletID int, kernelID int) error {
		if want, got := 42, dropletID; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		if want, got := 43, kernelID; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		return nil
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.droplets.actions;
pkg.change_kernel(42, 43);
	`)
}

func TestDropletActionEnableIPv6(t *testing.T) {
	cloud := mockcloud.Client(nil)
	mock := cloud.MockDroplets.MockDropletActions

	mock.EnableIPv6Fn = func(ctx context.Context, dropletID int) error {
		if want, got := 42, dropletID; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		return nil
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.droplets.actions;
pkg.enable_ipv6(42);
	`)
}

func TestDropletActionEnablePrivateNetworking(t *testing.T) {
	cloud := mockcloud.Client(nil)
	mock := cloud.MockDroplets.MockDropletActions

	mock.EnablePrivateNetworkingFn = func(ctx context.Context, dropletID int) error {
		if want, got := 42, dropletID; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		return nil
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.droplets.actions;
pkg.enable_private_networking(42);
	`)
}

func TestDropletActionUpgrade(t *testing.T) {
	cloud := mockcloud.Client(nil)
	mock := cloud.MockDroplets.MockDropletActions

	mock.UpgradeFn = func(ctx context.Context, dropletID int) error {
		if want, got := 42, dropletID; got != want {
			t.Fatalf("want %v got %v", want, got)
		}
		return nil
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.droplets.actions;
pkg.upgrade(42);
	`)
}
