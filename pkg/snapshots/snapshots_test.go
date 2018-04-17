package snapshots_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aybabtme/godotto/pkg/extra/vmtest"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/snapshots"
	"github.com/aybabtme/godotto/pkg/extra/do/mockcloud"
	"github.com/digitalocean/godo"
)

var testSnap *godo.Snapshot = &godo.Snapshot{}

func TestSnapshotApply(t *testing.T) {
	cloud := mockcloud.Client(nil)

	vmtest.Run(t, cloud, `
	var pkg = cloud.snapshots;
	assert(pkg != null, "package should be loaded");
	assert(pkg.list != null, "list function shouled be defined");
	assert(pkg.list_droplet != null, "list_droplet function shouled be defined");
	assert(pkg.list_volume != null, "list_volume function shouled be defined");
	assert(pkg.get != null, "get function shouled be defined");
	assert(pkg.delete != null, "delete function should be defined");
	`)
}

func TestSnapshotThrows(t *testing.T) {
	cloud := mockcloud.Client(nil)

	cloud.MockSnapshots.ListFn = func(_ context.Context) (<-chan snapshots.Snapshot, <-chan error) {
		sc := make(chan snapshots.Snapshot)
		close(sc)
		ec := make(chan error, 1)
		ec <- errors.New("throw me")
		close(ec)
		return sc, ec
	}

	cloud.MockSnapshots.ListDropletFn = func(_ context.Context) (<-chan snapshots.Snapshot, <-chan error) {
		sc := make(chan snapshots.Snapshot)
		close(sc)
		ec := make(chan error, 1)
		ec <- errors.New("throw me")
		close(ec)
		return sc, ec
	}

	cloud.MockSnapshots.ListVolumeFn = func(_ context.Context) (<-chan snapshots.Snapshot, <-chan error) {
		sc := make(chan snapshots.Snapshot)
		close(sc)
		ec := make(chan error, 1)
		ec <- errors.New("throw me")
		close(ec)
		return sc, ec
	}

	cloud.MockSnapshots.GetFn = func(_ context.Context, _ string) (snapshots.Snapshot, error) {
		return nil, errors.New("throw me")
	}

	cloud.MockSnapshots.DeleteFn = func(_ context.Context, _ string) error {
		return errors.New("throw me")
	}

	vmtest.Run(t, cloud, `
		var pkg = cloud.snapshots;


		var ss = {
		    "created_at": "2017-06-08T09:11:06Z",
		    "id": "11223344",
		    "min_disk_size": 20,
		    "name": "example-server-007",
		    "regions": [
		      "nyc1"
		    ],
		    "resource_id": "44332211",
		    "resource_type": "droplet",
		    "size": 2.24
		};

		[
			{name: "get", fn: function() { pkg.get(ss.id) }},
			{name: "delete", fn: function() { pkg.delete(ss.id) }},
			{name: "list", fn: function() { pkg.list() }},
			{name: "list_droplet", fn: function() { pkg.list_droplet() }},
			{name: "list_volume", fn: function() { pkg.list_volume() }},
		 ].forEach(function(kv) {
			var name = kv.name;
			var fn = kv.fn;

			try {
				fn(); throw "don't catch me";
			} catch(e) {
				equals("throw me", e.message, name + "should send the right exception!");
			}
		 });
	`)
}

var (
	sd = &godo.Snapshot{ID: "11223344", Name: "example-server-007", ResourceID: "44332211", ResourceType: "droplet", Regions: []string{"nyc1"}, MinDiskSize: 20, SizeGigaBytes: 2.24, Created: "2017-06-08T09:11:06Z"}
	sv = &godo.Snapshot{ID: "11223345", Name: "example-server-007", ResourceID: "44332210", ResourceType: "volume", Regions: []string{"nyc1"}, MinDiskSize: 20, SizeGigaBytes: 2.24, Created: "2017-06-08T09:11:06Z"}
)

type snapshot struct {
	*godo.Snapshot
}

func (k *snapshot) Struct() *godo.Snapshot { return k.Snapshot }

