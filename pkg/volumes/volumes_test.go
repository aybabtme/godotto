package volumes_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aybabtme/godotto/pkg/extra/vmtest"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/volumes"
	"github.com/aybabtme/godotto/pkg/extra/do/mockcloud"
	"github.com/digitalocean/godo"
)

func TestApply(t *testing.T) {
	cloud := mockcloud.Client(nil)
	vmtest.Run(t, cloud, `
var pkg = cloud.volumes;

assert(pkg != null, "package should be loaded");
assert(pkg.list_volumes != null, "list_volumes function should be defined");
assert(pkg.get_volume != null, "get_volume function should be defined");
assert(pkg.create_volume != null, "create_volume function should be defined");
assert(pkg.delete_volume != null, "delete_volume function should be defined");

assert(pkg.list_snapshots != null, "list_snapshots function should be defined");
assert(pkg.get_snapshot != null, "get_snapshot function should be defined");
assert(pkg.delete_snapshot != null, "delete_snapshot function should be defined");
assert(pkg.create_snapshot != null, "create_snapshot function should be defined");
    `)
}

func TestThrows(t *testing.T) {
	cloud := mockcloud.Client(nil)

	cloud.MockVolumes.ListVolumesFn = func(_ context.Context) (<-chan volumes.Volume, <-chan error) {
		lc := make(chan volumes.Volume)
		close(lc)
		ec := make(chan error, 1)
		ec <- errors.New("throw me")
		close(ec)
		return lc, ec
	}
	cloud.MockVolumes.GetVolumeFn = func(_ context.Context, _ string) (volumes.Volume, error) {
		return nil, errors.New("throw me")
	}
	cloud.MockVolumes.CreateVolumeFn = func(_ context.Context, _, _ string, _ int64, _ ...volumes.CreateOpt) (volumes.Volume, error) {
		return nil, errors.New("throw me")
	}
	cloud.MockVolumes.DeleteVolumeFn = func(_ context.Context, _ string) error {
		return errors.New("throw me")
	}
	cloud.MockVolumes.ListSnapshotsFn = func(_ context.Context, _ string) (<-chan volumes.Snapshot, <-chan error) {
		lc := make(chan volumes.Snapshot)
		close(lc)
		ec := make(chan error, 1)
		ec <- errors.New("throw me")
		close(ec)
		return lc, ec
	}
	cloud.MockVolumes.GetSnapshotFn = func(_ context.Context, _ string) (volumes.Snapshot, error) {
		return nil, errors.New("throw me")
	}
	cloud.MockVolumes.CreateSnapshotFn = func(_ context.Context, _, _ string, _ ...volumes.SnapshotOpt) (volumes.Snapshot, error) {
		return nil, errors.New("throw me")
	}
	cloud.MockVolumes.DeleteSnapshotFn = func(_ context.Context, _ string) error {
		return errors.New("throw me")
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.volumes;

var r = { name: "", slug: "", sizes: [], available: true, features: [] };
var dr = {
	id: "",
	region: "",
	name: "",
	size: "",
	desc: "",
	droplet_ids: []
};
var snap = {
	volume: "",
	name: "",
	desc: ""
};

[
	{ name: "list_volumes",         fn: function() { pkg.list_volumes() } },
	{ name: "get_volume",           fn: function() { pkg.get_volume("hello.com") } },
	{ name: "create_volume",        fn: function() { pkg.create_volume(dr) } },
	{ name: "delete_volume",        fn: function() { pkg.delete_volume("") } },

	{ name: "list_snapshots",      fn: function() { pkg.list_snapshots("") } },
	{ name: "get_snapshot",        fn: function() { pkg.get_snapshot("") } },
	{ name: "create_snapshot",     fn: function() { pkg.create_snapshot(snap) } },
	{ name: "delete_snapshot",     fn: function() { pkg.delete_snapshot("")  } }

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

type volume struct {
	*godo.Volume
}

func (k *volume) Struct() *godo.Volume { return k.Volume }

type snapshot struct {
	*godo.Snapshot
}

func (k *snapshot) Struct() *godo.Snapshot { return k.Snapshot }

var (
	region = &godo.Region{Name: "newyork3", Slug: "nyc3", Sizes: []string{"small"}, Available: true, Features: []string{"all"}}
)

func TestVolumeList(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockVolumes.ListVolumesFn = func(_ context.Context) (<-chan volumes.Volume, <-chan error) {
		lc := make(chan volumes.Volume, 1)
		lc <- &volume{&godo.Volume{
			ID:            "lol",
			Region:        region,
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
var pkg = cloud.volumes;

var list = pkg.list_volumes();
assert(list != null, "should have received a list");
assert(list.length > 0, "should have received some elements")

var region = { name: "newyork3", slug: "nyc3", sizes: ["small"], available: true, features: ["all"] };

var d = list[0];
var want = {
	id:          "lol",
	region:      region,
	name:        "my_name",
	size:        100,
	description: "lolz",
	droplet_ids: [42],
};
equals(d, want, "should have proper object");
`)
}

func TestVolumeGet(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockVolumes.GetVolumeFn = func(_ context.Context, _ string) (volumes.Volume, error) {
		return &volume{&godo.Volume{
			ID:            "lol",
			Region:        region,
			Name:          "my_name",
			SizeGigaBytes: 100,
			Description:   "lolz",
			DropletIDs:    []int{42},
		}}, nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.volumes;

var region = { name: "newyork3", slug: "nyc3", sizes: ["small"], available: true, features: ["all"] };

var d = pkg.get_volume("my_name")
var want = {
	id:          "lol",
	region:      region,
	name:        "my_name",
	size:        100,
	description: "lolz",
	droplet_ids: [42],
};
equals(d, want, "should have proper object");
`)
}

func TestVolumeCreate(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockVolumes.CreateVolumeFn = func(_ context.Context, _, _ string, _ int64, _ ...volumes.CreateOpt) (volumes.Volume, error) {
		return &volume{&godo.Volume{
			ID:            "lol",
			Region:        region,
			Name:          "my_name",
			SizeGigaBytes: 100,
			Description:   "lolz",
			DropletIDs:    []int{42},
		}}, nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.volumes;

var region = { name: "newyork3", slug: "nyc3", sizes: ["small"], available: true, features: ["all"] };

var d = pkg.create_volume({
	name: "my_name",
	size: 100,
	region: region
});
var want = {
	id:          "lol",
	region:      region,
	name:        "my_name",
	size:        100,
	description: "lolz",
	droplet_ids: [42],
};
equals(d, want, "should have proper object");
`)
}

func TestVolumeDelete(t *testing.T) {
	wantName := "my name"
	cloud := mockcloud.Client(nil)
	cloud.MockVolumes.DeleteVolumeFn = func(_ context.Context, gotName string) error {
		if gotName != wantName {
			t.Fatalf("want %q got %q", wantName, gotName)
		}
		return nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.volumes;

pkg.delete_volume("my name");
`)
}

func TestListSnapshot(t *testing.T) {
	wantName := "my name"
	cloud := mockcloud.Client(nil)
	cloud.MockVolumes.ListSnapshotsFn = func(_ context.Context, gotName string) (<-chan volumes.Snapshot, <-chan error) {
		if gotName != wantName {
			t.Fatalf("want %q got %q", wantName, gotName)
		}
		lc := make(chan volumes.Snapshot, 1)
		lc <- &snapshot{&godo.Snapshot{
			ID:            "lol",
			ResourceID:    "lolzzzz",
			Regions:       []string{region.Slug},
			Name:          "my_name",
			SizeGigaBytes: 100,
		}}
		close(lc)
		ec := make(chan error)
		close(ec)
		return lc, ec
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.volumes;

var snapshots = pkg.list_snapshots("my name");
assert(snapshots != null, "should have received a snapshots");
assert(snapshots.length > 0, "should have received some elements")

var d = snapshots[0];
var want = {
	id:          "lol",
	volume_id:    "lolzzzz",
	regions:     ["nyc3"],
	name:        "my_name",
	size:        100,
};
equals(d, want, "should have proper object");
`)
}

func TestGetSnapshot(t *testing.T) {
	wantName := "my_name"
	cloud := mockcloud.Client(nil)
	cloud.MockVolumes.GetSnapshotFn = func(_ context.Context, gotName string) (volumes.Snapshot, error) {
		if gotName != wantName {
			t.Fatalf("want %q got %q", wantName, gotName)
		}
		return &snapshot{&godo.Snapshot{
			ID:            "lol",
			ResourceID:    "lolzzzz",
			Regions:       []string{region.Slug},
			Name:          wantName,
			SizeGigaBytes: 100,
		}}, nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.volumes;

var region = { name: "newyork3", slug: "nyc3", sizes: ["small"], available: true, features: ["all"] };

var d = pkg.get_snapshot("my_name", 42)
var want = {
	id:          "lol",
	volume_id:    "lolzzzz",
	regions:      ["nyc3"],
	name:        "my_name",
	size:        100
};
equals(d, want, "should have proper object");
`)
}

func TestCreateSnapshot(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockVolumes.CreateSnapshotFn = func(_ context.Context, _, _ string, _ ...volumes.SnapshotOpt) (volumes.Snapshot, error) {
		return &snapshot{&godo.Snapshot{
			ID:            "lol",
			ResourceID:    "lolzzzz",
			Regions:       []string{region.Slug},
			Name:          "my_name",
			SizeGigaBytes: 100,
		}}, nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.volumes;

var region = { name: "newyork3", slug: "nyc3", sizes: ["small"], available: true, features: ["all"] };

var arg = {
	id:          "lol",
	volume:       "lolzzzz",
	region:      region,
	name:        "my_name",
	description: "lolz",
};
var want = {
	id:          "lol",
	volume_id:	 "lolzzzz",
	regions:      ["nyc3"],
	name:        "my_name",
	size:        100,
};

var d = pkg.create_snapshot(arg);
equals(d, want, "should have proper object");
`)
}

func TestDeleteSnapshot(t *testing.T) {

	wantID := "my id"

	cloud := mockcloud.Client(nil)
	cloud.MockVolumes.DeleteSnapshotFn = func(_ context.Context, gotID string) error {

		if gotID != wantID {
			t.Fatalf("want %v got %v", wantID, gotID)
		}
		return nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.volumes;

pkg.delete_snapshot("my id");
`)
}
