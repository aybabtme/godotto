package drives_test

import (
	"errors"
	"testing"

	"golang.org/x/net/context"

	"github.com/aybabtme/godotto/internal/vmtest"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/drives"
	"github.com/aybabtme/godotto/pkg/extra/do/mockcloud"
	"github.com/digitalocean/godo"
)

func TestApply(t *testing.T) {
	cloud := mockcloud.Client(nil)
	vmtest.Run(t, cloud, `
var pkg = cloud.drives;

assert(pkg != null, "package should be loaded");
assert(pkg.list_drives != null, "list_drives function should be defined");
assert(pkg.get_drive != null, "get_drive function should be defined");
assert(pkg.create_drive != null, "create_drive function should be defined");
assert(pkg.delete_drive != null, "delete_drive function should be defined");

assert(pkg.list_snapshots != null, "list_snapshots function should be defined");
assert(pkg.get_snapshot != null, "get_snapshot function should be defined");
assert(pkg.delete_snapshot != null, "delete_snapshot function should be defined");
assert(pkg.create_snapshot != null, "create_snapshot function should be defined");
    `)
}

func TestThrows(t *testing.T) {
	cloud := mockcloud.Client(nil)

	cloud.MockDrives.ListDrivesFn = func(_ context.Context) (<-chan drives.Drive, <-chan error) {
		lc := make(chan drives.Drive)
		close(lc)
		ec := make(chan error, 1)
		ec <- errors.New("throw me")
		close(ec)
		return lc, ec
	}
	cloud.MockDrives.GetDriveFn = func(_ context.Context, _ string) (drives.Drive, error) {
		return nil, errors.New("throw me")
	}
	cloud.MockDrives.CreateDriveFn = func(_ context.Context, _, _ string, _ int64, _ ...drives.CreateOpt) (drives.Drive, error) {
		return nil, errors.New("throw me")
	}
	cloud.MockDrives.DeleteDriveFn = func(_ context.Context, _ string) error {
		return errors.New("throw me")
	}
	cloud.MockDrives.ListSnapshotsFn = func(_ context.Context, _ string) (<-chan drives.Snapshot, <-chan error) {
		lc := make(chan drives.Snapshot)
		close(lc)
		ec := make(chan error, 1)
		ec <- errors.New("throw me")
		close(ec)
		return lc, ec
	}
	cloud.MockDrives.GetSnapshotFn = func(_ context.Context, _ string) (drives.Snapshot, error) {
		return nil, errors.New("throw me")
	}
	cloud.MockDrives.CreateSnapshotFn = func(_ context.Context, _, _ string, _ ...drives.SnapshotOpt) (drives.Snapshot, error) {
		return nil, errors.New("throw me")
	}
	cloud.MockDrives.DeleteSnapshotFn = func(_ context.Context, _ string) error {
		return errors.New("throw me")
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.drives;

var dr = {
	type: "",
	name: "",
	data: "",
	priority: "",
	port: "",
	weight: "",
};

[
	{ name: "list_drives",         fn: function() { pkg.list_drives() } },
	{ name: "get_drive",           fn: function() { pkg.get_drive("hello.com") } },
	{ name: "create_drive",        fn: function() { pkg.create_drive({}) } },
	{ name: "delete_drive",        fn: function() { pkg.delete_drive({}) } },

	{ name: "list_snapshots",      fn: function() { pkg.list_snapshots("hello.com") } },
	{ name: "get_snapshot",        fn: function() { pkg.get_snapshot("hello.com", 1) } },
	{ name: "create_snapshot",     fn: function() { pkg.create_snapshot("hello.com", {}) } },
	{ name: "delete_snapshot",     fn: function() { pkg.delete_snapshot("hello.com", 1)  } }

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

type drive struct {
	*godo.Drive
}

func (k *drive) Struct() *godo.Drive { return k.Drive }

type snapshot struct {
	*godo.Snapshot
}

func (k *snapshot) Struct() *godo.Snapshot { return k.Snapshot }

func TestList(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockDrives.ListDrivesFn = func(_ context.Context) (<-chan drives.Drive, <-chan error) {
		lc := make(chan drives.Drive, 1)
		lc <- &drive{&godo.Drive{
			ID:            "lol",
			Region:        &godo.Region{Slug: "nyc1"},
			Name:          "my_name",
			SizeGigaBytes: 100,
			Description:   "lolz",
			DropletIDs:    []int{42},
		}}
		close(lc)
		ec := make(chan error)
		close(ec)
		return lc, ec
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.drives;

var list = pkg.list_drives();
assert(list != null, "should have received a list");
assert(list.length > 0, "should have received some elements")

var d = list[0];
var want = {
	id:          "lol",
	region:      "nyc1",
	name:        "my_name",
	size:        100,
	description: "lolz",
	droplet_ids: [42],
};
equals(d, want, "should have proper object");
`)
}

func TestGet(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockDrives.GetDriveFn = func(_ context.Context, _ string) (drives.Drive, error) {
		return &drive{&godo.Drive{
			ID:            "lol",
			Region:        &godo.Region{Slug: "nyc1"},
			Name:          "my_name",
			SizeGigaBytes: 100,
			Description:   "lolz",
			DropletIDs:    []int{42},
		}}, nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.drives;

var d = pkg.get_drive("my_name")
var want = {
	id:          "lol",
	region:      "nyc1",
	name:        "my_name",
	size:        100,
	description: "lolz",
	droplet_ids: [42],
};
equals(d, want, "should have proper object");
`)
}

func TestCreate(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockDrives.CreateDriveFn = func(_ context.Context, _, _ string, _ int64, _ ...drives.CreateOpt) (drives.Drive, error) {
		return &drive{&godo.Drive{
			ID:            "lol",
			Region:        &godo.Region{Slug: "nyc1"},
			Name:          "my_name",
			SizeGigaBytes: 100,
			Description:   "lolz",
			DropletIDs:    []int{42},
		}}, nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.drives;

var d = pkg.create_drive({
	name: "my_name",
	ip: "127.0.0.1"
});
var want = {
	id:          "lol",
	region:      "nyc1",
	name:        "my_name",
	size:        100,
	description: "lolz",
	droplet_ids: [42],
};
equals(d, want, "should have proper object");
`)
}

func TestDelete(t *testing.T) {
	wantName := "my name"
	cloud := mockcloud.Client(nil)
	cloud.MockDrives.DeleteDriveFn = func(_ context.Context, gotName string) error {
		if gotName != wantName {
			t.Fatalf("want %q got %q", wantName, gotName)
		}
		return nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.drives;

pkg.delete_drive("my name");
`)
}

func TestListSnapshot(t *testing.T) {
	wantName := "my name"
	cloud := mockcloud.Client(nil)
	cloud.MockDrives.ListSnapshotsFn = func(_ context.Context, gotName string) (<-chan drives.Snapshot, <-chan error) {
		if gotName != wantName {
			t.Fatalf("want %q got %q", wantName, gotName)
		}
		lc := make(chan drives.Snapshot, 1)
		lc <- &snapshot{&godo.Snapshot{
			ID:            "lol",
			DriveID:       "lolzzzz",
			Region:        &godo.Region{Slug: "nyc1"},
			Name:          "my_name",
			SizeGibiBytes: 100,
			Description:   "lolz",
		}}
		close(lc)
		ec := make(chan error)
		close(ec)
		return lc, ec
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.drives;

var snapshots = pkg.list_snapshots("my name");
assert(snapshots != null, "should have received a snapshots");
assert(snapshots.length > 0, "should have received some elements")

var d = snapshots[0];
var want = {
	id:          "lol",
	drive_id:    "lolzzzz",
	region:      "nyc1",
	name:        "my_name",
	size:        100,
	description: "lolz",
};
equals(d, want, "should have proper object");
`)
}

func TestGetSnapshot(t *testing.T) {
	wantName := "my_name"
	cloud := mockcloud.Client(nil)
	cloud.MockDrives.GetSnapshotFn = func(_ context.Context, gotName string) (drives.Snapshot, error) {
		if gotName != wantName {
			t.Fatalf("want %q got %q", wantName, gotName)
		}
		return &snapshot{&godo.Snapshot{
			ID:            "lol",
			DriveID:       "lolzzzz",
			Region:        &godo.Region{Slug: "nyc1"},
			Name:          wantName,
			SizeGibiBytes: 100,
			Description:   "lolz",
		}}, nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.drives;

var d = pkg.get_snapshot("my_name", 42)
var want = {
	id:          "lol",
	drive_id:    "lolzzzz",
	region:      "nyc1",
	name:        "my_name",
	size:        100,
	description: "lolz",
};
equals(d, want, "should have proper object");
`)
}

func TestCreateSnapshot(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockDrives.CreateSnapshotFn = func(_ context.Context, _, _ string, _ ...drives.SnapshotOpt) (drives.Snapshot, error) {
		return &snapshot{&godo.Snapshot{
			ID:            "lol",
			DriveID:       "lolzzzz",
			Region:        &godo.Region{Slug: "nyc1"},
			Name:          "my_name",
			SizeGibiBytes: 100,
			Description:   "lolz",
		}}, nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.drives;

var want = {
	id:          "lol",
	drive_id:    "lolzzzz",
	region:      "nyc1",
	name:        "my_name",
	size:        100,
	description: "lolz",
};

var d = pkg.create_snapshot(want.name, want);
equals(d, want, "should have proper object");
`)
}

func TestDeleteSnapshot(t *testing.T) {

	wantID := "my id"

	cloud := mockcloud.Client(nil)
	cloud.MockDrives.DeleteSnapshotFn = func(_ context.Context, gotID string) error {

		if gotID != wantID {
			t.Fatalf("want %v got %v", wantID, gotID)
		}
		return nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.drives;

pkg.delete_snapshot("my id");
`)
}