func TestSnapshotsList(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockSnapshots.ListFn = func(_ context.Context) (<-chan snapshots.Snapshot, <-chan error) {
		sc := make(chan snapshots.Snapshot, 1)
		sc <- &snapshot{sd}
		close(sc)
		ec := make(chan error)
		close(ec)
		return sc, ec
	}

	vmtest.Run(t, cloud, `
	var pkg = cloud.snapshots;
	var list = pkg.list();
	assert(list != null, "should have received a list");
	assert(list.length > 0, "should have received some elements");

	var want = {
		    "created_at": "2017-06-08T09:11:06Z",
		    "id": "11223344",
		    "min_disk_size": 20,
		    "name": "example-server-007",
		    "regions": [
		      "nyc1"
		    ],
		    "resource_id": "44332211",
		    "resource_type": "droplet",
		    "size": 2.24
	};


	var s = list[0];

	equals(s, want, "should have proper object");
	`)
}

func TestSnapshotsListDroplet(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockSnapshots.ListDropletFn = func(_ context.Context) (<-chan snapshots.Snapshot, <-chan error) {
		sc := make(chan snapshots.Snapshot, 1)
		sc <- &snapshot{sd}
		close(sc)
		ec := make(chan error)
		close(ec)
		return sc, ec
	}

	vmtest.Run(t, cloud, `
	var pkg = cloud.snapshots;
	var list_droplet = pkg.list_droplet();
	assert(list_droplet != null, "should have received a list");
	assert(list_droplet.length > 0, "should have received some elements");

	var want = {
		    "created_at": "2017-06-08T09:11:06Z",
		    "id": "11223344",
		    "min_disk_size": 20,
		    "name": "example-server-007",
		    "regions": [
		      "nyc1"
		    ],
		    "resource_id": "44332211",
		    "resource_type": "droplet",
		    "size": 2.24
	};


	var s = list_droplet[0];

	equals(s, want, "should have proper object");
	`)
}

func TestSnapshotsListVolume(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockSnapshots.ListVolumeFn = func(_ context.Context) (<-chan snapshots.Snapshot, <-chan error) {
		sc := make(chan snapshots.Snapshot, 1)
		sc <- &snapshot{sv}
		close(sc)
		ec := make(chan error)
		close(ec)
		return sc, ec
	}

	vmtest.Run(t, cloud, `
	var pkg = cloud.snapshots;
	var list_volume = pkg.list_volume();
	assert(list_volume != null, "should have received a list");
	assert(list_volume.length > 0, "should have received some elements");

	var want = {
		    "created_at": "2017-06-08T09:11:06Z",
		    "id": "11223345",
		    "min_disk_size": 20,
		    "name": "example-server-007",
		    "regions": [
		      "nyc1"
		    ],
		    "resource_id": "44332210",
		    "resource_type": "volume",
		    "size": 2.24
	};


	var s = list_volume[0];

	equals(s, want, "should have proper object");
	`)
}
func TestSnapshotDelete(t *testing.T) {
	wantId := "11223344"
	cloud := mockcloud.Client(nil)

	cloud.MockSnapshots.DeleteFn = func(_ context.Context, gotId string) error {
		if gotId != wantId {
			t.Fatalf("want %v got %v", wantId, gotId)
		}

		return nil
	}

	vmtest.Run(t, cloud, `
			var pkg = cloud.snapshots;
			pkg.delete("11223344");
	`)
}

func TestSnapshotGet(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockSnapshots.GetFn = func(_ context.Context, id string) (snapshots.Snapshot, error) {
		return &snapshot{sd}, nil
	}

	vmtest.Run(t, cloud, `
		var pkg = cloud.snapshots;

		var want = {
		    "created_at": "2017-06-08T09:11:06Z",
		    "id": "11223344",
		    "min_disk_size": 20,
		    "name": "example-server-007",
		    "regions": [
		      "nyc1"
		    ],
		    "resource_id": "44332211",
		    "resource_type": "droplet",
		    "size": 2.24
		};


		var l = pkg.get('11223344');

		equals(l, want, "should have proper object");
	`)
}
